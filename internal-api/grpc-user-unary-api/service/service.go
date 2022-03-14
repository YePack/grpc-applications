package service

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"grpc-applications/crud-mongodb"
	upb "grpc-applications/protoc/protobuf-user"
	"time"
)

const timeLayout = `2006-01-02 15:04:25 Monday`

type service struct {
	logger   *zap.Logger
	dbClient crud_mongodb.MongoDBInterface
}

type UnaryService interface {
	SaveNewUser(ctx context.Context, r *upb.SaveUserRequest) (*upb.SaveUserResponse, error)
	DeleteUser(ctx context.Context, r *upb.DeleteUserRequest) (*upb.DeleteUserResponse, error)
}

func NewUnaryService(logger *zap.Logger, client crud_mongodb.MongoDBInterface) UnaryService {
	return &service{logger: logger, dbClient: client}
}

func (s *service) SaveNewUser(ctx context.Context, r *upb.SaveUserRequest) (*upb.SaveUserResponse, error) {
	s.logger.Debug("Processing new user saving!")

	r.User.UserId = uuid.New().String()
	r.User.RegisterDate = time.Now().Format(timeLayout)
	if err := s.dbClient.InsertUser(ctx, r); err != nil {
		return &upb.SaveUserResponse{UserId: "", Saved: false}, err
	}

	s.logger.Debug("New user was processed!")
	return &upb.SaveUserResponse{
		UserId: r.User.UserId,
		Saved:  true,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, r *upb.DeleteUserRequest) (*upb.DeleteUserResponse, error) {
	s.logger.Debug("Processing user for delete!")

	if err := s.dbClient.DeleteUser(ctx, r.UserId); err != nil {
		return &upb.DeleteUserResponse{Deleted: false}, err
	}

	s.logger.Debug("Requested user was processed!")
	return &upb.DeleteUserResponse{Deleted: true}, nil
}
