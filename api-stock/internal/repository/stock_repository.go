package repository

import (
	"api-stock/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// stockRepository implementa la interfaz domain.StockRepository y maneja las operaciones con la base de datos.
type stockRepository struct {
	db *sql.DB // Conexión a la base de datos SQL
}

// NewStockRepository crea una nueva instancia de stockRepository con la conexión a la base de datos proporcionada.
func NewStockRepository(db *sql.DB) domain.StockRepository {
	return &stockRepository{db: db}
}

// GetRecommendations obtiene recomendaciones paginadas desde la base de datos, con filtro opcional por ticker.
// Recibe el contexto para control de tiempo y cancelación, el ticker para filtrar, y los parámetros de paginación (página y límite).
// Devuelve la lista de recomendaciones, el total de recomendaciones para la consulta, y un error si ocurre alguno.
func (r *stockRepository) GetRecommendations(ctx context.Context, ticker string, page, limit int) ([]domain.StockRecommendation, int, error) {
	offset := (page - 1) * limit // Calcula el offset para paginación

	// Consulta SQL que selecciona las recomendaciones, filtrando por ticker si se pasa (si ticker es cadena vacía, no filtra)
	query := `SELECT ticker, target_from, target_to, company, action, 
              brokerage, rating_from, rating_to, time
              FROM recommendations
              WHERE ($1 = '' OR ticker = $1)
              ORDER BY time DESC
              LIMIT $2 OFFSET $3`

	// Ejecuta la consulta con los parámetros recibidos
	rows, err := r.db.QueryContext(ctx, query, ticker, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error in SQL query: %v", err)
	}
	defer rows.Close()

	var recommendations []domain.StockRecommendation

	// Itera las filas obtenidas y las escanea en structs StockRecommendation
	for rows.Next() {
		var rec domain.StockRecommendation
		err := rows.Scan(
			&rec.Ticker,
			&rec.TargetFrom,
			&rec.TargetTo,
			&rec.Company,
			&rec.Action,
			&rec.Brokerage,
			&rec.RatingFrom,
			&rec.RatingTo,
			&rec.Time,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %v", err)
		}
		recommendations = append(recommendations, rec)
	}

	// Consulta para contar el total de recomendaciones que cumplen el filtro (para paginación)
	var total int
	countQuery := `SELECT COUNT(*) FROM recommendations WHERE ($1 = '' OR ticker = $1)`
	err = r.db.QueryRowContext(ctx, countQuery, ticker).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting rows: %v", err)
	}

	// Devuelve la lista de recomendaciones, el total y nil si no hay error
	return recommendations, total, nil
}

