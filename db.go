package main

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	databaseName      = "vision_recognizer"
	ClassesCollection = "classes"
	ImagesCollection  = "images"
)

type repo struct {
	Vision *mgo.Database
}

func openConnection() (*repo, error) {

	schema := &mgo.DialInfo{
		Addrs:    []string{Conf.Mongo.Address},
		Timeout:  Conf.Mongo.Timeout,
		Database: Conf.Mongo.Database,
		Username: Conf.Mongo.Username,
		Password: Conf.Mongo.Password,
	}

	ms, err := mgo.DialWithInfo(schema)
	if err != nil {
		return nil, err
	}

	ms.SetMode(mgo.Monotonic, true)
	return &repo{Vision: ms.DB(databaseName)}, nil
}

// ImageProb class id
type ImageProb struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	ClassID    int32         `json:"class_id" bson:"class_id"`
	Name       string        `json:"name" bson:"name"`
	ImageURL   string        `json:"image_url" bson:"image_url"`
	Path       string        `json:"path" bson:"path"`
	Descripter *[]float32    `json:"descripter,omitempty" bson:"descripter,omitempty"`
	ProcessAt  *time.Time    `json:"process_at,omitempty" bson:"process_at,omitempty"`
}

// GetAllImages get image by class
func (r *repo) GetAllImages() ([]ImageProb, error) {
	data := make([]ImageProb, 0)

	selector := bson.M{"process_at": bson.M{"$exists": false}}
	if err := r.Vision.C(ImagesCollection).Find(selector).All(&data); err != nil {
		return data, err
	}

	return data, nil
}

// UpdateFaceInfo update face info
// func (r *repo) UpdateFaceInfo(info *face.Face, id bson.ObjectId) error {
// 	selector := bson.M{"_id": id}
// 	updater := bson.M{"$set": bson.M{"descripter": info.Descriptor, "process_at": time.Now()}}
// 	return r.Vision.C(ImagesCollection).Update(selector, updater)
// }
