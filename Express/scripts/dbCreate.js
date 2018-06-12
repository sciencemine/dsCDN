// node requires
const assert = require('assert'),
    path = require('path'),
    fs = require('fs');

// package requires
const MongoClient = require('mongodb').MongoClient;

// local requires
const Asset = require('../classes/asset'),
    CE = require('../classes/ce'),
    DSM = require('../classes/dsm');

// constants
const mongoURL = 'mongodb://localhost:27017',
    dbName = 'ds',
    testDataDir = './testData',
    assetColName = 'assets',
    ceColName = 'ces',
    dsmColName = 'dsms';

MongoClient.connect(mongoURL, { useNewUrlParser: true }, (err, client) => {
    assert.equal(null, err);
    console.log(`Connected to the server at: ${mongoURL}`);

    const db = client.db(dbName);

    console.log('Dropping the collections.');

    // first we need to drop all the documents
    Promise.all([ assetColName, ceColName, dsmColName ].map((collectionName) => {
        return drop(db, collectionName)
        .then(() => {
            console.log(`Dropped the collection: ${collectionName}`);
        })
        .catch(() => {
            console.log(`Failed to drop the collection: ${collectionName}`);
        })
    }))
    // fall through if there is an error
    .catch()
    // create all the assets once done dropping the collections
    .then(() => {
        // first insert the assets
        console.log('Inserting assets.');

        fs.readdir(path.join(testDataDir, 'assets'), (err, assetFiles) => {
            if (err) throw err;

            let assetMap = { };

            Promise.all(assetFiles.map((file) => {
                let asset = JSON.parse(fs.readFileSync(path.join(testDataDir,
                        'assets', file))),
                    assetId = asset.id;

                // there is an id in the json and we don't want to push it into the db
                delete asset.id;
                
                let assetObj = new Asset(asset.version, asset.url,
                        asset.type, asset);

                return add(db, assetColName, assetObj)
                .then((result) => {
                    assetMap[assetId] = result.insertedId.toString();
                    console.log(`Inserted an asset with _id: ${result.insertedId}`);
                });
            }))
            // then insert the ces
            .then(() => {
                console.log('Inserting the ces.');

                fs.readdir(path.join(testDataDir, 'ces'), (err, ceFiles) => {
                    if (err) throw err;

                    let ceMap = { };

                    Promise.all(ceFiles.map((file) => {
                        let ce = JSON.parse(fs.readFileSync(path.join(testDataDir,
                                'ces', file))),
                                ceId = ce.id;

                        // there is an id in the json and we don't want to push it into the db
                        delete ce.id;

                        // fix the ce to point to the correct id of asset
                        ce.playlist.forEach((item, index) => {
                            if (!Array.isArray(item)) {
                                ce.playlist[index] = assetMap[item];
                            }
                            else {
                                ce.playlist[index][0] = assetMap[item[0]];

                                item[1].forEach((concurrent, conIndex) => {
                                    ce.playlist[index][1][conIndex] = assetMap[concurrent];
                                });
                            }
                        });
                        
                        let ceObj = new CE(ce.title, ce.version, ce.playlist, ce);

                        return add(db, ceColName, ceObj)
                        .then((result) => {
                            ceMap[ceId] = result.insertedId.toString();
                            console.log(`Inserted a ce with _id: ${result.insertedId}`);
                        })
                    }))
                    // then insert the dsm
                    .then(() => {
                        console.log('Inserting the dsm.');

                        let dsm = JSON.parse(fs.readFileSync(path.join(testDataDir,
                                'hst.json'))),
                            ce_set = { };;

                        // want to manually put this in here so i can read it
                        dsm._id = dsm.id;

                        // delete the old one
                        delete dsm.id;

                        // fix the references to ces
                        dsm.idle_backgrounds.forEach((ce, index) => {
                            dsm.idle_backgrounds[index] = ceMap[ce];
                        });

                        dsm.video_select_backgrounds.forEach((ce, index) => {
                            dsm.video_select_backgrounds[index] = ceMap[ce];
                        });

                        for (let ce in dsm.ce_set) {
                            ce_set[ceMap[ce]] = dsm.ce_set[ce];
                        }

                        dsm.ce_set = ce_set;
            
                        let dsmObj = new DSM(dsm._id, dsm.title, dsm.description,
                                dsm.version, dsm.author, dsm.config,
                                dsm.stylesheet, dsm.style, dsm.contributors,
                                dsm.idle_backgrounds,
                                dsm.video_select_backgrounds, dsm.ce_set,
                                dsm.attributes);

                        return add(db, dsmColName, dsmObj)
                        .then((result) => {
                            console.log(`Inserted a dsm with _id: ${result.insertedId}`);
                        })
                        .then(() => {
                            client.close();
                        });
                    });
                });
            });
        });
    });
});

function drop(db, collectionName) {
    return new Promise((resolve, reject) => {
        let collection = db.collection(collectionName);

        collection.drop().then(resolve).catch(reject);
    });
}

function add(db, collectionName, obj, opts = { }) {
    return new Promise((resolve, reject) => {
        let collection = db.collection(collectionName);

        collection.insertOne(obj, opts).then(resolve).catch(reject);
    });
}