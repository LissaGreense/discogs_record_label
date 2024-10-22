package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mockedReleaseJSON = `{
		"id": 123456,
		"artists": [{"name": "Some Artist"}],
		"styles": ["Rock"],
		"genres": ["Pop"]
	}`

	mockedReleasesJSON = `{
		"releases": [{
			"resource_url": "https://api.discogs.com/releases/123456"
		}],
		"pagination": {
			"urls": {
				"next": null
			}
		}
	}`
)

func TestParseReleaseResponse(t *testing.T) {
	body := []byte(mockedReleaseJSON)
	release, err := parseReleaseResponse(nil, body)

	require.NoError(t, err)
	assert.Equal(t, int32(123456), release.Id)
	assert.Equal(t, []string{"Some Artist"}, release.Artists)
	assert.Equal(t, []string{"Rock"}, release.Styles)
	assert.Equal(t, []string{"Pop"}, release.Genres)
}

func TestParseReleases(t *testing.T) {
	body := []byte(mockedReleasesJSON)
	releaseUrls, err := parseReleases(body)

	require.NoError(t, err)
	assert.Equal(t, 1, len(releaseUrls))
	assert.Equal(t, "https://api.discogs.com/releases/123456", releaseUrls[0])
}

func TestGetNextPageURL(t *testing.T) {
	body := []byte(mockedReleasesJSON)
	nextURL, hasNext := getNextPageURL(body)

	assert.False(t, hasNext)
	assert.Empty(t, nextURL)
}

func TestExtractArtists(t *testing.T) {
	releaseMap := map[string]interface{}{
		"artists": []interface{}{
			map[string]interface{}{"name": "Artist 1"},
			map[string]interface{}{"name": "Artist 2"},
		},
	}

	artists := extractArtists(releaseMap)
	assert.Equal(t, []string{"Artist 1", "Artist 2"}, artists)
}

func TestExtractStyles(t *testing.T) {
	releaseMap := map[string]interface{}{
		"styles": []interface{}{
			"Style 1",
			"Style 2",
		},
	}

	styles := extractStyles(releaseMap)
	assert.Equal(t, []string{"Style 1", "Style 2"}, styles)
}

func TestExtractGenres(t *testing.T) {
	releaseMap := map[string]interface{}{
		"genres": []interface{}{
			"Genre 1",
			"Genre 2",
		},
	}

	genres := extractGenres(releaseMap)
	assert.Equal(t, []string{"Genre 1", "Genre 2"}, genres)
}
