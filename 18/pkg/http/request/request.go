package request

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"wb_l2/18/internal/model"
	"wb_l2/18/pkg/http/response"
)

var invalidRequest = fmt.Errorf("Invalid Request")

func ReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.Header().Add("Accept-Post", "application/json")
		response.Response(w, http.StatusUnsupportedMediaType, model.ErrorResp("Content-Type must be application/json"))
		return []byte{}, invalidRequest
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Response(w, http.StatusBadRequest, model.ErrorResp("Invalid body"))
		return []byte{}, invalidRequest
	}

	return body, nil
}
