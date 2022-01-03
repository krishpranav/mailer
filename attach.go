package mailer

type File struct {
	FilePath string
	Name     string
	MimeType string
	B64Data  string
	Data     []byte
	Inline   bool
}

type attachType int
