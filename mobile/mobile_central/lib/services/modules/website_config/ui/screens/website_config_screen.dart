import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/website_config_provider.dart';

class WebsiteConfigScreen extends StatefulWidget {
  final int? businessId;

  const WebsiteConfigScreen({super.key, this.businessId});

  @override
  State<WebsiteConfigScreen> createState() => _WebsiteConfigScreenState();
}

class _WebsiteConfigScreenState extends State<WebsiteConfigScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadConfig();
    });
  }

  void _loadConfig() {
    context.read<WebsiteConfigProvider>().fetchConfig();
  }

  Future<void> _toggleSection(String field, bool value) async {
    final provider = context.read<WebsiteConfigProvider>();
    final dto = UpdateWebsiteConfigDTO(
      showHero: field == 'hero' ? value : null,
      showAbout: field == 'about' ? value : null,
      showFeaturedProducts: field == 'featured_products' ? value : null,
      showFullCatalog: field == 'full_catalog' ? value : null,
      showTestimonials: field == 'testimonials' ? value : null,
      showLocation: field == 'location' ? value : null,
      showContact: field == 'contact' ? value : null,
      showSocialMedia: field == 'social_media' ? value : null,
      showWhatsapp: field == 'whatsapp' ? value : null,
    );
    final success = await provider.updateConfig(dto);
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(success ? 'Actualizado' : 'Error al actualizar'),
          duration: const Duration(seconds: 1),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Configuracion del Sitio Web')),
      body: Consumer<WebsiteConfigProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (provider.error != null) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.error_outline,
                        size: 48, color: Colors.red),
                    const SizedBox(height: 16),
                    Text(provider.error!,
                        textAlign: TextAlign.center,
                        style: const TextStyle(color: Colors.red)),
                    const SizedBox(height: 16),
                    FilledButton.icon(
                      onPressed: _loadConfig,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          final config = provider.config;
          if (config == null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.web_outlined,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('Sin configuracion',
                      style:
                          TextStyle(fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadConfig(),
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      children: [
                        Icon(Icons.palette_outlined,
                            color: Theme.of(context).colorScheme.primary),
                        const SizedBox(width: 12),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('Template',
                                style: TextStyle(
                                    fontWeight: FontWeight.w600,
                                    fontSize: 14)),
                            Text(
                                config.template.isNotEmpty
                                    ? config.template
                                    : 'Default',
                                style: TextStyle(
                                    fontSize: 13,
                                    color: Colors.grey.shade600)),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                const Text('Secciones del Sitio',
                    style:
                        TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
                const SizedBox(height: 8),
                _SectionSwitch(
                  icon: Icons.view_carousel_outlined,
                  label: 'Hero',
                  description: 'Banner principal',
                  value: config.showHero,
                  onChanged: (v) => _toggleSection('hero', v),
                ),
                _SectionSwitch(
                  icon: Icons.info_outline,
                  label: 'Acerca de',
                  description: 'Seccion informativa',
                  value: config.showAbout,
                  onChanged: (v) => _toggleSection('about', v),
                ),
                _SectionSwitch(
                  icon: Icons.star_outline,
                  label: 'Productos Destacados',
                  description: 'Productos destacados',
                  value: config.showFeaturedProducts,
                  onChanged: (v) =>
                      _toggleSection('featured_products', v),
                ),
                _SectionSwitch(
                  icon: Icons.grid_view,
                  label: 'Catalogo Completo',
                  description: 'Todos los productos',
                  value: config.showFullCatalog,
                  onChanged: (v) => _toggleSection('full_catalog', v),
                ),
                _SectionSwitch(
                  icon: Icons.format_quote_outlined,
                  label: 'Testimonios',
                  description: 'Opiniones de clientes',
                  value: config.showTestimonials,
                  onChanged: (v) => _toggleSection('testimonials', v),
                ),
                _SectionSwitch(
                  icon: Icons.location_on_outlined,
                  label: 'Ubicacion',
                  description: 'Mapa y direccion',
                  value: config.showLocation,
                  onChanged: (v) => _toggleSection('location', v),
                ),
                _SectionSwitch(
                  icon: Icons.contact_mail_outlined,
                  label: 'Contacto',
                  description: 'Formulario de contacto',
                  value: config.showContact,
                  onChanged: (v) => _toggleSection('contact', v),
                ),
                _SectionSwitch(
                  icon: Icons.share_outlined,
                  label: 'Redes Sociales',
                  description: 'Links a redes',
                  value: config.showSocialMedia,
                  onChanged: (v) => _toggleSection('social_media', v),
                ),
                _SectionSwitch(
                  icon: Icons.chat_outlined,
                  label: 'WhatsApp',
                  description: 'Boton flotante',
                  value: config.showWhatsapp,
                  onChanged: (v) => _toggleSection('whatsapp', v),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _SectionSwitch extends StatelessWidget {
  final IconData icon;
  final String label;
  final String description;
  final bool value;
  final ValueChanged<bool> onChanged;

  const _SectionSwitch({
    required this.icon,
    required this.label,
    required this.description,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: SwitchListTile(
        secondary: Icon(icon,
            color: value
                ? Theme.of(context).colorScheme.primary
                : Colors.grey.shade400),
        title: Text(label, style: const TextStyle(fontWeight: FontWeight.w500)),
        subtitle: Text(description,
            style: TextStyle(fontSize: 12, color: Colors.grey.shade600)),
        value: value,
        onChanged: onChanged,
      ),
    );
  }
}
