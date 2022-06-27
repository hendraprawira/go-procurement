package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hendraprawira/go-procurement/database"
	"github.com/hendraprawira/go-procurement/graph/generated"
	"github.com/hendraprawira/go-procurement/graph/model"
)

func (r *authOpsResolver) Login(ctx context.Context, obj *model.AuthOps, creds model.CredsLogin) (interface{}, error) {
	return db.Login(creds)
}

func (r *authOpsResolver) Register(ctx context.Context, obj *model.AuthOps, input model.NewUser) (interface{}, error) {
	return db.CreateUser(input), nil
}

func (r *mutationResolver) CreateItem(ctx context.Context, input model.NewItem) (*model.Item, error) {
	return db.Save(input), nil
}

func (r *mutationResolver) DeleteItem(ctx context.Context, id string) (*bool, error) {
	return db.Delete(id), nil
}

func (r *mutationResolver) Auth(ctx context.Context) (*model.AuthOps, error) {
	return &model.AuthOps{}, nil
}

func (r *queryResolver) Items(ctx context.Context, input *model.FilterItem) ([]*model.Item, error) {
	if input == nil {
		input = &model.FilterItem{}
	}
	return db.Find(input), nil
}

func (r *queryResolver) Item(ctx context.Context, id string) (*model.Item, error) {
	return db.FindByID(id), nil
}

// AuthOps returns generated.AuthOpsResolver implementation.
func (r *Resolver) AuthOps() generated.AuthOpsResolver { return &authOpsResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type authOpsResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var db = database.Connect()

type itemsResolver struct{ *Resolver }
