const MongoClient = require('mongodb').MongoClient;
const assert = require('assert');
const path = require('path');
const url = 'mongodb://localhost:27017';
const dbName = 'ds';
const remoteURL = '66.62.91.16:60080/smds';

MongoClient.connect(url, (err, client) => {
    assert.equal(null, err);
    console.log('Connected to server');

    const db = client.db(dbName);

    createAssets(db, (res) => {
        console.log(res.insertedIds);
        createCEs(db, res.insertedIds, (res) => {
            console.log(`Created CEs ${res.insertedIds}`);
            client.close();
        });
    });
});

function asset(url, type, title, description, version = '0.0.0', options) {
    this.url = url;
    this.type = type
    this.version = version;

    if (title) {
        this.title = title;
    }
    if (description) {
        this.description = description;
    }
    if  (options) {
        this.options = options;
    }
}

function ce(title, playlist, description, version = '0.0.0') {
    this.version = version;
    this.title = title;
    this.playlist = playlist;

    if (description) {
        this.description = description;
    }
}

function createAssets(db, callback) {
    const collection = db.collection('assets');

    collection.insertMany([
        new asset(path.join(remoteURL, 'SM-HSt-Critr-1s.mp4'), 'video/mp4',
                'Intro to Critters'),
        new asset(path.join(remoteURL, 'SM-HSt-Critr-2.mp4'), 'video/mp4',
                'More on Critters'),
        new asset(path.join(remoteURL, 'SM-HSt-Critr2-1s.mp4'), 'video/mp4',
                'Advanced Critters')
    ], (err, res) => {
        callback(res);
    });
}

function createCEs(db, assets, callback) {
    const collection = db.collection('ces');

    collection.insertMany([
        new ce('Test ce', [[ assets[1], [] ]]),
        new ce('Has Teaser',
            [
                assets[0],
                [
                    assets[1],
                    [ assets[2] ]
                ]
            ]),
        new ce('Lot of Stuff',
            [
                assets[0],
                [
                    assets[1],
                    [ assets[0], assets[1], assets[2] ]
                ],
                [
                    assets[2],
                    [ assets[2], assets[1], assets[0] ]
                ]
            ])
    ], (err, res) => {
        console.log('Inserted the items');
        callback(res);
    });
}
