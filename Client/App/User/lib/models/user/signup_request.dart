import 'package:json_annotation/json_annotation.dart';

part 'signup_request.g.dart';

@JsonSerializable()
class SignupRequest {
  final String username;
  final String password;
  final String phoneNumber;
  final String fullName;
  final String address;
  final String gender;
  final String otp;

  SignupRequest({
    required this.username,
    required this.password,
    required this.phoneNumber,
    required this.fullName,
    required this.address,
    required this.gender,
    required this.otp,
  });

  factory SignupRequest.fromJson(Map<String, dynamic> json) =>
      _$SignupRequestFromJson(json);

  Map<String, dynamic> toJson() => _$SignupRequestToJson(this);
}
