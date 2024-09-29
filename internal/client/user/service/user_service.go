package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
	"github.com/fatih/color"
	"os"
	"strings"
)

// UserService is an interface that defines methods for user registration and login.
// Implementations of this interface should provide the actual functionality.
type UserService interface {
	RegisterUser(ctx context.Context, login, password string) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

// UserProvider is a struct that provides user-related functionalities.
// It contains a reference to a UserService and a ClientState.
type UserProvider struct {
	userService UserService        // The user service implementation
	state       *state.ClientState // Client state management
}

// NewUserService initializes a new UserProvider with the given user service and client state.
// It returns a pointer to the newly created UserProvider.
func NewUserService(u UserService, state *state.ClientState) *UserProvider {
	return &UserProvider{
		userService: u,
		state:       state,
	}
}

// RegisterUser prompts the user for their login credentials to register.
// It handles input validation and calls the user service to perform the registration.
func (u *UserProvider) RegisterUser(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	red := color.New(color.FgRed).SprintFunc()

	var login, password string

	yellowBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(yellowBold("Input 'login password' to register:"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input login as %s: ", yellow("'valid email'"))
	scanner.Scan()
	data := scanner.Text()
	login = strings.TrimSpace(data)
	if len(login) == 0 {
		fmt.Println(red("Login must not be empty please try again"))
		return
	}

	fmt.Printf("Input password as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	password = strings.TrimSpace(data)
	if len(password) == 0 {
		fmt.Println(red("Password must not be empty please try again"))
		return
	}

	token, err := u.userService.RegisterUser(ctx, login, password)
	if err != nil {
		lib.UnpackGRPCError(err)
	} else {
		u.state.SetToken(token)
		u.state.SetIsAuthorized(true)
		u.state.SetLogin(login)
	}
}

// LoginUser prompts the user for their login credentials to log in.
// It handles input validation and calls the user service to perform the login.
func (u *UserProvider) LoginUser(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	red := color.New(color.FgRed).SprintFunc()

	var login, password string

	yellowBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(yellowBold("Input 'login password' to login:"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input login as %s: ", yellow("'valid email'"))
	scanner.Scan()
	data := scanner.Text()
	login = strings.TrimSpace(data)
	if len(login) == 0 {
		fmt.Println(red("Login must not be empty please try again"))
		return
	}

	fmt.Printf("Input password as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	password = strings.TrimSpace(data)
	if len(password) == 0 {
		fmt.Println(red("Password must not be empty please try again"))
		return
	}

	token, err := u.userService.LoginUser(ctx, login, password)
	if err != nil {
		lib.UnpackGRPCError(err)
	} else {
		u.state.SetToken(token)
		u.state.SetIsAuthorized(true)
		u.state.SetLogin(login)
	}
}
