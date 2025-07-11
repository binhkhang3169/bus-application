// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_info.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserInfo _$UserInfoFromJson(Map<String, dynamic> json) => UserInfo(
  id: (json['id'] as num).toInt(),
  username: json['username'] as String,
  phoneNumber: json['phoneNumber'] as String,
  fullName: json['fullName'] as String,
  address: json['address'] as String,
  gender: json['gender'] as String,
  active: (json['active'] as num).toInt(),
  image: json['image'] as String?,
);

Map<String, dynamic> _$UserInfoToJson(UserInfo instance) => <String, dynamic>{
  'id': instance.id,
  'username': instance.username,
  'phoneNumber': instance.phoneNumber,
  'fullName': instance.fullName,
  'address': instance.address,
  'gender': instance.gender,
  'active': instance.active,
  'image': instance.image,
};
