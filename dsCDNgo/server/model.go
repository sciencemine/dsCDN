package server

//The model is where we define the data types that we work with and store in the database

import (
	"gopkg.in/mgo.v2/bson"
)

//DSM defines what a digital signage model (dsm) is for an exhibit
type DSM struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Version     string        `json:"version" bson:"version"`
	Author      string        `json:"author" bson:"author"`
	Config      struct {
		Idle             float32 `json:"idle" bson:"idle"`
		MenuDwell        float32 `json:"menuDwell" bson:"menuDwell"`
		PopoverDwell     float32 `json:"popoverDwell" bson:"popoverDwell"`
		PopoverShowDelay float32 `json:"popoverShowDelay" bson:"popovershowDelay"`
	} `json:"config" bson:"config"`
	Tactile struct {
		Select   []string `json:"select" bson:"select"`
		Previous []string `json:"previous" bson:"previous"`
		Cancel   []string `json:"cancel" bson:"cancel"`
		Next     []string `json:"next" bson:"next"`
	} `json:"tactile" bson:"tactile"`
	Stylesheet string `json:"stylesheet" bson:"stylesheet"`
	Style      struct {
		Idle    string `json:"idle" bson:"idle"`
		Menu    string `json:"menu" bson:"menu"`
		Select  string `json:"select" bson:"select"`
		Playing string `json:"playing" bson:"playing"`
	} `json:"style" bson:"style"`
	Contributors           []string `json:"contributors" bson:"contributors"`
	IdleBackgrounds        []string `json:"idle_backgrounds" bson:"idle_backgrounds"`
	VideoSelectBackgrounds []string `json:"video_select_backgrounds" bson:"video_select_backgrounds"`
	CEset                  map[string]struct {
		ID            string `json:"id" bson:"id"`
		Attributes    []int  `json:"attributes" bson:"attributes"`
		Relationships []struct {
			From         string `json:"from" bson:"from"`
			To           string `json:"to" bson:"to"`
			Weight       int    `json:"weight" bson:"weight"`
			AttributeIdx int    `json:"attribute" bson:"attribute"`
		} `json:"relationships" bson:"relationships"`
	} `json:"ce_set" bson:"ce_set"`
	Attributes []struct {
		Title       string `json:"title" bson:"title"`
		Description string `json:"description" bson:"description"`
		Image       string `json:"image" bson:"image"`
	} `json:"attributes" bson:"attributes"`
}

//CE defines one of the Content Element nodes in the semantic web
type CE struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	Version     string        `json:"version" bson:"version"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Playlist    struct {
		Teaser interface{} `json:"teaser,omitempty" bson:"teaser"`
		Queue  []struct {
			Primary     interface{} `json:"primary" bson:"primary"`
			Backgrounds []desc      `json:"backgrounds" bson:"backgrounds"`
			Tracks      []desc      `json:"tracks" bson:"tracks"`
			Overlays    []desc      `json:"overlays" bson:"overlays"`
		} `json:"queue" bson:"queue"`
	} `json:"playlist" bson:"playlist"`
}

type desc struct {
	Asset    interface{} `json:"asset" bson:"asset"`
	X        float64     `json:"x" bson:"x"`
	Y        float64     `json:"y" bson:"y"`
	Start    float64     `json:"start" bson:"start"`
	Duration float64     `json:"duration" bson:"duration"`
	Width    float64     `json:"width" bson:"width"`
	Height   float64     `json:"height" bson:"height"`
}

//Asset defines what an asset should be
type Asset struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	Version     string        `json:"version" bson:"version"`
	Attribution string        `json:"attribution" bson:"attribution"`
	URL         string        `json:"url" bson:"url"`
	MIMEType    string        `json:"type" bson:"type"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Options     struct {
		Start    string `json:"start" bson:"start"`
		End      string `json:"end" bson:"end"`
		Duration string `json:"duration" bson:"duration"`
	} `json:"options"`
}

//PATH defines a path in the semantic web
type PATH struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Model struct {
		ModelID          string `json:"id" bson:"id"`
		ModelVersion     string `json:"version_model" bson:"version_model"`
		ModelDescription string `json:"description_model" bson:"description_model"`
		Author           string `json:"author" bson:"author"`
	} `json:"model" bson:"model"`
	Relations []Relation `json:"relations" bson:"relations"`
}

//Relation describes a edge on the semantic web between two ce nodes
type Relation struct {
	AtrTitle       string `json:"title_attr" bson:"title_attr"`
	AtrDescription string `json:"description_attr" bson:"description_attr"`
	Weight         int    `json:"weight" bson:"weight"`
	CEList         []struct {
		CEID          string `json:"id_ce" bson:"id_ce"`
		CEVersion     string `json:"version_ce" bson:"version_ce"`
		CETitle       string `json:"title_ce" bson:"title_ce"`
		CEDescription string `json:"description_ce" bson:"description_ce"`
		Events        []struct {
			AssetID          string  `json:"id_asset" bson:"id_asset"`
			AssetVersion     string  `json:"version_asset" bson:"version_asset"`
			AssetTitle       string  `json:"title_asset" bson:"title_asset"`
			AssetDescription string  `json:"description_asset" bson:"description_asset"`
			Type             string  `json:"type" bson:"type"`
			URL              string  `json:"url" bson:"url"`
			Playspec         float32 `json:"playspec" bson:"playspec"`
		} `json:"events" bson:"events"`
	} `json:"ce_list" bson:"ce_list"`
}

//Dsms is an array of dsm ids
type Dsms []struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`
}

//Ces is an array of ce ids
type Ces []struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`
}

//Paths is an array of path ids
type Paths []struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`
}

//Assets is an array of Asset ids
type Assets []struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`
}
