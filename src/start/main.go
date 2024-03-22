package main

import (
	"html/template"
	"net/http"
	"work/src/comand"
	"work/src/getssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

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

type user struct {
	ID       string
	Name     string
	UserName string
}

var users = []user{
	{ID: "1", Name: "Yura", UserName: "hj"},
	{ID: "2", Name: "Vasia", UserName: "vs"},
	{ID: "3", Name: "Lida", UserName: "ld"},
}

// apiHandler отправка API
func apiHandler(c *gin.Context) {
	// Отправка JSON-ответа
	c.JSON(http.StatusOK, gin.H{
		"message": users,
	})
}

type RequestData struct {
	Role    string `json:"role"`
	Message string `json:"message_post"`
}

func apiPost(c *gin.Context) {
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseMessage := "Получено сообщение: " + requestData.Message
	switch requestData.Role {
	case "user":
		// Логика для пользователя
		responseMessage += " (Обработано как пользователь)"
		c.JSON(http.StatusOK, gin.H{"response": responseMessage})
	case "admin":
		// Логика для администратора
		responseMessage += " (Обработано как администратор)"
		c.JSON(http.StatusOK, gin.H{"response": responseMessage})
	case "getut":
		t, err := comand.GetUptime()
		if err != nil {
			responseMessage += " (Ошибка выполнеия команды uptime) " + err.Error()
		} else {
			responseMessage += " время работы linux " + t
		}
		c.JSON(http.StatusOK, gin.H{"response": responseMessage})
	case "getssh":
		testText, err := getssh.GetSSH() // Получаем SSH ключи
		if err != nil {
			responseMessage += " (Получения ssh ключей) " + err.Error()
			c.JSON(http.StatusOK, gin.H{"response": responseMessage})
		} else {
			c.JSON(http.StatusOK, gin.H{"response": testText})
		}
	default:
		responseMessage = "Неизвестная роль"
		c.JSON(http.StatusOK, gin.H{"response": responseMessage})
	}
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
	// Обработчик API
	r.GET("/api/data", apiHandler)
	r.POST("/api/data", apiPost)
	r.GET("/ws", wshandler)

	// Запускаем сервер
	if err := r.Run(":64000"); err != nil {
		panic(err)
	}
}
