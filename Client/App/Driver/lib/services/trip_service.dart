// lib/services/trip_service.dart

import 'package:dio/dio.dart';
import 'package:taixe/models/trip.dart';
import 'package:taixe/services/auth_service.dart'; // Import AuthRepository để lấy token

class TripService {
  final Dio _dio = Dio();
  final AuthRepository _authRepository = AuthRepository();
  // Thay thế bằng base URL của bạn
  final String _baseUrl = 'http://57.155.76.74/api/v1'; 

  Future<List<Trip>> getDriverTrips() async {
    try {
      // Lấy access token hợp lệ từ AuthRepository
      final accessToken = await _authRepository.getValidAccessToken();
      if (accessToken == null) {
        throw Exception('Người dùng chưa đăng nhập hoặc phiên đã hết hạn.');
      }

      final response = await _dio.get(
        '$_baseUrl/trips/driver',
        options: Options(
          headers: {
            'Authorization': 'Bearer $accessToken',
          },
        ),
      );

      if (response.statusCode == 200 && response.data['code'] == 200) {
        // Lấy danh sách dữ liệu từ key 'data'
        List<dynamic> tripData = response.data['data'];
        // Chuyển đổi mỗi phần tử JSON trong danh sách thành một đối tượng Trip
        return tripData.map((json) => Trip.fromJson(json)).toList();
      } else {
        // Ném ra lỗi nếu API trả về mã lỗi
        throw Exception(response.data['message'] ?? 'Lấy dữ liệu chuyến đi thất bại.');
      }
    } on DioException catch (e) {
      // Xử lý các lỗi liên quan đến Dio (mạng, timeout, v.v.)
      print('DioException in getDriverTrips: $e');
      throw Exception('Lỗi mạng hoặc máy chủ không phản hồi. Vui lòng thử lại.');
    } catch (e) {
      // Bắt các lỗi khác
      print('Error in getDriverTrips: $e');
      throw Exception('Đã xảy ra lỗi không mong muốn.');
    }
  }
}