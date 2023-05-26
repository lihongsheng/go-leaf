package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

// auth 依赖于client需要以x-md-global格式注入头信息进行传递
// client 的例子 ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
//func auth(config conf.Secret) middleware.Middleware {
//	return func(handler middleware.Handler) middleware.Handler {
//		return func(ctx context.Context, req interface{}) (interface{}, error) {
//			pathMap := make(map[string]struct{})
//			pathMap["/api.user.v1.User/Get"] = struct{}{}
//			if !config.Switch {
//				return handler(ctx, req)
//			}
//			if tr, ok := transport.FromServerContext(ctx); ok {
//				if _, ok := pathMap[tr.Operation()]; ok {
//					if md, exist := md.FromServerContext(ctx); exist {
//						sign := md.Get(types.MetaAppSign)
//						ts := md.Get(types.MetaAppTime)
//						appID := md.Get(types.MetaAppID)
//						if sign == "" || appID == "" {
//							return nil, apiError.ErrorUnauthorized("auth fail,sign is empty")
//						}
//						if secret, exist := config.AppSecret[appID]; exist {
//							if tools.Md5(secret+ts) != sign {
//								return nil, apiError.ErrorUnauthorized("auth fail,sign is fail.")
//							}
//						} else {
//							return nil, apiError.ErrorUnauthorized("auth fail,appID  is not find secret")
//						}
//					} else {
//						return nil, apiError.ErrorUnauthorized("auth fail,not find sign")
//					}
//				}
//			}
//			return handler(ctx, req)
//		}
//	}
//}
