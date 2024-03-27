package main

import (
	"fmt"
	"html/template"
	"net/http"
	"work/src/comand"
	"work/src/getssh"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
		if string(msg) == "get_journal" {
			journal, err := comand.GetJournal()
			if err != nil {
				// Обрабатываем ошибку
				fmt.Println("Ошибка при получении журнала:", err)
				continue
			}
			// Отправляем данные журнала клиенту
			if err := conn.WriteMessage(websocket.TextMessage, []byte(journal)); err != nil {
				// Обрабатываем ошибку
				fmt.Println("Ошибка при отправке журнала:", err)
				continue
			}
		} else {
			conn.WriteMessage(t, msg)
		}
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

var avtorization = false

// indexHandler Выводит основную страницу сайта на экран.
func indexHandler(c *gin.Context) {
	if !avtorization {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		c.HTML(http.StatusOK, "index.html", nil)
		avtorization = false
	}
	// c.HTML(http.StatusOK, "index.html", nil)
}

func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func log(c *gin.Context) {
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

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Определяем обработчики
	r.GET("/getssh", getSSHHandler)
	r.GET("/", indexHandler)
	r.GET("/login", loginHandler)
	r.POST("/log", log)
	// Обработчик API
	r.GET("/api/data", apiHandler)
	r.POST("/api/data", apiPost)
	r.GET("/ws", wshandler)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"user1": "love",
		"user2": "god",
		"user3": "sex",
	}))

	authorized.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"secret": "The secret ingredient to the BBQ sauce is stiring it in an old whiskey barrel.",
		})
	})

	// Запускаем сервер
	if err := r.Run(":64000"); err != nil {
		panic(err)
	}
}
