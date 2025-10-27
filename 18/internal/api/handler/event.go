package handler

import (
	"net/http"
	"wb_l2/18/internal/model"
	"wb_l2/18/internal/repository/inmemory/event"
	"wb_l2/18/internal/service"
	"wb_l2/18/pkg/http/request"
	"wb_l2/18/pkg/http/response"
)

type withId struct {
	Id int `json:"id"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.MethodNotAllowed(w, "POST")
		return
	}

	body, err := request.ReadBody(w, r)
	if err != nil {
		return
	}

	id, err := h.service.Event.Create(body)
	if err != nil {
		switch err {
		case model.InvalidFormat:
			response.Response(w, http.StatusBadRequest, model.ErrorResp(err.Error()))
		default:
			response.InternalServerError(w)
		}
		return
	}

	response.Response(
		w,
		http.StatusCreated,
		model.ResultWithDataResp("Event created successfully", withId{Id: id}),
	)
}

func (h *Handler) ListEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response.MethodNotAllowed(w, "GET")
		return
	}

	events, err := h.service.Event.List(r.URL.Query(), service.Day)
	if err != nil {
		switch err {
		case service.InvalidQuery:
			response.Response(w, http.StatusBadRequest, model.ErrorResp(err.Error()))
		default:
			response.InternalServerError(w)
		}
		return
	}

	response.Response(
		w,
		http.StatusOK,
		model.ResultWithDataResp("List of events for a day", events),
	)
}

func (h *Handler) ListEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response.MethodNotAllowed(w, "GET")
		return
	}

	events, err := h.service.Event.List(r.URL.Query(), service.Week)
	if err != nil {
		switch err {
		case service.InvalidQuery:
			response.Response(w, http.StatusBadRequest, model.ErrorResp(err.Error()))
		default:
			response.InternalServerError(w)
		}
		return
	}

	response.Response(
		w,
		http.StatusOK,
		model.ResultWithDataResp("List of events for a week", events),
	)
}

func (h *Handler) ListEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response.MethodNotAllowed(w, "GET")
		return
	}

	events, err := h.service.Event.List(r.URL.Query(), service.Month)
	if err != nil {
		switch err {
		case service.InvalidQuery:
			response.Response(w, http.StatusBadRequest, model.ErrorResp(err.Error()))
		default:
			response.InternalServerError(w)
		}
		return
	}

	response.Response(
		w,
		http.StatusOK,
		model.ResultWithDataResp("List of events for a month", events),
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.MethodNotAllowed(w, "POST")
		return
	}

	body, err := request.ReadBody(w, r)
	if err != nil {
		return
	}

	events, err := h.service.Event.Update(body)
	if err != nil {
		switch err {
		case model.InvalidFormat:
			response.Response(w, http.StatusBadRequest, model.ErrorResp(err.Error()))
		case event.ErrorEventNotFound:
			response.Response(w, http.StatusNotFound, model.ErrorResp(err.Error()))
		default:
			response.InternalServerError(w)
		}
		return
	}

	response.Response(
		w,
		http.StatusOK,
		model.ResultWithDataResp("Event updated successfully", events),
	)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.MethodNotAllowed(w, "POST")
		return
	}

	body, err := request.ReadBody(w, r)
	if err != nil {
		return
	}

	if err := h.service.Event.Delete(body); err != nil {
		switch err {
		case model.InvalidFormat:
			response.Response(w, http.StatusBadRequest, model.ErrorResp(err.Error()))
		case event.ErrorEventNotFound:
			response.Response(w, http.StatusNotFound, model.ErrorResp(err.Error()))
		default:
			response.InternalServerError(w)
		}
		return
	}

	response.Response(
		w,
		http.StatusOK,
		model.ResultResp("Event deleted successfully"),
	)
}
