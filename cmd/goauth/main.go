package main

import (
	"fmt"
	"github.com/mohammadrezajavid-lab/goauth/cmd/goauth/command"
	"os"
)

// @title GoAuth Service API
// @version 1.0
// @description This is a sample server for OTP-based authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
