package handler

import (
	"context"

	svc "user-unary/service"

	upb "grpc-applications/protoc/protobuf-user"

	"go.uber.org/zap"
)

type handler struct {
	service svc.UnaryService
	logger  *zap.Logger
}

func NewUnaryHandler(l *zap.Logger, s svc.UnaryService) *handler {
	return &handler{service: s, logger: l}
}

func (h *handler) SaveUserFunc(ctx context.Context, r *upb.SaveUserRequest) (*upb.SaveUserResponse, error) {
	h.logger.Info("Message was applied!")

	resp, err := h.service.SaveNewUser(ctx, r)
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

func (h *handler) UpdateUserFunc(ctx context.Context, r *upb.UpdateUserRequest) (*upb.UpdateUserResponse, error) {
	return nil, nil
}
