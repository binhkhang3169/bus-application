import 'package:flutter/material.dart';
import 'package:taixe/main_page.dart';
import 'package:taixe/ui/login/login_page.dart';
import 'calendar_page.dart';
import 'checkin_page.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Driver App',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(primaryColor: Colors.blue),
      initialRoute: '/',
      routes: {
        '/': (context) => LoginScreen(),
        '/main': (context) => MainPage(),
        '/calendar': (context) => CalendarPage(),
        // '/checkin':
        //     (context) => CheckInPage(
        //       trip: {
        //         "startLocation": "Bến xe Miền Tây",
        //         "endLocation": "Bến xe Lai Vung",
        //         "startTime": "08:00",
        //       },
        //       selectedDate: DateTime(2025, 5, 21),
        //     ),
      },
      onGenerateRoute: (settings) {
        if (settings.name == '/checkin') {
          final args = settings.arguments as Map<String, dynamic>?;
          return MaterialPageRoute(
            builder:
                (context) => CheckInPage(
                  trip:
                      args?['trip'] ??
                      {
                        "startLocation": "Bến xe Miền Tây",
                        "endLocation": "Bến xe Lai Vung",
                        "startTime": "08:00",
                      },
                  // selectedDate: args?['selectedDate'] ?? DateTime(2025, 5, 21),
                ),
          );
        }
        return null;
      },
    );
  }
}
