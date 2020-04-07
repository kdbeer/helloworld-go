package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// func getImageToProcess(rp *repo) ([]ImageProb, error) {
// 	processImage := make([]ImageProb, 0)

// 	images, err := rp.GetAllImages()
// 	if err != nil {
// 		return processImage, err
// 	}
// 	log.Println("Total images: ", len(images))

// 	for i, v := range images {
// 		log.Infof("Downloading [%d/%d] => %s\n", i+1, len(images), v.Name)

// 		response, e := http.Get(v.ImageURL)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 		defer response.Body.Close()

// 		if response.StatusCode != 200 {
// 			log.Fatalln("Failed: ", response.Status)
// 			continue
// 		}

// 		//open a file for writing
// 		file, err := os.Create(path.Join("./images", v.Name))
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		defer file.Close()

// 		// Use io.Copy to just dump the response body to the file. This supports huge files
// 		_, err = io.Copy(file, response.Body)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		processImage = append(processImage, v)
// 	}

// 	return processImage, nil
// }

// // ProcessVisionBatch download and adjust weight images
// func ProcessVisionBatch() error {
// 	rp, err := newConnection()
// 	if err != nil {
// 		log.Fatalln("[ProcessVisionBatch] new connection error: ", err)
// 		return err
// 	}

// 	images, err := getImageToProcess(rp)
// 	if err != nil {
// 		log.Fatalln("[ProcessVisionBatch] get image error: ", err)
// 		return err
// 	}
// 	log.Println("[ProcessVisionBatch] Images before process: ", len(images))

// 	rec, err := face.NewRecognizer("./models")
// 	if err != nil {
// 		log.Fatalln(err)
// 		return err
// 	}

// 	for _, v := range images {
// 		fullPath := path.Join("/home/faceuser/vision/images", v.Name)
// 		face, err := rec.RecognizeSingleFileCNN(fullPath)
// 		if err != nil {
// 			log.Fatalln("recognize file error: ", err)
// 			continue
// 		}

// 		log.Println("==============================")
// 		log.Println(face.Descriptor)
// 		log.Println(face.Rectangle)
// 		log.Println(face.Shapes)
// 		log.Println("==============================")

// 		if err := rp.UpdateFaceInfo(face, v.ID); err != nil {
// 			log.Fatalln("[ProcessVisionBatch] process vision batch error: ", err)
// 			return err
// 		}
// 	}

// 	log.Println(rec)

// 	return nil
// }

type response struct {
	MessageCode        string      `json:"messageCode"`
	MessageDescription string      `json:"messageDescription"`
	Data               []ImageProb `json:"data"`
}

func getImageFromAPI() ([]ImageProb, error) {
	var resp struct {
		Data []ImageProb `json:"data"`
	}

	imp := make([]ImageProb, 0)
	response, err := http.Get(apiGetImage)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return imp, err
	}

	dataByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("read data error %s\n", err)
		return imp, err
	}

	if err := json.Unmarshal(dataByte, &resp); err != nil {
		fmt.Printf("Unmarshal data error %s\n", err)
		return imp, err
	}

	imp = resp.Data

	fmt.Println(imp)
	return imp, nil
}
