package mgodb

import (
	"gopkg.in/mgo.v2"
)

var Session *mgo.Session

func Connect() {

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		defer session.Close()
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	Session = session
}
