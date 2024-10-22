package graphQL

import (
	"database/sql"
	"github.com/LissaGreense/discogs_record_label/backend/storage"
	"github.com/graphql-go/graphql"
)

func ReleaseCountsResolver(db *sql.DB) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		var artist, style, genre string

		if artistArg, ok := params.Args["artist"].(string); ok {
			artist = artistArg
		}
		if styleArg, ok := params.Args["style"].(string); ok {
			style = styleArg
		}
		if genreArg, ok := params.Args["genre"].(string); ok {
			genre = genreArg
		}

		return storage.FetchReleaseCounts(db, artist, style, genre)
	}
}

func UniqueArtistsResolver(db *sql.DB) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		return storage.FetchUniqueNames(db, storage.ArtistsTableName)
	}
}

func UniqueGenresResolver(db *sql.DB) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		return storage.FetchUniqueNames(db, storage.GenresTableName)
	}
}

func UniqueStylesResolver(db *sql.DB) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		return storage.FetchUniqueNames(db, storage.StylesTableName)
	}
}
