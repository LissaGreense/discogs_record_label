package graphQL

import (
	"database/sql"
	"github.com/graphql-go/graphql"
)

func NewQueryType(db *sql.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"releaseCounts": &graphql.Field{
				Type: CountResultType,
				Args: graphql.FieldConfigArgument{
					"artist": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"style": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"genre": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: ReleaseCountsResolver(db),
			},
			"uniqueArtists": &graphql.Field{
				Type:    graphql.NewList(UniqueNameType),
				Resolve: UniqueArtistsResolver(db),
			},
			"uniqueGenres": &graphql.Field{
				Type:    graphql.NewList(UniqueNameType),
				Resolve: UniqueGenresResolver(db),
			},
			"uniqueStyles": &graphql.Field{
				Type:    graphql.NewList(UniqueNameType),
				Resolve: UniqueStylesResolver(db),
			},
		},
	})
}
