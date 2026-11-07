package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nxkh4ng/debui/driver"
)

func TableData(d driver.Driver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := 100, 0
		if s := r.URL.Query().Get("limit"); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 500 {
				writeError(
					w, http.StatusBadRequest,
					fmt.Sprintf("limit must be in range 1..500, got %q", s),
				)
				return
			}
			limit = n
		}
		if s := r.URL.Query().Get("offset"); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil || n < 0 {
				writeError(
					w, http.StatusBadRequest,
					fmt.Sprintf("offset must be non-negative integer, got %q", s),
				)
				return
			}
			offset = n
		}

		table := r.PathValue("table")
		data, err := d.GetData(r.Context(), table, limit, offset)
		if err != nil {
			if errors.Is(err, driver.ErrTableNotFound) {
				writeError(
					w, http.StatusNotFound,
					fmt.Sprintf("table %q not found", table),
				)
				return
			}

			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, data)
	}
}
