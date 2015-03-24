package app

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	j "github.com/brianseitel/snitch/app/jobs"
	"github.com/brianseitel/snitch/app/models"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/jobs/app/jobs"
)

var DB *gorp.DbMap

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(InitDB)
	revel.OnAppStart(func() { jobs.Every(1*time.Minute, j.TrackerJob{}) })

	// revel.OnAppStart(FillCache)
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

/** SET UP DATABASE STUFF **/

var InitDB func() = func() {
	connectionString := getConnectionString()
	if db, err := sql.Open("mysql", connectionString); err != nil {
		revel.ERROR.Fatal(err)
	} else {
		DB = &gorp.DbMap{
			Db:      db,
			Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	}

	// Defines the table for use by GORP
	// This is a function we will create soon.
	_ = DB.AddTableWithName(models.Checkin{}, "checkins").SetKeys(false, "Id")
	_ = DB.AddTableWithName(models.Script{}, "scripts").SetKeys(false, "Id")

	DB.CreateTablesIfNotExists()
}

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func getConnectionString() string {
	host := getParamString("db.host", "")
	port := getParamString("db.port", "3306")
	user := getParamString("db.user", "")
	pass := getParamString("db.password", "")
	dbname := getParamString("db.name", "auction")
	protocol := getParamString("db.protocol", "tcp")
	dbargs := getParamString("dbargs", " ")

	if strings.Trim(dbargs, " ") != "" {
		dbargs = "?" + dbargs
	} else {
		dbargs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbname, dbargs)
}
