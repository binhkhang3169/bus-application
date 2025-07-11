import 'package:caoky/models/api_response.dart';
import 'package:caoky/models/user/login_request.dart';
import 'package:caoky/models/user/login_response.dart';
import 'package:caoky/services/notification_service.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'api_user_service.dart';
import '../models/user/user_info.dart';

// Define constants for SharedPreferences keys to avoid typos
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
}


class AuthRepository {
  final Dio _dio = Dio();
  late final ApiUserService _apiUserService;

  SharedPreferences? _prefsInstance;

  Future<SharedPreferences> get _prefs async {
    _prefsInstance ??= await SharedPreferences.getInstance();
    return _prefsInstance!;
  }

  AuthRepository() {
    _apiUserService = ApiUserService(_dio);
  }

  /// Performs login and returns the response.
  Future<LoginResponse> login(LoginRequest request) async {
    return await _apiUserService.login(request);
  }

  /// Saves tokens and user data from the login response.
  Future<void> saveLoginData({
    required String accessToken,
    required String refreshToken,
    required String role,
    required String userId,
    required String username,
    required bool rememberMe,
    int accessTokenLifetimeMinutes = 10, // Default lifetime
  }) async {
    try {
      await NotificationService.initializeNotification(userId: userId);
      print("Notification service initialized for user: $userId");
    } catch (e) {
      print("Failed to initialize notification service: $e");
    }

    final prefs = await _prefs;
    await prefs.setString(_PrefsKeys.accessToken, accessToken);
    await prefs.setString(_PrefsKeys.refreshToken, refreshToken);
    await prefs.setString(_PrefsKeys.role, role);
    await prefs.setString(_PrefsKeys.userId, userId);
    await prefs.setString(_PrefsKeys.username, username);
    await prefs.setBool(_PrefsKeys.rememberMe, rememberMe);

    final expiryTime = DateTime.now().add(Duration(minutes: accessTokenLifetimeMinutes));
    await prefs.setInt(_PrefsKeys.accessTokenExpiry, expiryTime.millisecondsSinceEpoch);

    print("Saved login data. Token expires at: $expiryTime, RememberMe: $rememberMe");
  }

  /// The primary method for getting a valid access token.
  Future<String?> getValidAccessToken() async {
    final prefs = await _prefs;
    final expiryMillis = prefs.getInt(_PrefsKeys.accessTokenExpiry);
    final accessToken = prefs.getString(_PrefsKeys.accessToken);

    if (accessToken != null && expiryMillis != null) {
      final expiryTime = DateTime.fromMillisecondsSinceEpoch(expiryMillis);
      if (DateTime.now().isBefore(expiryTime)) {
        print("Access token is valid.");
        return accessToken;
      }
    }

    print("Access token expired or missing. Attempting to refresh...");
    return await _refreshToken();
  }

  /// Attempts to refresh the access token using the stored refresh token.
  Future<String?> _refreshToken() async {
    final prefs = await _prefs;
    final storedRefreshToken = prefs.getString(_PrefsKeys.refreshToken);

    if (storedRefreshToken == null) {
      print("No refresh token available. Cannot refresh.");
      await clearAuthData();
      return null;
    }

    try {
      final response = await _apiUserService.refreshToken(storedRefreshToken);

      if (response.code == 200 && response.data != null) {
        final newAccessToken = response.data!['accessToken'];

        if (newAccessToken != null) {
          const int lifetimeMinutes = 10;
          final newExpiryTime = DateTime.now().add(Duration(minutes: lifetimeMinutes));

          await prefs.setString(_PrefsKeys.accessToken, newAccessToken);
          await prefs.setInt(_PrefsKeys.accessTokenExpiry, newExpiryTime.millisecondsSinceEpoch);
          
          print("Token refreshed successfully. New token expires at $newExpiryTime.");
          return newAccessToken;
        }
      }

      print("Failed to refresh token (API Error Code: ${response.code}): ${response.message}.");
      await clearAuthData();
      return null;
    } catch (e) {
      print("Error during token refresh network call: $e. Clearing auth data.");
      await clearAuthData();
      return null;
    }
  }

  /// Checks authentication status on app start.
  Future<String?> attemptAutoLogin() async {
    final prefs = await _prefs;
    final wasRememberMe = prefs.getBool(_PrefsKeys.rememberMe) ?? false;
    final refreshToken = prefs.getString(_PrefsKeys.refreshToken);

    if (wasRememberMe && refreshToken != null) {
      print("Attempting auto-login...");
      final token = await getValidAccessToken();
      if (token != null) {
        await fetchAndSaveDetailedUserInfo(token);
        return prefs.getString(_PrefsKeys.role);
      }
    } else if (!wasRememberMe && refreshToken != null) {
       print("Session was not set to be remembered. Clearing data.");
       await clearAuthData();
    }
    
    return null;
  }

  /// Informs the server of logout and clears all local data.
  Future<void> logout() async {
    final prefs = await _prefs;
    final refreshToken = prefs.getString(_PrefsKeys.refreshToken);
    if (refreshToken != null) {
      try {
        print("Calling server logout with refresh token.");
        await _apiUserService.logout(refreshToken);
      } catch (e) {
        print("Error logging out from server (proceeding with local cleanup): $e");
      }
    }
    await clearAuthData();
    print("User logged out, local data cleared.");
  }

  /// Clears all stored authentication and user profile data.
  Future<void> clearAuthData() async {
    final prefs = await _prefs;
    await prefs.clear(); // Clears all data in SharedPreferences
    print("All local authentication and user profile data have been cleared.");
  }

  /// Fetches and saves detailed user information.
  Future<bool> fetchAndSaveDetailedUserInfo(String accessToken) async {
    try {
      UserInfo userInfo = await _apiUserService.getUserInfo("Bearer $accessToken");
      final prefs = await _prefs;

      await prefs.setString(_PrefsKeys.userId, userInfo.id.toString());
      await prefs.setString(_PrefsKeys.username, userInfo.username);
      await prefs.setString(_PrefsKeys.phoneNumber, userInfo.phoneNumber ?? '');
      await prefs.setString(_PrefsKeys.fullName, userInfo.fullName ?? '');
      await prefs.setString(_PrefsKeys.address, userInfo.address ?? '');
      await prefs.setString(_PrefsKeys.gender, userInfo.gender ?? '');
      await prefs.setBool(_PrefsKeys.userActiveStatus, userInfo.active == 1);

      print("Fetched and saved detailed user info for ${userInfo.username}.");
      return true;
    } catch (e) {
      print("Failed to fetch detailed user info: $e");
      if (e is DioException && (e.response?.statusCode == 401 || e.response?.statusCode == 403)) {
        print("getUserInfo failed due to unauthorized/forbidden access.");
      }
      return false;
    }
  }
}