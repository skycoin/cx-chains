package api

import (
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"

	"github.com/skycoin/cx-chains/src/readable"
)

// versionHandler returns the application version info
// URI: /api/v1/version
// Method: GET
func versionHandler(bi readable.BuildInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, bi)
	}
}
