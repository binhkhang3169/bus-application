// lib/models/transaction_history_model.dart

import 'package:flutter/foundation.dart';

class TransactionHistory {
  final int id;
  final int accountId;
  final String transactionType;
  final int? amount; // Có thể là null, ví dụ khi tạo tài khoản
  final String? currency; // Có thể là null
  final DateTime transactionTimestamp;
  final String description;

  TransactionHistory({
    required this.id,
    required this.accountId,
    required this.transactionType,
    this.amount,
    this.currency,
    required this.transactionTimestamp,
    required this.description,
  });

  factory TransactionHistory.fromJson(Map<String, dynamic> json) {
    try {
      return TransactionHistory(
        id: json['id'],
        accountId: json['account_id'],
        transactionType: json['transaction_type'] as String,
        amount: json['amount'] as int?, // Trình phân tích JSON của Dart sẽ tự xử lý null
        currency: json['currency'] as String?,
        // API trả về string, cần parse thành DateTime
        transactionTimestamp: DateTime.parse(json['transaction_timestamp'] as String),
        description: json['description'] ?? '', // Xử lý trường hợp description là null
      );
    } catch (e) {
      debugPrint('Error parsing TransactionHistory: $e');
      debugPrint('Problematic JSON: $json');
      // Ném lỗi để báo hiệu dữ liệu không hợp lệ
      throw FormatException('Failed to parse transaction history item.', e);
    }
  }
}