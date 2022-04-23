package service

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"grpc-applications/crud-mongodb"
	m "grpc-applications/model"
	upb "grpc-applications/protoc/protobuf-user"
	"time"
)

const timeLayout = `2006-01-02 15:04:25 Monday`

type service struct {
	logger   *zap.Logger
	dbClient crud_mongodb.MongoDBInterface
}

type UnaryService interface {
	SaveNewUser(ctx context.Context, creds *upb.Credentials, userName string, root bool) (*upb.SaveUserResponse, error)
	DeleteUser(ctx context.Context, r *upb.DeleteUserRequest) (*upb.DeleteUserResponse, error)

	// GetUser - Stream
	GetUser(ctx context.Context, userId string) (*upb.ReadUserResponse, error)
}

func NewUnaryService(logger *zap.Logger, client crud_mongodb.MongoDBInterface) UnaryService {
	return &service{logger: logger, dbClient: client}
}

func (s *service) SaveNewUser(ctx context.Context, creds *upb.Credentials, userName string, root bool) (*upb.SaveUserResponse, error) {
	s.logger.Debug("Processing new user saving!")

	user := &m.User{
		UserId:   uuid.New().String(),
		UserName: userName,
		Credentials: m.Credentials{
			Login:    creds.GetLogin(),
			Password: creds.GetPassword(),
		},
		RegisterDate: time.Now().Format(timeLayout),
		Root:         root,
	}

	if err := s.dbClient.InsertUser(ctx, user); err != nil {
		return &upb.SaveUserResponse{UserId: "", Saved: false}, err
	}

	s.logger.Debug("New user was processed!")
	return &upb.SaveUserResponse{
		UserId: user.UserId,
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

func (s *service) GetUser(ctx context.Context, userId string) (*upb.ReadUserResponse, error) {
	s.logger.Debug("Processing for retrieving specific user")

	user, err := s.dbClient.ReadUser(ctx, userId)
	if err != nil {
		s.logger.With(zap.Error(err))
		return nil, err
	}

	s.logger.Debug("Requested user was retrieved!")
	return &upb.ReadUserResponse{User: &upb.User{
		UserId:       userId,
		UserName:     user.UserName,
		Credentials:  &upb.Credentials{Login: user.Credentials.Login, Password: user.Credentials.Password},
		RegisterDate: user.RegisterDate,
		Root:         user.Root,
	}}, nil
}
