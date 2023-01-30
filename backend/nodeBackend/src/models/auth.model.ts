import { Collection, Schema, model } from 'mongoose';

export interface Auth {
    email: string;
    password: string;
}

const authSchema = new Schema<Auth>(
    {
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

export const Auth = model<Auth>('Auth', authSchema);
