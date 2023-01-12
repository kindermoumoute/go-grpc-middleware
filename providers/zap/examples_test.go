// Copyright (c) The go-grpc-middleware Authors.
// Licensed under the Apache License 2.0.

package zap_test

import (
	"context"
	"testing"
	"time"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/providers/zap/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var (
	customFunc             logging.CodeToLevel
	customDurationToFields logging.DurationToFields
)

func Example_initializationWithCustomLevels() {
	// Logger is used, allowing pre-definition of certain fields by the user.
	logger := zap.NewNop()
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []logging.Option{
		logging.WithLevels(customFunc),
	}
	// Create a server, make sure we put the tags context before everything else.
	_ = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczap.InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(grpczap.InterceptorLogger(logger), opts...),
		),
	)
}

func Example_initializationWithDurationFieldOverride() {
	// Logger is used, allowing pre-definition of certain fields by the user.
	logger := zap.NewNop()
	// Shared options for the logger, with a custom duration to log field function.
	opts := []logging.Option{
		logging.WithDurationField(customDurationToFields),
	}
	// Create a server, make sure we put the tags context before everything else.
	_ = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczap.InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(grpczap.InterceptorLogger(logger), opts...),
		),
	)
}

func ExampleWithDecider() {
	// Logger is used, allowing pre-definition of certain fields by the user.
	logger := zap.NewNop()
	// Shared options for the logger, with a custom decider that log everything except successful calls from "/blah.foo.healthcheck/Check" method.
	opts := []logging.Option{
		logging.WithDecider(func(methodFullName string, _ error) logging.Decision {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if methodFullName == "/blah.foo.healthcheck/Check" {
				return logging.NoLogCall
			}

			// by default you will log all calls
			return logging.LogStartAndFinishCall
		}),
	}
	// Create a server, make sure we put the tags context before everything else.
	_ = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczap.InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(grpczap.InterceptorLogger(logger), opts...),
		),
	}
}

func ExampleServerPayloadLoggingDecider() {
	// Logger is used, allowing pre-definition of certain fields by the user.
	logger := zap.NewNop()
	// Expect payload from  "/blah.foo.healthcheck/Check" call to be logged.
	payloadDecider := func(ctx context.Context, fullMethodName string, servingObject interface{}) logging.PayloadDecision {
		if fullMethodName == "/blah.foo.healthcheck/Check" {
			return logging.LogPayloadRequestAndResponse
		}
		return logging.NoPayloadLogging
	}

	// Create a server, make sure we put the tags context before everything else.
	_ = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczap.InterceptorLogger(logger)),
			logging.PayloadUnaryServerInterceptor(grpczap.InterceptorLogger(logger), payloadDecider, time.RFC3339),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(grpczap.InterceptorLogger(logger)),
			logging.PayloadStreamServerInterceptor(grpczap.InterceptorLogger(logger), payloadDecider, time.RFC3339),
		),
	}
}

func TestExamplesBuildable(t *testing.T) {
	Example_initializationWithCustomLevels()
	Example_initializationWithDurationFieldOverride()
	ExampleWithDecider()
	ExampleServerPayloadLoggingDecider()
}
