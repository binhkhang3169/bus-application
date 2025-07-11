import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:lottie/lottie.dart';
import 'package:taixe/global/toast.dart';
import 'package:taixe/main_page.dart';
import 'package:taixe/models/user/login_request.dart';
import 'package:taixe/models/user/login_response.dart';
import 'package:taixe/services/auth_service.dart';
import 'package:taixe/ui/login/forgot_password_page.dart';

class LoginScreen extends StatefulWidget {
  // Chuyển sang const constructor vì không còn tham số
  const LoginScreen({super.key});

  @override
  _LoginScreenState createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  late AuthRepository _authRepository;

  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;
  bool _rememberMe = false;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _authRepository = AuthRepository();
    // KHÔNG còn `_checkLoginStatusAndAttemptAutoLogin()` ở đây nữa.
    // AuthWrapper đã xử lý việc này.
  }

  // ĐÃ XÓA: Phương thức `_checkLoginStatusAndAttemptAutoLogin` không còn cần thiết ở đây.

  void _onLoginSuccess(LoginResponse loginResponse) async {
    if (loginResponse.data == null || loginResponse.data!.user == null) {
      ToastUtils.show("Dữ liệu đăng nhập không hợp lệ.");
      if (mounted) setState(() => _isLoading = false);
      return;
    }

    final loginData = loginResponse.data!;
    await _authRepository.saveLoginData(
      accessToken: loginData.accessToken,
      refreshToken: loginData.refreshToken,
      role: loginData.user!.role,
      userId: loginData.user!.id.toString(),
      username: loginData.user!.username,
      rememberMe: _rememberMe,
    );

    print(
      "Đăng nhập thành công. Vai trò: ${loginData.user!.role}, UserID: ${loginData.user!.id}, Ghi nhớ: $_rememberMe",
    );

    // Sau khi lưu dữ liệu, lấy thông tin chi tiết người dùng
    bool userInfoFetched = await _authRepository.fetchAndSaveDetailedUserInfo(
      loginData.accessToken,
    );

    if (mounted) {
      if (userInfoFetched) {
        print("Đã lấy và lưu thông tin chi tiết người dùng thành công sau khi đăng nhập.");
        _navigateToMainPage(loginData.user!.role);
      } else {
        print(
          "Không thể lấy thông tin chi tiết người dùng sau khi đăng nhập.",
        );
        ToastUtils.show("Không thể xác thực thông tin chi tiết người dùng.");
        await _authRepository.logout();
        setState(() => _isLoading = false);
      }
    }
  }

  void _navigateToMainPage(String role) {
    if (!mounted) return;

    Widget destinationPage;
    if (role == 'ROLE_DRIVER') {
      destinationPage = MainPage();
    } else {
      print("Vai trò không xác định: $role. Không thể điều hướng.");
      ToastUtils.show("Vai trò người dùng không xác định: $role");
      if (mounted) setState(() => _isLoading = false);
      return;
    }

    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (context) => destinationPage),
      (route) => false,
    );
  }

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void _login() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    if (!mounted) return;
    setState(() {
      _isLoading = true;
    });

    try {
      final request = LoginRequest(
        username: _usernameController.text.trim(),
        password: _passwordController.text.trim(),
        rememberMe: _rememberMe, // rememberMe vẫn được gửi đi nếu API của bạn cần
      );

      final LoginResponse response = await _authRepository.login(request);

      if (response.code == 200) {
        _onLoginSuccess(response);
      } else {
        ToastUtils.show(
          response.message ?? "Đăng nhập không thành công (mã lỗi ${response.code})",
        );
        if (mounted) setState(() => _isLoading = false);
      }
    } catch (e) {
      String errorMessage = "Lỗi không xác định. Vui lòng thử lại.";
      if (e is DioException) {
        if (e.response?.data is Map && e.response!.data.containsKey("message")) {
          errorMessage = e.response!.data["message"];
        } else if (e.response?.data is String && (e.response!.data as String).isNotEmpty) {
          errorMessage = e.response!.data as String;
        } else if (e.response?.statusCode != null) {
          errorMessage = "Lỗi ${e.response!.statusCode}: Không thể kết nối đến máy chủ.";
        } else {
          errorMessage = "Lỗi mạng hoặc máy chủ không phản hồi.";
        }
      } else {
        errorMessage = "Lỗi đăng nhập: ${e.toString()}";
      }
      ToastUtils.show(errorMessage);
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    // Toàn bộ phần UI (build method) giữ nguyên như cũ, không cần thay đổi.
    // ... Dán lại toàn bộ Widget build(BuildContext context) của bạn vào đây ...
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: Stack(
          children: [
            Center(
              child: SingleChildScrollView(
                padding: EdgeInsets.symmetric(horizontal: 24.0),
                child: Form(
                  key: _formKey,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Lottie.asset(
                        'assets/images/driver.json',
                        width: 180,
                        height: 180,
                        fit: BoxFit.contain,
                      ),
                      SizedBox(height: 16),
                      Text(
                        'Thượng lộ bình an',
                        style: TextStyle(
                          fontFamily: 'Pacifico',
                          fontSize: 32,
                          color: Colors.black,
                        ),
                      ),
                      SizedBox(height: 32),
                      TextFormField(
                        controller: _usernameController,
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Tên đăng nhập',
                          prefixIcon: Icon(
                            Icons.person,
                            color: Colors.grey,
                            size: 20,
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30),
                            borderSide: BorderSide.none,
                          ),
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Vui lòng nhập tên người dùng';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      TextFormField(
                        controller: _passwordController,
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Mật khẩu',
                          prefixIcon: Icon(
                            Icons.lock,
                            color: Colors.grey,
                            size: 20,
                          ),
                          suffixIcon: IconButton(
                            icon: Icon(
                              _obscurePassword
                                  ? Icons.visibility
                                  : Icons.visibility_off,
                              color: Colors.grey,
                              size: 20,
                            ),
                            onPressed: () {
                              if (mounted)
                                setState(
                                  () => _obscurePassword = !_obscurePassword,
                                );
                            },
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30),
                            borderSide: BorderSide.none,
                          ),
                        ),
                        obscureText: _obscurePassword,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Vui lòng nhập mật khẩu';
                          }
                          if (value.length < 6) {
                            return 'Mật khẩu phải có ít nhất 6 ký tự';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 8),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Row(
                            children: [
                              Transform.scale(
                                scale: 0.8,
                                child: Checkbox(
                                  value: _rememberMe,
                                  onChanged: (value) {
                                    if (mounted)
                                      setState(
                                        () => _rememberMe = value ?? false,
                                      );
                                  },
                                  activeColor: Colors.blue.shade700,
                                ),
                              ),
                              Text(
                                'Nhớ mật khẩu',
                                style: TextStyle(
                                  color: Colors.grey,
                                  fontSize: 14,
                                ),
                              ),
                            ],
                          ),
                          GestureDetector(
                            onTap: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (context) => ForgotPasswordScreen(),
                                ),
                              );
                            },
                            child: Text(
                              'Quên mật khẩu?',
                              style: TextStyle(
                                color: Colors.blue.shade700,
                                fontSize: 14,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: 24),
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: _isLoading ? null : _login,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.blue.shade700,
                            padding: EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(30),
                            ),
                            elevation: 5,
                          ),
                          child: Text(
                            'ĐĂNG NHẬP',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              color: Colors.white,
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
            if (_isLoading)
              Positioned.fill(
                child: Container(
                  color: Colors.black.withOpacity(0.3),
                  child: Center(
                    child: CircularProgressIndicator(
                      color: Colors.white,
                      strokeWidth: 3,
                    ),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}