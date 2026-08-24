import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/integration_provider.dart';

const Map<String, String> integrationCategoryLabels = {
  'ecommerce': 'Tiendas',
  'invoicing': 'Facturacion',
  'messaging': 'Mensajeria',
  'shipping': 'Transporte',
  'payment': 'Pagos',
  'platform': 'Plataforma',
};

class IntegrationListScreen extends StatefulWidget {
  const IntegrationListScreen({super.key});

  @override
  State<IntegrationListScreen> createState() => _IntegrationListScreenState();
}

class _IntegrationListScreenState extends State<IntegrationListScreen> {
  final _searchController = TextEditingController();
  String _category = '';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _refresh() {
    context.read<IntegrationProvider>().fetchIntegrationTypes();
  }

  String _categoryOf(IntegrationType type) =>
      type.category?.code ?? type.integrationCategory?.code ?? '';

  @override
  Widget build(BuildContext context) {
    return Consumer<IntegrationProvider>(
      builder: (context, provider, _) {
        if (provider.isLoadingTypes && provider.integrationTypes.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.integrationTypes.isEmpty) {
          return AppEmptyState(
            icon: Icons.hub_outlined,
            title: 'Catalogo vacio',
            message: 'No hay integraciones disponibles para conectar.',
            actionLabel: 'Actualizar',
            onAction: _refresh,
          );
        }

        final categories = <String>{
          ...provider.integrationTypes.map(_categoryOf).where((c) => c.isNotEmpty),
        }.toList()
          ..sort();

        final query = _searchController.text.trim().toLowerCase();
        final rows = provider.integrationTypes.where((type) {
          final matchesCategory = _category.isEmpty || _categoryOf(type) == _category;
          final matchesQuery = query.isEmpty || type.name.toLowerCase().contains(query);
          return matchesCategory && matchesQuery;
        }).toList();

        return RefreshIndicator(
          onRefresh: () async => _refresh(),
          color: AppColors.primary,
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
                child: AppSearchField(
                  controller: _searchController,
                  hintText: 'Buscar integracion',
                  onChanged: (_) => setState(() {}),
                ),
              ),
              AppFilterChips(
                options: [
                  (value: '', label: 'Todas'),
                  ...categories.map(
                    (code) => (
                      value: code,
                      label: integrationCategoryLabels[code] ?? code,
                    ),
                  ),
                ],
                selected: _category,
                onSelected: (value) => setState(() => _category = value),
              ),
              const SizedBox(height: 4),
              Expanded(
                child: rows.isEmpty
                    ? const AppEmptyState(
                        icon: Icons.search_off_rounded,
                        title: 'Sin resultados',
                        message: 'Ninguna integracion coincide con el filtro.',
                      )
                    : GridView.builder(
                        physics: const AlwaysScrollableScrollPhysics(),
                        padding: AppSpacing.page,
                        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2,
                          mainAxisSpacing: 11,
                          crossAxisSpacing: 11,
                          mainAxisExtent: 132,
                        ),
                        itemCount: rows.length,
                        itemBuilder: (context, index) => _CatalogTile(
                          type: rows[index],
                          categoryCode: _categoryOf(rows[index]),
                        ),
                      ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _CatalogTile extends StatelessWidget {
  const _CatalogTile({required this.type, required this.categoryCode});

  final IntegrationType type;
  final String categoryCode;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final unavailable = type.inDevelopment == true || !type.isActive;

    return AppCard(
      padding: const EdgeInsets.all(13),
      onTap: unavailable
          ? null
          : () {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Conectar ${type.name} desde la plataforma web')),
              );
            },
      child: Opacity(
        opacity: unavailable ? 0.55 : 1,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                BrandLogo(
                  name: type.name,
                  imageUrl: type.imageUrl,
                  size: 42,
                  radius: 11,
                  padding: 6,
                ),
                const Spacer(),
                if (type.inDevelopment == true)
                  const AppStatusChip(
                    dense: true,
                    label: 'Pronto',
                    tone: AppStatusTone.warning,
                  ),
              ],
            ),
            const Spacer(),
            Text(
              type.name,
              style: theme.textTheme.titleSmall,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 3),
            Text(
              integrationCategoryLabels[categoryCode] ?? categoryCode,
              style: theme.textTheme.labelSmall,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }
}
