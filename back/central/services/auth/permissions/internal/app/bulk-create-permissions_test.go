package app

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBulkCreatePermissions_AllSuccess(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockLogger := new(mocks.LoggerMock)

	useCase := New(mockRepo, mockLogger)

	ctx := context.Background()
	businessTypeID := uint(1)

	dtos := []domain.CreatePermissionDTO{
		{
			Name:           "Crear Usuarios",
			Code:           "users.create",
			Description:    "Permite crear usuarios",
			ResourceID:     1,
			ActionID:       1,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
		{
			Name:           "Editar Usuarios",
			Code:           "users.edit",
			Description:    "Permite editar usuarios",
			ResourceID:     1,
			ActionID:       2,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
	}

	// Configurar mocks: ambos permisos no existen y se crean exitosamente
	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything).Maybe()

	mockRepo.On("PermissionExistsByName", ctx, "Crear Usuarios").Return(false, nil)
	mockRepo.On("PermissionExistsByName", ctx, "Editar Usuarios").Return(false, nil)

	mockRepo.On("CreatePermission", ctx, mock.AnythingOfType("domain.Permission")).
		Return("permiso creado exitosamente", nil).Times(2)

	// Act
	results, err := useCase.BulkCreatePermissions(ctx, dtos)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.True(t, results[0].Success)
	assert.Equal(t, "Crear Usuarios", results[0].Name)
	assert.Empty(t, results[0].Error)

	assert.True(t, results[1].Success)
	assert.Equal(t, "Editar Usuarios", results[1].Name)
	assert.Empty(t, results[1].Error)

	mockRepo.AssertExpectations(t)
}

func TestBulkCreatePermissions_PartialFailure(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockLogger := new(mocks.LoggerMock)

	useCase := New(mockRepo, mockLogger)

	ctx := context.Background()
	businessTypeID := uint(1)

	dtos := []domain.CreatePermissionDTO{
		{
			Name:           "Crear Usuarios",
			Code:           "users.create",
			ResourceID:     1,
			ActionID:       1,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
		{
			Name:           "Permiso Duplicado",
			Code:           "users.duplicate",
			ResourceID:     1,
			ActionID:       2,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
	}

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything).Maybe()

	// El primero se crea exitosamente
	mockRepo.On("PermissionExistsByName", ctx, "Crear Usuarios").Return(false, nil)
	mockRepo.On("CreatePermission", ctx, mock.AnythingOfType("domain.Permission")).
		Return("permiso creado exitosamente", nil).Once()

	// El segundo falla porque el nombre ya existe
	mockRepo.On("PermissionExistsByName", ctx, "Permiso Duplicado").Return(true, nil)

	// Act
	results, err := useCase.BulkCreatePermissions(ctx, dtos)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.True(t, results[0].Success)
	assert.Equal(t, "Crear Usuarios", results[0].Name)
	assert.Empty(t, results[0].Error)

	assert.False(t, results[1].Success)
	assert.Equal(t, "Permiso Duplicado", results[1].Name)
	assert.Contains(t, results[1].Error, "Permiso Duplicado")

	mockRepo.AssertExpectations(t)
}

func TestBulkCreatePermissions_AllFail(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockLogger := new(mocks.LoggerMock)

	useCase := New(mockRepo, mockLogger)

	ctx := context.Background()
	businessTypeID := uint(1)

	dtos := []domain.CreatePermissionDTO{
		{
			Name:           "Permiso A",
			Code:           "perm.a",
			ResourceID:     1,
			ActionID:       1,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
		{
			Name:           "Permiso B",
			Code:           "perm.b",
			ResourceID:     1,
			ActionID:       2,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
	}

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything).Maybe()

	// Ambos nombres ya existen
	mockRepo.On("PermissionExistsByName", ctx, "Permiso A").Return(true, nil)
	mockRepo.On("PermissionExistsByName", ctx, "Permiso B").Return(true, nil)

	// Act
	results, err := useCase.BulkCreatePermissions(ctx, dtos)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.False(t, results[0].Success)
	assert.Equal(t, "Permiso A", results[0].Name)
	assert.NotEmpty(t, results[0].Error)

	assert.False(t, results[1].Success)
	assert.Equal(t, "Permiso B", results[1].Name)
	assert.NotEmpty(t, results[1].Error)

	// Verificar que CreatePermission nunca fue llamado
	mockRepo.AssertNotCalled(t, "CreatePermission", mock.Anything, mock.Anything)
}

func TestBulkCreatePermissions_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockLogger := new(mocks.LoggerMock)

	useCase := New(mockRepo, mockLogger)

	ctx := context.Background()

	dtos := []domain.CreatePermissionDTO{}

	// Act
	results, err := useCase.BulkCreatePermissions(ctx, dtos)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 0)

	// Verificar que no se llamó al repositorio
	mockRepo.AssertNotCalled(t, "PermissionExistsByName", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "CreatePermission", mock.Anything, mock.Anything)
}

func TestBulkCreatePermissions_ValidationError(t *testing.T) {
	// Arrange: un item sin nombre falla validacion antes de llegar al repositorio
	mockRepo := new(mocks.RepositoryMock)
	mockLogger := new(mocks.LoggerMock)

	useCase := New(mockRepo, mockLogger)

	ctx := context.Background()
	businessTypeID := uint(1)

	dtos := []domain.CreatePermissionDTO{
		{
			// Name vacio: debe fallar la validacion
			Name:           "",
			Code:           "perm.noname",
			ResourceID:     1,
			ActionID:       1,
			ScopeID:        1,
			BusinessTypeID: &businessTypeID,
		},
	}

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()

	// Act
	results, err := useCase.BulkCreatePermissions(ctx, dtos)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	assert.False(t, results[0].Success)
	assert.Equal(t, "", results[0].Name)
	assert.Contains(t, results[0].Error, "nombre")

	// El repositorio no debe ser llamado porque la validacion falla primero
	mockRepo.AssertNotCalled(t, "PermissionExistsByName", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "CreatePermission", mock.Anything, mock.Anything)
}
