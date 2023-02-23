import prisma from '../client';

export const checkIfUserExists = async (email: string): Promise<[boolean, number | null]> => {
    let user = await prisma.user.findUnique({
        where: {
            email,
        },
    });
    if (!user) {
        return [false, null];
    }
    return [true, user.id];
};
