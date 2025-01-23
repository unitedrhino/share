package utils

import (
	"crypto/tls"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

func SenEmail(c conf.Email, to []string, subject string, body string) error {
	if c.Port == 0 {
		c.Port = 465
	}
	auth := LoginAuth(c.From, c.Secret)
	e := email.NewEmail()
	if c.Nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", c.Nickname, c.From)
	} else {
		e.From = c.From
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	var err error
	hostAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	if c.IsSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: c.Host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	if err != nil && strings.HasPrefix(err.Error(), "short response:") {
		return nil
	}
	return err
}

type loginAuth struct {
	username, password string
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.Default.AddDetail(fromServer)
		}
	}
	return nil, nil
}
