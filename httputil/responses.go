package httputil

import (
	"encoding/json"
	"io"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, data interface{}, status ...int) {
	w.Header().Set(`Content-Type`, `application/json`)
	headerSent := false

	if err, ok := data.(error); ok && err != nil {
		data = map[string]interface{}{
			`success`: false,
			`error`:   err.Error(),
		}

		if len(status) == 0 {
			status = []int{http.StatusInternalServerError}
		}
	}

	if len(status) > 0 {
		w.WriteHeader(status[0])
		headerSent = true
	}

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			Logger.Warningf("Failed to encode response body: %v", err)
		}
	} else if !headerSent {
		w.WriteHeader(http.StatusNoContent)
	}
}

func ParseJSON(r io.Reader, into interface{}) error {
	return json.NewDecoder(r).Decode(into)
}
