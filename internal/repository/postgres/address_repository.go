package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &addressRepository{db: db}
}

// Basic CRUD operations
func (r *addressRepository) Create(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *addressRepository) GetByID(ctx context.Context, id int) (*domain.Address, error) {
	var address domain.Address
	err := r.db.WithContext(ctx).First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *addressRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Address{}, id).Error
}

// Google Places specific operations
func (r *addressRepository) GetByPlaceID(ctx context.Context, placeID string) (*domain.Address, error) {
	var address domain.Address
	err := r.db.WithContext(ctx).Where("place_id = ?", placeID).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) ExistsByPlaceID(ctx context.Context, placeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Address{}).Where("place_id = ?", placeID).Count(&count).Error
	return count > 0, err
}

func (r *addressRepository) CreateOrUpdateByPlaceID(ctx context.Context, address *domain.Address) (*domain.Address, error) {
	var existingAddress domain.Address
	err := r.db.WithContext(ctx).Where("place_id = ?", address.PlaceID).First(&existingAddress).Error

	if err == gorm.ErrRecordNotFound {
		// Create new address
		if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
			return nil, err
		}
		return address, nil
	} else if err != nil {
		return nil, err
	}

	// Update existing address
	existingAddress.FullAddress = address.FullAddress
	existingAddress.Country = address.Country
	existingAddress.City = address.City
	existingAddress.District = address.District
	existingAddress.Street = address.Street
	existingAddress.PostalCode = address.PostalCode
	existingAddress.Latitude = address.Latitude
	existingAddress.Longitude = address.Longitude
	existingAddress.DoorNumber = address.DoorNumber

	if err := r.db.WithContext(ctx).Save(&existingAddress).Error; err != nil {
		return nil, err
	}

	return &existingAddress, nil
}

// Search operations
func (r *addressRepository) SearchByCity(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error) {
	var addresses []*domain.Address
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Address{}).Where("city ILIKE ?", "%"+city+"%")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("city ASC, full_address ASC").
		Find(&addresses).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return addresses, paginationResponse, nil
}

func (r *addressRepository) SearchByCountry(ctx context.Context, country string, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error) {
	var addresses []*domain.Address
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Address{}).Where("country ILIKE ?", "%"+country+"%")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("country ASC, city ASC, full_address ASC").
		Find(&addresses).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return addresses, paginationResponse, nil
}

func (r *addressRepository) SearchByFullAddress(ctx context.Context, query string, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error) {
	var addresses []*domain.Address
	var total int64

	searchQuery := r.db.WithContext(ctx).Model(&domain.Address{}).
		Where("full_address ILIKE ? OR city ILIKE ? OR country ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")

	// Count total
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := searchQuery.
		Offset(offset).
		Limit(pageSize).
		Order("country ASC, city ASC, full_address ASC").
		Find(&addresses).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return addresses, paginationResponse, nil
}

// Location-based operations
func (r *addressRepository) GetByCoordinates(ctx context.Context, latitude, longitude float64, radiusKm int, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error) {
	var addresses []*domain.Address
	var total int64

	// Using Haversine formula for distance calculation
	query := r.db.WithContext(ctx).Model(&domain.Address{}).
		Where("(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) <= ?",
			latitude, longitude, latitude, radiusKm)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("(6371 * acos(cos(radians(" + fmt.Sprintf("%f", latitude) + ")) * cos(radians(latitude)) * cos(radians(longitude) - radians(" + fmt.Sprintf("%f", longitude) + ")) + sin(radians(" + fmt.Sprintf("%f", latitude) + ")) * sin(radians(latitude)))) ASC").
		Find(&addresses).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return addresses, paginationResponse, nil
}

