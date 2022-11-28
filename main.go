package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func list() {
	files, err := ioutil.ReadDir("./filestore")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
func listAllFiles() {
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "ls") {
			list()
		}
	}
}
func removeFiles() {
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "rm") {
			for i, f := range os.Args {
				if i > 1 {
					os.Remove("./filestore/" + f)
				}
			}
		}
	}
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func addFiles() {
	urlStr := "http://localhost:8080/upload"
	isExist := false
	files, err := ioutil.ReadDir("./filestore")
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "add") {
			for i, f := range os.Args {
				if i > 1 {
					client := &http.Client{}
					for _, file := range files {
						trimmedName := strings.TrimLeft(f, ".")
						trimmedName = strings.TrimLeft(trimmedName, "\\")
						if file.Name() == f || file.Name() == trimmedName {
							isExist = true
							break
						}
					}
					if !isExist {
						r, err := newfileUploadRequest(urlStr, nil, "myFile", f)
						if err != nil {
							fmt.Println("Error encountered")
						}
						fmt.Println("--------------------------------")
						resp, _ := client.Do(r)
						fmt.Println(resp)
					} else {
						fmt.Println("File already exist with the same name: ", f)
						isExist = false
					}
				}
			}
		}
	}
}

func updateFileContent() {
	urlStr := "http://localhost:8080/upload"
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "update") {
			for i, f := range os.Args {
				if i > 1 {
					client := &http.Client{}

					r, err := newfileUploadRequest(urlStr, nil, "myFile", f)
					r.Method = "PUT"
					if err != nil {
						fmt.Println("Error encountered")
					}
					fmt.Println("--------------------------------")
					resp, _ := client.Do(r)
					fmt.Println(resp)

				}
			}
		}
	}
}

func wc() {
	count := 0
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "wc") {
			files, err := ioutil.ReadDir("./filestore")
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				file, err := os.Open("./filestore" + f.Name())
				if err != nil {
					log.Fatal(err)
				}
				Scanner := bufio.NewScanner(file)
				Scanner.Split(bufio.ScanWords)
				for Scanner.Scan() {
					count++
				}
				if err := Scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
			fmt.Println("Total words ", count)
		}
	}
}

func freq() {
	wordsCount := make(map[string]int)
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "freq-words") {
			files, err := ioutil.ReadDir("./filestore")
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				file, err := os.Open("./filestore/" + f.Name())
				if err != nil {
					log.Fatal(err)
				}
				Scanner := bufio.NewScanner(file)
				Scanner.Split(bufio.ScanWords)
				for Scanner.Scan() {
					wordsCount[Scanner.Text()]++
				}
				if err := Scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
			keys := make([]string, 0, len(wordsCount))

			for key := range wordsCount {

				keys = append(keys, key)
			}
			limit := 0
			if strings.EqualFold(os.Args[2], "--limit") {
				limit, err = strconv.Atoi(os.Args[3])
				if err != nil {
					fmt.Println("Error occurred")
				}
			}
			isDescending := true
			if strings.EqualFold(os.Args[4], "--order=dsc") {
				isDescending = true
			} else {
				isDescending = false
			}
			if isDescending {
				sort.Slice(keys, func(i, j int) bool {

					return wordsCount[keys[i]] > wordsCount[keys[j]]
				})
			} else {
				sort.Slice(keys, func(i, j int) bool {

					return wordsCount[keys[i]] < wordsCount[keys[j]]
				})
			}
			for idx, key := range keys {

				fmt.Printf("%s %d\n", key, wordsCount[key])

				if idx == limit {
					break
				}
			}
		}
	}
}

func main() {
	if err := os.Mkdir("filestore", os.ModePerm); err != nil {
		fmt.Println("Store Already exists")
	}
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/update", updateFile)
	server := &http.Server{Addr: ":8080"}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		addFiles()
		listAllFiles()
		removeFiles()
		updateFileContent()
		wc()
		freq()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}
