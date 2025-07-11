import 'package:caoky/global/toast.dart';
import 'package:caoky/models/user/login_request.dart';
import 'package:caoky/models/user/login_response.dart';
// ApiUserService direct import might not be needed if all calls go via AuthRepository
// import 'package:caoky/services/api_user_service.dart';
import 'package:caoky/services/auth_service.dart'; // Your AuthRepository
import 'package:caoky/ui/login/forgot_password_page.dart';
import 'package:caoky/ui/login/register_page.dart';
import 'package:caoky/ui/main_page.dart';
// import 'package:caoky/ui/admin/admin_main_page.dart'; // Example for ROLE_ADMIN
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:lottie/lottie.dart';
// SharedPreferences direct import not needed here

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
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
    // Pre-fill for testing, you can remove this
    _usernameController.text = 'uchihaha3169@gmail.com';
    _passwordController.text = 'Password123!';
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
        rememberMe: _rememberMe,
      );

      final LoginResponse response = await _authRepository.login(request);

      if (response.code == 200) {
        _onLoginSuccess(response);
      } else {
        ToastUtils.show(response.message ?? "Đăng nhập không thành công (mã lỗi ${response.code})");
        if (mounted) setState(() => _isLoading = false);
      }
    } catch (e) {
      String errorMessage = "Lỗi không xác định. Vui lòng thử lại.";
      if (e is DioException) {
        // ... (Error handling logic remains the same)
      } else {
        errorMessage = "Lỗi đăng nhập: ${e.toString()}";
      }
      ToastUtils.show(errorMessage);
      if (mounted) setState(() => _isLoading = false);
    }
  }

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

    print("Login successful. Role: ${loginData.user!.role}, RememberMe: $_rememberMe");

    bool userInfoFetched = await _authRepository.fetchAndSaveDetailedUserInfo(loginData.accessToken);

    if (mounted) {
      if (userInfoFetched) {
        _navigateToMainPage(loginData.user!.role);
      } else {
        ToastUtils.show("Không thể xác thực thông tin chi tiết người dùng.");
        await _authRepository.logout();
        setState(() => _isLoading = false);
      }
    }
  }

  void _navigateToMainPage(String role) {
    if (!mounted) return;
    
    Widget destinationPage;
    if (role == 'ROLE_CUSTOMER' || role == 'ROLE_ADMIN') {
        destinationPage = MainPage();
    } else {
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


  void _loginWithGoogle() {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Đăng nhập với Google... (chưa triển khai)')),
    );
  }

  void _loginWithFacebook() {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Đăng nhập với Facebook... (chưa triển khai)')),
    );
  }

  @override
  Widget build(BuildContext context) {
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
                        'assets/images/bus.json',
                        width: 180,
                        height: 180,
                        fit: BoxFit.contain,
                      ),
                      SizedBox(height: 16),
                      Text(
                        'Chào mừng bạn',
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
                      SizedBox(height: 24), // Increased spacing
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: _isLoading ? null : _login,
                          style: ElevatedButton.styleFrom(
                            backgroundColor:
                                Colors.blue.shade700, // Darker blue
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
                      SizedBox(height: 20),
                      Row(
                        children: [
                          Expanded(
                            child: Divider(
                              color: Colors.grey.shade300,
                              thickness: 1,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 8.0,
                            ),
                            child: Text(
                              'Đăng nhập với',
                              style: TextStyle(
                                color: Colors.grey,
                                fontSize: 14,
                              ),
                            ),
                          ),
                          Expanded(
                            child: Divider(
                              color: Colors.grey.shade300,
                              thickness: 1,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: 20),
                      Row(
                        mainAxisAlignment:
                            MainAxisAlignment.center, // Centered buttons
                        children: [
                          _buildSocialLoginButton(
                            'assets/images/google_icon.png',
                            'Google',
                            _loginWithGoogle,
                          ),
                          SizedBox(width: 24), // Spacing between buttons
                          _buildSocialLoginButton(
                            'assets/images/facebook_icon.png',
                            'Facebook',
                            _loginWithFacebook,
                          ),
                        ],
                      ),
                      SizedBox(height: 24), // Increased spacing
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "Bạn chưa có tài khoản? ",
                            style: TextStyle(color: Colors.grey),
                          ),
                          GestureDetector(
                            onTap: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (context) => RegisterPage(),
                                ),
                              );
                            },
                            child: Text(
                              'Đăng ký',
                              style: TextStyle(
                                color: Colors.blue.shade700,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: 20), // Bottom padding
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

  Widget _buildSocialLoginButton(
    String imagePath,
    String text,
    VoidCallback onPressed,
  ) {
    return OutlinedButton(
      onPressed: onPressed,
      style: OutlinedButton.styleFrom(
        padding: EdgeInsets.symmetric(
          horizontal: 20,
          vertical: 12,
        ), // Adjusted padding
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(30)),
        side: BorderSide(color: Colors.grey.shade300),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Image.asset(imagePath, width: 20, height: 20),
          SizedBox(width: 10),
          Text(
            text,
            style: TextStyle(
              color: Colors.black87,
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
