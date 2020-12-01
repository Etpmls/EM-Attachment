package middleware

import (
	"context"
	"github.com/Etpmls/EM-Attachment/src/application"
	em "github.com/Etpmls/Etpmls-Micro"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	em_utils "github.com/Etpmls/Etpmls-Micro/utils"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)



func (this *middleware) Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// fullMethodName: /protobuf.User/GetCurrent
		service := em_library.NewGrpc().GetServiceName(info.FullMethod)

		// Get token
		// 获取令牌
		g := em_library.NewGrpc()
		token, err := g.ExtractHeader(ctx, "token")
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_GetToken"))
		}

		// Get Claims
		// 获取Claims
		tmp, err := em_library.JwtToken.ParseToken(token)
		tk, ok := tmp.(*jwt.Token)
		if !ok || err != nil {
			return nil, status.Error(codes.Unauthenticated, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_TokenVerificationFailed"))
		}

		// Determine whether the role has the corresponding permissions
		// 判断所属角色是否有相应的权限
		if claims,ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
			if userId, ok := claims["jti"].(string); ok {
				id, err := strconv.Atoi(userId)
				if err == nil {
					b := em.NewClient().AuthCheck(application.Service_AuthService, service, uint(id))
					if b {
						// Pass the token to the method
						// 把token传递到方法中
						ctx = context.WithValue(ctx,"token", token)
						return handler(ctx, req)
					}
				} else {
					em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
				}
			}
		}

		return nil, status.Error(codes.InvalidArgument, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
	}
}