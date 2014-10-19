package annonce

import (
	"gopkg.in/mgo.v2"
)

var (
	session *Session
)

const (
	url = "mongodb://localhost"
)

func Session() *Session {
	if session == nil {
		s, err = mgo.Dial(url)
		if err {
			panic(err)
		}
		session = s
	}
	return s
}
