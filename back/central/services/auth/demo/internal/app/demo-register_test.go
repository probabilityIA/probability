package app

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/mocks"
)

type entorno struct {
	repo  *mocks.RepositoryMock
	email *mocks.EmailSenderMock
	otp   *mocks.OTPPublisherMock
	cfg   *mocks.ConfigMock
	uc    IUseCase
}

func nuevoEntorno(t *testing.T) *entorno {
	t.Helper()
	e := &entorno{
		repo:  &mocks.RepositoryMock{},
		email: &mocks.EmailSenderMock{},
		otp:   &mocks.OTPPublisherMock{},
		cfg:   &mocks.ConfigMock{Values: map[string]string{}},
	}
	e.uc = New(e.repo, e.email, e.otp, mocks.NewSilentLogger(), e.cfg)
	return e
}

func registroValido() domain.DemoRegisterRequest {
	return domain.DemoRegisterRequest{
		FullName:     "Sebastian Camacho",
		BusinessName: "Tienda Probability",
		Email:        "demo@tienda.com",
		Password:     "secreta123",
	}
}

func TestDemoRegister_CreaLaCuentaConLosDatosNormalizados(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.FullName = "  Sebastian Camacho  "
	req.BusinessName = "  Tienda Probability  "
	req.Email = "  DEMO@Tienda.COM  "

	resp, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, resp.Success)
	require.Len(t, e.repo.CuentasCreadas, 1)
	creada := e.repo.CuentasCreadas[0]
	assert.Equal(t, "Sebastian Camacho", creada.FullName)
	assert.Equal(t, "Tienda Probability", creada.BusinessName)
	assert.Equal(t, "demo@tienda.com", creada.Email,
		"el correo se guarda en minusculas para que el login y el reenvio siempre encuentren la misma cuenta")
}

func TestDemoRegister_GuardaLaContrasenaHasheadaNoEnClaro(t *testing.T) {
	e := nuevoEntorno(t)

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	hash := e.repo.CuentasCreadas[0].PasswordHash
	assert.NotEqual(t, "secreta123", hash)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash), []byte("secreta123")),
		"el hash guardado debe validar la contrasena original con bcrypt")
}

func TestDemoRegister_ExigeNombreNegocioYCorreo(t *testing.T) {
	casos := map[string]func(*domain.DemoRegisterRequest){
		"sin nombre":           func(r *domain.DemoRegisterRequest) { r.FullName = "" },
		"sin negocio":          func(r *domain.DemoRegisterRequest) { r.BusinessName = "" },
		"sin correo":           func(r *domain.DemoRegisterRequest) { r.Email = "" },
		"nombre en blanco":     func(r *domain.DemoRegisterRequest) { r.FullName = "   " },
		"correo solo espacios": func(r *domain.DemoRegisterRequest) { r.Email = "  " },
	}

	for nombre, romper := range casos {
		t.Run(nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			req := registroValido()
			romper(&req)

			_, err := e.uc.DemoRegister(context.Background(), req)

			require.Error(t, err)
			assert.Empty(t, e.repo.CuentasCreadas)
		})
	}
}

func TestDemoRegister_ExigeContrasenaDeAlMenosSeisCaracteres(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Password = "12345"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "6 caracteres")
	assert.Empty(t, e.repo.CuentasCreadas)
}

func TestDemoRegister_UnaContrasenaDeSeisCaracteresYaEsAceptada(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Password = "123456"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err,
		"el minimo es 6 y no hay exigencia de mayusculas ni simbolos: es una cuenta demo, no produccion")
}

func TestDemoRegister_DistingueCuentaActivaDeCuentaPendiente(t *testing.T) {
	casos := []struct {
		nombre   string
		activa   bool
		esperado error
	}{
		{"cuenta ya activada", true, domain.ErrEmailAlreadyRegistered},
		{"cuenta pendiente de verificar", false, domain.ErrEmailPendingVerification},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			e.repo.GetDemoUserByEmailFn = func(ctx context.Context, email string) (*domain.PendingDemoUser, error) {
				return &domain.PendingDemoUser{UserID: 5, IsActive: caso.activa}, nil
			}

			_, err := e.uc.DemoRegister(context.Background(), registroValido())

			assert.ErrorIs(t, err, caso.esperado,
				"el mensaje distinto revela al front si el correo existe: es enumeracion de usuarios, aceptada aqui para poder ofrecer el reenvio")
			assert.Empty(t, e.repo.CuentasCreadas)
		})
	}
}

