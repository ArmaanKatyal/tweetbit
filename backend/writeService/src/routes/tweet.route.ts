import express from 'express';
import { verifyToken } from '../middlewares/auth.middleware';
import { createTweet } from '../controllers/tweet.controller';

const router = express.Router();

router.post('/create', verifyToken, createTweet);
router.post('/delete/:id', verifyToken);

export { router as tweetRouter };
