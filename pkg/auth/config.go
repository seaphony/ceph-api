package auth

import "time"

type Config struct {
	AccessTokenLifespan  time.Duration `yaml:"accessTokenLifespan"`
	RefreshTokenLifespan time.Duration `yaml:"refreshTokenLifespan"`
	ClientID             string        `yaml:"clientID"`
	Issuer               string        `yaml:"issuer"`
}
