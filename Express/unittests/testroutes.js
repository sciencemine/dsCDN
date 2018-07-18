//local requires
const Server  = require('../server')

// package requires
const request = require('supertest').agent(Server)
const chai = require('chai')
const chaiHttp = require('chai-http')
const should = chai.should()

chai.use(chaiHttp)

describe('Routes', () => {
    let server

    beforeEach(() => {
        server = Server()
    })

    afterEach((done) => {
        server.close()
        done()
    })

    describe('[GET] /ce', (done) => {
        //test get reqs for all endpoints
        it('should get all the CEs', (done) =>  {
            chai.request(server)
            .get('/ce')
            .end((err, res) => {
                res.should.have.status(200)
                res.body.should.be.a('array')
                done() 
            })
        })
    })

    describe('[GET] /dsm', (done) => {
        it('should get a list of all dsms', (done) => {
            chai.request(server)
            .get('/dsm')
            .end((err, res) => {
                res.should.have.status(200)
                res.body.should.be.a('array')
                done()
            })
        })
    })

    describe('[GET] /asset', (done) => {
        it('should get a list of all assets', (done) => {
            chai.request(server)
            .get('/asset')
            .end((err, res) => {
                res.should.have.status(200)
                res.body.should.be.a('array')
                done()
            })
        })
    })

    describe('[GET] /path', (done) => {
        it('Should get a list of all paths', (done) => {
            chai.request(server)
            .get('/path')
            .end((err, res) => {
                res.should.have.status(200)
                res.body.should.be.a('array')
                done()
            })
        })

        it('should get a path with specific id', (done) => {
            chai.request(server)
            .get('/path/5b3e9560affc0b6f593c4ff3')
            .end((err, res) => {
                res.should.have.status(200)
                res.body.should.be.a('object')
                done()
            })
        })
    })

    describe('[POST] /ce', (done) => {
        //test post requests of all endpoints
        ce = {
            "id": "t1000x",
            "version": "0.0.0",
            "title": "testificate",
            "description": "testing woo",
            "playlist": [
                [
                    "t1000x",
                    []
                ]
            ]	
        }

        it('post a ce to the server', (done) => {
            chai.request(server)
            .post('/ce')
            .send(ce)
            .end((err, res) => {
                res.should.have.status(201)
                done()
            })
        })
    })

    describe('[POST] /dsm', (done) => {
        dsm = {
            "id": "",
            "title": "this is just a test",
            "description": "this is just to be used for testing",
            "version": "0.0.0",
            "author": "Justin Bak",
            "config": {
                "idle": 0,
                "menuDwell": 0,
                "popoverDwell": 0,
                "popoverShowDelay": 0
            },
            "stylesheet": "the best stylesheet",
            "style": "gangnam style",
            "contributors": ["Justin", "Me", "Myself", "I"],
            "idle_backgrounds": [],
            "video_select_backgrounds": [],
            "ce_set": {	},
            "attributes": []
        }
        
        it('post a dsm to the server', (done) => {
            chai.request(server)
            .post('/dsm')
            .send(dsm)
            .end((err, res) => {
                res.should.have.status(201)
                res.body.should.be.a('object')
                done()
            })
        })
    })

    describe('[POST] /asset', (done) => {

        asset = {
            "id": "t1000x",
            "version": "0.0.0",
            "url": "www.test.test",
            "type": "testing",
            "options": {}	
        }

        it('posts an asset to the db', (done) => {
            chai.request(server)
            .post('/asset')
            .send(asset)
            .end((err, res) => {
                res.should.have.status(201)
                res.body.should.be.a('object')
                done()
            })
        })
    })

    describe ('[POST] /path', (done) => {
        path = {
            "id": "test",
            "model": {},
            "relations": []
        }
 
        it('posts a path to the db', (done) => {
            chai.request(server)
            .post('/path')
            .send(path)
            .end((err, res) => {
                res.should.have.status(201)
                res.body.should.be.a('object')
                done()
            })
        })
    })

    describe ('[PUT] /path', (done) => {
        update = {
            "relations": [
                {
                    "ce_list": [
                        null,
                        {
                        "id_ce": "5b47c6b8132bb67388ee2530",
                        "version_ce": "0.0.0",
                        "title_ce": "What are \"the Five Cs\"?",
                        "description_ce": "The 5Cs",
                        "events": []
                        }
                    ]
                }
            ] 
        }
        it('[PUT] /path/:id', (done) => {
            chai.request(server)
            .put('/path/5b3e9560affc0b6f593c4ff3')
            .send(update)
            .end((err, res) => {
                res.should.have.status(200)
                res.body.should.be.a('object')
                done()
            })
        })
    })

    //all other requests should error
    describe('[GET] /unicorns/farts', (done) => {
        it('should fail with 404', (done) => {
            chai.request(server)
            .get('/unicorns/farts')
            .end((err, res) => {
                res.should.have.status(404)
                done()
            })
        })
    })
})