import { PrismaClient } from '@prisma/client';

declare global {
    namespace NodeJS {
        interface Global {}
    }
}

// add prisma to the NodeJS global type
interface CustomNodeJsGlobal extends NodeJS.Global {
    prisma: PrismaClient;
}

declare const global: CustomNodeJsGlobal;
const prisma = global.prisma || new PrismaClient();

export default async () => {
    await prisma.$transaction([
        prisma.tweet.deleteMany(),
        prisma.user.deleteMany(),
        prisma.user_Followers.deleteMany(),
        prisma.tweet_Likes.deleteMany(),
        prisma.tweet_Comments.deleteMany(),
    ]);
};
