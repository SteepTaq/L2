package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wb_l2/18/internal/model"
)

func Response(w http.ResponseWriter, status int, body any) {
	json, err := json.Marshal(body)
	if err != nil {
		slog.Error(err.Error())

		send(
			w,
			http.StatusInternalServerError,
			model.ErrorResp("Something went wrong, try again later").ToJSON(),
		)
		return
	}

	send(w, status, json)
}

func InternalServerError(w http.ResponseWriter) {
	Response(w, http.StatusInternalServerError, model.ErrorResp("Something went wrong, try again later"))
}

func MethodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", "POST")
	Response(w, http.StatusMethodNotAllowed, model.ErrorResp("Method is not allowed"))
}

func send(w http.ResponseWriter, status int, json []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(json)
}
