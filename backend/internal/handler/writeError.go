package handler

import (
	appErr "github.com/valeriapadilla/stock-insights/internal/errors"
	"net/http"
	"encoding/json"

)

func WriteHTTPError(w http.ResponseWriter, err error){
	httpErr, ok := err.(*appErr.HTTPError)
	if !ok{
		httpErr = appErr.ErrInternalServer
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(W).Endode(map[string]string{
		"error": httpErr.Message,
	})
}