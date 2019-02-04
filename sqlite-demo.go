package main

import (
	"./db"
	"fmt"
	"os"
)

var testDB db.SqliteDB

const projTab = `pid int primary key not null, proj_name varchar(50) not null`
const reportTab = `pid int not null,
										mod_date date not null,
										report_type varchar(10) not null,
										rid int not null primary key,
										report_file varchar(255) not null,
										FOREIGN KEY(pid) REFERENCES hc1projects(pid)`
const workerTab = `rid int not null, finding varchar(10),
										FOREIGN KEY(rid) REFERENCES hc1reports(rid)`

func main() {

	// *** Init test database
	testDB.DbName = "/home/pj001268/hctrail"
	testDB.DbDriver = "sqlite3"

	// *** Create test database
	db.CreateDB(&testDB)
	if testDB.DbError != nil {
		fmt.Println(testDB.DbError)
		os.Exit(1)
	}
	fmt.Println(testDB)

	//// *** Create test tables
	//testDB.CreateTable("hc1test","t1 int, t2 char(10)")
	//if testDB.DbError != nil {
	//	fmt.Println(testDB.DbError)
	//}

	//// *** Insert data into test table
	//testDB.InsertData("hc1test","11,'td1'")
	//if testDB.DbError != nil {
	//	fmt.Println(testDB.DbError)
	//}

	// *** Query data
	testDB.QueryDB("SELECT * FROM hc1test")
	if testDB.DbError != nil {
		fmt.Println(testDB.DbError)
	}

	// *** Close test database
	db.CloseDB(&testDB)
	if testDB.DbError != nil {
		fmt.Println(testDB.DbError)
	}

}
