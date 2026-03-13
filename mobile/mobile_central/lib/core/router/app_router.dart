import 'package:go_router/go_router.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';
import '../../services/auth/login/ui/screens/login_screen.dart';
import '../../services/auth/users/ui/screens/user_list_screen.dart';
import '../../shared/widgets/home_screen.dart';

class AppRouter {
  final LoginProvider loginProvider;

  AppRouter({required this.loginProvider});

  late final GoRouter router = GoRouter(
    refreshListenable: loginProvider,
    initialLocation: '/login',
    redirect: (context, state) {
      final isLoggedIn = loginProvider.isLoggedIn;
      final isLoginRoute = state.matchedLocation == '/login';

      if (!isLoggedIn && !isLoginRoute) return '/login';
      if (isLoggedIn && isLoginRoute) return '/home';
      return null;
    },
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomeScreen(),
      ),
      GoRoute(
        path: '/users',
        builder: (context, state) => const UserListScreen(),
      ),
    ],
  );
}
