import { Request, Response } from 'express';
import { checkIfUserExists } from '../helpers/verifyUser.helper';
import nodeConfig from 'config';
import prisma from '../../prisma/client';
import { userClient } from '../services/user.service';

import opentelemetry from '@opentelemetry/api';
import { initTracer } from '../utils/opentelemetry.util';

/**
 * follow the user with the given email
 * @param req {Request} {params: {userEmail: string}
 * @param res {Response}
 */
export const followUser = async (req: Request, res: Response) => {
    let { email, uuid } = (req as any).token;
    let { userEmail: userToFollowEmail } = req.params;
    try {
        let [userExists, user_id] = await checkIfUserExists(email);
        if (!userExists) {
            req.log.error({
                message: 'User not found',
                email,
                uuid,
            });
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        let [userToFollowExists, userToFollowId] = await checkIfUserExists(userToFollowEmail);
        if (!userToFollowExists) {
            req.log.error({
                message: 'User to follow not found',
                email,
                uuid,
            });
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // Increase the follower count
        let userWithIncreasedFollowerCount = await prisma.user.update({
            where: {
                id: userToFollowId!,
            },
            data: {
                followers_count: {
                    increment: 1,
                },
            },
        });

        // Increase the following count
        let followerWithIncreasedFollowingCount = await prisma.user.update({
            where: {
                id: user_id!,
            },
            data: {
                following_count: {
                    increment: 1,
                },
            },
        });

        // create a new follower in the database
        let newFollower = await prisma.user_Followers.create({
            data: {
                user_id: userWithIncreasedFollowerCount.id,
                follower: {
                    connect: {
                        id: followerWithIncreasedFollowingCount.id,
                    },
                },
            },
        });

        const span = initTracer().startSpan('/follow');
        opentelemetry.context.with(opentelemetry.trace.setSpan(opentelemetry.context.active(), span), () => {
            span.setAttribute('userId', userWithIncreasedFollowerCount.id.toString());
            span.setAttribute('followerId', followerWithIncreasedFollowingCount.id.toString());

            // Contact the fanout service
            userClient.FollowUser(
                {
                    userId: userWithIncreasedFollowerCount.id.toString(),
                    followerId: followerWithIncreasedFollowingCount.id.toString(),
                },
                (error) => {
                    span.end();
                    if (error) {
                        req.log.error({
                            message: error.message,
                            email,
                            uuid,
                        });
                        return res
                            .status(500)
                            .json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
                    }
                }
            );
        });

        req.log.info({
            message: 'User followed',
            email,
            uuid,
        });
        res.status(200).json({
            newFollower,
            user: followerWithIncreasedFollowingCount,
            user_followed: userWithIncreasedFollowerCount,
        });
    } catch (error: Error | any) {
        req.log.error({
            message: error.message,
            email,
            uuid,
        });
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * unfollow the user with the given email
 * @param req {Request} {params: {userEmail: string}}
 * @param res {Response}
 */
export const unfollowUser = async (req: Request, res: Response) => {
    let { email, uuid } = (req as any).token;
    let { userEmail: userToUnfollowEmail } = req.params;
    try {
        let [userExists, user_id] = await checkIfUserExists(email);
        if (!userExists) {
            req.log.error({
                message: 'User not found',
                email,
                uuid,
            });
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        let [userToUnfollowExists, userToUnfollowId] = await checkIfUserExists(userToUnfollowEmail);
        if (!userToUnfollowExists) {
            req.log.error({
                message: 'User to unfollow not found',
                email,
                uuid,
            });
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        let userWithDecreasedFollowerCount = await prisma.user.update({
            where: {
                id: userToUnfollowId!,
            },
            data: {
                followers_count: {
                    decrement: 1,
                },
            },
        });

        let userWithDecreasedFollowingCount = await prisma.user.update({
            where: {
                id: user_id!,
            },
            data: {
                following_count: {
                    decrement: 1,
                },
            },
        });

        let deletedFollower = await prisma.user_Followers.deleteMany({
            where: {
                user_id: userToUnfollowId!,
                follower_id: user_id!,
            },
        });

        const span = initTracer().startSpan('/unfollow');
        opentelemetry.context.with(opentelemetry.trace.setSpan(opentelemetry.context.active(), span), () => {
            span.setAttribute('userId', userToUnfollowId!.toString());
            span.setAttribute('followerId', user_id!.toString());

            // contact the fanout service
            userClient.UnfollowUser(
                {
                    userId: user_id!.toString(),
                    followerId: userToUnfollowId!.toString(),
                },
                (error) => {
                    if (error) {
                        req.log.error({
                            message: error.message,
                            email,
                            uuid,
                        });
                        return res
                            .status(500)
                            .json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
                    }
                }
            );
        });

        req.log.info({
            message: 'User unfollowed',
            email,
            uuid,
        });

        res.status(200).json({
            deletedFollower,
            user: userWithDecreasedFollowingCount,
            user_unfollowed: userWithDecreasedFollowerCount,
        });
    } catch (error: Error | any) {
        req.log.error({
            message: error.message,
            email,
            uuid,
        });
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};
