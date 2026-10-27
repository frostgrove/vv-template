package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrackingMode string

const (
	TrackingModeQuantity TrackingMode = "quantity"
	TrackingModeInstance TrackingMode = "instance"
)

type StockItem struct {
	Id               uuid.UUID
	Name             string
	Code             string
	VendorStockItems []*VendorStockItem
	TrackingMode     TrackingMode
}

type VendorStockItem struct {
	Id          uuid.UUID
	StockItemId uuid.UUID
	VendorId    uuid.UUID
	VendorCode  string

	StockItem StockItem
	Vendor    Vendor
}

type StockItemAvailability struct {
	Id                uuid.UUID
	VendorStockItemId uuid.UUID
	WarehouseId       uuid.UUID
	Stock             int

	VendorStockItem VendorStockItem
}

type StockItemAvailabilityInstance struct {
	Id        uuid.UUID
	ExpiresAt time.Time
}

type Warehouse struct {
	Id        uuid.UUID
	AddressId uuid.UUID
	Name      string
	Address   Address
}

type Address struct {
	Id         uuid.UUID
	Latitude   *float64
	Longitude  *float64
	Line1      string
	Line2      string
	City       string
	Region     string
	PostalCode string
	Country    string
}

type Vendor struct {
	Id   uuid.UUID
	Name *string
}
