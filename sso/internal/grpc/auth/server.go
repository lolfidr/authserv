package auth

import (
	"context"
	"errors"

	ssov1 "github.com/lolfidr/authserv/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth *Auth
}

func Register(gRPCServer *grpc.Server, auth *Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLoginRequest(in); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		return handleLoginError(err)
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func validateLoginRequest(in *ssov1.LoginRequest) error {
	if in.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if in.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func handleLoginError(err error) (*ssov1.LoginResponse, error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, ErrAppNotFound):
		return nil, status.Error(codes.NotFound, "app not found")
	default:
		return nil, status.Error(codes.Internal, "failed to login")
	}
}
