package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LissaGreense/discogs_record_label/backend/models"
)

func TestStoreRelease(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	release := &models.Release{
		Id:      1,
		Artists: []string{"Artist 1"},
		Genres:  []string{"Genre 1"},
		Styles:  []string{"Style 1"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO releases").WithArgs(release.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO artists").WithArgs(release.Id, release.Artists[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO genres").WithArgs(release.Id, release.Genres[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO styles").WithArgs(release.Id, release.Styles[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := StoreRelease(db, release); err != nil {
		t.Fatalf("failed to store release: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchReleaseCounts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"releaseCount", "artistName", "styleName", "genreName"}).
		AddRow(1, "SomeArtist", "SomeStyle", "SomeGenre").
		AddRow(2, "SomeArtist", "SomeStyleTwo", "SomeGenreTwo")

	query := `SELECT COUNT\(DISTINCT r.id\) as releaseCount, 
                    COALESCE\(a.name, ''\) as artistName, 
                    COALESCE\(s.name, ''\) as styleName, 
                    COALESCE\(g.name, ''\) as genreName 
              FROM releases r 
              LEFT JOIN artists a ON r.id = a.release_id 
              LEFT JOIN styles s ON r.id = s.release_id 
              LEFT JOIN genres g ON r.id = g.release_id 
              WHERE 1=1 
              AND a.name ILIKE \$1 
              GROUP BY a.name, s.name, g.name`

	mock.ExpectQuery(query).WithArgs("%SomeArtist%").WillReturnRows(rows)

	countResult, err := FetchReleaseCounts(db, "SomeArtist", "", "")
	if err != nil {
		t.Fatalf("failed to fetch release counts: %v", err)
	}
	if countResult.ReleaseCount != 2 {
		t.Fatalf("expected release count 2, got %d", countResult.ReleaseCount)
	}

	if len(countResult.ArtistCounts) != 1 {
		t.Fatalf("expected 1 ArtistCounts object, got %d", len(countResult.ArtistCounts))
	}
	if countResult.ArtistCounts[0].Name != "SomeArtist" || countResult.ArtistCounts[0].Count != 2 {
		t.Errorf("expected ArtistCounts to have name 'SomeArtist' and count 1, got name '%s' and count %d", countResult.ArtistCounts[0].Name, countResult.ArtistCounts[0].Count)
	}

	if len(countResult.StyleCounts) != 2 {
		t.Fatalf("expected 2 StyleCounts, got %d", len(countResult.StyleCounts))
	}
	expectedStyles := map[string]int{
		"SomeStyle":    1,
		"SomeStyleTwo": 1,
	}
	for _, style := range countResult.StyleCounts {
		if count, ok := expectedStyles[style.Name]; ok {
			if style.Count != count {
				t.Errorf("expected style '%s' count %d, got %d", style.Name, count, style.Count)
			}
		} else {
			t.Errorf("unexpected style name '%s'", style.Name)
		}
	}

	if len(countResult.GenreCounts) != 2 {
		t.Fatalf("expected 2 GenreCounts, got %d", len(countResult.GenreCounts))
	}
	expectedGenres := map[string]int{
		"SomeGenre":    1,
		"SomeGenreTwo": 1,
	}
	for _, genre := range countResult.GenreCounts {
		if count, ok := expectedGenres[genre.Name]; ok {
			if genre.Count != count {
				t.Errorf("expected genre '%s' count %d, got %d", genre.Name, count, genre.Count)
			}
		} else {
			t.Errorf("unexpected genre name '%s'", genre.Name)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("ArtistOne").
		AddRow("ArtistTwo").
		AddRow("ArtistThree")

	mock.ExpectQuery("SELECT DISTINCT name FROM artists ORDER BY name").WillReturnRows(rows)

	uniqueNames, err := FetchUniqueNames(db, ArtistsTableName)
	if err != nil {
		t.Fatalf("failed to fetch unique names: %v", err)
	}
	if len(uniqueNames) != 3 {
		t.Fatalf("expected 3 unique artist, got %d", len(uniqueNames))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
