package main

import (
	"github.com/Etpmls/EM-Attachment/v3/database"
	"github.com/Etpmls/EM-Attachment/v3/proto/pb"
	"github.com/Etpmls/EM-Attachment/v3/service"
	em "github.com/Etpmls/Etpmls-Micro/v3"
	em_define "github.com/Etpmls/Etpmls-Micro/v3/define"
	"google.golang.org/grpc"
)

const (
	AppName    = "EM-Attachment"
	AppVersion = "v3.0.0-beta"
)

func main()  {
	var reg = em.Register{
		Version:                map[string]string{AppName + " Version": AppVersion},
		EnabledFeature:     []string{
			em_define.EnableValidator,
			em_define.EnableTranslate,
			em_define.EnableCircuitBreaker,
			em_define.EnableServiceDiscovery,
		},
		RegisterService: func(s *grpc.Server) {
			pb.RegisterAttachmentServer(s, &service.ServiceAttachment{})
		},
		OverrideInterface:         em.OverrideInterface{},
		OverrideFunction:          em.OverrideFunction{},
	}
	reg.Init()
	database.NewDatabase().Init()
	reg.Run()
}