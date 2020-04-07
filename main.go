package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/Kagami/go-face"
)

var (
	repository *repo
)

const apiGetImage = "https://primitiveai.com/api/vision/batch/images"

func main() {
	// config
	stage := flag.String("stage", "local", "set working environment")
	flag.Parse()

	// initialize constant
	log.Print("Config database ...")
	Conf.initViper(*stage)
	log.Print("Done")

	http.HandleFunc("/", processBatchHandler)
	// http.HandleFunc("/process", processBatchHandler)

	log.Print("Hello world sample started.")

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }

	port := "8080"

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// func handler(w http.ResponseWriter, r *http.Request) {

// 	images, err := getImageToProcess(rp)
// 	if err != nil {
// 		log.Errorln("[ProcessVisionBatch] get image error: ", err)
// 		return err
// 	}
// 	log.Infoln("[ProcessVisionBatch] Images before process: ", len(images))

// 	res, err := json.Marshal(images)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(res)
// }

func processBatchHandler(w http.ResponseWriter, r *http.Request) {

	processImage := make([]ImageProb, 0)

	images, err := getImageFromAPI()
	if err != nil {
		log.Fatalln("[ProcessVisionBatch] get image error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[ProcessVisionBatch] Images before process: ", len(images))

	for i, v := range images {
		log.Printf("Downloading [%d/%d] => %s\n", i+1, len(images), v.Name)

		response, e := http.Get(v.ImageURL)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			log.Fatalln("Failed: ", response.Status)
			continue
		}

		// open a file for writing
		file, err := os.Create(path.Join("./images", v.Name))
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		processImage = append(processImage, v)
	}

	rec, err := face.NewRecognizer("./models")
	if err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("-------------------------------------------")
	for i, v := range images {
		fullPath := path.Join("./images", v.Name)
		face, err := rec.RecognizeSingleFileCNN(fullPath)
		if err != nil {
			log.Println("[ERROR] recognize file error: ", err)
			continue
		}

		fmt.Println("[Proceed] image: ", fullPath)

		dc := make([]float32, 0)
		for _, v := range face.Descriptor {
			dc = append(dc, v)
		}

		images[i].Descripter = &dc
	}

	fmt.Println("[processBatchHandler] make response")
	res, err := json.Marshal(images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
