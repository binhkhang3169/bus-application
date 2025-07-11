import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:taixe/ui/login/change_password_page.dart';
import 'package:taixe/ui/login/login_page.dart';
import 'package:taixe/ui/login/profile/profile_image_picker.dart';
import 'package:taixe/ui/login/profile/user_profile_page.dart';

class SettingPage extends StatefulWidget {
  const SettingPage({super.key});

  @override
  State<SettingPage> createState() => _SettingPageState();
}

class _SettingPageState extends State<SettingPage> {
  bool check = false;
  String fullName = "";
  String token = "";
  String image_url = "";
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadUserData();
  }

  Future<void> _loadUserData() async {
    setState(() {
      _isLoading = true;
    });
    SharedPreferences prefs = await SharedPreferences.getInstance();
    setState(() {
      fullName = prefs.getString('fullName') ?? "";
      token = prefs.getString('accessToken') ?? "";
      image_url = prefs.getString('image') ?? "";
      check = token.isNotEmpty;
      _isLoading = false;
    });
  }

  void _logout() async {
    setState(() {
      check = false;
    });
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.clear();
    Navigator.pushReplacement(
      context,
      MaterialPageRoute(builder: (context) => LoginScreen()),
    );
  }

  void _confirmLogout() {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            backgroundColor: Colors.white,
            title: const Text("Xác nhận đăng xuất"),
            content: const Text("Bạn có chắc chắn muốn đăng xuất không?"),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text("Hủy", style: TextStyle(color: Colors.blue)),
              ),
              TextButton(
                onPressed: () {
                  Navigator.pop(context);
                  _logout();
                },
                child: const Text(
                  "Đăng xuất",
                  style: TextStyle(color: Colors.red),
                ),
              ),
            ],
          ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 243, 243, 243),
      body: Stack(
        children: [
          Positioned(
            top: 120,
            left: 0,
            right: 0,
            child: Container(
              height: 40,
              decoration: BoxDecoration(color: Colors.white),
            ),
          ),
          Column(
            children: [
              _buildHeader(),
              Expanded(
                child: SingleChildScrollView(
                  physics: NeverScrollableScrollPhysics(),
                  child: Column(children: _buildMenuItems()),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Stack(
      children: [
        Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(vertical: 40, horizontal: 20),
          decoration: const BoxDecoration(
            color: Colors.blue,
            borderRadius: BorderRadius.only(
              bottomLeft: Radius.circular(20),
              bottomRight: Radius.circular(20),
            ),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              ClipRRect(
                borderRadius: const BorderRadius.only(
                  topRight: Radius.circular(20),
                  bottomRight: Radius.circular(20),
                  bottomLeft: Radius.circular(20),
                ),
                child: ProfileImagePicker(imageUrl: image_url),
              ),
              const SizedBox(width: 10),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    fullName,
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: Colors.white,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ],
    );
  }

  List<Widget> _buildMenuItems() {
    return [
      _buildMenuItem(Icons.person, "Thông tin cá nhân", () {
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => UserProfilePage()),
        ).then((_) => _loadUserData()); // Reload data after returning
      }),
      _buildMenuItem(Icons.notification_important, "Thông báo", () {}),
      _buildMenuItem(Icons.info, "Phiên bản", () {}),
      _buildMenuItem(Icons.password_rounded, "Thay đổi mật khẩu", () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => ChangePasswordScreen(token: token),
          ),
        ).then((_) => _loadUserData()); // Reload data after returning
      }),
      _buildMenuItem(Icons.security, "Chính sách và điều khoản", () {}),
      const SizedBox(height: 10),
      check
          ? Column(
            children: [
              _buildMenuItem(Icons.logout, "Đăng xuất", () {
                _confirmLogout();
              }),
            ],
          )
          : _buildMenuItem(Icons.account_circle_outlined, "Đăng nhập", () {
            Navigator.push(
              context,
              MaterialPageRoute(builder: (context) => LoginScreen()),
            );
          }),
    ];
  }

  Widget _buildMenuItem(IconData icon, String title, VoidCallback onTap) {
    return Container(
      color: Colors.white,
      child: ListTile(
        leading: Icon(icon, color: Colors.black54),
        title: Text(title, style: const TextStyle(fontSize: 14)),
        trailing: const Icon(
          Icons.arrow_forward_ios,
          size: 14,
          color: Colors.black54,
        ),
        onTap: onTap,
      ),
    );
  }
}
