import 'package:flutter/material.dart';

class NotificationPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,

      appBar: AppBar(
        title: Center(
          child: Text(
            "Thông báo",
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 18,
              color: Colors.white,
            ),
          ),
        ),

        backgroundColor: Colors.blueAccent,
      ),
      body: ListView(
        children: [
          ListTile(
            title: Text("Thông báo mới"),
            subtitle: Text("Thay đổi số tổng đài..."),
            trailing: Text("20/03/2025"),
          ),
          ListTile(
            title: Text("Khuyến mãi"),
            subtitle: Text("Giảm giá 50%..."),
            trailing: Text("19/03/2025"),
          ),
        ],
      ),
    );
  }
}
