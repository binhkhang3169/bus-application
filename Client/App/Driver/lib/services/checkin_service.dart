import 'package:dio/dio.dart';
import 'package:taixe/models/checkin.dart';
import 'package:taixe/services/auth_service.dart';

class CheckinService {
  final Dio _dio = Dio();
  final AuthRepository _authRepository = AuthRepository();
  final String _baseUrl = 'http://57.155.76.74/api/v1'; // <<--- THAY URL CỦA BẠN VÀO ĐÂY

  /// Lấy danh sách các vé đã check-in cho một chuyến đi
  Future<List<Checkin>> getCheckedInSeats(int tripId) async {
    final token = await _authRepository.getValidAccessToken();
    if (token == null) throw Exception('Phiên đăng nhập hết hạn.');

    try {
      final response = await _dio.get(
        '$_baseUrl/checkin/trip/$tripId',
        options: Options(headers: {'Authorization': 'Bearer $token'}),
      );
      if (response.statusCode == 200 && response.data['code'] == 200) {
        if (response.data['data'] == null) return [];
        List<dynamic> data = response.data['data'];
        return data.map((json) => Checkin.fromJson(json)).toList();
      } else {
        throw Exception(response.data['message'] ?? 'Lấy dữ liệu check-in thất bại.');
      }
    } on DioException catch (e) {
      throw Exception('Lỗi mạng: ${e.message}');
    }
  }

  /// Gửi yêu cầu check-in một vé
  Future<Checkin> performCheckin({required String qrContent, required int tripId}) async {
    final token = await _authRepository.getValidAccessToken();
    if (token == null) throw Exception('Phiên đăng nhập hết hạn.');
    
    final requestBody = {
      'qr_content': qrContent,
      'trip_id': tripId.toString(),
    };
    
  print("$requestBody");

    try {
      final response = await _dio.post(
        '$_baseUrl/checkin/',
        data: requestBody,
        options: Options(headers: {'Authorization': 'Bearer $token'}),
      );
      if (response.statusCode == 200 && response.data['code'] == 200) {
         return Checkin.fromJson(response.data['data']);
      } else {
        // Trường hợp API trả về 200 nhưng code bên trong là lỗi
        throw Exception(response.data['message'] ?? 'Check-in thất bại.');
      }
    } on DioException catch (e) {
        // Xử lý các mã lỗi HTTP cụ thể
        if (e.response != null) {
             final errorMessage = e.response!.data['message'] ?? 'Lỗi không xác định từ máy chủ.';
             throw Exception(errorMessage);
        }
        throw Exception('Lỗi mạng khi thực hiện check-in.');
    }
  }
}