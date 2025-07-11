import 'dart:convert';
import 'package:caoky/global/toast.dart';
import 'package:caoky/models/user/signup_request.dart';
import 'package:caoky/services/api_user_service.dart';
import 'package:caoky/ui/login/keys/shipping.dart';
import 'package:caoky/ui/login/otp_page.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:fluttertoast/fluttertoast.dart';
import 'package:intl/intl.dart';
import 'package:flutter/gestures.dart';
import 'package:shared_preferences/shared_preferences.dart';

class RegisterPage extends StatefulWidget {
  @override
  _RegisterPageState createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  bool _isLoading = false;
  bool _isPasswordVisible = false;
  bool _isAgreed = false;
  String address = "";
  String codes = "";

  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _phoneNumberController = TextEditingController();
  final TextEditingController _fullNameController = TextEditingController();
  final TextEditingController _specificAddressController =
      TextEditingController();
  final TextEditingController _birthdayController = TextEditingController();

  final _usernameFocusNode = FocusNode();
  final _passwordFocusNode = FocusNode();
  final _phoneNumberFocusNode = FocusNode();
  final _emailFocusNode = FocusNode();
  final _fullNameFocusNode = FocusNode();
  final _specificAddressFocusNode = FocusNode();
  final _birthdayFocusNode = FocusNode();
  final _provinceFocusNode = FocusNode();
  final _districtFocusNode = FocusNode();
  final _wardFocusNode = FocusNode();

  List<dynamic> provinces = [];
  List<dynamic> districts = [];
  List<dynamic> wards = [];
  String? selectedProvince;
  String? selectedDistrict;
  String? selectedWard;
  String? selectedGender;

  @override
  void initState() {
    super.initState();
    fetchProvinces();
  }

  void getNameFromIds(String idString) {
    List<String> ids = idString.split(',');

    if (ids.length != 3) {
      print("Invalid ID string");
      return;
    }

    String wardId = ids[0].trim();
    String districtId = ids[1].trim();
    String provinceId = ids[2].trim();

    String provinceName = getProvinceName(provinceId);
    String districtName = getDistrictName(districtId);

    String wardName = getWardName(wardId);

    setState(() {
      address =
          "${_specificAddressController.text}, $wardName, $districtName, $provinceName";
    });
  }

  String getProvinceName(String provinceId) {
    var province = provinces.firstWhere(
      (element) => element['ProvinceID'].toString() == provinceId,
      orElse: () => null,
    );
    return province != null ? province['ProvinceName'] : "Không tìm thấy tỉnh";
  }

  String getDistrictName(String districtId) {
    var district = districts.firstWhere(
      (element) => element['DistrictID'].toString() == districtId,
      orElse: () => null,
    );
    return district != null ? district['DistrictName'] : "Không tìm thấy huyện";
  }

  String getWardName(String wardId) {
    var ward = wards.firstWhere(
      (element) => element['WardCode'].toString() == wardId,
      orElse: () => null,
    );
    return ward != null ? ward['WardName'] : "Không tìm thấy xã";
  }

  Future<void> fetchProvinces() async {
    Shipping shipping = Shipping();
    final url = Uri.parse(
      "https://dev-online-gateway.ghn.vn/shiip/public-api/master-data/province",
    );
    try {
      final response = await http.get(url, headers: {"Token": shipping.apiKey});
      if (response.statusCode == 200) {
        final data = jsonDecode(utf8.decode(response.bodyBytes));
        setState(() {
          provinces = data['data'] ?? [];
        });
      } else {
        ToastUtils.show("Lỗi khi tải danh sách tỉnh/thành");
      }
    } catch (e) {
      ToastUtils.show("Không thể kết nối đến server");

      print("Lỗi kết nối: $e");
    }
  }

