const MongoClient = require('mongodb').MongoClient;
const assert = require('assert');
const url = 'mongodb://localhost:27017';
const dbName = 'ds';

MongoClient.connect(url, (err, client) => {
    assert.equal(null, err);
    console.log('Connected to server');

    const db = client.db(dbName);

    insert(db, () => {
        client.close();
    });
});

function ce(title, playlist, description, version = '0.0.0') {
    this.version = version;
    this.title = title;
    this.playlist = playlist;

    if (description) {
        this.description = description;
    }
}

function insert(db, callback) {
    const collection = db.collection('ces');

    collection.insertMany([
        new ce('Test ce', [[ '5b0f0882ac30b02022a8e609', [] ]]),
        new ce('Has Teaser',
            [
                '5b0f0882ac30b02022a8e60b',
                [
                    '5b0f0882ac30b02022a8e609',
                    [ '5b0f0882ac30b02022a8e60a' ]
                ]
            ])
    ], (err, res) => {
        console.log('Inserted the items');
        callback(res);
    });
}
