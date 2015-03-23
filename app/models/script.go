package models

import (

)

type Script struct {
    Id int64
    Name string
    Interval string `db:"interval"`
    LastCheckin string `db:"last_checkin"`
    Severity int64
}

func (s *Script) Late() bool {
	return s.Severity != 0
}