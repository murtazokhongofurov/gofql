package gofql

import "database/sql"

func main() {

}

const (
	p = "postgres"
	m = "mysql"
	s = "sqlite3"
)

type ORM struct {
	db *sql.DB
}

func New(driver string, dataSourceName string) (*ORM, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &ORM{db: db}, nil
}

func (o *ORM) Close() error {
	return o.db.Close()
}



