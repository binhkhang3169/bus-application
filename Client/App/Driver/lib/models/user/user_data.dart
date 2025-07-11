import 'package:json_annotation/json_annotation.dart';

part 'user_data.g.dart';

@JsonSerializable()
class UserData {
  @JsonKey(name: 'id') // Assuming id can be String or int, adjust if necessary
  final dynamic id;
  final String username;
  final String role;

  UserData({required this.id, required this.username, required this.role});

  factory UserData.fromJson(Map<String, dynamic> json) =>
      _$UserDataFromJson(json);

  Map<String, dynamic> toJson() => _$UserDataToJson(this);
}