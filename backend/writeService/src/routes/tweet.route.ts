import express from 'express';
import { verifyToken } from '../middlewares/auth.middleware';
import {
    commentTweet,
    createTweet,
    deleteTweet,
    likeTweet,
    retweetTweet,
    unlikeTweet,
} from '../controllers/tweet.controller';

const router = express.Router();

router.post('/create', verifyToken, createTweet);
router.post('/delete/:tweetId', verifyToken, deleteTweet);
router.post('/like/:tweetId', verifyToken, likeTweet);
router.post('/unlike/:tweetId', verifyToken, unlikeTweet);
router.post('/retweet/:tweetId', verifyToken, retweetTweet);
router.post('/comment/:tweetId', verifyToken, commentTweet);

export { router as tweetRouter };
