package application


const (
	Version_Service = "1.2.5"
)

const (
	Service_Attachment = "AttachmentRpcService"
)


/*
	Config
*/
type serviceConfiguration struct {
	Service struct {
		Host              string
		FileStorageMethod string `yaml:"file-storage-method"`
	}
}

var ServiceConfig serviceConfiguration