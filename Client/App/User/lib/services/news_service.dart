import 'package:caoky/services/auth_service.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

import 'package:caoky/models/announcement_model.dart';
import 'package:shared_preferences/shared_preferences.dart'; // Đảm bảo đường dẫn đúng

class NewsService {
  // !!! QUAN TRỌNG: Thay đổi địa chỉ này thành URL của backend Go của bạn
  static const String _baseUrl =
      'http://57.155.76.74/api/v1'; // 10.0.2.2 là localhost cho máy ảo Android
  final AuthRepository _authRepository = AuthRepository();
  Future<List<Announcement>> getAnnouncements({
    int limit = 5,
    int offset = 0,
  }) async {
    final prefs = await SharedPreferences.getInstance();
    final token = _authRepository.getValidAccessToken();

    // Assuming you also store the user ID in shared preferences upon login
    final String? userId = prefs.getString('userId');

    if (token == null || userId == null) {
      throw Exception('Authentication token or User ID not found.');
    }

    final response = await http.get(
      Uri.parse('$_baseUrl/news?limit=$limit&offset=$offset'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      // API trả về mảng rỗng [] nếu không có tin tức
      if (response.body.trim() == '[]' || response.body.isEmpty) {
        return [];
      }
      // UTF-8 decode để hiển thị tiếng Việt đúng
      final String responseBody = utf8.decode(response.bodyBytes);
      return announcementFromJson(responseBody);
    } else {
      // Xử lý lỗi
      throw Exception(
        'Failed to load announcements. Status code: ${response.statusCode}',
      );
    }
  }
}
