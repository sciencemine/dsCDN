class CE {
    constructor(title, version, playlist, options) {
        this.title = title;
        this.version = version;
        this.playlist = playlist;

        for (let option in options) {
            this[option] = options[option];
        }
    }
}

module.exports = CE;