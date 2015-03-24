package models

type Checkin struct {
	Id          int64
	ScriptId    int64  `db:"script_id"`
	LastCheckin string `db:"last_checkin"`
}
