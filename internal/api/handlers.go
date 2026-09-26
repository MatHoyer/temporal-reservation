package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	serviceerror "go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"

	"github.com/mathoyer/temporal-reservation/internal/reservation"
)

type Server struct {
	temporal  client.Client
	taskQueue string
}

func NewServer(temporal client.Client, taskQueue string) *Server {
	return &Server{temporal: temporal, taskQueue: taskQueue}
}

type createReservationResponse struct {
	ID string `json:"id"`
}

func (s *Server) createReservation(w http.ResponseWriter, r *http.Request) {
	var in reservation.Input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := "reservation-" + uuid.NewString()
	_, err := s.temporal.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: s.taskQueue,
	}, reservation.ReservationWorkflow, in)
	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, "failed to start reservation", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, createReservationResponse{ID: id})
}

func (s *Server) confirmPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := s.temporal.SignalWorkflow(r.Context(), id, "", reservation.ConfirmPaymentSignal, nil)
	if err != nil {
		var notFound *serviceerror.NotFound
		if errors.As(err, &notFound) {
			http.Error(w, "reservation not found", http.StatusNotFound)
			return
		}
		log.Printf("failed to signal workflow %q: %v", id, err)
		http.Error(w, "failed to submit payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) getReservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	resp, err := s.temporal.QueryWorkflow(r.Context(), id, "", reservation.StatusQuery)
	if err != nil {
		var notFound *serviceerror.NotFound
		if errors.As(err, &notFound) {
			http.Error(w, "reservation not found", http.StatusNotFound)
			return
		}
		log.Printf("failed to query workflow %q: %v", id, err)
		http.Error(w, "failed to fetch reservation status", http.StatusInternalServerError)
		return
	}

	var result reservation.QueryResult
	if err := resp.Get(&result); err != nil {
		log.Printf("failed to decode query result for %q: %v", id, err)
		http.Error(w, "failed to fetch reservation status", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
