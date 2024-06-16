package middleware

import (
	"context"
	"fmt"
	"net/http"

	m "de.server/app/model"
	s "de.server/app/store"
)

type key string

type Middleware func(http.Handler) http.Handler

const (
	role key = "role"
)

func BasicAuth(userStore *s.CrudStore[m.User]) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				http.Error(w, "basic auth not set", http.StatusUnauthorized)
				return
			}
			userModel := m.User{
				Email:    &username,
				Password: &password,
			}
			userRole, err := userStore.FindOne(&userModel, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), role, userRole.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithRole(roles []m.Role) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value(role)
			fmt.Println(userRole)
			if len(roles) == 0 {
				next.ServeHTTP(w, r)
			}
			for _, acceptedRole := range roles {
				if userRole == acceptedRole {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "you don't have the role to access this endpoint", http.StatusUnauthorized)
		})
	}
}

func GetMiddlewares(handler []Middleware) []func(http.Handler) http.Handler {
	var handlers []func(http.Handler) http.Handler
	for _, h := range handler {
		handlers = append(handlers, h)
	}
	return handlers
}
