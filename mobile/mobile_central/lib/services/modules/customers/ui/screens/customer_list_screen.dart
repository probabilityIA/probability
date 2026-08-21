import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/customer_provider.dart';
import 'customer_detail_screen.dart';

class CustomerListScreen extends StatefulWidget {
  const CustomerListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<CustomerListScreen> createState() => _CustomerListScreenState();
}

class _CustomerListScreenState extends State<CustomerListScreen> {
  final _searchController = TextEditingController();
  final _scrollController = ScrollController();
  Timer? _debounce;

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final position = _scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 320) {
      context.read<CustomerProvider>().loadMore(businessId: widget.businessId);
    }
  }

  void _refresh() {
    final provider = context.read<CustomerProvider>();
    provider.setPage(1);
    provider.fetchCustomers(businessId: widget.businessId);
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<CustomerProvider>().setSearch(value);
      _refresh();
    });
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Clientes',
      subtitle: 'Directorio y compras',
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
            child: AppSearchField(
              controller: _searchController,
              hintText: 'Nombre, correo, telefono o documento',
              onChanged: _onSearch,
            ),
          ),
          Expanded(
            child: Consumer<CustomerProvider>(
              builder: (context, provider, _) {
                if (provider.isLoading && provider.customers.isEmpty) {
                  return const AppListSkeleton();
                }
                if (provider.error != null && provider.customers.isEmpty) {
                  return AppErrorState(message: provider.error!, onRetry: _refresh);
                }
                if (provider.customers.isEmpty) {
                  return AppEmptyState(
                    icon: Icons.people_alt_outlined,
                    title: 'Sin clientes',
                    message: 'Los clientes se crean solos cuando entran ordenes de tus canales.',
                    actionLabel: 'Actualizar',
                    onAction: _refresh,
                  );
                }

                return RefreshIndicator(
                  onRefresh: () async => _refresh(),
                  color: AppColors.primary,
                  child: ListView.separated(
                    controller: _scrollController,
                    physics: const AlwaysScrollableScrollPhysics(),
                    padding: AppSpacing.page,
                    itemCount: provider.customers.length + 1,
                    separatorBuilder: (context, index) => const SizedBox(height: 10),
                    itemBuilder: (context, index) {
                      if (index == provider.customers.length) {
                        return Padding(
                          padding: const EdgeInsets.symmetric(vertical: 18),
                          child: Center(
                            child: provider.isLoadingMore
                                ? const SizedBox(
                                    height: 22,
                                    width: 22,
                                    child: CircularProgressIndicator(strokeWidth: 2.2),
                                  )
                                : Text(
                                    provider.hasMore
                                        ? 'Desliza para cargar mas'
                                        : '${provider.customers.length} de ${provider.pagination?.total ?? provider.customers.length} clientes',
                                    style: Theme.of(context).textTheme.labelSmall,
                                  ),
                          ),
                        );
                      }
                      final customer = provider.customers[index];
                      return _CustomerCard(
                        customer: customer,
                        onTap: () => Navigator.of(context).push(
                          MaterialPageRoute(
                            builder: (_) => CustomerDetailScreen(
                              customerId: customer.id,
                              businessId: widget.businessId,
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

class _CustomerCard extends StatelessWidget {
  const _CustomerCard({required this.customer, required this.onTap});

  final CustomerInfo customer;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final detail = customer is CustomerDetail ? customer as CustomerDetail : null;

    return AppCard(
      padding: const EdgeInsets.all(13),
      onTap: onTap,
      child: Row(
        children: [
          Container(
            width: 42,
            height: 42,
            alignment: Alignment.center,
            decoration: const BoxDecoration(
              color: AppColors.primarySoft,
              shape: BoxShape.circle,
            ),
            child: Text(
              AppFormat.initials(customer.name),
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 13,
                fontWeight: FontWeight.w700,
                color: AppColors.primary,
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  customer.name,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  customer.email ?? customer.phone,
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          const SizedBox(width: 10),
          if (detail != null)
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  AppFormat.money(detail.totalSpent),
                  style: theme.textTheme.titleSmall?.copyWith(fontSize: 14),
                ),
                const SizedBox(height: 2),
                Text('${detail.orderCount} ordenes', style: theme.textTheme.labelSmall),
              ],
            ),
        ],
      ),
    );
  }
}
