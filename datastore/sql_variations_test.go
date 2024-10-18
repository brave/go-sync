package datastore_test

import (
	"os"
	"testing"

	"github.com/brave/go-sync/datastore"
	"github.com/stretchr/testify/suite"
)

type SQLVariationsSuite struct {
	suite.Suite
	variations *datastore.SQLVariations
}

func (s *SQLVariationsSuite) SetupTest() {
	s.T().Setenv(datastore.SQLSaveRolloutsEnvKey, "1=0.5,2=0.75")
	s.T().Setenv(datastore.SQLMigrateRolloutsEnvKey, "1=0.25,3=1.0")
	var err error
	s.variations, err = datastore.LoadSQLVariations()
	s.Require().NoError(err)
}

func (s *SQLVariationsSuite) TestShouldSaveToSQL() {
	s.True(s.variations.ShouldSaveToSQL(1, 0.4))
	s.False(s.variations.ShouldSaveToSQL(1, 0.6))
	s.True(s.variations.ShouldSaveToSQL(2, 0.7))
	s.False(s.variations.ShouldSaveToSQL(2, 0.8))
	s.False(s.variations.ShouldSaveToSQL(3, 0.5)) // Non-existent key
}

func (s *SQLVariationsSuite) TestShouldMigrateToSQL() {
	s.True(s.variations.ShouldMigrateToSQL(1, 0.2))
	s.False(s.variations.ShouldMigrateToSQL(1, 0.3))
	s.True(s.variations.ShouldMigrateToSQL(3, 0.9))
	s.False(s.variations.ShouldMigrateToSQL(2, 0.5)) // Non-existent key
}

func (s *SQLVariationsSuite) TestVariationHashDecimal() {
	hash1 := datastore.VariationHashDecimal("test1")
	hash2 := datastore.VariationHashDecimal("test2")
	s.NotEqual(hash1, hash2)
	s.Less(hash1, float32(1.0))
	s.Less(hash2, float32(1.0))
	s.GreaterOrEqual(hash1, float32(0.0))
	s.GreaterOrEqual(hash2, float32(0.0))
}

func (s *SQLVariationsSuite) TestParseRolloutsError() {
	os.Setenv(datastore.SQLSaveRolloutsEnvKey, "invalid=format")
	_, err := datastore.LoadSQLVariations()
	s.Error(err)
}

func TestSQLVariationsSuite(t *testing.T) {
	suite.Run(t, new(SQLVariationsSuite))
}
