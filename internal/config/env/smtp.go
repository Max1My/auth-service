package env

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

const (
	SmtpHostEnvName     = "SMTP_HOST"
	SmtpPortEnvName     = "SMTP_PORT"
	SmtpUserEnvName     = "SMTP_USER"
	SmtpPasswordEnvName = "SMTP_PASSWORD"
)

type SmtpConfig interface {
}

type SmtpConfigData struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
}

func NewSmtpConfig() (*SmtpConfigData, error) {
	smtpHost := os.Getenv(SmtpHostEnvName)
	if len(smtpHost) == 0 {
		return nil, errors.New("smtp host not found")
	}
	smtpPortStr := os.Getenv(SmtpPortEnvName)
	if len(smtpPortStr) == 0 {
		return nil, errors.New("smtp port not found")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr) // изменено на Atoi для преобразования в int
	if err != nil {
		return nil, errors.New("invalid smtp port value")
	}
	smtpUser := os.Getenv(SmtpUserEnvName)
	if len(smtpUser) == 0 {
		return nil, errors.New("smtp user not found")
	}
	smtpPassword := os.Getenv(SmtpPasswordEnvName)
	if len(smtpPassword) == 0 {
		return nil, errors.New("smtp password not found")
	}

	return &SmtpConfigData{
		SmtpHost:     smtpHost,
		SmtpPort:     smtpPort,
		SmtpUser:     smtpUser,
		SmtpPassword: smtpPassword,
	}, nil
}
