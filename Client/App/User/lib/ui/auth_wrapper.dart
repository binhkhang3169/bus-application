import 'package:caoky/services/auth_service.dart';
import 'package:caoky/ui/login/login_page.dart';
import 'package:caoky/ui/main_page.dart';
import 'package:flutter/material.dart';

class AuthWrapper extends StatefulWidget {
  const AuthWrapper({super.key});

  @override
  State<AuthWrapper> createState() => _AuthWrapperState();
}

class _AuthWrapperState extends State<AuthWrapper> {
  final AuthRepository _authRepository = AuthRepository();

  @override
  void initState() {
    super.initState();
    _decideNextScreen();
  }

  Future<void> _decideNextScreen() async {
    // Use the attemptAutoLogin method from your refactored repository
    final role = await _authRepository.attemptAutoLogin();
    
    // Ensure the widget is still mounted before navigating
    if (!mounted) return;

    Widget destinationPage;
    if (role != null) {
      print("AuthWrapper: Auto login successful. Role: $role. Navigating to MainPage.");
      // You can add logic for different roles here if needed
      destinationPage = MainPage();
    } else {
      print("AuthWrapper: No valid session found. Navigating to LoginScreen.");
      destinationPage = LoginScreen();
    }

    // Use pushAndRemoveUntil to replace the loading screen with the destination
    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (context) => destinationPage),
      (route) => false, // This predicate removes all previous routes
    );
  }

  @override
  Widget build(BuildContext context) {
    // Show a loading indicator while we're checking the auth state
    return const Scaffold(
      backgroundColor: Colors.white,
      body: Center(
        child: CircularProgressIndicator(
          color: Colors.orange,
        ),
      ),
    );
  }
}