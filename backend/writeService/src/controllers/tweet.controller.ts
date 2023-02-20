import { Request, Response } from 'express';
import prisma from '../../prisma/client';
import nodeConfig from 'config';
import { tweetClient } from '../services/tweet.service';

// Convert BigInt to string
(BigInt.prototype as any).toJSON = function () {
    return this.toString();
};

// What happens when a tweet is created
// 1. Tweet is stored in the sql database on the users timeline
// 2. Contacts the fanoutService using gRPC to publish the tweet to the fanout queue
// 3. broadcast on the kafka/RabbitMQ fanout queue which is consumed by the other services
export const createTweet = async (req: Request, res: Response) => {
    // create a tweet for the user in the database
    let { email, uuid } = (req as any).token;
    let { content } = req.body;
    try {
        // Check if the user exists
        let user = await prisma.user.findUnique({
            where: {
                email,
            },
        });
        // If the user does not exist, return an error
        if (!user) {
            return res.status(404).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }
        // Create the tweet
        let tweet = await prisma.tweet.create({
            data: {
                uuid,
                content,
                user: {
                    connect: {
                        id: user.id,
                    },
                },
            },
        });
        res.status(201).json(tweet);

        // contact the fanout service using gRPC
        tweetClient.createTweet(
            {
                id: tweet.id.toString(),
                content: tweet.content,
                userId: tweet.user_id.toString(),
                uuid: tweet.uuid,
                createdAt: tweet.created_at.toString(),
                likesCount: tweet.likes_count,
                retweetsCount: tweet.retweets_count,
            },
            (err) => {
                if (err) {
                    console.log(err);
                }
            }
        );
    } catch (error: Error | any) {
        console.log(error);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};
