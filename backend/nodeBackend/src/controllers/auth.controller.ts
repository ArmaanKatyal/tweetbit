import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { SECRET_KEY } from '../middleware/auth.middleware';

const login = async (req: Request, res: Response) => {
    
};

const logout = async (req: Request, res: Response) => {
   
};

export { login, logout };
