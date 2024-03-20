package comand

import (
	"fmt"
	"os/exec"
)

func GetUptime() (string, error) {
	// создаем команду
	cmd := exec.Command("uptime", "-p")
	// Запускае команду, получаем ее вывод
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды:", err)
		return "", err
	}
	return string(output), err
}
