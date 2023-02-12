import jwt, { Secret, JwtPayload } from 'jsonwebtoken';
import { Request, Response, NextFunction } from 'express';
import * as dotenv from 'dotenv';
import nodeConfig from 'config';
dotenv.config();

// Define the secret key
export const SECRET_KEY: Secret = process.env.JWT_SECRET_KEY!;

// Define the custom request interface
export interface CustomRequest extends Request {
    token: string | JwtPayload;
}

// Define the payload interface
export interface TokenPayload {
    email: string;
    type: string;
    iat: number;
    exp: number;
}

/**
 * Verify the auth token or return an appropriate error code
 */
export const verifyToken = async (req: Request, res: Response, next: NextFunction) => {
    const token = req.header('Authorization')?.replace('Bearer ', '');
    if (!token) {
        return res.status(401).json({ error: nodeConfig.get('error_codes.TOKEN_NOT_FOUND') });
    }
    try {
        const decoded = await jwt.verify(token, SECRET_KEY);
        // Check if the token has the correct type
        if ((decoded as TokenPayload).type !== 'access') {
            return res.status(401).json({
                error: nodeConfig.get('error_codes.INVALID_TOKEN'),
            });
        }
        (req as CustomRequest).token = decoded;
        next();
    } catch (error: Error | any) {
        // if the token is expired, return a new access token and refresh token
        if (error instanceof jwt.TokenExpiredError) {
            return res.status(401).json({
                error: nodeConfig.get('error_codes.TOKEN_EXPIRED'),
            });
        } else if (error instanceof jwt.JsonWebTokenError) {
            return res.status(401).json({
                error: nodeConfig.get('error_codes.INVALID_TOKEN'),
            });
        } else {
            return res.status(500).json({
                error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'),
            });
        }
    }
};

export const verifyRefreshToken = async (req: Request, res: Response, next: NextFunction) => {
    const token = req.cookies.refresh_token;
    if (!token) {
        return res.status(401).json({ error: nodeConfig.get('error_codes.TOKEN_NOT_FOUND') });
    }
    try {
        const decoded = await jwt.verify(token, SECRET_KEY);
        // Check if the token has the correct type
        if ((decoded as TokenPayload).type !== 'refresh') {
            return res.status(401).json({
                error: nodeConfig.get('error_codes.INVALID_TOKEN'),
            });
        }
        (req as CustomRequest).token = decoded;
        next();
    } catch (error: Error | any) {
        // if the token is expired, return a new access token and refresh token
        if (error instanceof jwt.TokenExpiredError) {
            return res.status(401).json({
                error: nodeConfig.get('error_codes.TOKEN_EXPIRED'),
            });
        } else if (error instanceof jwt.JsonWebTokenError) {
            return res.status(401).json({
                error: nodeConfig.get('error_codes.INVALID_TOKEN'),
            });
        } else {
            return res.status(500).json({
                error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'),
            });
        }
    }
};
