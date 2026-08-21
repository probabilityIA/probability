import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/website_config_provider.dart';

class WebsiteConfigScreen extends StatefulWidget {
  const WebsiteConfigScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<WebsiteConfigScreen> createState() => _WebsiteConfigScreenState();
}

class _WebsiteConfigScreenState extends State<WebsiteConfigScreen> {
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  void _refresh() {
    context.read<WebsiteConfigProvider>().fetchConfig(businessId: widget.businessId);
  }

  Future<void> _toggle(WebsiteConfigData config, String section, bool value) async {
    setState(() => _saving = true);

    final dto = UpdateWebsiteConfigDTO(
      showHero: section == 'hero' ? value : config.showHero,
      showAbout: section == 'about' ? value : config.showAbout,
      showFeaturedProducts: section == 'featured' ? value : config.showFeaturedProducts,
      showFullCatalog: section == 'catalog' ? value : config.showFullCatalog,
      showTestimonials: section == 'testimonials' ? value : config.showTestimonials,
      showLocation: section == 'location' ? value : config.showLocation,
      showContact: section == 'contact' ? value : config.showContact,
      showSocialMedia: section == 'social' ? value : config.showSocialMedia,
      showWhatsapp: section == 'whatsapp' ? value : config.showWhatsapp,
    );

    final provider = context.read<WebsiteConfigProvider>();
    final ok = await provider.updateConfig(dto, businessId: widget.businessId);

    if (!mounted) return;
    setState(() => _saving = false);
    if (!ok) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(provider.error ?? 'No se pudo guardar')),
      );
    }
    _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<WebsiteConfigProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.config == null) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.config == null) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }

        final config = provider.config;
        if (config == null) {
          return AppEmptyState(
            icon: Icons.web_outlined,
            title: 'Sin configuracion',
            message: 'Aun no hay un sitio configurado para este negocio.',
            actionLabel: 'Actualizar',
            onAction: _refresh,
          );
        }

        return RefreshIndicator(
          onRefresh: () async => _refresh(),
          color: AppColors.primary,
          child: ListView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: AppSpacing.page,
            children: [
              AppCard(
                child: Row(
                  children: [
                    Container(
                      width: 44,
                      height: 44,
                      alignment: Alignment.center,
                      decoration: BoxDecoration(
                        color: AppColors.primarySoft,
                        borderRadius: AppRadius.mdAll,
                      ),
                      child: const Icon(Icons.palette_outlined, size: 21, color: AppColors.primary),
                    ),
                    const SizedBox(width: 13),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Plantilla', style: Theme.of(context).textTheme.labelSmall),
                          const SizedBox(height: 2),
                          Text(
                            config.template.isEmpty ? 'Sin plantilla' : config.template,
                            style: Theme.of(context).textTheme.titleMedium,
                          ),
                        ],
                      ),
                    ),
                    if (_saving)
                      const SizedBox(
                        height: 18,
                        width: 18,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 18),
              const AppSectionHeader(
                title: 'Secciones del sitio',
                subtitle: 'Lo que ven tus clientes al entrar',
              ),
              AppCard(
                padding: const EdgeInsets.symmetric(horizontal: 14),
                child: Column(
                  children: [
                    _SectionToggle(
                      icon: Icons.wallpaper_outlined,
                      label: 'Portada',
                      detail: _stringFrom(config.heroContent, 'title'),
                      value: config.showHero,
                      onChanged: (v) => _toggle(config, 'hero', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.star_outline_rounded,
                      label: 'Productos destacados',
                      value: config.showFeaturedProducts,
                      onChanged: (v) => _toggle(config, 'featured', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.grid_view_outlined,
                      label: 'Catalogo completo',
                      value: config.showFullCatalog,
                      onChanged: (v) => _toggle(config, 'catalog', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.info_outline_rounded,
                      label: 'Sobre nosotros',
                      detail: _stringFrom(config.aboutContent, 'text'),
                      value: config.showAbout,
                      onChanged: (v) => _toggle(config, 'about', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.format_quote_rounded,
                      label: 'Testimonios',
                      value: config.showTestimonials,
                      onChanged: (v) => _toggle(config, 'testimonials', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.place_outlined,
                      label: 'Ubicacion',
                      detail: _stringFrom(config.locationContent, 'address'),
                      value: config.showLocation,
                      onChanged: (v) => _toggle(config, 'location', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.mail_outline_rounded,
                      label: 'Contacto',
                      detail: _stringFrom(config.contactContent, 'email'),
                      value: config.showContact,
                      onChanged: (v) => _toggle(config, 'contact', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.share_outlined,
                      label: 'Redes sociales',
                      detail: _stringFrom(config.socialMediaContent, 'instagram'),
                      value: config.showSocialMedia,
                      onChanged: (v) => _toggle(config, 'social', v),
                    ),
                    const Divider(height: 1),
                    _SectionToggle(
                      icon: Icons.chat_outlined,
                      label: 'Boton de WhatsApp',
                      detail: _stringFrom(config.whatsappContent, 'phone'),
                      value: config.showWhatsapp,
                      onChanged: (v) => _toggle(config, 'whatsapp', v),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 26),
            ],
          ),
        );
      },
    );
  }

  String? _stringFrom(Map<String, dynamic>? content, String key) {
    final value = content?[key];
    if (value == null) return null;
    final text = value.toString().trim();
    return text.isEmpty ? null : text;
  }
}

class _SectionToggle extends StatelessWidget {
  const _SectionToggle({
    required this.icon,
    required this.label,
    required this.value,
    required this.onChanged,
    this.detail,
  });

  final IconData icon;
  final String label;
  final bool value;
  final ValueChanged<bool> onChanged;
  final String? detail;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 5),
      child: Row(
        children: [
          Icon(icon, size: 19, color: value ? AppColors.primary : AppColors.textDisabled),
          const SizedBox(width: 13),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(label, style: theme.textTheme.titleSmall),
                if (detail != null) ...[
                  const SizedBox(height: 2),
                  Text(
                    detail!,
                    style: theme.textTheme.labelSmall,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ],
            ),
          ),
          Switch(value: value, onChanged: onChanged),
        ],
      ),
    );
  }
}
