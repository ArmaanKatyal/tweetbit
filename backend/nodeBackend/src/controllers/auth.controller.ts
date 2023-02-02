import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { SECRET_KEY } from '../middlewares/auth.middleware';
import { validateLogin } from '../validation/auth.validate';
import { Auth } from '../models/auth.model';
import nodeConfig from 'config';
import * as bcrpyt from 'bcrypt';

const login = async (req: Request, res: Response) => {
    // Validate input
    const { error } = validateLogin(req.body);
    if (error) {
        return res
            .status(400)
            .json({ error: nodeConfig.get('error_codes.INVALID_INPUT') });
    }
    // Check if user exists
    const checkUser = await Auth.findOne({ email: req.body.email });
    if (!checkUser) {
        return res
            .status(400)
            .json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
    }
    // Check if password is correct
    const validPassword = await bcrpyt.compare(
        req.body.password,
        checkUser.password
    );
    if (!validPassword) {
        return res
            .status(400)
            .json({ error: nodeConfig.get('error_codes.INVALID_PASSWORD') });
    }
    // Create and assign a token
    const access_token = jwt.sign(
        {
            id: checkUser._id,
            email: checkUser.email,
            type: 'access',
        },
        SECRET_KEY,
        {
            expiresIn: nodeConfig.get('token.expire.access'),
        }
    );

    const refresh_token = jwt.sign(
        {
            id: checkUser._id,
            email: checkUser.email,
            type: 'refresh',
        },
        SECRET_KEY,
        {
            expiresIn: nodeConfig.get('token.expire.refresh'),
        }
    );

    // Set cookie
    res.cookie('refresh_token', refresh_token, {
        httpOnly: true,
        secure: false,
        sameSite: 'none',
    });

    // TODO: Make a call to the User Database and attach the information with the response

    // Send token
    res.status(200).json({
        access_token,
        refresh_token,
    });
};

const logout = async (req: Request, res: Response) => {
    res.clearCookie('refresh_token');
    res.status(200).json({
        message: 'Logged out successfully',
    });
};

const refresh = async (req: Request, res: Response) => {
    const access_token = jwt.sign(
        {
            id: req.body.id,
            email: req.body.email,
        },
        SECRET_KEY,
        {
            expiresIn: nodeConfig.get('token.expire.access'),
        }
    );

    res.status(200).json({
        access_token,
    });
};

const temp = async (req: Request, res: Response) => {
    res.status(200).json({
        message: 'You are logged in',
    });
};

export { login, logout, refresh, temp };
