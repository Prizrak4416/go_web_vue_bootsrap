package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

// check Обработка ошибок
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// viewHandler читает записи из файла и выводит их
func viewHandler(write http.ResponseWriter, request *http.Request) {
	testText := getStrings("../template/test.txt")
	html, err := template.ParseFiles("../template/view.html")
	check(err)
	guestbook := Guestbook{
		SignatureCount: len(testText),
		Signatures:     testText,
	}
	err = html.Execute(write, guestbook)
	check(err)
}

// getStrings возвращает сегмент строк, прочитанный из fileName
func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	// если будет получена ошибка что файл не существует
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	sampleData := []string{"Element 1", "Element 2", "Element 3"}
	tmpl, err := newTemplate("index.html", "../template/index.html")
	// tmpl, err := template.ParseFiles("../template/index2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, sampleData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// newTemplate меняет шаблон кавычек в переданом шаблоне [[]]
func newTemplate(name string, files ...string) (*template.Template, error) {
	tmpl := template.New(name).Delims("[[", "]]")
	return tmpl.ParseFiles(files...)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/", indexHandler)

	// Указание папки template как источника статических файлов
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// fs := http.FileServer(http.Dir("../template"))
	// http.Handle("/", http.StripPrefix("/", fs))

	// Обслуживание статических файлов
	staticPath := "../template/static/"
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
