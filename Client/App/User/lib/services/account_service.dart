import 'dart:convert';
import 'dart:io';
import 'package:caoky/models/transaction_history_model.dart';

import '../../models/account_model.dart'; // Replace 'package_name' with your actual package name
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:caoky/services/auth_service.dart';

class AccountService {
  // Replace with your actual API base URL
  final AuthRepository _authRepository = AuthRepository();

  static const String _baseUrl = "http://57.155.76.74/api/v1";

  Future<Account> getMyAccount() async {
    final prefs = await SharedPreferences.getInstance();
    final token = await _authRepository.getValidAccessToken();

    // Assuming you also store the user ID in shared preferences upon login
    final String? userId = prefs.getString('userId');

    if (token == null || userId == null) {
      throw Exception('Authentication token or User ID not found.');
    }

    final response = await http.get(
      Uri.parse('$_baseUrl/accounts/me'),
      headers: {
        HttpHeaders.contentTypeHeader: 'application/json',
        HttpHeaders.authorizationHeader: 'Bearer $token',
        'X-User-ID': userId, // Custom header for your Go backend
      },
    );

    if (response.statusCode == 200) {
      return Account.fromJson(jsonDecode(response.body));
    } else if (response.statusCode == 404) {
      // Use a specific exception for "not found" to handle it in the UI
      throw AccountNotFoundException('Account not found for the user.');
    } else {
      // Handle other potential errors (500, 401, etc.)
      throw Exception(
        'Failed to load account. Status code: ${response.statusCode}',
      );
    }
  }

  Future<Account> createAccount({
    required String ownerName,
    required int balance,
    required String currency,
  }) async {
    final prefs = await SharedPreferences.getInstance();

    final token = await _authRepository.getValidAccessToken();
    final String? userId = prefs.getString('userId');
    print('🔑 Token: $token');
    if (token == null || userId == null) {
      throw Exception('Authentication token or User ID not found.');
    }
    final response = await http.post(
      Uri.parse('$_baseUrl/accounts'),
      headers: {
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      },
      body: jsonEncode({
        'owner_name': userId,
        'balance': balance,
        'currency': currency,
      }),
    );

    final responseBody = jsonDecode(utf8.decode(response.bodyBytes));

    if (response.statusCode == 201) {
      // Created
      return Account.fromJson(responseBody);
    } else {
      // Ném lỗi với thông báo từ server
      throw Exception(responseBody['message'] ?? 'Không thể tạo tài khoản');
    }
  }

  Future<List<TransactionHistory>> getTransactionHistory({
    required int pageId,
    required int pageSize,
  }) async {
    final token = _authRepository.getValidAccessToken();

    if (token == null) {
      throw Exception('User not authenticated');
    }

    final uri = Uri.parse('$_baseUrl/accounts/history').replace(
      queryParameters: {
        'page_id': pageId.toString(),
        'page_size': pageSize.toString(),
      },
    );

    try {
      final response = await http.get(
        uri,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final List<dynamic> body = jsonDecode(utf8.decode(response.bodyBytes));
        if (body is List) {
          return body
              .map(
                (dynamic item) =>
                    TransactionHistory.fromJson(item as Map<String, dynamic>),
              )
              .toList();
        }
        throw Exception('Invalid response format');
      } else {
        // Xử lý các lỗi HTTP khác
        throw Exception(
          'Failed to load transaction history: ${response.statusCode}',
        );
      }
    } catch (e) {
      // Xử lý lỗi mạng hoặc parsing
      throw Exception('An error occurred while fetching history: $e');
    }
  }

  // Giả sử bạn có hàm helper này để lấy token
  // Future<String?> _getAuthToken() async {
  //   final prefs = await SharedPreferences.getInstance();
  //   return prefs.getString('accessToken');
  // }
}

// Custom exception for clarity
class AccountNotFoundException implements Exception {
  final String message;
  AccountNotFoundException(this.message);
}
