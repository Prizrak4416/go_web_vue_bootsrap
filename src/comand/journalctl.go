package comand

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetJournal() (string, error) {
	// создаем команду
	cmd := exec.Command("journalctl", "-e")
	// Запускае команду, получаем ее вывод
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды:", err)
		return "", err
	}
	outputStr := string(output)
	outputStr = strings.ReplaceAll(outputStr, "\n", "<br>")
	return outputStr, err
}
