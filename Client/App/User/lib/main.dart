import 'package:caoky/firebase_options.dart';
import 'package:caoky/services/notification_service.dart';
import 'package:caoky/ui/auth_wrapper.dart'; // Import the new wrapper
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_stripe/flutter_stripe.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Stripe
  Stripe.publishableKey = 'pk_test_51RU9va4EIvWCmk48xBDoCluY603octv5fw7tV4CQOS7GuZwFtPkv9SLobNdlQd6ashpmubrLH62tzvNPh7lbkOqb00KNmFSSgb';
  await Stripe.instance.applySettings();

  // Initialize Firebase
  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );
  
  // Firebase background message handler
  FirebaseMessaging.onBackgroundMessage(
    NotificationService.firebaseMessagingBackgroundHandler,
  );

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'DACNTT TDTU',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(primarySwatch: Colors.orange),
      locale: const Locale('vi', 'VN'),
      supportedLocales: const [Locale('vi', 'VN')],
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        DefaultWidgetsLocalizations.delegate,
      ],
      // The key change: start with AuthWrapper
      home: const AuthWrapper(),
    );
  }
}