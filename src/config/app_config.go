package config

type Postgres struct {
	Uri string `yaml:"uri" env:"POSTGRES_URI"`
}

type AppConfig struct {
	Port     int      `yaml:"port" env:"PORT"`
	Postgres Postgres `yaml:"postgres"`
}
