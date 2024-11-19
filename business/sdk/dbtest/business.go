package dbtest

import (
	"time"

	"github.com/golang/mock/gomock"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/business/domain/todobus/mocks"
	"github.com/himynamej/todo/business/domain/todobus/stores/itemdb"
	"github.com/himynamej/todo/business/domain/userbus"
	"github.com/himynamej/todo/business/domain/userbus/stores/usercache"
	"github.com/himynamej/todo/business/domain/userbus/stores/userdb"
	"github.com/himynamej/todo/business/sdk/delegate"
	"github.com/himynamej/todo/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain APIs needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	User     *userbus.Business
	Todo     *todobus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB, ctrl *gomock.Controller) BusDomain {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, usercache.NewStore(log, userdb.NewStore(log, db), time.Hour))

	// Create mocked dependencies for Todo
	todostore := itemdb.NewStore(log, db)
	mockSQSClient := mocks.NewMockSQSClient(ctrl)
	mockS3Client := mocks.NewMockS3Client(ctrl)

	mockS3Client.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("mock-file-id", nil).AnyTimes() // Adjust as necessary

	mockSQSClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes() // You can adjust the return value and times as needed.
	// Construct the Todo business logic

	todoBus := todobus.NewBusiness(log, todostore, mockSQSClient, mockS3Client)

	return BusDomain{
		Delegate: delegate,
		User:     userBus,
		Todo:     todoBus,
	}
}
