package server

import (
	"context"
	"os"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/rs/zerolog/log"
)

const (
	lastRolloutStateCacheKey string = "last-rollout-state"
	rolloutConfirmChannelKey string = "rollout-confirm"
	sqlDisableRolloutConfirm string = "SQL_DISABLE_ROLLOUT_CONFIRM"
)

func maybeWaitOnRolloutConfigChange(sqlVariations *datastore.SQLVariations, cache *cache.Cache) {
	currentDigest := sqlVariations.GetStateDigest()

	lastDigest, err := cache.Get(context.Background(), lastRolloutStateCacheKey, false)
	if err != nil {
		log.Fatal().Msgf("failed to get last rollout state: %v", err)
		return
	}

	rolloutConfirmDisabled := os.Getenv(sqlDisableRolloutConfirm) != ""
	if !rolloutConfirmDisabled && currentDigest != lastDigest {
		log.Info().Msg("Rollout configuration detected. Commits/writes disabled until Redis confirmation event is received...")
		err = cache.SubscribeAndWait(context.Background(), rolloutConfirmChannelKey)
		if err != nil {
			log.Fatal().Msgf("failed to subscribe and wait for rollout confirmation: %v", err)
			return
		}

		err = cache.Set(context.Background(), lastRolloutStateCacheKey, currentDigest, 0)
		if err != nil {
			log.Fatal().Msgf("failed to update last rollout state: %v", err)
			return
		}
		log.Info().Msg("Confirmation event received")
	}

	sqlVariations.Ready = true
}
