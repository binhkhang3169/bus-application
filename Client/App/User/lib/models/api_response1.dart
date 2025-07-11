class ApiResponse1<T> {
  final String message;
  final T data;
  final int code;

  ApiResponse1({required this.message, required this.data, required this.code});

  factory ApiResponse1.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic) fromJsonT,
  ) {
    final rawData = json['data'];
    return ApiResponse1(
      message: json['message'] as String,
      data: fromJsonT(rawData),
      code: json['code'] as int,
    );
  }
}
