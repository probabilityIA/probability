import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'core/network/api_client.dart';
import 'core/router/app_router.dart';
import 'core/storage/token_storage.dart';
import 'services/auth/actions/ui/providers/action_provider.dart';
import 'services/auth/business/ui/providers/business_provider.dart';
import 'services/auth/login/ui/providers/login_provider.dart';
import 'services/auth/permissions/ui/providers/permission_provider.dart';
import 'services/auth/resources/ui/providers/resource_provider.dart';
import 'services/auth/roles/ui/providers/role_provider.dart';
import 'services/auth/users/ui/providers/user_provider.dart';
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
