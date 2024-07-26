package helper

import (
	"fmt"

	"github.com/joho/godotenv"
)

var SigningKey = []byte("secret")

func Configure(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		fmt.Printf("error at loading %s", filename)
		return err
	}
	return nil
}
