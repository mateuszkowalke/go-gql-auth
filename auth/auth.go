package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"example.com/go-graphql-auth/graph/model"
	"example.com/go-graphql-auth/jwt"
	"example.com/go-graphql-auth/users"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func Middleware(next http.Handler) http.Handler {
	fmt.Println("---------------", next, "----------------")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}
		//validate jwt token
		tokenStr := header
		email, err := jwt.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		// create user and check if user exists in db
		user := model.User{Email: email}
		id, err := users.GetUserIdByEmail(email)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user.ID = strconv.Itoa(id)
		// put it in context
		ctx := context.WithValue(r.Context(), userCtxKey, &user)
		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// FromContext finds the user from the context. REQUIRES Middleware to have run.
func FromContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}
