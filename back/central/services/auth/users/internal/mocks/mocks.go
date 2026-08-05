package mocks

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UserRepositoryMock struct {
	GetUserByEmailFn                   func(ctx context.Context, email string) (*domain.UserAuthInfo, error)
	GetUserByIDFn                      func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error)
	GetUserRolesFn                     func(ctx context.Context, userID uint) ([]domain.Role, error)
	GetUserBusinessesFn                func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error)
	GetUserRoleByBusinessFn            func(ctx context.Context, userID uint, businessID uint) (*domain.Role, error)
	GetRoleByIDFn                      func(ctx context.Context, id uint) (*domain.Role, error)
	GetUsersFn                         func(ctx context.Context, filters domain.UserFilters) ([]domain.UserQueryDTO, int64, error)
	CreateUserFn                       func(ctx context.Context, user domain.UsersEntity) (uint, error)
	UpdateUserFn                       func(ctx context.Context, id uint, user domain.UsersEntity) (string, error)
	DeleteUserFn                       func(ctx context.Context, id uint) (string, error)
	AssignBusinessStaffRelationshipsFn func(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error
	GetBusinessStaffRelationshipsFn    func(ctx context.Context, userID uint) ([]domain.BusinessRoleAssignmentDetailed, error)
	AssignRoleToUserBusinessFn         func(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error
	AssignBusinessesToUserFn           func(ctx context.Context, userID uint, businessIDs []uint) error
	AssignRolesToUserFn                func(ctx context.Context, userID uint, roleIDs []uint) error
	GetRoleIDByNameAndScopeFn          func(ctx context.Context, roleName string, scopeID uint) (uint, error)
}

var _ domain.IUserRepository = (*UserRepositoryMock)(nil)

func (m *UserRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *UserRepositoryMock) GetUserByID(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
	if m.GetUserByIDFn != nil {
		return m.GetUserByIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *UserRepositoryMock) GetUserRoles(ctx context.Context, userID uint) ([]domain.Role, error) {
	if m.GetUserRolesFn != nil {
		return m.GetUserRolesFn(ctx, userID)
	}
	return nil, nil
}

func (m *UserRepositoryMock) GetUserBusinesses(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
	if m.GetUserBusinessesFn != nil {
		return m.GetUserBusinessesFn(ctx, userID)
	}
	return nil, nil
}

func (m *UserRepositoryMock) GetUserRoleByBusiness(ctx context.Context, userID uint, businessID uint) (*domain.Role, error) {
	if m.GetUserRoleByBusinessFn != nil {
		return m.GetUserRoleByBusinessFn(ctx, userID, businessID)
	}
	return nil, nil
}

func (m *UserRepositoryMock) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	if m.GetRoleByIDFn != nil {
		return m.GetRoleByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *UserRepositoryMock) GetUsers(ctx context.Context, filters domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
	if m.GetUsersFn != nil {
		return m.GetUsersFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user domain.UsersEntity) (uint, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return 0, nil
}

func (m *UserRepositoryMock) UpdateUser(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, id, user)
	}
	return "", nil
}

func (m *UserRepositoryMock) DeleteUser(ctx context.Context, id uint) (string, error) {
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, id)
	}
	return "", nil
}

