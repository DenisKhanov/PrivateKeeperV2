package state

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

const (
	perm = 0o755 // Default permission for created directories (rwxr-xr-x)
)

// ClientState holds the state information of the client,
// including authorization status, token, login, and working directory.
type ClientState struct {
	token        string // Token for authorized access
	isAuthorized bool   // Flag to indicate if the user is authorized
	login        string // User login
	dirPath      string // Path to the working directory
}

// NewClientState creates and returns a new instance of ClientState.
func NewClientState() *ClientState {
	return &ClientState{}
}

// IsAuthorized returns whether the client is authorized.
func (c *ClientState) IsAuthorized() bool {
	return c.isAuthorized
}

// SetIsAuthorized sets the authorization status of the client.
func (c *ClientState) SetIsAuthorized(isAuthorized bool) {
	c.isAuthorized = isAuthorized
}

// GetToken retrieves the current token of the client.
func (c *ClientState) GetToken() string {
	return c.token
}

// SetToken sets the token for the client.
func (c *ClientState) SetToken(token string) {
	c.token = token
}

// GetLogin retrieves the current login of the client.
func (c *ClientState) GetLogin() string {
	return c.login
}

// SetLogin sets the login for the client.
func (c *ClientState) SetLogin(login string) {
	c.login = login
}

// SetWorkingDirectory prompts the user to enter a path for the working directory.
// It creates the directory if it doesn't exist.
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
		err = os.MkdirAll(path, perm) // 0755 - права доступа к каталогу
		if err != nil {
			fmt.Println("Error creating working directory, please try again")
			return // Возвращаемся, если ошибка создания
		}
	}

	// Устанавливаем путь в структуру
	c.dirPath = path
	fmt.Println("Working directory set to:", path)
}

// GetDirPath retrieves the current working directory path.
func (c *ClientState) GetDirPath() string {
	return c.dirPath
}
