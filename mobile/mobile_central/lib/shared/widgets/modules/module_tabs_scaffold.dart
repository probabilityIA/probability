import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../services/auth/business/ui/providers/business_provider.dart';
import '../../../services/auth/login/ui/providers/login_provider.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';
import '../ui/ui.dart';

typedef ModuleTabsBuilder = List<Widget> Function(BuildContext context, int? businessId);

class ModuleTabsScaffold extends StatefulWidget {
  const ModuleTabsScaffold({
    super.key,
    required this.title,
    required this.tabs,
    required this.builder,
    this.subtitle,
    this.initialTab = 0,
    this.actions = const [],
  });

  final String title;
  final String? subtitle;
  final List<String> tabs;
  final ModuleTabsBuilder builder;
  final int initialTab;
  final List<Widget> actions;

  @override
  State<ModuleTabsScaffold> createState() => _ModuleTabsScaffoldState();
}

class _ModuleTabsScaffoldState extends State<ModuleTabsScaffold> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      final login = context.read<LoginProvider>();
      if (!login.isSuperAdmin) return;
      final biz = context.read<BusinessProvider>();
      if (biz.businessesSimple.isEmpty) biz.fetchBusinessesSimple();
    });
  }

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();
    final isSuperAdmin = login.isSuperAdmin;
    final businessId = context.watch<BusinessProvider>().selectedBusinessId;

    if (isSuperAdmin && businessId == null) {
      return const AppScaffold(body: _BusinessGate());
    }

    final children = widget.builder(context, isSuperAdmin ? businessId : null);
    final singleTab = widget.tabs.length < 2;

    return DefaultTabController(
      length: widget.tabs.length,
      initialIndex: widget.initialTab.clamp(0, widget.tabs.length - 1),
      child: AppScaffold(
        title: widget.title,
        subtitle: widget.subtitle,
        actions: widget.actions,
        bottom: singleTab
            ? null
            : PreferredSize(
                preferredSize: const Size.fromHeight(46),
                child: Align(
                  alignment: Alignment.centerLeft,
                  child: TabBar(
                    isScrollable: true,
                    tabAlignment: TabAlignment.start,
                    padding: const EdgeInsets.symmetric(horizontal: 8),
                    tabs: widget.tabs.map((label) => Tab(height: 44, text: label)).toList(),
                  ),
                ),
              ),
        body: singleTab ? children.first : TabBarView(children: children),
      ),
    );
  }
}

class _BusinessGate extends StatelessWidget {
  const _BusinessGate();

  @override
  Widget build(BuildContext context) {
    return Consumer<BusinessProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading) return const AppLoading(label: 'Cargando negocios');

        if (provider.businessesSimple.isEmpty) {
          return AppEmptyState(
            icon: Icons.apartment_outlined,
            title: 'Sin negocios disponibles',
            message: provider.error ?? 'Selecciona el negocio con el que quieres operar.',
            actionLabel: 'Cargar negocios',
            onAction: provider.fetchBusinessesSimple,
          );
        }

        return ListView.separated(
          padding: AppSpacing.page,
          itemCount: provider.businessesSimple.length,
          separatorBuilder: (context, index) => const SizedBox(height: 10),
          itemBuilder: (context, index) {
            final business = provider.businessesSimple[index];
            return AppCard(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 13),
              onTap: () => context
                  .read<BusinessProvider>()
                  .setSelectedBusinessId(business.id),
              child: Row(
                children: [
                  BrandLogo(
                    name: business.name,
                    imageUrl: business.logoUrl,
                    size: 40,
                    radius: AppRadius.sm,
                    padding: 4,
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(business.name, style: Theme.of(context).textTheme.titleSmall),
                        const SizedBox(height: 2),
                        Text('ID ${business.id}', style: Theme.of(context).textTheme.labelSmall),
                      ],
                    ),
                  ),
                  const Icon(Icons.chevron_right_rounded, color: AppColors.textDisabled),
                ],
              ),
            );
          },
        );
      },
    );
  }
}
