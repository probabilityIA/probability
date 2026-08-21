import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../services/auth/business/ui/providers/business_provider.dart';
import '../../../services/auth/login/ui/providers/login_provider.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';
import '../../utils/formatters.dart';
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
  int? _businessId;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      final login = context.read<LoginProvider>();
      if (!login.isSuperAdmin) return;
      final biz = context.read<BusinessProvider>();
      if (biz.businessesSimple.isEmpty) biz.fetchBusinessesSimple();
      if (biz.selectedBusinessId != null) {
        setState(() => _businessId = biz.selectedBusinessId);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();
    final isSuperAdmin = login.isSuperAdmin;

    if (isSuperAdmin && _businessId == null) {
      return AppScaffold(
        title: widget.title,
        subtitle: 'Selecciona un negocio para continuar',
        body: _BusinessGate(onSelected: (id) => setState(() => _businessId = id)),
      );
    }

    final children = widget.builder(context, isSuperAdmin ? _businessId : null);
    final singleTab = widget.tabs.length < 2;

    return DefaultTabController(
      length: widget.tabs.length,
      initialIndex: widget.initialTab.clamp(0, widget.tabs.length - 1),
      child: AppScaffold(
        title: widget.title,
        subtitle: widget.subtitle,
        actions: [
          ...widget.actions,
          if (isSuperAdmin) _BusinessChip(businessId: _businessId, onClear: () => setState(() => _businessId = null)),
        ],
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

class _BusinessChip extends StatelessWidget {
  const _BusinessChip({required this.businessId, required this.onClear});

  final int? businessId;
  final VoidCallback onClear;

  @override
  Widget build(BuildContext context) {
    return Consumer<BusinessProvider>(
      builder: (context, provider, _) {
        final business = provider.businessesSimple
            .where((b) => b.id == businessId)
            .firstOrNull;
        final name = business?.name ?? 'Negocio $businessId';
        return Padding(
          padding: const EdgeInsets.only(right: 4),
          child: TextButton.icon(
            onPressed: () {
              context.read<BusinessProvider>().setSelectedBusinessId(null);
              onClear();
            },
            icon: const Icon(Icons.swap_horiz_rounded, size: 16),
            label: Text(
              name.length > 14 ? '${name.substring(0, 13)}...' : name,
              style: const TextStyle(fontSize: 12.5, fontWeight: FontWeight.w600),
            ),
          ),
        );
      },
    );
  }
}

class _BusinessGate extends StatelessWidget {
  const _BusinessGate({required this.onSelected});

  final ValueChanged<int> onSelected;

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
              onTap: () {
                context.read<BusinessProvider>().setSelectedBusinessId(business.id);
                onSelected(business.id);
              },
              child: Row(
                children: [
                  Container(
                    width: 38,
                    height: 38,
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: AppColors.primarySoft,
                      borderRadius: AppRadius.smAll,
                    ),
                    child: Text(
                      AppFormat.initials(business.name),
                      style: const TextStyle(
                        fontFamily: 'Inter',
                        fontWeight: FontWeight.w700,
                        fontSize: 13,
                        color: AppColors.primary,
                      ),
                    ),
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
