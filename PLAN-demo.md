# Plan: Autoregistro Demo (rama feat/demo-autoregistro, worktree)

Objetivo: boton "Crea tu demo" en el login -> modal flotante (nombre, negocio, email,
contrasena) -> crea Business + User + rol "demo" + verificacion de correo. Al verificar,
el usuario entra y la plataforma ya tiene integraciones en MODO TEST (ordenes Probability,
facturacion, guias) apuntando a los MOCKS locales (back/testing). Solo ve modulos:
Ordenes, Facturacion, Envios(guias), Inventario. No puede crear usuarios.

Correo de pruebas: secamc93@gmail.com (puedo leer su inbox via Gmail MCP para sacar el link).

## Datos reales de DB (RDS prod) ya verificados
- scope: platform=1, business=2, operator=3
- action: Create=1, Read=2, Update=3, Delete=4, Manage=5
- business_type: 1=Probability (unico)
- resource: Usuarios=1, Ordenes=6, Productos=7, Integraciones=8, Envios=9, Facturacion=10, Clientes=19, Inventario=22, Bodegas=23
- role existentes: SuperAdmin=1, Operador=2, Administrador=4(scope2), cliente_final=5(scope2), OperatorAdmin=6
- permisos scope=2 existentes a reusar para demo:
  - Ordenes: Ver=1, Editar=2  (NO existe "Crear Ordenes" scope2 -> crear en seed por si la UI lo exige)
  - Productos: Ver=5
  - Envios: Crear=7, Ver=6, Actualizar=9, Eliminar=8
  - Facturacion: Create=18, Read=20, Update=21, Delete=19
  - Inventario: Create=66, Read=67, Update=68, Delete=69
- integration_types (base_url || base_url_test):
  - 5 softpymes: prod || http://back-testing:9090
  - 6 platform: vacio || vacio
  - 8 siigo: prod || http://back-testing:9095
  - 12 envioclick: prod || http://back-testing:9091/api/v2
- OJO: base_url_test usa hostname Docker `back-testing`. Para correr LOCAL (go run) hay que
  mapear back-testing -> 127.0.0.1 en /etc/hosts, o el backend no resuelve los mocks.

