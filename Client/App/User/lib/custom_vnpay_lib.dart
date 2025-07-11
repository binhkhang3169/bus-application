// import 'dart:convert';
// import 'package:crypto/crypto.dart';
// import 'package:flutter/foundation.dart';
// import 'package:flutter/material.dart';
// import 'package:webview_flutter/webview_flutter.dart';
// import 'package:intl/intl.dart';
// import 'package:url_launcher/url_launcher.dart';

// //[VNPayHashType] List of Hash Type in VNPAY, default is HMACSHA512
// enum VNPayHashType {
//   SHA256,
//   HMACSHA512,
// }

// //[BankCode] List of valid payment bank in VNPAY, if not provide, it will be manual select, default is null
// enum BankCode { VNPAYQR, VNBANK, INTCARD }

// //[VNPayHashTypeExt] Extension to convert from HashType Enum to valid string of VNPAY
// extension VNPayHashTypeExt on VNPayHashType {
//   String toValueString() {
//     switch (this) {
//       case VNPayHashType.SHA256:
//         return 'SHA256';
//       case VNPayHashType.HMACSHA512:
//         return 'HmacSHA512';
//     }
//   }
// }

// //[VNPAYFlutter] instance class VNPAY Flutter - Custom version
// class CustomVNPAYFlutter {
//   static final CustomVNPAYFlutter _instance = CustomVNPAYFlutter();

//   //[instance] Single Ton Init
//   static CustomVNPAYFlutter get instance => _instance;

//   //[generatePaymentUrl] Generate payment Url with input parameters - Fixed to match Go version
//   String generatePaymentUrl({
//     String url = 'https://sandbox.vnpayment.vn/paymentv2/vpcpay.html',
//     required String version,
//     String command = 'pay',
//     required String tmnCode,
//     String locale = 'vn',
//     String currencyCode = 'VND',
//     required String txnRef,
//     String orderInfo = 'Pay Order',
//     required double amount,
//     required String returnUrl,
//     required String ipAdress,
//     DateTime? createAt,
//     required String vnpayHashKey,
//     VNPayHashType vnPayHashType = VNPayHashType.HMACSHA512,
//     String vnpayOrderType = 'other',
//     BankCode? bankCode,
//     required DateTime vnpayExpireDate,
//   }) {
//     // Create input data map like in Go version
//     final inputData = <String, String>{
//       'vnp_Version': version,
//       'vnp_TmnCode': tmnCode,
//       'vnp_Amount': (amount * 100).toStringAsFixed(0), // Same as Go: int(req.Amount * 100)
//       'vnp_Command': command,
//       'vnp_CreateDate': DateFormat('yyyyMMddHHmmss')
//           .format(createAt ?? DateTime.now())
//           .toString(),
//       'vnp_CurrCode': currencyCode,
//       'vnp_IpAddr': ipAdress,
//       'vnp_Locale': locale,
//       'vnp_OrderInfo': orderInfo,
//       'vnp_OrderType': vnpayOrderType,
//       'vnp_ReturnUrl': returnUrl,
//       'vnp_TxnRef': txnRef,
//       'vnp_ExpireDate':
//           DateFormat('yyyyMMddHHmmss').format(vnpayExpireDate).toString(),
//     };

//     // Add bank code if provided
//     if (bankCode != null) {
//       inputData['vnp_BankCode'] = bankCode.name;
//     }

//     // Sort keys like in Go version
//     final keys = inputData.keys.toList()..sort();

//     // Build query string and hash data like in Go version
//     final queryBuilder = StringBuffer();
//     final hashDataBuilder = StringBuffer();

//     for (int i = 0; i < keys.length; i++) {
//       final key = keys[i];
//       final value = inputData[key]!;
      
//       // URL encode for query string
//       final encodedKey = Uri.encodeQueryComponent(key);
//       final encodedValue = Uri.encodeQueryComponent(value);

//       if (i > 0) {
//         queryBuilder.write('&');
//         hashDataBuilder.write('&');
//       }

//       // Query string uses encoded values
//       queryBuilder.write(encodedKey);
//       queryBuilder.write('=');
//       queryBuilder.write(encodedValue);

//       // Hash data uses encoded values (same as Go version)
//       hashDataBuilder.write(encodedKey);
//       hashDataBuilder.write('=');
//       hashDataBuilder.write(encodedValue);
//     }

//     // Build VNPay URL
//     final vnpURL = '$url?${queryBuilder.toString()}';
//     final hashData = hashDataBuilder.toString();

//     // Generate secure hash using HMAC-SHA512 like in Go version
//     String vnpSecureHash = "";
//     if (vnPayHashType == VNPayHashType.HMACSHA512) {
//       final hmacSha512 = Hmac(sha512, utf8.encode(vnpayHashKey));
//       final digest = hmacSha512.convert(utf8.encode(hashData));
//       vnpSecureHash = digest.toString();
//     } else if (vnPayHashType == VNPayHashType.SHA256) {
//       final bytes = utf8.encode(vnpayHashKey + hashData);
//       vnpSecureHash = sha256.convert(bytes).toString();
//     }