func TestDemoRegister_SiFallaLaConsultaDeCorreoNoFiltraElErrorInterno(t *testing.T) {
	e := nuevoEntorno(t)
	e.repo.GetDemoUserByEmailFn = func(ctx context.Context, email string) (*domain.PendingDemoUser, error) {
		return nil, errors.New("connection refused a postgres")
	}

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error(),
		"el detalle del fallo de base no se expone al cliente publico")
}

func TestDemoRegister_SinRolDemoConfiguradoNoCreaLaCuenta(t *testing.T) {
	casos := map[string]func() (uint, error){
		"el rol no existe":      func() (uint, error) { return 0, nil },
		"falla la consulta rol": func() (uint, error) { return 0, errors.New("timeout") },
	}

	for nombre, resultado := range casos {
		t.Run(nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			e.repo.GetDemoRoleIDFn = func(ctx context.Context) (uint, error) { return resultado() }

			_, err := e.uc.DemoRegister(context.Background(), registroValido())

			require.Error(t, err)
			assert.Empty(t, e.repo.CuentasCreadas,
				"sin rol demo la cuenta quedaria sin permisos: mejor no crearla")
		})
	}
}

func TestDemoRegister_AsignaElRolDemoQueDevuelveElRepositorio(t *testing.T) {
	e := nuevoEntorno(t)
	e.repo.GetDemoRoleIDFn = func(ctx context.Context) (uint, error) { return 42, nil }

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	assert.EqualValues(t, 42, e.repo.CuentasCreadas[0].RoleID)
}

func TestDemoRegister_GeneraUnCodigoDeNegocioDerivadoDelNombre(t *testing.T) {
	e := nuevoEntorno(t)

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	codigo := e.repo.CuentasCreadas[0].BusinessCode
	assert.True(t, strings.HasPrefix(codigo, "tienda-probability-"),
		"el codigo arranca con el slug del negocio para que sea reconocible en la base, got %q", codigo)
	assert.Len(t, codigo, len("tienda-probability-")+6,
		"el sufijo son 3 bytes en hex, o sea 6 caracteres")
}

func TestDemoRegister_ReintentaHastaEncontrarUnCodigoLibre(t *testing.T) {
	e := nuevoEntorno(t)
	intentos := 0
	e.repo.BusinessCodeExistsFn = func(ctx context.Context, code string) (bool, error) {
		intentos++
		return intentos < 3, nil
	}

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	assert.Equal(t, 3, intentos)
	assert.Len(t, e.repo.CodigosVerificados, 3,
		"cada intento genera un sufijo aleatorio nuevo hasta que uno no colisione")
}

func TestDemoRegister_SeRindeDespuesDeSeisColisionesDeCodigo(t *testing.T) {
	e := nuevoEntorno(t)
	e.repo.BusinessCodeExistsFn = func(ctx context.Context, code string) (bool, error) { return true, nil }

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.Error(t, err)
	assert.Len(t, e.repo.CodigosVerificados, 6,
		"tope de 6 intentos: evita un bucle infinito si la tabla quedo en un estado raro")
	assert.Empty(t, e.repo.CuentasCreadas)
}

func TestDemoRegister_SiFallaLaConsultaDeCodigoAbortaElRegistro(t *testing.T) {
	e := nuevoEntorno(t)
	e.repo.BusinessCodeExistsFn = func(ctx context.Context, code string) (bool, error) {
		return false, errors.New("timeout")
	}

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.Error(t, err)
	assert.Empty(t, e.repo.CuentasCreadas)
}

func TestDemoRegister_UnNombreDeNegocioSinLetrasUsaElCodigoBaseDemo(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.BusinessName = "123 456"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(e.repo.CuentasCreadas[0].BusinessCode, "123-456-"),
		"los digitos si cuentan para el slug")
}

func TestDemoRegister_UnNombreSinCaracteresValidosCaeAlSlugDemo(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.BusinessName = "***"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(e.repo.CuentasCreadas[0].BusinessCode, "demo-"))
}

