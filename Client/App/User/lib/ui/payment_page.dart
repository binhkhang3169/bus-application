// lib/ui/payment_page.dart
import 'dart:async'; // For StreamSubscription
import 'dart:convert'; // For jsonEncode, jsonDecode
import 'dart:developer'; // For log
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_stripe/flutter_stripe.dart' hide Card;
import 'package:get/get_connect/http/src/utils/utils.dart';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart'; // For NumberFormat
import 'package:caoky/services/auth_service.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart'; // Required for generating UUIDs if needed
import 'package:web_socket_channel/io.dart'; // For WebSocket communication

// --- ENUM FOR MANAGING THE BOOKING FLOW STATE ---
enum BookingFlowStatus {
  idle, // Form is visible and ready for user input
  initiating, // "Initiate" button pressed, sending initial request to backend
  waitingForServer, // Waiting for WebSocket response after successful initiation
  paymentSelection, // WebSocket success, showing payment options (Stripe/Bank)
  processingPayment, // A payment method was selected and is being processed
  error, // An error occurred at any stage of the flow
}

class DashedLinePainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint =
        Paint()
          ..color = Colors.grey
          ..strokeWidth = 2
          ..style = PaintingStyle.stroke;
    const dashHeight = 5;
    const dashSpace = 3;
    double startY = 0;
    double centerX = size.width / 2;
    while (startY < size.height) {
      canvas.drawLine(
        Offset(centerX, startY),
        Offset(centerX, startY + dashHeight),
        paint,
      );
      startY += dashHeight + dashSpace;
    }
  }

  @override
  bool shouldRepaint(CustomPainter oldDelegate) => false;
}

class PaymentScreen extends StatefulWidget {
  // Departure Trip Info
  final String selectedRoute;
  final int rawTotalPrice;
  final String seatType;
  final String time;
  final String initialPickupLocation;
  final String initialDropoffLocation;
  final String initialPickupLocationEnd;
  final String initialDropoffLocationEnd;
  final List<String> selectedSeatsNames;
  final List<int> selectedSeatIds;
  final int tripId;
  final String fullRoute;
  final String estimatedDistance;

  // Return Trip Info (Nullable)
  final int? returnTripId;
  final List<String>? selectedReturnSeatsNames;
  final List<int>? selectedReturnSeatIds;
  final dynamic returnTripInfo;

  PaymentScreen({
    super.key,
    required this.selectedRoute,
    required this.rawTotalPrice,
    required this.seatType,
    required this.time,
    required this.initialPickupLocation,
    required this.initialDropoffLocation,
    required this.selectedSeatsNames,
    required this.selectedSeatIds,
    required this.tripId,
    required this.fullRoute,
    required this.estimatedDistance,
    this.returnTripId,
    this.selectedReturnSeatsNames,
    this.selectedReturnSeatIds,
    this.returnTripInfo,
    this.initialPickupLocationEnd = "",
    this.initialDropoffLocationEnd = "",
  });

  @override
  _PaymentScreenState createState() => _PaymentScreenState();
}

class _PaymentScreenState extends State<PaymentScreen> {
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final TextEditingController _emailController = TextEditingController();

  // --- NEW STATE MANAGEMENT FOR ASYNC BOOKING FLOW ---
  BookingFlowStatus _bookingStatus = BookingFlowStatus.idle;
  String? _bookingId; // Stores the ID received from initiate-booking endpoint
  IOWebSocketChannel? _webSocketChannel; // WebSocket channel instance
  StreamSubscription? _webSocketListener;
  String _errorMessage = ""; // To display error messages in the UI

  // User and location data
  String fullName = "";
  String phoneNumber = "";
  String email = "";
  String transferCode = "";
  String selectedPickupType = "Bến xe/VP";
  late String selectedPickupLocation;
  bool showPickupLocationField = true;
  String selectedDropoffType = "Bến xe/VP";
  late String selectedDropoffLocation;
  bool showDropoffLocationField = true;
  List<String> _parsedPickupLocations = [];
  List<String> _parsedDropoffLocations = [];
  String selectedReturnPickupType = "Bến xe/VP";
  String? selectedReturnPickupLocation;
  bool showReturnPickupLocationField = true;
  String selectedReturnDropoffType = "Bến xe/VP";
  String? selectedReturnDropoffLocation;
  bool showReturnDropoffLocationField = true;
  List<String> _parsedReturnPickupLocations = [];
  List<String> _parsedReturnDropoffLocations = [];
  bool _isLoading = true;

  // Payment related state
  final AuthRepository _authRepository = AuthRepository();
  final NumberFormat currencyFormatter = NumberFormat.currency(locale: 'vi_VN', symbol: '₫');
  String? _activeClientSecret;
  String? _activePaymentIntentId;
  String? _retrievedTicketId; // The final ticket ID from WebSocket

  final String _apiBaseUrl = "http://57.155.76.74/api/v1";
  final String _wsBaseUrl = "ws://57.155.76.74/api/v1";


  @override
  void initState() {
    super.initState();
    _initializeScreenData();
  }

  @override
  void dispose() {
    // Clean up controllers and WebSocket connection
    _nameController.dispose();
    _phoneController.dispose();
    _emailController.dispose();
    _webSocketListener?.cancel();
    _webSocketChannel?.sink.close();
    super.dispose();
  }

  // --- CORE BOOKING FLOW LOGIC ---

  /// Starts the entire booking process.
  /// 1. Validates user input.
  /// 2. Calls the backend to initiate the booking.
  Future<void> _startBookingFlow() async {
    // --- VALIDATION ---
    final currentPickupInfo = _extractStationInfo(selectedPickupLocation);
    final currentDropoffInfo = _extractStationInfo(selectedDropoffLocation);

    if (widget.selectedSeatsNames.isEmpty ||
        currentPickupInfo['id']!.isEmpty ||
        currentDropoffInfo['id']!.isEmpty ||
        _nameController.text.trim().isEmpty ||
        _phoneController.text.trim().isEmpty ||
        _emailController.text.trim().isEmpty) {
      _showInfoDialog("Vui lòng điền đầy đủ tất cả thông tin cần thiết.");
      return;
    }

    if (widget.returnTripId != null) {
      if (selectedReturnPickupLocation == null || selectedReturnDropoffLocation == null) {
        _showInfoDialog("Vui lòng chọn điểm đón và trả cho chuyến về.");
        return;
      }
      final returnPickupInfo = _extractStationInfo(selectedReturnPickupLocation!);
      final returnDropoffInfo = _extractStationInfo(selectedReturnDropoffLocation!);
      if (returnPickupInfo['id']!.isEmpty || returnDropoffInfo['id']!.isEmpty) {
        _showInfoDialog("Vui lòng chọn điểm đón và trả hợp lệ cho chuyến về.");
        return;
      }
    }

    setState(() {
      _bookingStatus = BookingFlowStatus.initiating;
    });

    await _initiateBookingOnBackend();
  }

