package domain

import (
	"context"
	"mime/multipart"
)

type IBusinessRepository interface {
	GetBusinesses(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]Business, int64, error)
	GetBusinessByID(ctx context.Context, id uint) (*Business, error)
	GetBusinessByCode(ctx context.Context, code string) (*Business, error)
	GetBusinessByCustomDomain(ctx context.Context, domain string) (*Business, error)
	CreateBusiness(ctx context.Context, business Business) (uint, error)
	UpdateBusiness(ctx context.Context, id uint, business Business) (string, error)
	DeleteBusiness(ctx context.Context, id uint) (string, error)
	GetBusinessTypes(ctx context.Context) ([]BusinessType, error)
	GetBusinessTypeByID(ctx context.Context, id uint) (*BusinessType, error)
	GetBusinessTypeByCode(ctx context.Context, code string) (*BusinessType, error)
	GetBusinessTypeByName(ctx context.Context, name string) (*BusinessType, error)
	CreateBusinessType(ctx context.Context, businessType BusinessType) (string, error)
	UpdateBusinessType(ctx context.Context, id uint, businessType BusinessType) (string, error)
	DeleteBusinessType(ctx context.Context, id uint) (string, error)

	// Métodos para BusinessTypeResourcePermitted
	GetBusinessTypeResourcesPermitted(ctx context.Context, businessTypeID uint) ([]BusinessTypeResourcePermitted, error)
	GetResourceByID(ctx context.Context, resourceID uint) (*Resource, error)

	// Métodos para BusinessResourceConfigured (recursos asignados a un business)
	GetBusinessesWithConfiguredResourcesPaginated(ctx context.Context, page, perPage int, businessID *uint, businessTypeID *uint) ([]BusinessWithConfiguredResourcesResponse, int64, error)
	GetBusinessByIDWithConfiguredResources(ctx context.Context, businessID uint) (*BusinessWithConfiguredResourcesResponse, error)
	ToggleBusinessResourceActive(ctx context.Context, businessID uint, resourceID uint, active bool) error
}

// IS3Service define las operaciones de almacenamiento en S3
type IS3Service interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) // Retorna path relativo
	GetImageURL(filename string) string                                                         // Genera URL completa
	DeleteImage(ctx context.Context, filename string) error
	ImageExists(ctx context.Context, filename string) (bool, error)
}