func TestDemoRegister_ElCodigoDeNegocioNoPasaDe50Caracteres(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.BusinessName = strings.Repeat("nombrelargo", 10)

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.Len(t, e.repo.CuentasCreadas[0].BusinessCode, 50,
		"la columna business_code esta limitada a 50 caracteres")
}

func TestDemoRegister_DerivaElPrefijoDeOrdenesConTresLetras(t *testing.T) {
	casos := []struct {
		negocio  string
		esperado string
	}{
		{"Tienda Probability", "TIE"},
		{"AB", "ABX"},
		{"A", "AXX"},
		{"12345", "DEM"},
		{"a b c", "ABC"},
	}

	for _, caso := range casos {
		t.Run(caso.negocio, func(t *testing.T) {
			e := nuevoEntorno(t)
			req := registroValido()
			req.BusinessName = caso.negocio

			_, err := e.uc.DemoRegister(context.Background(), req)

			require.NoError(t, err)
			assert.Equal(t, caso.esperado, e.repo.CuentasCreadas[0].OrderPrefix,
				"el prefijo va en el numero de orden, por eso siempre son 3 letras rellenadas con X")
		})
	}
}

func TestDemoRegister_PorCorreoMandaElEnlaceDeVerificacionYNoUsaWhatsApp(t *testing.T) {
	e := nuevoEntorno(t)
	e.cfg.Values["FRONTEND_BASE_URL"] = "https://app.probability.com/"

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	require.Len(t, e.email.Enviados, 1)
	assert.Equal(t, "demo@tienda.com", e.email.Enviados[0].To)
	assert.Contains(t, e.email.Enviados[0].HTML, "https://app.probability.com/verify-email?token=",
		"la barra final de la variable de entorno se recorta para no generar una URL con doble slash")
	assert.Empty(t, e.otp.Publicados)
}

func TestDemoRegister_SinFrontendBaseURLUsaLocalhost(t *testing.T) {
	e := nuevoEntorno(t)

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	assert.Contains(t, e.email.Enviados[0].HTML, "http://localhost:3000/verify-email?token=")
}

func TestDemoRegister_ElTokenViajaEnClaroPorCorreoPeroSeGuardaHasheado(t *testing.T) {
	e := nuevoEntorno(t)

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	guardado := e.repo.CuentasCreadas[0].TokenHash
	html := e.email.Enviados[0].HTML
	assert.NotContains(t, html, guardado,
		"si se filtra la base, el hash no sirve para verificar cuentas ajenas")
	assert.Len(t, guardado, 64, "sha256 en hex son 64 caracteres")
}

func TestDemoRegister_ElTokenDeCorreoCaducaEn24Horas(t *testing.T) {
	e := nuevoEntorno(t)
	antes := time.Now()

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	expira := e.repo.CuentasCreadas[0].ExpiresAt
	assert.WithinDuration(t, antes.Add(24*time.Hour), expira, time.Minute)
}

func TestDemoRegister_PorWhatsAppPublicaUnOTPDeSeisDigitosYNoMandaCorreo(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Channel = "WhatsApp"
	req.Phone = "  573001112233  "

	resp, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.Contains(t, resp.Message, "WhatsApp")
	require.Len(t, e.otp.Publicados, 1)
	evento := e.otp.Publicados[0]
	assert.Len(t, evento.Code, 6)
	assert.Regexp(t, `^\d{6}$`, evento.Code)
	assert.Equal(t, "573001112233", evento.Phone, "el telefono se recorta antes de mandarlo a Meta")
	assert.Equal(t, 15, evento.ExpiresMinutes)
	assert.Empty(t, e.email.Enviados)
}

func TestDemoRegister_ElOTPSeGuardaHasheadoJuntoConElCorreo(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Channel = "whatsapp"
	req.Phone = "573001112233"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	esperado := hashOTP("demo@tienda.com", e.otp.Publicados[0].Code)
	assert.Equal(t, esperado, e.repo.CuentasCreadas[0].TokenHash,
		"el hash mezcla correo y codigo, asi un OTP de 6 digitos no vale para la cuenta de otro")
}

func TestDemoRegister_ElOTPCaducaEn15Minutos(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Channel = "whatsapp"
	req.Phone = "573001112233"
	antes := time.Now()

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.WithinDuration(t, antes.Add(15*time.Minute), e.repo.CuentasCreadas[0].ExpiresAt, time.Minute,
		"la ventana del OTP es mucho mas corta que la del enlace por ser solo 6 digitos")
}

