import 'package:caoky/services/account_service.dart';
import 'package:flutter/material.dart';

class CreateAccountPage extends StatefulWidget {
  @override
  _CreateAccountPageState createState() => _CreateAccountPageState();
}

class _CreateAccountPageState extends State<CreateAccountPage> {
  // Trạng thái cho checkbox và nút bấm
  bool _isTermsAgreed = false;
  bool _isLoading = false;
  String? _errorMessage;

  final AccountService _accountService = AccountService();

  // Hàm được gọi khi nhấn nút tạo tài khoản
  Future<void> _submitCreateAccount() async {
    // Chỉ thực hiện khi đã đồng ý điều khoản
    if (!_isTermsAgreed) return;

    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    try {
      // Gọi service để tạo tài khoản với các giá trị mặc định
      await _accountService.createAccount(
        ownerName: "",
        balance: 0,
        currency: "VND",
      );

      // Quay lại trang trước và báo hiệu thành công
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Tạo tài khoản thành công!'),
            backgroundColor: Colors.green,
          ),
        );
        Navigator.pop(context, true);
      }
    } catch (e) {
      setState(() {
        _errorMessage = e.toString().replaceFirst("Exception: ", "");
      });
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: Text("Điều khoản & Điều kiện"),
        backgroundColor: Colors.blueAccent,
        foregroundColor: Colors.white,
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Khung chứa nội dung điều khoản
            Expanded(
              child: Container(
                padding: const EdgeInsets.all(12.0),
                decoration: BoxDecoration(
                  border: Border.all(color: Colors.grey.shade300),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: SingleChildScrollView(
                  child: Text(
                    // Thay thế bằng nội dung điều khoản thực tế của bạn
                    '''Chào mừng bạn đến với DACNTT!

1.  **Chấp nhận Điều khoản**: Bằng việc tạo tài khoản, bạn xác nhận đã đọc, hiểu và đồng ý bị ràng buộc bởi các điều khoản và điều kiện này.

2.  **Thông tin Tài khoản**: Tài khoản của bạn sẽ được tạo với các thông tin sau:
    -   **Tên chủ tài khoản**: Sẽ được liên kết với ID người dùng của bạn.
    -   **Số dư ban đầu**: 0 (Không) đồng.
    -   **Loại tiền tệ**: VND (Việt Nam Đồng).

3.  **Trách nhiệm của Người dùng**: Bạn có trách nhiệm bảo mật thông tin đăng nhập và mọi hoạt động diễn ra dưới tài khoản của mình.

4.  **Giao dịch**: Mọi giao dịch nạp tiền, thanh toán sẽ được ghi lại trong lịch sử giao dịch và tuân thủ các quy định của DACNTT.

5.  **Chấm dứt tài khoản**: Chúng tôi có quyền tạm ngưng hoặc chấm dứt tài khoản của bạn nếu phát hiện vi phạm các điều khoản này.

Cảm ơn bạn đã sử dụng dịch vụ của chúng tôi!
                    ''',
                    style: TextStyle(
                      fontSize: 14,
                      height: 1.5,
                      color: Colors.grey[700],
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 16),
            // Checkbox đồng ý điều khoản
            Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Checkbox(
                  value: _isTermsAgreed,
                  onChanged: (bool? value) {
                    setState(() {
                      _isTermsAgreed = value ?? false;
                    });
                  },
                  checkColor:
                      Colors
                          .white, // Dấu kiểm màu trắng để tương phản với nền xanh
                  activeColor:
                      Colors.blueAccent, // Màu nền khi Checkbox được chọn
                  fillColor: MaterialStateProperty.resolveWith((states) {
                    if (states.contains(MaterialState.selected)) {
                      return Colors.blueAccent; // Màu xanh khi được chọn
                    }
                    return Colors
                        .grey
                        .shade300; // Màu xám nhạt khi không được chọn
                  }),
                ),
                Expanded(
                  child: Text(
                    "Tôi đã đọc, hiểu và đồng ý với các điều khoản trên.",
                    style: TextStyle(fontSize: 14),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            // Hiển thị lỗi nếu có
            if (_errorMessage != null)
              Padding(
                padding: const EdgeInsets.only(bottom: 16.0),
                child: Text(
                  _errorMessage!,
                  style: TextStyle(color: Colors.red, fontSize: 14),
                  textAlign: TextAlign.center,
                ),
              ),
            // Nút Tạo tài khoản
            ElevatedButton(
              // Nút bị vô hiệu hóa khi đang tải hoặc chưa đồng ý điều khoản
              onPressed:
                  _isLoading || !_isTermsAgreed ? null : _submitCreateAccount,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.blueAccent,
                disabledBackgroundColor: Colors.grey,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child:
                  _isLoading
                      ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                      : const Text(
                        "Xác nhận tạo tài khoản",
                        style: TextStyle(fontSize: 16, color: Colors.white),
                      ),
            ),
          ],
        ),
      ),
    );
  }
}
