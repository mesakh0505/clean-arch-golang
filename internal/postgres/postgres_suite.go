package postgres

import (
	"database/sql"

	// This is imported for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/suite"
)

const (
	postgres = "postgres"
)

// Suite struct for MySQL Suite
type Suite struct {
	suite.Suite
	DSN                     string
	DBConn                  *sql.DB
	Migration               *migration
	MigrationLocationFolder string
	DBName                  string
}

// SetupSuite setup at the beginning of test
func (s *Suite) SetupSuite() {
	var err error
	s.DBConn, err = sql.Open(postgres, s.DSN)
	s.Require().NoError(err)
	err = s.DBConn.Ping()
	s.Require().NoError(err)
	s.Migration, err = runMigration(s.DBConn, s.MigrationLocationFolder)
	s.Require().NoError(err)
}

// TearDownSuite teardown at the end of test
func (s *Suite) TearDownSuite() {
	err := s.DBConn.Close()
	s.Require().NoError(err)
}
