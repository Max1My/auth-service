package model

type MailMessageInfo struct {
	From         string
	To           string
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
	Token        string
}
