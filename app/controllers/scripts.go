package controllers	

import (
	"github.com/brianseitel/snitch/app"
    "github.com/brianseitel/snitch/app/models"
    "github.com/revel/revel"
    "time"
)

type Scripts struct {
    GorpController
}

func (c Scripts) Add(name string, interval string) revel.Result {
	script := models.Script{}
	script.Name = name
	script.Interval = interval
	script.Url = script.GenerateUrl()
	script.LastCheckin = time.Now().Local().Format("2006-01-02 15:04:05")

	if err := app.DB.Insert(&script); err != nil {
		panic(err)
		return c.RenderText("Error inserting record.")
	} else {
		return c.Redirect("/list")
	}
}

func (c Scripts) List() revel.Result {
    lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
    limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
    scripts, err := app.DB.Select(models.Script{}, 
        `SELECT * FROM scripts WHERE Id > ? LIMIT ?`, lastId, limit)
    if err != nil {
    	panic(err)
        return c.RenderText(
            "Error trying to get records from DB.")
    }
    return c.Render(scripts)
}

func (c Scripts) DoCheckin(url string) revel.Result {

	var script models.Script
	err := app.DB.SelectOne(&script,
		`SELECT * FROM scripts WHERE url = ?`, url)

	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}


	now := time.Now().Local().Format("2006-01-02 15:04:05")
	checkin := models.Checkin{Id: 0, ScriptId: script.Id, LastCheckin: now}

	if err := app.DB.Insert(&checkin); err != nil {
		return c.RenderText(
			"Failed to record checkin.")
	} else {
		script.LastCheckin = now
		script.Severity = 0

		if result, err := app.DB.Update(&script); err != nil {
			return c.RenderText(
				"Failed to update script info.")
		} else {
			return c.RenderJson(result)
		}
	}

	return c.RenderJson(script)
}