// lib/services/dio_client.dart
import 'package:dio/dio.dart';

class DioClient {
  static Dio createDio() {
    final dio = Dio(
      BaseOptions(
        // The baseUrl is specified in the @RestApi annotation in ApiTripService
        // So, it's not strictly needed here unless you want a default/fallback
        // or for other services that don't use Retrofit's baseUrl.
        // baseUrl: "YOUR_FALLBACK_BASE_URL_IF_NEEDED",
        connectTimeout: Duration(seconds: 20), // Increased timeout
        receiveTimeout: Duration(seconds: 20), // Increased timeout
        headers: {
          'Content-Type': 'application/json',
          // Add other common headers here if needed
          // e.g., 'Accept': 'application/json',
        },
      ),
    );

    // Add interceptors for logging, error handling, or authentication
    dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
      requestHeader: true,
      responseHeader: false,
      error: true,
      logPrint: (object) => print(object.toString()), // Ensure logs are visible
    ));
    return dio;
  }
}