## Backend de testing / mocks (back/testing)
- Levantar: `cd back/testing && go run cmd/main.go` (o dev-services start testing). Puerto API 9092.
- Spawnea mocks: softpymes :9090, envioclick :9091, shopify :9093, bold :9094, siigo :9095.
- Aceptan credenciales dummy. EnvioClick: POST /api/v2/quotation, /api/v2/shipment (devuelve tracker EC-xx"), /api/v2/track. Siigo/Softpymes: crear factura.

## Rol demo: permisos a asignar (IDs existentes)
[1(Ordenes Ver), 2(Ordenes Editar), 5(Productos Ver), 67/66/68/69(Inventario CRUD), 6(Envios Ver), 20(Facturacion Ver)]
+ crear permiso "Crear Ordenes" (resource6, action1, scope2) e incluirlo.
Resultado en sidebar: muestra Ordenes, Inventario (via Productos Ver), Facturacion (via Facturacion Read),
y Envios/guias si agrego item. NO muestra Usuarios/IAM/Integraciones/Clientes/Wallet/etc.
Frontend gating: front/central/src/shared/ui/sidebar.tsx + permissions-context.tsx (hasPermission).

## Fases

### Fase 0 - Migracion (back/migration) [additiva, idempotente]
- [ ] modelo email_verification_token.go (igual a password_reset_token: UserID, TokenHash, ExpiresAt, UsedAt)
- [ ] migrate_email_verification_tokens.go (AutoMigrate)
- [ ] seed_demo_role.go: crear role "demo" (scope=2, business_type=1, level=4, is_system=false)
      idempotente; crear permiso "Crear Ordenes" si falta; insertar role_permissions con los IDs.
- [ ] registrar en constructor.go del migration repo
- Correr: cd back/migration && go run cmd/main.go  (toca RDS prod, additivo)

### Fase 1 - Backend: registro demo + verificacion (services/auth/login o nuevo modulo demo)
- [ ] POST /auth/demo-register (publico): body {full_name, business_name, email, password}
      - crea Business (name=business_name, business_type_id=1, order_prefix auto) -> reusar usecase business
      - crea User (email, password propia hasheada bcrypt, is_active=false) 
      - crea business_staff(user, business, role demo=ID)
      - crea email_verification_token + envia correo SES con link {FRONTEND_BASE_URL}/verify-email?token=...
      - respuesta generica 200
- [ ] POST /auth/verify-email: valida token -> user.is_active=true, marca token usado
      -> dispara provisioning de integraciones test (Fase 2)
- Reusar: storefront register (publico) como patron; business CreateBusiness usecase; password_reset pattern.
- Cross-module: el modulo no debe importar repos de otros; replicar SELECT/creates necesarios o
  via RabbitMQ/usecase expuesto. Evaluar: lo mas simple es un modulo `demo` con su repo propio
  que escribe business/user/business_staff/integrations (todas tablas via models compartidos).

### Fase 2 - Auto-provisioning integraciones modo test
Al verificar (o al registrar), crear para el business:
- [ ] integration platform (type 6) is_default, is_testing=false, creds {} 
- [ ] integration siigo (type 8) is_testing=true, creds dummy {username, access_key, partner_id}
- [ ] integration envioclick (type 12) is_testing=true, config {auto_generate_guide_enabled:true}, creds {api_key:demo}
- [ ] invoicing_config: enabled=true, invoicing_integration_id = (siigo integration id), integration_ids null (todas)
- Reusar usecaseintegrations.CreateIntegration (cifra creds). base_url_test se inyecta del integration_type.
- Verificar trigger guia: shipments consumer order_created_autogen lee config auto_generate_guide_enabled.
- Verificar trigger factura: invoicing create_invoice usa invoicing_config.

### Fase 3 - Frontend
- [ ] Boton "Crea tu demo" en LoginForm (front/central .../login/ui/components/LoginForm.tsx)
- [ ] Modal flotante DemoRegisterForm (nombre, negocio, email, password) -> server action demoRegisterAction
- [ ] Pagina /verify-email (lee token, llama verify-email action, muestra resultado, link a login)
- [ ] Rutas publicas en (auth)/layout.tsx: agregar /verify-email (demo-register es modal en /login)
- [ ] Sidebar: el gating por permisos del rol demo ya limita modulos. Verificar que Envios/guias
      aparezca (revisar si hay item "Envios"/guias o esta dentro de shipments/delivery).

### Fase 4 - Pruebas locales (mocks) - ITERAR
- Levantar: infra + back/testing (mocks) + central + frontend. Mapear back-testing->127.0.0.1 en /etc/hosts.
- Registrar demo con secamc93@gmail.com, leer correo via Gmail MCP, abrir verify link, login.
- Crear orden -> ver factura (mock siigo) + guia (mock envioclick) generadas. Validar en DB y UI.
- Corregir hasta que el flujo completo funcione. Solo al final levantar todo junto.

## Estado actual
- [x] worktree creado /home/cam/Desktop/probability-demo (rama feat/demo-autoregistro) + env copiado
- [x] exploracion + datos DB
- [x] Fase 0 CODIGO escrito y compila (go build migration OK):
      - shared/models/email_verification_token.go
      - internal/infra/repository/migrate_demo_autoregistro.go (AutoMigrate + seed rol demo + permisos + permiso "Crear Ordenes")
      - registrado en constructor.go (migrateDemoAutoregistro)
      - PENDIENTE correr: cd back/migration && go run cmd/main.go (toca RDS prod, additivo) -> correr antes de probar Fase 1
- [x] Fase 1 CODIGO escrito y compila (go build central worktree OK):
      modulo NUEVO services/auth/demo (bundle.go cableado en services/auth/bundle.go):
      - domain/dtos.go, ports.go (IDemoRepository, IEmailSender)
      - infra/secondary/repository/repository.go: EmailExists, BusinessCodeExists, GetDemoRoleID,
        CreateDemoAccount (TX: business minimal sin resource_config + user is_active=false +
        user_businesses + business_staff(rol demo) + email_verification_token),
        GetValidEmailVerificationToken, ActivateUserAndConsumeToken
      - app: demo-register.go (valida, bcrypt, slug code unico, prefix, token, envia SES),
        verify-email.go (valida token, activa user), email templates inline
      - infra/primary/handlers: POST /auth/demo-register, POST /auth/verify-email (PUBLICOS)
      - Business minimal: NO crea business_resource_configured -> GetUserRolesPermissions
        permite todos los permisos del rol demo (sin filtrado por recursos activos). VERIFICAR esto en prueba.
- [~] Migracion corriendo (background) para crear tabla + rol demo. VERIFICAR que quede:
      email_verification_tokens, role demo, permiso "Crear Ordenes", role_permissions del demo.
- [x] Fase 0 migracion CORRIDA y verificada: email_verification_tokens OK, rol demo id=7,
      permiso "Crear Ordenes" id=128, demo tiene 10 permisos.
- [x] Fase 1 PROBADA E2E (backend worktree corrido en :3050, binario /tmp/demo-central):
      - secamc93@gmail.com YA existe como user 4 -> usar alias secamc93+demoN@gmail.com
        (llega al mismo inbox; hay que verificar el alias en SES antes: aws sesv2 create-email-identity,
        leer correo AWS via Gmail MCP, GET al link de verificacion).
      - demo-register +demo2 -> 200, crea business 49 "Cafe Demo Dos" + user 52 (is_active=false)
        + business_staff rol demo + email_verification_token; SES envia correo (verificado en Gmail).
      - login con inactivo -> 403 "usuario inactivo" (gate OK).
      - verify-email con token del correo -> 200, activa user.
      - login tras verificar -> 200 (devuelve user+business+token cookie).
      - BUG ENCONTRADO Y CORREGIDO: GORM omitia is_active=false (zero value) y aplicaba default:true.
        Fix en repository.go CreateDemoAccount: Update("is_active", false) explicito tras Create. RECOMPILADO.
      - PENDIENTE refinamiento menor: login devuelve require_password_change=true en primer login
        (isFirstLogin = last_login_at==nil). Para demo no deberia forzar cambio (ya eligio su pass).
        Opcion: en ActivateUserAndConsumeToken setear last_login_at, o manejar en frontend. NO bloquea.
- [x] Fase 3 frontend CODIGO escrito:
      - server actions demoRegisterAction + verifyEmailAction en services/auth/login/infra/actions/index.ts
        (fetch directo a API_BASE_URL/auth/demo-register y /auth/verify-email)
      - componente DemoRegisterModal.tsx (modal flotante: full_name,business_name,email,password)
      - LoginForm.tsx: estado showDemoModal + boton "Crea tu demo gratis" bajo el submit + render modal
      - pagina app/(auth)/verify-email/page.tsx (auto-llama verifyEmailAction con ?token=, muestra resultado)
      - (auth)/layout.tsx: /verify-email agregado a isPublicPage
      - PENDIENTE PROBAR: levantar frontend worktree. PROBLEMA node_modules: worktree no los tiene y el
        symlink a main rompe Turbopack ("Symlink invalid, points out of filesystem root"). SOLUCION:
        pnpm install en el worktree (corriendo). Luego pnpm dev (ojo puerto 3000: matar next-server viejos).
      - Frontend worktree .env: API_BASE_URL=http://localhost:3050/api/v1 (server actions OK).
- [x] Fase 3 PROBADA E2E con Playwright (frontend worktree :3000, backend worktree :3050):
      - /login muestra boton "Crea tu demo gratis" -> modal abre con 4 campos.
      - registro UI con secamc93+demo3@gmail.com -> "Cuenta creada. Revisa tu correo".
      - correo SES llega (leido via Gmail), /verify-email?token=... -> "Cuenta verificada".
      - login -> landing /home.
      - GATING CORRECTO: sidebar solo muestra /home, /products(Inventario), /orders(Ordenes con
        sub-nav Envios/guias+Cotizaciones), /invoicing(Facturacion), + always-on /tickets y /subscription.
        Ocultos: Usuarios, IAM, Integraciones, Wallet, Clientes, Notificaciones, etc. = lo pedido.
      - BUG ENCONTRADO Y CORREGIDO (beneficia a TODO negocio vacio): dashboard GetOrdersByWeek panic
        "index out of range [0] with length 0" en dashboard/.../repository.go:1220 (loop de reversa
        accedia results[0] con slice vacio). Fix: reversa estandar segura (for i,j:=0,len-1; i<j). Recompilado.
        Ahora /home renderiza dashboard con ceros para el demo.
      - Setup worktree front: pnpm install en worktree (symlink node_modules rompe Turbopack). Frontend
        corre con: cd front/central && nohup pnpm dev (puerto 3000; matar next-server viejos antes).
      - require_password_change=true en login NO bloquea (frontend no fuerza cambio); cosmetico, opcional.
      - Pendiente menor UX: Tickets y Suscripcion salen siempre (hardcoded always-on). El usuario pidio solo
        ordenes/factura/guias/inventario -> opcional ocultarlos para rol demo.
- [x] Fase 2 PROVISIONING HECHO Y VERIFICADO EN DB:
      - services/auth/demo/internal/infra/secondary/repository/provisioning.go: replica cifrado AES-GCM
        (ENCRYPTION_KEY del env, formato {"encrypted":base64}); ProvisionDemoIntegrations(businessID,userID)
        crea: platform(tipo6,is_testing=false,is_default), siigo(tipo8,is_testing=true,creds dummy cifradas),
        envioclick(tipo12,is_testing=true,config {"auto_generate_guide_enabled":true}) + invoicing_config
        (enabled=true, auto_invoice=true, invoicing_integration_id=siigo). Idempotente por code/business.
      - repo.New ahora recibe ENCRYPTION_KEY (bundle pasa cfg.Get). CreatedByID=userID en integrations e
        invoicing_config (FK fk_integrations_created_by y fk_invoicing_configs_created_by exigen user valido).
      - verify-email AHORA idempotente/auto-reparable: corre provisionDemo() tanto en verify fresco como en
        token ya usado (self-heal si fallo antes). GetBusinessIDByUserID resuelve el negocio.
      - PROBADO: business 51 (secamc93+demo4@gmail.com / demo123456, user 54) tiene integraciones 200/201/202
        + invoicing_config 19. base_url_test usa back-testing (ya mapeado a 127.0.0.1 en /etc/hosts).
- [x] Fase 2 E2E FACTURA FUNCIONA (probado): demo4 crea orden manual -> factura auto via mock Siigo
      -> invoice 24498 status="issued" FE-0-1001 (business 51). Detalles del flujo:
      - Mock backend: cd back/testing && /tmp/demo-testing (build: go build -o /tmp/demo-testing cmd/main.go),
        levanta mocks 9090(softpymes) 9091(envioclick) 9095(siigo) etc. .env copiado de main. back-testing->127.0.0.1 ok.
      - Provisioning AHORA tambien crea BODEGA (warehouses) por defecto (Bogota, DANE 11001000) -> requerida
        para crear ordenes. Y hace UPSERT de integraciones (actualiza creds/config si ya existen = auto-heal).
      - Siigo creds: REQUIERE campos username, access_key, ACCOUNT_ID, partner_id (helpers.go decripta c/u).
        Faltaba account_id -> agregado. Mi cifrado AES-GCM es compatible con el core (username/access_key
        decriptaron ok).
      - Orden manual: POST /orders con platform="manual", invoiceable:true (sin esto invoicing dice
        "Orden no es facturable"), items formato {"sku","name","price","quantity"} (NO product_sku), y el
        negocio necesita bodega. payment_method_id=6 (COD). El IntegrationID se auto-rellena del platform.
      - OJO CACHE: al actualizar creds via re-verify hay que REINICIAR el backend (cache de creds en Redis se
        calienta al arrancar; update de DB no lo refresca).
      - OJO PROCESOS: cada restart dejaba backends zombies -> SIEMPRE pkill -9 -f "/tmp/demo-central" antes.
- [x] Fase 2 E2E GUIA FUNCIONA (probado): generacion MANUAL via POST /api/v1/shipments/generate
      con body {order_uuid, carrier:"envioclick", origin{daneCode,city,state,street,name,phone},
      destination{...}, packages:[{weight,height,width,length}], contentValue}. Como usuario demo (JWT
      business_id) resuelve business y GetActiveShippingCarrier -> envioclick test -> mock :9091
      POST /api/v2/shipment 200. Resultado: shipment 34739 status=pending tracking=SRV-665605808
      guide_url=...EC-005001.pdf, evento shipment.guide_generated. ESTE es el flujo realista del demo
      (el usuario genera la guia desde la orden). El autogen (quote) NO aplica a ordenes manuales, ok.

## *** DEMO FUNCIONAL COMPLETO E2E (MVP logrado) ***
Register(boton+modal)->correo SES->verify->login->modulos limitados(ordenes/inventario/factura/guias)
-> auto-provision(integraciones test+invoicing_config+bodega) -> crear orden -> FACTURA auto (mock Siigo)
-> generar GUIA (mock EnvioClick). Todo en modo test contra back/testing. Cuenta: secamc93+demo4@gmail.com.

## Pendiente (refinamientos / cierre)
- [ ] Verificar en UI (Playwright) que el demo genere la guia desde /orders (boton generar guia) y vea
      factura en /invoicing — para confirmar el flujo visual completo del usuario.
- [ ] Refinamiento: require_password_change=true en primer login (no forzar para demo). Opcion: setear
      last_login_at al activar, o manejar en frontend.
- [x] Refinamiento HECHO: ocultar Tickets y Suscripcion para rol demo. sidebar.tsx: isDemo =
      permissions?.role_name === 'demo'; canViewTickets = !isDemo; Suscripcion envuelto en {!isDemo && ...}.
      ADEMAS fix bug real: permissions-context no poblaba role_name al recargar (la respuesta trae
      role:{name} pero el front espera role_name plano) -> mapeado role_name = raw.role?.name. Esto
      arregla tambien el check role_name==='cliente_final' del layout. VERIFICADO en UI: demo muestra
      solo Inicio/Inventario/Ordenes/Facturacion (4 items).
- [ ] LIMPIEZA datos de prueba en RDS prod: se crearon businesses/users demo (46-51 aprox), integraciones,
      invoices, shipments de prueba. Decidir si limpiar (con permiso) o dejar.
- [ ] COMMIT cuando el usuario lo pida (rama feat/demo-autoregistro). NO push sin permiso. Incluye back
      (modulo demo, migracion, fix dashboard) + front (modal, paginas, gating). Recordar reglas: sin
      comentarios en Go/TS, ASCII en archivos largos.
- [ ] (viejo) Fase 2 E2E GUIA: maybeAutoGenerate (shipments order_created_autogen)
      exige (1) flag auto_generate_guide_enabled en la integracion DE LA ORDEN (platform 200, no envioclick)
      y (2) un QUOTE de envio guardado (SavedQuote con carrier/rate) — una orden MANUAL no tiene quote.
      El flujo esta disenado para checkout/storefront. Opciones para el demo:
      (a) generar la guia MANUALMENTE desde la orden (endpoint/boton "generar guia" -> mock envioclick 9091),
          que es realista ("el usuario genera la guia"). Buscar ese endpoint/flow.
      (b) cotizar+guardar quote para la orden y poner el flag en platform.
      Revisar como la UI genera guia manual (services/modules/shipments, /orders Envios, request-confirmation
      o un generate-guide). El mock envioclick responde POST /api/v2/quotation y /api/v2/shipment (tracker EC-xx).
- [ ] (viejo) Fase 2 E2E: levantar back/testing (mocks 9090/9091/9095):
      cd /home/cam/Desktop/probability-demo/back/testing && go run cmd/main.go (o copiar back/testing/.env).
      Login como secamc93+demo4@gmail.com en :3000, crear una ORDEN manual, y verificar:
      (a) factura auto-generada via invoicing (auto_invoice) -> mock siigo :9095 (back-testing:9095)
      (b) guia auto-generada via shipments consumer order_created_autogen -> mock envioclick :9091.
      Validar en UI (/orders, /invoicing, /orders Envios) y en DB (invoices, shipments). Iterar/corregir.
      OJO: el flujo factura/guia es async (consumers RabbitMQ) y depende de mas config (payment_method, etc.)
      -> probablemente requiera ajustes. Revisar logs del backend al crear la orden.
- [ ] (viejo) Fase 2: auto-provisionar integraciones TEST al verificar (o al registrar):
      platform(6), siigo(8) is_testing, envioclick(12) is_testing + config auto_generate_guide_enabled,
      invoicing_config enabled apuntando a siigo. Para que el demo cree orden -> factura(mock siigo 9095)
      + guia(mock envioclick 9091). NECESITA: levantar back/testing (mocks) y mapear back-testing->127.0.0.1
      en /etc/hosts (base_url_test usa hostname docker). Credenciales: reusar usecaseintegrations.CreateIntegration
      (cifra creds) o insertar via repo del modulo demo replicando el cifrado (mejor reusar el usecase del core).
- [ ] (viejo) Fase 3 PRUEBA: con frontend worktree en :3000 (o el que tome), Playwright: /login -> ver boton
      "Crea tu demo gratis" -> abrir modal -> registrar con secamc93+demoN@gmail.com (verificar alias SES antes)
      -> ver mensaje exito -> abrir /verify-email?token=... (sacar token del correo via Gmail) -> ver verificado
      -> login -> VERIFICAR sidebar solo muestra Ordenes/Inventario/Facturacion (+Envios/guias) y NO Usuarios/IAM.
- [ ] (resto) Fase 3 detalle previo: boton "Crea tu demo" en LoginForm -> modal DemoRegisterForm
      (full_name,business_name,email,password) -> server action POST /auth/demo-register; pagina /verify-email
      (lee token, POST /auth/verify-email); ruta publica /verify-email en (auth)/layout.tsx.
      Permite verificar visualmente el gating de modulos del rol demo (ordenes, factura, guias, inventario).
- [ ] Fase 2 integraciones test (platform/siigo/envioclick + invoicing_config + auto_generate_guide).
      Recordar: base_url_test usa hostname Docker back-testing -> para local mapear back-testing->127.0.0.1
      en /etc/hosts, o levantar back/testing (mocks 9090/9091/9095) y ajustar.
- NOTA estado servicios: backend del WORKTREE corre en :3050 (binario /tmp/demo-central). El backend main
  (forgot-password) fue detenido. Frontend main aun puede estar en :3000.

## Notas
- NO push. NO commit hasta que el usuario lo pida.
- Migraciones tocan RDS prod (additivas). Confirmar antes de correr cada una.
- back/central ya tiene shared/email (SES) y password_reset (referencia de tokens).
