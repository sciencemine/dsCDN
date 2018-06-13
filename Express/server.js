//---------------------------------
//- Built-In Node Module Requires -
//---------------------------------
const path = require('path');

//------------------------
//- Node Module Requires -
//------------------------
const express = require('express');

//---------------------
//- Additional Routes -
//---------------------
const ce = require('./routes/ce');
const dsm = require('./routes/dsm')

//--------------------
//- Script Constants -
//--------------------
const app = express();
const hostname = 'csdept26.mtech.edu';
const port = 30120;
// path to the public files (this is the ember front-end)
const static_dir = '../WebApp/dist';
const static_path = path.join(__dirname, static_dir);

//-------------------
//- Server Listener -
//-------------------
app.listen(port, () => {
    console.log(`Listening at "${hostname}:${port}"`);
});

//--------------------
//- App Static Files -
//--------------------
app.use(ce);
app.use(dsm);
// app.use(express.static(static_path, { dotfiles: 'ignore' }));

//-----------------
//- Set up Routes -
//-----------------
// Bad request
function badRequest(req, res) {
    res.status(400).send('Bad Request\n');
}

// All other routes
app.route('*')
// Define get requests to 404
.get((req, res) => {
    res.status(404).send('Page not Found\n');
})
.all(badRequest);

