import express from 'express';
import * as auth from '../controllers/auth.controller';
import { verifyRefreshToken } from '../middlewares/auth.middleware';

const router = express.Router();

router.post('/login', auth.login);
router.post('/logout', auth.logout);
router.get('/refresh', verifyRefreshToken, auth.refresh);
router.post('/register', auth.register);

export { router as authRouter };
