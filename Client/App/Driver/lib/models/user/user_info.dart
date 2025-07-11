import 'package:json_annotation/json_annotation.dart';

part 'user_info.g.dart';

@JsonSerializable()
class UserInfo {
  final int id;
  final String username;
  final String phoneNumber;
  final String fullName;
  final String address;
  final String gender;
  final int active;
  final String? image;

  UserInfo({
    required this.id,
    required this.username,
    required this.phoneNumber,
    required this.fullName,
    required this.address,
    required this.gender,
    required this.active,
    required this.image,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) =>
      _$UserInfoFromJson(json);
  Map<String, dynamic> toJson() => _$UserInfoToJson(this);
}
