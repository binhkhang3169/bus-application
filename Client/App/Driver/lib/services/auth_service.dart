import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:taixe/models/api_response.dart';
import 'package:taixe/models/user/login_request.dart';
import 'package:taixe/models/user/login_response.dart';
import 'package:taixe/models/user/user_info.dart';
import 'api_user_service.dart';

// Định nghĩa các hằng số cho key của SharedPreferences để tránh lỗi chính tả
class _PrefsKeys {
  static const String accessToken = 'accessToken';
  static const String accessTokenExpiry = 'accessTokenExpiry';
  static const String refreshToken = 'refreshToken';
  static const String role = 'role';
  static const String userId = 'userId';
  static const String username = 'username';
  static const String rememberMe = 'rememberMe';
  static const String phoneNumber = 'phoneNumber';
  static const String fullName = 'fullName';
  static const String address = 'address';
  static const String gender = 'gender';
  static const String userActiveStatus = 'userActiveStatus';
  static const String image = 'image';
}

class AuthRepository {
  final Dio _dio = Dio();
  late final ApiUserService _apiUserService;
  SharedPreferences? _prefsInstance;

  // Sử dụng một instance duy nhất của SharedPreferences
  Future<SharedPreferences> get _prefs async {
    _prefsInstance ??= await SharedPreferences.getInstance();
    return _prefsInstance!;
  }

  AuthRepository() {
    _apiUserService = ApiUserService(_dio);
  }

  /// Thực hiện đăng nhập và trả về phản hồi.
  Future<LoginResponse> login(LoginRequest request) async {
    return await _apiUserService.login(request);
  }

  /// Lưu trữ token và dữ liệu người dùng từ phản hồi đăng nhập.
  Future<void> saveLoginData({
    required String accessToken,
    required String refreshToken,
    required String role,
    required String userId,
    required String username,
    required bool rememberMe,
    int accessTokenLifetimeMinutes = 60, // Đặt thời gian sống của token (ví dụ: 60 phút)
  }) async {
    final prefs = await _prefs;
    await prefs.setString(_PrefsKeys.accessToken, accessToken);
    await prefs.setString(_PrefsKeys.refreshToken, refreshToken);
    await prefs.setString(_PrefsKeys.role, role);
    await prefs.setString(_PrefsKeys.userId, userId);
    await prefs.setString(_PrefsKeys.username, username);
    await prefs.setBool(_PrefsKeys.rememberMe, rememberMe);

    // Lưu thời gian hết hạn của access token
    final expiryTime = DateTime.now().add(Duration(minutes: accessTokenLifetimeMinutes));
    await prefs.setInt(_PrefsKeys.accessTokenExpiry, expiryTime.millisecondsSinceEpoch);

    print("Đã lưu dữ liệu đăng nhập. Token hết hạn lúc: $expiryTime, Ghi nhớ đăng nhập: $rememberMe");
  }

  /// Phương thức chính để lấy một access token hợp lệ.
  /// Tự động làm mới nếu token đã hết hạn.
  Future<String?> getValidAccessToken() async {
    final prefs = await _prefs;
    final expiryMillis = prefs.getInt(_PrefsKeys.accessTokenExpiry);
    final accessToken = prefs.getString(_PrefsKeys.accessToken);

    if (accessToken != null && expiryMillis != null) {
      final expiryTime = DateTime.fromMillisecondsSinceEpoch(expiryMillis);
      // Kiểm tra xem token hiện tại còn hợp lệ không (trừ đi 1 phút để an toàn)
      if (DateTime.now().isBefore(expiryTime.subtract(const Duration(minutes: 1)))) {
        print("Access token còn hợp lệ.");
        return accessToken;
      }
    }

    print("Access token đã hết hạn hoặc không tồn tại. Đang thử làm mới...");
    return await _refreshToken();
  }

  /// Cố gắng làm mới access token bằng refresh token đã lưu.
  Future<String?> _refreshToken() async {
    final prefs = await _prefs;
    final storedRefreshToken = prefs.getString(_PrefsKeys.refreshToken);

    if (storedRefreshToken == null) {
      print("Không có refresh token. Không thể làm mới.");
      await clearAuthData();
      return null;
    }

    try {
      final response = await _apiUserService.refreshToken(storedRefreshToken);

      if (response.code == 200 && response.data != null) {
        final newAccessToken = response.data!['accessToken'];

        if (newAccessToken != null) {
          const int lifetimeMinutes = 60; // Thời gian sống của token mới
          final newExpiryTime = DateTime.now().add(const Duration(minutes: lifetimeMinutes));

          await prefs.setString(_PrefsKeys.accessToken, newAccessToken);
          await prefs.setInt(_PrefsKeys.accessTokenExpiry, newExpiryTime.millisecondsSinceEpoch);
          
          print("Làm mới token thành công. Token mới hết hạn lúc $newExpiryTime.");
          return newAccessToken;
        }
      }

      print("Làm mới token thất bại (Mã lỗi API: ${response.code}): ${response.message}.");
      await clearAuthData();
      return null;
    } catch (e) {
      print("Lỗi mạng trong quá trình làm mới token: $e. Xóa dữ liệu xác thực.");
      await clearAuthData();
      return null;
    }
  }

