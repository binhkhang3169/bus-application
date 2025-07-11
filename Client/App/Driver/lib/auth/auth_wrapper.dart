import 'package:flutter/material.dart';
import 'package:taixe/services/auth_service.dart';

import 'package:taixe/main_page.dart';
import 'package:taixe/ui/login/login_page.dart';

class AuthWrapper extends StatefulWidget {
  const AuthWrapper({super.key});

  @override
  State<AuthWrapper> createState() => _AuthWrapperState();
}

class _AuthWrapperState extends State<AuthWrapper> {
  final AuthRepository _authRepository = AuthRepository();

  @override
  void initState() {
    super.initState();
    _decideNextScreen();
  }

  Future<void> _decideNextScreen() async {
    // Sử dụng phương thức attemptAutoLogin đã được tái cấu trúc từ repository
    final role = await _authRepository.attemptAutoLogin();

    // Đảm bảo widget vẫn còn trong cây widget trước khi điều hướng
    if (!mounted) return;

    Widget destinationPage;
    if (role != null) {
      print(
        "AuthWrapper: Tự động đăng nhập thành công. Vai trò: $role. Điều hướng đến MainPage.",
      );
      // Bạn có thể thêm logic cho các vai trò khác nhau ở đây nếu cần
      destinationPage = MainPage();
    } else {
      print(
        "AuthWrapper: Không tìm thấy phiên hợp lệ. Điều hướng đến LoginScreen.",
      );
      destinationPage = LoginScreen();
    }

    // Sử dụng pushAndRemoveUntil để thay thế màn hình loading bằng màn hình đích
    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (context) => destinationPage),
      (route) => false, // Predicate này xóa tất cả các route trước đó
    );
  }

  @override
  Widget build(BuildContext context) {
    // Hiển thị một chỉ báo loading trong khi kiểm tra trạng thái xác thực
    return const Scaffold(
      backgroundColor: Colors.white,
      body: Center(
        child: CircularProgressIndicator(
          color: Colors.blue, // Màu sắc phù hợp với app driver
        ),
      ),
    );
  }
}
