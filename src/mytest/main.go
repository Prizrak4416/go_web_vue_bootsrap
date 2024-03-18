package main

import (
	"fmt"
	"log"
	"os/user"
)

func main() {
	// Получаем данные текущего пользователя.
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Выводим папку пользователя.
	fmt.Printf("Домашняя папка пользователя %s является %s\n", usr.Username, usr.HomeDir)
}
