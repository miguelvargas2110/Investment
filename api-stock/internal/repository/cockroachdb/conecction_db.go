package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

// Connect establece una conexión con la base de datos CockroachDB usando pgx.
// Configura parámetros de conexión y valida que la conexión sea exitosa mediante un ping.
func Connect(connStr string) (*sql.DB, error) {
	// Parseo de la cadena de conexión
	connConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la cadena de conexión: %v", err)
	}

	// Apertura de la base de datos como *sql.DB usando stdlib
	db := stdlib.OpenDB(*connConfig)

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(50)                  // Número máximo de conexiones abiertas
	db.SetMaxIdleConns(30)                  // Número máximo de conexiones inactivas
	db.SetConnMaxLifetime(10 * time.Minute) // Tiempo máximo de vida de una conexión

	// Validación de conexión con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error al verificar la conexión: %v", err)
	}

	log.Println("Conexión exitosa a CockroachDB")
	return db, nil
}
