package mailer

import (
	"errors"
	"io/ioutil"
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
	if len(file.Data) > 0 {
		if len(file.Name) == 0 {
			return 0, errors.New("attach from bytes requires a name")
		}
		return attachData, nil
	}

	if len(file.B64Data) > 0 {
		if len(file.Name) == 0 {
			return 0, errors.New("attach from base64 string requires a name")
		}
		return attachB64, nil
	}

	if len(file.FilePath) > 0 {
		return attachFile, nil
	}

	return 0, errors.New("empty attachment")
}

func (email *Email) attachB64(file *File) error {
	return nil
}

func (email *Email) attachFile(file *File) error {
	data, err := ioutil.ReadFile(file.FilePath)
	if err != nil {
		return errors.New("Mail Error: Failed to add file error: " + err.Error())
	}

	email.attachData(&File{
		Name:     file.Name,
		MimeType: file.MimeType,
		Data:     data,
		Inline:   file.Inline,
	})

	return nil
}

func (email *Email) attachData(file *File) {
	if file.Inline {
		email.inlines = append(email.inlines, file)
	} else {
		email.attachments = append(email.attachments, file)
	}
}
