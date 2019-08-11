package database

import (
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestPing_onValidDB() {
	err := s.db.Ping()
	assert.NoError(s.T(), err)
}

func (s *DatabaseSuite) TestPing_onClosedDB() {
	s.db.Close()
	err := s.db.Ping()
	assert.Error(s.T(), err)
}
