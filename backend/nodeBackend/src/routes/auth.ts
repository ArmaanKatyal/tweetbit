import express from 'express';
import * as auth from '../controllers/auth';
const router = express.Router();

router.post('/login', auth.login);
router.post('/logout', auth.logout);

export { router as authRouter };