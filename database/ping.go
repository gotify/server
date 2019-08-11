package database

// Ping pings the database to verify the connection.
func (d *GormDatabase) Ping() error {
	return d.DB.DB().Ping()
}