func (m *UserRepositoryMock) AssignBusinessStaffRelationships(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error {
	if m.AssignBusinessStaffRelationshipsFn != nil {
		return m.AssignBusinessStaffRelationshipsFn(ctx, userID, assignments)
	}
	return nil
}

func (m *UserRepositoryMock) GetBusinessStaffRelationships(ctx context.Context, userID uint) ([]domain.BusinessRoleAssignmentDetailed, error) {
	if m.GetBusinessStaffRelationshipsFn != nil {
		return m.GetBusinessStaffRelationshipsFn(ctx, userID)
	}
	return nil, nil
}

func (m *UserRepositoryMock) AssignRoleToUserBusiness(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error {
	if m.AssignRoleToUserBusinessFn != nil {
		return m.AssignRoleToUserBusinessFn(ctx, userID, assignments)
	}
	return nil
}

func (m *UserRepositoryMock) AssignBusinessesToUser(ctx context.Context, userID uint, businessIDs []uint) error {
	if m.AssignBusinessesToUserFn != nil {
		return m.AssignBusinessesToUserFn(ctx, userID, businessIDs)
	}
	return nil
}

func (m *UserRepositoryMock) AssignRolesToUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if m.AssignRolesToUserFn != nil {
		return m.AssignRolesToUserFn(ctx, userID, roleIDs)
	}
	return nil
}

func (m *UserRepositoryMock) GetRoleIDByNameAndScope(ctx context.Context, roleName string, scopeID uint) (uint, error) {
	if m.GetRoleIDByNameAndScopeFn != nil {
		return m.GetRoleIDByNameAndScopeFn(ctx, roleName, scopeID)
	}
	return 0, nil
}

type S3ServiceMock struct {
	GetImageURLFn   func(filename string) string
	DeleteImageFn   func(ctx context.Context, filename string) error
	ImageExistsFn   func(ctx context.Context, filename string) (bool, error)
	UploadFileFn    func(ctx context.Context, file io.ReadSeeker, filename string) (string, error)
	DownloadFileFn  func(ctx context.Context, filename string) (io.ReadSeeker, error)
	FileExistsFn    func(ctx context.Context, filename string) (bool, error)
	GetFileURLFn    func(ctx context.Context, filename string) (string, error)
	UploadImageFn   func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	DeleteImageCall []string
}

var _ domain.IS3Service = (*S3ServiceMock)(nil)

func (m *S3ServiceMock) GetImageURL(filename string) string {
	if m.GetImageURLFn != nil {
		return m.GetImageURLFn(filename)
	}
	return ""
}

func (m *S3ServiceMock) DeleteImage(ctx context.Context, filename string) error {
	m.DeleteImageCall = append(m.DeleteImageCall, filename)
	if m.DeleteImageFn != nil {
		return m.DeleteImageFn(ctx, filename)
	}
	return nil
}

func (m *S3ServiceMock) ImageExists(ctx context.Context, filename string) (bool, error) {
	if m.ImageExistsFn != nil {
		return m.ImageExistsFn(ctx, filename)
	}
	return false, nil
}

func (m *S3ServiceMock) UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error) {
	if m.UploadFileFn != nil {
		return m.UploadFileFn(ctx, file, filename)
	}
	return "", nil
}

func (m *S3ServiceMock) DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error) {
	if m.DownloadFileFn != nil {
		return m.DownloadFileFn(ctx, filename)
	}
	return nil, nil
}

func (m *S3ServiceMock) FileExists(ctx context.Context, filename string) (bool, error) {
	if m.FileExistsFn != nil {
		return m.FileExistsFn(ctx, filename)
	}
	return false, nil
}

func (m *S3ServiceMock) GetFileURL(ctx context.Context, filename string) (string, error) {
	if m.GetFileURLFn != nil {
		return m.GetFileURLFn(ctx, filename)
	}
	return "", nil
}

func (m *S3ServiceMock) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	if m.UploadImageFn != nil {
		return m.UploadImageFn(ctx, file, folder)
	}
	return "", nil
}

type ConfigMock struct {
	Values map[string]string
}

var _ env.IConfig = (*ConfigMock)(nil)

func NewConfigMock(values map[string]string) *ConfigMock {
	if values == nil {
		values = map[string]string{}
	}
	return &ConfigMock{Values: values}
}

func (m *ConfigMock) Get(key string) string {
	return m.Values[key]
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger {
	return &SilentLogger{}
}

func (l *SilentLogger) nop() zerolog.Logger {
	return zerolog.Nop()
}

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Info()
}

func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Error()
}

func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Warn()
}

func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Debug()
}

func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Fatal()
}

func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event {
	n := l.nop()
	return n.Panic()
}

func (l *SilentLogger) With() zerolog.Context {
	n := l.nop()
	return n.With()
}

func (l *SilentLogger) WithService(service string) log.ILogger {
	return l
}

func (l *SilentLogger) WithModule(module string) log.ILogger {
	return l
}

func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger {
	return l
}
