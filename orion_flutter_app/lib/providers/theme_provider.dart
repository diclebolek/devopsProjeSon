import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../constants/colors.dart';

class ThemeProvider extends ChangeNotifier {
  bool _isDarkMode = false;
  bool _isInitialized = false;

  bool get isDarkMode => _isDarkMode;
  bool get isInitialized => _isInitialized;

  ThemeProvider() {
    _loadThemePreference();
  }

  void toggleTheme() {
    _isDarkMode = !_isDarkMode;
    _saveThemePreference();
    notifyListeners();
  }

  void setTheme(bool isDark) {
    _isDarkMode = isDark;
    _saveThemePreference();
    notifyListeners();
  }

  void _saveThemePreference() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('isDarkMode', _isDarkMode);
  }

  void _loadThemePreference() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final savedTheme = prefs.getBool('isDarkMode');
      _isDarkMode = savedTheme ?? false;
      _isInitialized = true;
      notifyListeners();
    } catch (e) {
      _isDarkMode = false;
      _isInitialized = true;
      notifyListeners();
    }
  }

  ThemeData get darkTheme => ThemeData(
        brightness: Brightness.dark,
        primarySwatch: Colors.blue,
        scaffoldBackgroundColor: Colors.black,
        fontFamily: 'Cormorant',
        dialogTheme: DialogThemeData(
          backgroundColor: SiriusColors.dialogBackground(true),
          shape: SiriusColors.dialogShape,
          elevation: 10,
          shadowColor: SiriusColors.accent.withValues(alpha: 0.25),
        ),
        colorScheme: const ColorScheme.dark(
          primary: SiriusColors.accent,
          secondary: SiriusColors.accent,
        ),
      );

  ThemeData get lightTheme => ThemeData(
        brightness: Brightness.light,
        primarySwatch: Colors.blue,
        scaffoldBackgroundColor: const Color(0xFFE5E2DB),
        fontFamily: 'Cormorant',
        dialogTheme: DialogThemeData(
          backgroundColor: SiriusColors.dialogBackground(false),
          shape: SiriusColors.dialogShape,
          elevation: 10,
          shadowColor: SiriusColors.accent.withValues(alpha: 0.25),
        ),
        colorScheme: const ColorScheme.light(
          primary: SiriusColors.accent,
          secondary: SiriusColors.accent,
        ),
      );

  ThemeData get currentTheme => _isDarkMode ? darkTheme : lightTheme;
}
