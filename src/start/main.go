package main

import (
	"html/template"
	"net/http"
	"work/src/getssh"

	"github.com/gin-gonic/gin"
)

// SshKey структура для передачи данных в шаблон
type SshKey struct {
	SshCount int
	SshList  []string
}

// loadTemplates загружает и парсит шаблоны с измененными разделителями
func loadTemplates(templatesDir string) (*template.Template, error) {
	tmpl := template.New("").Delims("[[", "]]") // Задаем новые разделители
	// Загружаем шаблоны
	return tmpl.ParseGlob(templatesDir + "/*.html")
}

// getSSHHandler Получаем ssh ключи и выводим их а экран
func getSSHHandler(c *gin.Context) {
	testText, err := getssh.GetSSH() // Получаем SSH ключи
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sshbook := SshKey{
		SshCount: len(testText),
		SshList:  testText,
	}

	c.HTML(http.StatusOK, "view.html", sshbook)
}

// indexHandler Выводит основную страницу сайта на экран.
func indexHandler(c *gin.Context) {
	sampleData := []string{"Element 1", "Element 2", "Element 3"}
	c.HTML(http.StatusOK, "index.html", sampleData)
}

func main() {
	r := gin.Default()
	r.Static("/template/", "../template/")
	// Загружаем и парсим шаблоны
	tmpl, err := loadTemplates("../template")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl) // Устанавливаем шаблоны для Gin

	// Определяем обработчики
	r.GET("/getssh", getSSHHandler)
	r.GET("/", indexHandler)

	// Запускаем сервер
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
