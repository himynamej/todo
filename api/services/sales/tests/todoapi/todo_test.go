package todoapi

import (
	"testing"

	"github.com/himynamej/todo/app/sdk/apitest"
)

func Test_TodoApp(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_TodoApp")

	// Seed data if necessary
	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------
	// Run test cases for CreateTodoItem
	// -------------------------------------------------------------------------

	test.Run(t, createTodoItem200(sd), "createtodoitem-200")
	test.Run(t, createTodoItem400(sd), "createtodoitem-400")
	test.Run(t, createTodoItem401(), "createtodoitem-401")

	// -------------------------------------------------------------------------
	// Run test cases for File Upload and Download
	// -------------------------------------------------------------------------

}
