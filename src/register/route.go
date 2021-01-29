package register

import (
	"github.com/Etpmls/EM-Attachment/src/application/service"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	em_library "github.com/Etpmls/Etpmls-Micro/v2/library"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"os"
	"strings"
)

// Register Route
func RegisterRoute(mux *runtime.ServeMux)  {
	e, _ := em.Kv.ReadKey(define.MakeServiceConfField(em_library.Config.Service.RpcId, define.KvServiceCheckUrl))

	mux.HandlePath("GET", e, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
