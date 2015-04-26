package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	UPLOAD_DIR = "./files"
)

const ()

func checkError(str string, err error) {
	if err != nil {
		log.Fatal("ListenAndServen: ", err.Error())
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	pageHead := `
	<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8">
	<title>Upload</title>
	</head>
	<body>
	`
	pageForm := `
	<form method="POST" action="/upload" enctype="multipart/form-data">
	Choose a file to upload: <input name="file" type="file">
	<input type="submit" value="Upload">
	</form>
	`
	pageFoot := `
	</body>
	</html>
	`

	fileInfoArr, err := ioutil.ReadDir("./files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var listHtml string
	for _, fileInfo := range fileInfoArr {
		fileId := fileInfo.Name()
		listHtml += "<li><a href=\"/files/" + fileId + "\">" + fileId + "</a></li>"
	}
	listHtml = "<ol>" + listHtml + "</ol>"

	if r.Method == "GET" {
		io.WriteString(w, pageHead+pageForm+listHtml+pageFoot)
		return
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

const (
	ListDir = 0x0001
)

func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & ListDir) == 0 {
			if exists := isExists(file); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
	})
}
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/files/", "./files", 0)
	mux.HandleFunc("/", uploadHandler)
	err := http.ListenAndServe(":80", mux)
	checkError("ListenAndServe: ", err)
}
