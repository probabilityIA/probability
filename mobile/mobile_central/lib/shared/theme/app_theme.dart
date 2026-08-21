import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'app_colors.dart';
import 'app_tokens.dart';

class AppTheme {
  const AppTheme._();

  static const String fontFamily = 'Inter';

  static const ColorScheme _scheme = ColorScheme(
    brightness: Brightness.light,
    primary: AppColors.primary,
    onPrimary: Colors.white,
    primaryContainer: AppColors.primarySoft,
    onPrimaryContainer: AppColors.primaryDark,
    secondary: AppColors.accent,
    onSecondary: AppColors.textPrimary,
    secondaryContainer: AppColors.accentSoft,
    onSecondaryContainer: AppColors.textPrimary,
    tertiary: AppColors.secondary,
    onTertiary: Colors.white,
    error: AppColors.error,
    onError: Colors.white,
    errorContainer: AppColors.errorSoft,
    onErrorContainer: Color(0xFF7F1D1D),
    surface: AppColors.surface,
    onSurface: AppColors.textPrimary,
    surfaceContainerLowest: Colors.white,
    surfaceContainerLow: AppColors.background,
    surfaceContainer: AppColors.surfaceMuted,
    onSurfaceVariant: AppColors.textMuted,
    outline: AppColors.borderStrong,
    outlineVariant: AppColors.border,
    shadow: Color(0x1A000000),
    scrim: Color(0x66000000),
    inverseSurface: AppColors.textPrimary,
    onInverseSurface: Colors.white,
    inversePrimary: AppColors.primarySoft,
  );

  static TextTheme get _text {
    const base = TextStyle(fontFamily: fontFamily, color: AppColors.textPrimary);
    return TextTheme(
      displayLarge: base.copyWith(fontSize: 32, fontWeight: FontWeight.w700, letterSpacing: -0.6),
      displayMedium: base.copyWith(fontSize: 28, fontWeight: FontWeight.w700, letterSpacing: -0.5),
      displaySmall: base.copyWith(fontSize: 24, fontWeight: FontWeight.w700, letterSpacing: -0.4),
      headlineMedium: base.copyWith(fontSize: 22, fontWeight: FontWeight.w600, letterSpacing: -0.3),
      headlineSmall: base.copyWith(fontSize: 20, fontWeight: FontWeight.w600, letterSpacing: -0.2),
      titleLarge: base.copyWith(fontSize: 18, fontWeight: FontWeight.w600, letterSpacing: -0.2),
      titleMedium: base.copyWith(fontSize: 16, fontWeight: FontWeight.w600),
      titleSmall: base.copyWith(fontSize: 14, fontWeight: FontWeight.w600),
      bodyLarge: base.copyWith(fontSize: 16, fontWeight: FontWeight.w400, height: 1.45),
      bodyMedium: base.copyWith(fontSize: 14, fontWeight: FontWeight.w400, height: 1.45),
      bodySmall: base.copyWith(fontSize: 12, fontWeight: FontWeight.w400, height: 1.4, color: AppColors.textMuted),
      labelLarge: base.copyWith(fontSize: 14, fontWeight: FontWeight.w600, letterSpacing: 0.1),
      labelMedium: base.copyWith(fontSize: 12, fontWeight: FontWeight.w500, color: AppColors.textMuted),
      labelSmall: base.copyWith(fontSize: 11, fontWeight: FontWeight.w500, color: AppColors.textMuted, letterSpacing: 0.2),
    );
  }

