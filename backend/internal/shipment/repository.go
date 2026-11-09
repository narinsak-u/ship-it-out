package shipment


// ShipmentFilter holds all optional query parameters for listing shipments.
// Page and Limit control pagination (defaults: page=1, limit=10). Use limit=-1 for all rows.
// Search runs a case-insensitive ILIKE across order_id, tracking_number, customer_name, and destination.
// Status and ExcludeStatus filter/suppress by exact status match.
type ShipmentFilter struct {
	Page          int
	Limit         int
	Search        string
	Status        string
	ExcludeStatus string
}

