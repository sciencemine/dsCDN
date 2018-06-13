class DSM {
    constructor(id, title, description, version, author, config, stylesheet, style,
            contributors, idle_backgrounds, video_select_backgrounds, ce_set,
            attributes) {
        this._id = id;
        this.title = title;
        this.description = description;
        this.version = version;
        this.author = author;
        this.config = config;
        this.stylesheet = stylesheet;
        this.style = style;
        this.contributors = contributors;
        this.idle_backgrounds = idle_backgrounds;
        this.video_select_backgrounds = video_select_backgrounds;
        this.ce_set = ce_set;
        this.attributes = attributes;
    }
}

module.exports = DSM;