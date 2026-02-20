package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env           string      `yaml:"env" env-default:"local"`
	HTTPServer    HTTPServer  `yaml:"http_server"`
	Postgres      PostgresCfg `yaml:"postgres"`
	MigrationsDir string      `yaml:"migrations_dir" env-default:"./database/migrations"`
	DBDialect     string      `yaml:"db_dialect" env-default:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresCfg struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_INTERNAL_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DB" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
	TZ       string `yaml:"tz" env:"TZ" env-default:"Europe/Moscow"`

	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"50"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"10"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"30m"`
}

func (p PostgresCfg) DSN() string {
	return "host=" + p.Host +
		" user=" + p.User +
		" password=" + p.Password +
		" dbname=" + p.DBName +
		" port=" + p.Port +
		" sslmode=" + p.SSLMode +
		" TimeZone=" + p.TZ
}

func MustLoad() *Config {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %s", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check is file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