  Future<void> fetchDistricts(int provinceId) async {
    Shipping shipping = Shipping();
    final url = Uri.parse(
      "https://dev-online-gateway.ghn.vn/shiip/public-api/master-data/district",
    );
    try {
      final response = await http.post(
        url,
        headers: {"Content-Type": "application/json", "Token": shipping.apiKey},
        body: jsonEncode({"province_id": provinceId}),
      );
      if (response.statusCode == 200) {
        final data = jsonDecode(utf8.decode(response.bodyBytes));
        setState(() {
          districts = data['data'] ?? [];
          selectedDistrict = null;
          wards = [];
          selectedWard = null;
        });
      } else {
        ToastUtils.show("Lỗi khi tải danh sách quận/huyện");
      }
    } catch (e) {
      ToastUtils.show("Không thể kết nối đến server");
      print("Lỗi kết nối: $e");
    }
  }

  Future<void> fetchWards(int districtId) async {
    Shipping shipping = Shipping();
    final url = Uri.parse(
      "https://dev-online-gateway.ghn.vn/shiip/public-api/master-data/ward",
    );
    try {
      final response = await http.post(
        url,
        headers: {"Content-Type": "application/json", "Token": shipping.apiKey},
        body: jsonEncode({"district_id": districtId}),
      );
      if (response.statusCode == 200) {
        final data = jsonDecode(utf8.decode(response.bodyBytes));
        setState(() {
          wards = data['data'] ?? [];
          selectedWard = null;
        });
      } else {
        ToastUtils.show("Lỗi khi tải danh sách phường/xã");
      }
    } catch (e) {
      ToastUtils.show("Không thể kết nối đến server");
      print("Lỗi kết nối: $e");
    }
  }

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime(1900),
      lastDate: DateTime.now(),
      builder: (BuildContext context, Widget? child) {
        return Theme(
          data: Theme.of(context).copyWith(
            dialogBackgroundColor: Colors.white,
            colorScheme: ColorScheme.light(
              primary: Colors.blue,
              onPrimary: Colors.white,
              onSurface: Colors.black,
              background: Colors.grey[200],
            ),
          ),
          child: child!,
        );
      },
    );

    if (picked != null) {
      setState(() {
        _birthdayController.text = DateFormat('dd/MM/yyyy').format(picked);
      });
    }
  }

  Future<void> _register() async {
    if (_formKey.currentState!.validate()) {
      if (selectedProvince == null ||
          selectedDistrict == null ||
          selectedWard == null) {
        ToastUtils.show("Vui lòng chọn đầy đủ địa chỉ");
        return;
      }
      if (selectedGender == null) {
        ToastUtils.show("Vui lòng chọn giới tính");
        return;
      }
      if (!_isAgreed) {
        ToastUtils.show("Vui lòng đồng ý với điều khoản");
        return;
      }
      setState(() {
        _isLoading = true;
      });

      setState(() {
        codes = "$selectedWard,$selectedDistrict,$selectedProvince";
        if (codes != null) {
          getNameFromIds(codes);
        }
      });

      final dio = Dio();
      final apiService = ApiUserService(dio);

      try {
        final request = SignupRequest(
          username: _usernameController.text.trim(),
          password: _passwordController.text.trim(),
          phoneNumber: _phoneNumberController.text.trim(),
          fullName: _fullNameController.text.trim(),
          address:
              _specificAddressController.text.trim() +
              ", " +
              selectedWard.toString() +
              ", " +
              selectedDistrict.toString() +
              ", " +
              selectedProvince.toString(),
          gender: selectedGender.toString(),
          otp: "",
        );

        final response = await apiService.register(request);

        if (response.code == 200) {
          final prefs = await SharedPreferences.getInstance();
          await prefs.setString('username', _usernameController.text.trim());
          await prefs.setString('password', _passwordController.text.trim());
          await prefs.setString(
            'phoneNumber',
            _phoneNumberController.text.trim(),
          );
          await prefs.setString('address', address);
          await prefs.setString('fullName', _fullNameController.text.trim());
          await prefs.setString('gender', selectedGender.toString());

          Navigator.pushReplacement(
            context,
            MaterialPageRoute(builder: (context) => VerifyOtpScreen()),
          );

          ToastUtils.show(response.message);
        } else {
          ToastUtils.show(response.message);
        }
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
    } else {
      ToastUtils.show("Vui lòng kiểm tra lại thông tin!");
    }
  }

  //   Future<void> _register() async {
  //   if (selectedProvince == null ||
  //       selectedDistrict == null ||
  //       selectedWard == null) {
  //     Fluttertoast.showToast(msg: "Vui lòng chọn đầy đủ địa chỉ");
  //     return;
  //   }

  //   setState(() {
  //     _isLoading = true;
  //   });

  //   String email = _emailController.text.trim();
  //   String password = _passwordController.text.trim();
  //   String fullName = _fullNameController.text.trim();

  //   if (email.isEmpty) {
  //     Fluttertoast.showToast(msg: "Email không hợp lệ");
  //     setState(() {
  //       _isLoading = false;
  //     });
  //     return;
  //   }

  //   if (password.isEmpty) {
  //     Fluttertoast.showToast(msg: "Vui lòng nhập mật khẩu");
  //     setState(() {
  //       _isLoading = false;
  //     });
  //     return;
  //   }

  //   setState(() {
  //     codes = "$selectedWard,$selectedDistrict,$selectedProvince";
  //     if (codes != null) {
  //       getNameFromIds(codes);
  //     }
  //   });

  // }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Đăng ký tài khoản",
          style: TextStyle(fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
        foregroundColor: Colors.white,
        backgroundColor: Colors.blueAccent,
      ),
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              color: Colors.white,
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      TextFormField(
                        controller: _usernameController,
                        focusNode: _usernameFocusNode,
                        style: TextStyle(color: Colors.black, fontSize: 14),
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
                            return 'Vui lòng nhập tên đăng nhập';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      TextFormField(
                        controller: _passwordController,
                        focusNode: _passwordFocusNode,
                        obscureText: !_isPasswordVisible,
                        style: TextStyle(color: Colors.black, fontSize: 14),
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
                              _isPasswordVisible
                                  ? Icons.visibility
                                  : Icons.visibility_off,
                              color: Colors.grey,
                            ),
                            onPressed: () {
                              setState(() {
                                _isPasswordVisible = !_isPasswordVisible;
                              });
                            },
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30),
                            borderSide: BorderSide.none,
                          ),
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Vui lòng nhập mật khẩu';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      TextFormField(
                        controller: _phoneNumberController,
                        focusNode: _phoneNumberFocusNode,
                        keyboardType: TextInputType.phone,
                        style: TextStyle(color: Colors.black, fontSize: 14),
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Số điện thoại',
                          prefixIcon: Icon(
                            Icons.phone,
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
                            return 'Vui lòng nhập số điện thoại';
                          }
                          if (!RegExp(r'^[0-9]{10}$').hasMatch(value)) {
                            return 'Số điện thoại không hợp lệ';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      TextFormField(
                        controller: _fullNameController,
                        focusNode: _fullNameFocusNode,
                        style: TextStyle(color: Colors.black, fontSize: 14),
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Họ và tên',
                          prefixIcon: Icon(
                            Icons.person_outline,
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
                            return 'Vui lòng nhập họ và tên';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      DropdownButtonFormField(
                        dropdownColor: Colors.white,
                        focusNode: _provinceFocusNode,
                        style: TextStyle(color: Colors.black, fontSize: 14),
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Tỉnh/Thành',
                          prefixIcon: Icon(
                            Icons.location_city,
                            color: Colors.grey,
                            size: 20,
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30),
                            borderSide: BorderSide.none,
                          ),
                        ),
                        value: selectedProvince,
                        items:
                            provinces.map((province) {
                              return DropdownMenuItem(
                                value: province['ProvinceID'].toString(),
                                child: Text(
                                  province['ProvinceName'],
                                  style: TextStyle(
                                    fontSize: 14,
                                    color: Colors.black,
                                  ),
                                ),
                              );
                            }).toList(),
                        onChanged: (value) {
                          setState(() {
                            selectedProvince = value.toString();
                            fetchDistricts(int.parse(value.toString()));
                          });
                        },
                        validator: (value) {
                          if (value == null) {
                            return 'Chọn tỉnh/thành';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      DropdownButtonFormField(
                        dropdownColor: Colors.white,
                        focusNode: _districtFocusNode,
                        style: TextStyle(color: Colors.black, fontSize: 14),
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Quận/Huyện',
                          prefixIcon: Icon(
                            Icons.location_on,
                            color: Colors.grey,
                            size: 20,
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30),
                            borderSide: BorderSide.none,
                          ),
                        ),
                        value: selectedDistrict,
                        items:
                            districts.map((district) {
                              return DropdownMenuItem(
                                value: district['DistrictID'].toString(),
                                child: Text(
                                  district['DistrictName'],
                                  style: TextStyle(
                                    fontSize: 14,
                                    color: Colors.black,
                                  ),
                                ),
                              );
                            }).toList(),
                        onChanged: (value) {
                          setState(() {
                            selectedDistrict = value.toString();
                            fetchWards(int.parse(value.toString()));
                          });
                        },
                        validator: (value) {
                          if (value == null) {
                            return 'Chọn quận/huyện';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      DropdownButtonFormField(
                        dropdownColor: Colors.white,
                        focusNode: _wardFocusNode,
                        style: TextStyle(color: Colors.black, fontSize: 14),
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Phường/Xã',
                          prefixIcon: Icon(
                            Icons.home,
                            color: Colors.grey,
                            size: 20,
                          ),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30),
                            borderSide: BorderSide.none,
                          ),
                        ),
                        value: selectedWard,
                        items:
                            wards.map((ward) {
                              return DropdownMenuItem(
                                value: ward['WardCode'].toString(),
                                child: Text(
                                  ward['WardName'],
                                  style: TextStyle(
                                    fontSize: 14,
                                    color: Colors.black,
                                  ),
                                ),
                              );
                            }).toList(),
                        onChanged: (value) {
                          setState(() {
                            selectedWard = value.toString();
                          });
                        },
                        validator: (value) {
                          if (value == null) {
                            return 'Chọn phường/xã';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      TextFormField(
                        controller: _specificAddressController,
                        focusNode: _specificAddressFocusNode,
                        style: TextStyle(color: Colors.black, fontSize: 14),
                        decoration: InputDecoration(
                          filled: true,
                          fillColor: Colors.grey.shade100,
                          hintText: 'Địa chỉ cụ thể',
                          prefixIcon: Icon(
                            Icons.map,
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
                            return 'Vui lòng nhập địa chỉ cụ thể';
                          }
                          return null;
                        },
                      ),
                      SizedBox(height: 16),
                      Row(
                        children: [
                          Expanded(
                            flex: 4,
                            child: TextFormField(
                              controller: _birthdayController,
                              focusNode: _birthdayFocusNode,
                              readOnly: true,
                              style: TextStyle(
                                color: Colors.black,
                                fontSize: 14,
                              ),
                              decoration: InputDecoration(
                                filled: true,
                                fillColor: Colors.grey.shade100,
                                hintText: 'Ngày sinh',
                                prefixIcon: Icon(
                                  Icons.cake,
                                  color: Colors.grey,
                                  size: 20,
                                ),
                                suffixIcon: Icon(
                                  Icons.calendar_today,
                                  color: Colors.grey,
                                  size: 20,
                                ),
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(30),
                                  borderSide: BorderSide.none,
                                ),
                              ),
                              onTap: () => _selectDate(context),
                              validator: (value) {
                                if (value == null || value.isEmpty) {
                                  return 'Vui lòng chọn ngày sinh';
                                }
                                return null;
                              },
                            ),
                          ),
                          SizedBox(width: 8),
                          Expanded(
                            flex: 3,
                            child: DropdownButtonFormField(
                              dropdownColor: Colors.white,

                              style: TextStyle(
                                color: Colors.black,
                                fontSize: 14,
                              ),
                              decoration: InputDecoration(
                                filled: true,
                                fillColor: Colors.grey.shade100,
                                hintText: 'Giới tính',
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
                              value: selectedGender,
                              items:
                                  ['Nam', 'Nữ', 'Khác'].map((gender) {
                                    return DropdownMenuItem(
                                      value: gender,
                                      child: Text(
                                        gender,
                                        style: TextStyle(
                                          fontSize: 14,
                                          color: Colors.black,
                                        ),
                                      ),
                                    );
                                  }).toList(),
                              onChanged: (value) {
                                setState(() {
                                  selectedGender = value.toString();
                                });
                              },
                              validator: (value) {
                                if (value == null) {
                                  return 'Chọn giới tính';
                                }
                                return null;
                              },
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: 8),
                      Row(
                        children: [
                          Theme(
                            data: Theme.of(context).copyWith(
                              checkboxTheme: CheckboxThemeData(
                                fillColor: MaterialStateProperty.resolveWith((
                                  states,
                                ) {
                                  if (states.contains(MaterialState.selected)) {
                                    return Colors.blue;
                                  }
                                  return Colors.white;
                                }),
                                side: BorderSide(color: Colors.grey, width: 2),
                              ),
                            ),
                            child: Checkbox(
                              value: _isAgreed,
                              onChanged: (value) {
                                setState(() {
                                  _isAgreed = value ?? false;
                                });
                              },
                            ),
                          ),
                          Expanded(
                            child: Text(
                              "Tôi đồng ý với điều khoản sử dụng",
                              style: TextStyle(
                                fontSize: 14,
                                color: Colors.grey.shade600,
                              ),
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: 16),
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: _isLoading ? null : _register,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.blue,
                            foregroundColor: Colors.white,
                            padding: EdgeInsets.symmetric(vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(30),
                            ),
                          ),
                          child:
                              _isLoading
                                  ? CircularProgressIndicator(
                                    color: Colors.white,
                                  )
                                  : Text(
                                    'ĐĂNG KÝ',
                                    style: TextStyle(fontSize: 18),
                                  ),
                        ),
                      ),
                      SizedBox(height: 16),
                      Center(
                        child: RichText(
                          text: TextSpan(
                            text: "Bạn đã có tài khoản? ",
                            style: TextStyle(fontSize: 14, color: Colors.black),
                            children: [
                              TextSpan(
                                text: "Đăng nhập",
                                style: TextStyle(
                                  fontSize: 14,
                                  color: Colors.blue,
                                  fontWeight: FontWeight.bold,
                                ),
                                recognizer:
                                    TapGestureRecognizer()
                                      ..onTap = () {
                                        Navigator.pop(context);
                                      },
                              ),
                            ],
                          ),
                        ),
                      ),
                      SizedBox(height: 16),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    _phoneNumberController.dispose();
    _fullNameController.dispose();
    _specificAddressController.dispose();
    _birthdayController.dispose();
    _usernameFocusNode.dispose();
    _passwordFocusNode.dispose();
    _phoneNumberFocusNode.dispose();
    _emailFocusNode.dispose();
    _fullNameFocusNode.dispose();
    _specificAddressFocusNode.dispose();
    _birthdayFocusNode.dispose();
    _provinceFocusNode.dispose();
    _districtFocusNode.dispose();
    _wardFocusNode.dispose();
    super.dispose();
  }
}

class WaveClipper extends CustomClipper<Path> {
  @override
  Path getClip(Size size) {
    var path = Path();
    path.lineTo(0, size.height - 40);

    var firstControlPoint = Offset(size.width / 4, size.height);
    var firstEndPoint = Offset(size.width / 2, size.height - 40);
    path.quadraticBezierTo(
      firstControlPoint.dx,
      firstControlPoint.dy,
      firstEndPoint.dx,
      firstEndPoint.dy,
    );

    var secondControlPoint = Offset(size.width * 3 / 4, size.height - 80);
    var secondEndPoint = Offset(size.width, size.height - 40);
    path.quadraticBezierTo(
      secondControlPoint.dx,
      secondControlPoint.dy,
      secondEndPoint.dx,
      secondEndPoint.dy,
    );

    path.lineTo(size.width, 0);
    path.close();
    return path;
  }

  @override
  bool shouldReclip(CustomClipper<Path> oldClipper) => false;
}
