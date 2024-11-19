package todoapp

import "net/http"

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	filter := queryParams{
		ID: values.Get("user_id"),
	}

	return filter, nil
}