  /// Kiểm tra trạng thái xác thực khi khởi động ứng dụng.
  Future<String?> attemptAutoLogin() async {
    final prefs = await _prefs;
    final wasRememberMe = prefs.getBool(_PrefsKeys.rememberMe) ?? false;
    final refreshTokenExists = prefs.getString(_PrefsKeys.refreshToken) != null;

    if (wasRememberMe && refreshTokenExists) {
      print("Đang thử tự động đăng nhập...");
      // Lấy token hợp lệ (có thể sẽ làm mới nếu cần)
      final token = await getValidAccessToken();
      if (token != null) {
        // Nếu lấy được token, cập nhật lại thông tin chi tiết người dùng
        await fetchAndSaveDetailedUserInfo(token);
        return prefs.getString(_PrefsKeys.role);
      }
    } else if (!wasRememberMe && refreshTokenExists) {
      print("Phiên đăng nhập không được thiết lập để ghi nhớ. Xóa dữ liệu.");
      await clearAuthData();
    }
    
    return null; // Không thể tự động đăng nhập
  }

  /// Thông báo cho máy chủ về việc đăng xuất và xóa tất cả dữ liệu cục bộ.
  Future<void> logout() async {
    final prefs = await _prefs;
    final refreshToken = prefs.getString(_PrefsKeys.refreshToken);
    if (refreshToken != null) {
      try {
        print("Gọi API đăng xuất của máy chủ với refresh token.");
        await _apiUserService.logout(refreshToken);
      } catch (e) {
        print("Lỗi khi đăng xuất khỏi máy chủ (vẫn tiếp tục dọn dẹp cục bộ): $e");
      }
    }
    await clearAuthData();
    print("Người dùng đã đăng xuất, dữ liệu cục bộ đã được xóa.");
  }

  /// Xóa tất cả dữ liệu xác thực và hồ sơ người dùng đã lưu.
  Future<void> clearAuthData() async {
    final prefs = await _prefs;
    await prefs.clear(); // Xóa tất cả dữ liệu trong SharedPreferences
    print("Tất cả dữ liệu xác thực và hồ sơ người dùng cục bộ đã được xóa.");
  }

  /// Lấy và lưu thông tin chi tiết của người dùng.
  Future<bool> fetchAndSaveDetailedUserInfo(String accessToken) async {
    try {
      ApiResponse<UserInfo> response = await _apiUserService.getUserInfo("Bearer $accessToken");

      if (response.code != 200 || response.data == null) {
        print("Lỗi API hoặc không có dữ liệu người dùng: ${response.message}");
        return false;
      }
      
      UserInfo userInfo = response.data!;
      final prefs = await _prefs;

      // Lưu thông tin chi tiết vào SharedPreferences
      await prefs.setString(_PrefsKeys.userId, userInfo.id.toString());
      await prefs.setString(_PrefsKeys.username, userInfo.username);
      await prefs.setString(_PrefsKeys.phoneNumber, userInfo.phoneNumber ?? '');
      await prefs.setString(_PrefsKeys.fullName, userInfo.fullName ?? '');
      await prefs.setString(_PrefsKeys.address, userInfo.address ?? '');
      await prefs.setString(_PrefsKeys.gender, userInfo.gender ?? '');
      await prefs.setString(_PrefsKeys.image, userInfo.image ?? "");
      await prefs.setBool(_PrefsKeys.userActiveStatus, userInfo.active == 1);

      print("Đã lấy và lưu thông tin chi tiết cho người dùng ${userInfo.username}.");
      return true;
    } catch (e) {
      print("Thất bại khi lấy thông tin chi tiết người dùng: $e");
      if (e is DioException && (e.response?.statusCode == 401 || e.response?.statusCode == 403)) {
        print("getUserInfo thất bại do token không hợp lệ hoặc bị cấm.");
      }
      return false;
    }
  }

  // Getters để các phần khác của ứng dụng có thể sử dụng nếu cần
  Future<String?> getRole() async => (await _prefs).getString(_PrefsKeys.role);
  Future<String?> getAccessToken() async => await getValidAccessToken();
  Future<bool> isLoggedIn() async => await getAccessToken() != null;
}