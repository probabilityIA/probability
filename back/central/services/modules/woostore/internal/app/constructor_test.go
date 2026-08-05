package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/woostore/internal/domain"
)

type powerManagerMock struct {
	StartFn  func(ctx context.Context) (*domain.PowerState, error)
	StopFn   func(ctx context.Context) (*domain.PowerState, error)
	StatusFn func(ctx context.Context) (*domain.PowerState, error)

	Llamadas []string
}

var _ domain.IPowerManager = (*powerManagerMock)(nil)

func (m *powerManagerMock) Start(ctx context.Context) (*domain.PowerState, error) {
	m.Llamadas = append(m.Llamadas, "start")
	if m.StartFn != nil {
		return m.StartFn(ctx)
	}
	return &domain.PowerState{InstanceID: "i-123", State: "pending"}, nil
}

func (m *powerManagerMock) Stop(ctx context.Context) (*domain.PowerState, error) {
	m.Llamadas = append(m.Llamadas, "stop")
	if m.StopFn != nil {
		return m.StopFn(ctx)
	}
	return &domain.PowerState{InstanceID: "i-123", State: "stopping"}, nil
}

func (m *powerManagerMock) Status(ctx context.Context) (*domain.PowerState, error) {
	m.Llamadas = append(m.Llamadas, "status")
	if m.StatusFn != nil {
		return m.StatusFn(ctx)
	}
	return &domain.PowerState{InstanceID: "i-123", State: "running", PublicIP: "54.1.2.3"}, nil
}

func TestStart_DelegaEnElAdministradorDeEnergia(t *testing.T) {
	mgr := &powerManagerMock{}

	estado, err := New(mgr).Start(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []string{"start"}, mgr.Llamadas)
	assert.Equal(t, "pending", estado.State)
}

func TestStop_DelegaEnElAdministradorDeEnergia(t *testing.T) {
	mgr := &powerManagerMock{}

	estado, err := New(mgr).Stop(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []string{"stop"}, mgr.Llamadas)
	assert.Equal(t, "stopping", estado.State)
}

func TestStatus_DevuelveElEstadoYLaIPPublicaDeLaInstancia(t *testing.T) {
	mgr := &powerManagerMock{}

	estado, err := New(mgr).Status(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []string{"status"}, mgr.Llamadas)
	assert.Equal(t, "running", estado.State)
	assert.Equal(t, "54.1.2.3", estado.PublicIP,
		"la IP publica cambia en cada arranque de la EC2, por eso se devuelve junto con el estado")
}

func TestPower_NoHayGuardaDeNegocioAntesDeApagarLaTiendaDePruebas(t *testing.T) {
	mgr := &powerManagerMock{
		StatusFn: func(ctx context.Context) (*domain.PowerState, error) {
			return &domain.PowerState{State: "running"}, nil
		},
	}
	uc := New(mgr)

	_, err := uc.Stop(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []string{"stop"}, mgr.Llamadas,
		"el caso de uso no consulta el estado antes de apagar: apagar una instancia ya apagada es un no-op en EC2, pero tampoco hay validacion de permisos aqui")
}

func TestStart_PropagaElErrorDeAWS(t *testing.T) {
	fallo := errors.New("UnauthorizedOperation")
	mgr := &powerManagerMock{
		StartFn: func(ctx context.Context) (*domain.PowerState, error) { return nil, fallo },
	}

	estado, err := New(mgr).Start(context.Background())

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, estado)
}

func TestStop_PropagaElErrorDeAWS(t *testing.T) {
	fallo := errors.New("IncorrectInstanceState")
	mgr := &powerManagerMock{
		StopFn: func(ctx context.Context) (*domain.PowerState, error) { return nil, fallo },
	}

	estado, err := New(mgr).Stop(context.Background())

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, estado)
}

func TestStatus_PropagaElErrorDeAWS(t *testing.T) {
	fallo := errors.New("InvalidInstanceID.NotFound")
	mgr := &powerManagerMock{
		StatusFn: func(ctx context.Context) (*domain.PowerState, error) { return nil, fallo },
	}

	estado, err := New(mgr).Status(context.Background())

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, estado)
}
