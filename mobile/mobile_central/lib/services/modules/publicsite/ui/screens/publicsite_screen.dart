import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/publicsite_provider.dart';

class PublicSiteScreen extends StatefulWidget {
  const PublicSiteScreen({super.key});

  @override
  State<PublicSiteScreen> createState() => _PublicSiteScreenState();
}

class _PublicSiteScreenState extends State<PublicSiteScreen> {
  final TextEditingController _slugCtrl = TextEditingController();
  bool _hasSearched = false;

  @override
  void dispose() {
    _slugCtrl.dispose();
    super.dispose();
  }

  void _searchBusiness() {
    final slug = _slugCtrl.text.trim();
    if (slug.isEmpty) return;
    setState(() => _hasSearched = true);
    context.read<PublicSiteProvider>().fetchBusinessPage(slug);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Sitio Publico')),
      body: Consumer<PublicSiteProvider>(
        builder: (context, provider, _) {
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _slugCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Codigo del negocio',
                          hintText: 'ej: mi-tienda',
                          border: OutlineInputBorder(),
                          prefixIcon: Icon(Icons.search),
                        ),
                        onSubmitted: (_) => _searchBusiness(),
                      ),
                    ),
                    const SizedBox(width: 8),
                    FilledButton(
                      onPressed: _searchBusiness,
                      child: const Text('Buscar'),
                    ),
                  ],
                ),
                const SizedBox(height: 24),
                if (provider.isLoading)
                  const Center(
                      child: Padding(
                    padding: EdgeInsets.all(32),
                    child: CircularProgressIndicator(),
                  ))
                else if (provider.error != null)
                  Center(
                    child: Padding(
                      padding: const EdgeInsets.all(24),
                      child: Column(
                        children: [
                          const Icon(Icons.error_outline,
                              size: 48, color: Colors.red),
                          const SizedBox(height: 16),
                          Text(provider.error!,
                              textAlign: TextAlign.center,
                              style: const TextStyle(color: Colors.red)),
                        ],
                      ),
                    ),
                  )
                else if (provider.business != null)
                  _BusinessInfo(business: provider.business!)
                else if (_hasSearched)
                  Center(
                    child: Padding(
                      padding: const EdgeInsets.all(32),
                      child: Column(
                        children: [
                          Icon(Icons.search_off,
                              size: 48, color: Colors.grey.shade400),
                          const SizedBox(height: 16),
                          Text('Negocio no encontrado',
                              style: TextStyle(
                                  fontSize: 16, color: Colors.grey.shade600)),
                        ],
                      ),
                    ),
                  )
                else
                  Center(
                    child: Padding(
                      padding: const EdgeInsets.all(32),
                      child: Column(
                        children: [
                          Icon(Icons.web_outlined,
                              size: 64, color: Colors.grey.shade400),
                          const SizedBox(height: 16),
                          Text(
                              'Ingresa el codigo de un negocio para ver su pagina publica',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                  fontSize: 14, color: Colors.grey.shade600)),
                        ],
                      ),
                    ),
                  ),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _BusinessInfo extends StatelessWidget {
  final PublicBusiness business;
  const _BusinessInfo({required this.business});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                if (business.logoUrl.isNotEmpty)
                  ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: Image.network(
                      business.logoUrl,
                      width: 60,
                      height: 60,
                      fit: BoxFit.contain,
                      errorBuilder: (context, error, stackTrace) => CircleAvatar(
                        radius: 30,
                        backgroundColor: Colors.grey.shade200,
                        child: const Icon(Icons.business, size: 30),
                      ),
                    ),
                  )
                else
                  CircleAvatar(
                    radius: 30,
                    backgroundColor: Colors.grey.shade200,
                    child: const Icon(Icons.business, size: 30),
                  ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(business.name,
                          style: const TextStyle(
                              fontWeight: FontWeight.bold, fontSize: 18)),
                      Text(business.code,
                          style: TextStyle(
                              fontSize: 13, color: Colors.grey.shade600)),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        if (business.description.isNotEmpty) ...[
          const SizedBox(height: 12),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Descripcion',
                      style:
                          TextStyle(fontWeight: FontWeight.w600, fontSize: 14)),
                  const SizedBox(height: 8),
                  Text(business.description,
                      style: TextStyle(
                          fontSize: 14, color: Colors.grey.shade700)),
                ],
              ),
            ),
          ),
        ],
        if (business.featuredProducts.isNotEmpty) ...[
          const SizedBox(height: 16),
          const Text('Productos Destacados',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
          const SizedBox(height: 8),
          ...business.featuredProducts.take(5).map((p) => Card(
                margin: const EdgeInsets.only(bottom: 8),
                child: ListTile(
                  leading: p.imageUrl.isNotEmpty
                      ? ClipRRect(
                          borderRadius: BorderRadius.circular(4),
                          child: Image.network(
                            p.imageUrl,
                            width: 48,
                            height: 48,
                            fit: BoxFit.cover,
                            errorBuilder: (context, error, stackTrace) => Container(
                              width: 48,
                              height: 48,
                              color: Colors.grey.shade200,
                              child: const Icon(Icons.image_not_supported,
                                  size: 20),
                            ),
                          ),
                        )
                      : Container(
                          width: 48,
                          height: 48,
                          color: Colors.grey.shade200,
                          child:
                              const Icon(Icons.inventory_2_outlined, size: 20),
                        ),
                  title: Text(p.name,
                      style: const TextStyle(fontWeight: FontWeight.w600),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis),
                  subtitle: Text('\$${p.price.toStringAsFixed(0)} ${p.currency}'),
                ),
              )),
        ],
        if (business.websiteConfig != null) ...[
          const SizedBox(height: 16),
          const Text('Secciones del Sitio',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
          const SizedBox(height: 8),
          _SectionToggle(
              label: 'Hero', enabled: business.websiteConfig!.showHero),
          _SectionToggle(
              label: 'Acerca de', enabled: business.websiteConfig!.showAbout),
          _SectionToggle(
              label: 'Productos Destacados',
              enabled: business.websiteConfig!.showFeaturedProducts),
          _SectionToggle(
              label: 'Catalogo Completo',
              enabled: business.websiteConfig!.showFullCatalog),
          _SectionToggle(
              label: 'Testimonios',
              enabled: business.websiteConfig!.showTestimonials),
          _SectionToggle(
              label: 'Ubicacion',
              enabled: business.websiteConfig!.showLocation),
          _SectionToggle(
              label: 'Contacto',
              enabled: business.websiteConfig!.showContact),
        ],
      ],
    );
  }
}

class _SectionToggle extends StatelessWidget {
  final String label;
  final bool enabled;
  const _SectionToggle({required this.label, required this.enabled});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: ListTile(
        dense: true,
        title: Text(label, style: const TextStyle(fontSize: 14)),
        trailing: Icon(
          enabled ? Icons.check_circle : Icons.cancel_outlined,
          color: enabled ? Colors.green : Colors.grey.shade400,
          size: 20,
        ),
      ),
    );
  }
}