  /// Step 1: Send booking details to the `/initiate-booking` endpoint.
  Future<void> _initiateBookingOnBackend() async {
    final currentPickupInfo = _extractStationInfo(selectedPickupLocation);
    final currentDropoffInfo = _extractStationInfo(selectedDropoffLocation);

    final Map<String, dynamic> bookingPayload = {
      "ticket_type": widget.returnTripId != null ? 1 : 0,
      "price": widget.rawTotalPrice.toDouble(),
      "status": 1,
      "payment_status": 0,
      "booking_channel": 0, // 0 for Mobile App
      "policy_id": 1,
      "phone": _phoneController.text.trim(),
      "email": _emailController.text.trim(),
      "name": _nameController.text.trim(),
      "booked_by": "customer",
      "trip_id_begin": widget.tripId.toString(),
      "seat_id_begin": widget.selectedSeatIds,
      "pickup_location_begin": int.parse(currentPickupInfo['id']!),
      "dropoff_location_begin": int.parse(currentDropoffInfo['id']!),
    };

    if (widget.returnTripId != null) {
      final returnPickupInfo = _extractStationInfo(selectedReturnPickupLocation!);
      final returnDropoffInfo = _extractStationInfo(selectedReturnDropoffLocation!);
      bookingPayload.addAll({
        "trip_id_end": widget.returnTripId.toString(),
        "seat_id_end": widget.selectedReturnSeatIds,
        "pickup_location_end": int.parse(returnPickupInfo['id']!),
        "dropoff_location_end": int.parse(returnDropoffInfo['id']!),
      });
    }

    try {
      final String jsonPayload = jsonEncode(bookingPayload);
      log("--- Sending Payload to Backend ($_apiBaseUrl/initiate-booking) ---");
      log(jsonPayload);
      final token = await _getAuthToken();
      final response = await http.post(
        Uri.parse("$_apiBaseUrl/initiate-booking"),
        headers: {
          "Content-Type": "application/json",
          if (token != null) "Authorization": "Bearer $token",
        },
        body: jsonPayload,
      );

      log("Backend Initiate Booking API Response Status: ${response.statusCode}");
      log("Backend Initiate Booking API Response Body: ${response.body}");

      if (response.statusCode == 202) { // 202 Accepted is the expected response
        final responseData = jsonDecode(response.body) as Map<String, dynamic>;
        final bookingId = responseData['data']?['bookingId']?.toString();
        if (bookingId != null) {
          log("Booking initiated successfully. Booking ID: $bookingId");
          setState(() {
            _bookingId = bookingId;
            _bookingStatus = BookingFlowStatus.waitingForServer;
          });
          _connectAndListenToWebSocket(bookingId);
        } else {
          throw Exception("Server did not return a valid bookingId.");
        }
      } else {
        final errorBody = jsonDecode(response.body);
        final errorMessage = errorBody['detail'] ?? "Lỗi không xác định từ máy chủ.";
        throw Exception("Lỗi khởi tạo đặt vé: ${response.reasonPhrase} ($errorMessage)");
      }
    } catch (e, s) {
      log("Exception calling initiate booking API: $e\n$s");
      setState(() {
        _bookingStatus = BookingFlowStatus.error;
        _errorMessage = "Lỗi kết nối khi khởi tạo đặt vé: ${e.toString()}";
      });
    }
  }

  /// Step 2: Connect to WebSocket to track the booking status.
  void _connectAndListenToWebSocket(String bookingId) {
    try {
      final wsUrl = Uri.parse('$_wsBaseUrl/ws/track/$bookingId');
      _webSocketChannel = IOWebSocketChannel.connect(wsUrl);
      log("Connecting to WebSocket: $wsUrl");

      _webSocketListener = _webSocketChannel!.stream.listen(
        (message) {
          log("WebSocket message received: $message");
          final decodedMessage = jsonDecode(message) as Map<String, dynamic>;

          if (decodedMessage['type'] == 'result') {
            final payload = decodedMessage['payload'] as Map<String, dynamic>?;
            final ticketId = payload?['ticket_id']?.toString();

            if (ticketId != null) {
              _retrievedTicketId = ticketId;
              log("WebSocket SUCCESS: Ticket created with ID: $ticketId");
              _webSocketListener?.cancel(); // Stop listening
              _webSocketChannel?.sink.close();
              
              // Proceed to payment
              if(mounted){
                setState(() {
                  _bookingStatus = BookingFlowStatus.paymentSelection;
                });
                _showPaymentMethodSelection();
              }

            } else {
              throw Exception("WebSocket result payload is missing 'ticket_id'.");
            }
          } else if (decodedMessage['type'] == 'error') {
            final payload = decodedMessage['payload'] as Map<String, dynamic>?;
            final errorMessage = payload?['message']?.toString() ?? "Lỗi không xác định từ WebSocket.";
            throw Exception(errorMessage);
          }
        },
        onError: (error) {
          log("WebSocket Error: $error");
          if(mounted){
            setState(() {
              _bookingStatus = BookingFlowStatus.error;
              _errorMessage = "Lỗi kết nối WebSocket: ${error.toString()}";
            });
          }
           _webSocketListener?.cancel();
          _webSocketChannel?.sink.close();
        },
        onDone: () {
          log("WebSocket connection closed.");
           if (_bookingStatus == BookingFlowStatus.waitingForServer) {
              if(mounted){
                 setState(() {
                   _bookingStatus = BookingFlowStatus.error;
                   _errorMessage = "Kết nối với máy chủ đã đóng trước khi có kết quả.";
                 });
              }
           }
        },
        cancelOnError: true,
      );
    } catch (e) {
       if(mounted){
         setState(() {
           _bookingStatus = BookingFlowStatus.error;
           _errorMessage = "Không thể kết nối tới máy chủ: ${e.toString()}";
         });
       }
    }
  }

  /// Resets the booking flow to the initial state.
  void _resetBookingFlow() {
    _webSocketListener?.cancel();
    _webSocketChannel?.sink.close();
    setState(() {
      _bookingStatus = BookingFlowStatus.idle;
      _bookingId = null;
      _retrievedTicketId = null;
      _errorMessage = "";
    });
  }

  // --- PAYMENT LOGIC (Largely unchanged, but now triggered by WebSocket) ---

