package service

import (
	"context"
	"fmt"

	"github.com/louco-event/internal/domain"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/repository"
	"github.com/rs/zerolog"
)

type AddressService interface {
	// Basic CRUD operations
	CreateAddress(ctx context.Context, req dto.CreateAddressRequest) (*dto.AddressResponse, error)
	GetAddressByID(ctx context.Context, id int) (*dto.AddressResponse, error)
	UpdateAddress(ctx context.Context, id int, req dto.CreateAddressRequest) (*dto.AddressResponse, error)
	DeleteAddress(ctx context.Context, id int) error

	// Google Places integration
	GetAddressByPlaceID(ctx context.Context, placeID string) (*dto.AddressResponse, error)
	CreateOrUpdateByPlaceID(ctx context.Context, req dto.CreateAddressRequest) (*dto.AddressResponse, error)

	// Location-based operations
	GetAddressesByCity(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error)
	GetAddressesByCountry(ctx context.Context, country string, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error)
	GetAddressesByCoordinates(ctx context.Context, latitude, longitude float64, radiusKm int, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error)
	GetNearbyAddresses(ctx context.Context, latitude, longitude float64, radiusKm int, limit int) ([]*dto.AddressResponse, error)

	// Search and filtering
	SearchAddresses(ctx context.Context, query string, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error)
	GetAddressesWithFilters(ctx context.Context, filters dto.AddressFilterRequest, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error)

	// Validation operations
	ValidateAddress(ctx context.Context, address *domain.Address) error
	ValidateCoordinates(ctx context.Context, latitude, longitude float64) error

	// Statistics operations
	GetPopularCities(ctx context.Context, limit int) ([]string, error)
	GetPopularCountries(ctx context.Context, limit int) ([]string, error)
}

type addressService struct {
	addressRepo repository.AddressRepository
	logger      zerolog.Logger
}

func NewAddressService(addressRepo repository.AddressRepository, logger zerolog.Logger) AddressService {
	return &addressService{
		addressRepo: addressRepo,
		logger:      logger.With().Str("service", "address").Logger(),
	}
}

// Basic CRUD operations
func (s *addressService) CreateAddress(ctx context.Context, req dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	// Validate request
	if err := s.validateCreateAddressRequest(&req); err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Invalid create address request")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create domain entity
	address := &domain.Address{
		PlaceID:     req.PlaceID,
		FullAddress: req.FullAddress,
		Country:     req.Country,
		City:        req.City,
		District:    req.District,
		Street:      req.Street,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		DoorNumber:  req.DoorNumber,
	}

	// Validate domain entity
	if err := s.ValidateAddress(ctx, address); err != nil {
		return nil, err
	}

	// Create address
	if err := s.addressRepo.Create(ctx, address); err != nil {
		s.logger.Error().Err(err).Interface("address", address).Msg("Failed to create address")
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	s.logger.Info().Int("address_id", address.ID).Str("place_id", address.PlaceID).Msg("Address created successfully")

	return s.addressToResponse(address), nil
}

func (s *addressService) GetAddressByID(ctx context.Context, id int) (*dto.AddressResponse, error) {
	address, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("address_id", id).Msg("Failed to get address by ID")
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	return s.addressToResponse(address), nil
}

func (s *addressService) UpdateAddress(ctx context.Context, id int, req dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	// Get existing address
	existingAddress, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing address: %w", err)
	}

	// Validate request
	if err := s.validateCreateAddressRequest(&req); err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Invalid update address request")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update fields
	existingAddress.PlaceID = req.PlaceID
	existingAddress.FullAddress = req.FullAddress
	existingAddress.Country = req.Country
	existingAddress.City = req.City
	existingAddress.District = req.District
	existingAddress.Street = req.Street
	existingAddress.PostalCode = req.PostalCode
	existingAddress.Latitude = req.Latitude
	existingAddress.Longitude = req.Longitude
	existingAddress.DoorNumber = req.DoorNumber

	// Validate updated address
	if err := s.ValidateAddress(ctx, existingAddress); err != nil {
		return nil, err
	}

	// Update address
	if err := s.addressRepo.Update(ctx, existingAddress); err != nil {
		s.logger.Error().Err(err).Int("address_id", id).Msg("Failed to update address")
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	s.logger.Info().Int("address_id", id).Msg("Address updated successfully")

	return s.addressToResponse(existingAddress), nil
}

func (s *addressService) DeleteAddress(ctx context.Context, id int) error {
	// Check if address exists
	exists, err := s.addressRepo.ExistsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check address existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("address not found")
	}

	// Note: In a real implementation, you would check if address is being used by events
	// For now, we'll skip this check since the repository doesn't have this method

	// Delete address
	if err := s.addressRepo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("address_id", id).Msg("Failed to delete address")
		return fmt.Errorf("failed to delete address: %w", err)
	}

	s.logger.Info().Int("address_id", id).Msg("Address deleted successfully")
	return nil
}

