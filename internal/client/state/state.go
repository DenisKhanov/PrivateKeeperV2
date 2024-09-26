package state

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

const (
	perm = 0o755
)

type ClientState struct {
	token        string
	isAuthorized bool
	login        string
	dirPath      string
}

func NewClientState() *ClientState {
	return &ClientState{}
}

func (c *ClientState) IsAuthorized() bool {
	return c.isAuthorized
}

func (c *ClientState) SetIsAuthorized(isAuthorized bool) {
	c.isAuthorized = isAuthorized
}

func (c *ClientState) GetToken() string {
	return c.token
}

func (c *ClientState) SetToken(token string) {
	c.token = token
}

func (c *ClientState) GetLogin() string {
	return c.login
}

func (c *ClientState) SetLogin(login string) {
	c.login = login
}

func (c *ClientState) SetWorkingDirectory() {
	// Создаем новый сканер для чтения ввода с консоли
	scanner := bufio.NewScanner(os.Stdin)

	// Настройка для цветного текста
	yellowBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(yellowBold("Write path to your working directory (will be created if it doesn't exist)"))

	// Запрос ввода пути
	fmt.Printf("Input path to working directory: ")
	scanner.Scan() // Считываем введенный текст
	path := scanner.Text()

	// Преобразуем слеши в соответствии с текущей операционной системой
	path = filepath.FromSlash(path)

	// Проверяем, существует ли каталог
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Создаем все необходимые директории по пути, если они не существуют
		err = os.MkdirAll(path, 0755) // 0755 - права доступа к каталогу
		if err != nil {
			fmt.Println("Error creating working directory, please try again")
			return // Возвращаемся, если ошибка создания
		}
	}

	// Устанавливаем путь в структуру
	c.dirPath = path
	fmt.Println("Working directory set to:", path)
}

func (c *ClientState) GetDirPath() string {
	return c.dirPath
}
