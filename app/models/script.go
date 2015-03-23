package models

import (
	"log"
)

type Script struct {
    Id int64
    Name string
    Interval string `db:"interval"`
    LastCheckin string `db:"last_checkin"`
    Severity int64
}

func FindScript(id int) (bool, *Script) {
	obj, err := dbmap.Get(Script{}, id)
	script := obj.(*Script)

	if err != nil {
		log.Print("ERROR findScript: ")
		log.Println(err)
	}

	return (err == nil), script
}

func (s *Script) Late() bool {
	return s.Severity != 0
}