class PATH {
    constructor(id, model, relations) {
        this.id = id
        this.model = model
        this.relations = relations
    }
}

class MODEL {
    constructor(id, version, description, author) {
        this.id = id
        this.version = version
        this.description = description
        this.author = author
    }
}

class RELATIONS {
    constructor(title, desc, weight, ce_list) {
        this.title = title 
        this.desc = desc
        this.weight = weight 
        this.ce_list = ce_list
    }
}

module.exports = PATH, MODEL, RELATIONS