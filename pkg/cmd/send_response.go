package cmd

import (
	"encoding/json"
	"net/http"

	zlog "github.com/rs/zerolog/log"
)

func SendResponse(w http.ResponseWriter, code int, body any, err error, msg string) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(code)

	empty := true
	logEvent := zlog.Debug()

	if err != nil {
		logEvent = logEvent.Err(err)
		empty = false
	}

	if body != nil {
		resp, err1 := json.Marshal(body)
		if err1 != nil {
			logEvent.Err(err1).Msg("Error during response sending.")
			return
		}
		w.Write(resp)

		logEvent = logEvent.RawJSON("response", resp)
		empty = false
	}

	if msg != "" {
		logEvent.Msg(msg)
	} else {
		if !empty {
			logEvent.Send()
		}
	}
}
