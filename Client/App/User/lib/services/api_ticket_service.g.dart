// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_ticket_service.dart';

// **************************************************************************
// RetrofitGenerator (Manual)
// **************************************************************************

class _ApiTicketService implements ApiTicketService {
  _ApiTicketService(this._dio, {this.baseUrl}) {
    baseUrl ??= 'http://57.155.76.74/api/v1';
  }

  final Dio _dio;
  String? baseUrl;

  @override
  Future<ApiResponse1<List<TripInfo>>> getTicketsAvailable(int id) async {
    final response = await _dio.get<Map<String, dynamic>>(
      '/tickets-available/$id',
      options: Options(method: 'GET', responseType: ResponseType.json),
    );

    return ApiResponse1<List<TripInfo>>.fromJson(
      response.data!,
      (json) =>
          (json as List)
              .map((e) => TripInfo.fromJson(e as Map<String, dynamic>))
              .toList(),
    );
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
