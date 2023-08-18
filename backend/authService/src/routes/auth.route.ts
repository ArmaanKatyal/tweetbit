import express from 'express';
import * as auth from '../controllers/auth.controller';

const router = express.Router();

router.post('/login', auth.login);
router.post('/logout', auth.logout);
router.get('/check_token', auth.checkToken);
router.get('/refresh', auth.refresh);
router.post('/register', auth.register);

export { router as authRouter };
