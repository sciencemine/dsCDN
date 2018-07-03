const assert  = require('assert')
const Server  = require('../server')
const asset   = require('../routes/asset')
const ce      = require('../routes/ce')
const dsm     = require('../routes/dsm')
const request = require('supertest')
const url = 'http://localhost:30120'

describe('Routes', () => {

    //test get reqs for all endpoints
    it('[GET] /ce', (done) =>  {
        request(url)
        .get('/ce')
        .expect(200, done)
    })

    it('[GET] /dsm', (done) => {
        request(url)
        .get('/dsm')
        .expect(200, done)
    })

    it('[GET] /asset', (done) => {
        request(url)
        .get('/asset')
        .expect(200, done)
    })

    //test post requests of all endpoints
    it('[POST] /ce', (done) => {
        request(url)
        .post('/ce')
        .expect(200, done)
    })

    it('[POST] /dsm', (done) => {
        request(url)
        .post('/dsm')
        .expect(200, done)
    })

    it('[POST] /asset', (done) => {
        request(url)
        .post('/asset')
        .expect(200, done)
    })

    //all other requests should error
    
    it('[GET] FAIL', (done) => {
        request(url)
        .post('/unicorns/farts')
        .expect(400, done)
    })




})