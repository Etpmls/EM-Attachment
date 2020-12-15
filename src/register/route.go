package register

import (
	"github.com/Etpmls/EM-Attachment/src/application/service"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"os"
	"strings"
)

// Register Route
func RegisterRoute(mux *runtime.ServeMux)  {
	mux.HandlePath("GET", em_library.Config.ServiceDiscovery.Service.CheckUrl, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("hello"))
	})
	mux.HandlePath("POST", "/api/attachment/v1/attachment/uploadImage", service.ServiceAttachment{}.UploadImage)

	mux.HandlePath("GET", "/storage/upload/**", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		f, err := os.Stat(strings.TrimLeft(r.URL.String(), "/"))
		if err != nil || f.IsDir() {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fileServer := http.StripPrefix("/storage/upload", http.FileServer(http.Dir("./storage/upload")))
		fileServer.ServeHTTP(w, r)
	})
}
