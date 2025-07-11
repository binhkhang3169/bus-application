import 'package:json_annotation/json_annotation.dart';
import 'login_data.dart'; // Import the LoginData class

part 'login_response.g.dart';

@JsonSerializable()
class LoginResponse {
  final int code;
  final String message;
  final LoginData data; // Changed from String? token to LoginData

  LoginResponse({required this.code, required this.message, required this.data});

  factory LoginResponse.fromJson(Map<String, dynamic> json) =>
      _$LoginResponseFromJson(json);

  Map<String, dynamic> toJson() => _$LoginResponseToJson(this);
}