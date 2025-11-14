package app

import (
	"net/http"
	"strconv"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/client"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/config"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/handler"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	http.Server
}

func New(cfg config.Config) *App {
	dtClient, _ := client.NewDataTransferClient(cfg.GRPC.DataTransfer)
	h := handler.NewDataTransferHandler(dtClient)
	hClinic := handler.NewClinicHandler(dtClient)

	uClient, _ := client.NewUserClient(cfg.GRPC.User)
	hUser := handler.NewUserHandler(uClient.API)

	r := chi.NewRouter()

	// Добавляем CORS middleware (должен быть первым!)
	r.Use(middleware.CORS)

	// Добавляем middleware для извлечения Authorization header
	r.Use(middleware.ExtractAuthToken)

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

	// Clinics
	r.Route("/clinics", func(r chi.Router) {
		r.Post("/", hClinic.CreateClinic)
		r.Get("/", hClinic.ListClinics)
		r.Get("/{id}", hClinic.GetClinic)
		r.Patch("/{id}", hClinic.UpdateClinic)
		r.Delete("/{id}", hClinic.DeleteClinic)
	})

	//Users
	r.Route("/users", func(r chi.Router) {
		r.Post("/", hUser.CreateUser)       // публичная рега пациента или по token-у (employee/admin)
		r.Get("/", hUser.ListUsers)         // требует Authorization
		r.Get("/{id}", hUser.GetUser)       // требует Authorization
		r.Patch("/{id}", hUser.UpdateUser)  // требует Authorization
		r.Delete("/{id}", hUser.DeleteUser) // требует Authorization
	})
	r.Post("/users/login", hUser.Login) // публичный

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
