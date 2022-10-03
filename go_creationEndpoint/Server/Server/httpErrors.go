package Server

import (
	logger "Server/Logger"
	"encoding/json"
	"net/http"
)

type ErrorResp struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func generateError(msg string) []byte {
	errorResp := &ErrorResp{
		Error:        true,
		ErrorMessage: msg,
	}
	res, err := json.Marshal(errorResp)
	if err != nil {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to marshal http error with err %v", err)
		return nil
	}
	return res
}

func failure(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	logger.LogError(logger.SERVER, logger.ESSENTIAL, "Failure with code %v to req %v", status, r.RequestURI)

	_, err := w.Write(generateError(msg))
	if err != nil {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to write error to http request with err %v", err)
	}
}
