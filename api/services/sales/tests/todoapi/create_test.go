package todoapi

import (
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/himynamej/todo/app/domain/todoapp"
	"github.com/himynamej/todo/app/sdk/apitest"
	"github.com/himynamej/todo/app/sdk/errs"
)

func createTodoItem200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/todo",
			Token:      sd.Admins[1].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &todoapp.NewTodoItem{
				Description: "Test Todo Item",
				DueDate:     time.Now().Add(72 * time.Hour).Format(time.RFC3339),
				FileID:      "mock-file-id",
			},
			GotResp: &todoapp.TodoItem{},
			ExpResp: &todoapp.TodoItem{
				Description: "Test Todo Item",
				DueDate:     time.Now().Add(72 * time.Hour).Format(time.RFC3339),
				FileID:      "mock-file-id",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*todoapp.TodoItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*todoapp.TodoItem)

				// Adjust dynamic fields
				expResp.ID = gotResp.ID

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func createTodoItem400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/v1/todo",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &todoapp.NewTodoItem{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"description\",\"error\":\"description is a required field\"},{\"field\":\"dueDate\",\"error\":\"dueDate is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func createTodoItem401() []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/todo",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "expected authorization header format: Bearer <token>"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "invalidtoken",
			URL:        "/v1/todo",
			Token:      "invalid-token",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
