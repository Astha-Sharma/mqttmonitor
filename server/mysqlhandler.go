package server

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	_ "github.com/vivekvasvani/mqttmonitor/config"
	sql "github.com/vivekvasvani/mqttmonitor/sql"
	"strings"
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

type Builds struct {
	BuildVersion []string `json:"version"`
}

type Ids struct {
	IdsE []string `json:"ids"`
}

type AndroidIssuesByVersion struct {
	CrashDetails  []SCrashDetails  `json:"crashDetails"`
	IssueTypes    []SIssuesTypes   `json:"issueTypes"`
	AreaWiseSplit []SAreaWiseSplit `json:"areaWiseSplit"`
}

type SCrashDetails struct {
	Title         string `json:"title" gorm:"Column:title"`
	Subtitle      string `json:"subtitle" gorm:"Column:subtitle"`
	Area          string `json:"area" gorm:"Column:area"`
	IssueType     string `json:"issueType" gorm:"Column:issueType"`
	ImpactLevel   int    `json:"impactLevel" gorm:"Column:impactLevel"`
	Occurances    int    `json:"occurances" gorm:"Column:occurrences"`
	UsersAffected int    `json:"usersAffected" gorm:"Column:usersAffected"`
	InfocusRatio  string `json:"infocusRatio" gorm:"Column:infocusRatio"`
	JiraId        string `json:"jiraId" gorm:"Column:jiraId"`
	Status        string `json:"status" gorm:"Column:status"`
	FixVersion    string `json:"fixVersion" gorm:"Column:fixVersion"`
	Assignee      string `json:"assignee" gorm:"Column:assignee"`
	StackTrace    string `json:"stackTrace" gorm:"Column:stackTrace"`
	FirstBuild    string `json:"firstBuild" gorm:"Column:firstBuild"`
	LastBuild     string `json:"lastBuild" gorm:"Column:lastBuild"`
}

type SIssuesTypes struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type SAreaWiseSplit struct {
	Text  string `json:"text"`
	Value string `json:"value"`
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

func AndroidBuildVersion(ctx *fasthttp.RequestCtx) {
	var (
		versionArray []string
		version      string
	)
	db := sql.GetDBConnection()
	rows, err := db.Raw(`SHOW TABLES`).Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&version)
		if strings.HasPrefix(version, "Issues_") {
			versionArray = append(versionArray, strings.SplitN(version, "_", 2)[1])
		}

	}
	finalResp := Builds{versionArray}
	tests, _ := json.Marshal(&finalResp)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(string(tests))
}

func AndroidCrashByVersion(ctx *fasthttp.RequestCtx) {
	db := sql.GetDBConnection()
	version := string(ctx.QueryArgs().Peek("version"))
	tableName := "Issues_" + version
	var (
		externalId    string
		externalIds   []string
		crashDetails  []SCrashDetails
		issueType     []SIssuesTypes
		areaWiseSplit []SAreaWiseSplit
	)

	rows, err := db.Raw("SELECT original_externalId From " + tableName).Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&externalId)
		externalIds = append(externalIds, externalId)

	}

	if err := db.Table("crashdetails").Where("externalId in (?)", externalIds).Find(&crashDetails).Error; err != nil {
		fmt.Println("Error ----#### ", err.Error)
	}
	crashes, _ := json.Marshal(&AndroidIssuesByVersion{crashDetails, issueType, areaWiseSplit})
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(string(crashes))
}
