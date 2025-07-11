// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_trip_service.dart';

// **************************************************************************
// RetrofitGenerator
// **************************************************************************

// ignore_for_file: unnecessary_brace_in_string_interps,no_leading_underscores_for_local_identifiers

class _ApiTripService implements ApiTripService {
  _ApiTripService(this._dio, {this.baseUrl}) {
    baseUrl ??= 'http://57.155.76.74/api/v1';
  }

  final Dio _dio;

  String? baseUrl;

  @override
  Future<ApiResponse1<List<TripInfo>>> getTrips(
    fromCityName,
    fromCityId,
    toCityName,
    toCityId,
    departureDate,
    quantity,
  ) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{
      r'from': fromCityName,
      r'fromId': fromCityId,
      r'to': toCityName,
      r'toId': toCityId,
      r'fromTime': departureDate,
      r'quantity': quantity,
    };
    final _headers = <String, dynamic>{};
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse1<List<TripData>>>(
        Options(method: 'GET', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/trips/search',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse1<List<TripInfo>>.fromJson(
      _result
          .data!, // This is the entire response map: {"code": ..., "message": ..., "data": ...}
      (jsonData) {
        // jsonData here is the value of the 'data' key from the API response
        if (jsonData is List) {
          return jsonData
              .map((item) => TripInfo.fromJson(item as Map<String, dynamic>))
              .toList();
        }
        // Optional: Handle cases where jsonData is not a list,
        // though for this specific endpoint, it's expected to be a list.
        // You could return an empty list or throw an error.
        return <
          TripInfo
        >[]; // Default to an empty list if data is not in the expected List format
      },
    );
    return value;
  }

  @override
  Future<ApiResponse<TripInfo>> getTripDetails(int tripId) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<TripInfo>>(
        Options(method: 'GET', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/trips/${tripId}/seats',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<TripInfo>.fromJson(
      _result.data!,
      (json) => TripInfo.fromJson(json as Map<String, dynamic>),
    );
    return value;
  }

  @override
  Future<ApiResponse<SeatsData>> getAvailableSeats(int tripId) async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<SeatsData>>(
        Options(method: 'GET', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/tickets-available/${tripId}',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<SeatsData>.fromJson(
      _result.data!,
      (json) => SeatsData.fromJson(json as Map<String, dynamic>),
    );
    return value;
  }

  @override
  Future<ApiResponse<List<AddressInfo>>> getListAddress() async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<List<AddressInfo>>>(
        Options(method: 'GET', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/provinces',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<List<AddressInfo>>.fromJson(
      _result
          .data!, // This is the entire response map: {"code": ..., "message": ..., "data": ...}
      (jsonData) {
        // jsonData here is the value of the 'data' key from the API response
        if (jsonData is List) {
          return jsonData
              .map((item) => AddressInfo.fromJson(item as Map<String, dynamic>))
              .toList();
        }
        // Optional: Handle cases where jsonData is not a list,
        // though for this specific endpoint, it's expected to be a list.
        // You could return an empty list or throw an error.
        return <
          AddressInfo
        >[]; // Default to an empty list if data is not in the expected List format
      },
    );
    return value;
  }

  @override
  Future<ApiResponse<List<Station>>> getAllStations() async {
    const _extra = <String, dynamic>{};
    final queryParameters = <String, dynamic>{};
    final _headers = <String, dynamic>{};
    final Map<String, dynamic>? _data = null;
    final _result = await _dio.fetch<Map<String, dynamic>>(
      _setStreamType<ApiResponse<List<Station>>>(
        Options(method: 'GET', headers: _headers, extra: _extra)
            .compose(
              _dio.options,
              '/stations',
              queryParameters: queryParameters,
              data: _data,
            )
            .copyWith(baseUrl: baseUrl ?? _dio.options.baseUrl),
      ),
    );
    final value = ApiResponse<List<Station>>.fromJson(
      _result
          .data!, // This is the entire response map: {"code": ..., "message": ..., "data": ...}
      (jsonData) {
        // jsonData here is the value of the 'data' key from the API response
        if (jsonData is List) {
          return jsonData
              .map((item) => Station.fromJson(item as Map<String, dynamic>))
              .toList();
        }
        // Optional: Handle cases where jsonData is not a list,
        // though for this specific endpoint, it's expected to be a list.
        // You could return an empty list or throw an error.
        return <
          Station
        >[]; // Default to an empty list if data is not in the expected List format
      },
    );
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
