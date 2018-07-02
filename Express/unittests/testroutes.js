const assert = require('assert')
const Server = require('../server')
const asset  = require('../routes/asset')
const ce     = require('../routes/ce')
const dsm    = require('../routes/dsm')

describe('Routes', () => {
    let server

    beforeEach(() => {
        server = Server()
    })

    afterEach((done) => {
        server.close(done)
    })

    it('[GET] /ce', function testIndexGet(done) {
        request(server)
        .get('/ce')
        .expect(200, done)
    })
})