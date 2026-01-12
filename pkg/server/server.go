package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alanshaw/1up-service/pkg/build"
	"github.com/alanshaw/ucantone/ucan"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("server")

type ServerInfo struct {
	ID    string    `json:"id"`
	Build BuildInfo `json:"build"`
}

type BuildInfo struct {
	Version string `json:"version"`
	Repo    string `json:"repo"`
}

func NewRootHandler(id ucan.Principal) http.Handler {
	info := ServerInfo{
		ID: id.DID().String(),
		Build: BuildInfo{
			Version: build.Version,
			Repo:    "https://github.com/alanshaw/1up-service",
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			data, err := json.Marshal(&info)
			if err != nil {
				log.Errorf("failed JSON marshal server info: %w", err)
				http.Error(w, "failed JSON marshal server info", http.StatusInternalServerError)
				return
			}
			w.Write(data)
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("▗     \n"))
			w.Write([]byte("▜ ▌▌▛▌\n"))
			w.Write([]byte("▟▖▙▌▙▌\n"))
			w.Write([]byte("    ▌ \n"))
			w.Write([]byte("\n"))
			w.Write([]byte(fmt.Sprintf("🍄 1up-service %s\n", info.Build.Version)))
			w.Write([]byte(fmt.Sprintf("🆔 %s\n", info.ID)))
			w.Write([]byte("🐙 https://github.com/alanshaw/1up-service\n"))
		}
	})
}
