package mysql

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vivekvasvani/mqttmonitor/config"
)

type mysqlConnector struct {
	DbPool *gorm.DB
}

var connector *mysqlConnector

func InitMysql() *mysqlConnector {
	if connector != nil {
		log.Info("Warning Mysql: DataBase is already initialized, skipping")
		return connector
	}
	log.Info("Warning Mysql: DataBase was not initialized ....initializing again")
	var err error
	connector, err = initDB()
	if err != nil {
		panic(err)
	}
	return connector
}

func GetDBConnection() *gorm.DB {
	return connector.DbPool
}

// DB Initialization
func initDB() (*mysqlConnector, error) {
	fmt.Println("data==", config.Get("database"))
	dbCfg := config.Get("database").(config.DBCfgTemplate)
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbCfg.Username, dbCfg.Password, dbCfg.Server, dbCfg.Port, dbCfg.Schema)
	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	if maxCons := dbCfg.MaxConnection; maxCons > 0 {
		db.DB().SetMaxOpenConns(maxCons)
		db.DB().SetMaxIdleConns(maxCons / 3)
	}

	return &mysqlConnector{db}, nil
}
