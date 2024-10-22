package storage

import (
	"database/sql"
	"fmt"
	"github.com/LissaGreense/discogs_record_label/backend/models"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Table names
const (
	ArtistsTableName  = "artists"
	GenresTableName   = "genres"
	StylesTableName   = "styles"
	releasesTableName = "releases"
)

// SQL statements for creating tables
const (
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS %s (
		%s
	);`
	releasesColumnDef   = `id INT PRIMARY KEY`
	attributesColumnDef = `id SERIAL PRIMARY KEY,
		release_id INT REFERENCES %s(id) ON DELETE CASCADE,
		name TEXT NOT NULL`
)

// SQL queries for insertion and fetching
const (
	insertReleaseSQL = `
		INSERT INTO %s (id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING;
	`

	insertAttributeSQL = `
		INSERT INTO %s (release_id, name)
		VALUES ($1, $2);
	`
	fetchAttrsNamesSQL = `
		SELECT
			COUNT(DISTINCT r.id) as releaseCount,
			COALESCE(a.name, '') as artistName,
			COALESCE(s.name, '') as styleName,
			COALESCE(g.name, '') as genreName
		FROM %s r
		LEFT JOIN %s a ON r.id = a.release_id
		LEFT JOIN %s s ON r.id = s.release_id
		LEFT JOIN %s g ON r.id = g.release_id
		WHERE 1=1
	`

	fetchUniqueNamesSQL = `
		SELECT DISTINCT name FROM %s ORDER BY name
	`
)

func InitDatabase() (*sql.DB, error) {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")

	if pgUser == "" || pgPassword == "" || pgDB == "" || pgHost == "" || pgPort == "" {
		return nil, fmt.Errorf("PostgreSQL environment variables are not set")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateSchema(db *sql.DB) error {
	creationFailedMsg := "failed to create %s table: %v"

	if err := createTable(db, releasesColumnDef, releasesTableName); err != nil {
		return fmt.Errorf(creationFailedMsg, releasesTableName, err)
	}

	if err := createTable(db, attributesColumnDef, ArtistsTableName, releasesTableName); err != nil {
		return fmt.Errorf(creationFailedMsg, ArtistsTableName, err)
	}

	if err := createTable(db, attributesColumnDef, GenresTableName, releasesTableName); err != nil {
		return fmt.Errorf(creationFailedMsg, GenresTableName, err)
	}

	if err := createTable(db, attributesColumnDef, StylesTableName, releasesTableName); err != nil {
		return fmt.Errorf(creationFailedMsg, StylesTableName, err)
	}

	log.Println("Tables created successfully")
	return nil
}

func createTable(db *sql.DB, columnDef string, tableName string, additionalArgs ...any) error {
	filledColumDefs := fmt.Sprintf(columnDef, additionalArgs...)
	finalCreateTableStatement := fmt.Sprintf(createTableSQL, tableName, filledColumDefs)

	_, err := db.Exec(finalCreateTableStatement)
	return err
}

func StoreRelease(db *sql.DB, release *models.Release) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	if err := insertRelease(tx, release); err != nil {
		return err
	}

	if err := insertAttributes(tx, release.Id, release.Artists, ArtistsTableName); err != nil {
		return fmt.Errorf("failed to insert artists: %v", err)
	}
	if err := insertAttributes(tx, release.Id, release.Genres, GenresTableName); err != nil {
		return fmt.Errorf("failed to insert genres: %v", err)
	}
	if err := insertAttributes(tx, release.Id, release.Styles, StylesTableName); err != nil {
		return fmt.Errorf("failed to insert styles: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Stored release ID: %d with artists, genres, and styles", release.Id)
	return nil
}

func insertRelease(tx *sql.Tx, release *models.Release) error {
	releaseQuery := fmt.Sprintf(insertReleaseSQL, releasesTableName)

	_, err := tx.Exec(releaseQuery, release.Id)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to insert release: %v", err)
	}
	return nil
}

func insertAttributes(tx *sql.Tx, releaseID int32, attributes []string, tableName string) error {
	attributeQuery := fmt.Sprintf(insertAttributeSQL, tableName)

	for _, attr := range attributes {
		_, err := tx.Exec(attributeQuery, releaseID, attr)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return fmt.Errorf("failed to insert into table %s: %v", tableName, err)
		}
	}

	return nil
}

func FetchReleaseCounts(db *sql.DB, artist, style, genre string) (models.CountResult, error) {
	query := fmt.Sprintf(fetchAttrsNamesSQL, releasesTableName, ArtistsTableName, StylesTableName, GenresTableName)

	args, query := createFilterQueries(artist, query, style, genre)

	query += " GROUP BY a.name, s.name, g.name"

	rows, err := db.Query(query, args...)
	if err != nil {
		return models.CountResult{}, fmt.Errorf("failed to execute query: %v", err)
	}

	defer rows.Close()

	countResult, err := fetchReleaseCountFromRows(rows)
	if err != nil {
		return models.CountResult{}, err
	}

	return countResult, nil
}

func fetchReleaseCountFromRows(rows *sql.Rows) (models.CountResult, error) {
	artistMap := make(map[string]int)
	styleMap := make(map[string]int)
	genreMap := make(map[string]int)
	totalReleases := 0

	for rows.Next() {
		var releaseId int
		var artistName, styleName, genreName sql.NullString

		if err := rows.Scan(&releaseId, &artistName, &styleName, &genreName); err != nil {
			return models.CountResult{}, fmt.Errorf("failed to scan row: %v", err)
		}

		totalReleases++

		if artistName.Valid && artistName.String != "" {
			artistMap[artistName.String]++
		}

		if styleName.Valid && styleName.String != "" {
			styleMap[styleName.String]++
		}

		if genreName.Valid && genreName.String != "" {
			genreMap[genreName.String]++
		}
	}

	artistCounts := make([]models.NameCount, 0, len(artistMap))
	for name, count := range artistMap {
		artistCounts = append(artistCounts, models.NameCount{Name: name, Count: count})
	}

	styleCounts := make([]models.NameCount, 0, len(styleMap))
	for name, count := range styleMap {
		styleCounts = append(styleCounts, models.NameCount{Name: name, Count: count})
	}

	genreCounts := make([]models.NameCount, 0, len(genreMap))
	for name, count := range genreMap {
		genreCounts = append(genreCounts, models.NameCount{Name: name, Count: count})
	}

	countResult := models.CountResult{
		ReleaseCount: totalReleases,
		ArtistCounts: artistCounts,
		StyleCounts:  styleCounts,
		GenreCounts:  genreCounts,
	}
	return countResult, nil
}

func createFilterQueries(artist string, query string, style string, genre string) ([]interface{}, string) {
	var args []interface{}
	argIndex := 1

	if artist != "" {
		query += fmt.Sprintf(" AND a.name ILIKE $%d", argIndex)
		args = append(args, "%"+artist+"%")
		argIndex++
	}
	if style != "" {
		query += fmt.Sprintf(" AND s.name ILIKE $%d", argIndex)
		args = append(args, "%"+style+"%")
		argIndex++
	}
	if genre != "" {
		query += fmt.Sprintf(" AND g.name ILIKE $%d", argIndex)
		args = append(args, "%"+genre+"%")
		argIndex++
	}
	return args, query
}

func FetchUniqueNames(db *sql.DB, tableName string) ([]*models.UniqueName, error) {
	query := fmt.Sprintf(fetchUniqueNamesSQL, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unique names from %s: %v", tableName, err)
	}
	defer rows.Close()

	var uniqueNames []*models.UniqueName
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan name: %v", err)
		}
		uniqueNames = append(uniqueNames, &models.UniqueName{Name: name})
	}

	return uniqueNames, nil
}
