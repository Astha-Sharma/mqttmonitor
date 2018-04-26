package server

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	_ "github.com/vivekvasvani/mqttmonitor/config"
	sql "github.com/vivekvasvani/mqttmonitor/sql"
)

type AndroidCrashFreeTrends struct {
	Dailytreand []Trends `json:"dailytreand"`
	Weeklytrend []Weekly `json:"weeklytrend"`
}

type Trends struct {
	ID         int     `json:"id"`
	Date       string  `json:"date"`
	Percentage float64 `json:"percentage"`
	Delta      float64 `json:"delta"`
}

type Weekly struct {
	Date       string  `json:"date"`
	Week       int     `json:"week"`
	Percentage float64 `json:"percentage"`
}

func CrashFreeTrendsAndroid(ctx *fasthttp.RequestCtx) {
	db := sql.GetDBConnection()
	var (
		trendsA []Trends
		trendsB []Weekly
		date    string
		week    int
		per     float64
	)

	if err := db.Table("crash_free_users").
		Order("crash_free_users.id asc").
		Find(&trendsA).Error; err != nil {
		fmt.Println("Error ----#### ", err.Error)
	}
	rows, err := db.Raw(`SELECT CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2)) AS DATE, WEEK(CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2))) AS WEEK, percentage as percentage FROM crash_free_users WHERE CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2)) > DATE_SUB(NOW(), INTERVAL 20 WEEK) GROUP BY WEEK(CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2))) ORDER BY CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2));
`).Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&date, &week, &per)
		trendsB = append(trendsB, Weekly{date, week, per})
	}

	finalT := AndroidCrashFreeTrends{trendsA, trendsB}
	tests, _ := json.Marshal(&finalT)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(string(tests))
}

func CrashFreeTrendsIOS(ctx *fasthttp.RequestCtx) {
	db := sql.GetDBConnection()
	var (
		trendsA []Trends
		trendsB []Weekly
		date    string
		week    int
		per     float64
	)

	if err := db.Table("crash_free_users_ios").
		Order("crash_free_users_ios.id asc").
		Find(&trendsA).Error; err != nil {
		fmt.Println("Error ----#### ", err.Error)
	}
	rows, err := db.Raw(`SELECT CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2)) AS DATE, WEEK(CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2))) AS WEEK, percentage as percentage FROM crash_free_users WHERE CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2)) > DATE_SUB(NOW(), INTERVAL 20 WEEK) GROUP BY WEEK(CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2))) ORDER BY CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2));
`).Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&date, &week, &per)
		trendsB = append(trendsB, Weekly{date, week, per})
	}

	finalT := AndroidCrashFreeTrends{trendsA, trendsB}
	tests, _ := json.Marshal(&finalT)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(string(tests))
}

func CrashFreeTrendsAndroidTopbuilds(ctx *fasthttp.RequestCtx) {
	db := sql.GetDBConnection()
	var (
		trendsA []Trends
		trendsB []Weekly
		date    string
		week    int
		per     float64
	)

	if err := db.Table("crash_free_users_topbuilds").
		Order("crash_free_users_topbuilds.id asc").
		Find(&trendsA).Error; err != nil {
		fmt.Println("Error ----#### ", err.Error)
	}
	rows, err := db.Raw(`SELECT CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2)) AS DATE, WEEK(CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2))) AS WEEK, percentage as percentage FROM crash_free_users_topbuilds WHERE CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2)) > DATE_SUB(NOW(), INTERVAL 20 WEEK) GROUP BY WEEK(CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2))) ORDER BY CONCAT(RIGHT(date,4), "-", SUBSTRING(date,4,2),"-", LEFT(date,2));
`).Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&date, &week, &per)
		trendsB = append(trendsB, Weekly{date, week, per})
	}

	finalT := AndroidCrashFreeTrends{trendsA, trendsB}
	tests, _ := json.Marshal(&finalT)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(string(tests))
}
