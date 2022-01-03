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

func (c *smtpClient) close() error {
	return c.text.Close()
}

func (c *smtpClient) hello() error {
	if !c.didHello {
		c.didHello = true
		err := c.ehlo()
		if err != nil {
			c.helloError = c.helo()
		}
	}
	return c.helloError
}

func (c *smtpClient) helo() error {
	c.ext = nil
	_, _, err := c.cmd(250, "HELO %s", c.localName)
	return err
}
