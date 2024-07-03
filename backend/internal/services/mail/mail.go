package mail

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"strconv"
)

type Mailer struct {
	auth smtp.Auth
	cfg  Config
}

type Config struct {
	Host        string `env:"SMTP_HOST"`
	Port        int    `env:"SMTP_PORT"`
	Username    string `env:"SMTP_USERNAME"`
	Password    string `env:"SMTP_PASSWORD"`
	FrontendURL string `env:"FRONTEND_URL"`
}

func NewClient() Mailer {
	mailPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	cfg := Config{
		Host:        os.Getenv("SMTP_HOST"),
		Port:        mailPort,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return Mailer{auth: auth, cfg: cfg}
}

func (m Mailer) Send(to []string, subject string, message string) error {
	msg := []byte(
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n" +
			"From: Dokedu Support <support@dokedu.org>\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			message + "\r\n")

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	err := smtp.SendMail(addr, m.auth, "support@dokedu.org", to, msg)
	if err != nil {
		slog.Error("error while trying to send email", "err", err)
		return err
	}

	return nil
}

func (m Mailer) SendToken(to string, name string, token string) error {
	link := fmt.Sprintf("%s/login#token=%s", m.cfg.FrontendURL, token)
	subject := "Dokedu Drive Login Link"

	template, err := LoginLinkMailTemplate(name, link)
	if err != nil {
		slog.Error("error while trying to generate password reset mail template", "err", err)
		return err
	}

	return m.Send([]string{to}, subject, template)
}
