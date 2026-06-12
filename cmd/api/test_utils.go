package main

import (
	"SOCIAL/internal/auth"
	"SOCIAL/internal/store"
	"SOCIAL/internal/store/cache"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// Uncomment to enable logs
	// logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	log.Printf("cfg creating the app: %v", cfg)

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	log.Printf("Execute for user %v", req)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
