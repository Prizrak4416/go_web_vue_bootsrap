package main

import (
	// "fmt"
	"html/template"
	"net/http"
	"work/src/comand"
	"work/src/getssh"
	"work/src/wbsocket"

	"github.com/gin-gonic/gin"
)

var avtorization = false

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

// indexHandler Выводит основную страницу сайта на экран.
func indexHandler(c *gin.Context) {
	if !avtorization {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		c.HTML(http.StatusOK, "index.html", nil)
		avtorization = false
	}
}

// loginHandler Страница для отображения авторизации и проверки пароля.
func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
  }
  
func log_pass(c *gin.Context) {
	var userInput struct {
		Password string `json:"password"`
	}
	pass := "123456"
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	if pass == userInput.Password {
		avtorization = true
		c.JSON(http.StatusOK, gin.H{"response": "Авторизация успешна"})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Необходимо ввести учетные данные"})
	}
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

// RequestData структура для API запросов
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
	r.GET("/", indexHandler)
	r.GET("/login", loginHandler)
 	r.POST("/log", log_pass)
	// Обработчик API
	r.GET("/api/data", apiHandler)
	r.POST("/api/data", apiPost)
	wsHandler := wbsocket.NewWebSocketHandler()
	r.GET("/ws", wsHandler.Handle)

	// Запускаем сервер
	if err := r.Run(":64000"); err != nil {
		panic(err)
	}
}
