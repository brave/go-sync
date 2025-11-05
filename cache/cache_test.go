package cache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/brave/go-sync/cache"
)

type CacheTestSuite struct {
	suite.Suite
	cache *cache.Cache
}

func (suite *CacheTestSuite) SetupSuite() {
	suite.cache = cache.NewCache(cache.NewRedisClient())
}

func (suite *CacheTestSuite) TestSetTypeMtime() {
	suite.cache.SetTypeMtime(context.Background(), "id", 123, 12345678)
	val, err := suite.cache.Get(context.Background(), "id#123", false)
	suite.Require().NoError(err)
	suite.Require().Equal("12345678", val)
}

func (suite *CacheTestSuite) TestIsTypeMtimeUpdated() {
	suite.cache.SetTypeMtime(context.Background(), "id", 123, 12345678)

	tests := map[string]struct {
		clientID string
		dataType int
		mtime    int64
		ret      bool
	}{
		"cache miss for mismatch clientID": {
			clientID: "id2",
			dataType: 123,
			mtime:    12345677,
			ret:      true,
		},
		"cache miss for mismatch dataType": {
			clientID: "id",
			dataType: 456,
			mtime:    12345677,
			ret:      true,
		},
		"cache hit with older client token": {
			clientID: "id",
			dataType: 123,
			mtime:    12345677,
			ret:      true,
		},
		"cache hit with the same client token": {
			clientID: "id",
			dataType: 123,
			mtime:    12345678,
			ret:      false,
		},
		"cache hit with the newer client token": {
			clientID: "id",
			dataType: 123,
			mtime:    12345679,
			ret:      false,
		},
	}

	for testName, test := range tests {
		ret := suite.cache.IsTypeMtimeUpdated(
			context.Background(), test.clientID, test.dataType, test.mtime)
		suite.Require().Equal(test.ret, ret,
			"unexpected return value for %s test case", testName)
	}
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
