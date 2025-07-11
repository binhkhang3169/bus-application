import 'package:caoky/global/toast.dart';
import 'package:caoky/services/api_user_service.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';

class ChangePasswordScreen extends StatefulWidget {
  final String? token;

  const ChangePasswordScreen({super.key, this.token});
  @override
  _ChangePasswordScreenState createState() => _ChangePasswordScreenState();
}

class _ChangePasswordScreenState extends State<ChangePasswordScreen> {
  bool _isOldPasswordVisible = false;
  bool _isNewPasswordVisible = false;
  bool _isConfirmPasswordVisible = false;
  bool _isLoading = false;

  final _oldPasswordFocusNode = FocusNode();
  final _newPasswordFocusNode = FocusNode();
  final _confirmPasswordFocusNode = FocusNode();

  final TextEditingController _oldPasswordController = TextEditingController();
  final TextEditingController _newPasswordController = TextEditingController();
  final TextEditingController _confirmPasswordController =
      TextEditingController();

  String? _confirmPasswordError;

  final _formKey = GlobalKey<FormState>();

  Future<void> _changePassword(
    String token,
    String oldPass,
    String newPass,
  ) async {
    setState(() {
      _isLoading = true;
    });

    final dio = Dio();
    final apiService = ApiUserService(dio);

    try {
      final response = await apiService.changePassword(
        "Bearer $token",
        oldPass,
        newPass,
      );

      if (response.code == 200) {
        ToastUtils.show(response.message);
        Navigator.pop(context);
      } else {
        ToastUtils.show(response.message);
      }
      setState(() {
        _isLoading = false;
      });
    } catch (e) {
      if (e is DioException) {
        ToastUtils.show(e.response?.data['message'] ?? "Lỗi server");
      } else {
        ToastUtils.show("Lỗi kết nối đến server: $e");
      }
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  void initState() {
    super.initState();
    _oldPasswordFocusNode.addListener(() {
      setState(() {});
    });
    _newPasswordFocusNode.addListener(() {
      setState(() {});
    });
    _confirmPasswordFocusNode.addListener(() {
      setState(() {});
    });

    _confirmPasswordController.addListener(() {
      setState(() {
        _confirmPasswordError = _validateConfirmPassword(
          _confirmPasswordController.text.trim(),
        );
      });
    });

    _newPasswordController.addListener(() {
      setState(() {
        _confirmPasswordError = _validateConfirmPassword(
          _confirmPasswordController.text.trim(),
        );
      });
    });
  }

  String? _validateConfirmPassword(String confirmPassword) {
    if (confirmPassword.isEmpty) {
      return null;
    }
    if (confirmPassword != _newPasswordController.text.trim()) {
      return "Mật khẩu không khớp";
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Đổi mật khẩu", style: TextStyle(color: Colors.white)),
        centerTitle: true,
        backgroundColor: Colors.blueAccent,
        iconTheme: IconThemeData(color: Colors.white),
      ),
      body: Container(
        color: Colors.white,
        child: Padding(
          padding: EdgeInsets.all(20),
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                _buildPasswordField(
                  "Mật khẩu hiện tại",
                  _isOldPasswordVisible,
                  () {
                    setState(
                      () => _isOldPasswordVisible = !_isOldPasswordVisible,
                    );
                  },
                  _oldPasswordFocusNode,
                  _oldPasswordController,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Vui lòng nhập mật khẩu hiện tại';
                    }
                    return null;
                  },
                ),
                SizedBox(height: 15),
                _buildPasswordField(
                  "Mật khẩu mới",
                  _isNewPasswordVisible,
                  () {
                    setState(
                      () => _isNewPasswordVisible = !_isNewPasswordVisible,
                    );
                  },
                  _newPasswordFocusNode,
                  _newPasswordController,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Vui lòng nhập mật khẩu mới';
                    }
                    return null;
                  },
                ),
                SizedBox(height: 15),
                _buildPasswordField(
                  "Xác nhận mật khẩu",
                  _isConfirmPasswordVisible,
                  () {
                    setState(
                      () =>
                          _isConfirmPasswordVisible =
                              !_isConfirmPasswordVisible,
                    );
                  },
                  _confirmPasswordFocusNode,
                  _confirmPasswordController,
                  errorText: _confirmPasswordError,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Vui lòng xác nhận mật khẩu';
                    }
                    return _validateConfirmPassword(value.trim());
                  },
                ),
                SizedBox(height: 30),
                ElevatedButton(
                  onPressed: () {
                    if (_formKey.currentState!.validate()) {
                      _changePassword(
                        widget.token!,
                        _oldPasswordController.text.trim(),
                        _newPasswordController.text.trim(),
                      );
                    }
                  },
                  child: Text(
                    "Xác nhận",
                    style: TextStyle(fontSize: 18, color: Colors.white),
                  ),
                  style: ElevatedButton.styleFrom(
                    minimumSize: Size(double.infinity, 50),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    backgroundColor: Colors.blueAccent,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildPasswordField(
    String label,
    bool isVisible,
    VoidCallback toggleVisibility,
    FocusNode focusNode,
    TextEditingController controller, {
    String? errorText,
    String? Function(String?)? validator,
  }) {
    return TextFormField(
      controller: controller,
      obscureText: !isVisible,
      focusNode: focusNode,
      decoration: InputDecoration(
        labelText: label,
        labelStyle: TextStyle(
          color: focusNode.hasFocus ? Colors.blueAccent : Colors.black,
        ),
        floatingLabelBehavior: FloatingLabelBehavior.auto,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(10),
          borderSide: BorderSide(color: Colors.blueAccent, width: 2),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(10),
          borderSide: BorderSide(color: Colors.grey, width: 1),
        ),
        suffixIcon: IconButton(
          icon: Icon(
            isVisible ? Icons.visibility : Icons.visibility_off,
            color: focusNode.hasFocus ? Colors.blueAccent : Colors.grey,
          ),
          onPressed: toggleVisibility,
        ),
        errorText: errorText,
      ),
      onTapOutside: (event) {
        focusNode.unfocus();
      },
      validator: validator,
    );
  }

  @override
  void dispose() {
    _oldPasswordController.dispose();
    _newPasswordController.dispose();
    _confirmPasswordController.dispose();
    _oldPasswordFocusNode.dispose();
    _newPasswordFocusNode.dispose();
    _confirmPasswordFocusNode.dispose();
    super.dispose();
  }
}
