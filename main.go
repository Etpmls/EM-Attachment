package main

import (
	"github.com/Etpmls/EM-Attachment/src/application"
	"github.com/Etpmls/EM-Attachment/src/application/database"
	"github.com/Etpmls/EM-Attachment/src/register"
	"github.com/Etpmls/Etpmls-Micro/v2"
)

func main() {
	var reg = em.Register{
		AppVersion: 		map[string]string{"EM-Attachment Version": application.Version_Service},
		AppEnabledFeatureName:		[]string{em.EnableCircuitBreaker, em.EnableDatabase, em.EnableI18n, em.EnableServiceDiscovery, em.EnableValidator},
		RpcServiceFunc:    	register.RegisterRpcService,
		HttpServiceFunc:    	register.RegisterHttpService,
		HttpRouteFunc:          	register.RegisterRoute,
		DatabaseMigrate:		[]interface{}{
			&database.Attachment{},
		},
	}
	reg.Init()
	reg.Run()
}
