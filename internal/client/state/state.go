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
	scanner := bufio.NewScanner(os.Stdin)

	yellowBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(yellowBold("Write path to your working directory (will be created if it doesn't exist)"))

	fmt.Printf("Input path to working directory : ")
	scanner.Scan()
	path := scanner.Text()

	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)
	fmt.Println(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, perm)
		if err != nil {
			fmt.Println("Error creating working directory, please try again")
		}
	}

	c.dirPath = dir
}

func (c *ClientState) GetDirPath() string {
	return c.dirPath
}