  void _showPaymentMethodSelection() {
    // This is now called after WebSocket success
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (ctx) => AlertDialog(
        backgroundColor: Colors.white,
        title: Text("Thanh toán"),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: Icon(Icons.credit_card, color: Colors.indigo),
              title: Text("Thanh toán qua Stripe"),
              onTap: () {
                Navigator.of(ctx).pop();
                _startStripePaymentFlow();
              },
            ),
            ListTile(
              leading: Icon(Icons.account_balance, color: Colors.green),
              title: Text("Thanh toán qua Ngân hàng"),
              onTap: () {
                Navigator.of(ctx).pop();
                _startBankPaymentFlow();
              },
            ),
          ],
        ),
        actions: [
          TextButton(
            child: Text("Hủy bỏ", style: TextStyle(color: Colors.red)),
            onPressed: () {
              log("Payment selection cancelled. Ticket ID $_retrievedTicketId may need to be voided.");
              Navigator.of(ctx).pop();
              _resetBookingFlow(); // Allow user to restart
            },
          ),
        ],
      ),
    );
  }

  Future<void> _startStripePaymentFlow() async {
    setState(() {
      _bookingStatus = BookingFlowStatus.processingPayment;
    });
    final paymentIntentData = await _createStripePaymentIntent();
    if (mounted && paymentIntentData != null) {
      _activeClientSecret = paymentIntentData['client_secret'];
      _activePaymentIntentId = paymentIntentData['payment_intent_id'];
      log("Received from backend - Client Secret: ${_activeClientSecret != null}, Payment Intent ID: $_activePaymentIntentId");
      if (_activeClientSecret != null && _activePaymentIntentId != null) {
        await _initializeAndPresentStripeSheet();
      } else {
        _showErrorDialog("Không nhận được đủ thông tin thanh toán từ máy chủ.");
      }
    }
    if (mounted && _bookingStatus != BookingFlowStatus.idle) { // Avoid resetting if payment was successful
       setState(() {
         // Reset to payment selection if payment fails/is cancelled
         _bookingStatus = BookingFlowStatus.paymentSelection; 
       });
    }
  }

  Future<void> _startBankPaymentFlow() async {
    setState(() {
      _bookingStatus = BookingFlowStatus.processingPayment;
    });

    final bankRequestData = await _createBankPaymentRequest();

    if (mounted && bankRequestData != null) {
      final userConfirmedPayment = await _showBankTransferDetailsDialog(bankRequestData);

      if (userConfirmedPayment == true) {
        final String? invoiceId = bankRequestData['invoice_id']?.toString();
        if (invoiceId == null) {
          _showErrorDialog("Không tìm thấy ID hoá đơn để xác nhận thanh toán.");
        } else {
          final confirmationSuccess = await _confirmBankPaymentOnBackend(invoiceId);
          if (mounted && confirmationSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(SnackBar(
              content: Text("Đã ghi nhận yêu cầu thanh toán. Vé sẽ được xác nhận sớm. Mã vé: ${_retrievedTicketId ?? 'N/A'}"),
              duration: Duration(seconds: 5),
            ));
            Navigator.of(context).popUntil((route) => route.isFirst);
          }
        }
      } else {
        log("Bank payment was cancelled by the user.");
        _showInfoDialog("Bạn đã hủy thanh toán qua ngân hàng.");
      }
    }

    if (mounted) {
       setState(() {
          // Reset to payment selection if payment fails/is cancelled
         _bookingStatus = BookingFlowStatus.paymentSelection;
       });
    }
  }

  // NOTE: The rest of the functions (_createStripePaymentIntent, _confirmPaymentOnBackend, etc.)
  // are assumed to be correct and are kept as they were, since the refactor is about *how* we get
  // to the payment step, not the payment logic itself. They will now correctly use the `_retrievedTicketId`
  // that was set by the WebSocket listener.

  // --- UI BUILDER METHODS ---

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        appBar: AppBar(title: Text("Đang tải...")),
        body: Center(child: CircularProgressIndicator()),
      );
    }

    // Main UI switcher based on the booking status
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () {
            if (_bookingStatus == BookingFlowStatus.idle || _bookingStatus == BookingFlowStatus.paymentSelection) {
              Navigator.pop(context);
            } else {
              // Offer to cancel the ongoing process
              _showCancelConfirmationDialog();
            }
          },
        ),
        title: Text(
          widget.selectedRoute,
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
          overflow: TextOverflow.ellipsis,
        ),
        centerTitle: true,
        backgroundColor: Colors.blueAccent,
        elevation: 3,
      ),
      body: AnimatedSwitcher(
        duration: const Duration(milliseconds: 300),
        child: _buildContentForStatus(),
      ),
      bottomNavigationBar: _bookingStatus == BookingFlowStatus.idle || _bookingStatus == BookingFlowStatus.initiating
          ? _buildPaymentBottomBar()
          : null, // Only show bottom bar when in the form view
    );
  }

  /// Builds the appropriate widget for the current booking status.
  Widget _buildContentForStatus() {
    switch (_bookingStatus) {
      case BookingFlowStatus.waitingForServer:
        return _buildWaitingScreen();
      case BookingFlowStatus.error:
        return _buildErrorScreen();
      case BookingFlowStatus.idle:
      case BookingFlowStatus.initiating:
      case BookingFlowStatus.paymentSelection:
      case BookingFlowStatus.processingPayment:
      default:
        return _buildBookingForm();
    }
  }

  /// The main form for user input.
  Widget _buildBookingForm() {
    return SingleChildScrollView(
      key: ValueKey('BookingForm'),
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // All the original form cards
            _buildCustomerInfoCard(),
            _buildTripInfoWidget(),
            _buildSeatInfoWidget(),
            _buildPickupDropoffWidget(),
            if (widget.returnTripId != null) ...[
              _buildSeatInfoWidget(isReturn: true),
              _buildPickupDropoffWidget(isReturn: true),
            ],
            SizedBox(height: 20),
          ],
        ),
      ),
    );
  }

  /// The bottom bar with the main action button.
  Widget _buildPaymentBottomBar() {
      final currentPickupInfo = _extractStationInfo(selectedPickupLocation);
      final currentDropoffInfo = _extractStationInfo(selectedDropoffLocation);

      bool canProceed = 
        widget.selectedSeatsNames.isNotEmpty &&
        currentPickupInfo['id']!.isNotEmpty &&
        currentDropoffInfo['id']!.isNotEmpty &&
        _nameController.text.trim().isNotEmpty &&
        _phoneController.text.trim().isNotEmpty &&
        _emailController.text.trim().isNotEmpty &&
        _bookingStatus != BookingFlowStatus.initiating; // Disable while initiating


    return Container(
      key: ValueKey('BottomBar'),
      padding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.1), blurRadius: 8, offset: Offset(0, -2))],
      ),
      child: SizedBox(
        width: double.infinity,
        height: 50,
        child: ElevatedButton(
          style: ElevatedButton.styleFrom(
            backgroundColor: canProceed ? Colors.orange : Colors.grey,
            padding: EdgeInsets.symmetric(vertical: 12),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
            elevation: 2,
          ),
          onPressed: canProceed ? _startBookingFlow : null,
          child: _bookingStatus == BookingFlowStatus.initiating
              ? CircularProgressIndicator(valueColor: AlwaysStoppedAnimation<Color>(Colors.white))
              : Text(
                  "Tiếp tục", // Simplified button text
                  style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
                ),
        ),
      ),
    );
  }

  /// Screen shown while waiting for the WebSocket response.
  Widget _buildWaitingScreen() {
    return Center(
      key: ValueKey('WaitingScreen'),
      child: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(),
            SizedBox(height: 24),
            Text(
              "Đang xử lý yêu cầu đặt vé...",
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
            ),
            SizedBox(height: 8),
            Text(
              "Vui lòng không đóng ứng dụng.",
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 14, color: Colors.grey[600]),
            ),
            if (_bookingId != null) ...[
              SizedBox(height: 16),
              Text(
                "Mã yêu cầu của bạn:",
                style: TextStyle(fontSize: 12, color: Colors.grey[700]),
              ),
              Text(
                _bookingId!,
                style: TextStyle(fontSize: 14, fontWeight: FontWeight.bold, fontFamily: 'monospace'),
              ),
            ]
          ],
        ),
      ),
    );
  }

  /// Screen shown when an error occurs.
  Widget _buildErrorScreen() {
    return Center(
      key: ValueKey('ErrorScreen'),
      child: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, color: Colors.red, size: 60),
            SizedBox(height: 20),
            Text(
              "Đã xảy ra lỗi",
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
            SizedBox(height: 12),
            Text(
              _errorMessage,
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 15, color: Colors.grey[700]),
            ),
            SizedBox(height: 24),
            ElevatedButton(
              onPressed: _resetBookingFlow,
              style: ElevatedButton.styleFrom(backgroundColor: Colors.blueAccent),
              child: Text("Thử lại", style: TextStyle(color: Colors.white)),
            ),
          ],
        ),
      ),
    );
  }

  void _showCancelConfirmationDialog() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text("Hủy bỏ quy trình?"),
        content: Text("Bạn có chắc muốn hủy bỏ quy trình đặt vé hiện tại không?"),
        actions: [
          TextButton(
            child: Text("Không"),
            onPressed: () => Navigator.of(ctx).pop(),
          ),
          TextButton(
            child: Text("Có, Hủy bỏ", style: TextStyle(color: Colors.red)),
            onPressed: () {
              Navigator.of(ctx).pop(); // Close dialog
              Navigator.of(context).pop(); // Go back to previous screen
              _resetBookingFlow();
            },
          ),
        ],
      ),
    );
  }

  // --- HELPER & EXISTING UNCHANGED FUNCTIONS ---
  // All other functions like _loadUserData, _fetchTripFullRoute, _parseFullRoute,
  // _showLocationPicker, _build...Widget methods are kept as they were.
  // The payment methods (_createStripePaymentIntent, etc.) are also kept.
  // Only the trigger points and state management around them have changed.

  Future<String?> _getAuthToken() async {
    try {
      final token = _authRepository.getValidAccessToken();
      if (token != null) {
        log("Retrieved Auth Token: $token");
        return token;
      } else {
        log("Auth Token not found in SharedPreferences. Using placeholder.");
        return "YOUR_FALLBACK_OR_TEST_BEARER_TOKEN";
      }
    } catch (e) {
      log("Error fetching token from SharedPreferences: $e. Using placeholder.");
      return "YOUR_FALLBACK_OR_TEST_BEARER_TOKEN";
    }
  }

  Future<Map<String, dynamic>?> _createStripePaymentIntent() async {
    if (_retrievedTicketId == null) {
      _showErrorDialog("Không có mã vé để tạo thanh toán.");
      return null;
    }
    final Map<String, dynamic> payload = {
      "amount": widget.rawTotalPrice,
      "currency": "vnd",
      "customer_id": _emailController.text.trim(),
      "ticket_id": _retrievedTicketId,
    };
    try {
      final String jsonPayload = jsonEncode(payload);
      log("--- Creating Payment Intent ($_apiBaseUrl/stripe/create-payment-intent) ---");
      log(jsonPayload);
      final token = await _getAuthToken();
      final response = await http.post(
        Uri.parse("$_apiBaseUrl/stripe/create-payment-intent"),
        headers: {
          "Content-Type": "application/json",
          if (token != null) "Authorization": "Bearer $token",
        },
        body: jsonPayload,
      );
      log("Create Payment Intent API Response Status: ${response.statusCode}");
      log("Create Payment Intent API Response Body: ${response.body}");
      if (response.statusCode == 200) {
        final responseData = jsonDecode(response.body) as Map<String, dynamic>;
        if (responseData['data']?['client_secret'] != null && responseData['data']?['payment_intent_id'] != null) {
          return responseData['data'];
        } else {
          _showErrorDialog("Dữ liệu thanh toán không hợp lệ từ máy chủ (thiếu client_secret/payment_intent_id).");
          return null;
        }
      } else {
        _showErrorDialog("Lỗi tạo phiên thanh toán: ${response.reasonPhrase} (${response.statusCode}).");
        return null;
      }
    } catch (e, s) {
      log("Exception calling create payment intent API: $e\n$s");
      _showErrorDialog("Lỗi kết nối khi tạo phiên thanh toán: ${e.toString()}");
      return null;
    }
  }

  Future<void> _initializeAndPresentStripeSheet() async {
    if (_activeClientSecret == null) {
      log("Error: Client secret is null before initializing payment sheet.");
      _showErrorDialog("Không thể khởi tạo thanh toán (thiếu client secret). Vui lòng thử lại.");
      return;
    }
    try {
      await Stripe.instance.initPaymentSheet(
        paymentSheetParameters: SetupPaymentSheetParameters(
          paymentIntentClientSecret: _activeClientSecret!,
          merchantDisplayName: "CAOKY",
          style: ThemeMode.light,
          allowsDelayedPaymentMethods: true,
        ),
      );
      await _presentStripePaymentSheet();
    } on StripeException catch (e) {
      log("StripeException during initPaymentSheet: ${e.error.code} - ${e.error.message}");
      _showErrorDialog("Lỗi Stripe khi khởi tạo: ${e.error.localizedMessage ?? e.error.message}");
    } catch (e, s) {
      log("Error initializing Stripe payment sheet: $e\n$s");
      _showErrorDialog("Lỗi khởi tạo thanh toán: ${e.toString()}");
    }
  }

  Future<void> _presentStripePaymentSheet() async {
    try {
      await Stripe.instance.presentPaymentSheet();
      final paymentIntent = await Stripe.instance.retrievePaymentIntent(_activeClientSecret!);
      log("Payment Intent status after sheet dismissal: ${paymentIntent.status}");
      if (paymentIntent.status == PaymentIntentsStatus.Succeeded) {
        log('Payment successful via Stripe SDK! PI ID: ${paymentIntent.id}');
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text("Thanh toán qua Stripe thành công! Đang xác nhận với máy chủ...")));
        await _confirmPaymentOnBackend(paymentIntent.id);
      } else if (paymentIntent.status == PaymentIntentsStatus.Canceled) {
        log('Payment canceled by user. PI ID: ${paymentIntent.id}');
        _showInfoDialog("Thanh toán đã bị hủy.");
      } else {
        log('Payment failed or requires action. Status: ${paymentIntent.status}, PI ID: ${paymentIntent.id}');
        _showErrorDialog("Thanh toán thất bại. Trạng thái: ${paymentIntent.status}");
      }
    } on StripeException catch (e) {
      log("StripeException during presentPaymentSheet: ${e.error.code} - ${e.error.message}");
      if (e.error.code == FailureCode.Canceled) {
        _showInfoDialog("Bạn đã hủy thanh toán.");
      } else {
        _showErrorDialog("Lỗi thanh toán Stripe: ${e.error.localizedMessage ?? e.error.message}");
      }
    } catch (e, s) {
      log("Generic error in _presentStripePaymentSheet: $e\n$s");
      _showErrorDialog("Đã xảy ra lỗi không xác định trong quá trình thanh toán: ${e.toString()}");
    }
  }

  Future<void> _confirmPaymentOnBackend(String paymentIntentIdToConfirm) async {
    final token = await _getAuthToken();
    if (token == null || token.contains("YOUR_FALLBACK_OR_TEST_BEARER_TOKEN") || token.isEmpty) {
      log("Auth token not available or is a placeholder. Cannot confirm payment on backend.");
      _showErrorDialog("Lỗi xác thực tài khoản. Không thể xác nhận vé với máy chủ. Vui lòng đăng nhập lại và thử.");
      return;
    }
    try {
      log("--- Confirming Payment on Backend ($_apiBaseUrl/stripe/confirm-payment) ---");
      log(jsonEncode({"payment_intent_id": paymentIntentIdToConfirm}));
      final response = await http.post(
        Uri.parse("$_apiBaseUrl/stripe/confirm-payment"),
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer $token",
        },
        body: jsonEncode({"payment_intent_id": paymentIntentIdToConfirm}),
      );
      log("Backend Confirm Payment API Response Status: ${response.statusCode}");
      log("Backend Confirm Payment API Response Body: ${response.body}");
      if (response.statusCode == 200) {
        log("Payment confirmed successfully on backend. Ticket ID: $_retrievedTicketId. PI Confirmed: $paymentIntentIdToConfirm");
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text("Đặt vé thành công! Mã vé của bạn là: ${_retrievedTicketId ?? 'N/A'}"),
          duration: Duration(seconds: 5),
        ));
        Navigator.of(context).pop(true);
      } else {
        log("Failed to confirm payment on backend. Status: ${response.statusCode}, Body: ${response.body}");
        _showErrorDialog("Lỗi xác nhận vé với máy chủ: ${response.reasonPhrase} (${response.statusCode}). Mã vé tạm thời: ${_retrievedTicketId ?? 'N/A'}. Chi tiết: ${response.body}");
      }
    } catch (e, s) {
      log("Exception calling backend confirm payment API: $e\n$s");
      _showErrorDialog("Lỗi kết nối khi xác nhận vé: ${e.toString()}. Mã vé tạm thời: ${_retrievedTicketId ?? 'N/A'}");
    }
  }

  Future<Map<String, dynamic>?> _createBankPaymentRequest() async {
    final prefs = await SharedPreferences.getInstance();
    final String? userId = prefs.getString('userId');
    if (_retrievedTicketId == null) {
      _showErrorDialog("Không có mã vé để tạo yêu cầu thanh toán.");
      return null;
    }
    final payload = {
      "customer_id": userId,
      "ticket_id": _retrievedTicketId,
      "amount": widget.rawTotalPrice.toDouble(),
      "currency": "VND",
    };
    try {
      final String jsonPayload = jsonEncode(payload);
      log("--- Creating Bank Payment Request ($_apiBaseUrl/bank/create-payment-request) ---");
      log(jsonPayload);
      final token = await _getAuthToken();
      final response = await http.post(
        Uri.parse("$_apiBaseUrl/bank/create-payment-request"),
        headers: {
          "Content-Type": "application/json",
          if (token != null) "Authorization": "Bearer $token",
        },
        body: jsonPayload,
      );
      log("Create Bank Payment Request API Response Status: ${response.statusCode}");
      log("Create Bank Payment Request API Response Body: ${response.body}");
      if (response.statusCode >= 200 && response.statusCode < 300) {
        final responseData = jsonDecode(response.body);
        if (responseData['data']?['invoice_id'] != null) {
          return responseData['data'];
        } else {
          _showErrorDialog("Phản hồi từ máy chủ không hợp lệ (thiếu invoice_id).");
          return null;
        }
      } else {
        _showErrorDialog("Lỗi tạo yêu cầu thanh toán qua ngân hàng: ${response.reasonPhrase} (${response.statusCode})");
        return null;
      }
    } catch (e, s) {
      log("Exception creating bank payment request: $e\n$s");
      _showErrorDialog("Lỗi kết nối khi tạo yêu cầu thanh toán: ${e.toString()}");
      return null;
    }
  }

  Future<bool?> _showBankTransferDetailsDialog(Map<String, dynamic> bankData) {
    final bankName = bankData['our_bank_name'] ?? 'Không rõ';
    final accountName = bankData['our_bank_account_name'] ?? 'Không rõ';
    final accountNumber = bankData['our_bank_account_number'] ?? 'Không rõ';
    final payableAmount = bankData['payable_amount'] ?? widget.rawTotalPrice;
    transferCode = bankData['bank_transfer_code'] ?? _retrievedTicketId;
    return showDialog<bool>(
      context: context,
      barrierDismissible: false,
      builder: (ctx) => AlertDialog(
        backgroundColor: Colors.white,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: Text("Thông tin chuyển khoản", style: TextStyle(fontWeight: FontWeight.bold)),
        content: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text("Vui lòng chuyển khoản chính xác số tiền và nội dung dưới đây:"),
              SizedBox(height: 16),
              _buildInfoRow("Ngân hàng", bankName),
              _buildInfoRow("Chủ tài khoản", accountName),
              _buildInfoRow("Số tài khoản", accountNumber),
              _buildInfoRow("Số tiền", currencyFormatter.format(payableAmount), valueColor: Colors.red),
              SizedBox(height: 16),
              Text("Nội dung chuyển khoản:", style: TextStyle(fontWeight: FontWeight.bold)),
              SizedBox(height: 4),
              Text(transferCode, style: TextStyle(color: Colors.blueAccent, fontWeight: FontWeight.bold, fontSize: 16)),
              SizedBox(height: 16),
              Text("Sau khi hoàn tất chuyển khoản, vui lòng nhấn nút bên dưới để xác nhận.", style: TextStyle(color: Colors.grey[700])),
            ],
          ),
        ),
        actionsPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        actionsAlignment: MainAxisAlignment.spaceBetween,
        actions: <Widget>[
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: Text("Hủy bỏ"),
            style: TextButton.styleFrom(foregroundColor: Colors.red, textStyle: TextStyle(fontWeight: FontWeight.bold)),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: Text("Tôi đã thanh toán"),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.green,
              foregroundColor: Colors.white,
              textStyle: TextStyle(fontWeight: FontWeight.bold),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value, {Color? valueColor}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(flex: 2, child: Text("$label:", style: TextStyle(fontWeight: FontWeight.bold))),
          Expanded(flex: 3, child: Text(value, style: TextStyle(fontWeight: FontWeight.w500, color: valueColor ?? Colors.black))),
        ],
      ),
    );
  }

  Future<bool> _confirmBankPaymentOnBackend(String invoiceId) async {
    final prefs = await SharedPreferences.getInstance();
    final token = _authRepository.getValidAccessToken();
    final String? userId = prefs.getString('userId');
    final String? fullName = prefs.getString('fullName');
    if (token == null || userId == null) {
      throw Exception('Authentication token or User ID not found.');
    }
    final payload = {
      "invoice_id": invoiceId,
      "amount_received": widget.rawTotalPrice.toDouble(),
      "currency_received": "VND",
      "confirmation_timestamp": DateTime.now().toUtc().toIso8601String(),
      "bank_transfer_code": transferCode,
      "payer_account_name": null,
      "payer_account_number": userId,
      "payer_bank_name": fullName,
      "confirmed_by": userId,
    };
    try {
      final String jsonPayload = jsonEncode(payload);
      log("--- Confirming Bank Payment ($_apiBaseUrl/bank/confirm-payment) ---");
      log(jsonPayload);
      final token = await _getAuthToken();
      final response = await http.post(
        Uri.parse("$_apiBaseUrl/bank/confirm-payment"),
        headers: {
          "Content-Type": "application/json",
          if (token != null) "Authorization": "Bearer $token",
        },
        body: jsonPayload,
      );
      log("Confirm Bank Payment API Response Status: ${response.statusCode}");
      log("Confirm Bank Payment API Response Body: ${response.body}");
      if (response.statusCode >= 200 && response.statusCode < 300) {
        log("Bank payment confirmation sent successfully.");
        return true;
      } else {
        _showErrorDialog("Lỗi gửi xác nhận thanh toán: ${response.reasonPhrase} (${response.statusCode})");
        return false;
      }
    } catch (e, s) {
      log("Exception confirming bank payment: $e\n$s");
      _showErrorDialog("Lỗi kết nối khi gửi xác nhận thanh toán: ${e.toString()}");
      return false;
    }
  }

  void _showErrorDialog(String message) {
    if (!mounted) return;
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text("Đã xảy ra lỗi"),
        content: Text(message),
        actions: <Widget>[
          TextButton(
            child: Text("OK"),
            onPressed: () {
              Navigator.of(ctx).pop();
            },
          ),
        ],
      ),
    );
  }

  void _showInfoDialog(String message) {
    if (!mounted) return;
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text("Thông báo"),
        content: Text(message),
        actions: <Widget>[
          TextButton(
            child: Text("OK"),
            onPressed: () {
              Navigator.of(ctx).pop();
            },
          ),
        ],
      ),
    );
  }

  Future<void> _loadUserData() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    fullName = prefs.getString('fullName') ?? "";
    phoneNumber = prefs.getString('phoneNumber') ?? "";
    email = prefs.getString('username') ?? "";
    _nameController.text = fullName;
    _phoneController.text = phoneNumber;
    _emailController.text = email;
  }

  Future<void> _initializeScreenData() async {
    await _loadUserData();
    try {
      final departureFullRoute = await _fetchTripFullRoute(widget.tripId);
      _parseFullRoute(
        departureFullRoute,
        isReturn: false,
        initialPickup: widget.initialPickupLocation,
        initialDropoff: widget.initialDropoffLocation,
      );
      if (widget.returnTripId != null) {
        log("Fetched fullRoute for trip ${widget.returnTripId}: ${widget.returnTripInfo}");
        final returnFullRoute = await _fetchTripFullRoute(widget.returnTripId!);
        _parseFullRoute(
          returnFullRoute,
          isReturn: true,
          initialPickup: widget.initialPickupLocationEnd,
          initialDropoff: widget.initialDropoffLocationEnd,
        );
      }
    } catch (e) {
      log("Error fetching trip data: $e");
      _showErrorDialog("Không thể tải dữ liệu chuyến đi. Vui lòng thử lại.");
      _parseFullRoute(
        widget.fullRoute,
        isReturn: false,
        initialPickup: widget.initialPickupLocation,
        initialDropoff: widget.initialDropoffLocation,
      );
    }
    setState(() {
      _isLoading = false;
    });
  }

  Future<String> _fetchTripFullRoute(int tripId) async {
    final response = await http.get(Uri.parse('$_apiBaseUrl/trips/$tripId/seats'));
    if (response.statusCode == 200) {
      final responseData = jsonDecode(response.body);
      final fullRoute = responseData['data']?['fullRoute'] as String?;
      if (fullRoute != null) {
        log("Fetched fullRoute for trip $tripId: $fullRoute");
        return fullRoute;
      }
    }
    throw Exception('Failed to load trip data for tripId $tripId');
  }

  Map<String, String> _extractStationInfo(String fullNameWithId) {
    if (fullNameWithId.isEmpty || fullNameWithId == "N/A") {
      return {'name': fullNameWithId, 'id': '', 'fullName': fullNameWithId};
    }
    RegExp regex = RegExp(r"^(.*) \((\d+)\)$");
    Match? match = regex.firstMatch(fullNameWithId);
    if (match != null && match.groupCount == 2) {
      return {'name': match.group(1)!.trim(), 'id': match.group(2)!.trim(), 'fullName': fullNameWithId};
    }
    log("Warning: Could not parse station ID from '$fullNameWithId'. Using full name as fallback.");
    return {'name': fullNameWithId, 'id': '', 'fullName': fullNameWithId};
  }

  void _parseFullRoute(String fullRoute, {required bool isReturn, required String initialPickup, required String initialDropoff}) {
    final stopsInRoute = fullRoute.split('→').map((e) => e.trim()).where((s) => s.isNotEmpty).toList();
    List<String> parsedPickupLocations = [];
    List<String> parsedDropoffLocations = [];
    if (stopsInRoute.isNotEmpty) {
      if (stopsInRoute.length > 1) {
        parsedPickupLocations.addAll(stopsInRoute.sublist(0, stopsInRoute.length - 1));
        parsedDropoffLocations.addAll(stopsInRoute.sublist(1));
      } else {
        parsedPickupLocations.add(stopsInRoute.first);
        parsedDropoffLocations.add(stopsInRoute.first);
      }
    }
    String findFullName(List<String> list, String name) {
      return list.firstWhere((item) => _extractStationInfo(item)['name'] == name, orElse: () => list.isNotEmpty ? list.first : "N/A");
    }

    String currentSelectedPickup = findFullName(parsedPickupLocations, _extractStationInfo(initialPickup)['name']!);
    String currentSelectedDropoff = findFullName(parsedDropoffLocations, _extractStationInfo(initialDropoff)['name']!);
    if (isReturn) {
      _parsedReturnPickupLocations = parsedPickupLocations;
      _parsedReturnDropoffLocations = parsedDropoffLocations;
      selectedReturnPickupLocation = currentSelectedPickup;
      selectedReturnDropoffLocation = currentSelectedDropoff;
    } else {
      _parsedPickupLocations = parsedPickupLocations;
      _parsedDropoffLocations = parsedDropoffLocations;
      selectedPickupLocation = currentSelectedPickup;
      selectedDropoffLocation = currentSelectedDropoff;
    }
  }

  void _showLocationPicker(BuildContext context, bool isPickup, {bool isReturn = false}) {
    final locationsFullName = isReturn ? (isPickup ? _parsedReturnPickupLocations : _parsedReturnDropoffLocations) : (isPickup ? _parsedPickupLocations : _parsedDropoffLocations);
    log("Locations Full Name: $_parsedPickupLocations");
    if (locationsFullName.isEmpty || (locationsFullName.length == 1 && _extractStationInfo(locationsFullName.first)['name'] == "N/A")) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text("Không có điểm ${isPickup ? 'đón' : 'trả'} khả dụng.")));
      return;
    }
    String? currentSelectedFullName = isReturn ? (isPickup ? selectedReturnPickupLocation : selectedReturnDropoffLocation) : (isPickup ? selectedPickupLocation : selectedDropoffLocation);
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return ConstrainedBox(
              constraints: BoxConstraints(maxHeight: MediaQuery.of(context).size.height * 0.5),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.only(topLeft: Radius.circular(16), topRight: Radius.circular(16)),
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(isPickup ? "Chọn điểm đón" : "Chọn điểm trả", style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                          IconButton(icon: Icon(Icons.close, color: Colors.grey), onPressed: () => Navigator.pop(context)),
                        ],
                      ),
                    ),
                    Divider(height: 1),
                    Expanded(
                      child: ListView.builder(
                        shrinkWrap: true,
                        itemCount: locationsFullName.length,
                        itemBuilder: (context, index) {
                          final locationFullName = locationsFullName[index];
                          if (locationFullName == "N/A" && locationsFullName.length > 1) return SizedBox.shrink();
                          final stationInfo = _extractStationInfo(locationFullName);
                          final displayName = stationInfo['name']!;
                          if (displayName == "N/A" && locationsFullName.length > 1) return SizedBox.shrink();
                          final isSelected = currentSelectedFullName == locationFullName;
                          return ListTile(
                            title: Text(displayName),
                            trailing: Container(
                              width: 24,
                              height: 24,
                              decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                border: Border.all(color: Colors.grey.shade400),
                                color: isSelected ? Colors.blueAccent : Colors.transparent,
                              ),
                              child: isSelected ? Icon(Icons.check, size: 16, color: Colors.white) : null,
                            ),
                            onTap: () {
                              setModalState(() {
                                currentSelectedFullName = locationFullName;
                              });
                              setState(() {
                                if (isReturn) {
                                  if (isPickup) {
                                    selectedReturnPickupLocation = locationFullName;
                                  } else {
                                    selectedReturnDropoffLocation = locationFullName;
                                  }
                                } else {
                                  if (isPickup) {
                                    selectedPickupLocation = locationFullName;
                                  } else {
                                    selectedDropoffLocation = locationFullName;
                                  }
                                }
                              });
                              Navigator.pop(context);
                            },
                          );
                        },
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }

  void _showEditCustomerInfoDialog() {
    final tempNameController = TextEditingController(text: _nameController.text);
    final tempPhoneController = TextEditingController(text: _phoneController.text);
    final tempEmailController = TextEditingController(text: _emailController.text);
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text("Chỉnh sửa thông tin"),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(controller: tempNameController, decoration: InputDecoration(labelText: "Họ tên")),
                TextField(controller: tempPhoneController, decoration: InputDecoration(labelText: "Số điện thoại"), keyboardType: TextInputType.phone),
                TextField(controller: tempEmailController, decoration: InputDecoration(labelText: "Email"), keyboardType: TextInputType.emailAddress),
              ],
            ),
          ),
          actions: [
            TextButton(onPressed: () => Navigator.pop(context), child: Text("Hủy", style: TextStyle(color: Colors.grey[700]))),
            TextButton(
              onPressed: () {
                setState(() {
                  _nameController.text = tempNameController.text.trim();
                  _phoneController.text = tempPhoneController.text.trim();
                  _emailController.text = tempEmailController.text.trim();
                });
                Navigator.pop(context);
              },
              child: Text("Lưu", style: TextStyle(color: Colors.blueAccent)),
            ),
          ],
        );
      },
    );
  }

  Widget _buildCustomerInfoCard(){
      return Card(
      color: Colors.white,
      elevation: 2,
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text("Thông tin hành khách", style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.blueAccent)),
                IconButton(
                  icon: Icon(Icons.edit, color: Colors.orange),
                  tooltip: "Chỉnh sửa thông tin",
                  onPressed: _showEditCustomerInfoDialog,
                )
              ],
            ),
            SizedBox(height: 8),
            Row(children: [
              Icon(Icons.person_outline, color: Colors.grey[700], size: 20),
              SizedBox(width: 8),
              Expanded(child: Text(_nameController.text, style: TextStyle(fontSize: 15))),
            ]),
            SizedBox(height: 6),
            Row(children: [
              Icon(Icons.phone_outlined, color: Colors.grey[700], size: 20),
              SizedBox(width: 8),
              Text(_phoneController.text, style: TextStyle(fontSize: 15)),
            ]),
            SizedBox(height: 6),
            Row(children: [
              Icon(Icons.email_outlined, color: Colors.grey[700], size: 20),
              SizedBox(width: 8),
              Expanded(child: Text(_emailController.text, style: TextStyle(fontSize: 15))),
            ]),
          ],
        ),
      ),
    );
  }

  Widget _buildTripInfoWidget() {
    final pickupInfo = _extractStationInfo(selectedPickupLocation);
    final dropoffInfo = _extractStationInfo(selectedDropoffLocation);
    return Card(
      color: Colors.white,
      elevation: 2,
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text("Thông tin chuyến đi", style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.blueAccent)),
            SizedBox(height: 12),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text("Tổng tiền:", style: TextStyle(color: Colors.grey[700])),
              Text(currencyFormatter.format(widget.rawTotalPrice), style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
            ]),
            SizedBox(height: 8),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text("Loại ghế:", style: TextStyle(color: Colors.grey[700])),
              Text(widget.seatType, style: TextStyle(fontWeight: FontWeight.w500)),
            ]),
            SizedBox(height: 8),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text("Thời gian:", style: TextStyle(color: Colors.grey[700])),
              Flexible(child: Text(widget.time, textAlign: TextAlign.end, style: TextStyle(fontWeight: FontWeight.w500))),
            ]),
            SizedBox(height: 16),
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Column(mainAxisAlignment: MainAxisAlignment.center, children: [
                  SizedBox(height: 4),
                  Icon(Icons.my_location, color: Colors.green, size: 20),
                  SizedBox(width: 24, height: 40, child: CustomPaint(painter: DashedLinePainter())),
                  Icon(Icons.location_on, color: Colors.red, size: 20),
                  SizedBox(height: 4),
                ]),
                SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(pickupInfo['name']!, style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500)),
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 8.0),
                        child: Text(widget.estimatedDistance, style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                      ),
                      Text(dropoffInfo['name']!, style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500)),
                    ],
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLocationSelectionRadio(String title, String groupValue, String value, Function(String?) onChanged) {
    return GestureDetector(
      onTap: () => onChanged(value),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Radio<String>(value: value, groupValue: groupValue, onChanged: onChanged, materialTapTargetSize: MaterialTapTargetSize.shrinkWrap, activeColor: Colors.blueAccent),
          Text(title, style: TextStyle(fontSize: 15)),
        ],
      ),
    );
  }

  Widget _buildPickupDropoffWidget({bool isReturn = false}) {
    final String currentPickupLoc = isReturn ? (selectedReturnPickupLocation ?? "N/A") : selectedPickupLocation;
    final String currentDropoffLoc = isReturn ? (selectedReturnDropoffLocation ?? "N/A") : selectedDropoffLocation;
    final currentPickupDisplayInfo = _extractStationInfo(currentPickupLoc);
    final currentDropoffDisplayInfo = _extractStationInfo(currentDropoffLoc);
    return Card(
      color: Colors.white,
      elevation: 2,
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(isReturn ? "Thông tin đón trả chuyến về" : "Thông tin đón trả", style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.blueAccent)),
            SizedBox(height: 16),
            Text("Điểm đón hành khách", style: TextStyle(fontSize: 15, fontWeight: FontWeight.w600)),
            SizedBox(height: 4),
            Row(children: [
              _buildLocationSelectionRadio("Bến xe/VP", isReturn ? selectedReturnPickupType : selectedPickupType, "Bến xe/VP", (val) {
                if (val == null) return;
                setState(() {
                  if (isReturn) {
                    selectedReturnPickupType = val;
                    showReturnPickupLocationField = true;
                  } else {
                    selectedPickupType = val;
                    showPickupLocationField = true;
                  }
                });
              }),
              SizedBox(width: 10),
              _buildLocationSelectionRadio("Trung chuyển", isReturn ? selectedReturnPickupType : selectedPickupType, "Trung chuyển", (val) {
                if (val == null) return;
                setState(() {
                  if (isReturn) {
                    selectedReturnPickupType = val;
                    showReturnPickupLocationField = true;
                  } else {
                    selectedPickupType = val;
                    showPickupLocationField = true;
                  }
                });
              }),
            ]),
            if (isReturn ? showReturnPickupLocationField : showPickupLocationField) ...[
              SizedBox(height: 8),
              GestureDetector(
                onTap: () => _showLocationPicker(context, true, isReturn: isReturn),
                child: Container(
                  padding: EdgeInsets.symmetric(vertical: 12, horizontal: 16),
                  decoration: BoxDecoration(border: Border.all(color: Colors.grey.shade400), borderRadius: BorderRadius.circular(8), color: Colors.grey[50]),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Expanded(
                        child: Text(
                          currentPickupDisplayInfo['name']! == "N/A" ? "Chọn điểm đón" : currentPickupDisplayInfo['name']!,
                          overflow: TextOverflow.ellipsis,
                          style: TextStyle(fontSize: 15, color: currentPickupDisplayInfo['name']! == "N/A" ? Colors.grey : Colors.black),
                        ),
                      ),
                      Icon(Icons.arrow_drop_down, color: Colors.grey[700]),
                    ],
                  ),
                ),
              ),
            ],
            SizedBox(height: 20),
            Text("Điểm trả hành khách", style: TextStyle(fontSize: 15, fontWeight: FontWeight.w600)),
            SizedBox(height: 4),
            Row(children: [
              _buildLocationSelectionRadio("Bến xe/VP", isReturn ? selectedReturnDropoffType : selectedDropoffType, "Bến xe/VP", (val) {
                if (val == null) return;
                setState(() {
                  if (isReturn) {
                    selectedReturnDropoffType = val;
                    showReturnDropoffLocationField = true;
                  } else {
                    selectedDropoffType = val;
                    showDropoffLocationField = true;
                  }
                });
              }),
              SizedBox(width: 10),
              _buildLocationSelectionRadio("Trung chuyển", isReturn ? selectedReturnDropoffType : selectedDropoffType, "Trung chuyển", (val) {
                if (val == null) return;
                setState(() {
                  if (isReturn) {
                    selectedReturnDropoffType = val;
                    showReturnDropoffLocationField = true;
                  } else {
                    selectedDropoffType = val;
                    showDropoffLocationField = true;
                  }
                });
              }),
            ]),
            if (isReturn ? showReturnDropoffLocationField : showDropoffLocationField) ...[
              SizedBox(height: 8),
              GestureDetector(
                onTap: () => _showLocationPicker(context, false, isReturn: isReturn),
                child: Container(
                  padding: EdgeInsets.symmetric(vertical: 12, horizontal: 16),
                  decoration: BoxDecoration(border: Border.all(color: Colors.grey.shade400), borderRadius: BorderRadius.circular(8), color: Colors.grey[50]),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Expanded(
                        child: Text(
                          currentDropoffDisplayInfo['name']! == "N/A" ? "Chọn điểm trả" : currentDropoffDisplayInfo['name']!,
                          overflow: TextOverflow.ellipsis,
                          style: TextStyle(fontSize: 15, color: currentDropoffDisplayInfo['name']! == "N/A" ? Colors.grey : Colors.black),
                        ),
                      ),
                      Icon(Icons.arrow_drop_down, color: Colors.grey[700]),
                    ],
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildSeatInfoWidget({bool isReturn = false}) {
    final seatNames = isReturn ? widget.selectedReturnSeatsNames : widget.selectedSeatsNames;
    return Card(
      color: Colors.white,
      elevation: 2,
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    isReturn ? "Ghế đã chọn (chuyến về) (${seatNames?.length ?? 0})" : "Ghế đã chọn (${seatNames?.length ?? 0})",
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.blueAccent),
                  ),
                  SizedBox(height: 8),
                  Text(
                    seatNames != null && seatNames.isNotEmpty ? seatNames.join(", ") : "Chưa chọn ghế",
                    style: TextStyle(fontSize: 15, color: seatNames == null || seatNames.isEmpty ? Colors.grey[600] : Colors.black),
                    overflow: TextOverflow.ellipsis,
                    maxLines: 2,
                  ),
                ],
              ),
            ),
            IconButton(
              icon: Icon(Icons.edit_note, color: Colors.orange, size: 28),
              tooltip: "Thay đổi ghế",
              onPressed: () {
                if (Navigator.canPop(context)) {
                  Navigator.pop(context);
                }
              },
            ),
          ],
        ),
      ),
    );
  }
}
