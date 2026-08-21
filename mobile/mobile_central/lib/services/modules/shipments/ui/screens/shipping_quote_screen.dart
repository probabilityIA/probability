import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/utils/rate_pricing.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/shipment_provider.dart';
import '../widgets/rate_card.dart';

class ShippingQuoteScreen extends StatefulWidget {
  const ShippingQuoteScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<ShippingQuoteScreen> createState() => _ShippingQuoteScreenState();
}

class _ShippingQuoteScreenState extends State<ShippingQuoteScreen> {
  final _destinationController = TextEditingController();
  final _daneController = TextEditingController(text: '11001');
  final _weightController = TextEditingController(text: '1000');
  final _valueController = TextEditingController(text: '150000');

  OriginAddress? _origin;
  bool _cod = false;
  bool _insured = false;
  bool _quoting = false;
  String? _error;
  List<EnvioClickRate> _rates = const [];
  EnvioClickRate? _selected;

  RatePricingOptions get _options => RatePricingOptions(cod: _cod, insured: _insured);

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ShipmentProvider>().fetchOriginAddresses(businessId: widget.businessId);
    });
  }

  @override
  void dispose() {
    _destinationController.dispose();
    _daneController.dispose();
    _weightController.dispose();
    _valueController.dispose();
    super.dispose();
  }

  Future<void> _quote() async {
    final origin = _origin;
    if (origin == null) {
      setState(() => _error = 'Selecciona la bodega de origen');
      return;
    }
    if (_destinationController.text.trim().isEmpty) {
      setState(() => _error = 'Escribe la direccion de destino');
      return;
    }

    setState(() {
      _quoting = true;
      _error = null;
      _rates = const [];
      _selected = null;
    });

    final request = EnvioClickQuoteRequest(
      businessId: widget.businessId,
      description: 'Cotizacion desde la app',
      contentValue: double.tryParse(_valueController.text) ?? 0,
      codValue: _cod ? (double.tryParse(_valueController.text) ?? 0) : null,
      includeGuideCost: true,
      codPaymentMethod: 'cash',
      insurance: _insured,
      packages: [
        EnvioClickPackage(
          weight: double.tryParse(_weightController.text) ?? 1000,
          height: 20,
          width: 20,
          length: 20,
        ),
      ],
      origin: EnvioClickAddress(
        company: origin.company,
        firstName: origin.firstName,
        lastName: origin.lastName,
        email: origin.email,
        phone: origin.phone,
        address: origin.street,
        suburb: origin.suburb,
        daneCode: origin.cityDaneCode,
      ),
      destination: EnvioClickAddress(
        address: _destinationController.text.trim(),
        suburb: 'Centro',
        daneCode: _daneController.text.trim(),
      ),
    );

    final provider = context.read<ShipmentProvider>();
    final result = await provider.quoteShipment(request);

    if (!mounted) return;
    setState(() {
      _quoting = false;
      if (result == null) {
        _error = provider.error ?? 'No se pudo cotizar';
        return;
      }
      _rates = provider.quotes;
      if (_rates.isEmpty) _error = 'Ninguna transportadora cubre ese destino';
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Expanded(
          child: ListView(
            padding: AppSpacing.page,
            children: [
              const AppSectionHeader(
                title: 'Origen',
                subtitle: 'Bodega desde donde sale el paquete',
              ),
              _OriginPicker(
                selected: _origin,
                onSelected: (value) => setState(() => _origin = value),
              ),
              const SizedBox(height: 18),
              const AppSectionHeader(title: 'Destino'),
              TextField(
                controller: _destinationController,
                decoration: const InputDecoration(
                  hintText: 'Direccion de entrega',
                  prefixIcon: Icon(Icons.place_outlined, size: 20),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: _daneController,
                keyboardType: TextInputType.number,
                inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                decoration: const InputDecoration(
                  hintText: 'Codigo DANE de la ciudad',
                  prefixIcon: Icon(Icons.map_outlined, size: 20),
                ),
              ),
              const SizedBox(height: 18),
              const AppSectionHeader(title: 'Paquete'),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _weightController,
                      keyboardType: TextInputType.number,
                      inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                      decoration: const InputDecoration(
                        labelText: 'Peso (g)',
                        prefixIcon: Icon(Icons.scale_outlined, size: 20),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextField(
                      controller: _valueController,
                      keyboardType: TextInputType.number,
                      inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                      decoration: const InputDecoration(
                        labelText: 'Valor declarado',
                        prefixIcon: Icon(Icons.attach_money, size: 20),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 14),
              AppCard(
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 4),
                child: Column(
                  children: [
                    SwitchListTile(
                      contentPadding: EdgeInsets.zero,
                      value: _cod,
                      onChanged: (value) => setState(() => _cod = value),
                      title: Text('Contra entrega', style: Theme.of(context).textTheme.titleSmall),
                      subtitle: Text(
                        'Suma la comision de la transportadora al total',
                        style: Theme.of(context).textTheme.labelSmall,
                      ),
                    ),
                    const Divider(height: 1),
                    SwitchListTile(
                      contentPadding: EdgeInsets.zero,
                      value: _insured,
                      onChanged: (value) => setState(() => _insured = value),
                      title: Text('Asegurar envio', style: Theme.of(context).textTheme.titleSmall),
                      subtitle: Text(
                        'Agrega el seguro extra sobre el valor declarado',
                        style: Theme.of(context).textTheme.labelSmall,
                      ),
                    ),
                  ],
                ),
              ),
              if (_error != null) ...[
                const SizedBox(height: 14),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                  decoration: BoxDecoration(
                    color: AppColors.errorSoft,
                    borderRadius: AppRadius.mdAll,
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.error_outline, size: 18, color: AppColors.error),
                      const SizedBox(width: 10),
                      Expanded(
                        child: Text(
                          _error!,
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: const Color(0xFFB91C1C),
                                fontWeight: FontWeight.w500,
                              ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              const SizedBox(height: 16),
              FilledButton.icon(
                onPressed: _quoting ? null : _quote,
                style: FilledButton.styleFrom(minimumSize: const Size(0, 50)),
                icon: _quoting
                    ? const SizedBox(
                        height: 17,
                        width: 17,
                        child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                      )
                    : const Icon(Icons.calculate_outlined, size: 19),
                label: Text(_quoting ? 'Cotizando' : 'Cotizar envio'),
              ),
              if (_rates.isNotEmpty) ...[
                const SizedBox(height: 22),
                AppSectionHeader(
                  title: 'Tarifas disponibles',
                  subtitle: '${_rates.length} opciones para ese destino',
                ),
                ..._rates.map(
                  (rate) => Padding(
                    padding: const EdgeInsets.only(bottom: 10),
                    child: RateCard(
                      rate: rate,
                      options: _options,
                      selected: _selected?.idRate == rate.idRate,
                      onTap: () => setState(() => _selected = rate),
                    ),
                  ),
                ),
              ],
              const SizedBox(height: 24),
            ],
          ),
        ),
        if (_selected != null) _SelectionBar(rate: _selected!, options: _options),
      ],
    );
  }
}

class _SelectionBar extends StatelessWidget {
  const _SelectionBar({required this.rate, required this.options});

  final EnvioClickRate rate;
  final RatePricingOptions options;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final total = RatePricing.totalCost(rate, options);

    return Container(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 14),
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(top: AppBorders.hairline),
      ),
      child: SafeArea(
        top: false,
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(rate.carrier, style: theme.textTheme.titleSmall),
                  Text(AppFormat.money(total), style: theme.textTheme.titleLarge),
                ],
              ),
            ),
            FilledButton.icon(
              onPressed: () {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text('Generar guia con ${rate.carrier} por ${AppFormat.money(total)}'),
                  ),
                );
              },
              icon: const Icon(Icons.local_shipping_outlined, size: 18),
              label: const Text('Generar guia'),
            ),
          ],
        ),
      ),
    );
  }
}

