import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../services/auth/business/ui/providers/business_provider.dart';
import '../theme/app_colors.dart';
import '../theme/app_tokens.dart';
import 'ui/ui.dart';

Future<void> showBusinessSwitcher(BuildContext context) async {
  final provider = context.read<BusinessProvider>();
  if (provider.businessesSimple.isEmpty) provider.fetchBusinessesSimple();

  await showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.surface,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(AppRadius.xl)),
    ),
    builder: (_) => const _BusinessSheet(),
  );
}

class _BusinessSheet extends StatelessWidget {
  const _BusinessSheet();

  @override
  Widget build(BuildContext context) {
    final maxHeight = MediaQuery.sizeOf(context).height * 0.78;

    return SafeArea(
      top: false,
      child: ConstrainedBox(
        constraints: BoxConstraints(maxHeight: maxHeight),
        child: Consumer<BusinessProvider>(
          builder: (context, provider, _) {
            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 14, 20, 10),
                  child: Row(
                    children: [
                      Expanded(
                        child: Text(
                          'Cambiar de negocio',
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                      ),
                      if (provider.isLoading)
                        const SizedBox(
                          height: 16,
                          width: 16,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        ),
                    ],
                  ),
                ),
                if (provider.businessesSimple.isEmpty && !provider.isLoading)
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 8, 20, 24),
                    child: Text(
                      provider.error ?? 'No hay negocios disponibles.',
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  )
                else
                  Flexible(
                    child: ListView.separated(
                      shrinkWrap: true,
                      padding: const EdgeInsets.fromLTRB(16, 4, 16, 18),
                      itemCount: provider.businessesSimple.length,
                      separatorBuilder: (context, index) =>
                          const SizedBox(height: 8),
                      itemBuilder: (context, index) {
                        final business = provider.businessesSimple[index];
                        final active =
                            business.id == provider.selectedBusinessId;
                        return AppCard(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 13, vertical: 12),
                          onTap: () {
                            provider.setSelectedBusinessId(business.id);
                            Navigator.pop(context);
                          },
                          child: Row(
                            children: [
                              BrandLogo(
                                name: business.name,
                                imageUrl: business.logoUrl,
                                size: 38,
                                radius: AppRadius.sm,
                                padding: 4,
                              ),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Text(
                                  business.name,
                                  style:
                                      Theme.of(context).textTheme.titleSmall,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                ),
                              ),
                              if (active)
                                const Icon(Icons.check_circle_rounded,
                                    size: 20, color: AppColors.primary),
                            ],
                          ),
                        );
                      },
                    ),
                  ),
              ],
            );
          },
        ),
      ),
    );
  }
}
