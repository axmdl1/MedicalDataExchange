package service

type Claims struct {
	UserID   int64
	Role     string // "patient"|"employee"|"admin"
	ClinicID *int64 // nil для patient/admin, обязателен для employee
	Email    string
}
