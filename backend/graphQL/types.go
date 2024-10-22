package graphQL

import "github.com/graphql-go/graphql"

var UniqueNameType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UniqueName",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var CountResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CountResult",
	Fields: graphql.Fields{
		"releaseCount": &graphql.Field{
			Type: graphql.Int,
		},
		"styleCounts": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "StyleCount",
				Fields: graphql.Fields{
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"count": &graphql.Field{
						Type: graphql.Int,
					},
				},
			})),
		},
		"genreCounts": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "GenreCount",
				Fields: graphql.Fields{
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"count": &graphql.Field{
						Type: graphql.Int,
					},
				},
			})),
		},
		"artistCounts": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "ArtistCount",
				Fields: graphql.Fields{
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"count": &graphql.Field{
						Type: graphql.Int,
					},
				},
			})),
		},
	},
})
