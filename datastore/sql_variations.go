package datastore

import (
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"strconv"
	"strings"
)

const sqlSaveRolloutsEnvKey = "SQL_SAVE_ROLLOUTS"
const sqlMigrateRolloutsEnvKey = "SQL_MIGRATE_ROLLOUTS"

func VariationHashDecimal(input string) float32 {
	h := fnv.New32a()
	h.Write([]byte(input))
	hashValue := h.Sum32()

	// Convert hash to a decimal between 0 and 1
	return float32(hashValue) / math.MaxUint32
}

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

func LoadSQLVariations() (*SQLVariations, error) {
	sqlSaveRollouts, err := parseRollouts(sqlSaveRolloutsEnvKey)
	if err != nil {
		return nil, err
	}
	sqlMigrateRollouts, err := parseRollouts(sqlMigrateRolloutsEnvKey)
	if err != nil {
		return nil, err
	}

	return &SQLVariations{
		sqlSaveRollouts:    sqlSaveRollouts,
		sqlMigrateRollouts: sqlMigrateRollouts,
		Ready:              false,
	}, nil
}

func (sqlVariations *SQLVariations) ShouldSaveToSQL(dataType int, variationHashDecimal float32) bool {
	rolloutPercent, exists := sqlVariations.sqlSaveRollouts[dataType]
	return exists && variationHashDecimal <= rolloutPercent
}

func (sqlVariations *SQLVariations) ShouldMigrateToSQL(dataType int, variationHashDecimal float32) bool {
	rolloutPercent, exists := sqlVariations.sqlMigrateRollouts[dataType]
	return exists && variationHashDecimal <= rolloutPercent
}

func (sqlVariations *SQLVariations) GetStateDigest() string {
	return sqlSaveRolloutsEnvKey + ":" + os.Getenv(sqlSaveRolloutsEnvKey) + ";" +
		sqlMigrateRolloutsEnvKey + ":" + os.Getenv(sqlMigrateRolloutsEnvKey)
}

func (sqlDB *SQLDB) Variations() *SQLVariations {
	return sqlDB.variations
}
