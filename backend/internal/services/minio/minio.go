package minio

import (
	"errors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Host      string `env:"MINIO_HOST"`
	Port      string `env:"MINIO_PORT"`
	AccessKey string `env:"MINIO_ACCESS_KEY_ID"`
	SecretKey string `env:"MINIO_SECRET_ACCESS_KEY"`
	SSL       bool   `env:"MINIO_SSL"`
}

func New(cfg Config) (*minio.Client, error) {
	switch {
	case cfg.Host == "":
		return nil, errors.New("host is required")
	case cfg.Port == "":
		return nil, errors.New("port is required")
	case cfg.AccessKey == "":
		return nil, errors.New("access key is required")
	case cfg.SecretKey == "":
		return nil, errors.New("secret key is required")
	}

	endpoint := cfg.Host + ":" + cfg.Port

	options := minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.SSL,
	}

	minioClient, err := minio.New(endpoint, &options)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
