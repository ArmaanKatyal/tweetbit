const {MongoClient} = require('mongodb');
const {MongoMemoryServer} = require('mongodb-memory-server');

describe('Single MongoDB Test', () => {
    let con;
    let mongoServer;

    beforeAll(async() => {
        mongoServer = await MongoMemoryServer.create();
        con = await MongoClient.connect(mongoServer.getUri(), {});
    });

    afterAll(async() => {
        if (con) {
            await con.close();
        }
        if (mongoServer) {
            await mongoServer.stop();
        }
    });

    it('should successfully set & get information from a MongoDB instance', async() => {
        const db = con.db(mongoServer.instanceInfo.dbName);
        expect(db).toBeDefined();
        const col = db.collection('test');
        const result = await col.insertMany([{a: 1}, {b: 1}]);
        expect(result.insertedCount).toStrictEqual(2);
        expect(await col.countDocuments({})).toBe(2);
    });
});
