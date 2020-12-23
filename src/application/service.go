package application


const (
	Version_Service = "1.2.4"
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