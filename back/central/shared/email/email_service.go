package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Interfaz genérica de envío de correo
type IEmailService interface {
	SendHTML(ctx context.Context, to, subject, html string) error
}

type EmailService struct {
	config env.IConfig
	logger log.ILogger
}

func New(config env.IConfig, logger log.ILogger) IEmailService {
	return &EmailService{
		config: config,
		logger: logger,
	}
}

func (e *EmailService) SendHTML(ctx context.Context, to, subject, html string) error {
	return e.sendEmail(ctx, to, subject, html)
}

func (e *EmailService) sendEmail(ctx context.Context, to, subject, body string) error {
	// Configuración SMTP
	smtpHost := e.config.Get("SMTP_HOST")
	smtpPort := e.config.Get("SMTP_PORT")
	smtpUser := e.config.Get("SMTP_USER")
	smtpPass := e.config.Get("SMTP_PASS")
	fromEmail := e.config.Get("FROM_EMAIL")
	useTLS := e.config.Get("SMTP_USE_TLS") == "true"
	useSTARTTLS := e.config.Get("SMTP_USE_STARTTLS") == "true"

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || fromEmail == "" {
		e.logger.Error().Msg("Configuración SMTP incompleta")
		return fmt.Errorf("configuración SMTP incompleta")
	}

	// Crear la dirección completa del servidor SMTP
	addr := smtpHost + ":" + smtpPort

	// Crear el mensaje
	message := e.buildMessage(fromEmail, to, subject, body)

	// Configurar autenticación
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Enviar según el método de seguridad
	var err error
	if useTLS {
		err = e.sendWithTLS(addr, auth, fromEmail, []string{to}, message)
	} else if useSTARTTLS {
		err = e.sendWithSTARTTLS(addr, auth, fromEmail, []string{to}, message)
	} else {
		err = smtp.SendMail(addr, auth, fromEmail, []string{to}, message)
	}

	if err != nil {
		e.logger.Error().
			Err(err).
			Str("to", to).
			Str("subject", subject).
			Str("smtp_host", smtpHost).
			Str("smtp_port", smtpPort).
			Str("security", e.getSecurityMethod(useTLS, useSTARTTLS)).
			Msg("Error enviando email")
		return fmt.Errorf("error enviando email: %v", err)
	}

	e.logger.Info().
		Str("to", to).
		Str("subject", subject).
		Str("smtp_host", smtpHost).
		Str("smtp_port", smtpPort).
		Str("security", e.getSecurityMethod(useTLS, useSTARTTLS)).
		Msg("Email enviado exitosamente")

	return nil
}

func (e *EmailService) buildMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "quoted-printable"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return []byte(message)
}

func (e *EmailService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Crear configuración TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         strings.Split(addr, ":")[0],
	}

	// Establecer conexión TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("error estableciendo conexión TLS: %v", err)
	}
	defer conn.Close()

	// Crear cliente SMTP sobre la conexión TLS
	client, err := smtp.NewClient(conn, strings.Split(addr, ":")[0])
	if err != nil {
		return fmt.Errorf("error creando cliente SMTP TLS: %v", err)
	}
	defer client.Quit()

	// Autenticar
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("error en autenticación TLS: %v", err)
	}

	// Establecer remitente
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("error estableciendo remitente TLS: %v", err)
	}

	// Establecer destinatarios
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("error estableciendo destinatario TLS %s: %v", recipient, err)
		}
	}

	// Enviar datos
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("error iniciando datos TLS: %v", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("error escribiendo mensaje TLS: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error cerrando datos TLS: %v", err)
	}

	return nil
}

func (e *EmailService) sendWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Establecer conexión normal
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("error estableciendo conexión STARTTLS: %v", err)
	}
	defer client.Quit()

	// Configurar TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         strings.Split(addr, ":")[0],
	}

	// Iniciar STARTTLS
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("error iniciando STARTTLS: %v", err)
	}

	// Autenticar
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("error en autenticación STARTTLS: %v", err)
	}

	// Establecer remitente
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("error estableciendo remitente STARTTLS: %v", err)
	}

	// Establecer destinatarios
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("error estableciendo destinatario STARTTLS %s: %v", recipient, err)
		}
	}

	// Enviar datos
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("error iniciando datos STARTTLS: %v", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("error escribiendo mensaje STARTTLS: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error cerrando datos STARTTLS: %v", err)
	}

	return nil
}

func (e *EmailService) getSecurityMethod(useTLS, useSTARTTLS bool) string {
	if useTLS {
		return "TLS"
	} else if useSTARTTLS {
		return "STARTTLS"
	}
	return "PLAIN"
}
