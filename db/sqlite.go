package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// SQLite data structure
type SqliteDB struct {
	DbName   string
	DbDriver string
	DbHandle *sql.DB
	DbError  error
}

// Structure of query result
type qryRes struct {
	// Column structure as list [column name]column type
	ColInfo map[string]string
	// Query data as Array [row index][column data]
	QryData [][]interface{}
}

// *** SQLite functions for handling the database

/* NewDb
#
# DESCRIPTION: open/create SQLite DB
# PARAMETER:   <SqliteDB.DbName>    name of the handled database
#              <SqliteDB.DbDriver>  name of the database driver
# RETURN:      <SqliteDB.DbHandle>  handle of the current used database
#              <SqliteDB.DbError>   currently saved error of DB operations
#
*/
func (sdb *SqliteDB) NewDb() {

	// Check if necessary params are available
	if sdb.DbDriver == "" || sdb.DbName == "" {
		// If not: save error
		sdb.DbError = fmt.Errorf("missing parameter(s) for NewDb")
	} else {
		// If yes: open existing database or create a new one
		db, err := sql.Open(sdb.DbDriver, sdb.DbName)
		if err != nil {
			sdb.DbError = err
		} else {
			sdb.DbHandle = db
		}
	}

}

/* CloseDb
#
# DESCRIPTION: close opened SQLite DB
# PARAMETER:   <SqliteDB.DbHandle>  handle of the database to be closed
# RETURN:      <SqliteDB.DbError>   currently saved error of DB operations
#
*/
func (sdb *SqliteDB) CloseDb() {

	if sdb.DbHandle != nil {
		sdb.DbError = sdb.DbHandle.Close()
	}

}

/* CreateTable
#
# DESCRIPTION: create Sqlite table (if not exists)
# PARAMETER:   tblName  - name of the new table
#              tblData  - schema of new table
# RETURN:      <SqliteDB.DbError>   currently saved error of DB operations
#
*/
func (sdb *SqliteDB) CreateTable(tblName string, tblData string) {

	if sdb.DbHandle != nil && tblData != "" && tblName != "" {
		// Check conn to database
		if err := sdb.DbHandle.Ping(); err != nil {
			sdb.DbError = err
		} else {
			sqlcmd := "CREATE TABLE " + tblName + "(" + tblData + ");"
			stmt, err := sdb.DbHandle.Prepare(sqlcmd)
			if err != nil {
				sdb.DbError = err
			} else {
				defer stmt.Close()
				_, sdb.DbError = stmt.Exec()
			}
		}
	}

}

/* InsertData
#
# DESCRIPTION: insert data record into Sqlite table
# PARAMETER:   tblName  - name of the table
#              tblRec   - new record
# RETURN:      <SqliteDB.DbError>   currently saved error of DB operations
#
*/
func (sdb *SqliteDB) InsertData(tblName string, tblRec string) {

	if sdb.DbHandle != nil && tblRec != "" && tblName != "" {
		// Check conn to database
		if err := sdb.DbHandle.Ping(); err != nil {
			sdb.DbError = err
		} else {
			sqlcmd := "INSERT INTO " + tblName + " VALUES (" + tblRec + ");"
			fmt.Println("cmd: ", sqlcmd)
			stmt, err := sdb.DbHandle.Prepare(sqlcmd)
			if err != nil {
				sdb.DbError = err
			} else {
				defer stmt.Close()
				_, sdb.DbError = stmt.Exec()
			}
		}
	}

}

/* QueryDB
#
# DESCRIPTION: query Sqlite table (SELECT)
# PARAMETER:   qryCmd  - SQL command
# RETURN:      <SqliteDB.DbError>   currently saved error of DB operations
#
*/
func (sdb *SqliteDB) QueryDB(qryCmd string) {

	var qryResult qryRes
	var colData = make([][]interface{}, 1)

	if sdb.DbHandle != nil && qryCmd != "" {
		// Check conn to database
		if err := sdb.DbHandle.Ping(); err != nil {
			sdb.DbError = err
		} else {
			rows, err := sdb.DbHandle.Query(qryCmd)
			if err != nil {
				sdb.DbError = err
				return
			}
			defer rows.Close()

			// *** Get column info of query result
			qryResult.ColInfo = GetColInfo(rows)

			noOfCols := len(qryResult.ColInfo)

			//tableData := make([]map[string]interface{}, 0)
			values := make([]interface{}, noOfCols)
			valuePtrs := make([]interface{}, noOfCols)
			for rows.Next() {
				for i := 0; i < noOfCols; i++ {
					valuePtrs[i] = &values[i]
				}
				rows.Scan(valuePtrs...)
				colData = append(colData, values)
				//entry := make(map[string]interface{})
				//for i, col := range columns {
				//	var v interface{}
				//	val := values[i]
				//	b, ok := val.([]byte)
				//	if ok {
				//		v = string(b)
				//	} else {
				//		v = val
				//	}
				//	entry[strconv.Itoa(i)+"_"+col] = v
				//}
				//tableData = append(tableData, entry)
			}
			qryResult.QryData = colData
			// TODO: weiter
		}
	} else {
		sdb.DbError = fmt.Errorf("missing params for queryDB")
	}

}

/* GetColInfo
#
# DESCRIPTION: Get column names and types of query result
# PARAMETER:   qryResult - *sql.Rows
# RETURN:      colInfo    - map[string]string (map[col name]col type)
#
*/
func GetColInfo(qryResult *sql.Rows) map[string]string {

	var ci = make(map[string]string)
	var colName string
	var colType string

	// First check if valid query results are available
	if qryResult != nil {
		// If yes: get column info
		cols, err := qryResult.ColumnTypes()
		// If column info is not available, exit with empty map
		if err != nil {
			return ci
		}
		// For each found column ...
		for _, col := range cols {
			// Get column name
			colName = col.Name()
			// Get column type
			colType = col.DatabaseTypeName()
			ci[colName] = colType
		}
	}

	return ci

}
