const db = connect("localhost:27017/ds")

db.createCollection( "dsm", {
    validator: { $jsonschema: {
        bsonType: "object",
        required: [
            "_id",
            "title",
            "description",
            "version",
            "author",
            "config",
            "stylesheet",
            "style",
            "contributors",
            "idle_backgrounds",
            "video_select_backgrounds",
            "ce_set",
            "attributes"
        ],
        properties: {
            _id: {
                bsonType: "string",
                description: "must be a string, required"
            },
            title: {
                bsonType: "string",
                description: "title of the dsm, required"
            },
            description: {
                bsonType: "string",
                description: "the description of this dsm"
            },
            version: {
                bsonType: "string",
                description: "the description of this dsm"
            },
            author: {
                bsonType: "string",
                description: "the original creator of the dsm"
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
            },
            style: {
                bsonType: "object"
            },
            contributors: {
                bsonType: "array",
                description: "a list of all the authors that worked on the dsm"
            },
            idle_backgrounds: {
                bsonType: "array",
                description: "a list of ces to be played in idle mode"
            },
            video_select_background: {
                bsonType: "array",
                description: "a list of ces that are to be rendered in video"
            },
            ce_set: {
                bsonType: "object",
                description: "the set of CEs in the semantic web of the dsm"
            },
            attributes: {
                bsonType: "array",
                description: "the list of attributes the dsm has"
            }
        }
    }}
})

db.createCollection( "ces", {
    validator: { $jsonschema: {
        bsonType: "object",
        required: [
            "_id",
            "version",
            "title",
            "playlist",
        ],
        properties: {
            _id: {
                bsonType: "string",
                description: "identifies the ce"
            },
            version: {
                bsonType: "string",
                description: "the current version of the ce"
            },
            description: {
                bsonType: "string",
                description: "is optional"
            },
            title: {
                bsonType: "string"
            },
            playlist: {
                bsonType: "array",
                description: "an array of assets to play, optional teaser asset at index 0"
            }
        }

    }}
})

db.createCollection("assets", {
    validator: { $jsonschema: {
        bsonType: "object:",
        required: [
            "_id",
            "version",
            "url",
            "type",
            "options"
        ],
        properties: {
            _id: {
                bsonType: "string",
                description: "identifies the asset"
            },
            version: {
                bsonType: "string"
            },
            url: {
                bsonType: "string"
            },
            type: {
                bsonType: "string",
                description: "mime type"
            },
            options: {
                bsonType: "object",
                properties: {
                    start: {
                        bsonType: "string"
                    },
                    end: {
                        bsonType: "string"
                    },
                    duration: {
                        bsonType: "string"
                    }
                }
            }
        }
    }}
})