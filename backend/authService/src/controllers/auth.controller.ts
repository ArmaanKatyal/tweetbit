import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { SECRET_KEY, TokenPayload } from '../middlewares/auth.middleware';
import { validateLogin, validateRegister } from '../validation/auth.validate';
import { Auth } from '../models/auth.model';
import nodeConfig from 'config';
import * as bcrypt from 'bcrypt';
import { v4 as uuidv4 } from 'uuid';
import * as dotenv from 'dotenv';
import prisma from '../../prisma/client';
import {
    IncHttpTransaction,
    MetricsCode,
    MetricsMethod,
    ObserveHttpResponseTime,
} from '../internal/prometheus';
dotenv.config();

export const salt: number = parseInt(process.env.SALT_ROUNDS!);

const login = async (req: Request, res: Response) => {
    let start = Date.now();

    // Validate input
    const { error } = validateLogin(req.body);
    if (error) {
        req.log.info({
            message: 'Invalid input',
            userEmail: req.body.email,
            service: 'auth',
            function: 'login',
            error: error.details[0].message,
        });
        collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
        return res.status(400).json({ error: nodeConfig.get('error_codes.INVALID_INPUT') });
    }
    // Check if user exists
    let checkAuth = await Auth.findOne({ email: req.body.email });
    if (!checkAuth) {
        req.log.info({
            message: 'Auth not found',
            userEmail: req.body.email,
            service: 'auth',
            function: 'login',
        });
        collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
        return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
    }
    // Check if password is correct
    const validPassword = await bcrypt.compare(req.body.password, checkAuth.password);
    if (!validPassword) {
        req.log.info({
            message: 'Invalid password',
            userEmail: req.body.email,
            service: 'auth',
            function: 'login',
        });
        collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
        return res.status(400).json({ error: nodeConfig.get('error_codes.INVALID_PASSWORD') });
    }
    // Create and assign a token
    const access_token = jwt.sign(
        {
            email: checkAuth.email,
            uuid: checkAuth.uuid,
            type: 'access',
        },
        SECRET_KEY,
        {
            expiresIn: nodeConfig.get('token.expire.access'),
        }
    );

    const refresh_token = jwt.sign(
        {
            email: checkAuth.email,
            uuid: checkAuth.uuid,
            type: 'refresh',
        },
        SECRET_KEY,
        {
            expiresIn: nodeConfig.get('token.expire.refresh'),
        }
    );

    res.cookie('access_token', access_token, {
        httpOnly: true,
        secure: false,
        sameSite: 'none',
    });

    // Set cookie
    res.cookie('refresh_token', refresh_token, {
        httpOnly: true,
        secure: false,
        sameSite: 'none',
    });

    let checkUser;
    try {
        checkUser = await prisma.user.findUnique({
            where: {
                email: checkAuth.email,
            },
        });
        if (!checkUser) {
            req.log.info({
                message: 'User not found',
                userEmail: req.body.email,
                service: 'auth',
                function: 'login',
            });
            collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }
    } catch (err) {
        req.log.info({
            message: 'Error while fetching user from user database',
            userEmail: req.body.email,
            service: 'auth',
            function: 'login',
        });
        collectMetrics(MetricsCode.InternalServerError, MetricsMethod.Post, start);
        return res.status(500).json({
            error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'),
        });
    }

    // Remove the uuid and _id from the payload
    let { uuid, id, ...newPayload } = checkUser as any;

    // Send token
    res.status(200).json({
        access_token,
        refresh_token,
        user: newPayload,
    });
    collectMetrics(MetricsCode.Ok, MetricsMethod.Post, start);
};

const logout = async (_: Request, res: Response) => {
    let start = Date.now();
    res.clearCookie('refresh_token');
    res.clearCookie('access_token');
    res.status(200).json({
        message: 'Logged out successfully',
    });
    collectMetrics(MetricsCode.Ok, MetricsMethod.Post, start);
};

