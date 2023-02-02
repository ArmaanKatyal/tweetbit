import express from 'express';
import * as auth from '../controllers/auth.controller';
import { verifyRefreshToken, verifyToken } from '../middleware/auth.middleware';

const router = express.Router();

router.post('/login', auth.login);
router.post('/logout', auth.logout);
router.get('/refresh', verifyRefreshToken, auth.refresh);
router.get('/temp', verifyToken, auth.temp);

export { router as authRouter };