// Google Places integration
func (s *addressService) GetAddressByPlaceID(ctx context.Context, placeID string) (*dto.AddressResponse, error) {
	address, err := s.addressRepo.GetByPlaceID(ctx, placeID)
	if err != nil {
		s.logger.Error().Err(err).Str("place_id", placeID).Msg("Failed to get address by place ID")
		return nil, fmt.Errorf("failed to get address by place ID: %w", err)
	}

	return s.addressToResponse(address), nil
}

func (s *addressService) CreateOrUpdateByPlaceID(ctx context.Context, req dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	// Validate request
	if err := s.validateCreateAddressRequest(&req); err != nil {
		s.logger.Error().Err(err).Interface("request", req).Msg("Invalid create address request")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create domain entity
	address := &domain.Address{
		PlaceID:     req.PlaceID,
		FullAddress: req.FullAddress,
		Country:     req.Country,
		City:        req.City,
		District:    req.District,
		Street:      req.Street,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		DoorNumber:  req.DoorNumber,
	}

	// Create or update address
	updatedAddress, err := s.addressRepo.CreateOrUpdateByPlaceID(ctx, address)
	if err != nil {
		s.logger.Error().Err(err).Interface("address", address).Msg("Failed to create or update address")
		return nil, fmt.Errorf("failed to create or update address: %w", err)
	}

	s.logger.Info().Int("address_id", updatedAddress.ID).Str("place_id", updatedAddress.PlaceID).Msg("Address created or updated successfully")

	return s.addressToResponse(updatedAddress), nil
}

// Location-based operations
func (s *addressService) GetAddressesByCity(ctx context.Context, city string, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error) {
	addresses, paginationResp, err := s.addressRepo.SearchByCity(ctx, city, pagination)
	if err != nil {
		s.logger.Error().Err(err).Str("city", city).Msg("Failed to get addresses by city")
		return nil, nil, fmt.Errorf("failed to get addresses by city: %w", err)
	}

	var responses []*dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.addressToResponse(address))
	}

	return responses, paginationResp, nil
}

func (s *addressService) GetAddressesByCountry(ctx context.Context, country string, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error) {
	addresses, paginationResp, err := s.addressRepo.SearchByCountry(ctx, country, pagination)
	if err != nil {
		s.logger.Error().Err(err).Str("country", country).Msg("Failed to get addresses by country")
		return nil, nil, fmt.Errorf("failed to get addresses by country: %w", err)
	}

	var responses []*dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.addressToResponse(address))
	}

	return responses, paginationResp, nil
}

func (s *addressService) GetNearbyAddresses(ctx context.Context, latitude, longitude float64, radiusKm int, limit int) ([]*dto.AddressResponse, error) {
	// Validate coordinates
	if err := s.ValidateCoordinates(ctx, latitude, longitude); err != nil {
		return nil, err
	}

	addresses, err := s.addressRepo.GetNearbyAddresses(ctx, latitude, longitude, radiusKm, limit)
	if err != nil {
		s.logger.Error().Err(err).Float64("latitude", latitude).Float64("longitude", longitude).Int("radius", radiusKm).Msg("Failed to get nearby addresses")
		return nil, fmt.Errorf("failed to get nearby addresses: %w", err)
	}

	var responses []*dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.addressToResponse(address))
	}

	return responses, nil
}

