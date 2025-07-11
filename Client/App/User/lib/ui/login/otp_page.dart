import 'dart:async';
import 'package:caoky/global/toast.dart';
import 'package:caoky/models/user/resend_otp.dart';
import 'package:caoky/models/user/signup_request.dart';
import 'package:caoky/services/api_user_service.dart';
import 'package:caoky/ui/login/login_page.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:pin_code_fields/pin_code_fields.dart';
import 'package:shared_preferences/shared_preferences.dart';

class VerifyOtpScreen extends StatefulWidget {
  @override
  _VerifyOtpScreenState createState() => _VerifyOtpScreenState();
}

class _VerifyOtpScreenState extends State<VerifyOtpScreen> {
  late ApiUserService apiService;

  String username = "";
  String password = "";
  String phoneNumber = "";
  String email = "";
  String fullName = "";
  String address = "";
  String gender = "";
  String otp = "";

  bool isButtonEnabled = false;
  int _timeRemaining = 60;
  late Timer _timer;

  @override
  void initState() {
    super.initState();
    apiService = ApiUserService(Dio());
    _loadUserData();
    _startTimer();
  }

  Future<void> _loadUserData() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    setState(() {
      username = prefs.getString('username') ?? "";
      password = prefs.getString('password') ?? "";
      phoneNumber = prefs.getString('phoneNumber') ?? "";
      email = prefs.getString('email') ?? "";
      fullName = prefs.getString('fullName') ?? "";
      address = prefs.getString('address') ?? "";
      gender = prefs.getString('gender') ?? "";
      print("Name: " + fullName);
    });
  }

  void _startTimer() {
    _timer = Timer.periodic(Duration(seconds: 1), (timer) {
      if (_timeRemaining > 0) {
        setState(() {
          _timeRemaining--;
        });
      } else {
        _timer.cancel();
      }
    });
  }

  void _validateCode(String value) {
    setState(() {
      otp = value;
      isButtonEnabled = otp.length == 6;
    });
  }

  void _resendOtp() {
    if (_timeRemaining == 0) {
      setState(() {
        _timeRemaining = 60;
      });
      _startTimer();

      apiService.resendOtp(ResendOtpRequest(email: email)).then((response) {
        ToastUtils.show(
          response.code == 200
              ? "Mã OTP mới đã được gửi!"
              : "Lỗi khi gửi lại mã OTP!",
        );
      });
    }
  }

  String maskEmail(String email) {
    if (!email.contains("@")) return email;

    List<String> parts = email.split("@");
    String firstPart = parts[0];
    String domain = parts[1];
    String maskedFirstPart = firstPart[0] + '*' * (firstPart.length - 1);
    return "$maskedFirstPart@$domain";
  }

  String _formatTime(int seconds) {
    int minutes = seconds ~/ 60;
    int remainingSeconds = seconds % 60;
    return "${minutes.toString().padLeft(2, '0')}:${remainingSeconds.toString().padLeft(2, '0')}";
  }

  @override
  void dispose() {
    _timer.cancel();
    super.dispose();
  }

  Future<void> _verifyOtp() async {
    try {
      final response = await apiService.verifyOtp(
        SignupRequest(
          username: username,
          password: password,
          phoneNumber: phoneNumber,
          fullName: fullName,
          address: address,
          gender: gender,
          otp: otp,
        ),
      );
      ToastUtils.show(response.message);

      SharedPreferences prefs = await SharedPreferences.getInstance();
      await prefs.clear();
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => LoginScreen()),
      );
    } on DioException catch (e) {
      if (e.response != null && e.response?.statusCode == 400) {
        final message = e.response?.data['message'] ?? "Lỗi không xác định";
        ToastUtils.show(message);
      } else {
        ToastUtils.show(e.message!);
      }
    } catch (e) {
      ToastUtils.show(e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.blue,
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () {
            Navigator.pop(context);
          },
        ),
        title: Text(
          "Xác minh",
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
        elevation: 7,
        shadowColor: Colors.black.withOpacity(0.3),
      ),
      body: SafeArea(
        child: Container(
          color: Colors.white,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  padding: EdgeInsets.symmetric(vertical: 30),
                  child: Center(
                    child: Column(
                      children: [
                        Text(
                          "Chúng tôi đã gửi mã xác nhận qua email",
                          style: TextStyle(fontSize: 16, color: Colors.black),
                        ),
                        SizedBox(height: 12),
                        Text(
                          "${maskEmail(email)}",
                          style: TextStyle(
                            fontSize: 16,
                            color: Colors.black,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                PinCodeTextField(
                  appContext: context,
                  length: 6,
                  keyboardType: TextInputType.number,
                  animationType: AnimationType.fade,
                  pinTheme: PinTheme(
                    shape: PinCodeFieldShape.box,
                    borderRadius: BorderRadius.circular(12),
                    fieldHeight: 50,
                    fieldWidth: 50,
                    activeFillColor: Colors.white,
                    inactiveFillColor: Colors.white,
                    selectedFillColor: Colors.white,
                    inactiveColor: Colors.grey.shade400,
                    selectedColor: Colors.blue.shade300,
                    activeColor: Colors.blue.shade900,
                    borderWidth: 1,
                  ),
                  textStyle: TextStyle(
                    color: Colors.black,
                    fontSize: 22,
                    fontWeight: FontWeight.bold,
                  ),
                  animationDuration: Duration(milliseconds: 300),
                  backgroundColor: Colors.transparent,
                  enableActiveFill: true,
                  onChanged: _validateCode,
                ),
                SizedBox(height: 20),
                Center(
                  child: Column(
                    children: [
                      _timeRemaining > 0
                          ? Text(
                            "Mã hết hạn sau: ${_formatTime(_timeRemaining)}",
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                              color: Colors.red,
                            ),
                          )
                          : Text(
                            "Mã đã hết hạn, vui lòng yêu cầu mã mới!",
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                              color: Colors.red,
                            ),
                          ),
                    ],
                  ),
                ),
                SizedBox(height: 30),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed:
                        isButtonEnabled && _timeRemaining > 0
                            ? _verifyOtp
                            : null,
                    style: ElevatedButton.styleFrom(
                      minimumSize: Size(double.infinity, 50),
                      backgroundColor:
                          (isButtonEnabled && _timeRemaining > 0)
                              ? Colors.blue
                              : Colors.grey[300],
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(30),
                      ),
                    ),
                    child: Text("Xác nhận", style: TextStyle(fontSize: 18)),
                  ),
                ),
                SizedBox(height: 25),
                SizedBox(
                  width: double.infinity,
                  child: TextButton(
                    onPressed: _timeRemaining == 0 ? _resendOtp : null,
                    child: Text(
                      "Bạn đã nhận mã chưa? Gửi lại",
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: _timeRemaining == 0 ? Colors.blue : Colors.grey,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