func TestDemoRegister_WhatsAppSinTelefonoSeRechazaAntesDeCrearLaCuenta(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Channel = "whatsapp"
	req.Phone = "   "

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "telefono")
	assert.Empty(t, e.repo.CuentasCreadas)
}

func TestDemoRegister_UnCanalDesconocidoCaeAlFlujoDeCorreo(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Channel = "telegram"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.Len(t, e.email.Enviados, 1,
		"cualquier canal que no sea whatsapp se trata como correo, no se rechaza")
	assert.Empty(t, e.otp.Publicados)
}

func TestDemoRegister_SiFallaElInsertNoSeNotificaANadie(t *testing.T) {
	e := nuevoEntorno(t)
	e.repo.CreateDemoAccountFn = func(ctx context.Context, params domain.CreateDemoAccountParams) (uint, error) {
		return 0, errors.New("unique violation")
	}

	_, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.Error(t, err)
	assert.Empty(t, e.email.Enviados)
	assert.Empty(t, e.otp.Publicados)
}

func TestDemoRegister_SiFallaElCorreoLaCuentaIgualQuedaCreadaYSeReportaExito(t *testing.T) {
	e := nuevoEntorno(t)
	e.email.SendHTMLFn = func(ctx context.Context, to, subject, html string) error {
		return errors.New("SMTP 550")
	}

	resp, err := e.uc.DemoRegister(context.Background(), registroValido())

	require.NoError(t, err)
	assert.True(t, resp.Success,
		"la cuenta quedo creada: se responde exito y el usuario debe usar el reenvio, aunque nunca reciba el primer correo")
	assert.Len(t, e.repo.CuentasCreadas, 1)
}

func TestDemoRegister_SiFallaLaPublicacionDelOTPTambienSeReportaExito(t *testing.T) {
	e := nuevoEntorno(t)
	e.otp.PublishDemoOTPFn = func(ctx context.Context, event domain.DemoOTPEvent) error {
		return errors.New("rabbitmq caido")
	}
	req := registroValido()
	req.Channel = "whatsapp"
	req.Phone = "573001112233"

	resp, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, resp.Success,
		"queda una cuenta pendiente cuyo codigo nunca se envio: solo el reenvio la rescata")
}

func TestDemoRegister_GuardaElTelefonoAunEnElFlujoDeCorreo(t *testing.T) {
	e := nuevoEntorno(t)
	req := registroValido()
	req.Phone = "573001112233"

	_, err := e.uc.DemoRegister(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "573001112233", e.repo.CuentasCreadas[0].Phone)
}

func TestSlugify_ColapsaSeparadoresYDescartaAcentosYSimbolos(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{"Tienda Probability", "tienda-probability"},
		{"Tienda  Probability", "tienda-probability"},
		{"Tienda_-_Probability", "tienda-probability"},
		{"  Tienda  ", "tienda"},
		{"Café Bar", "caf-bar"},
		{"***", ""},
		{"", ""},
		{"---", ""},
	}

	for _, caso := range casos {
		assert.Equal(t, caso.esperado, slugify(caso.entrada),
			"los acentos se descartan (no se transliteran), por eso Cafe queda como caf")
	}
}

func TestGenerateToken_DevuelveUnTokenAleatorioConSuHashSha256(t *testing.T) {
	raw1, hash1, err := generateToken()
	require.NoError(t, err)
	raw2, _, err := generateToken()
	require.NoError(t, err)

	assert.Len(t, raw1, 64, "32 bytes en hex")
	assert.NotEqual(t, raw1, raw2, "dos registros seguidos no pueden compartir token")
	assert.Equal(t, hash1, hashToken(raw1),
		"el hash guardado debe reproducirse desde el token que llega por la URL de verificacion")
}

func TestHashOTP_NormalizaElCorreoAntesDeMezclarlo(t *testing.T) {
	assert.Equal(t, hashOTP("demo@tienda.com", "123456"), hashOTP("  DEMO@Tienda.COM  ", "123456"),
		"el correo llega distinto desde el registro y desde el formulario de verificacion")
	assert.NotEqual(t, hashOTP("demo@tienda.com", "123456"), hashOTP("otro@tienda.com", "123456"),
		"el mismo codigo no puede servir para dos cuentas")
}
