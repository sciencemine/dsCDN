// This route is for defining ce requests and stuff
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
        mongoURL = 'mongodb://localhost:27017',
        dbName = 'ds',
        ceColName = 'ces'

let router = express.Router();

router.route('/ce/:id')
.get((req, res) => {
    let ceID = new mongo.ObjectID(req.params.id);

    // get the ce from mongo
    MongoClient.connect(mongoURL, (err, client) => {
        if (err) {
            console.error(err);

            res.status(500);

            return;
        }

        const db = client.db(dbName);

        // do the query
        query(db, 'ces', ceID)
        .then((doc) => {
            let ceObj = doc,
                    promises = [],
                    hasTeaser = false;

            // if the playlist has a teaser query for it
            if (!Array.isArray(ceObj.playlist[0])) {
                query(db, 'assets', new mongo.ObjectID(ceObj.playlist[0]))
                .then((doc) => {
                    ceObj.playlist[0] = doc;
                })
                .catch(queryErr);

                hasTeaser = true;
            }

            // replace all assets in the playlist
            for (let i = (hasTeaser ? 1 : 0); i < ceObj.playlist.length; i++) {
                let assetList = ceObj.playlist[i]

                // push the query for the primary asset of this sequence into the set
                promises.push(query(db, 'assets', new mongo.ObjectID(assetList[0]))
                        .then((doc) => {
                            assetList[0] = doc;
                        })
                        .catch(queryErr));

                // replace all concurrent assets
                assetList[1].forEach((conAsset, index) => {
                    promises.push(query(db, 'assets', new mongo.ObjectID(conAsset))
                            .then((doc) => {
                                assetList[1][index] = doc;
                            })
                            .catch(queryErr));
                });
            }

            Promise.all(promises)
            .then((values) => {
                res.status(200).json(ceObj);
            })
            .catch(() => {
                res.status(500).send('Something went wrong\n');
            });
        })
        .catch((err) => {
            console.error(err);

            res.status(500);
        });
    });
})
.all((req, res) => {
    res.status(405);
});

router.route('/ce')
.get((req, res) => {
    MongoClient.connect(mongoURL, (err, client) => {
        if (err) {
            console.error(err);

            res.status(500);

            return;
        }


        const db = client.db(dbName);

        let collection = db.collection('ces');

        collection.find({}, { projection: { _id: 1 } }).toArray((err, docs) => {
            if (err) {
                console.error(err);

                res.status(500);

                return;
            }

            res.status(200).json(docs);
        });
    });
})
.post((req, res) => {
    let body = req.body

    MongoClient.connect(mongoURL, (err, client) =>  {
        if (err) {
            console.error(err)

            res.status(500)

            return
        }

        const db = client.db(dbName)

        let ce = JSON.parse(body)

        let ceObj = new CE(ce.title, ce.version, ce.playlist, ce)

        return add(db, ceColName, ceObj)


    })
})
.all((req, res) => {
    res.status(405);
});

function query(db, collectionName, id) {
    return new Promise((resolve, reject) => {
        let collection = db.collection(collectionName);

        collection.find({ _id: id }).toArray((err, doc) => {
            if (err) reject(err);

            resolve(doc[0]);
        });
    });
}

function queryErr(err) {
    console.error(err);
}

//Adds a new item to the database
function add(db, collectionName, obj, opts = { }) {
    return new Promise((resolve, reject) => {
        let collection = db.collection(collectionName)

        collection.insertOne(obj, opts).then(resolve).catch(reject)
    })
}

module.exports = router;
