import 'package:caoky/global/toast.dart';
import 'package:caoky/services/api_user_service.dart';
import 'package:caoky/ui/profile/manager_address.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class UserProfilePage extends StatefulWidget {
  @override
  _UserProfilePageState createState() => _UserProfilePageState();
}

class _UserProfilePageState extends State<UserProfilePage> {
  late TextEditingController usernameController;
  late TextEditingController phoneNumberController;
  late TextEditingController fullNameController;
  late TextEditingController addressController;
  String selectedGender = 'Nam';
  bool isLoading = true;
  String token = "";
  late ApiUserService _apiService;

  @override
  void initState() {
    super.initState();
    _apiService = ApiUserService(Dio());
    usernameController = TextEditingController();
    phoneNumberController = TextEditingController();
    fullNameController = TextEditingController();
    addressController = TextEditingController();
    loadData();
  }

  Future<void> loadData() async {
    setState(() {
      isLoading = true;
    });

    final prefs = await SharedPreferences.getInstance();

    setState(() {
      token = prefs.getString('accessToken') ?? '';
      usernameController.text = prefs.getString('username') ?? '';
      phoneNumberController.text = prefs.getString('phoneNumber') ?? '';
      fullNameController.text = prefs.getString('fullName') ?? '';
      addressController.text = prefs.getString('address') ?? '';
      selectedGender = prefs.getString('gender') ?? 'Nam';
      isLoading = false;
    });
  }

  Future<void> saveData() async {
    setState(() {
      isLoading = true;
    });
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('phoneNumber', phoneNumberController.text);
    await prefs.setString('fullName', fullNameController.text);
    await prefs.setString('address', addressController.text);
    await prefs.setString('gender', selectedGender);
    Map<String, dynamic> data = {
      'phoneNumber': phoneNumberController.text,
      'fullName': fullNameController.text,
      'address': addressController.text,
      'gender': selectedGender,
    };
    try {
      final response = await _apiService.changeInfo("Bearer $token", data);
      if (response.code == 200) {
        ToastUtils.show("Đã lưu thay đổi");
      } else {
        ToastUtils.show("Lỗi khi lưu thay đổi");
      }
      setState(() {
        isLoading = false;
      });
    } catch (e) {
      print('Lỗi khi lưu dữ liệu: $e');
      ToastUtils.show("Lỗi khi lưu thay đổi");
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void dispose() {
    usernameController.dispose();
    phoneNumberController.dispose();
    fullNameController.dispose();
    addressController.dispose();
    super.dispose();
  }

  Widget buildTextField(
    String label,
    TextEditingController controller,
    bool enabled,
  ) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: TextFormField(
        controller: controller,
        enabled: enabled,
        style: TextStyle(color: Colors.black),
        decoration: InputDecoration(
          filled: true,
          fillColor: Colors.white,
          labelText: label,
          labelStyle: TextStyle(color: Colors.blue.shade700),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: BorderSide(color: Colors.blue, width: 2),
          ),
          border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
        ),
      ),
    );
  }

  Widget buildGenderDropdown() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: DropdownButtonFormField<String>(
        dropdownColor: Colors.white,
        value: selectedGender,
        onChanged: (String? newValue) {
          setState(() {
            selectedGender = newValue!;
          });
        },
        items:
            <String>['Nam', 'Nữ', 'Khác'].map<DropdownMenuItem<String>>((
              String value,
            ) {
              return DropdownMenuItem<String>(
                value: value,
                child: Text(value, style: TextStyle(color: Colors.black)),
              );
            }).toList(),

        decoration: InputDecoration(
          labelText: "Giới tính",

          // Viền các trạng thái
          border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
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
          floatingLabelStyle: TextStyle(color: Colors.blue), // khi focus
          hintStyle: TextStyle(color: Colors.grey.shade400), // màu gợi ý
          suffixStyle: TextStyle(color: Colors.blue), // màu "VND"
        ),
      ),
    );
  }

  Future<void> _navigateToManagerAddressPage() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => ManagerAddressPage()),
    );
    if (result != null) {
      setState(() {
        addressController.text = result;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: Text('Thông tin cá nhân'),
        backgroundColor: Colors.blue,
        foregroundColor: Colors.white,
        centerTitle: true,
      ),
      body: Stack(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: ListView(
              children: [
                buildTextField('Username', usernameController, false),
                buildTextField('Số điện thoại', phoneNumberController, true),
                buildTextField('Họ tên', fullNameController, true),
                GestureDetector(
                  onTap: _navigateToManagerAddressPage,
                  child: AbsorbPointer(
                    child: buildTextField('Địa chỉ', addressController, false),
                  ),
                ),
                buildGenderDropdown(),
                SizedBox(height: 20),
                ElevatedButton(
                  onPressed: saveData,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.blue,
                    padding: EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: Text(
                    'Lưu thay đổi',
                    style: TextStyle(fontSize: 16, color: Colors.white),
                  ),
                ),
              ],
            ),
          ),
          if (isLoading)
            Container(
              color: Colors.white.withOpacity(0.7),
              child: Center(
                child: CircularProgressIndicator(color: Colors.blue),
              ),
            ),
        ],
      ),
    );
  }
}
