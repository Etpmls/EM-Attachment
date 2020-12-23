package main

import (
	"github.com/Etpmls/EM-Attachment/src/application"
	"github.com/Etpmls/EM-Attachment/src/application/database"
	"github.com/Etpmls/EM-Attachment/src/register"
	"github.com/Etpmls/Etpmls-Micro"
)

func main() {
	var reg = em.Register{
		Version_Service: 		map[string]string{"EM-Attachment Version": application.Version_Service},
		GrpcServiceFunc:    	register.RegisterRpcService,
		HttpServiceFunc:    	register.RegisterHttpService,
		RouteFunc:          	register.RegisterRoute,
		DatabaseMigrate:		[]interface{}{
			&database.Attachment{},
		},
		CustomConfiguration: struct {
			Path       string
			DebugPath  string
			StructAddr interface{}
		}{Path: "storage/config/attachment.yaml", DebugPath: "storage/config/attachment_debug.yaml", StructAddr: &application.ServiceConfig},
	}
	reg.Init()
	reg.Run()
}
