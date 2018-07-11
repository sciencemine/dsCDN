//A set of untility functions to reduce code repition

const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
		bodyParser = require('body-parser')
		path = require('path')
        assert = require('assert')

// constants
const mongoURL = 'mongodb://localhost:27017',
    dbName = 'ds',
    pathColName = 'paths'


function query(db, collectionName, id) {
	return new Promise((resolve, reject) => {
		let collection = db.collection(collectionName)

		collection.find({_id: id }).toArray((err, doc) => {
			if (err) {
				reject(err)
			}

			resolve(doc[0])
		})
	})
}

function queryErr(err) {
	console.error(err)
}

function add(db, collectionName, obj, opts = { }) {
	return new Promise((resolve, reject) => {
		console.log(typeof(db))
		let collection = db.collection(collectionName)

		collection.insertOne(obj, opts).then(resolve).catch((err) => { console.log(err); reject(err); })
	})
}

module.exports = add, query, queryErr