package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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

	storyTpl := template.Must(template.New("").Parse(storyTempl))
	hc := cyoa.NewHandler(story, cyoa.WithTemplate(storyTpl), cyoa.WithPathFunc(pathFn))
	mux := http.NewServeMux()
	//Custom : This story handler is using a custom function and template
	// Because we use /story/ (trailing slash) all web requests
	// whose path has the /story/ prefix will be routed here.
	mux.Handle("/story/", hc)
	//Default:  This story handler is using the default functions and templates
	// Because we use / (base path) all incoming requests not
	// mapped elsewhere will be sent here.
	mux.Handle("/", cyoa.NewHandler(story))
	fmt.Println("Starting the server on port \n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
	// fmt.Println("%+v\n", story)

}

// Custom paths and template
func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	path = path[len("/story/"):]
	return path
}

var storyTempl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p> {{.}} Some para</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section> 
    <style>
    body {
      font-family: helvetica, arial;
    }
    h1 {
      text-align:center;
      position:relative;
    }
    .page {
      width: 80%;
      max-width: 500px;
      margin: auto;
      margin-top: 40px;
      margin-bottom: 40px;
      padding: 80px;
      background: #FFFCF6;
      border: 1px solid #eee;
      box-shadow: 0 10px 6px -6px #777;
    }
    ul {
      border-top: 1px dotted #ccc;
      padding: 10px 0 0 0;
      -webkit-padding-start: 0;
    }
    li {
      padding-top: 10px;
    }
    a,
    a:visited {
      text-decoration: none;
      color: #6295b5;
    }
    a:active,
    a:hover {
      color: #7792a2;
    }
    p {
      text-indent: 1em;
    }
  </style>  
  </body>
</html>`
