import { Request, Response } from 'express';
import prisma from '../../prisma/client'
import nodeConfig from 'config';

// Convert BigInt to string
(BigInt.prototype as any).toJSON = function () {
    return this.toString();
};

export const createTweet = async (req: Request, res: Response) => {
    // create a tweet for the user in the database
    let { email, uuid } = (req as any).token;
    let { content } = req.body;
    try {
        let user = await prisma.user.findUnique({
            where: {
                email
            }
        });
        if (!user) {
            return res.status(404).json({ error: nodeConfig.get('error_codes.USER_NOT_FOUND') });
        }
        let tweet = await prisma.tweet.create({
            data: {
                uuid,
                content,
                user: {
                    connect: {
                        id: user.id
                    }
                }
            }
        });
        return res.status(201).json({ tweet });
    } catch (error: Error | any) {
        console.log(error);
        return res.status(500).json({ error: nodeConfig.get('error_codes.INTERNAL_SERVER_ERROR') });
    }
};
