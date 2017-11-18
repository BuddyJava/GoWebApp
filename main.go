package main

import (
	"fmt"
	"net/http"
	"./models"
	"crypto/rand"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var posts map[string]*models.Post

func indexHandler(rnd render.Render) {
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML(200, "write", nil)
}

func editHandler(rnd render.Render, writer http.ResponseWriter, request *http.Request, params martini.Params) {
	id := params["id"]
	post, found := posts[id]
	if !found {
		http.NotFound(writer, request)
	}
	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, request *http.Request) {
	id := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")

	var post *models.Post
	if id == "" {
		id = generateId()
	}

	post = models.NewPost(id, title, content)
	posts[post.Id] = post

	rnd.Redirect("/")
}

func deleteHandler(rnd render.Render, params martini.Params) {
	id := params["id"]

	if id == "" {
		rnd.Redirect("/")
	}

	delete(posts, id)
	rnd.Redirect("/")
}

func main() {
	fmt.Println("Listen on port: 3000")

	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		//Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,    // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	m.Run()
}

func generateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprint("%x", b)
}
