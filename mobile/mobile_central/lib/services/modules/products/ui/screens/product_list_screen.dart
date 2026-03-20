import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/product_provider.dart';
import '../../domain/entities.dart';

class ProductListScreen extends StatefulWidget {
  final int? businessId;

  const ProductListScreen({super.key, this.businessId});

  @override
  State<ProductListScreen> createState() => _ProductListScreenState();
}

class _ProductListScreenState extends State<ProductListScreen> {
  final _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProductProvider>().fetchProducts(businessId: widget.businessId);
    });
  }

  @override
  void didUpdateWidget(ProductListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      context.read<ProductProvider>().resetFilters();
      context
          .read<ProductProvider>()
          .fetchProducts(businessId: widget.businessId);
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _onSearch(String value) {
    final provider = context.read<ProductProvider>();
    provider.setFilters(name: value);
    provider.fetchProducts(businessId: widget.businessId);
  }

  void _goToPage(int page) {
    final provider = context.read<ProductProvider>();
    provider.setPage(page);
    provider.fetchProducts(businessId: widget.businessId);
  }

  Future<void> _showCreateEditSheet({Product? product}) async {
    final result = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      builder: (context) => _ProductFormSheet(
        product: product,
        businessId: widget.businessId,
      ),
    );

    if (result == true && mounted) {
      context.read<ProductProvider>().fetchProducts(businessId: widget.businessId);
    }
  }

  Future<void> _confirmDelete(Product product) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Eliminar producto'),
        content: Text(
            'Estas seguro de que deseas eliminar "${product.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(context).colorScheme.error,
            ),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      final provider = context.read<ProductProvider>();
      final success = await provider.deleteProduct(product.id);
      if (success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Producto eliminado')),
        );
        provider.fetchProducts(businessId: widget.businessId);
      } else if (mounted && provider.error != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(provider.error!)),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Productos'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showCreateEditSheet(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<ProductProvider>(
        builder: (context, provider, child) {
          return Column(
            children: [
              // Search bar
              Padding(
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Buscar por nombre...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              _onSearch('');
                            },
                          )
                        : null,
                    isDense: true,
                  ),
                  onSubmitted: _onSearch,
                  textInputAction: TextInputAction.search,
                ),
              ),

              const SizedBox(height: 4),

              // Content
              Expanded(
                child: _buildContent(provider),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _buildContent(ProductProvider provider) {
    final colorScheme = Theme.of(context).colorScheme;

    if (provider.isLoading && provider.products.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.products.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 48, color: colorScheme.error),
              const SizedBox(height: 16),
              Text(
                provider.error!,
                textAlign: TextAlign.center,
                style: TextStyle(color: colorScheme.error),
              ),
              const SizedBox(height: 16),
              FilledButton.icon(
                onPressed: () =>
                    provider.fetchProducts(businessId: widget.businessId),
                icon: const Icon(Icons.refresh),
                label: const Text('Reintentar'),
              ),
            ],
          ),
        ),
      );
    }

    if (provider.products.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.inventory_2_outlined,
                size: 64, color: colorScheme.outline),
            const SizedBox(height: 16),
            Text(
              'No hay productos',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    color: colorScheme.outline,
                  ),
            ),
            const SizedBox(height: 8),
            FilledButton.icon(
              onPressed: () => _showCreateEditSheet(),
              icon: const Icon(Icons.add),
              label: const Text('Crear producto'),
            ),
          ],
        ),
      );
    }

    return Column(
      children: [
        Expanded(
          child: RefreshIndicator(
            onRefresh: () =>
                provider.fetchProducts(businessId: widget.businessId),
            child: ListView.builder(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: provider.products.length,
              itemBuilder: (context, index) {
                final product = provider.products[index];
                return _ProductCard(
                  product: product,
                  onEdit: () => _showCreateEditSheet(product: product),
                  onDelete: () => _confirmDelete(product),
                );
              },
            ),
          ),
        ),

        // Pagination
        if (provider.pagination != null) _buildPagination(provider),
      ],
    );
  }

  Widget _buildPagination(ProductProvider provider) {
    final pagination = provider.pagination!;
    if (pagination.lastPage <= 1) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
          top: BorderSide(
            color: Theme.of(context).colorScheme.outlineVariant,
          ),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            '${pagination.total} productos',
            style: Theme.of(context).textTheme.bodySmall,
          ),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed: pagination.hasPrev
                    ? () => _goToPage(pagination.currentPage - 1)
                    : null,
                iconSize: 20,
              ),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Text(
                  '${pagination.currentPage} / ${pagination.lastPage}',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: pagination.hasNext
                    ? () => _goToPage(pagination.currentPage + 1)
                    : null,
                iconSize: 20,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _ProductCard extends StatelessWidget {
  final Product product;
  final VoidCallback onEdit;
  final VoidCallback onDelete;

  const _ProductCard({
    required this.product,
    required this.onEdit,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        onTap: onEdit,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              // Product image
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: SizedBox(
                  width: 56,
                  height: 56,
                  child: product.imageUrl != null &&
                          product.imageUrl!.isNotEmpty
                      ? Image.network(
                          product.imageUrl!,
                          fit: BoxFit.cover,
                          errorBuilder: (context, error, stackTrace) =>
                              _buildImagePlaceholder(colorScheme),
                        )
                      : _buildImagePlaceholder(colorScheme),
                ),
              ),

              const SizedBox(width: 12),

              // Product info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      product.name,
                      style: textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    if (product.sku.isNotEmpty)
                      Text(
                        'SKU: ${product.sku}',
                        style: textTheme.bodySmall?.copyWith(
                          color: colorScheme.outline,
                        ),
                      ),
                    const SizedBox(height: 6),
                    Row(
                      children: [
                        // Price
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8, vertical: 2),
                          decoration: BoxDecoration(
                            color: colorScheme.primaryContainer,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            '${product.currency} ${_formatAmount(product.price)}',
                            style: textTheme.labelMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                              color: colorScheme.onPrimaryContainer,
                            ),
                          ),
                        ),

                        const SizedBox(width: 8),

                        // Stock
                        Icon(
                          Icons.inventory_2_outlined,
                          size: 14,
                          color: colorScheme.outline,
                        ),
                        const SizedBox(width: 2),
                        Text(
                          product.manageStock
                              ? '${product.stock}'
                              : 'Ilimitado',
                          style: textTheme.bodySmall?.copyWith(
                            color: product.manageStock && product.stock <= 0
                                ? colorScheme.error
                                : colorScheme.outline,
                            fontWeight: product.manageStock && product.stock <= 0
                                ? FontWeight.bold
                                : FontWeight.normal,
                          ),
                        ),

                        const Spacer(),

                        // Active status
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8, vertical: 2),
                          decoration: BoxDecoration(
                            color: product.isActive
                                ? const Color(0xFFDCFCE7)
                                : const Color(0xFFFEE2E2),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            product.isActive ? 'Activo' : 'Inactivo',
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                              color: product.isActive
                                  ? const Color(0xFF166534)
                                  : const Color(0xFF991B1B),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),

              const SizedBox(width: 4),

              // Actions menu
              PopupMenuButton<String>(
                onSelected: (value) {
                  if (value == 'edit') {
                    onEdit();
                  } else if (value == 'delete') {
                    onDelete();
                  }
                },
                itemBuilder: (context) => [
                  const PopupMenuItem(
                    value: 'edit',
                    child: Row(
                      children: [
                        Icon(Icons.edit_outlined, size: 18),
                        SizedBox(width: 8),
                        Text('Editar'),
                      ],
                    ),
                  ),
                  PopupMenuItem(
                    value: 'delete',
                    child: Row(
                      children: [
                        Icon(Icons.delete_outline,
                            size: 18, color: colorScheme.error),
                        const SizedBox(width: 8),
                        Text('Eliminar',
                            style: TextStyle(color: colorScheme.error)),
                      ],
                    ),
                  ),
                ],
                icon: Icon(Icons.more_vert, color: colorScheme.outline),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildImagePlaceholder(ColorScheme colorScheme) {
    return Container(
      color: colorScheme.surfaceContainerHighest,
      child: Center(
        child: Icon(
          Icons.image_outlined,
          color: colorScheme.outline,
          size: 24,
        ),
      ),
    );
  }

  static String _formatAmount(double amount) {
    if (amount == amount.roundToDouble()) {
      return amount.toStringAsFixed(0).replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},');
    }
    final parts = amount.toStringAsFixed(2).split('.');
    final integer = parts[0].replaceAllMapped(
        RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},');
    return '$integer.${parts[1]}';
  }
}

