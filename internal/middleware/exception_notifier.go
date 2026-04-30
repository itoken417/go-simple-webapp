package middleware

import "github.com/itoken417/goutils/mailsender"

var errSender *mailsender.Sender

func InitErrorNotifier(sender *mailsender.Sender) {
	errSender = sender
}
