package main

import (
	"fmt"
	"net/http"
	"html/template"
	"./models"
	"crypto/rand"
)

var posts map[string]*models.Post

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("templates/index.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprint(writer, err.Error())
	}

	temp.ExecuteTemplate(writer, "index", posts)
}

func writeHandler(writer http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("templates/write.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprint(writer, err.Error())
	}
	temp.ExecuteTemplate(writer, "write", nil)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("templates/write.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprint(writer, err.Error())
	}

	id := request.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(writer, request)
	}

	temp.ExecuteTemplate(writer, "write", post)
}

func savePostHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")

	var post *models.Post
	if id == "" {
		id = generateId()
	}

	post = models.NewPost(id, title, content)
	posts[post.Id] = post

	http.Redirect(writer, request, "/", 302)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")

	if id == "" {
		http.NotFound(writer, request)
	}

	delete(posts, id)

	http.Redirect(writer, request, "/", 302)
}

func main() {
	fmt.Println("Listen on port: 3000")

	posts = make(map[string]*models.Post, 0)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/SavePost", savePostHandler)

	http.ListenAndServe(":3000", nil)
}

func generateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprint("%x", b)
}
