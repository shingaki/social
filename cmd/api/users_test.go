package main

import (
	"SOCIAL/internal/store/cache"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {

	t.Setenv("REDIS_ENABLED", "true")

	t.Logf("REDIS_ENABLED raw value: %q", os.Getenv("REDIS_ENABLED"))

	withRedis := config{
		redisCfg: redisConfig{
			enabled: true,
		},
	}

	t.Logf("before creating the new app withRedis: %v", withRedis)

	app := newTestApplication(t, withRedis)
	mux := app.mount()

	t.Logf("redis enabled from config: %v", app.config.redisCfg.enabled)

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		// check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/140", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authorized requests", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		mockCacheStore.On("Get", int64(140)).
			Return(nil, nil)

		mockCacheStore.On("Set", mock.Anything, mock.Anything).
			Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/140", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})

	t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {

		log.Printf("withRedis enabled: %v", withRedis)
		log.Printf("config enabled: %v", app.config.redisCfg.enabled)

		withRedis := config{
			redisCfg: redisConfig{
				enabled: true,
			},
		}

		app := newTestApplication(t, withRedis)
		mux := app.mount()

		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		mockCacheStore.On("Get", int64(139)).Return(nil, nil)
		mockCacheStore.On("Get", int64(140)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/140", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)

		mockCacheStore.Calls = nil // Reset mock expectations
	})

	t.Run("should NOT hit the cache if it is not enabled", func(t *testing.T) {
		withRedis := config{
			redisCfg: redisConfig{
				enabled: false,
			},
		}

		app := newTestApplication(t, withRedis)
		mux := app.mount()

		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/140", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNotCalled(t, "Get")

		mockCacheStore.Calls = nil // Reset mock expectations
	})
}
