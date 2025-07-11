// lib/services/dio_client.dart
import 'package:dio/dio.dart';

class DioClient {
  static Dio createDio() {
    final dio = Dio(
      BaseOptions(
        connectTimeout: Duration(seconds: 20), // Increased timeout
        receiveTimeout: Duration(seconds: 20), // Increased timeout
        headers: {'Content-Type': 'application/json'},
      ),
    );

    dio.interceptors.add(
      LogInterceptor(
        requestBody: true,
        responseBody: true,
        requestHeader: true,
        responseHeader: false,
        error: true,
        logPrint:
            (object) => print(object.toString()), // Ensure logs are visible
      ),
    );
    return dio;
  }
}
