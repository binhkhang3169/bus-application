import 'package:flutter/material.dart';
import 'package:lottie/lottie.dart';
import 'package:intl/intl.dart';

void main() {
  runApp(
    MaterialApp(
      debugShowCheckedModeBanner: false,
      home: PaymentSuccessPage(
        customerName: "Nguyễn Cao Kỳ",
        amount: 1700000,
        transactionId: "GD202405310001",
      ),
    ),
  );
}

class PaymentSuccessPage extends StatelessWidget {
  final String customerName;
  final double amount;
  final String transactionId;

  PaymentSuccessPage({
    required this.customerName,
    required this.amount,
    required this.transactionId,
    Key? key,
  }) : super(key: key);

  final String paymentTime = DateFormat(
    'dd/MM/yyyy HH:mm',
  ).format(DateTime.now());

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.green[50],
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 60),
          child: Container(
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(20),
              boxShadow: const [
                BoxShadow(
                  color: Colors.black12,
                  blurRadius: 12,
                  offset: Offset(0, 6),
                ),
              ],
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Lottie.asset(
                  'assets/images/success.json',
                  width: 150,
                  repeat: false,
                ),
                const SizedBox(height: 16),
                const Text(
                  "Thanh toán thành công!",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: Colors.green,
                  ),
                ),
                const SizedBox(height: 12),
                Text(
                  "Cảm ơn bạn, $customerName!",
                  style: const TextStyle(fontSize: 18),
                ),
                const SizedBox(height: 20),
                Divider(thickness: 1, color: Colors.grey[300]),
                const SizedBox(height: 12),
                _buildInfoRow("Mã giao dịch:", transactionId),
                const SizedBox(height: 8),
                _buildInfoRow("Thời gian:", paymentTime),
                const SizedBox(height: 8),
                _buildInfoRow(
                  "Giá trị:",
                  '${NumberFormat.currency(locale: 'vi_VN', symbol: '₫').format(amount)}',
                ),
                const SizedBox(height: 24),
                ElevatedButton.icon(
                  onPressed: () => Navigator.pop(context),
                  icon: const Icon(Icons.home),
                  label: const Text("Về trang chính"),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.green,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(
                      horizontal: 32,
                      vertical: 14,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(30),
                    ),
                    elevation: 4,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(fontWeight: FontWeight.w500, fontSize: 16),
        ),
        Flexible(
          child: Text(
            value,
            textAlign: TextAlign.right,
            style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
          ),
        ),
      ],
    );
  }
}
