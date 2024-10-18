package datastoretest

import "github.com/brave/go-sync/datastore"

// ResetSQLTables clears SQL tables.
func ResetSQLTables(sqlDB *datastore.SQLDB) error {
	_, err := sqlDB.Exec("DELETE FROM chains")
	return err
}
