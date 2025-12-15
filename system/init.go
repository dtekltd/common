package system

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("> init system")
	// load .env file
	if _, err := os.Stat(".env"); err != nil {
		// NOTE: please do not use Logger here
		// it was not initialized
		fmt.Printf("Config .env file does not exist! %s", err.Error())
	} else {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Printf("Error loading .env file! %s", err.Error())
		}
	}

	// init logger
	initAntiglossLogger()
}
