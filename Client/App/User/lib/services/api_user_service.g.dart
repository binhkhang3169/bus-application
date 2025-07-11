// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_user_service.dart';

// **************************************************************************
// RetrofitGenerator
// **************************************************************************

// ignore_for_file: unnecessary_brace_in_string_interps,no_leading_underscores_for_local_identifiers

class _ApiUserService implements ApiUserService {
  _ApiUserService(this._dio, {this.baseUrl}) {
    baseUrl ??= 'http://57.155.76.74/api/v1';
  }

  final Dio _dio;

  String? baseUrl;

  @override
  Future<LoginResponse> login(request) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final _data = <String, dynamic>{};
    _data.addAll(request.toJson());
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<LoginResponse>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/auth/login',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = LoginResponse.fromJson(_result.data!);
    return value;
  }

  @override
  Future<ApiResponse<dynamic>> changePassword(
    token,
    oldPassword,
    newPassword,
  ) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{
      r'oldPassword': oldPassword,
      r'newPassword': newPassword,
    };
    final _headers = <String, dynamic>{r'Authorization': token};
    _headers.removeWhere((k, v) => v == null);
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<dynamic>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/change-password',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<dynamic>.fromJson(_result.data!, (json) => json);
    return value;
  }

  @override
  Future<void> sendResetPassword(email) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{r'email': email};
    final _headers = <String, dynamic>{};
    final Map<String, dynamic>? _data = null;
    await _dio.fetch<void>(
      _setStreamType<void>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/forgot-password',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
  }

  @override
  Future<ApiResponse<dynamic>> register(request) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final _data = <String, dynamic>{};
    _data.addAll(request.toJson());
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<dynamic>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/signup',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<dynamic>.fromJson(_result.data!, (json) => json);
    return value;
  }

  @override
  Future<ApiResponse<dynamic>> verifyOtp(request) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final _data = <String, dynamic>{};
    _data.addAll(request.toJson());
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<dynamic>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/verify-otp',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<dynamic>.fromJson(_result.data!, (json) => json);
    return value;
  }

  @override
  Future<ApiResponse<dynamic>> resendOtp(request) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final _data = <String, dynamic>{};
    _data.addAll(request.toJson());
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<dynamic>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/resend-otp',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<dynamic>.fromJson(_result.data!, (json) => json);
    return value;
  }

  @override
  Future<UserInfo> getUserInfo(token) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{r'Authorization': token};
    _headers.removeWhere((k, v) => v == null);
    final Map<String, dynamic>? _data = null;

    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<UserInfo>(
        Options(method: 'GET', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/customer/info',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );

    // Lấy phần `data` trong response
    final value = UserInfo.fromJson(_result.data!['data']);
    return value;
  }

  @override
  Future<ApiResponse<dynamic>> changeInfo(token, data) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{r'Authorization': token};
    _headers.removeWhere((k, v) => v == null);
    final _data = <String, dynamic>{};
    _data.addAll(data);
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<dynamic>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/customer/change-info',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<dynamic>.fromJson(_result.data!, (json) => json);
    return value;
  }

  @override
  Future<ApiResponse<Map<String, String>>> refreshToken(refreshToken) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{r'X-Refresh-Token': refreshToken};
    _headers.removeWhere((k, v) => v == null);
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<Map<String, String>>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/auth/refresh-token',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<Map<String, String>>.fromJson(
      _result.data!,
      (json) => json,
    );
    return value;
  }

  @override
  Future<ApiResponse<dynamic>> logout(refreshToken) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{r'X-Refresh-Token': refreshToken};
    _headers.removeWhere((k, v) => v == null);
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<dynamic>>(
        Options(method: 'POST', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/auth/logout',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<dynamic>.fromJson(_result.data!, (json) => json);
    return value;
  }

  RequestOptions _setStreamType<T>(RequestOptions requestOptions) {
    if (T != dynamic &&
        !(requestOptions.responseType == ResponseType.bytes ||
            requestOptions.responseType == ResponseType.stream)) {
      if (T == String) {
        requestOptions.responseType = ResponseType.plain;
      } else {
        requestOptions.responseType = ResponseType.json;
      }
    }
    return requestOptions;
  }
}
