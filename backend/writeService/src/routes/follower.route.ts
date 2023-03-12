import express from 'express';
import { verifyToken } from '../middlewares/auth.middleware';
import { followUser, unfollowUser } from '../controllers/follower.controller';

const router = express.Router();

router.post('/follow/:userEmail', verifyToken, followUser);
router.post('/unfollow/:userEmail', verifyToken, unfollowUser);

export { router as followRouter };
