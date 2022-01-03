package mailer

import (
	"errors"
	"mime"
	"path/filepath"
)

type File struct {
	FilePath string
	Name     string
	MimeType string
	B64Data  string
	Data     []byte
	Inline   bool
}

type attachType int

const (
	attachData attachType = iota
	attachB64
	attachFile
)

func (email *Email) Attach(file *File) *Email {
	if email.Error != nil {
		return email
	}

	var name = file.Name
	var mimeType = file.MimeType

	if len(name) == 0 && len(file.FilePath) > 0 {
		_, name = filepath.Split(file.FilePath)
	}

	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(name))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	attachTy, err := getAttachmentType(file)
	if err != nil {
		email.Error = errors.New("Mail Error: Failed to add attachment with following error: " + err.Error())
		return email
	}

	file.Name = name
	file.MimeType = mimeType

	switch attachTy {
	case attachData:
		email.attachData(file)
	case attachB64:
		email.Error = email.attachB64(file)
	case attachFile:
		email.Error = email.attachFile(file)
	}

	return email
}

func getAttachmentType(file *File) (attachType, error) {
	return nil
}

func (email *Email) attachB64(file *File) error {
	return nil
}

func (email *Email) attachFile(file *File) error {
	return nil
}

func (email *Email) attachData(file *File) {
}
