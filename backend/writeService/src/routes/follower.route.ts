import express from 'express';
import { verifyToken } from '../middlewares/auth.middleware';
import { followUser } from '../controllers/follower.controller';

const router = express.Router();

router.post('/follow/:userId', verifyToken, followUser);
router.post('/unfollow/:userId', verifyToken);

export { router as followRouter };
