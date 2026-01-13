package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds configuration for SMTP email sending.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

// Sender is an email sender that uses SMTP.
type Sender struct {
	config SMTPConfig
}

// NewSender creates a new SMTP email sender.
func NewSender(config SMTPConfig) *Sender {
	return &Sender{config: config}
}

// Send sends an email using SMTP.
// to: recipient email address
// subject: email subject line
// htmlBody: HTML version of the email
// textBody: plain text version of the email (fallback)
func (s *Sender) Send(to, subject, htmlBody, textBody string) error {
	// Build email headers and body
	from := s.config.From
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/alternative; boundary=\"boundary\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n--boundary\r\n"
	message += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	message += "Content-Transfer-Encoding: 7bit\r\n\r\n"
	message += textBody + "\r\n"

	message += "--boundary\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	message += "Content-Transfer-Encoding: 7bit\r\n\r\n"
	message += htmlBody + "\r\n"

	message += "--boundary--"

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Set up authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Send email
	if s.config.UseTLS {
		// Use TLS
		return s.sendWithTLS(addr, auth, s.config.From, []string{to}, []byte(message))
	}

	// Use plain SMTP
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(message))
}

// sendWithTLS sends email using TLS encryption.
func (s *Sender) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Connect to SMTP server with TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return nil
}

// Helper function to convert HTML to plain text (simple implementation)
func htmlToText(html string) string {
	// Very basic HTML to text conversion
	text := html
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "<p>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "<h1>", "\n")
	text = strings.ReplaceAll(text, "</h1>", "\n")
	text = strings.ReplaceAll(text, "<h2>", "\n")
	text = strings.ReplaceAll(text, "</h2>", "\n")

	// Remove all other HTML tags (simple regex would be better but this works)
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	return text
}