func (s *addressService) GetAddressesByCoordinates(ctx context.Context, latitude, longitude float64, radiusKm int, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error) {
	// Validate coordinates
	if err := s.ValidateCoordinates(ctx, latitude, longitude); err != nil {
		return nil, nil, err
	}

	addresses, paginationResp, err := s.addressRepo.GetByCoordinates(ctx, latitude, longitude, radiusKm, pagination)
	if err != nil {
		s.logger.Error().Err(err).Float64("latitude", latitude).Float64("longitude", longitude).Int("radius", radiusKm).Msg("Failed to get addresses by coordinates")
		return nil, nil, fmt.Errorf("failed to get addresses by coordinates: %w", err)
	}

	var responses []*dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.addressToResponse(address))
	}

	return responses, paginationResp, nil
}

// Search and filtering
func (s *addressService) SearchAddresses(ctx context.Context, query string, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error) {
	addresses, paginationResp, err := s.addressRepo.SearchByFullAddress(ctx, query, pagination)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query).Msg("Failed to search addresses")
		return nil, nil, fmt.Errorf("failed to search addresses: %w", err)
	}

	var responses []*dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.addressToResponse(address))
	}

	return responses, paginationResp, nil
}

func (s *addressService) GetAddressesWithFilters(ctx context.Context, filters dto.AddressFilterRequest, pagination dto.PaginationRequest) ([]*dto.AddressResponse, *dto.PaginationResponse, error) {
	addresses, paginationResp, err := s.addressRepo.GetAddressesWithFilters(ctx, filters, pagination)
	if err != nil {
		s.logger.Error().Err(err).Interface("filters", filters).Msg("Failed to get addresses with filters")
		return nil, nil, fmt.Errorf("failed to get addresses with filters: %w", err)
	}

	var responses []*dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.addressToResponse(address))
	}

	return responses, paginationResp, nil
}

// Validation operations
func (s *addressService) ValidateAddress(ctx context.Context, address *domain.Address) error {
	// Validate coordinates
	if err := s.ValidateCoordinates(ctx, address.Latitude, address.Longitude); err != nil {
		return err
	}

	// Note: Place ID validation would be done with Google Places API
	// For now, we'll skip this validation

	// Check for duplicate addresses
	exists, err := s.addressRepo.ExistsByPlaceID(ctx, address.PlaceID)
	if err != nil {
		return fmt.Errorf("failed to check address existence: %w", err)
	}
	if exists && address.ID == 0 { // Only check for new addresses
		return fmt.Errorf("address with place ID %s already exists", address.PlaceID)
	}

	return nil
}

func (s *addressService) ValidateCoordinates(ctx context.Context, latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", longitude)
	}
	return nil
}

// Statistics operations
func (s *addressService) GetPopularCities(ctx context.Context, limit int) ([]string, error) {
	cities, err := s.addressRepo.GetPopularCities(ctx, limit)
	if err != nil {
		s.logger.Error().Err(err).Int("limit", limit).Msg("Failed to get popular cities")
		return nil, fmt.Errorf("failed to get popular cities: %w", err)
	}

	return cities, nil
}

func (s *addressService) GetPopularCountries(ctx context.Context, limit int) ([]string, error) {
	countries, err := s.addressRepo.GetPopularCountries(ctx, limit)
	if err != nil {
		s.logger.Error().Err(err).Int("limit", limit).Msg("Failed to get popular countries")
		return nil, fmt.Errorf("failed to get popular countries: %w", err)
	}

	return countries, nil
}

// Helper methods
func (s *addressService) validateCreateAddressRequest(req *dto.CreateAddressRequest) error {
	if req.PlaceID == "" {
		return fmt.Errorf("place ID is required")
	}
	if req.FullAddress == "" {
		return fmt.Errorf("full address is required")
	}
	if req.Country == "" {
		return fmt.Errorf("country is required")
	}
	if req.City == "" {
		return fmt.Errorf("city is required")
	}

	// Validate coordinates
	if err := s.ValidateCoordinates(context.Background(), req.Latitude, req.Longitude); err != nil {
		return err
	}

	return nil
}

func (s *addressService) addressToResponse(address *domain.Address) *dto.AddressResponse {
	return &dto.AddressResponse{
		ID:          address.ID,
		PlaceID:     address.PlaceID,
		FullAddress: address.FullAddress,
		Country:     address.Country,
		City:        address.City,
		District:    address.District,
		Street:      address.Street,
		PostalCode:  address.PostalCode,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
		DoorNumber:  address.DoorNumber,
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
	}
}
