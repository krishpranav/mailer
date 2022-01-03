package mailer

import (
	"bytes"
	"net/textproto"
	"time"
)

type Email struct {
	from        string
	sender      string
	replyTo     string
	returnPath  string
	recipients  []string
	headers     textproto.MIMEHeader
	parts       []part
	attachments []*File
	inlines     []*File
	Charset     string
	Encoding    encoding
	Error       error
	SMTPServer  *smtpClient
	DkimMsg     string
}

type SMTPClient struct {
	Client      *smtpClient
	KeepAlive   bool
	SendTimeout time.Duration
}

type part struct {
	contentType string
	body        *bytes.Buffer
}