//     // Final payment URL with secure hash
//     final paymentUrl = '$vnpURL&vnp_SecureHash=$vnpSecureHash';
    
//     debugPrint("=====>[PAYMENT URL]: $paymentUrl");
//     debugPrint("=====>[HASH DATA]: $hashData");
//     debugPrint("=====>[SECURE HASH]: $vnpSecureHash");
    
//     return paymentUrl;
//   }

//   /// Show payment webview - Custom implementation using webview_flutter
//   ///
//   /// [onPaymentSuccess], [onPaymentError] callback when payment success, error on app
//   /// [onWebPaymentComplete] callback when payment complete on web
//   Future<void> show({
//     required BuildContext context,
//     required String paymentUrl,
//     String? appBarTitle,
//     TextStyle appBarTitleStyle = const TextStyle(
//       fontSize: 18,
//       fontWeight: FontWeight.w600,
//     ),
//     Function(Map<String, dynamic>)? onPaymentSuccess,
//     Function(Map<String, dynamic>)? onPaymentError,
//     Function()? onWebPaymentComplete,
//   }) async {
//     if (kIsWeb) {
//       await launchUrl(
//         Uri.parse(paymentUrl),
//         webOnlyWindowName: '_self',
//       );
//       if (onWebPaymentComplete != null) {
//         onWebPaymentComplete();
//       }
//     } else {
//       // Use custom WebView implementation
//       final result = await Navigator.of(context).push<Map<String, dynamic>>(
//         MaterialPageRoute(
//           builder: (context) => CustomVNPayWebView(
//             paymentUrl: paymentUrl,
//             appBarTitle: appBarTitle ?? 'VNPay Payment',
//             appBarTitleStyle: appBarTitleStyle,
//           ),
//         ),
//       );

//       if (result != null) {
//         if (result['vnp_ResponseCode'] == '00') {
//           if (onPaymentSuccess != null) {
//             onPaymentSuccess(result);
//           }
//         } else {
//           if (onPaymentError != null) {
//             onPaymentError(result);
//           }
//         }
//       }
//     }
//   }
// }

// // Custom WebView widget for VNPay
// class CustomVNPayWebView extends StatefulWidget {
//   final String paymentUrl;
//   final String appBarTitle;
//   final TextStyle appBarTitleStyle;

//   const CustomVNPayWebView({
//     Key? key,
//     required this.paymentUrl,
//     required this.appBarTitle,
//     required this.appBarTitleStyle,
//   }) : super(key: key);

//   @override
//   State<CustomVNPayWebView> createState() => _CustomVNPayWebViewState();
// }

// class _CustomVNPayWebViewState extends State<CustomVNPayWebView> {
//   late final WebViewController controller;
//   bool isLoading = true;

//   @override
//   void initState() {
//     super.initState();
//     controller = WebViewController()
//       ..setJavaScriptMode(JavaScriptMode.unrestricted)
//       ..setNavigationDelegate(
//         NavigationDelegate(
//           onPageStarted: (String url) {
//             debugPrint('VNPay WebView - Page started loading: $url');
//             if (url.contains('vnp_ResponseCode')) {
//               _handlePaymentResult(url);
//             }
//           },
//           onPageFinished: (String url) {
//             setState(() {
//               isLoading = false;
//             });
//             debugPrint('VNPay WebView - Page finished loading: $url');
//           },
//           onWebResourceError: (WebResourceError error) {
//             debugPrint('VNPay WebView - Error: ${error.description}');
//           },
//         ),
//       )
//       ..loadRequest(Uri.parse(widget.paymentUrl));
//   }

//   void _handlePaymentResult(String url) {
//     try {
//       final uri = Uri.parse(url);
//       final params = <String, dynamic>{};
     
//       uri.queryParameters.forEach((key, value) {
//         params[key] = value;
//       });

//       debugPrint('VNPay Payment Result: $params');
//       Navigator.of(context).pop(params);
//     } catch (e) {
//       debugPrint('Error parsing payment result: $e');
//       Navigator.of(context).pop({'vnp_ResponseCode': '99', 'error': e.toString()});
//     }
//   }

//   @override
//   Widget build(BuildContext context) {
//     return Scaffold(
//       backgroundColor: Colors.white,
//       appBar: AppBar(
//         title: Text(
//           widget.appBarTitle,
//           style: widget.appBarTitleStyle,
//         ),
//         leading: IconButton(
//           icon: const Icon(Icons.close),
//           onPressed: () {
//             Navigator.of(context).pop({'vnp_ResponseCode': '24'}); // User cancelled
//           },
//         ),
//       ),
//       body: Stack(
//         children: [
//           WebViewWidget(controller: controller),
//           if (isLoading)
//             const Center(
//               child: Column(
//                 mainAxisAlignment: MainAxisAlignment.center,
//                 children: [
//                   CircularProgressIndicator(),
//                   SizedBox(height: 16),
//                   Text('Đang tải trang thanh toán...'),
//                 ],
//               ),
//             ),
//         ],
//       ),
//     );
//   }
// }