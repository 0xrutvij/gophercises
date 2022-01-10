package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cyoa"
)

func main() {

	port := flag.Int("port", 3000, "The port to start the server CYOA web application on.")

	filename := flag.String("file", "gopher.json", "the JSON file for CYOA story to use.")
	flag.Parse()

	fmt.Printf("Using the story in file %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
