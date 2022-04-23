package handler

import (
	"context"
	svc "user-unary/service"

	upb "grpc-applications/protoc/protobuf-user"

	"go.uber.org/zap"
)

type handler struct {
	logger  *zap.Logger
	service svc.UnaryService
}

func NewUnaryHandler(l *zap.Logger, s svc.UnaryService) *handler {
	return &handler{service: s, logger: l}
}

func (h *handler) SaveUserFunc(ctx context.Context, r *upb.SaveUserRequest) (*upb.SaveUserResponse, error) {
	h.logger.Info("Message was applied!")

	resp, err := h.service.SaveNewUser(ctx, r.User.Credentials, r.User.GetUserName(), r.User.GetRoot())
	if err != nil {
		h.logger.With(zap.Error(err)).Error("failed to process")
		return resp, err
	}

	h.logger.Info("Message was successfully processed!")
	return resp, nil
}

func (h *handler) DeleteUserFunc(ctx context.Context, r *upb.DeleteUserRequest) (*upb.DeleteUserResponse, error) {
	h.logger.Info("Message was applied!")

	resp, err := h.service.DeleteUser(ctx, r)
	if err != nil {
		h.logger.With(zap.Error(err)).Error("failed to process")
		return resp, err
	}
	h.logger.Info("Message was successfully processed")
	return resp, nil
}

func (h *handler) ReadSpecificUsersFunc(r *upb.ReadSpecificUsersRequest,
	stream upb.UserService_ReadSpecificUsersFuncServer) error {
	h.logger.Info("Stream for user update is opened!")

	for _, u := range r.UserIds {
		res, err := h.service.GetUser(stream.Context(), u)
		if err != nil {
			h.logger.With(zap.Error(err), zap.String("user_id", u)).Error("Failed to retrieve user")
			continue
		}
		if err := stream.Send(res); err != nil {
			h.logger.With(zap.Error(err)).Error("failed to send result")
			continue
		}
	}

	h.logger.Info("Message was successfully processed!")
	return nil
}
