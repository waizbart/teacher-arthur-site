package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const PORT = 8080

const FILES_DIR = "./files"

func main() {
	addr := fmt.Sprintf("0.0.0.0:%d", PORT)
	fmt.Printf("Running the HTML server on %s\n", addr)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/{id}", fileHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(addr, r))
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	fileBytes := readFile("./static/index.html")
	if fileBytes == nil {
		w.Write([]byte("File not found :("))
	}
	fileContent := string(fileBytes)
	var finalFileListHTML strings.Builder

	filesList := listFiles()

	for _, fileName := range filesList {
		noHtml := strings.ReplaceAll(fileName, ".html", "")
		noDash := strings.ReplaceAll(noHtml, "-", " ")
		noUnderscore := strings.ReplaceAll(noDash, "_", " ")
		title := strings.Title(noUnderscore)
		fmt.Fprintf(&finalFileListHTML, "<li><a href=\"/%s\"><span class=\"file-id\">%s</span></a></li>", fileName, title)
	}

	fileContentWithFilesList := strings.Replace(fileContent, "{% files_list %}", finalFileListHTML.String(), 1)

	w.Write([]byte(fileContentWithFilesList))
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileBytes := readFile(FILES_DIR + "/" + id)
	if fileBytes == nil {
		w.Write([]byte("File not found :("))
	}
	w.Write(fileBytes)
}

func listFiles() []string {
	entries, err := os.ReadDir(FILES_DIR)
    if err != nil {
        log.Fatal(err)
    }

	var fileList []string
 
    for _, e := range entries {
        fileList = append(fileList, e.Name())
    }

	return fileList
}

func readFile(filename string) []byte {
	dataByte, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	return dataByte
}
