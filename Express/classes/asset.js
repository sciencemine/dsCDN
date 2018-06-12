class Asset {
    constructor(version, url, type, options = { }) {
        this.version = version;
        this.url = url;
        this.type = type;

        for (let option in options) {
            this[option] = options[option];
        }
    }
}

module.exports = Asset;