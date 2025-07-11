// lib/ui/account/deposit_page.dart
import 'dart:convert';
import 'dart:developer';
import 'package:caoky/services/auth_service.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_stripe/flutter_stripe.dart';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart';
import 'package:shared_preferences/shared_preferences.dart';

class DepositPage extends StatefulWidget {
  const DepositPage({super.key});

  @override
  _DepositPageState createState() => _DepositPageState();
}

class _DepositPageState extends State<DepositPage> {
  final _amountController = TextEditingController();
  final _formKey = GlobalKey<FormState>();
  bool _isProcessing = false;

  // *** QUAN TRỌNG: Đảm bảo URL này khớp với API của bạn ***
  final String _apiBaseUrl = "http://57.155.76.74/api/v1";

  // --- LOGIC TÁI SỬ DỤNG VÀ ĐIỀU CHỈNH TỪ PAYMENT_SCREEN ---
  final AuthRepository _authRepository = AuthRepository();
  Future<String?> _getAuthToken() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = _authRepository.getValidAccessToken();
      if (token != null) {
        log("Retrieved Auth Token for Deposit");
        return token;
      } else {
        log("Auth Token not found in SharedPreferences.");
        _showErrorDialog("Bạn chưa đăng nhập. Vui lòng đăng nhập để tiếp tục.");
        return null;
      }
    } catch (e) {
      log("Error fetching token from SharedPreferences: $e");
      _showErrorDialog("Lỗi truy xuất thông tin xác thực.");
      return null;
    }
  }

  // Bước 1: Tạo Payment Intent cho việc nạp tiền (KHÔNG CẦN ticket_id)
  Future<Map<String, dynamic>?> _createStripePaymentIntent(int amount) async {
    final prefs = await SharedPreferences.getInstance();
    final token =
        _authRepository
            .getValidAccessToken(); // Assuming you also store the user ID in shared preferences upon login
    final String? userId = prefs.getString('userId');

    if (token == null || userId == null) {
      throw Exception('Authentication token or User ID not found.');
    }

    final Map<String, dynamic> payload = {
      "amount": amount,
      "currency": "vnd",
      "customer_id": userId,
      "ticket_id": "Deposit",
      // Backend của bạn có thể cần customer_id hoặc không.
      // Nếu không cần, bạn có thể bỏ dòng này.
      // "customer_id": "some_customer_id_if_needed",
    };

    try {
      final String jsonPayload = jsonEncode(payload);
      log(
        "--- Creating Payment Intent for Deposit ($_apiBaseUrl/stripe/create-payment-intent) ---",
      );
      log(jsonPayload);

      final response = await http.post(
        Uri.parse("$_apiBaseUrl/stripe/create-payment-intent"),
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer $token",
        },
        body: jsonPayload,
      );

      log(
        "Create Payment Intent API Response: ${response.statusCode} ${response.body}",
      );

      if (response.statusCode == 200) {
        final responseData = jsonDecode(response.body) as Map<String, dynamic>;
        if (responseData['data']?['client_secret'] != null) {
          return responseData['data'];
        } else {
          _showErrorDialog("Dữ liệu thanh toán không hợp lệ từ máy chủ.");
          return null;
        }
      } else {
        _showErrorDialog(
          "Lỗi tạo phiên thanh toán: ${response.reasonPhrase} (${response.statusCode}).",
        );
        return null;
      }
    } catch (e) {
      log("Exception calling create payment intent API: $e");
      _showErrorDialog("Lỗi kết nối khi tạo phiên thanh toán: ${e.toString()}");
      return null;
    }
  }

  // Bước 2: Khởi tạo và hiển thị Stripe Payment Sheet
  Future<void> _initializeAndPresentStripeSheet(String clientSecret) async {
    try {
      await Stripe.instance.initPaymentSheet(
        paymentSheetParameters: SetupPaymentSheetParameters(
          paymentIntentClientSecret: clientSecret,
          merchantDisplayName: "CAOKY Nạp Tiền",
          style: ThemeMode.light,
        ),
      );
      await _presentStripePaymentSheet(clientSecret);
    } on StripeException catch (e) {
      log("StripeException during init: ${e.error.message}");
      _showErrorDialog(
        "Lỗi Stripe khi khởi tạo: ${e.error.localizedMessage ?? e.error.message}",
      );
    } catch (e) {
      log("Error initializing Stripe: $e");
      _showErrorDialog("Lỗi khởi tạo thanh toán: ${e.toString()}");
    }
  }

  // Bước 3: Xử lý kết quả từ Payment Sheet
  Future<void> _presentStripePaymentSheet(String clientSecret) async {
    try {
      await Stripe.instance.presentPaymentSheet();
      final paymentIntent = await Stripe.instance.retrievePaymentIntent(
        clientSecret,
      );

      if (paymentIntent.status == PaymentIntentsStatus.Succeeded) {
        log('Stripe payment successful! Amount: ${paymentIntent.amount}');
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text(
              "Thanh toán qua Stripe thành công! Đang cập nhật số dư...",
            ),
          ),
        );
        // Bước 4: Gọi API backend để xác nhận nạp tiền
        await _confirmDepositOnBackend(paymentIntent.amount.toInt());
      } else {
        log('Payment not successful. Status: ${paymentIntent.status}');
        _showInfoDialog("Thanh toán không thành công hoặc đã bị hủy.");
      }
    } on StripeException catch (e) {
      log("StripeException during present: ${e.error.message}");
      if (e.error.code != FailureCode.Canceled) {
        _showErrorDialog(
          "Lỗi thanh toán Stripe: ${e.error.localizedMessage ?? e.error.message}",
        );
      } else {
        _showInfoDialog("Bạn đã hủy thanh toán.");
      }
    } catch (e) {
      log("Generic error in presentPaymentSheet: $e");
      _showErrorDialog(
        "Đã xảy ra lỗi không xác định trong quá trình thanh toán.",
      );
    }
  }

  // *** BƯỚC QUAN TRỌNG NHẤT: GỌI ENDPOINT /DEPOSIT ***
  // Bước 4: Gọi API backend để xác nhận nạp tiền vào tài khoản
  Future<void> _confirmDepositOnBackend(int amount) async {
    final String? token = await _getAuthToken();
    if (token == null) return;

    // Backend yêu cầu Amount và Currency
    final Map<String, dynamic> payload = {
      "amount": amount, // Số tiền đã thanh toán thành công từ Stripe
      "currency": "VND",
    };

    try {
      final String jsonPayload = jsonEncode(payload);
      log(
        "--- Confirming Deposit on Backend ($_apiBaseUrl/accounts/deposit) ---",
      );
      log(jsonPayload);

      final response = await http.post(
        Uri.parse("$_apiBaseUrl/accounts/deposit"),
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer $token",
        },
        body: jsonPayload,
      );

      log(
        "Backend Deposit API Response: ${response.statusCode} ${response.body}",
      );

      if (response.statusCode == 200) {
        log("Deposit confirmed successfully on backend.");
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("Nạp tiền thành công!"),
            backgroundColor: Colors.green,
            duration: Duration(seconds: 3),
          ),
        );
        // Trở về màn hình chính và báo hiệu thành công
        if (mounted) {
          Navigator.of(context).pop(true);
        }
      } else {
        log("Failed to confirm deposit on backend.");
        _showErrorDialog(
          "Thanh toán Stripe thành công nhưng có lỗi khi cập nhật số dư. Vui lòng liên hệ hỗ trợ. Chi tiết: ${response.body}",
        );
      }
    } catch (e) {
      log("Exception calling backend deposit API: $e");
      _showErrorDialog(
        "Thanh toán Stripe thành công nhưng có lỗi kết nối khi cập nhật số dư. Vui lòng liên hệ hỗ trợ: ${e.toString()}",
      );
    }
  }

  // Hàm chính điều phối toàn bộ quá trình
  Future<void> _initiateDepositProcess() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    // Ẩn bàn phím
    FocusScope.of(context).unfocus();

    final amount = int.tryParse(_amountController.text.replaceAll(',', ''));
    if (amount == null || amount <= 0) {
      _showErrorDialog("Vui lòng nhập một số tiền hợp lệ.");
      return;
    }

    setState(() {
      _isProcessing = true;
    });

    // B1: Tạo Payment Intent
    final paymentIntentData = await _createStripePaymentIntent(amount);

    if (mounted && paymentIntentData?['client_secret'] != null) {
      // B2 & B3: Hiển thị Stripe Sheet và xử lý kết quả
      await _initializeAndPresentStripeSheet(
        paymentIntentData!['client_secret'],
      );
    }

    if (mounted) {
      setState(() {
        _isProcessing = false;
      });
    }
  }

  // --- CÁC WIDGET HỖ TRỢ VÀ GIAO DIỆN ---

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  void _showErrorDialog(String message) {
    if (!mounted) return;
    showDialog(
      context: context,
      builder:
          (ctx) => AlertDialog(
            title: const Text("Đã xảy ra lỗi"),
            content: Text(message),
            actions: [
              TextButton(
                child: const Text("OK"),
                onPressed: () => Navigator.of(ctx).pop(),
              ),
            ],
          ),
    );
  }

  void _showInfoDialog(String message) {
    if (!mounted) return;
    showDialog(
      context: context,
      builder:
          (ctx) => AlertDialog(
            title: const Text("Thông báo"),
            content: Text(message),
            actions: [
              TextButton(
                child: const Text("OK"),
                onPressed: () => Navigator.of(ctx).pop(),
              ),
            ],
          ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: const Text("Nạp tiền vào tài khoản"),
        backgroundColor: Colors.blueAccent,
        foregroundColor: Colors.white,
        centerTitle: true,
      ),
      body: GestureDetector(
        onTap:
            () =>
                FocusScope.of(
                  context,
                ).unfocus(), // Ẩn bàn phím khi chạm ra ngoài
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  "Nhập số tiền muốn nạp",
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 16),
                TextFormField(
                  controller: _amountController,
                  keyboardType: TextInputType.number,
                  inputFormatters: [
                    FilteringTextInputFormatter.digitsOnly,
                    // Format tiền tệ cho dễ nhìn
                    TextInputFormatter.withFunction((oldValue, newValue) {
                      if (newValue.text.isEmpty) return newValue;
                      final number = int.parse(newValue.text);
                      final formatter = NumberFormat("#,###");
                      final newString = formatter.format(number);
                      return TextEditingValue(
                        text: newString,
                        selection: TextSelection.collapsed(
                          offset: newString.length,
                        ),
                      );
                    }),
                  ],
                  decoration: InputDecoration(
                    labelText: "Số tiền",
                    hintText: "VD: 50.000",
                    suffixText: "VND",

                    // Viền các trạng thái
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(10),
                      borderSide: BorderSide(color: Colors.grey),
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(10),
                      borderSide: BorderSide(color: Colors.blue, width: 2),
                    ),
                    errorBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(10),
                      borderSide: BorderSide(color: Colors.red),
                    ),
                    focusedErrorBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(10),
                      borderSide: BorderSide(color: Colors.red, width: 2),
                    ),

                    // Màu chữ và label khi focus
                    labelStyle: TextStyle(color: Colors.grey), // mặc định
                    floatingLabelStyle: TextStyle(
                      color: Colors.blue,
                    ), // khi focus
                    hintStyle: TextStyle(
                      color: Colors.grey.shade400,
                    ), // màu gợi ý
                    suffixStyle: TextStyle(color: Colors.blue), // màu "VND"
                  ),

                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Vui lòng nhập số tiền';
                    }
                    final amount = int.tryParse(value.replaceAll(',', ''));
                    if (amount == null || amount < 10000) {
                      // Ví dụ: yêu cầu nạp tối thiểu 10.000
                      return 'Số tiền nạp tối thiểu là 10.000đ';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 32),
                const Text(
                  "Chọn phương thức thanh toán",
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.blueAccent, width: 1.5),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Row(
                    children: [
                      Image.asset(
                        'assets/images/stripe_logo.png', // **LƯU Ý:** Bạn cần thêm logo Stripe vào assets
                        height: 30,
                        errorBuilder:
                            (context, error, stackTrace) => const Icon(
                              Icons.credit_card,
                              color: Colors.blueAccent,
                              size: 30,
                            ),
                      ),
                      const SizedBox(width: 12),
                      const Expanded(
                        child: Text(
                          "Thẻ tín dụng/ghi nợ (Stripe)",
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                      const Icon(Icons.check_circle, color: Colors.blueAccent),
                    ],
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  "Các phương thức khác sẽ sớm được cập nhật.",
                  style: TextStyle(
                    color: Colors.grey,
                    fontStyle: FontStyle.italic,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
      bottomNavigationBar: Padding(
        padding: const EdgeInsets.all(16.0),
        child: ElevatedButton(
          onPressed: _isProcessing ? null : _initiateDepositProcess,
          style: ElevatedButton.styleFrom(
            backgroundColor: Colors.orange,
            padding: const EdgeInsets.symmetric(vertical: 16),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(10),
            ),
          ),
          child:
              _isProcessing
                  ? const SizedBox(
                    height: 24,
                    width: 24,
                    child: CircularProgressIndicator(
                      color: Colors.white,
                      strokeWidth: 3,
                    ),
                  )
                  : const Text(
                    "Nạp tiền",
                    style: TextStyle(
                      fontSize: 18,
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
        ),
      ),
    );
  }
}
