import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../services/auth/business/ui/providers/business_provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';
import 'app_colors.dart';
import 'brand_color.dart';

class ActiveBrand {
  const ActiveBrand._();

  static Color colorOf(LoginProvider login, BusinessProvider business) {
    if (login.isSuperAdmin) {
      final id = business.selectedBusinessId;
      if (id == null) return AppColors.primary;
      final match =
          business.businessesSimple.where((b) => b.id == id).firstOrNull;
      return BrandColor.resolve(match?.primaryColor);
    }

    final wanted = login.rolesPermissions?.businessId;
    final own = login.businesses
            .where((b) => wanted == null || wanted == 0 || b.id == wanted)
            .firstOrNull ??
        (login.businesses.isNotEmpty ? login.businesses.first : null);
    return BrandColor.resolve(own?.primaryColor);
  }

  static Color watch(BuildContext context) => colorOf(
        context.watch<LoginProvider>(),
        context.watch<BusinessProvider>(),
      );
}
