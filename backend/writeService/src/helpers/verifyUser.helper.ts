import { User } from '.prisma/client';
import prisma from '../../prisma/client';

export const checkIfUserExists = async (email: string): Promise<[boolean, User | null]> => {
    let user = await prisma.user.findUnique({
        where: {
            email,
        },
    });
    if (!user) {
        return [false, null];
    }
    return [true, user];
};
