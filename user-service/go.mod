module github.com/axmdl1/MedicalDataExchange/user-service

go 1.25

require (
	github.com/axmdl1/MedicalDataExchange/core-service v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/rs/zerolog v1.34.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.31.1
)

replace github.com/axmdl1/MedicalDataExchange/core-service => ../core-service
