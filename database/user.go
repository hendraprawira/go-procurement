package database

import (
	"context"
	"log"
	"time"

	"github.com/hendraprawira/Tripatra-procurement/graph/model"
	"github.com/hendraprawira/Tripatra-procurement/service"
	"github.com/hendraprawira/Tripatra-procurement/tools"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *DB) CreateUser(input model.NewUser) *model.User {
	print("Here")
	collection := db.client.Database("procurementTripatra").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	input.Password = tools.HashPassword(input.Password)
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	return &model.User{
		ID:    res.InsertedID.(primitive.ObjectID).Hex(),
		Name:  input.Name,
		Email: input.Email,
	}
}

func (db *DB) GetByEmail(email string) *model.User {
	collection := db.client.Database("procurementTripatra").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"email": email})
	user := model.User{}
	res.Decode(&user)
	return &user
}
func (db *DB) Login(creds model.CredsLogin) (interface{}, error) {
	user := db.GetByEmail(creds.Email)
	if user == nil {
		return nil, &gqlerror.Error{
			Message: "Email not found",
		}
	}

	var pw string
	pw = *user.Password

	if err := tools.ComparePassword(pw, creds.Password); err != nil {
		return nil, err
	}

	token, err := service.JwtGenerate(user.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}
