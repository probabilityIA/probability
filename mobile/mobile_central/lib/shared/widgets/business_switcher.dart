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

class _BusinessSheet extends StatefulWidget {
  const _BusinessSheet();

  @override
  State<_BusinessSheet> createState() => _BusinessSheetState();
}

class _BusinessSheetState extends State<_BusinessSheet> {
  final _searchController = TextEditingController();
  String _query = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final maxHeight = MediaQuery.sizeOf(context).height * 0.78;

    return SafeArea(
      top: false,
      child: ConstrainedBox(
        constraints: BoxConstraints(maxHeight: maxHeight),
        child: Consumer<BusinessProvider>(
          builder: (context, provider, _) {
            final all = provider.businessesSimple;
            final query = _query.trim().toLowerCase();
            final rows = query.isEmpty
                ? all
                : all
                    .where((b) => b.name.toLowerCase().contains(query))
                    .toList();

            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 14, 20, 6),
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
                        )
                      else
                        Text(
                          query.isEmpty
                              ? '${all.length} negocios'
                              : '${rows.length} de ${all.length}',
                          style: Theme.of(context).textTheme.labelSmall,
                        ),
                    ],
                  ),
                ),
                if (all.length > 6)
                  Padding(
                    padding: const EdgeInsets.fromLTRB(16, 6, 16, 4),
                    child: AppSearchField(
                      controller: _searchController,
                      hintText: 'Buscar negocio',
                      onChanged: (value) => setState(() => _query = value),
                    ),
                  ),
                if (rows.isEmpty && !provider.isLoading)
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 12, 20, 28),
                    child: Text(
                      all.isEmpty
                          ? (provider.error ?? 'No hay negocios disponibles.')
                          : 'Ningun negocio coincide con la busqueda.',
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  )
                else
                  Flexible(
                    child: ListView.separated(
                      shrinkWrap: true,
                      padding: const EdgeInsets.fromLTRB(16, 4, 16, 18),
                      itemCount: rows.length,
                      separatorBuilder: (context, index) =>
                          const SizedBox(height: 8),
                      itemBuilder: (context, index) {
                        final business = rows[index];
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
