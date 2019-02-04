package db

// *** DB interface for using different DB types
type Importer interface {
	NewDb()
	CloseDb()
}

func CreateDB(db Importer) {

	db.NewDb()

}

func CloseDB(db Importer) {

	db.CloseDb()

}
