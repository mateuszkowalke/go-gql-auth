package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"example.com/go-graphql-auth/companies"
	"example.com/go-graphql-auth/graph/generated"
	"example.com/go-graphql-auth/graph/model"
	"example.com/go-graphql-auth/jwt"
	"example.com/go-graphql-auth/users"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	newUser := users.User{&input}
	newUser.Create()
	token, err := jwt.GenerateToken(newUser.Email)
	if err != nil {
		return "", err
	}
	// we return default role for the new user which is "user"
	return token, nil
}

func (r *mutationResolver) CreateCompany(ctx context.Context, input model.NewCompany) (*model.Company, error) {
	newCompany := companies.Company{&input}
	companyID := newCompany.Save()
	return &model.Company{
		ID:      strconv.FormatInt(companyID, 10),
		Name:    newCompany.Name,
		Email:   newCompany.Email,
		Country: newCompany.Country,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input *model.Login) (string, error) {
	user := users.User{&model.NewUser{}}
	user.Email = input.Email
	user.Password = input.Password
	authenticated := user.Authenticate()
	if !authenticated {
		// todo
		// sort out the struct/interface to pass
		return "", &users.WrongEmailOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input *model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return users.GetAllUsers()
}

func (r *queryResolver) User(ctx context.Context, email string) (*model.User, error) {
	return users.GetUserByEmail(email)
}

func (r *queryResolver) Companies(ctx context.Context) ([]*model.Company, error) {
	return companies.GetAllCompanies()
}

func (r *queryResolver) Company(ctx context.Context, name string) (*model.Company, error) {
	return companies.GetCompanyByName(name)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
