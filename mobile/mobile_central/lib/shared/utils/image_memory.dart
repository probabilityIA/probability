import 'package:flutter/material.dart';

class ImageMemory {
  const ImageMemory._();

  static const int lowEndCacheBytes = 40 << 20;
  static const int lowEndCacheCount = 120;
  static const double maxDecodeRatio = 3;

  static int decodePixels(BuildContext context, double logicalSize) {
    final ratio = MediaQuery.devicePixelRatioOf(context);
    final effective = ratio.clamp(1.0, maxDecodeRatio);
    return (logicalSize * effective).round();
  }

  static void applyLowEndBudget() {
    final cache = PaintingBinding.instance.imageCache;
    cache.maximumSizeBytes = lowEndCacheBytes;
    cache.maximumSize = lowEndCacheCount;
  }
}
