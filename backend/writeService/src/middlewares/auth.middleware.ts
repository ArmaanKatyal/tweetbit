import jwt, { Secret, JwtPayload } from 'jsonwebtoken';
import { Request, Response, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import nodeConfig from 'config';
dotenv.config();

// Define the secret key
export const SECRET_KEY: Secret = process.env.JWT_SECRET_KEY!;

// Deine the custom request interface
export interface CustomRequest extends Request {
    token: string | JwtPayload;
}

// Define the payload interface
export interface TokenPayload {
    email: string;
    uuid: string;
    type: string;
    iat: number;
    exp: number;
}

export const verifyToken = (req: Request, res: Response, next: NextFunction) => {
    const token = req.header('Authorization')?.replace('Bearer ', '');
    if (!token) {
        return res.status(401).json({ error: nodeConfig.get('error_codes.TOKEN_NOT_FOUND') });
    }
    try {
        const decoded = jwt.verify(token, SECRET_KEY) as TokenPayload;
        if (decoded.type !== 'access') {
            return res.status(401).json({ error: nodeConfig.get('error_codes.INVALID_TOKEN') });
        }
        (req as CustomRequest).token = decoded;
        next();
    } catch (error: Error | any) {
        if (error instanceof jwt.TokenExpiredError) {
            return res.status(401).json({ error: nodeConfig.get('error_codes.TOKEN_EXPIRED') });
        } else if (error instanceof jwt.JsonWebTokenError) {
            return res.status(401).json({ error: nodeConfig.get('error_codes.INVALID_TOKEN') });
        } else {
            return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
        }
    }
};