package domain

import (
	"time"
)

type Address struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	PlaceID     string    `json:"place_id" gorm:"type:varchar(255);not null;uniqueIndex"`
	FullAddress string    `json:"full_address" gorm:"type:text;not null"`
	Country     string    `json:"country" gorm:"type:varchar(100);not null"`
	City        string    `json:"city" gorm:"type:varchar(100);not null"`
	District    *string   `json:"district" gorm:"type:varchar(100)"`
	Street      *string   `json:"street" gorm:"type:varchar(200)"`
	PostalCode  *string   `json:"postal_code" gorm:"type:varchar(20)"`
	Latitude    float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude   float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	DoorNumber  *string   `json:"door_number" gorm:"type:varchar(50)"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func NewAddress(placeID, fullAddress, country, city string, latitude, longitude float64) *Address {
	return &Address{
		PlaceID:     placeID,
		FullAddress: fullAddress,
		Country:     country,
		City:        city,
		Latitude:    latitude,
		Longitude:   longitude,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (a *Address) SetDistrict(district string) {
	a.District = &district
	a.UpdatedAt = time.Now()
}

func (a *Address) SetStreet(street string) {
	a.Street = &street
	a.UpdatedAt = time.Now()
}

func (a *Address) SetPostalCode(postalCode string) {
	a.PostalCode = &postalCode
	a.UpdatedAt = time.Now()
}

func (a *Address) SetDoorNumber(doorNumber string) {
	a.DoorNumber = &doorNumber
	a.UpdatedAt = time.Now()
}

func (a *Address) UpdateCoordinates(latitude, longitude float64) {
	a.Latitude = latitude
	a.Longitude = longitude
	a.UpdatedAt = time.Now()
}

func (a *Address) UpdateFromGooglePlaces(fullAddress, country, city string, district, street, postalCode *string, latitude, longitude float64) {
	a.FullAddress = fullAddress
	a.Country = country
	a.City = city
	a.District = district
	a.Street = street
	a.PostalCode = postalCode
	a.Latitude = latitude
	a.Longitude = longitude
	a.UpdatedAt = time.Now()
}

func (a *Address) ValidateRequiredFields() error {
	if a.PlaceID == "" {
		return ErrAddressPlaceIDRequired
	}
	if a.FullAddress == "" {
		return ErrAddressFullAddressRequired
	}
	if a.Country == "" {
		return ErrAddressCountryRequired
	}
	if a.City == "" {
		return ErrAddressCityRequired
	}
	if a.Latitude == 0 && a.Longitude == 0 {
		return ErrAddressCoordinatesRequired
	}
	return nil
}

func (a *Address) GetFormattedAddress() string {
	formatted := a.FullAddress
	if a.DoorNumber != nil && *a.DoorNumber != "" {
		formatted += " No: " + *a.DoorNumber
	}
	return formatted
}

func (a *Address) HasDoorNumber() bool {
	return a.DoorNumber != nil && *a.DoorNumber != ""
}

// Address domain errors
var (
	ErrAddressPlaceIDRequired     = NewDomainError("place ID is required")
	ErrAddressFullAddressRequired = NewDomainError("full address is required")
	ErrAddressCountryRequired     = NewDomainError("country is required")
	ErrAddressCityRequired        = NewDomainError("city is required")
	ErrAddressCoordinatesRequired = NewDomainError("latitude and longitude are required")
	ErrAddressNotFound            = NewDomainError("address not found")
	ErrAddressInvalidCoordinates  = NewDomainError("invalid coordinates")
)
