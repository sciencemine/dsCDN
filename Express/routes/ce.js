// This route is for defining ce requests and stuff
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
        mongoURL = 'mongodb://localhost:27017',
        dbName = 'ds';

let router = express.Router();

router.route('/ce/:id')
.get((req, res) => {
    let ceID = new mongo.ObjectID(req.params.id);

    // get the ce from mongo
    // let ceObj = testData;

    MongoClient.connect(mongoURL, (err, client) => {
        if (err) throw err;

        const db = client.db(dbName);

        // do the query
        query(db, 'ces', ceID,)
        .then((doc) => {
            let ceObj = doc,
                    promises = [],
                    hasTeaser = false;

            if (!Array.isArray(ceObj.playlist[0])) {
                query(db, 'assets', new mongo.ObjectID(ceObj.playlist[0]))
                .then((doc) => {
                    ceObj.playlist[0] = doc;
                })
                .catch((err) => {
                    throw err;
                });

                hasTeaser = true;
            }

            // replace all assets in the playlist
            for (let i = (hasTeaser ? 1 : 0); i < ceObj.playlist.length; i++) {
                let item = ceObj.playlist[i];

                promises.push(query(db, 'assets', new mongo.ObjectID(item[0]))
                    .then((doc) => {
                        item[0] = doc;
                    })
                    .catch((err) => {
                        throw err;
                    })
                );

                // replace all concurrent assets
                item[1].forEach((conAsset, index) => {
                    promises.push(query(db, 'assets', new mongo.ObjectID(conAsset))
                        .then((doc) => {
                            item[1][index] = doc;
                        })
                        .catch((err) => {
                            throw err;
                        })
                    );
                });
            }

            Promise.all(promises)
            .then((values) => {
                res.status(200).json(ceObj);
            });
        })
        .catch((err) => {
            throw err;
            res.status(404).json(err);
        });
    });
})
.all((req, res) => {
    res.status(403).send('Done Bork.\n')
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

module.exports = router;
