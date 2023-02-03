import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { SECRET_KEY } from '../middlewares/auth.middleware';
import { validateLogin, validateRegister } from '../validation/auth.validate';
import { Auth } from '../models/auth.model';
import { User } from '../models/user.model';
import nodeConfig from 'config';
import * as bcrypt from 'bcrypt';
import { v4 as uuidv4 } from 'uuid';
import * as dotenv from 'dotenv';
dotenv.config();

export const SALT_ROUNDS: string = process.env.SALT_ROUNDS!;

const login = async (req: Request, res: Response) => {
    // Validate input
    const { error } = validateLogin(req.body);
    if (error) {
        return res
            .status(400)
            .json({ error: nodeConfig.get('error_codes.INVALID_INPUT') });
    }
    // Check if user exists
    let checkAuth = await Auth.findOne({ email: req.body.email });
    if (!checkAuth) {
        return res
            .status(400)
            .json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
    }
    // Check if password is correct
    const validPassword = await bcrypt.compare(
        req.body.password,
        checkAuth.password
    );
    if (!validPassword) {
        return res
            .status(400)
            .json({ error: nodeConfig.get('error_codes.INVALID_PASSWORD') });
    }
    // Create and assign a token
    const access_token = jwt.sign(
        {
            uuid: checkAuth.uuid,
            email: checkAuth.email,
            type: 'access',
        },
        SECRET_KEY,
        {
            expiresIn: nodeConfig.get('token.expire.access'),
        }
    );

    const refresh_token = jwt.sign(
        {
            uuid: checkAuth.uuid,
            email: checkAuth.email,
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

    // Make a call to the User Database and attach the information with the response
    let checkUser = await User.findOne({ uuid: checkAuth.uuid, email: checkAuth.email });
    if (!checkUser) {
        return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
    }
    
    // Remove the uuid and _id from the payload
    let { uuid, _id, ...newPayload } = checkUser;

    // Send token
    res.status(200).json({
        access_token,
        refresh_token,
        user: newPayload,
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

const register = async (req: Request, res: Response) => {
    // validate the input data from the body
    const { error } = validateRegister(req.body);
    if (error) {
        return res.status(400).json({ error: nodeConfig.get('error_codes.INVALID_INPUT') });
    }

    // Check if user already exists
    let checkUser = await User.findOne({ email: req.body.email });
    if (checkUser) {
        return res.status(400).json({ error: nodeConfig.get('error_codes.USER_ALREADY_EXISTS') });
    }

    // Create a unique ID for the user
    const uniqueID = uuidv4();
    // Save the user to auth database
    const newAuth = new Auth({
        uuid: uniqueID,
        email: req.body.email,
        password: await bcrypt.hashSync(req.body.password, SALT_ROUNDS)
    });

    try {
        await newAuth.save();
        res.status(200).json({ message: 'Auth created successfully' });
    } catch (err) {
        res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }

    // Remove password from the payload
    let { password, ...newPayload } = req.body;
    
    // Save the user to user database
    const newUser = new User({
        uuid: uniqueID,
        ...newPayload
    });

    try {
        await newUser.save();
        res.status(200).json({ message: 'User created successfully' });
    } catch (err) {
        res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }

};

export { login, logout, refresh, register };
