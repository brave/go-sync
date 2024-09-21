package datastore

import (
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"strconv"
	"strings"
)

// SQLSaveRolloutsEnvKey defines the data types and rollout percentages for saving
// new items into the SQL database, instead of Dynamo.
const SQLSaveRolloutsEnvKey = "SQL_SAVE_ROLLOUTS"

// SQLSaveRolloutsEnvKey defines the data types and rollout percentages for periodic
// chunked migration from Dynamo to SQL.
const SQLMigrateRolloutsEnvKey = "SQL_MIGRATE_ROLLOUTS"

// VariationHashDecimal returns a decimal from 0.0 to 1.0 for a given client ID.
// The decimal is typically checked against a rollout percentage to determine if a user
// should be included in a rollout.
func VariationHashDecimal(input string) float32 {
	h := fnv.New32a()
	h.Write([]byte(input))
	hashValue := h.Sum32()

	// Convert hash to a decimal between 0 and 1
	return float32(hashValue) / math.MaxUint32
}

// SQLVariations handles SQL variation rollout functions
type SQLVariations struct {
	sqlSaveRollouts    map[int]float32
	sqlMigrateRollouts map[int]float32
	Ready              bool
}

func parseRollouts(envKey string) (map[int]float32, error) {
	rollouts := make(map[int]float32)
	envVal := os.Getenv(envKey)

	if len(envVal) > 0 {
		pairs := strings.Split(envVal, ",")

		for _, pair := range pairs {
			parts := strings.Split(strings.TrimSpace(pair), "=")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid format in %s: %s", envKey, pair)
			}

			key, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("Invalid integer in %s: %s", envKey, parts[0])
			}

			value, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 32)
			if err != nil {
				return nil, fmt.Errorf("Invalid float in %s: %s", envKey, parts[1])
			}

			rollouts[key] = float32(value)
		}
	}

	return rollouts, nil
}

// LoadSQLVariations creates a SQLVariations struct, configured by env vars
func LoadSQLVariations() (*SQLVariations, error) {
	sqlSaveRollouts, err := parseRollouts(SQLSaveRolloutsEnvKey)
	if err != nil {
		return nil, err
	}
	sqlMigrateRollouts, err := parseRollouts(SQLMigrateRolloutsEnvKey)
	if err != nil {
		return nil, err
	}

	return &SQLVariations{
		sqlSaveRollouts:    sqlSaveRollouts,
		sqlMigrateRollouts: sqlMigrateRollouts,
		Ready:              false,
	}, nil
}

// ShouldSaveToSQL returns true if a client should save the entity to the SQL database for a given data type
func (sqlVariations *SQLVariations) ShouldSaveToSQL(dataType int, variationHashDecimal float32) bool {
	rolloutPercent, exists := sqlVariations.sqlSaveRollouts[dataType]
	return exists && variationHashDecimal <= rolloutPercent
}

// ShouldMigrateToSQL returns true if chunked migration from Dynamo to SQL should occur for a given data type
func (sqlVariations *SQLVariations) ShouldMigrateToSQL(dataType int, variationHashDecimal float32) bool {
	rolloutPercent, exists := sqlVariations.sqlMigrateRollouts[dataType]
	return exists && variationHashDecimal <= rolloutPercent
}

// GetStateDigest returns a string that combines the env vars related to variations
func (sqlVariations *SQLVariations) GetStateDigest() string {
	return SQLSaveRolloutsEnvKey + ":" + os.Getenv(SQLSaveRolloutsEnvKey) + ";" +
		SQLMigrateRolloutsEnvKey + ":" + os.Getenv(SQLMigrateRolloutsEnvKey)
}

func (sqlDB *SQLDB) Variations() *SQLVariations {
	return sqlDB.variations
}
