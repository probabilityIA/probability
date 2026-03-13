import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'core/network/api_client.dart';
import 'core/router/app_router.dart';
import 'core/storage/token_storage.dart';
// Auth providers
import 'services/auth/actions/ui/providers/action_provider.dart';
import 'services/auth/business/ui/providers/business_provider.dart';
import 'services/auth/login/ui/providers/login_provider.dart';
import 'services/auth/permissions/ui/providers/permission_provider.dart';
import 'services/auth/resources/ui/providers/resource_provider.dart';
import 'services/auth/roles/ui/providers/role_provider.dart';
import 'services/auth/users/ui/providers/user_provider.dart';
// Module providers
import 'services/modules/orders/ui/providers/order_provider.dart';
import 'services/modules/products/ui/providers/product_provider.dart';
import 'services/modules/customers/ui/providers/customer_provider.dart';
import 'services/modules/invoicing/ui/providers/invoicing_provider.dart';
import 'services/modules/orderstatus/ui/providers/orderstatus_provider.dart';
import 'services/modules/paymentstatus/ui/providers/paymentstatus_provider.dart';
import 'services/modules/fulfillmentstatus/ui/providers/fulfillmentstatus_provider.dart';
import 'services/modules/shipments/ui/providers/shipment_provider.dart';
import 'services/modules/warehouses/ui/providers/warehouse_provider.dart';
import 'services/modules/inventory/ui/providers/inventory_provider.dart';
import 'services/modules/drivers/ui/providers/drivers_provider.dart';
import 'services/modules/vehicles/ui/providers/vehicle_provider.dart';
import 'services/modules/routes/ui/providers/route_provider.dart';
import 'services/modules/dashboard/ui/providers/dashboard_provider.dart';
import 'services/modules/pay/ui/providers/pay_provider.dart';
import 'services/modules/wallet/ui/providers/wallet_provider.dart';
import 'services/modules/notification_config/ui/providers/notification_config_provider.dart';
import 'services/modules/storefront/ui/providers/storefront_provider.dart';
import 'services/modules/publicsite/ui/providers/publicsite_provider.dart';
import 'services/modules/website_config/ui/providers/website_config_provider.dart';
import 'services/modules/my_integrations/ui/providers/my_integrations_provider.dart';
// Integration providers
import 'services/integrations/core/ui/providers/integration_provider.dart';
import 'shared/theme/app_theme.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  final tokenStorage = TokenStorage();
  final apiClient = ApiClient();

  final loginProvider = LoginProvider(
    tokenStorage: tokenStorage,
    apiClient: apiClient,
  );

  runApp(ProbabilityApp(
    apiClient: apiClient,
    loginProvider: loginProvider,
  ));
}

class ProbabilityApp extends StatefulWidget {
  final ApiClient apiClient;
  final LoginProvider loginProvider;

  const ProbabilityApp({
    super.key,
    required this.apiClient,
    required this.loginProvider,
  });

  @override
  State<ProbabilityApp> createState() => _ProbabilityAppState();
}

class _ProbabilityAppState extends State<ProbabilityApp> {
  late final AppRouter _appRouter;

  @override
  void initState() {
    super.initState();
    _appRouter = AppRouter(loginProvider: widget.loginProvider);
    widget.loginProvider.restoreSession();
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: widget.loginProvider),
        ChangeNotifierProvider(
          create: (_) => UserProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => RoleProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => PermissionProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => BusinessProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => ResourceProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => ActionProvider(apiClient: widget.apiClient),
        ),
        // Module providers
        ChangeNotifierProvider(
          create: (_) => OrderProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => ProductProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => CustomerProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => InvoicingProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => OrderStatusProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => PaymentStatusProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => FulfillmentStatusProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => ShipmentProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => WarehouseProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => InventoryProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => DriverProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => VehicleProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => RouteProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => DashboardProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => PayProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => WalletProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => NotificationConfigProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => StorefrontProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => PublicSiteProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => WebsiteConfigProvider(apiClient: widget.apiClient),
        ),
        ChangeNotifierProvider(
          create: (_) => MyIntegrationsProvider(apiClient: widget.apiClient),
        ),
        // Integration core provider
        ChangeNotifierProvider(
          create: (_) => IntegrationProvider(apiClient: widget.apiClient),
        ),
      ],
      child: MaterialApp.router(
        title: 'Probability Central',
        theme: AppTheme.light,
        darkTheme: AppTheme.dark,
        themeMode: ThemeMode.system,
        routerConfig: _appRouter.router,
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}
