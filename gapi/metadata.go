package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type MetaData struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *MetaData {
	mtdt := &MetaData{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// fmt.Printf("metadata: %v\n", md)
		if UserAgents := md.Get("user-agent"); len(UserAgents) > 0 {
			mtdt.UserAgent = UserAgents[0]
		}

		if UserAgents := md.Get("grpcgateway-user-agent"); len(UserAgents) > 0 {
			mtdt.UserAgent = UserAgents[0]
		}

		if ClientIP := md.Get("x-forwarded-for"); len(ClientIP) > 0 {
			mtdt.ClientIP = ClientIP[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt

}
