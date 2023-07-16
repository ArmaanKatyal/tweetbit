import { Request, Response } from 'express';
import prisma from '../../prisma/client';
import nodeConfig from 'config';
import { tweetClient } from '../services/tweet.service';
import { checkIfUserExists } from '../helpers/verifyUser.helper';

import opentelemetry, { SpanStatusCode } from '@opentelemetry/api';
import { getTracer } from '../utils/opentelemetry.util';
import { MetricsCode, MetricsMethod, collectMetrics } from '../internal/prometheus';
import Logger from '../internal/Logger';

// Convert BigInt to string
(BigInt.prototype as any).toJSON = function () {
    return this.toString();
};

// What happens when a tweet is created
// 1. Tweet is stored in the sql database on the users timeline
// 2. Contacts the fanoutService using gRPC to publish the tweet to the fanout queue
export const createTweet = async (req: Request, res: Response) => {
    let start = Date.now();
    // create a tweet for the user in the database
    let { email, uuid } = (req as any).token;
    let { content } = req.body;
    try {
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                uuid,
            });
            collectMetrics('createTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // Create the tweet
        let tweet = await prisma.tweet.create({
            data: {
                uuid,
                content,
                user: {
                    connect: {
                        id: user?.id!,
                    },
                },
            },
        });

        const span = getTracer().startSpan('createTweet');
        opentelemetry.context.with(opentelemetry.trace.setSpan(opentelemetry.context.active(), span), () => {
            span.setAttribute('tweetId', tweet.id.toString());
            span.setAttribute('uuid', tweet.uuid);
            // contact the fanout service using gRPC
            tweetClient.createTweet(
                {
                    ...tweet,
                    id: tweet.id.toString(),
                    likesCount: tweet.likes_count.toString(),
                    retweetsCount: tweet.retweets_count.toString(),
                    createdAt: tweet.created_at.toString(),
                },
                (err) => {
                    span.setStatus({ code: SpanStatusCode.ERROR });
                    span.end();
                    if (err) {
                        Logger.error({
                            message: 'Error trasmitting tweet to fanout service',
                            email,
                            uuid,
                        });
                        collectMetrics(
                            'createTweet',
                            MetricsCode.InternalServerError,
                            MetricsMethod.Post,
                            Date.now() - start
                        );
                        return res.status(500).json({
                            error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR'),
                        });
                    }
                }
            );
        });

        Logger.info({
            message: 'Tweet created',
            uuid,
        });
        span.setAttribute('span.kind', 'server');
        span.setStatus({ code: SpanStatusCode.OK });
        span.end();
        collectMetrics('createTweet', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        return res.status(201).json(tweet);
    } catch (error: Error | any) {
        console.log(error);
        Logger.error({
            message: 'Error creating tweet',
            error,
            uuid,
        });
        collectMetrics('createTweet', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * delete the tweet from the database permanently
 * @param req Request {params: {tweetId: string}
 * @param res Response
 */
export const deleteTweet = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { tweetId } = req.params;
    try {
        let [userExists, _] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('deleteTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // delete the tweet
        let tweet = await prisma.tweet.delete({
            where: {
                id: parseInt(tweetId, 10),
            },
        });
        // delete the likes
        let deleted_tweet_likes = await prisma.tweet_Likes.deleteMany({
            where: {
                tweet_id: parseInt(tweetId, 10),
            },
        });
        // delete the comments
        let deleted_tweet_comments = await prisma.tweet_Comments.deleteMany({
            where: {
                tweet_id: parseInt(tweetId, 10),
            },
        });
        collectMetrics('deleteTweet', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json({ tweet, deleted_tweet_likes, deleted_tweet_comments });
        Logger.info({
            message: 'Tweet deleted',
            email,
            uuid,
        });
    } catch (error: Error | any) {
        Logger.error({
            message: 'Error deleting tweet',
            email,
            uuid,
        });
        collectMetrics('deleteTweet', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * increase the likes count of the tweet by 1
 * @param req {Request} {params: {tweetId: string}
 * @param res {Response}
 */
export const likeTweet = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { tweetId } = req.params;
    try {
        // check if the user exists
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('likeTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }
        // check if the user has already liked the tweet
        let ifTweetLiked = await prisma.tweet_Likes.count({
            where: {
                tweet_id: parseInt(tweetId, 10),
                user_id: user?.id!,
            },
        });
        if (ifTweetLiked > 0) {
            // if yes, return an error
            Logger.error({
                message: 'Tweet already liked',
                email,
                uuid,
            });
            collectMetrics('likeTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.TWEET_ALREADY_LIKED') });
        }
        // increment the tweet likes count
        let tweet = await prisma.tweet.update({
            where: {
                id: parseInt(tweetId, 10),
            },
            data: {
                likes_count: {
                    increment: 1,
                },
            },
        });
        // add the user to the tweet likes
        await prisma.tweet_Likes.create({
            data: {
                user: {
                    connect: {
                        id: user?.id!,
                    },
                },
                tweet: {
                    connect: {
                        id: parseInt(tweetId, 10),
                    },
                },
            },
        });

        Logger.info({
            message: 'Tweet liked',
            email,
            uuid,
        });
        collectMetrics('likeTweet', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json(tweet);
    } catch (error: Error | any) {
        Logger.error({
            message: 'Error liking tweet',
            email,
            uuid,
            error,
        });
        collectMetrics('likeTweet', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * decrement the likes count of the tweet by 1
 * @param req {Request} {params: {tweetId: string}
 * @param res
 */
export const unlikeTweet = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { tweetId } = req.params;
    try {
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('unlikeTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // Not allowed to unlike a tweet if the user has not liked it
        let ifTweetLiked = await prisma.tweet_Likes.count({
            where: {
                tweet_id: parseInt(tweetId, 10),
                user_id: user?.id!,
            },
        });

        if (ifTweetLiked === 0) {
            Logger.error({
                message: 'Tweet not liked',
                email,
                uuid,
            });
            collectMetrics('unlikeTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.TWEET_NOT_LIKED') });
        }

        // decrement the tweet likes count
        let tweet = await prisma.tweet.update({
            where: {
                id: parseInt(tweetId, 10),
            },
            data: {
                likes_count: {
                    decrement: 1,
                },
            },
        });

        // delete the user from the tweet likes
        let deleted_like = await prisma.tweet_Likes.deleteMany({
            where: {
                tweet_id: parseInt(tweetId, 10),
                user_id: user?.id!,
            },
        });
        collectMetrics('unlikeTweet', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json({ tweet, deleted_like });
    } catch (error: Error | any) {
        console.log(error);
        Logger.error({
            message: 'Error unliking tweet',
            email,
            uuid,
            error,
        });
        collectMetrics('unlikeTweet', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * copies the original tweet and creates a new tweet with the same content but with new user_id
 * @param req {Request} {params: {tweetId: string}
 * @param res
 */
export const retweetTweet = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { tweetId } = req.params;
    try {
        // check if the user exists
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('retweetTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // The creater of the original tweet should not be able to retweet his own tweet
        let orignalTweet = await prisma.tweet.findFirst({
            where: {
                id: parseInt(tweetId, 10),
            },
        });
        // check if the tweet exists
        if (!orignalTweet) {
            Logger.error({
                message: 'Tweet not found',
                email,
                uuid,
            });
            collectMetrics('retweetTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.TWEET_NOT_FOUND') });
        }

        if (orignalTweet.user_id === user?.id!) {
            Logger.error({
                message: 'User cannot retweet his own tweet',
                email,
                uuid,
            });
            collectMetrics('retweetTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.TWEET_OWNER') });
        }

        // check if the user has already retweeted the tweet
        let ifTweetRetweeted = await prisma.tweet.count({
            where: {
                user_id: user?.id!,
                content: orignalTweet.content,
            },
        });

        if (ifTweetRetweeted > 0) {
            Logger.error({
                message: 'Tweet already retweeted',
                email,
                uuid,
            });
            collectMetrics('retweetTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.TWEET_ALREADY_RETWEETED') });
        }
        // increase the retweet count of the original tweet
        let existingTweet = await prisma.tweet.update({
            where: {
                id: parseInt(tweetId, 10),
            },
            data: {
                retweets_count: {
                    increment: 1,
                },
            },
        });
        // create a new tweet with the same content as the original tweet but with user_id of the retweeter
        let newTweet = await prisma.tweet.create({
            data: {
                uuid,
                content: existingTweet.content,
                user: {
                    connect: {
                        id: user?.id!,
                    },
                },
            },
        });

        // contact the fanout service using gRPC
        tweetClient.createTweet(
            {
                ...newTweet,
                id: newTweet.id.toString(),
                likesCount: newTweet.likes_count.toString(),
                retweetsCount: newTweet.retweets_count.toString(),
                createdAt: newTweet.created_at.toString(),
            },
            (err) => {
                if (err) {
                    Logger.error({
                        message: 'Error trasmitting tweet to fanout service',
                        email,
                        uuid,
                    });
                    collectMetrics(
                        'retweetTweet',
                        MetricsCode.InternalServerError,
                        MetricsMethod.Post,
                        Date.now() - start
                    );
                    return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
                }
            }
        );

        Logger.info({
            message: 'Tweet retweeted',
            email,
            uuid,
        });
        collectMetrics('retweetTweet', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json(existingTweet);
    } catch (error: Error | any) {
        Logger.error({
            message: 'Error retweeting tweet',
            email,
            uuid,
        });
        collectMetrics('retweetTweet', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * add a comment to a tweet
 * @param req {Request} {params: {tweetId: string}
 * @param res
 */
export const commentTweet = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { tweetId } = req.params;
    let { content } = req.body;

    try {
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('commentTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // Check if the tweet exists
        let ifTweetExists = await prisma.tweet.findFirst({
            where: {
                id: parseInt(tweetId, 10),
            },
        });
        if (!ifTweetExists) {
            Logger.error({
                message: 'Tweet not found',
                email,
                uuid,
            });
            collectMetrics('commentTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.TWEET_NOT_FOUND') });
        }

        // Check if the user has already commented on the tweet
        let ifUserCommented = await prisma.tweet_Comments.findFirst({
            where: {
                tweet_id: parseInt(tweetId, 10),
                user_id: user?.id!,
            },
        });

        if (ifUserCommented) {
            Logger.error({
                message: 'User already commented on the tweet',
                email,
                uuid,
            });
            collectMetrics('commentTweet', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_ALREADY_COMMENTED') });
        }

        // create a new comment
        let comment = await prisma.tweet_Comments.create({
            data: {
                content,
                user: {
                    connect: {
                        id: user?.id!,
                    },
                },
                tweet: {
                    connect: {
                        id: parseInt(tweetId, 10),
                    },
                },
            },
        });

        collectMetrics('commentTweet', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json(comment);

        // TODO: contact the fanout service using gRPC
    } catch (error: Error | any) {
        Logger.error({
            message: 'Error commenting tweet',
            email,
            uuid,
        });
        collectMetrics('commentTweet', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};
