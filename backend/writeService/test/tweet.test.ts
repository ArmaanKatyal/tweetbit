import { describe, expect, it } from "vitest";
import request from 'supertest'
import {app} from '../src/app'
import dotenv from 'dotenv';
dotenv.config();
const test_token = process.env.TEST_TOKEN || '';

describe('/api/tweet', async () => {
    describe('[POST] /api/tweet/create', () => {
        it('should return 201 and the created tweet', async () => {
            const {status, body} = await request(app).post('/api/tweet/create').send({
                content: 'TEST'
            }).set('Authorization', 'Bearer ' + test_token);
            console.log(body.error);
            expect(status).toBe(201);
        });
    });
});