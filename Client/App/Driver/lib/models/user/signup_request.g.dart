// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'signup_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SignupRequest _$SignupRequestFromJson(Map<String, dynamic> json) =>
    SignupRequest(
      username: json['username'] as String,
      password: json['password'] as String,
      phoneNumber: json['phoneNumber'] as String,
      fullName: json['fullName'] as String,
      address: json['address'] as String,
      gender: json['gender'] as String,
      otp: json['otp'] as String,
    );

Map<String, dynamic> _$SignupRequestToJson(SignupRequest instance) =>
    <String, dynamic>{
      'username': instance.username,
      'password': instance.password,
      'phoneNumber': instance.phoneNumber,
      'fullName': instance.fullName,
      'address': instance.address,
      'gender': instance.gender,
      'otp': instance.otp,
    };
