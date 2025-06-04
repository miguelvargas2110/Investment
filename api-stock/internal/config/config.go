package config

import (
	"github.com/joho/godotenv" // Permite cargar variables de entorno desde un archivo .env
	"log"
	"os"
	"strconv"
	"time"
)

// Config estructura las configuraciones que usará toda la aplicación.
type Config struct {
	Environment      string        // Entorno de ejecución (development, production, etc.)
	DBURL            string        // URL de conexión a la base de datos CockroachDB
	APIToken         string        // Token de autenticación para la API externa
	APIBaseURL       string        // URL base de la API de acciones
	HTTPPort         string        // Puerto en el que corre el servidor HTTP
	HTTPReadTimeout  time.Duration // Tiempo máximo de espera para lectura de peticiones
	HTTPWriteTimeout time.Duration // Tiempo máximo de espera para escritura de respuestas
	WorkerInterval   time.Duration // Intervalo entre ejecuciones del worker
	MaxPages         int           // Límite de páginas a consultar en la API
	MaxRetries       int           // Número máximo de reintentos para peticiones fallidas
	InitialDelay     time.Duration // Retardo inicial antes de comenzar a consultar la API
}

// Load carga las variables de entorno desde un archivo .env (si existe) y las encapsula en una instancia Config.
func Load() *Config {
	// Carga el archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env") // Si no se encuentra, no se considera un error fatal
	}

	// Retorna una instancia de Config con valores leídos de variables de entorno o valores por defecto
	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		DBURL:            getEnv("COCKROACHDB_URL", ""),
		APIToken:         getEnv("API_TOKEN", ""),
		APIBaseURL:       getEnv("API_BASE_URL", ""),
		HTTPPort:         getEnv("PORT", "8080"),
		HTTPReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", 15*time.Second),
		HTTPWriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", 30*time.Second),
		WorkerInterval:   getEnvAsDuration("WORKER_INTERVAL", 1*time.Hour),
		MaxPages:         getEnvAsInt("MAX_PAGES", 20),
		MaxRetries:       getEnvAsInt("MAX_RETRIES", 3),
		InitialDelay:     getEnvAsDuration("INITIAL_DELAY", 1*time.Second),
	}
}

// getEnv obtiene una variable de entorno como string, o retorna un valor por defecto si no existe.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsDuration obtiene una variable de entorno como time.Duration, o retorna un valor por defecto si no es válida.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
	}
	return defaultValue
}

// getEnvAsInt obtiene una variable de entorno como int, o retorna un valor por defecto si no es válida.
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
