package apitest

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	authbuild "github.com/himynamej/todo/api/services/auth/build/all"
	salesbuild "github.com/himynamej/todo/api/services/sales/build/all"
	"github.com/himynamej/todo/app/sdk/auth"
	"github.com/himynamej/todo/app/sdk/authclient"
	"github.com/himynamej/todo/app/sdk/mux"
	"github.com/himynamej/todo/business/domain/todobus/mocks"
	"github.com/himynamej/todo/business/sdk/dbtest"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {

	ctrl := gomock.NewController(t) // Create a new gomock controller
	defer ctrl.Finish()             // Ensures that all expected calls are made
	db := dbtest.New(t, testName)

	// -------------------------------------------------------------------------

	auth, err := auth.New(auth.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	// Create mock S3 and SQS clients using gomock
	mockS3Client := mocks.NewMockS3Client(ctrl)
	mockSQSClient := mocks.NewMockSQSClient(ctrl)
	mockS3Client.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-file-id", nil).AnyTimes() // Adjust as necessary

	mockSQSClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes() // You can adjust the return value and times as needed.
	// Construct the Todo business logic
	server := httptest.NewServer(mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		AuthConfig: mux.AuthConfig{
			Auth: auth,
		},
	}, authbuild.Routes()))

	authClient := authclient.New(db.Log, server.URL)

	// -------------------------------------------------------------------------

	mux := mux.WebAPI(mux.Config{
		Log:       db.Log,
		DB:        db.DB,
		S3Client:  mockS3Client,
		SQSClient: mockSQSClient,
		SalesConfig: mux.SalesConfig{
			AuthClient: authClient,
		},
	}, salesbuild.Routes())

	return &Test{
		DB:   db,
		Auth: auth,
		mux:  mux,
	}
}
