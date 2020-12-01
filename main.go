package main

import (
	"github.com/Etpmls/EM-Attachment/src/application/database"
	"github.com/Etpmls/EM-Attachment/src/register"
	"github.com/Etpmls/EM-Attachment/src/register/config"
	"github.com/Etpmls/Etpmls-Micro"
)

func main() {
	var reg = em.Register{
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
		}{Path: "storage/config/attachment.yaml", DebugPath: "storage/config/attachment_debug.yaml", StructAddr: &register_config.ServiceConfig},
	}

	reg.Run()
}
