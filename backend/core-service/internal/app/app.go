package app

import (
	"net/http"
	"strconv"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/blockchain"
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

	// --- gRPC client (настоящий)
	dtClient, _ := client.NewDataTransferClient(cfg.GRPC.DataTransfer)

	// --- Blockchain
	bc := blockchain.NewBlockchain()

	// --- Patient Access Adapter (поверх dtClient)
	dtAdapter := client.NewDataTransferAdapter(dtClient)

	// --- Handlers
	h := handler.NewDataTransferHandler(dtClient) // ← ТОЛЬКО dtClient
	hClinic := handler.NewClinicHandler(dtClient) // ← ТОЛЬКО dtClient

	hPatientAccess := handler.NewPatientAccessHandler(dtAdapter, bc) // ← АДАПТЕР + BLOCKCHAIN

	uClient, _ := client.NewUserClient(cfg.GRPC.User)
	hUser := handler.NewUserHandler(uClient.API)

	r := chi.NewRouter()

	r.Use(middleware.CORS)
	r.Use(middleware.ExtractAuthToken)
	r.Use(middleware.ExtractAuthenticatedUser)

	// --- MedicalData
	r.Route("/medical-data", func(r chi.Router) {
		r.Post("/", h.CreateMedicalData)
		r.Get("/", h.ListMedicalData)
		r.Get("/{id}", h.GetMedicalData)
		r.Delete("/{id}", h.DeleteMedicalData)
	})

	// --- Transfers
	r.Route("/transfers", func(r chi.Router) {
		r.Post("/", h.CreateDataTransfer)
		r.Get("/", h.ListDataTransfers)
		r.Get("/{id}", h.GetDataTransfer)
		r.Post("/decision", h.HandleDecision)
	})

	// --- Clinics
	r.Route("/clinics", func(r chi.Router) {
		r.Post("/", hClinic.CreateClinic)
		r.Get("/", hClinic.ListClinics)
		r.Get("/{id}", hClinic.GetClinic)
		r.Patch("/{id}", hClinic.UpdateClinic)
		r.Delete("/{id}", hClinic.DeleteClinic)
	})

	// --- Users
	r.Route("/users", func(r chi.Router) {
		r.Post("/", hUser.CreateUser)
		r.Get("/", hUser.ListUsers)
		r.Get("/{id}", hUser.GetUser)
		r.Patch("/{id}", hUser.UpdateUser)
		r.Delete("/{id}", hUser.DeleteUser)
	})
	r.Post("/users/login", hUser.Login)

	// --- Patient Access (adapter + blockchain)
	r.Route("/patient-access", func(r chi.Router) {
		r.Post("/request", hPatientAccess.CreateAccessRequest)
		r.Post("/approve/{id}", hPatientAccess.ApproveAccessRequest)
		r.Post("/reject/{id}", hPatientAccess.RejectAccessRequest)
		r.Get("/data", hPatientAccess.GetTemporaryData)
		r.Get("/requests", hPatientAccess.ListAccessRequests)
		//r.Get("/request/{id}", hPatientAccess.GetAccessRequest)
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
