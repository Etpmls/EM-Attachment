package register_config

type ServiceConfiguration struct {
	Service struct{
		Host string
		FileStorageMethod string 	`yaml:"file-storage-method"`
	}
}

var ServiceConfig ServiceConfiguration