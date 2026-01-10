package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alanshaw/1up-service/pkg/build"
	"github.com/alanshaw/1up-service/pkg/store/token"
	"github.com/alanshaw/ucantone/ipld/codec/dagcbor"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/container"
	"github.com/ipfs/go-cid"
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

func NewReceiptHandler(tokens token.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskLink, err := cid.Parse(r.PathValue("task"))
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid task CID: %v", err), http.StatusBadRequest)
			return
		}

		var invocations []ucan.Invocation
		var delegations []ucan.Delegation
		var receipts []ucan.Receipt
		for c, err := range tokens.FindByTask(r.Context(), taskLink) {
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to find receipt token: %v", err), http.StatusInternalServerError)
				return
			}
			invocations = append(invocations, c.Invocations()...)
			delegations = append(delegations, c.Delegations()...)
			receipts = append(receipts, c.Receipts()...)
		}

		out := container.New(
			container.WithInvocations(invocations...),
			container.WithDelegations(delegations...),
			container.WithReceipts(receipts...),
		)

		w.Header().Set("Content-Type", dagcbor.ContentType)
		err = out.MarshalCBOR(w)
		if err != nil {
			log.Errorw("marshaling receipt container", "error", err)
		}
	})
}
