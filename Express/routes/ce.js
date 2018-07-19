// This route is for defining ce requests and stuff
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
        mongoURL = 'mongodb://localhost:27017',
        dbName = 'ds',
        ceColName = 'ces',
		bodyParser = require('body-parser')

//local requires
const CE = require('../classes/ce'),
    utils = require('./util')

let router = express.Router();

router.use(bodyParser.json())
router.use(bodyParser.urlencoded({ extended: true }))

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
        utils.query(db, 'ces', ceID)
        .then((doc) => {
            let ceObj = doc,
                    promises = [],
                    hasTeaser = false;

            // if the playlist has a teaser query for it
            if (!Array.isArray(ceObj.playlist[0])) {
                utils.query(db, 'assets', new mongo.ObjectID(ceObj.playlist[0]))
                .then((doc) => {
                    ceObj.playlist[0] = doc;
                })
                .catch(utils.queryErr);

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
                        .catch(utils.queryErr));

                // replace all concurrent assets
                assetList[1].forEach((conAsset, index) => {
                    promises.push(query(db, 'assets', new mongo.ObjectID(conAsset))
                            .then((doc) => {
                                assetList[1][index] = doc;
                            })
                            .catch(utils.queryErr));
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
    MongoClient.connect(mongoURL, { useNewUrlParser: true }, (err, client) => {
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

                res.status(500).send();

                return;
            }

            res.status(200).json(docs);
        });
    });
})
.post((req, res) => {
    let ce = req.body

    MongoClient.connect(mongoURL, { useNewUrlParser: true }, (err, client) =>  {
        if (err) {
            console.error(err)

            res.status(500)

            return
        }

        const db = client.db(dbName)

        let ceObj = new CE(ce.title, ce.version, ce.playlist, ce)

        return utils.add(db, ceColName, ceObj)
        .then((results) => {
            console.log(`added a ce with _id: ${results.insertedId}`)
            res.status(201).send(`added a ce with _id: ${results.insertedId}`)
        })
        .catch(() => {
            res.status(500)
        })
    })
})
.all((req, res) => {
    res.status(405);
});

module.exports = router;
