package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shriya/cyoa"
)

func main() {
	fmt.Println("Hello ")
	port := flag.Int("port", 3000, "the port to start the CYOA web application on ")
	filename := flag.String("file", "../../gopher.json", "the JSON file with the cyoa story")
	flag.Parse()
	fmt.Printf("Using the story in %s \n", *filename)

	jsonFile, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Error occured ")
		panic(err)
	}

	story, err := cyoa.JsonStory(jsonFile)
	if err != nil {
		panic(err)
	}

	// If user need to pass it own template it can do by provide it in newHandler func
	h := cyoa.NewHandler(story, nil)
	fmt.Println("Starting the server on port \n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
	// fmt.Println("%+v\n", story)

}
