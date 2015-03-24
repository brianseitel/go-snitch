package controllers

import (
	"github.com/revel/revel"
	"github.com/brianseitel/snitch/app"
	"github.com/brianseitel/snitch/app/models"
)
type App struct {
	GorpController
}

func (c App) Index() revel.Result {
	results, err := app.DB.Select(models.Script{},
		"SELECT name as name, MAX(last_checkin) as last_checkin, `interval` FROM scripts GROUP BY name")
	if err != nil {
		panic(err)
	}

	var scripts []*models.Script
	for _, r := range results {
		s := r.(*models.Script)
		scripts = append(scripts, s)
	}

	return c.Render(scripts)
}
