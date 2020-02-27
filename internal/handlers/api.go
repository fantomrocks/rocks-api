package handlers

import (
	"fantomrocks-api/internal/common"
	"fantomrocks-api/internal/graphql/resolvers"
	gqlschema "fantomrocks-api/internal/graphql/schema"
	"fantomrocks-api/internal/repository"
	"fantomrocks-api/internal/services"
	"github.com/graph-gophers/graphql-go"
	"net/http"
)

// Construct and return the GraphQL API handler.
func ApiHandler(cfg *common.Config, repo *repository.Repository, log services.Logger) http.Handler {
	// we don't want to write a method for each type field if it could be matched directly
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	// create new parsed GraphQL schema
	schema := graphql.MustParseSchema(gqlschema.GetSchema(), resolvers.NewResolver(repo, log), opts...)

	// construct handlers chain for the API endpoint
	return LoggingHandler(log, CORSHandler(log, &CORSOptions{
		AllowOrigins:     cfg.Cors,
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400,
	}, GraphQLHandler(log, schema)))
}
