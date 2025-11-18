package config

type GRPC struct {
	Port string `yaml:"port"`
}
// db
type Postgres struct {
	DSN string `yaml:"dsn"`
}
// auth
type Auth struct {
	JWTSecret string `yaml:"jwt_secret"`
	TTLMin    int    `yaml:"ttl_min"`
}

// config
type Config struct {
	GRPC     GRPC     `yaml:"grpc"`
	Postgres Postgres `yaml:"postgres"`
	Auth     Auth     `yaml:"auth"`
}
