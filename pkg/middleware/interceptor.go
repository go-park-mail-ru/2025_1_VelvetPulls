package middleware

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, utils.REQUEST_ID_KEY, requestID)
		return handler(ctx, req)
	}
}

func AccessLogInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		requestID := utils.GetRequestIDFromCtx(ctx)
		logger := utils.Logger.With(
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
		)
		ctx = utils.WithLogger(ctx, logger)

		resp, err := handler(ctx, req)

		fields := []zap.Field{
			zap.Duration("execution_time", time.Since(start)),
		}

		if err != nil {
			logger.Error("gRPC request failed", append(fields, zap.Error(err))...)
		} else {
			logger.Info("gRPC request completed", fields...)
		}

		return resp, err
	}
}
