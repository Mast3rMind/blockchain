package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ipkg/blockchain"
)

// BlockStore http interface for inspecting the store

type server struct {
	st blockchain.BlockStore
}

func (svr *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := strings.Trim(r.URL.Path, "/")

	var resp interface{}
	switch upath {
	case "last":
		resp = svr.st.LastBlock()
	case "first":
		resp = svr.st.FirstBlock()
	case "size":
		resp = map[string]interface{}{"size": svr.st.BlockCount()}

	case "chain":
		s := svr.st

		b := s.LastBlock()
		if b == nil || b.BlockHeader == nil || isZeroBytes(b.PrevHash) {
			w.WriteHeader(404)
			w.Write([]byte("nil\n"))
			return
		}

		for {
			w.Write([]byte(fmt.Sprintf("%x\n", b.Hash())))
			for _, tx := range b.Transactions {
				w.Write([]byte(fmt.Sprintf("  %x\n", tx.Hash())))
			}

			b = s.Get(b.PrevHash)
			if b == nil || b.BlockHeader == nil || isZeroBytes(b.PrevHash) {
				break
			}

		}

		return

	default:
		w.WriteHeader(404)
		return
	}

	b, _ := json.Marshal(resp)
	w.Write(b)
}

func isZeroBytes(b []byte) bool {
	for _, e := range b {
		if e != 0 {
			return false
		}
	}
	return true
}
