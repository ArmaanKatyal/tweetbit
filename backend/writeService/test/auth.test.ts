import chai from 'chai';
import {SECRET_KEY, verifyToken} from '../src/middlewares/auth.middleware';
import * as sinon from 'sinon';
import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';

describe('Auth middleware', () => {
    let req: Request, res: Response, next: sinon.SinonSpy<any[], any>;
    beforeEach(() => {
        sinon.createSandbox();
        req = {
            header: sinon.stub().returns(`Bearer ${process.env.TEST_TOKEN}`),
        } as unknown as Request;
        // mock the response object using sinon
        res = {
            status: sinon.stub().returnsThis(),
            json: sinon.stub().returnsThis(),
        } as unknown as Response;
        next = sinon.spy();
    });
    afterEach(() => {
        sinon.restore();
    });

    it('should successfully authenticate user & call next function', () => {
        verifyToken(req, res, next);
        chai.expect(next.called).to.be.true;
        chai.expect((req as any).token).to.have.property('email');
        chai.expect((req as any).token).to.have.property('uuid');
        chai.expect((req as any).token).to.have.property('type');
    });

    it('should not call next function as token is missing', () => {
        req.header = sinon.stub().returns(undefined);
        verifyToken(req, res, next);
        chai.expect(next.called).to.be.false;
    });

    it('should not call next function as token is invalid', () => {
        req.header = sinon.stub().returns(`Bearer ${process.env.TEST_TOKEN}123`);
        verifyToken(req, res, next);
        chai.expect(next.called).to.be.false;
    });

    it('should not call next function as token is expired', () => {
        const expiredToken = jwt.sign({
            email: 'test@test.com',
            uuid: '123',
            type: 'access',
        }, SECRET_KEY, { expiresIn: '1ms' });

        req.header = sinon.stub().returns(`Bearer ${expiredToken}`);
        verifyToken(req, res, next);
        chai.expect(next.called).to.be.false;
    });

    it('should not call next function as token type is not access', () => {
        const refreshToken = jwt.sign({
            email: 'test@test.com',
            uuid: '123',
            type: 'refresh',
        }, SECRET_KEY, { expiresIn: '1h' });

        req.header = sinon.stub().returns(`Bearer ${refreshToken}`);
        verifyToken(req, res, next);
        chai.expect(next.called).to.be.false;
    });
});
