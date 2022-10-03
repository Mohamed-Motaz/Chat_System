package Server

import (
	logger "Server/Logger"
	"encoding/json"
	"net/http"
)

func generateSuccess(msg interface{}) []byte {

	res, err := json.Marshal(msg)
	if err != nil {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to marshal http success resp with err %v", err)
		return nil
	}
	return res
}

func success(w http.ResponseWriter, r *http.Request, msg interface{}, essential int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "Response ready to req %v", r.RequestURI)

	_, err := w.Write(generateSuccess(msg))
	if err != nil {
		logger.LogError(logger.SERVER, logger.ESSENTIAL, "Unable to write response to http request with err %v", err)
	}
}
