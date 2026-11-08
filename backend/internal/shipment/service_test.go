package shipment

import (
	"context"
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreatePersistsShipmentAndInitialEvent(t *testing.T) {
	repo := new(mockRepo)
	repo.On("CreateWithEvent", mock.Anything, mock.MatchedBy(func(event *models.ShipmentEvent) bool {
		return event.Status == "Label Created" && event.Description == "Awaiting pickup."
	})).Return(nil).Run(func(args mock.Arguments) {
		args.Get(0).(*models.Shipment).ID = 17
		args.Get(1).(*models.ShipmentEvent).ShipmentID = 17
	})
	service := NewService(repo, nil)
	shipment, err := service.Create(context.Background(), CreateRequest{
		Customer: models.ContactInfo{Name: "Sender", Province: "Bangkok", Zipcode: "10110"},
		Receiver: models.ContactInfo{Name: "Receiver", Province: "Chonburi", Zipcode: "20110"},
		Carrier: "Express", Weight: 2, Items: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, "pending", shipment.Status)
	assert.Equal(t, float64(0), shipment.Progress)
	repo.AssertExpectations(t)
}

func TestServiceCreateRejectsInvalidShipment(t *testing.T) {
	service := NewService(new(mockRepo), nil)
	_, err := service.Create(context.Background(), CreateRequest{})
	assert.EqualError(t, err, "customer name, province, and zipcode are required")
}
