package models

import (
	"math/rand"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/revel/revel"
)

type Script struct {
	Id          int64
	Name        string
	Interval    string `db:"interval"`
	LastCheckin string `db:"last_checkin"`
	Severity    int64
	Url         string
}

func (s *Script) ReadableInterval() string {
	switch s.Interval {
		case `1`:
			return "Every minute"
		case `30`:
			return "Every 30 minutes"
		case `60`:
			return "Hourly"
		case `120`:
			return "Every two hours"
		case `86400`: 
			return "Daily"
		case `604800`:
			return "Weekly"
		case `2592000`:
			return "Monthly"
	}

	return "Invalid Interval"
}

func (s *Script) Late() bool {
	return s.Severity > 0
}

func (s *Script) LastSeen() string {
	location, _ := time.LoadLocation("America/Los_Angeles")
	then, err := time.ParseInLocation("2006-01-02 15:04:05", s.LastCheckin, location)
	if err != nil {
		panic(err)
	}

	return humanize.Time(then)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *Script) GenerateUrl() string {
	b := make([]rune, 12) // 12 letters long
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
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
