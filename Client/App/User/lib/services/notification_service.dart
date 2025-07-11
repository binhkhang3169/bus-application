import 'dart:convert'; // Import để sử dụng jsonEncode

import 'package:caoky/firebase_options.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http; // Import thư viện http

class NotificationService {
  static final FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin =
      FlutterLocalNotificationsPlugin();
  static final FirebaseMessaging firebaseMessaging = FirebaseMessaging.instance;

  // HÀM MỚI: Gửi token lên server
  static Future<void> _registerTokenOnServer({
    required String token,
    required String userId,
  }) async {
    // !!! THAY THẾ BẰNG ĐỊA CHỈ API THỰC TẾ CỦA BẠN
    final url = Uri.parse(
      'http://57.155.76.74/api/v1/usersnoti/$userId/fcm-token',
    );

    try {
      final response = await http.post(
        url,
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          // Thêm các header khác nếu cần, ví dụ: Authorization
        },
        body: jsonEncode(<String, String>{'token': token}),
      );

      if (response.statusCode == 200) {
        print('FCM Token registered successfully for user: $userId');
      } else {
        print(
          'Failed to register FCM token. Status code: ${response.statusCode}',
        );
        print('Response body: ${response.body}');
      }
    } catch (e) {
      print('Error registering FCM token: $e');
    }
  }

  @pragma('vm:entry-point')
  static Future<void> firebaseMessagingBackgroundHandler(
    RemoteMessage message,
  ) async {
    await Firebase.initializeApp(
      options: DefaultFirebaseOptions.currentPlatform,
    );
    await _initializeLocalNotification();
    await _showFlutterNotification(message);
  }

  // SỬA ĐỔI: Hàm này giờ yêu cầu userId
  static Future<void> initializeNotification({required String userId}) async {
    await firebaseMessaging.requestPermission();

    FirebaseMessaging.onMessage.listen((RemoteMessage message) async {
      await _showFlutterNotification(message);
    });

    FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) async {
      print("User tap noti from background: ${message.data}");
    });

    // SỬA ĐỔI: Gọi hàm lấy và đăng ký token với userId
    await _getAndRegisterFcmToken(userId: userId);

    await _initializeLocalNotification();
    await _getInitialNotification();
  }

  // SỬA ĐỔI: Đổi tên và thêm logic gọi API
  static Future<void> _getAndRegisterFcmToken({required String userId}) async {
    String? token = await firebaseMessaging.getToken();
    print("FCM Token: $token");
    if (token != null) {
      // Gọi hàm để gửi token lên server của bạn
      await _registerTokenOnServer(token: token, userId: userId);
    }
  }

  static Future<void> _showFlutterNotification(RemoteMessage message) async {
    // Backend đã gửi data payload, nên chúng ta sẽ ưu tiên đọc từ data
    Map<String, dynamic> data = message.data;
    String title = data['title'] ?? "No title";
    String body = data['body'] ?? "No body";

    // RemoteNotification có thể là null khi app ở foreground trên iOS
    // Hoặc khi chỉ gửi data-only payload
    RemoteNotification? notification = message.notification;

    AndroidNotificationDetails androidDetails = AndroidNotificationDetails(
      'your_channel_id', // Thay đổi channel ID
      'your_channel_name', // Thay đổi channel name
      channelDescription: 'your_channel_description',
      importance: Importance.max,
      priority: Priority.high,
    );
    DarwinNotificationDetails iosDetails = DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );
    NotificationDetails notificationDetails = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );

    // Sử dụng hashcode của message để đảm bảo mỗi thông báo có ID duy nhất
    await flutterLocalNotificationsPlugin.show(
      message.hashCode,
      title,
      body,
      notificationDetails,
      // Gửi payload từ FCM để xử lý khi người dùng nhấn vào
      payload: jsonEncode(message.data),
    );
  }

  static Future<void> _initializeLocalNotification() async {
    const AndroidInitializationSettings androidInitializationSettings =
        AndroidInitializationSettings('@drawable/ic_launcher');

    const DarwinInitializationSettings iosInitializationSettings =
        DarwinInitializationSettings();

    final InitializationSettings initializationSettings =
        InitializationSettings(
          android: androidInitializationSettings,
          iOS: iosInitializationSettings,
        );

    await flutterLocalNotificationsPlugin.initialize(
      initializationSettings,
      onDidReceiveNotificationResponse: (
        NotificationResponse notificationResponse,
      ) {
        // Xử lý khi người dùng nhấn vào thông báo (khi app đang mở)
        print("User tap local noti payload: ${notificationResponse.payload}");
        if (notificationResponse.payload != null) {
          final Map<String, dynamic> data = jsonDecode(
            notificationResponse.payload!,
          );
          print("Decoded data from local noti tap: $data");
          // TODO: Thêm logic điều hướng hoặc xử lý dựa trên `data`
        }
      },
    );
  }

  static Future<void> _getInitialNotification() async {
    RemoteMessage? message = await firebaseMessaging.getInitialMessage();
    if (message != null) {
      // Xử lý khi người dùng nhấn vào thông báo và mở app từ trạng thái bị tắt
      print("Opened app from terminated state via noti: ${message.data}");
      // TODO: Thêm logic điều hướng hoặc xử lý dựa trên `message.data`
    }
  }
}
