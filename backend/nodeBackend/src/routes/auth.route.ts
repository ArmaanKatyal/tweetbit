import express from 'express';
import * as auth from '../controllers/auth.controller';
import { verifyToken } from '../middleware/auth.middleware';

const router = express.Router();

router.post('/login', auth.login);
router.post('/logout', verifyToken, auth.logout);

export { router as authRouter };
