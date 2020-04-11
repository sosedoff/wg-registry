package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type LetsEncryptConfig struct {
	Domain  string `json:"domain"`
	Email   string `json:"email"`
	Staging bool   `json:"staging"`
	Dir     string `json:"dir"`
	Enabled bool   `json:"enabled"`
}

func NewCertManager(config *LetsEncryptConfig) (*autocert.Manager, error) {
	policyFunc := func(ctx context.Context, host string) error {
		if host != config.Domain {
			return errors.New("invalid hostname")
		}
		return nil
	}

	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(config.Dir),
		Email:      config.Email,
		HostPolicy: policyFunc,
	}

	if config.Staging {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		manager.Client = &acme.Client{
			DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
			Key:          key,
		}
	}

	return &manager, nil
}
