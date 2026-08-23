class AppFormat {
  const AppFormat._();

  static String money(num? value, {String symbol = '\$', int decimals = 0}) {
    if (value == null) return '$symbol 0';
    final negative = value < 0;
    final fixed = value.abs().toStringAsFixed(decimals);
    final parts = fixed.split('.');
    final grouped = _group(parts[0]);
    final tail = parts.length > 1 ? ',${parts[1]}' : '';
    return '${negative ? '-' : ''}$symbol $grouped$tail';
  }

  static String number(num? value, {int decimals = 0}) {
    if (value == null) return '0';
    final fixed = value.abs().toStringAsFixed(decimals);
    final parts = fixed.split('.');
    final grouped = _group(parts[0]);
    final tail = parts.length > 1 ? ',${parts[1]}' : '';
    return '${value < 0 ? '-' : ''}$grouped$tail';
  }

  static String compact(num? value) {
    if (value == null) return '0';
    final abs = value.abs();
    if (abs >= 1000000000) return '${(value / 1000000000).toStringAsFixed(1)}B';
    if (abs >= 1000000) return '${(value / 1000000).toStringAsFixed(1)}M';
    if (abs >= 1000) return '${(value / 1000).toStringAsFixed(1)}K';
    return number(value);
  }

  static String _group(String digits) {
    final buffer = StringBuffer();
    for (var i = 0; i < digits.length; i++) {
      if (i > 0 && (digits.length - i) % 3 == 0) buffer.write('.');
      buffer.write(digits[i]);
    }
    return buffer.toString();
  }

  static const List<String> _months = [
    'ene', 'feb', 'mar', 'abr', 'may', 'jun',
    'jul', 'ago', 'sep', 'oct', 'nov', 'dic',
  ];

  static String date(DateTime? value) {
    if (value == null) return '-';
    final local = value.toLocal();
    return '${_two(local.day)} ${_months[local.month - 1]} ${local.year}';
  }

  static String dateTime(DateTime? value) {
    if (value == null) return '-';
    final local = value.toLocal();
    return '${date(local)} ${_two(local.hour)}:${_two(local.minute)}';
  }

  static String relative(DateTime? value) {
    if (value == null) return '-';
    final diff = DateTime.now().difference(value.toLocal());
    if (diff.inSeconds < 60) return 'hace un momento';
    if (diff.inMinutes < 60) return 'hace ${diff.inMinutes} min';
    if (diff.inHours < 24) return 'hace ${diff.inHours} h';
    if (diff.inDays < 7) return 'hace ${diff.inDays} d';
    return date(value);
  }

  static DateTime? parseDate(dynamic raw) {
    if (raw == null) return null;
    if (raw is DateTime) return raw;
    return DateTime.tryParse(raw.toString());
  }

  static String initials(String? name) {
    final clean = (name ?? '').trim();
    if (clean.isEmpty) return '?';
    final parts = clean.split(RegExp(r'\s+'));
    if (parts.length == 1) return parts.first.substring(0, 1).toUpperCase();
    return (parts[0].substring(0, 1) + parts[1].substring(0, 1)).toUpperCase();
  }

  static String _two(int value) => value.toString().padLeft(2, '0');
}
