class BrandAssets {
  const BrandAssets._();

  static const String mediaBaseUrl =
      'https://probability-media-assets.s3.us-east-1.amazonaws.com';

  static String? mediaUrl(String? value) {
    final raw = (value ?? '').trim();
    if (raw.isEmpty) return null;
    if (raw.startsWith('http://') || raw.startsWith('https://')) return raw;
    final key = raw.startsWith('/') ? raw.substring(1) : raw;
    return '$mediaBaseUrl/$key';
  }

  static const Map<String, String> carrierLogos = {
    'SERVIENTREGA': 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_servientrega.png',
    'COORDINADORA': 'https://probability-media-assets.s3.us-east-1.amazonaws.com/public/carriers/imagen_coordinadora.png',
    'DHLEXPRESS': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
    'DHL': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
    'FEDEX': 'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9d/FedEx_Express.svg/960px-FedEx_Express.svg.png',
    'INTERRAPIDISIMO': '$mediaBaseUrl/carriers/interrapidisimo.jpg',
    '472LOGISTICA': 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTnDF0ozRHf3s5BPqLsr7Vg-X8JRzECvFvwBQ&s',
    'SPEED': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
    'SPEEDCARGO': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
    'ENVIA': 'https://images.seeklogo.com/logo-png/31/1/envia-mensajeria-logo-png_seeklogo-311137.png',
    'PIBOX': 'https://play-lh.googleusercontent.com/r_zPLkaHZK4Odu1yp6dqIdUnVAmIiLc3s18F9gUFqcz8IyHqCb_aGHP4iJSesXxnUyU',
    'TCC': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
    'TRANSPORTADORADECARACOLOMBIA': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
    '99MINUTOS': 'https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Logo-99minutos.svg/3840px-Logo-99minutos.svg.png',
    'DEPRISA': 'https://www.specialcolombia.com/wp-content/uploads/2023/05/Logo_azul_concepto_azul-deprisa.png',
    'ENVIOCLICK': 'https://www.envioclickpro.com.co/assets/images/envioclick-logo.png',
    'ENVIAME': 'https://enviame.io/wp-content/uploads/2021/01/logo-enviame.svg',
    'MIPAQUETE': 'https://mipaquete.com/wp-content/uploads/2021/03/mipaquete-logo.png',
    'SHIPIT': '$mediaBaseUrl/integration-types/1785346702_shipit.png',
    'RAPPI': 'https://static.vecteezy.com/system/resources/previews/067/941/720/non_2x/rappi-logo-rounded-hd-free-png.png',
  };


  static const Map<String, String> integrationLogos = {
    'SHOPIFY': 'integration-types/1765744750_shopify.png',
    'WHATSAPP': 'integration-types/1765744972_whatsap.png',
    'WHASTAP': 'integration-types/1765744972_whatsap.png',
    'MERCADOLIBRE': 'integration-types/1771905467_mercado-libre-logo.png',
    'MELI': 'integration-types/1771905467_mercado-libre-logo.png',
    'WOOCOMMERCE': 'integration-types/1765745053_woocomerce.webp',
    'WOOCORMERCE': 'integration-types/1765745053_woocomerce.webp',
    'SOFTPYMES': 'integration-types/1769929713_sofpymes.png',
    'FACTUS': 'integration-types/1771736455_factus.png',
    'SIIGO': 'integration-types/1771738597_siigo.png',
    'ALEGRA': 'integration-types/1771738659_alegra.png',
    'WORLDOFFICE': 'integration-types/1771738901_worl-office.png',
    'HELISA': 'integration-types/1771739081_logo-helisa.png',
    'ENVIOCLICK': 'integration-types/1771901154_envioclik.png',
    'ENVIAME': 'integration-types/1771905179_enviame.png',
    'MIPAQUETE': 'integration-types/1771905337_mipaquete.png',
    'VTEX': 'integration-types/1771905400_vtex.png',
    'TIENDANUBE': 'integration-types/1771905372_tiendanube.png',
    'MAGENTO': 'integration-types/1771905298_magneto.png',
    'AMAZON': 'integration-types/1771905134_amazon.png',
    'FALABELLA': 'integration-types/1771905254_falabella.png',
    'EXITO': 'integration-types/1771905220_exito.png',
    'NEQUI': 'integration-types/1772147410_nequi.png',
    'BOLD': 'integration-types/1777354669_bold.png',
    'BOLDPAY': 'integration-types/1777354669_bold.png',
    'WOMPI': 'integration-types/1772147507_logowompi.png',
    'STRIPE': 'integration-types/1772147475_stripe.png',
    'PAYU': 'integration-types/1772147440_payu.png',
    'EPAYCO': 'integration-types/1772147322_epayco.png',
    'MELIPAGO': 'integration-types/1772147377_mercado-pago.png',
    'MERCADOPAGO': 'integration-types/1772147377_mercado-pago.png',
    'JUMPSELLER': 'integration-types/1784241131_jumpseller.png',
    'SHIPIT': 'integration-types/1785346702_shipit.png',
    'TIKTOK': 'integration-types/tiktok.png',
    'BANCOLOMBIAQR': 'integration-types/bancolombia.png',
  };

  static String? integrationLogo(String? name) {
    final key = _normalize(name);
    if (key.isEmpty) return null;
    final match = integrationLogos[key];
    return match == null ? null : mediaUrl(match);
  }

  static String _normalize(String? value) {
    final raw = (value ?? '').trim();
    if (raw.isEmpty) return '';
    return _stripDiacritics(raw.split(' - ').first)
        .toUpperCase()
        .replaceAll(RegExp(r'[\s\-_\.]'), '');
  }

  static String? carrierLogo(String? carrier) {
    final raw = (carrier ?? '').trim();
    if (raw.isEmpty) return null;
    final key = _normalize(raw);
    return carrierLogos[key];
  }

  static const Map<String, String> _diacritics = {
    '\u00E1': 'a', '\u00E9': 'e', '\u00ED': 'i', '\u00F3': 'o', '\u00FA': 'u',
    '\u00C1': 'A', '\u00C9': 'E', '\u00CD': 'I', '\u00D3': 'O', '\u00DA': 'U',
    '\u00F1': 'n', '\u00D1': 'N', '\u00FC': 'u', '\u00DC': 'U',
  };

  static String _stripDiacritics(String value) {
    final buffer = StringBuffer();
    for (final char in value.split('')) {
      buffer.write(_diacritics[char] ?? char);
    }
    return buffer.toString();
  }
}
