import 'package:badges/badges.dart' as badges;
import 'package:caoky/ui/screens/history_page.dart';
import 'package:caoky/ui/screens/home_page.dart';
import 'package:caoky/ui/screens/notification_page.dart';
import 'package:caoky/ui/screens/qr_page.dart';
import 'package:caoky/ui/screens/setting_page.dart';
import 'package:curved_navigation_bar/curved_navigation_bar.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MainPage extends StatefulWidget {
  final int initialIndex;

  const MainPage({Key? key, this.initialIndex = 0}) : super(key: key);

  @override
  State<MainPage> createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  late int _selectedIndex;

  int? userId;
  String token = "";

  final List<Widget> pages = [
    HomePage(),
    HistoryPage(),
    NotificationPage(),
    SettingPage(),
  ];

  final List<String> labels = ["Trang chủ", "Lịch sử", "Thông báo", "Cài đặt"];
  final List<IconData> icons = [
    Icons.home,
    Icons.history,
    Icons.notifications,
    Icons.person,
  ];

  @override
  void initState() {
    super.initState();
    _selectedIndex = widget.initialIndex;
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Scaffold(
        extendBody: true,
        backgroundColor: Colors.blueAccent,
        body: IndexedStack(index: _selectedIndex, children: pages),
        bottomNavigationBar: Stack(
          clipBehavior: Clip.none,
          alignment: Alignment.bottomCenter,
          children: [
            Positioned(
              bottom: 0,
              left: 0,
              right: 0,
              child: Container(
                padding: EdgeInsets.symmetric(vertical: 5),
                color: Colors.blueAccent,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: List.generate(
                    labels.length,
                    (index) => SizedBox(
                      width: 60,
                      child: Text(
                        labels[index],
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 12,
                          color:
                              _selectedIndex == index
                                  ? Colors.white
                                  : const Color.fromARGB(255, 245, 231, 231),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ),
            Positioned(
              bottom: 27,
              left: 0,
              right: 0,
              child: CurvedNavigationBar(
                height: 45,
                animationCurve: Curves.easeInOut,
                animationDuration: Duration(milliseconds: 300),
                backgroundColor: Colors.transparent,
                color: Colors.blueAccent,
                items: List.generate(
                  icons.length,
                  (index) => Icon(
                    icons[index],
                    size: 25,
                    color:
                        _selectedIndex == index
                            ? Colors.white
                            : const Color.fromARGB(255, 245, 231, 231),
                  ),
                ),
                onTap: (index) {
                  setState(() {
                    _selectedIndex = index;
                  });
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
