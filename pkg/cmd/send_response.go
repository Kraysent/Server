package cmd

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func SendResponse(w http.ResponseWriter, code int, body any, err error, msg string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(code)

	var logEvent *zerolog.Event
	if code == http.StatusOK {
		logEvent = zlog.Debug()
	} else if (code == http.StatusUnauthorized) || (code == http.StatusBadRequest) {
		logEvent = zlog.Warn()
	} else {
		logEvent = zlog.Error()
	}

	logEvent = logEvent.Int("response_code", code).Err(err)

	if body != nil {
		resp, err1 := json.Marshal(body)
		if err1 != nil {
			logEvent.Err(err1).Msg("Error during response sending.")
			return
		}
		w.Write(resp)

		logEvent = logEvent.RawJSON("response", resp)
	}

	logEvent.Msg(msg)
}
