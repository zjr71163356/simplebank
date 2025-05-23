package gapi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zjr71163356/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authorizationHeaderValue := md.Get(authorizationHeader)
	if len(authorizationHeaderValue) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeaderValue[0])
	if len(fields) < 2 {
		return nil, errors.New("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, errors.New("unsupported authorization type")
	}

	payload, err := server.tokenMaker.VerifyToken(fields[1])
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)

	}

	return payload, nil
}
