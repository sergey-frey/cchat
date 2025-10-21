package username

import (
	"github.com/brianvoe/gofakeit/v6"
)

func GenerateUsername() string {

	newUsername := gofakeit.Username()

	return newUsername
}
