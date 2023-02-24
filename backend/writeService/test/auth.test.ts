import chai from 'chai';
import {verifyToken} from '../src/middlewares/auth.middleware';
import * as sinon from 'sinon';
import { Request, Response } from 'express';

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
});