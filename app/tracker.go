package app

import (
	"fmt"
	"github.com/brianseitel/snitch/app/models"
)

type TrackerJob struct{}

type Result struct {
	Id   int
	Name string
	Late int
}

func (j TrackerJob) Run() {
	fmt.Println("Tracking!!!")

	sql := `SELECT s1.id, s1.name as name, CASE WHEN MAX(s2.last_checkin) IS NULL THEN 0 ELSE 1 END as late
FROM scripts s1
LEFT JOIN scripts s2 ON s2.id = s1.id AND s2.last_checkin < DATE_SUB(NOW(), INTERVAL s2.interval MINUTE)
GROUP BY s1.name`

	results := []Result{}

	_, err := DB.Select(&results, sql)
	if err != nil {
		panic(err)
	}
	for result := range results {
		message := "Checking " + results[result].Name + "... "
		if results[result].Late == 1 {
			var script models.Script
			_ = DB.SelectOne(&script, `SELECT * FROM scripts WHERE id = ?`, results[result].Id)

			script.Severity = script.Severity + 1
			_, err = DB.Update(&script)

			if err != nil {
				panic(err)
			}
			message += " LATE!"
		} else {
			message += " Success!"
		}

		fmt.Println(message)
	}

}
