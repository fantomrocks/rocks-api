package handlers

import (
	"encoding/json"
	"fantomrocks-api/internal/services"
	"github.com/graph-gophers/graphql-go"
	"net/http"
)

// define GraphQL query parameters structure
type QueryParams struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// Get new GraphQL HTTP leaf handler.
func GraphQLHandler(log services.Logger, schema *graphql.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// try to extract the request details from the HTTP request struct
		params := &QueryParams{}
		if err := json.NewDecoder(r.Body).Decode(params); err != nil {
			log.Errorf("GQL->ServeHTTP(): Request could not be decoded. Probably not a GraphQL request.", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// make sure to pass the current common to the resolvers
		ctx := r.Context()

		// get the response
		res, err := json.Marshal(schema.Exec(ctx, params.Query, params.OperationName, params.Variables))
		if err != nil {
			log.Criticalf("GQL->ServeHTTP(): Response could not be encoded to JSON.", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write the response
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(res); err != nil {
			log.Errorf("GQL->ServeHTTP(): Can not send response to remote client.", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
