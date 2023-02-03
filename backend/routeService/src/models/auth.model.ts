import { Collection, Schema, model } from 'mongoose';

export interface IAuth {
    uuid: string;
    email: string;
    password: string;
}

export interface ILogin {
    email: string;
    password: string;
}

const authSchema = new Schema<IAuth>(
    {
        uuid: {
            type: String,
            required: true,
            unique: true,
        },
        email: {
            type: String,
            required: true,
            unique: true,
        },
        password: {
            type: String,
            required: true,
        },
    },
    { collection: 'auth', versionKey: false }
);

export const Auth = model<IAuth>('Auth', authSchema);
