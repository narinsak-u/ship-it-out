package shipment

import (
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreatePersistsShipmentAndInitialEvent(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(0).(*models.Shipment).ID = 17
	})
	repo.On("CreateEvent", mock.MatchedBy(func(event *models.ShipmentEvent) bool {
		return event.ShipmentID == 17 && event.Status == "Label Created" && event.Description == "Awaiting pickup."
	})).Return(nil)
	service := NewService(repo, nil)

	shipment, err := service.Create(CreateRequest{
		Customer: models.ContactInfo{Name: "Sender", Province: "Bangkok", Zipcode: "10110"},
		Receiver: models.ContactInfo{Name: "Receiver", Province: "Chonburi", Zipcode: "20110"},
		Carrier: "Express", Weight: 2, Items: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, "pending", shipment.Status)
	assert.Equal(t, 0, shipment.Progress)
	repo.AssertExpectations(t)
}

func TestServiceCreateRejectsInvalidShipment(t *testing.T) {
	service := NewService(new(mockRepo), nil)
	_, err := service.Create(CreateRequest{})
	assert.EqualError(t, err, "customer name, province, and zipcode are required")
}
