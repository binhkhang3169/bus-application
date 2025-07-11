import 'package:caoky/global/toast.dart';
import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import '../../services/api_user_service.dart';

class ForgotPasswordScreen extends StatefulWidget {
  @override
  _ForgotPasswordScreenState createState() => _ForgotPasswordScreenState();
}

class _ForgotPasswordScreenState extends State<ForgotPasswordScreen> {
  final TextEditingController _controller = TextEditingController();
  String? _errorMessage;
  bool isButtonEnabled = false;
  bool isLoading = false;

  final _emailFocusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    _emailFocusNode.addListener(() {
      setState(() {});
    });

    _controller.addListener(() {
      _validateInput(_controller.text);
    });
  }

  bool isValidEmail(String email) {
    final RegExp emailRegex = RegExp(
      r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$",
    );
    return emailRegex.hasMatch(email);
  }

  void _validateInput(String input) {
    String? error;

    if (input.isEmpty) {
      error = "Vui lòng nhập email!";
    } else if (!isValidEmail(input)) {
      error = "Email không hợp lệ!";
    }

    setState(() {
      _errorMessage = error;
      isButtonEnabled = (error == null);
    });
  }

  void sendResetPassword(String email) async {
    final dio = Dio();
    final apiService = ApiUserService(dio);
    setState(() {
      isLoading = true;
    });
    try {
      await apiService.sendResetPassword(email);
      ToastUtils.show("Email đặt lại mật khẩu đã được gửi!");
      setState(() {
        isLoading = false;
      });
      Navigator.pop(context);
    } catch (e) {
      ToastUtils.show("Không thể gửi email!");
    } finally {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.blue,
        foregroundColor: Colors.white,
        leading: IconButton(
          icon: Icon(Icons.arrow_back),
          onPressed: () => Navigator.pop(context),
        ),
        title: Text(
          "Khôi phục mật khẩu",
          style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
        elevation: 7,
        shadowColor: Colors.black.withOpacity(0.3),
      ),
      body: Container(
        color: Colors.white,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextFormField(
                controller: _controller,
                focusNode: _emailFocusNode,
                decoration: InputDecoration(
                  labelText: "Email",
                  prefixIcon: Icon(
                    Icons.email,
                    color: _emailFocusNode.hasFocus ? Colors.blue : Colors.grey,
                  ),
                  labelStyle: TextStyle(
                    color:
                        _emailFocusNode.hasFocus ? Colors.blue : Colors.black,
                  ),
                  floatingLabelBehavior: FloatingLabelBehavior.auto,
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(30),
                    borderSide: BorderSide(color: Colors.blue, width: 2),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(30),
                    borderSide: BorderSide(color: Colors.grey, width: 1),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(30),
                  ),
                  errorText: _errorMessage,
                ),
                keyboardType: TextInputType.emailAddress,
                onChanged: _validateInput,
                onTapOutside: (event) {
                  _emailFocusNode.unfocus();
                },
              ),
              SizedBox(height: 20),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed:
                      isButtonEnabled && !isLoading
                          ? () => sendResetPassword(_controller.text.trim())
                          : null,
                  style: ElevatedButton.styleFrom(
                    minimumSize: Size(double.infinity, 55),
                    padding: EdgeInsets.symmetric(vertical: 16),
                    textStyle: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                    backgroundColor:
                        isButtonEnabled ? Colors.blue : Colors.grey[300],
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(30),
                    ),
                  ),
                  child:
                      isLoading
                          ? SizedBox(
                            height: 24,
                            width: 24,
                            child: CircularProgressIndicator(
                              color: Colors.white,
                              strokeWidth: 3,
                            ),
                          )
                          : Text("Xác nhận"),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    _emailFocusNode.dispose();
    super.dispose();
  }
}