class _OriginPicker extends StatelessWidget {
  const _OriginPicker({required this.selected, required this.onSelected});

  final OriginAddress? selected;
  final ValueChanged<OriginAddress> onSelected;

  @override
  Widget build(BuildContext context) {
    return Consumer<ShipmentProvider>(
      builder: (context, provider, _) {
        final addresses = provider.originAddresses;
        if (addresses.isEmpty) {
          return AppCard(
            child: Text(
              'No hay bodegas configuradas con datos de contacto',
              style: Theme.of(context).textTheme.bodySmall,
            ),
          );
        }

        return Column(
          children: addresses.map((address) {
            final active = selected?.id == address.id;
            return Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: AppCard(
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                borderColor: active ? AppColors.primary : AppColors.border,
                onTap: () => onSelected(address),
                child: Row(
                  children: [
                    Icon(
                      active ? Icons.radio_button_checked : Icons.radio_button_unchecked,
                      size: 19,
                      color: active ? AppColors.primary : AppColors.textDisabled,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(address.alias, style: Theme.of(context).textTheme.titleSmall),
                          const SizedBox(height: 2),
                          Text(
                            '${address.street}, ${address.city}',
                            style: Theme.of(context).textTheme.labelSmall,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            );
          }).toList(),
        );
      },
    );
  }
}
