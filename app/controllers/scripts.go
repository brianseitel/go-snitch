package controllers	

import (
	"encoding/json"
    "github.com/brianseitel/snitch/app/models"
    "github.com/revel/revel"
    "time"
)

type Scripts struct {
    GorpController
}

func (c Scripts) parseScriptItem() (models.Script, error) {
	script := models.Script{}
	err := json.NewDecoder(c.Request.Body).Decode(&script)

	if len(script.Url) == 0 {
		script.Url = script.GenerateUrl()
	}

	return script, err
}

func (c Scripts) Add() revel.Result {
	if script, err := c.parseScriptItem(); err != nil {
		return c.RenderText("Unable to parse entry.")
	} else {
		script.Validate(c.Validation)

		if c.Validation.HasErrors() {
			return c.RenderText("Input has errors.")
		} else {
			if err := c.Txn.Insert(&script); err != nil {
				return c.RenderText("Error inserting record.")
			} else {
				return c.RenderJson(script)
			}
		}
	}
}

func (c Scripts) List() revel.Result {
    lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
    limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
    scripts, err := c.Txn.Select(models.Script{}, 
        `SELECT * FROM scripts WHERE Id > ? LIMIT ?`, lastId, limit)
    if err != nil {
        return c.RenderText(
            "Error trying to get records from DB.")
    }
    return c.Render(scripts)
}

func (c Scripts) DoCheckin() revel.Result {
	name := c.Params.Get("name")

	var script models.Script
	err := c.Txn.SelectOne(&script,
		`SELECT * FROM scripts WHERE Name = ?`, name)

	if err != nil {
		return c.RenderText(
			"Error trying to get records from DB.")
	}


	now := time.Now().Format("2006-01-02 15:04:05")
	checkin := models.Checkin{Id: 0, ScriptId: script.Id, LastCheckin: now}

	if err := c.Txn.Insert(&checkin); err != nil {
		return c.RenderText(
			"Failed to record checkin.")
	} else {
		script.LastCheckin = now
		script.Severity = 0

		if result, err := c.Txn.Update(&script); err != nil {
			return c.RenderText(
				"Failed to update script info.")
		} else {
			return c.RenderJson(result)
		}
	}

	return c.RenderJson(script)
}