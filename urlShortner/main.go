package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"shurl/urlshort"
)

const BLOCK_SIZE int = 4096

func main() {

	yaml_filename := flag.String("yaml", "rules.yaml", "A .yaml/.yml file of the format '- path: \\n url: '")
	json_filename := flag.String("json", "rules.json", "A .json file of the format [{path: val, url: val}, ...]")

	flag.Parse()

	baseHandler := defaultHandler()

	baseHandler = updateHandler(yaml_filename, &baseHandler)

	baseHandler = updateHandler(json_filename, &baseHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", baseHandler)

}

func updateHandler(filePathStr *string, fallback *http.HandlerFunc) http.HandlerFunc {

	ext := filepath.Ext(*filePathStr)
	fileBytes, fileReadErr := ioutil.ReadFile(*filePathStr)

	if fileReadErr != nil {
		exit(fmt.Sprintf("Failed to open file %s", *filePathStr))
	}

	var handlerErr error = nil
	var handler http.HandlerFunc = *fallback

	switch ext {
	case ".yaml", ".yml":
		handler, handlerErr = urlshort.YAMLHandler(fileBytes, *fallback)
	case ".json":
		handler, handlerErr = urlshort.JSONHandler(fileBytes, *fallback)
	default:
		exit(fmt.Sprintf("Invalid file type %s for url shortener configuration.", ext))
	}

	if handlerErr != nil {
		panic(handlerErr)
	}

	return handler
}

func defaultHandler() http.HandlerFunc {
	mux := defaultMux()
	pathsToUrls := map[string]string{}
	return urlshort.MapHandler(pathsToUrls, mux)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
