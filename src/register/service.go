package register

import (
	"context"
	"github.com/Etpmls/EM-Attachment/src/application/protobuf"
	"github.com/Etpmls/EM-Attachment/src/application/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Register Rpc Service
func RegisterRpcService(s *grpc.Server)  {
	// protobuf.RegisterUserServer(s, &service.ServiceUser{})
	protobuf.RegisterAttachmentServer(s, &service.ServiceAttachment{})
	return
}

// Register Http Service
func RegisterHttpService(ctx context.Context, mux *runtime.ServeMux, grpcServerEndpoint *string, opts []grpc.DialOption) error {
	/*err := protobuf.RegisterUserHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}*/
	err := protobuf.RegisterAttachmentHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	return nil
}
