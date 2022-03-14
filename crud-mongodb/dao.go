package crud_mongodb

import (
	"context"
	"errors"
	"fmt"
	upb "grpc-applications/protoc/protobuf-user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

const (
	name            = "mongodb://%s:%s"
	dbName          = "applications"
	usersCollection = "users"
)

type mongoClient struct {
	client *mongo.Client
	logger *zap.Logger
}

type MongoDBInterface interface {
	InsertUser(ctx context.Context, r *upb.SaveUserRequest) error
	GetUsers(ctx context.Context) ([]*upb.User, error)
	ReadUser(ctx context.Context, userID string) (*upb.User, error)
	DeleteUser(ctx context.Context, userID string) error
	UpdateUser(ctx context.Context, r *upb.UpdateUserRequest) error
}

func NewMongoConnection(ctx context.Context, logger *zap.Logger, host, port string) (MongoDBInterface, error) {
	logger.Debug("Creation new connection at MongoDB!")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf(name, host, port)))
	if err != nil {
		logger.With(zap.Error(err)).Error("mongoDB connection error")
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logger.With(zap.Error(err)).Error("ping failed")
		return nil, err
	}

	logger.Debug("New MongoDB client was successfully created!")
	return &mongoClient{
		client: client,
		logger: logger,
	}, nil
}

func (c *mongoClient) InsertUser(ctx context.Context, r *upb.SaveUserRequest) error {
	c.logger.With(zap.Any("item", r)).Debug("Inserting new user!")

	collection := c.client.Database(dbName).Collection(usersCollection)
	if _, err := collection.InsertOne(ctx, r); err != nil {
		return err
	}

	c.logger.With(zap.Any("item", r)).Debug("New user was successfully inserted!")
	return nil
}

func (c *mongoClient) ReadUser(ctx context.Context, userID string) (*upb.User, error) {
	c.logger.With(zap.String("user_id", userID)).Debug("Retrieving data about user with specific ID!")

	var user *upb.User

	collection := c.client.Database(dbName).Collection(usersCollection)
	if err := collection.FindOne(ctx, bson.E{Key: userID}).Decode(&user); err != nil {
		return nil, err
	}

	c.logger.With(zap.String("user_id", userID)).Debug("User data was successfully parsed!")
	return user, nil
}

func (c *mongoClient) GetUsers(ctx context.Context) ([]*upb.User, error) {
	c.logger.Debug("Retrieving registered users!")

	var (
		user   upb.User
		result []*upb.User
	)

	collection := c.client.Database(dbName).Collection(usersCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		result = append(result, &user)
	}
	if cursor.Err() != nil {
		return nil, err
	}

	c.logger.Debug("Registered users were successfully parsed!")
	return result, nil
}

func (c *mongoClient) DeleteUser(ctx context.Context, userID string) error {
	c.logger.With(zap.String("user_id", userID)).Debug("Deleting user with specific ID!")

	collection := c.client.Database(dbName).Collection(usersCollection)
	res, err := collection.DeleteOne(ctx, bson.M{"user.userid": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("user with requested id wasn't found")
	}

	c.logger.Debug(fmt.Sprintf("Rows count deleted: %d!", res.DeletedCount))
	return nil
}

func (c *mongoClient) UpdateUser(ctx context.Context, r *upb.UpdateUserRequest) error {
	c.logger.With(zap.String("user_id", r.User.UserId)).Debug("Updating specific user data!")

	collection := c.client.Database(dbName).Collection(usersCollection)
	res, err := collection.UpdateOne(ctx, bson.M{"user.userid": r.User.UserId}, r)
	if err != nil {
		return err
	}
	if res.UpsertedCount < 1 || res.MatchedCount < 1 || res.ModifiedCount < 1 {
		c.logger.Debug("User was successfully updated!")
		return nil
	}
	c.logger.Debug("User was successfully updated!")
	return nil
}
