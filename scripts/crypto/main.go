package main

import (
	"fmt"
	"github.com/Dsmit05/avost-bot/internal/api/cryptography"
	"os"
	"time"
)

// Генерация токена для доступа к api
func main() {
	cr := cryptography.NewTokenJWT(os.Getenv("JWT"))
	token, err := cr.CreateToken("user@email.com", "user", time.Minute*15)
	fmt.Println(token, err)
}
