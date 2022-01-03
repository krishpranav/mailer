package mailer

type auth interface {
	start()
	next()
}
