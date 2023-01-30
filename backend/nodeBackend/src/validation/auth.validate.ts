import joi from 'joi';
import { Auth } from '../models/auth.model';

export const validateLogin = (data: Auth) => {
    const schema = joi.object({
        email: joi.string().required().email(),
        password: joi.string().required(),
    });
    return schema.validate(data);
};
