import jwt, { Secret, JwtPayload } from 'jsonwebtoken';
import { Request, Response, NextFunction } from 'express';
import * as dotenv from 'dotenv';
dotenv.config();

export const SECRET_KEY: Secret = process.env.JWT_SECRET_KEY!;

export interface CustomRequest extends Request {
    token: string | JwtPayload;
}

export const verifyToken = async (req: Request, res: Response, next: NextFunction) => {
    const token = req.header('Authorization');
    if (!token) {
        return res.status(401).json({ auth: false, message: 'No token provided' });
    }
    try {
        const decoded = await jwt.verify(token, SECRET_KEY);
        (req as CustomRequest).token = decoded;
        next();
    } catch (error) {
        return res.status(400).json({ auth: false, message: 'Failed to authenticate token' });
    }
};
