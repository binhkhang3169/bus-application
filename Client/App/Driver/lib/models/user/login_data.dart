import 'package:json_annotation/json_annotation.dart';
import 'user_data.dart'; // Import the UserData class

part 'login_data.g.dart';

@JsonSerializable()
class LoginData {
  final String accessToken;
  final String refreshToken;
  final UserData user;

  LoginData({
    required this.accessToken,
    required this.refreshToken,
    required this.user,
  });

  factory LoginData.fromJson(Map<String, dynamic> json) =>
      _$LoginDataFromJson(json);

  Map<String, dynamic> toJson() => _$LoginDataToJson(this);
}