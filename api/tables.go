package api

import (
	"net/http"

	"github.com/nxkh4ng/debui/driver"
)

func Tables(d driver.Driver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tables, err := d.ListTables(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, tables)
	}
}
