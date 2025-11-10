package app

import (
	"net/http"
	"strconv"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/client"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/config"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/handler"
	"github.com/go-chi/chi/v5"
)

type App struct {
	http.Server
}

func New(cfg config.Config) *App {
	dtClient, _ := client.NewDataTransferClient(cfg.GRPC.DataTransfer)
	h := handler.NewDataTransferHandler(dtClient)

	r := chi.NewRouter()

	// MedicalData
	r.Route("/medical-data", func(r chi.Router) {
		r.Post("/", h.CreateMedicalData)
		r.Get("/", h.ListMedicalData)
		r.Get("/{id}", h.GetMedicalData)
		r.Delete("/{id}", h.DeleteMedicalData)
	})

	// DataTransfer
	r.Route("/transfers", func(r chi.Router) {
		r.Post("/", h.CreateDataTransfer)
		r.Get("/", h.ListDataTransfers)
		r.Get("/{id}", h.GetDataTransfer)
		r.Post("/decision", h.HandleDecision)
	})

	return &App{
		Server: http.Server{
			Addr:    ":" + strconv.Itoa(cfg.REST.Port),
			Handler: r,
		},
	}
}

func (a *App) Run() error {
	println("REST API running on " + a.Addr)
	return a.ListenAndServe()
}
