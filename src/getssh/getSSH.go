package getssh

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
)

// getPuth создает путь строкой до публичных ключей authorized_keys
func getPuth() (string, error) {
	// Получаем данные текущего пользователя.
	usr, err := user.Current()
	if err != nil {
		return "", err // Возвращаем ошибку вместо завершения программы
	}
	return usr.HomeDir + "/.ssh/authorized_keys", nil
}

// sshString получает сам ключ между типом и названием пользователя
func sshString(arr []string) ([]string, error) {
	var lines []string
	for _, a := range arr {
		parts := strings.Split(a, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("некорректная строка: %s", a)
		}
		lines = append(lines, parts[1])
	}
	return lines, nil
}

// GetSSH возвращает сегмент строк, прочитанный из authorized_keys
func GetSSH() ([]string, error) {
	filePath, err := getPuth()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	// если будет получена ошибка что файл не существует
	if os.IsNotExist(err) {
		fmt.Printf("Файла %v не существует\n", filePath)
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	ssh, err := sshString(lines)
	if err != nil {
		return nil, err
	}

	return ssh, nil
}
