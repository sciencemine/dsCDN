//A set of untility functions to reduce code repition

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

module.exports = { add, queryErr, query}