package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 20)
	fmt.Println(r)
	file, handler, err := r.FormFile("myFile")
	fmt.Println(file)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// tempFile, err := ioutil.TempFile("filestore", r.FormValue("myFile"))
	// fmt.Println(tempFile.Name())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer tempFile.Close()
	tempFile, _ := os.Create("filestore/" + handler.Filename)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func fileExists(name string) bool {
	files, err := ioutil.ReadDir("./filestore")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if name == f.Name() {
			return true
		}
	}
	return false
}

func updateFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Update Endpoint Hit")

	r.ParseMultipartForm(10 << 20)
	fmt.Println(r)
	file, handler, err := r.FormFile("myFile")
	fmt.Println(file)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileExist := fileExists(handler.Filename)
	if fileExist {
		tempFile, _ := os.OpenFile("filestore/"+handler.Filename, os.O_RDWR, 0644)
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		tempFile.Write(fileBytes)
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
