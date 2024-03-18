package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"work/src/getssh"
)

// SshKey структура для передачи данных в шаблон
type SshKey struct {
	SshCount int
	SshList  []string
}

// viewHandler обрабатывает запросы на получение SSH ключей
func viewHandler(w http.ResponseWriter, r *http.Request) {
	testText, err := getssh.GetSSH() // Получаем SSH ключи
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("../template/view.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sshbook := SshKey{
		SshCount: len(testText),
		SshList:  testText,
	}

	if err := tmpl.Execute(w, sshbook); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// helloHandler обрабатывает запросы на /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

// indexHandler обрабатывает главную страницу
func indexHandler(w http.ResponseWriter, r *http.Request) {
	sampleData := []string{"Element 1", "Element 2", "Element 3"}

	tmpl, err := newTemplate("index.html", "../template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, sampleData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// newTemplate создает новый шаблон с измененными разделителями
func newTemplate(name string, files ...string) (*template.Template, error) {
	tmpl := template.New(name).Delims("[[", "]]")
	return tmpl.ParseFiles(files...)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/getssh", viewHandler)
	http.HandleFunc("/", indexHandler)

	staticPath := "../template/static/"
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
