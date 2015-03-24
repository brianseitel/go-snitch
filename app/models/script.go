package models

import (
	"math/rand"
    "github.com/revel/revel"
	
)

type Script struct {
    Id int64
    Name string
    Interval string `db:"interval"`
    LastCheckin string `db:"last_checkin"`
    Severity int64
    Url string
}

func (s *Script) Late() bool {
	return s.Severity != 0
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *Script) GenerateUrl() string {
	b := make([]rune, 32) // 32 letters long
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func (b *Script) Validate(v *revel.Validation) {

    v.Check(b.Name,
        revel.ValidRequired(),
        revel.ValidMaxSize(64))

    v.Check(b.Interval,
        revel.ValidRequired(),
        revel.ValidMaxSize(32))

    v.Check(b.LastCheckin,
        revel.ValidRequired())

    v.Check(b.Severity,
        revel.ValidRequired())
}