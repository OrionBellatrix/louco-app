package repository

import (
	"context"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
)

type AddressRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, address *domain.Address) error
	GetByID(ctx context.Context, id int) (*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id int) error

	// Google Places specific operations
	GetByPlaceID(ctx context.Context, placeID string) (*domain.Address, error)
	ExistsByPlaceID(ctx context.Context, placeID string) (bool, error)
	CreateOrUpdateByPlaceID(ctx context.Context, address *domain.Address) (*domain.Address, error)

	// Search operations
	SearchByCity(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error)
	SearchByCountry(ctx context.Context, country string, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error)
	SearchByFullAddress(ctx context.Context, query string, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error)

	// Location-based operations
	GetByCoordinates(ctx context.Context, latitude, longitude float64, radiusKm int, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error)
	GetNearbyAddresses(ctx context.Context, latitude, longitude float64, radiusKm int, limit int) ([]*domain.Address, error)

	// Statistics operations
	CountByCity(ctx context.Context, city string) (int64, error)
	CountByCountry(ctx context.Context, country string) (int64, error)
	GetPopularCities(ctx context.Context, limit int) ([]string, error)
	GetPopularCountries(ctx context.Context, limit int) ([]string, error)

	// Validation operations
	ExistsByID(ctx context.Context, id int) (bool, error)
	ValidateCoordinates(ctx context.Context, latitude, longitude float64) (bool, error)

	// Bulk operations
	GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Address, error)
	CreateMultiple(ctx context.Context, addresses []*domain.Address) error
	DeleteMultiple(ctx context.Context, ids []int) error

	// Advanced filtering
	GetAddressesWithFilters(ctx context.Context, filters dto.AddressFilterRequest, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error)
}
