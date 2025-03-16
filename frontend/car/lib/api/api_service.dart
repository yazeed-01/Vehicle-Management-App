import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  static const String baseUrl =
      'http://localhost:8080';
  static const storage = FlutterSecureStorage();

  Future<Map<String, dynamic>> login(
      String nationalNumber, String specificKey) async {
    final response = await http.post(
      Uri.parse('$baseUrl/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'national_number': nationalNumber,
        'specific_key': specificKey,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      await storage.write(key: 'jwt_token', value: data['token']);
      return data;
    } else {
      throw Exception('Login failed: ${response.body}');
    }
  }

  Future<String?> getToken() async {
    return await storage.read(key: 'jwt_token');
  }

  // Check if user is logged in
  Future<bool> isLoggedIn() async {
    final token = await getToken();
    return token != null;
  }

  Future<void> logout() async {
    await storage.delete(key: 'jwt_token');
  }

  Future<List<Map<String, dynamic>>> getUserVehicles() async {
    final token = await getToken();
    final response = await http.get(
      Uri.parse('$baseUrl/api/vehicles'),
      headers: {
        'Authorization': token ?? '',
      },
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return List<Map<String, dynamic>>.from(data['vehicles']);
    } else {
      throw Exception('Failed to fetch vehicles: ${response.body}');
    }
  }

  Future<Map<String, dynamic>> getUserInfo() async {
    final token = await getToken();
    final response = await http.get(
      Uri.parse('$baseUrl/api/user'),
      headers: {
        'Authorization': token ?? '',
      },
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return Map<String, dynamic>.from(data['user']);
    } else {
      throw Exception('Failed to fetch user info: ${response.body}');
    }
  }

  Future<Map<String, dynamic>> searchByPlate(String carPlate) async {
    final token = await getToken();
    final response = await http.post(
      Uri.parse('$baseUrl/api/search-plate'),
      headers: {
        'Authorization': token ?? '',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'car_plate': carPlate,
      }),
    );

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else {
      throw Exception('Failed to search by plate: ${response.body}');
    }
  }

  Future<void> updateUserInfo(String email, String phone) async {
    final token = await getToken();
    final response = await http.put(
      Uri.parse('$baseUrl/api/user'),
      headers: {
        'Authorization': token ?? '',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'email': email,
        'phone_number': phone,
      }),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to update user info: ${response.body}');
    }
  }
}
