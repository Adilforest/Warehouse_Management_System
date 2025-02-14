package utils

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body, attachmentPath string) error {
    from := os.Getenv("SMTP_EMAIL")
    password := os.Getenv("SMTP_PASSWORD")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")

    auth := smtp.PlainAuth("", from, password, smtpHost)

    msg := []byte("To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "MIME-Version: 1.0\r\n" +
        "Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
        body)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        return err
    }

    return nil
}