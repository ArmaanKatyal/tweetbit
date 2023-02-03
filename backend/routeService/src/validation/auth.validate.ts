import joi from 'joi';
import { ILogin } from '../models/auth.model';
import { IRegister } from '../models/user.model';

/**
 * validate the login data before to avoid errors and add security to the app
 */
export const validateLogin = (data: ILogin) => {
    const schema = joi.object({
        email: joi.string().required().email(),
        password: joi.string().required(),
    });
    return schema.validate(data);
};

/**
 * validate the register data before to avoid errors and add security to the app
 */
export const validateRegister = (data: IRegister) => {
    const schema = joi.object({
        username: joi.string().required(),
        first_name: joi.string().required(),
        last_name: joi.string().required(),
        email: joi.string().required().email(),
        password: joi.string().required(),
    });
    return schema.validate(data);
};