func (r *addressRepository) GetNearbyAddresses(ctx context.Context, latitude, longitude float64, radiusKm int, limit int) ([]*domain.Address, error) {
	var addresses []*domain.Address

	err := r.db.WithContext(ctx).
		Where("(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) <= ?",
			latitude, longitude, latitude, radiusKm).
		Order("(6371 * acos(cos(radians(" + fmt.Sprintf("%f", latitude) + ")) * cos(radians(latitude)) * cos(radians(longitude) - radians(" + fmt.Sprintf("%f", longitude) + ")) + sin(radians(" + fmt.Sprintf("%f", latitude) + ")) * sin(radians(latitude)))) ASC").
		Limit(limit).
		Find(&addresses).Error

	return addresses, err
}

// Statistics operations
func (r *addressRepository) CountByCity(ctx context.Context, city string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Address{}).
		Where("city ILIKE ?", "%"+city+"%").
		Count(&count).Error
	return count, err
}

func (r *addressRepository) CountByCountry(ctx context.Context, country string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Address{}).
		Where("country ILIKE ?", "%"+country+"%").
		Count(&count).Error
	return count, err
}

func (r *addressRepository) GetPopularCities(ctx context.Context, limit int) ([]string, error) {
	var cities []string
	err := r.db.WithContext(ctx).Model(&domain.Address{}).
		Select("city, COUNT(*) as count").
		Group("city").
		Order("count DESC").
		Limit(limit).
		Pluck("city", &cities).Error
	return cities, err
}

func (r *addressRepository) GetPopularCountries(ctx context.Context, limit int) ([]string, error) {
	var countries []string
	err := r.db.WithContext(ctx).Model(&domain.Address{}).
		Select("country, COUNT(*) as count").
		Group("country").
		Order("count DESC").
		Limit(limit).
		Pluck("country", &countries).Error
	return countries, err
}

// Validation operations
func (r *addressRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Address{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *addressRepository) ValidateCoordinates(ctx context.Context, latitude, longitude float64) (bool, error) {
	// Basic coordinate validation
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return false, nil
	}
	return true, nil
}

// Bulk operations
func (r *addressRepository) GetMultipleByIDs(ctx context.Context, ids []int) ([]*domain.Address, error) {
	var addresses []*domain.Address
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) CreateMultiple(ctx context.Context, addresses []*domain.Address) error {
	return r.db.WithContext(ctx).Create(&addresses).Error
}

func (r *addressRepository) DeleteMultiple(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Address{}).Error
}

// Advanced filtering
func (r *addressRepository) GetAddressesWithFilters(ctx context.Context, filters dto.AddressFilterRequest, pagination dto.PaginationRequest) ([]*domain.Address, *dto.PaginationResponse, error) {
	var addresses []*domain.Address
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Address{})

	// Apply filters
	if filters.City != nil && *filters.City != "" {
		query = query.Where("city ILIKE ?", "%"+*filters.City+"%")
	}
	if filters.Country != nil && *filters.Country != "" {
		query = query.Where("country ILIKE ?", "%"+*filters.Country+"%")
	}
	if filters.District != nil && *filters.District != "" {
		query = query.Where("district ILIKE ?", "%"+*filters.District+"%")
	}
	if filters.Query != nil && *filters.Query != "" {
		query = query.Where("(full_address ILIKE ? OR city ILIKE ? OR country ILIKE ?)",
			"%"+*filters.Query+"%", "%"+*filters.Query+"%", "%"+*filters.Query+"%")
	}

	// Coordinate range filters
	if filters.MinLatitude != nil {
		query = query.Where("latitude >= ?", *filters.MinLatitude)
	}
	if filters.MaxLatitude != nil {
		query = query.Where("latitude <= ?", *filters.MaxLatitude)
	}
	if filters.MinLongitude != nil {
		query = query.Where("longitude >= ?", *filters.MinLongitude)
	}
	if filters.MaxLongitude != nil {
		query = query.Where("longitude <= ?", *filters.MaxLongitude)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated results
	offset := pagination.GetOffset()
	pageSize := pagination.GetPageSizeWithDefault()

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("country ASC, city ASC, full_address ASC").
		Find(&addresses).Error

	if err != nil {
		return nil, nil, err
	}

	paginationResponse := &dto.PaginationResponse{
		Page:       pagination.GetPageWithDefault(),
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return addresses, paginationResponse, nil
}
