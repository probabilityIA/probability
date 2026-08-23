import 'package:flutter/material.dart';

class ChannelBrand {
  const ChannelBrand(this.color);

  final Color color;
}

const ChannelBrand _neutral = ChannelBrand(Color(0xFF9CA3AF));

const Map<String, ChannelBrand> _brands = <String, ChannelBrand>{
  'mercado libre': ChannelBrand(Color(0xFFD9B300)),
  'mercadolibre': ChannelBrand(Color(0xFFD9B300)),
  'meli': ChannelBrand(Color(0xFFD9B300)),
  'woocommerce': ChannelBrand(Color(0xFF7F54B3)),
  'woo': ChannelBrand(Color(0xFF7F54B3)),
  'siigo': ChannelBrand(Color(0xFF0EA5E9)),
  'shopify': ChannelBrand(Color(0xFF6B8E23)),
  'vtex': ChannelBrand(Color(0xFFF71963)),
  'jumpseller': ChannelBrand(Color(0xFFF97316)),
  'tiendanube': ChannelBrand(Color(0xFF2563EB)),
};

ChannelBrand channelBrand(String? name) {
  if (name == null || name.trim().isEmpty) return _neutral;
  final key = name.trim().toLowerCase();
  final direct = _brands[key];
  if (direct != null) return direct;
  for (final entry in _brands.entries) {
    if (key.contains(entry.key)) return entry.value;
  }
  return _neutral;
}
