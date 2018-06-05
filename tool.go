package main

// yum install freetds-devel
// go get github.com/minus5/gofreetds

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/minus5/gofreetds"
)

type resultImporExport struct {
	taskID          int
	taskType        string
	lifecycle       string
	createdAt       time.Time
	lastUpdated     time.Time
	databaseName    string
	s3Arn           string
	overwriteS3File string
	kmsKey          string
	taskProgress    int
	taskInfo        string
}

type resultStatus struct {
	taskID          int
	taskType        string
	databaseName    string
	percentComplete int
	dutation        string
	lifecycle       string
	taskInfo        string
	lastUpdated     time.Time
	createdAt       time.Time
	s3Arn           string
	overwriteS3File string
	kmsKey          string
}

type config struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"password"`
}

func loadConfiguration() config {

	var config config

	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(configFile, &config)
	return config
}

func dbConnection() *sql.DB {
	config := loadConfiguration()
	db, err := sql.Open("mssql", "host="+config.Host+";user="+config.User+";pwd="+config.Pass)
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func statusTask(dbName string) {
	db := dbConnection()
	defer db.Close()

	var result resultStatus

	row := db.QueryRow("exec msdb.dbo.rds_task_status @db_name=?", dbName)
	row.Scan(&result.taskID, &result.taskType, &result.databaseName, &result.percentComplete, &result.dutation, &result.lifecycle, &result.taskInfo, &result.lastUpdated, &result.createdAt, &result.s3Arn, &result.overwriteS3File, &result.kmsKey)

	fmt.Printf("%v [%s] %s %v %v\n%v\n", result.createdAt.Format("02-01-2006"), dbName, result.taskType, result.s3Arn, result.lifecycle, result.taskInfo)

}

func exportDB(dbName, s3Name, backupName string) {
	db := dbConnection()
	defer db.Close()

	if (dbName != "") && (s3Name != "") && (backupName != "") {

		var result resultImporExport

		query := `
		exec msdb.dbo.rds_backup_database
				@source_db_name=?,
				@s3_arn_to_backup_to='arn:aws:s3:::` + s3Name + `/` + backupName + `',
				@overwrite_S3_backup_file=1,
				@type='FULL';
		`
		row := db.QueryRow(query, dbName)
		row.Scan(&result.taskID, &result.taskType, &result.lifecycle, &result.createdAt, &result.lastUpdated, &result.databaseName, &result.s3Arn, &result.overwriteS3File, &result.kmsKey, &result.taskProgress, &result.taskInfo)
		fmt.Printf("%v [%s] %s %v %v%% %v\n", result.createdAt.Format("02-01-2006"), dbName, result.taskType, result.s3Arn, result.taskProgress, result.lifecycle)
	} else {
		fmt.Println("dbName,s3Name,backupName - one of this is empty!")
	}
}

func importDB(dbName, s3Name, restoreName string) {
	db := dbConnection()
	defer db.Close()

	if (dbName != "") && (s3Name != "") && (restoreName != "") {
		var result int
		var resultImport resultImporExport
		db.QueryRow("SELECT DB_ID('" + dbName + "')").Scan(&result)
		if result == 0 {
			query := `
			exec msdb.dbo.rds_restore_database
				@restore_db_name='?',
				@s3_arn_to_restore_from='arn:aws:s3:::` + s3Name + `/` + restoreName + `';
			`
			row := db.QueryRow(query, dbName)
			row.Scan(&resultImport.taskID, &resultImport.taskType, &resultImport.lifecycle, &resultImport.createdAt, &resultImport.lastUpdated, &resultImport.databaseName, &resultImport.s3Arn, &resultImport.overwriteS3File, &resultImport.kmsKey, &resultImport.taskProgress, &resultImport.taskInfo)
			fmt.Printf("%v [%s] %s %v %v%% %v\n%v\n", resultImport.createdAt.Format("02-01-2006"), dbName, resultImport.taskType, resultImport.s3Arn, resultImport.taskProgress, resultImport.lifecycle, resultImport.taskInfo)
		} else {
			fmt.Println("Cant restore backup to exist database!")
		}
	} else {
		fmt.Println("dbName,s3Name,restoreName - one of this is empty!")
	}
}

func usage() string {
	options := `Usage:
	tool status dbName
	tool export dbName s3Name backupName
	tool import dbName s3Name restoreName`
	return options
}

func main() {
	argsWithProg := os.Args
	totalArguments := len(argsWithProg)
	if totalArguments > 1 {
		switch argsWithProg[1] {
		case "status":
			if totalArguments == 3 {
				statusTask(argsWithProg[2])
			} else {
				fmt.Println(usage())
			}
		case "export":
			if totalArguments == 5 {
				exportDB(argsWithProg[2], argsWithProg[3], argsWithProg[4])
			} else {
				fmt.Println(usage())
			}
		case "import":
			if totalArguments == 5 {
				importDB(argsWithProg[2], argsWithProg[3], argsWithProg[4])
			} else {
				fmt.Println(usage())
			}
		default:
			fmt.Println(usage())
		}
	} else {
		fmt.Println(usage())
	}
}
