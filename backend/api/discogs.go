package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/LissaGreense/discogs_record_label/backend/models"
	"github.com/LissaGreense/discogs_record_label/backend/storage"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"time"
)

const (
	discogsLabelAPIURL = "https://api.discogs.com/labels/%d/releases?page=1&per_page=%d"
	perPage            = 100
)

func FetchAndStoreReleases(db *sql.DB, labelID int) error {
	client := createDiscogsClient()

	releaseUrls, err := getReleasesURLs(labelID, client)
	if err != nil {
		return err
	}

	err = fetchReleasesDetailAndSave(db, releaseUrls, client)
	if err != nil {
		return err
	}

	return nil
}

func createDiscogsClient() *resty.Client {
	client := resty.New()

	apiAppName := os.Getenv("DISCOGS_APP_NAME")
	discogsKey := os.Getenv("DISCOGS_KEY")
	discogsSecret := os.Getenv("DISCOGS_SECRET")

	if apiAppName != "" {
		client.SetHeader("User-Agent", apiAppName)
	}

	if discogsSecret != "" && discogsKey != "" {
		client.SetHeader("Authorization", fmt.Sprintf("Discogs key=%s, secret=%s", discogsKey, discogsSecret))
	}

	return client
}

func fetchReleasesDetailAndSave(db *sql.DB, releaseUrls []string, client *resty.Client) error {
	for releaseIndex := 0; releaseIndex < len(releaseUrls); {
		resp, err := client.R().
			SetHeader("Accept", "application/json").
			Get(releaseUrls[releaseIndex])

		if err != nil || resp.StatusCode() == 429 {
			handleRequestError(resp, err)
			time.Sleep(60 * time.Second)
			continue
		}
		releaseIndex++

		release, err := parseReleaseResponse(err, resp.Body())
		if err != nil {
			return err
		}

		err = storage.StoreRelease(db, release)
		if err != nil {
			log.Printf("Error storing release %d: %v", release.Id, err)
		}
	}
	return nil
}

func parseReleaseResponse(err error, body []byte) (*models.Release, error) {
	var releaseFromBody map[string]interface{}
	err = json.Unmarshal(body, &releaseFromBody)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v, body: %v", err, body)
	}

	release := &models.Release{
		Id:      int32(releaseFromBody["id"].(float64)),
		Artists: extractArtists(releaseFromBody),
		Styles:  extractStyles(releaseFromBody),
		Genres:  extractGenres(releaseFromBody),
	}
	return release, nil
}

func handleRequestError(resp *resty.Response, err error) {
	if err != nil {
		log.Printf("Failed to fetch release details: %v", err)
	}
	if resp != nil && resp.StatusCode() == 429 {
		log.Println("Received 429 Too Many Requests. Retrying after a delay...")
		headers := resp.Header()
		log.Printf("Rate Limiting: {X-Discogs-Ratelimit: %v, X-Discogs-Ratelimit-Used: %v, X-Discogs-Ratelimit-Remaining: %v}\n",
			headers.Get("X-Discogs-Ratelimit"), headers.Get("X-Discogs-Ratelimit-Used"), headers.Get("X-Discogs-Ratelimit-Remaining"))
	}
}

func getReleasesURLs(labelID int, client *resty.Client) ([]string, error) {
	var releaseUrls []string
	url := fmt.Sprintf(discogsLabelAPIURL, labelID, perPage)

	for {
		resp, err := client.R().
			SetHeader("Accept", "application/json").
			Get(url)

		if err != nil || resp.StatusCode() == 429 {
			handleRequestError(resp, err)
			time.Sleep(60 * time.Second)
			continue
		}

		releases, err := parseReleases(resp.Body())
		if err != nil {
			return nil, err
		}
		releaseUrls = append(releaseUrls, releases...)

		nextPageURL, hasNext := getNextPageURL(resp.Body())
		if !hasNext {
			break
		}
		url = nextPageURL
	}
	return releaseUrls, nil
}

func parseReleases(body []byte) ([]string, error) {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	items, ok := result["releases"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected format of releases data")
	}

	var releaseUrls []string
	for _, item := range items {
		releaseMap := item.(map[string]interface{})
		releaseUrl := releaseMap["resource_url"].(string)
		releaseUrls = append(releaseUrls, releaseUrl)
	}

	return releaseUrls, nil
}

func getNextPageURL(body []byte) (string, bool) {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", false
	}

	pagination, ok := result["pagination"].(map[string]interface{})
	if !ok {
		return "", false
	}
	paginationUrls, ok := pagination["urls"].(map[string]interface{})
	if !ok || paginationUrls["next"] == nil {
		return "", false
	}
	return paginationUrls["next"].(string), true
}

func extractArtists(releaseMap map[string]interface{}) []string {
	artistsRaw, ok := releaseMap["artists"].([]interface{})
	if ok && len(artistsRaw) > 0 {
		var artists []string
		for _, artistInterface := range artistsRaw {
			artistMap, ok := artistInterface.(map[string]interface{})
			if ok {
				name, ok := artistMap["name"].(string)
				if ok {
					artists = append(artists, name)
				}
			}
		}
		return artists
	}
	return nil
}

func extractStyles(releaseMap map[string]interface{}) []string {
	stylesRaw, ok := releaseMap["styles"].([]interface{})
	if ok && len(stylesRaw) > 0 {
		var styles []string
		for _, styleInterface := range stylesRaw {
			style, ok := styleInterface.(string)
			if ok {
				styles = append(styles, style)
			}
		}
		return styles
	}
	return nil
}

func extractGenres(releaseMap map[string]interface{}) []string {
	genresRaw, ok := releaseMap["genres"].([]interface{})
	if ok && len(genresRaw) > 0 {
		var genres []string
		for _, genreInterface := range genresRaw {
			genre, ok := genreInterface.(string)
			if ok {
				genres = append(genres, genre)
			}
		}
		return genres
	}
	return nil
}
