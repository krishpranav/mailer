package mailer

import (
	"crypto/tls"
	"net"
	"net/textproto"
)

type smtpClient struct {
	text       *textproto.Conn
	conn       net.Conn
	tls        bool
	serverName string
	ext        map[string]string

	a          []string
	localName  string
	didHello   bool
	helloError error
}

func newClient(conn net.Conn, host string) (*smtpClient, error) {
	text := textproto.NewConn(conn)
	_, _, err := text.ReadResponse(220)
	if err != nil {
		text.Close()
		return nil, err
	}
	c := &smtpClient{text: text, conn: conn, serverName: host, localName: "localhost"}
	_, c.tls = conn.(*tls.Conn)
	return c, nil
}
