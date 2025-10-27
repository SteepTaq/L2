package handler

import (
	"net/http"
	"wb_l2/18/internal/model"
	"wb_l2/18/internal/service"
	"wb_l2/18/pkg/http/response"
)

type Handler struct {
	mux     *http.ServeMux
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		mux:     http.NewServeMux(),
		service: service,
	}
}

func (h *Handler) HTTPHandler() http.Handler {
	return h.mux
}

func RegisterHandlers(h *Handler) {
	h.mux.HandleFunc("/ping", h.Ping)

	h.mux.HandleFunc("/create_event", h.CreateEvent)

	h.mux.HandleFunc("/events_for_day", h.ListEventsForDay)
	h.mux.HandleFunc("/events_for_week", h.ListEventsForWeek)
	h.mux.HandleFunc("/events_for_month", h.ListEventsForMonth)

	h.mux.HandleFunc("/update_event", h.Update)

	h.mux.HandleFunc("/delete_event", h.Delete)

	h.mux.HandleFunc("/", h.NotFound)
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	response.Response(w, http.StatusOK, model.ResultResp("Pong!"))
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	response.Response(w, http.StatusNotFound, model.ErrorResp("Unknown endpoint"))
}
