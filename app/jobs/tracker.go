package jobs

import (
	"fmt"
	"github.com/brianseitel/snitch/app"
)

type TrackerJob struct {}

type Result struct {
	id int
	late int
}

func (j TrackerJob) Run() {
	fmt.Println("Tracking!!!")

	sql := `SELECT s1.name as name, CASE WHEN MAX(s2.last_checkin) IS NULL THEN 0 ELSE 1 END as last_checkin
FROM scripts s1
LEFT JOIN scripts s2 ON s2.id = s1.id AND s2.last_checkin > DATE_SUB(NOW(), INTERVAL s2.interval MINUTE)
GROUP BY s1.name`
	
	results := []Result{}

	results, err := app.DB.Select(models.Script{}, sql)

	for result := range results {
		fmt.Println(result.id + " => " + result.late)
	}

}