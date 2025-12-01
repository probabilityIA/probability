package domain

import (
	"context"
	"io"
	"mime/multipart"
)

type IUserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*UserAuthInfo, error)
	GetUserByID(ctx context.Context, userID uint) (*UserAuthInfo, error)
	GetUserRoles(ctx context.Context, userID uint) ([]Role, error)
	GetUserBusinesses(ctx context.Context, userID uint) ([]BusinessInfoEntity, error)
	GetUserRoleByBusiness(ctx context.Context, userID uint, businessID uint) (*Role, error)
	GetRoleByID(ctx context.Context, id uint) (*Role, error)
	GetUsers(ctx context.Context, filters UserFilters) ([]UserQueryDTO, int64, error)
	CreateUser(ctx context.Context, user UsersEntity) (uint, error)
	UpdateUser(ctx context.Context, id uint, user UsersEntity) (string, error)
	DeleteUser(ctx context.Context, id uint) (string, error)
	AssignBusinessStaffRelationships(ctx context.Context, userID uint, assignments []BusinessRoleAssignment) error
	GetBusinessStaffRelationships(ctx context.Context, userID uint) ([]BusinessRoleAssignmentDetailed, error)
	AssignRoleToUserBusiness(ctx context.Context, userID uint, assignments []BusinessRoleAssignment) error // Asigna/actualiza roles a usuario en múltiples businesses
	AssignBusinessesToUser(ctx context.Context, userID uint, businessIDs []uint) error                     // Deprecated: usar AssignBusinessStaffRelationships
}

type IS3Service interface {
	GetImageURL(filename string) string
	DeleteImage(ctx context.Context, filename string) error
	ImageExists(ctx context.Context, filename string) (bool, error)
	UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error)
	DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error)
	FileExists(ctx context.Context, filename string) (bool, error)
	GetFileURL(ctx context.Context, filename string) (string, error)
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
}
