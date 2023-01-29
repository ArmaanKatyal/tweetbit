import {Schema, model} from 'mongoose';

interface IUser {
    username: string;
    first_name: string;
    last_name: string;
    email: string;
}

const userSchema = new Schema<IUser>({
    username: {
        type: String,
        required: true,
        unique: true
    },
    first_name: {
        type: String,
        required: true
    },
    last_name: {
        type: String,
        required: true
    },
    email: {
        type: String,
        required: true,
        unique: true
    }
});

export const User = model<IUser>('User', userSchema);
