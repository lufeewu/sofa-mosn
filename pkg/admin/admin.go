package admin

import (
	"encoding/json"
	"github.com/alipay/sofa-mosn/pkg/log"
	"net/http"
)

func configDump(w http.ResponseWriter, _ *http.Request) {
	if buf, err := json.Marshal(GetEffectiveConfig()); err == nil {
		w.Write(buf);
	} else {
		w.WriteHeader(500)
		w.Write([]byte(`{ error: "internal error" }`))
		log.DefaultLogger.Errorf("Admin API: ConfigDump failed, cause by %s", err)
	}
}

func Start(config interface{}) *http.Server {
	// merge MOSNConfig into global context
	MergeOriginalConf(config)
	srv := &http.Server{Addr: ":8888"}

	go func() {
		http.HandleFunc("/api/v1/config_dump", configDump)
		if err := srv.ListenAndServe(); err != nil {
			log.DefaultLogger.Errorf("Admin Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	return srv
}
