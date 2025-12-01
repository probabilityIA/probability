package domain

import "context"

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*UserAuthInfo, error)
	GetUserByID(ctx context.Context, userID uint) (*UserAuthInfo, error)
	GetUserRoles(ctx context.Context, userID uint) ([]Role, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]Permission, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
	ChangePassword(ctx context.Context, userID uint, newPassword string) error
	GetUserBusinesses(ctx context.Context, userID uint) ([]BusinessInfoEntity, error)
	GetUserRoleByBusiness(ctx context.Context, userID uint, businessID uint) (*Role, error)
	GetBusinessStaffRelation(ctx context.Context, userID uint, businessID *uint) (*BusinessStaffRelation, error)
	GetBusinessConfiguredResourcesIDs(ctx context.Context, businessID uint) ([]uint, error)
	GetBusinessByID(ctx context.Context, businessID uint) (*BusinessInfo, error)
	GetRoleByID(ctx context.Context, id uint) (*Role, error)
}
type IJWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string) (string, error)

	// Tokens para business
	GenerateBusinessToken(userID, businessID, businessTypeID, roleID uint) (string, error)
	ValidateBusinessToken(tokenString string) (*BusinessTokenClaims, error)
}
