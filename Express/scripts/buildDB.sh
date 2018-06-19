#!/bin/bash

mongo ds

#mongo ds --eval '(db.dsms.find()).pretty().shellPrint()'

mongo ds --eval 'db.createCollection( "dsms", {
	validator: {
		$jsonSchema: {
			bsonType: "object",
			required: ["id","title","version","author","config","stylesheet"],
			properties: {
				id: {
					bsonType: "string",
					description: "must be a string, required"
				},
				title: {
					bsonType: "string",
					description: "title of the dsm, required"
				},
				description: {
					bsonType: "string",
					description: "description of the dsm"
				},
				version: {
					bsonType: "string",
					description: "the current version of the DSM"
				},
				author: {
					bsonType: "string",
					description: "the author that originally created the dsm"
				},
				config: {
					bsonType: "object",
					required: [],
					properties: {
						idle: {
							bsonType: "int"
						},
						menuDwell: {
							bsonType: "int"
						},
						popoverDwell: {
							bsonType: "int"
						},
						popoverShowDelay: {
							bsonType: "int"
						}
					}
				},
				stylesheet: {
					bsonType: "string",
					description: "url to the style to use"
				}
			}
		}
	}
})'