// GetAvailableTickers devuelve una lista con los tickers distintos existentes en la tabla de recomendaciones.
// Se usa para conocer qué símbolos están disponibles para filtrar o mostrar.
func (r *stockRepository) GetAvailableTickers(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT ticker FROM recommendations ORDER BY ticker`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error en consulta SQL: %v", err)
	}
	defer rows.Close()

	var tickers []string
	// Itera sobre cada fila y escanea el ticker
	for rows.Next() {
		var ticker string
		if err := rows.Scan(&ticker); err != nil {
			return nil, fmt.Errorf("error escaneando ticker: %v", err)
		}
		tickers = append(tickers, ticker)
	}

	// Verifica errores que puedan haber ocurrido durante la iteración
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error después de iterar filas: %v", err)
	}

	return tickers, nil
}

// InsertRecommendations inserta un lote (bulk insert) de recomendaciones en la base de datos.
// Usa transacciones para asegurar que todas las inserciones ocurran juntas.
// En caso de conflicto (ticker y time duplicados), actualiza los datos existentes.
func (r *stockRepository) InsertRecommendations(ctx context.Context, recommendations []domain.StockRecommendation) error {
	if len(recommendations) == 0 {
		return nil // No hay nada que insertar
	}

	// Inicia una transacción
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback() // Rollback automático en caso de error

	// Prepara la consulta dinámica con los placeholders y los valores a insertar
	valueStrings := make([]string, 0, len(recommendations))
	valueArgs := make([]interface{}, 0, len(recommendations)*9) // 9 columnas por fila

	for i, rec := range recommendations {
		// Crea una parte de la query con placeholders ($1, $2, ... $9)
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9))

		// Agrega los valores en orden para cada fila
		valueArgs = append(valueArgs, rec.Ticker, rec.TargetFrom, rec.TargetTo,
			rec.Company, rec.Action, rec.Brokerage, rec.RatingFrom, rec.RatingTo, rec.Time)
	}

	// Construye la consulta SQL con ON CONFLICT para actualizar filas en caso de duplicados (ticker, time)
	stmt := fmt.Sprintf(`
        INSERT INTO recommendations (
            ticker, target_from, target_to, company, action, 
            brokerage, rating_from, rating_to, time
        ) VALUES %s
        ON CONFLICT (ticker, time) DO UPDATE SET
            target_from = EXCLUDED.target_from,
            target_to = EXCLUDED.target_to,
            company = EXCLUDED.company,
            action = EXCLUDED.action,
            brokerage = EXCLUDED.brokerage,
            rating_from = EXCLUDED.rating_from,
            rating_to = EXCLUDED.rating_to`,
		strings.Join(valueStrings, ","))

	// Ejecuta la consulta con todos los valores
	_, err = tx.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		return fmt.Errorf("error en bulk insert: %v", err)
	}

	// Hace commit si todo fue exitoso
	return tx.Commit()
}

// GetRecentRecommendations obtiene recomendaciones con fecha mayor a un intervalo de tiempo dado (desde ahora menos el intervalo).
// Útil para obtener recomendaciones recientes.
func (r *stockRepository) GetRecentRecommendations(ctx context.Context, since time.Duration) ([]domain.StockRecommendation, error) {
	query := `SELECT ticker, target_from, target_to, company, action, 
              brokerage, rating_from, rating_to, time
              FROM recommendations
              WHERE time > $1
              ORDER BY time DESC`

	// Calcula el timestamp límite con time.Now() - since
	rows, err := r.db.QueryContext(ctx, query, time.Now().Add(-since))
	if err != nil {
		return nil, fmt.Errorf("error en consulta SQL: %v", err)
	}
	defer rows.Close()

	var recommendations []domain.StockRecommendation

	// Escanea cada fila en una recomendación
	for rows.Next() {
		var rec domain.StockRecommendation
		err := rows.Scan(
			&rec.Ticker,
			&rec.TargetFrom,
			&rec.TargetTo,
			&rec.Company,
			&rec.Action,
			&rec.Brokerage,
			&rec.RatingFrom,
			&rec.RatingTo,
			&rec.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando fila: %v", err)
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

// GetLatestRecommendation obtiene la recomendación más reciente ordenada por fecha descendente.
// Retorna nil si no hay recomendaciones.
func (r *stockRepository) GetLatestRecommendation(ctx context.Context) (*domain.StockRecommendation, error) {
	query := `SELECT 
        ticker, target_from, target_to, company, action, 
        brokerage, rating_from, rating_to, time
        FROM recommendations
        ORDER BY time DESC
        LIMIT 1`

	row := r.db.QueryRowContext(ctx, query)

	var rec domain.StockRecommendation

	// Escanea la fila resultante
	err := row.Scan(
		&rec.Ticker,
		&rec.TargetFrom,
		&rec.TargetTo,
		&rec.Company,
		&rec.Action,
		&rec.Brokerage,
		&rec.RatingFrom,
		&rec.RatingTo,
		&rec.Time,
	)

	// Manejo de errores
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No hay datos
		}
		return nil, fmt.Errorf("error escaneando recomendación: %v", err)
	}

	return &rec, nil
}

// DeleteAllRecommendations elimina todas las recomendaciones de la tabla.
func (r *stockRepository) DeleteAllRecommendations(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM recommendations")
	if err != nil {
		return fmt.Errorf("error eliminando recomendaciones: %v", err)
	}
	return nil
}

// GetStockFeatures calcula características agregadas de las recomendaciones para un ticker dado.
// Devuelve un mapa con métricas numéricas para análisis o modelo de recomendaciones.
func (r *stockRepository) GetStockFeatures(ctx context.Context, ticker string) (map[string]float64, error) {
	query := `
    SELECT 
        COUNT(*) as total, -- Total de recomendaciones para el ticker
        AVG(CAST(REPLACE(target_to, '$', '') AS FLOAT8) - CAST(REPLACE(target_from, '$', '') AS FLOAT8)) as target_range, -- Promedio del rango objetivo (diferencia target_to - target_from)
        AVG(CASE WHEN action LIKE '%aumentado%' THEN 1 ELSE 0 END) as upgrade_prob, -- Probabilidad promedio de recomendaciones con acción "aumentado"
        AVG(CASE WHEN action LIKE '%bajado%' THEN 1 ELSE 0 END) as downgrade_prob, -- Probabilidad promedio de recomendaciones con acción "bajado"
        AVG(CASE WHEN rating_to LIKE '%comprar%' THEN 1 ELSE 0 END) as buy_rating, -- Promedio de recomendaciones con rating para "comprar"
        AVG(CASE WHEN rating_to LIKE '%vender%' THEN 1 ELSE 0 END) as sell_rating, -- Promedio de recomendaciones con rating para "vender"
        STDDEV(CAST(REPLACE(target_to, '$', '') AS FLOAT8) - CAST(REPLACE(target_from, '$', '') AS FLOAT8)) as target_volatility, -- Volatilidad (desviación estándar) del rango objetivo
        COUNT(DISTINCT brokerage) as unique_brokers -- Número de brokers distintos que han hecho recomendaciones
    FROM recommendations
    WHERE ticker = $1`

	row := r.db.QueryRowContext(ctx, query, ticker)

	features := make(map[string]float64)
	var total int
	var targetRange, upgradeProb, downgradeProb, buyRating, sellRating, targetVolatility *float64
	var uniqueBrokers int

	// Escanea los resultados en variables
	err := row.Scan(&total, &targetRange, &upgradeProb, &downgradeProb, &buyRating, &sellRating, &targetVolatility, &uniqueBrokers)
	if err != nil {
		return nil, err
	}

	// Agrega al mapa de resultados usando safeFloat para evitar nil
	features["total_recommendations"] = float64(total)
	features["target_range"] = safeFloat(targetRange)
	features["upgrade_probability"] = safeFloat(upgradeProb)
	features["downgrade_probability"] = safeFloat(downgradeProb)
	features["buy_rating"] = safeFloat(buyRating)
	features["sell_rating"] = safeFloat(sellRating)
	features["target_volatility"] = safeFloat(targetVolatility)

	features["unique_brokers"] = float64(uniqueBrokers)
	features["broker_diversity"] = float64(uniqueBrokers) / float64(total) // Diversidad de brokers

	return features, nil
}

// GetAllStockFeatures obtiene las características agregadas para todos los tickers disponibles.
// Devuelve una lista con ticker y sus características para análisis o procesamiento.
func (r *stockRepository) GetAllStockFeatures(ctx context.Context) ([]struct {
	Ticker   string
	Features map[string]float64
}, error) {
	tickers, err := r.GetAvailableTickers(ctx) // Obtiene todos los tickers
	if err != nil {
		return nil, err
	}

	var result []struct {
		Ticker   string
		Features map[string]float64
	}

	// Para cada ticker obtiene sus características y las agrega al resultado
	for _, ticker := range tickers {
		features, err := r.GetStockFeatures(ctx, ticker)
		if err != nil {
			continue // Si hay error con un ticker, continúa con los demás
		}
		result = append(result, struct {
			Ticker   string
			Features map[string]float64
		}{ticker, features})
	}

	return result, nil
}

// Ping verifica la conexión a la base de datos ejecutando un ping.
func (r *stockRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// RunMigrations crea la tabla recommendations y los índices necesarios si no existen.
// Esto asegura que la base de datos tenga la estructura mínima para almacenar datos.
func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS recommendations (
			ticker VARCHAR(10),
			target_from VARCHAR(20),
			target_to VARCHAR(20),
			company VARCHAR(100),
			action VARCHAR(50),
			brokerage VARCHAR(100),
			rating_from VARCHAR(50),
			rating_to VARCHAR(50),
			time TIMESTAMP,
			PRIMARY KEY (ticker, time)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_recommendations_ticker ON recommendations (ticker)`,
		`CREATE INDEX IF NOT EXISTS idx_recommendations_time ON recommendations (time)`,
	}

	// Ejecuta cada query de migración
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing migration: %v", err)
		}
	}
	return nil
}

// safeFloat es una función auxiliar que devuelve 0 si el puntero es nil, o el valor apuntado si no.
func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
