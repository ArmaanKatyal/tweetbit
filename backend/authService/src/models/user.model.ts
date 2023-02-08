import { Schema, model } from 'mongoose';

interface IUser {
    uuid: string;
    username: string;
    first_name: string;
    last_name: string;
    email: string;
}

export interface IRegister {
    username: string;
    first_name: string;
    last_name: string;
    email: string;
    password: string;
}

const userSchema = new Schema<IUser>(
    {
        uuid: {
            type: String,
            required: true,
            unique: true,
        },
        username: {
            type: String,
            required: true,
            unique: true,
        },
        first_name: {
            type: String,
            required: true,
        },
        last_name: {
            type: String,
            required: true,
        },
        email: {
            type: String,
            required: true,
            unique: true,
        },
    },
    { collection: 'users', versionKey: false }
);

export const User = model<IUser>('User', userSchema);
