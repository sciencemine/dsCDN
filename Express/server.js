//------------------------
//- Node Module Requires -
//------------------------
const express = require('express');

//---------------------
//- Additional Routes -
//---------------------
const ce    = require('./routes/ce');
const dsm   = require('./routes/dsm')
const asset = require('./routes/asset')
const graph_path  = require('./routes/path')

//--------------------
//- Script Constants -
//--------------------
const hostname = 'csdept26.mtech.edu';
const port = 30120;

var allowCrossDomain = function(req, res, next) {
    if ('GET' === req.method || 'POST' === req.method || 'PUT' === req.method) {
        res.header('Access-Control-Allow-Origin', '*');
        next();
    }
    else if ('OPTIONS' === req.method) {
        res.header('Access-Control-Allow-Origin', '*');
        res.header('Access-Control-Allow-Methods', 'GET,PUT,POST');
        res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization, Content-Length, X-Requested-With');
        res.sendStatus(200);
    }
};


function Server() {
    let app = express();
    //-------------------
    //- Server Listener -
    //-------------------
    

    //Allow CORS
    app.use(allowCrossDomain)

    //-----------------
    //- Set up Routes -
    //-----------------
    app.use(ce);
    app.use(dsm);
    app.use(asset)
    app.use(graph_path)

    // All other routes
    app.route('*')
    // Define get requests to 404
    .get((req, res) => {
        res.status(404).send('Page not Found\n');
    })
    .all(badRequest);

    return app.listen(port, () => {
        console.log(`Listening at "${hostname}:${port}"`);
    });
}

// Bad request
function badRequest(req, res) {
    res.status(400).send('Bad Request\n');
}

module.exports = Server
