import 'package:flutter/foundation.dart';

class Account {
  final int id;
  final String ownerName;
  final int balance;
  final String currency;
  final String status;
  final DateTime createdAt;
  final DateTime updatedAt;

  Account({
    required this.id,
    required this.ownerName,
    required this.balance,
    required this.currency,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      id: json['id'],
      ownerName: json['owner_name'],
      balance: json['balance'],
      currency: json['currency'],
      status: json['status'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}