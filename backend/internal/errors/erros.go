package errors

import ("errors")

var (
	ErrDBConnection = errors.New("Can't connect to the database")
	ErrDBPing       = errors.New("Can't ping the database")
	ErrConfig       = errors.New("DB_Connection_String not set in .env file")
	ErrDBSaveReports = errors.New("Failed to save stock reports to the database")
)
