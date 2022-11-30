package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const INDEX = `<!DOCTYPE html>
<html>
  <head>
    <title>Powered By Paketo Buildpacks</title>
  </head>
  <body>
    <h1>Powered By Paketo Buildpacks</h1>
		<img style="display: block; margin-left: auto; margin-right: auto; width: 50%;" src="https://paketo.io/images/paketo-logo-full-color.png"></img>
  </body>
</html>`

/* Commmenting out the block above and uncommenting the block below will swap the image on the homepage */
// const INDEX = `<!DOCTYPE html>
// <html>
//   <head>
//     <title>Powered By Paketo Buildpacks</title>
//   </head>
//   <body>
//     <h1>Powered By Paketo Buildpacks</h1>
//     <img style="display: block; margin-left: auto; margin-right: auto; width: 50%;" src="https://paketo.io/v2/images/buildpack-equation.svg"></img>
//   </body>
// </html>`

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, INDEX)
	})

	/* Uncommenting the line below will demo a file change that results in a new package import. */
	// fmt.Println(strings.Join([]string{"Hello", "developer"}, ", "))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
