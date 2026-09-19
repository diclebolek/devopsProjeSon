import 'dart:convert';

import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart' as p;

/// Geçici yerel veritabanı.
///
/// **Asıl backend: Supabase.** Bu katman yalnızca offline/demo ve
/// "beni hatırla" gibi geçici ihtiyaçlar içindir.
/// - Mobil/desktop: SQLite (`sqflite`)
/// - Web: SharedPreferences (tarayıcıda native SQLite yok)
class LocalDbService {
  LocalDbService._();
  static final LocalDbService instance = LocalDbService._();

  Database? _db;
  static const _prefsPrefix = 'orion_sqlite_fallback_';

  Future<void> init() async {
    if (kIsWeb) {
      await SharedPreferences.getInstance();
      return;
    }
    if (_db != null) return;
    final dbPath = await getDatabasesPath();
    _db = await openDatabase(
      p.join(dbPath, 'orion_local_temp.db'),
      version: 1,
      onCreate: (db, version) async {
        await db.execute('''
          CREATE TABLE IF NOT EXISTS kv (
            table_name TEXT NOT NULL,
            id TEXT NOT NULL,
            payload TEXT NOT NULL,
            updated_at INTEGER NOT NULL,
            PRIMARY KEY (table_name, id)
          )
        ''');
        await db.execute('''
          CREATE TABLE IF NOT EXISTS remembered_user (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL,
            role TEXT NOT NULL,
            remember INTEGER NOT NULL DEFAULT 1,
            updated_at INTEGER NOT NULL
          )
        ''');
      },
    );
  }

  Future<void> upsert(String table, String id, Map<String, dynamic> payload) async {
    await init();
    final json = jsonEncode(payload);
    final now = DateTime.now().millisecondsSinceEpoch;
    if (kIsWeb || _db == null) {
      final prefs = await SharedPreferences.getInstance();
      final idsKey = '$_prefsPrefix${table}__ids';
      final ids = prefs.getStringList(idsKey) ?? <String>[];
      if (!ids.contains(id)) {
        ids.add(id);
        await prefs.setStringList(idsKey, ids);
      }
      await prefs.setString('$_prefsPrefix${table}_$id', json);
      return;
    }
    await _db!.insert(
      'kv',
      {
        'table_name': table,
        'id': id,
        'payload': json,
        'updated_at': now,
      },
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  Future<Map<String, dynamic>?> get(String table, String id) async {
    await init();
    if (kIsWeb || _db == null) {
      final prefs = await SharedPreferences.getInstance();
      final raw = prefs.getString('$_prefsPrefix${table}_$id');
      if (raw == null) return null;
      return jsonDecode(raw) as Map<String, dynamic>;
    }
    final rows = await _db!.query(
      'kv',
      where: 'table_name = ? AND id = ?',
      whereArgs: [table, id],
      limit: 1,
    );
    if (rows.isEmpty) return null;
    return jsonDecode(rows.first['payload'] as String) as Map<String, dynamic>;
  }

  Future<void> saveRememberedUser({
    required String email,
    required String role,
    required bool remember,
  }) async {
    await init();
    final now = DateTime.now().millisecondsSinceEpoch;
    if (kIsWeb || _db == null) {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setBool('remember_me', remember);
      if (remember) {
        await prefs.setString('remembered_email', email);
        await prefs.setString('remembered_role', role);
      } else {
        await prefs.remove('remembered_email');
        await prefs.remove('remembered_role');
      }
      return;
    }
    await _db!.delete('remembered_user');
    if (remember) {
      await _db!.insert('remembered_user', {
        'email': email,
        'role': role,
        'remember': 1,
        'updated_at': now,
      });
    }
  }

  Future<({String email, String role, bool remember})?> loadRememberedUser() async {
    await init();
    if (kIsWeb || _db == null) {
      final prefs = await SharedPreferences.getInstance();
      final remember = prefs.getBool('remember_me') ?? false;
      if (!remember) return null;
      final email = prefs.getString('remembered_email') ?? '';
      final role = prefs.getString('remembered_role') ?? 'customer';
      if (email.isEmpty) return null;
      return (email: email, role: role, remember: true);
    }
    final rows = await _db!.query(
      'remembered_user',
      where: 'remember = 1',
      orderBy: 'updated_at DESC',
      limit: 1,
    );
    if (rows.isEmpty) return null;
    final row = rows.first;
    return (
      email: row['email'] as String,
      role: row['role'] as String,
      remember: true,
    );
  }

  String get backendLabel =>
      kIsWeb ? 'SharedPreferences (web geçici)' : 'SQLite (geçici yerel)';
}
