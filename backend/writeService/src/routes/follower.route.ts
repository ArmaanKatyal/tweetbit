import express from 'express';
import { verifyToken } from '../middlewares/auth.middleware';

const router = express.Router();

router.post('/follow', verifyToken);
router.post('/unfollow', verifyToken);

export { router as followRouter };
