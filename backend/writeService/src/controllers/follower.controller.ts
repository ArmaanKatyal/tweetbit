import { Request, Response } from 'express';
import { checkIfUserExists } from '../helpers/verifyUser.helper';
import nodeConfig from 'config';
import prisma from '../../prisma/client';
import { userClient } from '../services/user.service';

export const followUser = async (req: Request, res: Response) => {
    let { email, uuid } = (req as any).token;
    let { userToFollowEmail } = req.params;
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

        let newFollower = await prisma.user_Followers.create({
            data: {
                user_id: user_id!,
                follower: {
                    connect: {
                        id: userToFollowId!,
                    },
                },
            },
        });

        userClient.FollowUser(
            {
                userId: user_id!.toString(),
                followerId: userToFollowEmail,
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

export const unfollowUser = async (req: Request, res: Response) => {};