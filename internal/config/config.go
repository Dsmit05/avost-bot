package config

import (
	"os"
	"time"
)

var (
	// BuildVersion sets on compile time.
	BuildVersion = ""
)

const (
	tokenBotEnv     = "BOT_TOKEN"
	mainURLEnv      = "MAIN_URL"
	mirrorURLEnv    = "MIRROR_URL"
	secretKeyJWTEnv = "JWT"
	addressENV      = "ADDRESS"
)

type AnimeBotConfg struct {
	Token                  string
	LongPollerTimeout      time.Duration
	LastAnimeUpdateTimeout time.Duration
	APIConfig              APIConfig
	Address                string
	BuildVersion           string
	MainURL, MirrorURL     string
}

type APIConfig struct {
	SecretKeyJWT string
	Addr         string
}

func NewAnimeBotConfg(longPollerTimeout, lastAnimeUpdateTimeout time.Duration, addrAPI string) *AnimeBotConfg {
	token := os.Getenv(tokenBotEnv)
	mainURL := os.Getenv(mainURLEnv)
	mirrorURL := os.Getenv(mirrorURLEnv)
	secretKeyJWT := os.Getenv(secretKeyJWTEnv)
	Address := os.Getenv(addressENV)

	return &AnimeBotConfg{
		Token:                  token,
		LongPollerTimeout:      longPollerTimeout,
		LastAnimeUpdateTimeout: lastAnimeUpdateTimeout,
		APIConfig:              APIConfig{SecretKeyJWT: secretKeyJWT, Addr: addrAPI},
		BuildVersion:           BuildVersion,
		MainURL:                mainURL,
		MirrorURL:              mirrorURL,
		Address:                Address,
	}
}