class _ProductFormSheet extends StatefulWidget {
  final Product? product;
  final int? businessId;

  const _ProductFormSheet({this.product, this.businessId});

  @override
  State<_ProductFormSheet> createState() => _ProductFormSheetState();
}

class _ProductFormSheetState extends State<_ProductFormSheet> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _nameController;
  late final TextEditingController _skuController;
  late final TextEditingController _priceController;
  late final TextEditingController _stockController;
  late final TextEditingController _descriptionController;
  late String _currency;
  late bool _isActive;
  bool _isLoading = false;

  bool get _isEditing => widget.product != null;

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController(text: widget.product?.name ?? '');
    _skuController = TextEditingController(text: widget.product?.sku ?? '');
    _priceController = TextEditingController(
        text: widget.product?.price.toString() ?? '');
    _stockController = TextEditingController(
        text: widget.product?.stock.toString() ?? '0');
    _descriptionController =
        TextEditingController(text: widget.product?.description ?? '');
    _currency = widget.product?.currency ?? 'COP';
    _isActive = widget.product?.isActive ?? true;
  }

  @override
  void dispose() {
    _nameController.dispose();
    _skuController.dispose();
    _priceController.dispose();
    _stockController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    final provider = context.read<ProductProvider>();
    bool success;

    if (_isEditing) {
      success = await provider.updateProduct(
        widget.product!.id,
        UpdateProductDTO(
          name: _nameController.text.trim(),
          sku: _skuController.text.trim(),
          price: double.tryParse(_priceController.text) ?? 0,
          stock: int.tryParse(_stockController.text) ?? 0,
          description: _descriptionController.text.trim().isNotEmpty
              ? _descriptionController.text.trim()
              : null,
          currency: _currency,
          isActive: _isActive,
        ),
      );
    } else {
      final product = await provider.createProduct(
        CreateProductDTO(
          businessId: widget.businessId ?? 0,
          name: _nameController.text.trim(),
          sku: _skuController.text.trim(),
          price: double.tryParse(_priceController.text) ?? 0,
          stock: int.tryParse(_stockController.text) ?? 0,
          description: _descriptionController.text.trim().isNotEmpty
              ? _descriptionController.text.trim()
              : null,
          currency: _currency,
          isActive: _isActive,
        ),
      );
      success = product != null;
    }

    if (!mounted) return;

    setState(() => _isLoading = false);

    if (success) {
      Navigator.of(context).pop(true);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(_isEditing
              ? 'Producto actualizado'
              : 'Producto creado'),
        ),
      );
    } else if (provider.error != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(provider.error!)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return DraggableScrollableSheet(
      initialChildSize: 0.85,
      minChildSize: 0.5,
      maxChildSize: 0.95,
      expand: false,
      builder: (context, scrollController) {
        return Container(
          decoration: BoxDecoration(
            color: colorScheme.surface,
            borderRadius:
                const BorderRadius.vertical(top: Radius.circular(20)),
          ),
          child: Column(
            children: [
              // Handle
              Container(
                margin: const EdgeInsets.only(top: 8),
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: colorScheme.outlineVariant,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),

              // Header
              Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  children: [
                    Text(
                      _isEditing ? 'Editar Producto' : 'Nuevo Producto',
                      style:
                          Theme.of(context).textTheme.titleLarge?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                    ),
                    const Spacer(),
                    IconButton(
                      onPressed: () => Navigator.of(context).pop(),
                      icon: const Icon(Icons.close),
                    ),
                  ],
                ),
              ),

              const Divider(height: 1),

              // Form
              Expanded(
                child: SingleChildScrollView(
                  controller: scrollController,
                  padding: const EdgeInsets.all(16),
                  child: Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        TextFormField(
                          controller: _nameController,
                          decoration:
                              const InputDecoration(labelText: 'Nombre *'),
                          validator: (v) => (v == null || v.trim().isEmpty)
                              ? 'Nombre requerido'
                              : null,
                          textInputAction: TextInputAction.next,
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _skuController,
                          decoration:
                              const InputDecoration(labelText: 'SKU *'),
                          validator: (v) => (v == null || v.trim().isEmpty)
                              ? 'SKU requerido'
                              : null,
                          textInputAction: TextInputAction.next,
                        ),
                        const SizedBox(height: 16),
                        Row(
                          children: [
                            Expanded(
                              flex: 2,
                              child: TextFormField(
                                controller: _priceController,
                                decoration: const InputDecoration(
                                    labelText: 'Precio *'),
                                keyboardType:
                                    const TextInputType.numberWithOptions(
                                        decimal: true),
                                validator: (v) {
                                  if (v == null || v.trim().isEmpty) {
                                    return 'Precio requerido';
                                  }
                                  if (double.tryParse(v) == null) {
                                    return 'Precio invalido';
                                  }
                                  return null;
                                },
                                textInputAction: TextInputAction.next,
                              ),
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: DropdownButtonFormField<String>(
                                initialValue: _currency,
                                decoration: const InputDecoration(
                                    labelText: 'Moneda'),
                                items: const [
                                  DropdownMenuItem(
                                      value: 'COP', child: Text('COP')),
                                  DropdownMenuItem(
                                      value: 'USD', child: Text('USD')),
                                  DropdownMenuItem(
                                      value: 'MXN', child: Text('MXN')),
                                ],
                                onChanged: (v) {
                                  if (v != null) {
                                    setState(() => _currency = v);
                                  }
                                },
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _stockController,
                          decoration:
                              const InputDecoration(labelText: 'Stock'),
                          keyboardType: TextInputType.number,
                          textInputAction: TextInputAction.next,
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _descriptionController,
                          decoration: const InputDecoration(
                            labelText: 'Descripcion',
                            alignLabelWithHint: true,
                          ),
                          maxLines: 3,
                          textInputAction: TextInputAction.done,
                        ),
                        const SizedBox(height: 16),
                        SwitchListTile(
                          title: const Text('Activo'),
                          subtitle: Text(
                            _isActive
                                ? 'El producto esta visible'
                                : 'El producto esta oculto',
                          ),
                          value: _isActive,
                          onChanged: (v) => setState(() => _isActive = v),
                          contentPadding: EdgeInsets.zero,
                        ),
                        const SizedBox(height: 24),
                        FilledButton(
                          onPressed: _isLoading ? null : _submit,
                          child: _isLoading
                              ? const SizedBox(
                                  height: 20,
                                  width: 20,
                                  child: CircularProgressIndicator(
                                      strokeWidth: 2),
                                )
                              : Text(
                                  _isEditing ? 'Actualizar' : 'Crear'),
                        ),
                      ],
                    ),
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
