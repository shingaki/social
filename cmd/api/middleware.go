package main

import (
	"SOCIAL/internal/store"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
		authHeader := r.Header.Get("Authorization")
		log.Printf("authHeader %s:", r.Header.Get("Authorization"))
		log.Println(r.Header)
		if authHeader == "" {
			log.Println("missing authorization header 0:")
			app.unauthorizedResponse(w, r, fmt.Errorf("missing authorization header"))
			return
		}

		// parse the token to get the base64
		parts := strings.Split(authHeader, " ") // authorization: Bearer <token>
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		// validate the token
		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			log.Println("validate token error:")
			app.unauthorizedResponse(w, r, err)
			return
		}

		// parse the claim from the token and get the userID
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			log.Println("parse claims error:")

			app.unauthorizedResponse(w, r, err)
			return
		}

		// fetch the user from the dB
		ctx := r.Context()

		// call function to get user from cache instead of database
		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		// pick the user and put it in the context of the http request
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			log.Println("basic auth middleware")
			authHeader := r.Header.Get("Authorization")
			log.Printf("authHeader %s:", r.Header.Get("Authorization"))
			if authHeader == "" {
				log.Println("missing authorization header 2:")
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("missing authorization header"))
				return
			}

			// parse it to get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicResponse(w, r, err)
				return
			}

			// check the credentials against the credentials that are stored in the configuration, but we could have used a database
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Println(name, ":", value)
		}
	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromCtx(r)

		// if it is the users post
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)

	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}

// function to get user from cache
func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {

	log.Printf("middleware enabled %v ", app.config.redisCfg.enabled)
	if !app.config.redisCfg.enabled {
		log.Println("do not read cache")
		return app.store.Users.GetByID(ctx, userID)
	}

	// check cache first
	log.Printf("checking cache for user %v: ", userID)
	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// user is not on the cache then get user from the database
	if user == nil {
		log.Println("user not in cache, fetching from database")
		user, err = app.store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Set the user in the cache
		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("rate limiter middleware http.Request %v", r.RemoteAddr)
		log.Printf("rate limiter middleware enabled %v", app.config.rateLimiter.Enabled)
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
