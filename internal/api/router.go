package api

import (
	"io/fs"
	"net/http"
)

func (s *Server) Routes(webFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/reservations", s.createReservation)
	mux.HandleFunc("GET /api/reservations/{id}", s.getReservation)
	mux.HandleFunc("POST /api/reservations/{id}/pay", s.confirmPayment)
	mux.Handle("/", http.FileServer(http.FS(webFS)))

	return mux
}
