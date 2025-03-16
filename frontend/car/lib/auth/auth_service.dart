import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthService {
  static const _storage = FlutterSecureStorage();
  static const _keyAuthToken = 'auth_token';

  Future<void> login(String nationalNumber) async {
    await _storage.write(key: _keyAuthToken, value: nationalNumber);
  }

  Future<bool> isLoggedIn() async {
    final token = await _storage.read(key: _keyAuthToken);
    return token != null;
  }

  // Logout
  Future<void> logout() async {
    await _storage.delete(key: _keyAuthToken);
  }

  Future<String?> getNationalNumber() async {
    return await _storage.read(key: _keyAuthToken);
  }
}