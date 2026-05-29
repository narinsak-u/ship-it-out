// Package utils provides common utilities for the backend.
package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
)

var orderIDMutex sync.Mutex

// GenerateOrderID creates a human-readable order ID like "ORD-10251".
// Uses mutex to prevent race conditions under concurrent requests.
func GenerateOrderID() string {
	orderIDMutex.Lock()
	defer orderIDMutex.Unlock()

	var shipments []models.Shipment
	database.DB.Select("order_id").Find(&shipments)
	maxNum := 10245
	for _, s := range shipments {
		parts := strings.SplitN(s.OrderID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	return fmt.Sprintf("ORD-%d", maxNum+1)
}
