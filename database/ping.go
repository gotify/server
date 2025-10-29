package database

// Ping pings the database to verify the connection.
func (d *GormDatabase) Ping() error {
	sqldb, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqldb.Ping()
}
