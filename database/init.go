package database

import (
	em "github.com/Etpmls/Etpmls-Micro/v3"
	"github.com/Etpmls/Etpmls-Micro/v3/define"
	em_library "github.com/Etpmls/Etpmls-Micro/v3/library"
	"strings"
)

const (
	// Service Database
	KvServiceDatabase         = "/database/"		// /service/rpcName/database/
	KvServiceDatabaseEnable   = "/database/enable"
	KvServiceDatabaseHost     = "/database/host"
	KvServiceDatabaseUser     = "/database/user"
	KvServiceDatabasePassword = "/database/password"
	KvServiceDatabaseDbName   = "/database/dbname"
	KvServiceDatabasePort     = "/database/port"
	KvServiceDatabaseTimezone = "/database/timezone"
	KvServiceDatabasePrefix   = "/database/prefix"
)

var (
	host = em.MustGetServiceNameKvKey(KvServiceDatabaseHost)
	user = em.MustGetServiceNameKvKey(KvServiceDatabaseUser)
	password = em.MustGetServiceNameKvKey(KvServiceDatabasePassword)
	port = em.MustGetServiceNameKvKey(KvServiceDatabasePort)
	dbname = em.MustGetServiceNameKvKey(KvServiceDatabaseDbName)
	timezone = em.MustGetServiceNameKvKey(KvServiceDatabaseTimezone)
	prefix = em.MustGetServiceNameKvKey(KvServiceDatabasePrefix)

	migrate = []interface{}{
		&Attachment{},
	}
)

type database struct {

}

func NewDatabase() *database {
	return &database{}
}

func (this *database) Init()  {
	dbEnable, err := em.Kv.ReadKey(em_define.GetPathByFieldName(em_library.Config.Service.RpcName, KvServiceDatabaseEnable))
	if err != nil || strings.ToLower(dbEnable) != "true" {
		em_library.InitLog.Println("[WARNING]", em_define.GetPathByFieldName(em_library.Config.Service.RpcName, KvServiceDatabaseEnable), " is not configured or not enable!!")
	} else {
		// Init Database
		this.runDatabase()
		// Insert database initial data
		this.insertBasicDataToDatabase()
	}
}

func (this *database) insertBasicDataToDatabase()  {
	return
}