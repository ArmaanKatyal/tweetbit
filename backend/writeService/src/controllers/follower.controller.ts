import { Request, Response } from 'express';
import { checkIfUserExists } from '../helpers/verifyUser.helper';
import nodeConfig from 'config';
import prisma from '../../prisma/client';
import { userClient } from '../services/user.service';

import opentelemetry, { SpanStatusCode } from '@opentelemetry/api';
import { getTracer } from '../utils/opentelemetry.util';
import { MetricsCode, MetricsMethod, collectMetrics } from '../internal/prometheus';
import Logger from '../internal/Logger';

/**
 * follow the user with the given email
 * @param req {Request} {params: {userEmail: string}
 * @param res {Response}
 */
export const followUser = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { userEmail: userToFollowEmail } = req.params;
    try {
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('followUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        let [userToFollowExists, userToFollow] = await checkIfUserExists(userToFollowEmail);
        if (!userToFollowExists) {
            Logger.error({
                message: 'User to follow not found',
                email,
                uuid,
            });
            collectMetrics('followUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // Check if the user is trying to follow themselves
        if (user?.id === userToFollow?.id) {
            Logger.error({
                message: 'User cannot follow themselves',
                email,
                uuid,
            });
            collectMetrics('followUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: 'action_not_allowed' });
        }

        // Check if the user is already following the user
        let isAlreadyFollowing = await prisma.user_Followers.findFirst({
            where: {
                user_id: userToFollow?.id!,
                follower_id: user?.id!,
            },
        });
        if (isAlreadyFollowing) {
            Logger.error({
                message: 'User is already following the user',
                email,
                uuid,
            });
            collectMetrics('followUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: 'action_not_allowed' });
        }

        // Increase the follower count
        let userWithIncreasedFollowerCount = await prisma.user.update({
            where: {
                id: userToFollow?.id!,
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
                id: user?.id!,
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

        const span = getTracer().startSpan('/follow');
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
                        Logger.error({
                            message: error.message,
                            email,
                            uuid,
                        });
                        collectMetrics(
                            'followUser',
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
            message: 'User followed',
            email,
            uuid,
        });
        span.setStatus({ code: SpanStatusCode.OK });
        span.end();
        collectMetrics('followUser', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json({
            newFollower,
            user: followerWithIncreasedFollowingCount,
            user_followed: userWithIncreasedFollowerCount,
        });
    } catch (error: Error | any) {
        Logger.error({
            message: error.message,
            email,
            uuid,
        });
        collectMetrics('followUser', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};

/**
 * unfollow the user with the given email
 * @param req {Request} {params: {userEmail: string}}
 * @param res {Response}
 */
export const unfollowUser = async (req: Request, res: Response) => {
    let start = Date.now();
    let { email, uuid } = (req as any).token;
    let { userEmail: userToUnfollowEmail } = req.params;
    try {
        let [userExists, user] = await checkIfUserExists(email);
        if (!userExists) {
            Logger.error({
                message: 'User not found',
                email,
                uuid,
            });
            collectMetrics('unfollowUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        let [userToUnfollowExists, userToUnfollow] = await checkIfUserExists(userToUnfollowEmail);
        if (!userToUnfollowExists) {
            Logger.error({
                message: 'User to unfollow not found',
                email,
                uuid,
            });
            collectMetrics('unfollowUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }

        // Check if the user is trying to unfollow themselves
        if (user?.id === userToUnfollow?.id) {
            Logger.error({
                message: 'User cannot unfollow themselves',
                email,
                uuid,
            });
            collectMetrics('unfollowUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: 'action_not_allowed' });
        }

        // Check if the user is already following the user
        let deletedFollower = await prisma.user_Followers.deleteMany({
            where: {
                user_id: userToUnfollow?.id!,
                follower_id: user?.id!,
            },
        });
        if (!deletedFollower) {
            Logger.error({
                message: 'User is not following the user',
                email,
                uuid,
            });
            collectMetrics('unfollowUser', MetricsCode.BadRequest, MetricsMethod.Post, Date.now() - start);
            return res.status(400).json({ error: 'action_not_allowed' });
        }

        let userWithDecreasedFollowerCount = await prisma.user.update({
            where: {
                id: userToUnfollow?.id!,
            },
            data: {
                followers_count: {
                    decrement: 1,
                },
            },
        });

        let userWithDecreasedFollowingCount = await prisma.user.update({
            where: {
                id: user?.id!,
            },
            data: {
                following_count: {
                    decrement: 1,
                },
            },
        });

        const span = getTracer().startSpan('/unfollow');
        opentelemetry.context.with(opentelemetry.trace.setSpan(opentelemetry.context.active(), span), () => {
            span.setAttribute('userId', userWithDecreasedFollowerCount.id.toString());
            span.setAttribute('followerId', userWithDecreasedFollowingCount.id.toString());

            // contact the fanout service
            userClient.UnfollowUser(
                {
                    userId: userWithDecreasedFollowerCount.id.toString(),
                    followerId: userWithDecreasedFollowingCount.id.toString(),
                },
                (error) => {
                    if (error) {
                        Logger.error({
                            message: error.message,
                            email,
                            uuid,
                        });
                        collectMetrics(
                            'unfollowUser',
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
            message: 'User unfollowed',
            email,
            uuid,
        });
        span.setStatus({ code: SpanStatusCode.OK });
        span.end();
        collectMetrics('unfollowUser', MetricsCode.Ok, MetricsMethod.Post, Date.now() - start);
        res.status(200).json({
            deletedFollower,
            user: userWithDecreasedFollowingCount,
            user_unfollowed: userWithDecreasedFollowerCount,
        });
    } catch (error: Error | any) {
        Logger.error({
            message: error.message,
            email,
            uuid,
        });
        collectMetrics('unfollowUser', MetricsCode.InternalServerError, MetricsMethod.Post, Date.now() - start);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};