const refresh = async (req: Request, res: Response) => {
    let start = Date.now();

    // get the token from the cookies or the header of the request
    let token = req.cookies.refresh_token || req.headers['x-refresh-token'];
    if (!token) {
        req.log.info({
            message: 'No refresh token found',
            service: 'auth',
            function: 'refresh',
        });
        collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
        return res.status(400).json({ error: 'no_refresh_token_provided' });
    }

    // verify the token
    try {
        let decoded = jwt.verify(token, SECRET_KEY) as TokenPayload;
        if (decoded.type !== 'refresh') {
            req.log.info({
                message: 'Invalid refresh token',
                service: 'auth',
                function: 'refresh',
            });
            collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
            return res
                .status(400)
                .json({ error: 'invalid_refresh_token' });
        }

        // If the token is valid, create a new access token and send it
        let access_token = jwt.sign(
            {
                email: decoded.email,
                uuid: decoded.uuid,
                type: 'access',
            },
            SECRET_KEY,
            {
                expiresIn: nodeConfig.get('token.expire.access'),
            }
        );

        collectMetrics(MetricsCode.Ok, MetricsMethod.Post, start);
        return res.status(200).json({ access_token });
    } catch (error: any) {
        req.log.error({
            message: 'Error while verifying refresh token',
            service: 'auth',
            function: 'refresh',
        });
        if (error instanceof jwt.TokenExpiredError) {
            return res.status(401).json({
                error: 'token_expired',
            });
        } else if (error instanceof jwt.JsonWebTokenError) {
            return res.status(401).json({
                error: 'invalid_token',
            });
        } else {
            return res.status(500).json({
                error: 'internal_server_error',
            });
        }
    }
};

const register = async (req: Request, res: Response) => {
    let start = Date.now();
    // validate the input data from the body
    const { error } = validateRegister(req.body);
    if (error) {
        req.log.info({
            message: 'Invalid input',
            userEmail: req.body.email,
            service: 'auth',
            function: 'register',
            error: error.details[0].message,
        });
        collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
        return res.status(400).json({ error: nodeConfig.get('error_codes.INVALID_INPUT') });
    }

    // Check if user already exists
    // let checkUser = await Auth.findOne({ email: req.body.email });
    let checkUser = await prisma.user.findUnique({
        where: {
            email: req.body.email,
        },
    });
    if (checkUser) {
        req.log.info({
            message: 'User already exists',
            userEmail: req.body.email,
            service: 'auth',
            function: 'register',
        });
        collectMetrics(MetricsCode.BadRequest, MetricsMethod.Post, start);
        return res.status(400).json({ error: nodeConfig.get('error_codes.USER_ALREADY_EXISTS') });
    }

    // Create a unique ID for the user
    const uniqueID = uuidv4();
    // Save the user to auth database
    const newAuth = new Auth({
        uuid: uniqueID,
        email: req.body.email,
        password: bcrypt.hashSync(req.body.password, salt),
    });

    try {
        await newAuth.save();
    } catch (err) {
        req.log.info({
            message: 'Error while saving user to auth database',
            userEmail: req.body.email,
            service: 'auth',
            function: 'register',
        });
        collectMetrics(MetricsCode.InternalServerError, MetricsMethod.Post, start);
        return res.status(500).json({
            error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'),
        });
    }

    // Remove password from the payload
    // let { password, ...newPayload } = req.body;

    try {
        await prisma.user.create({
            data: {
                uuid: uniqueID,
                first_name: req.body.first_name,
                last_name: req.body.last_name,
                email: req.body.email,
            },
        });
    } catch (err) {
        req.log.info({
            message: 'Error while saving user to user database',
            userEmail: req.body.email,
            service: 'auth',
            function: 'register',
        });
        collectMetrics(MetricsCode.InternalServerError, MetricsMethod.Post, start);
        return res.status(500).json({
            error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'),
        });
    }

    res.status(200).json({
        message: 'User registered successfully',
    });
    collectMetrics(MetricsCode.Ok, MetricsMethod.Post, start);
};

const collectMetrics = (code: MetricsCode, method: MetricsMethod, time: number) => {
    IncHttpTransaction(code, method);
    ObserveHttpResponseTime(code, method, Date.now() - time);
};

export { login, logout, refresh, register };
