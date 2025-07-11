import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:retrofit/error_logger.dart';
import 'package:retrofit/http.dart';
import 'package:taixe/models/api_response.dart';
import 'package:taixe/models/user/login_request.dart';
import 'package:taixe/models/user/login_response.dart';
import 'package:taixe/models/user/resend_otp.dart';
import 'package:taixe/models/user/signup_request.dart';
import 'package:taixe/models/user/user_info.dart';

part 'api_user_service.g.dart'; // Make sure this is correct

@RestApi(baseUrl: "http://57.155.76.74/api/v1")
// @RestApi(baseUrl: "http://192.168.70.182:8081/api/v1")
abstract class ApiUserService {
  factory ApiUserService(Dio dio, {String baseUrl}) = _ApiUserService;

  @POST("/auth/login")
  Future<LoginResponse> login(@Body() LoginRequest request);

  @POST("/change-password")
  Future<ApiResponse> changePassword(
    // Consider specifying <dynamic> or a specific type if ApiResponse is generic
    @Header("Authorization") String token,
    @Query("oldPassword") String oldPassword,
    @Query("newPassword") String newPassword,
  );

  @POST("/forgot-password")
  Future<void> sendResetPassword(@Query("email") String email);

  @POST("/signup")
  Future<ApiResponse> register(@Body() SignupRequest request); // Consider <dynamic>

  @POST("/verify-otp")
  Future<ApiResponse> verifyOtp(@Body() SignupRequest request); // Consider <dynamic> (assuming request here might be OtpVerificationRequest)

  @POST("/resend-otp")
  Future<ApiResponse> resendOtp(@Body() ResendOtpRequest request); // Consider <dynamic>

  @GET("/employee/info")
  Future<ApiResponse<UserInfo>> getUserInfo(
    @Header("Authorization") String token,
  );

  @POST("/employee/change-info")
  Future<ApiResponse> changeInfo(
    // Consider <dynamic>
    @Header("Authorization") String token,
    @Body() Map<String, dynamic> data,
  );

  @POST("/employee/change-image")
  Future<void> changeImage(
    @Header("Authorization") String token,
    @Query("image") String image,
  );

  // New methods for refresh token and logout:

  /// Refreshes the access token using a refresh token.
  /// The refresh token is sent in the "X-Refresh-Token" header.
  ///
  /// Returns an ApiResponse containing a map with the new "accessToken"
  /// and the existing "refreshToken".
  @POST("/auth/refresh-token")
  Future<ApiResponse<Map<String, String>>> refreshToken(
    @Header("X-Refresh-Token") String refreshToken,
  );

  /// Logs out the user by invalidating the refresh token.
  /// The refresh token is sent in the "X-Refresh-Token" header.
  ///
  /// Returns an ApiResponse, typically with an empty string or null for data.
  @POST("/auth/logout")
  Future<ApiResponse<dynamic>> logout(
    // Using dynamic as data is often an empty string or not strictly typed
    @Header("X-Refresh-Token") String refreshToken,
  );

  // ... (other commented-out methods remain the same)
}
