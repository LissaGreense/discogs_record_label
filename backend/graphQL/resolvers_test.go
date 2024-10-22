package graphQL

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	params := graphql.Params{Schema: schema, RequestString: query}
	return graphql.Do(params)
}

func TestReleaseCountsResolver(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"releaseCount", "style", "genre", "artist"}).
		AddRow(5, "Rock", "Pop", "ArtistA"))

	query := NewQueryType(db)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	assert.NoError(t, err)

	queryString := `{
		releaseCounts(artist: "ArtistA", genre: "Pop", style: "Rock") {
			releaseCount
		}
	}`

	result := executeQuery(queryString, schema)

	assert.Nil(t, result.Errors)
	assert.Equal(t, 1, result.Data.(map[string]interface{})["releaseCounts"].(map[string]interface{})["releaseCount"])
}

func TestUniqueArtistsResolver(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("ArtistA").AddRow("ArtistB"))

	query := NewQueryType(db)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	assert.NoError(t, err)

	queryString := `{
		uniqueArtists {
			name
		}
	}`

	result := executeQuery(queryString, schema)

	assert.Nil(t, result.Errors)
	artists := result.Data.(map[string]interface{})["uniqueArtists"].([]interface{})
	assert.Len(t, artists, 2)
	assert.Equal(t, "ArtistA", artists[0].(map[string]interface{})["name"])
	assert.Equal(t, "ArtistB", artists[1].(map[string]interface{})["name"])
}

func TestUniqueGenresResolver(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Pop").AddRow("Rock"))

	query := NewQueryType(db)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	assert.NoError(t, err)

	queryString := `{
		uniqueGenres {
			name
		}
	}`

	result := executeQuery(queryString, schema)

	assert.Nil(t, result.Errors)
	genres := result.Data.(map[string]interface{})["uniqueGenres"].([]interface{})
	assert.Len(t, genres, 2)
	assert.Equal(t, "Pop", genres[0].(map[string]interface{})["name"])
	assert.Equal(t, "Rock", genres[1].(map[string]interface{})["name"])
}

func TestUniqueStylesResolver(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Jazz").AddRow("Blues"))

	query := NewQueryType(db)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	assert.NoError(t, err)

	queryString := `{
		uniqueStyles {
			name
		}
	}`

	result := executeQuery(queryString, schema)

	assert.Nil(t, result.Errors)
	styles := result.Data.(map[string]interface{})["uniqueStyles"].([]interface{})
	assert.Len(t, styles, 2)
	assert.Equal(t, "Jazz", styles[0].(map[string]interface{})["name"])
	assert.Equal(t, "Blues", styles[1].(map[string]interface{})["name"])
}