  static ThemeData get light {
    final text = _text;
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      colorScheme: _scheme,
      fontFamily: fontFamily,
      textTheme: text,
      scaffoldBackgroundColor: AppColors.background,
      canvasColor: AppColors.surface,
      dividerColor: AppColors.border,
      splashFactory: InkSparkle.splashFactory,
      visualDensity: VisualDensity.standard,
      appBarTheme: AppBarTheme(
        backgroundColor: AppColors.surface,
        foregroundColor: AppColors.textPrimary,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: text.titleLarge,
        iconTheme: const IconThemeData(color: AppColors.textSecondary, size: 22),
        systemOverlayStyle: SystemUiOverlayStyle.dark,
        shape: const Border(bottom: AppBorders.hairline),
      ),
      dividerTheme: const DividerThemeData(
        color: AppColors.border,
        thickness: 1,
        space: 1,
      ),
      cardTheme: CardThemeData(
        color: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: AppRadius.lgAll,
          side: AppBorders.hairline,
        ),
      ),
      listTileTheme: ListTileThemeData(
        contentPadding: AppSpacing.listItem,
        iconColor: AppColors.textMuted,
        titleTextStyle: text.titleSmall,
        subtitleTextStyle: text.bodySmall,
        shape: RoundedRectangleBorder(borderRadius: AppRadius.mdAll),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.surface,
        isDense: true,
        contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
        hintStyle: text.bodyMedium?.copyWith(color: AppColors.textDisabled),
        labelStyle: text.bodyMedium?.copyWith(color: AppColors.textMuted),
        floatingLabelStyle: text.labelLarge?.copyWith(color: AppColors.primary),
        prefixIconColor: AppColors.textMuted,
        suffixIconColor: AppColors.textMuted,
        border: OutlineInputBorder(
          borderRadius: AppRadius.mdAll,
          borderSide: AppBorders.hairline,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: AppRadius.mdAll,
          borderSide: AppBorders.hairline,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: AppRadius.mdAll,
          borderSide: const BorderSide(color: AppColors.primary, width: 1.6),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: AppRadius.mdAll,
          borderSide: const BorderSide(color: AppColors.error, width: 1.2),
        ),
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: AppRadius.mdAll,
          borderSide: const BorderSide(color: AppColors.error, width: 1.6),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: AppColors.primary,
          foregroundColor: Colors.white,
          disabledBackgroundColor: AppColors.surfaceMuted,
          disabledForegroundColor: AppColors.textDisabled,
          minimumSize: const Size(0, 48),
          padding: const EdgeInsets.symmetric(horizontal: 20),
          textStyle: text.labelLarge,
          elevation: 0,
          shape: RoundedRectangleBorder(borderRadius: AppRadius.mdAll),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.textPrimary,
          backgroundColor: AppColors.surface,
          minimumSize: const Size(0, 48),
          padding: const EdgeInsets.symmetric(horizontal: 20),
          textStyle: text.labelLarge,
          side: const BorderSide(color: AppColors.borderStrong),
          shape: RoundedRectangleBorder(borderRadius: AppRadius.mdAll),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.primary,
          textStyle: text.labelLarge,
          minimumSize: const Size(0, 40),
          shape: RoundedRectangleBorder(borderRadius: AppRadius.smAll),
        ),
      ),
      iconButtonTheme: IconButtonThemeData(
        style: IconButton.styleFrom(
          foregroundColor: AppColors.textSecondary,
          shape: RoundedRectangleBorder(borderRadius: AppRadius.smAll),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: AppColors.surfaceMuted,
        selectedColor: AppColors.primarySoft,
        side: BorderSide.none,
        labelStyle: text.labelMedium!.copyWith(color: AppColors.textSecondary),
        secondaryLabelStyle: text.labelMedium!.copyWith(color: AppColors.primaryDark),
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        shape: RoundedRectangleBorder(borderRadius: AppRadius.pillAll),
        showCheckmark: false,
      ),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
        backgroundColor: AppColors.surface,
        selectedItemColor: AppColors.primary,
        unselectedItemColor: AppColors.textDisabled,
        selectedLabelStyle: text.labelSmall!.copyWith(color: AppColors.primary, fontWeight: FontWeight.w600),
        unselectedLabelStyle: text.labelSmall,
        type: BottomNavigationBarType.fixed,
        elevation: 0,
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        indicatorColor: AppColors.primarySoft,
        indicatorShape: RoundedRectangleBorder(borderRadius: AppRadius.pillAll),
        height: 66,
        elevation: 0,
        labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
        labelTextStyle: WidgetStateProperty.resolveWith((states) {
          final selected = states.contains(WidgetState.selected);
          return text.labelSmall!.copyWith(
            color: selected ? AppColors.primary : AppColors.textDisabled,
            fontWeight: selected ? FontWeight.w600 : FontWeight.w500,
          );
        }),
        iconTheme: WidgetStateProperty.resolveWith((states) {
          final selected = states.contains(WidgetState.selected);
          return IconThemeData(
            size: 22,
            color: selected ? AppColors.primary : AppColors.textDisabled,
          );
        }),
      ),
      drawerTheme: const DrawerThemeData(
        backgroundColor: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        width: 296,
      ),
      dialogTheme: DialogThemeData(
        backgroundColor: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        titleTextStyle: text.titleLarge,
        contentTextStyle: text.bodyMedium,
        shape: RoundedRectangleBorder(borderRadius: AppRadius.lgAll),
      ),
      bottomSheetTheme: BottomSheetThemeData(
        backgroundColor: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        showDragHandle: true,
        dragHandleColor: AppColors.borderStrong,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(AppRadius.xl)),
        ),
      ),
      snackBarTheme: SnackBarThemeData(
        backgroundColor: AppColors.textPrimary,
        contentTextStyle: text.bodyMedium!.copyWith(color: Colors.white),
        behavior: SnackBarBehavior.floating,
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: AppRadius.mdAll),
      ),
      tabBarTheme: TabBarThemeData(
        labelColor: AppColors.primary,
        unselectedLabelColor: AppColors.textMuted,
        labelStyle: text.titleSmall,
        unselectedLabelStyle: text.titleSmall!.copyWith(fontWeight: FontWeight.w500),
        indicatorColor: AppColors.primary,
        indicatorSize: TabBarIndicatorSize.label,
        dividerColor: AppColors.border,
        dividerHeight: 1,
        overlayColor: WidgetStatePropertyAll(AppColors.primarySoft.withValues(alpha: 0.4)),
      ),
      switchTheme: SwitchThemeData(
        thumbColor: WidgetStateProperty.resolveWith((s) =>
            s.contains(WidgetState.selected) ? Colors.white : Colors.white),
        trackColor: WidgetStateProperty.resolveWith((s) =>
            s.contains(WidgetState.selected) ? AppColors.primary : AppColors.borderStrong),
        trackOutlineColor: const WidgetStatePropertyAll(Colors.transparent),
      ),
      progressIndicatorTheme: const ProgressIndicatorThemeData(
        color: AppColors.primary,
        linearTrackColor: AppColors.surfaceMuted,
        circularTrackColor: AppColors.surfaceMuted,
      ),
      floatingActionButtonTheme: FloatingActionButtonThemeData(
        backgroundColor: AppColors.primary,
        foregroundColor: Colors.white,
        elevation: 2,
        shape: RoundedRectangleBorder(borderRadius: AppRadius.lgAll),
      ),
      popupMenuTheme: PopupMenuThemeData(
        color: AppColors.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 2,
        textStyle: text.bodyMedium,
        shape: RoundedRectangleBorder(
          borderRadius: AppRadius.mdAll,
          side: AppBorders.hairline,
        ),
      ),
    );
  }

  static ThemeData get dark => light;
}
