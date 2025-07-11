// lib/ui/stripe_payment_screen.dart
import 'dart:convert';
import 'dart:developer';
import 'package:flutter/material.dart';
import 'package:flutter_stripe/flutter_stripe.dart';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart'; // For currency formatting

class StripePaymentScreen extends StatefulWidget {
  final int amount; // Amount in the smallest currency unit (e.g., 100000 for 100,000 VND)
  final String currency; // e.g., 'VND'
  final String? customerEmail;
  final String? customerName;
  final Map<String, dynamic> ticketBookingPayload; // Payload for your backend after payment

  const StripePaymentScreen({
    Key? key,
    required this.amount,
    required this.currency,
    this.customerEmail,
    this.customerName,
    required this.ticketBookingPayload,
  }) : super(key: key);

  @override
  _StripePaymentScreenState createState() => _StripePaymentScreenState();
}

class _StripePaymentScreenState extends State<StripePaymentScreen> {
  bool _isLoading = true;
  Map<String, dynamic>? _paymentIntentData;
  final NumberFormat currencyFormatter = NumberFormat.currency(locale: 'vi_VN', symbol: '₫');

  @override
  void initState() {
    super.initState();
    log("[StripeScreen] initState: Amount=${widget.amount}, Currency=${widget.currency}, Email=${widget.customerEmail}");
    _initiatePaymentProcess();
  }

  Future<void> _initiatePaymentProcess() async {
    if (!mounted) return;
    setState(() => _isLoading = true);
    log("[StripeScreen] Starting payment process...");

    try {
      // STEP 1: Create Payment Intent on your backend
      log("[StripeScreen] Attempting to create Payment Intent on backend.");
      _paymentIntentData = await _createPaymentIntentOnBackend(
        widget.amount,
        widget.currency,
        widget.customerEmail,
      );

      if (_paymentIntentData == null || _paymentIntentData!['clientSecret'] == null) {
        log("[StripeScreen] Error: Failed to create Payment Intent from backend or clientSecret is null.");
        _showPaymentDialog("Payment Setup Failed", "Could not initialize payment with the server. Please check logs and try again.");
        if (mounted) setState(() => _isLoading = false);
        return;
      }
      log("[StripeScreen] Backend response for Payment Intent (raw): ${jsonEncode(_paymentIntentData)}");
      log("[StripeScreen] Extracted clientSecret: ${_paymentIntentData!['clientSecret']}");
      log("[StripeScreen] Extracted customerId (if any): ${_paymentIntentData!['customer']}");
      log("[StripeScreen] Extracted ephemeralKey (if any): ${_paymentIntentData!['ephemeralKey']}");


      // STEP 2: Initialize Stripe Payment Sheet
      log("[StripeScreen] Initializing Payment Sheet with clientSecret: ${_paymentIntentData!['clientSecret']}");
      await Stripe.instance.initPaymentSheet(
        paymentSheetParameters: SetupPaymentSheetParameters(
          merchantDisplayName: 'Cao Ky Bus Service', // Your business name
          paymentIntentClientSecret: _paymentIntentData!['clientSecret'],
          customerId: _paymentIntentData!['customer'], // Optional: if your backend provides it
          customerEphemeralKeySecret: _paymentIntentData!['ephemeralKey'], // Optional: if your backend provides it
          style: ThemeMode.system,
          // testEnv: true, // For Google Pay in test environment
          // applePay: PaymentSheetApplePay(merchantCountryCode: 'VN'), // If using Apple Pay
          // googlePay: PaymentSheetGooglePay(merchantCountryCode: 'VN', testEnv: true), // If using Google Pay
        ),
      );
      log("[StripeScreen] Stripe Payment Sheet initialized successfully.");

      // STEP 3: Present Payment Sheet
      if (mounted) setState(() => _isLoading = false); // Allow UI to update before presenting sheet
      await _presentPaymentSheet();

    } catch (e, s) {
      log("[StripeScreen] Exception during payment initialization or sheet presentation: $e\n$s", error: e, stackTrace: s);
      if (mounted) {
        _showPaymentDialog("Payment Error", "An unexpected error occurred: ${e.toString()}. Check logs.");
        setState(() => _isLoading = false);
      }
    }
  }

  // !!! THIS IS A PLACEHOLDER - YOU MUST IMPLEMENT YOUR BACKEND LOGIC !!!
  Future<Map<String, dynamic>?> _createPaymentIntentOnBackend(
    int amount, String currency, String? email) async {
    
    // Replace with your actual backend endpoint
    const String backendUrl = 'YOUR_BACKEND_ENDPOINT/create-payment-intent'; 
    log("[StripeScreen_BackendCall] Requesting Payment Intent from: $backendUrl");
    log("[StripeScreen_BackendCall] Request body: amount=$amount, currency=$currency, email=$email");

    // --- SIMULATED BACKEND RESPONSE (FOR TESTING WITHOUT A REAL BACKEND) ---
    // Remove or comment this out when you have a real backend
    if (backendUrl.contains('YOUR_BACKEND_ENDPOINT')) {
      log("[StripeScreen_BackendCall] WARNING: Using SIMULATED backend response. Replace with actual backend call.");
      await Future.delayed(Duration(seconds: 1)); // Simulate network delay
      // This is what your backend should ideally return:
      return {
        "clientSecret": "pi_3 esimerkki_secret_esimerkki", // Replace with a test PaymentIntent client_secret from Stripe dashboard if needed for UI testing ONLY
        "ephemeralKey": "ek_test_esimerkki", // Replace with a test ephemeral key
        "customer": "cus_esimerkki", // Replace with a test customer ID
        "publishableKey": Stripe.publishableKey // Often not needed from backend if already set
      };
    }
    // --- END OF SIMULATED BACKEND RESPONSE ---


    try {
      final response = await http.post(
        Uri.parse(backendUrl),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'amount': amount, // Your backend should handle this as per Stripe's requirements for the currency
          'currency': currency,
          if (email != null) 'email': email,
          'customer_name': widget.customerName, // Example of sending more data
          // Add any other parameters your backend needs (e.g., items, metadata)
        }),
      );

      log("[StripeScreen_BackendCall] Backend Response Status: ${response.statusCode}");
      log("[StripeScreen_BackendCall] Backend Response Body: ${response.body}");

      if (response.statusCode == 200) {
        final Map<String, dynamic> responseData = jsonDecode(response.body);
        log("[StripeScreen_BackendCall] Successfully fetched Payment Intent data from backend.");
        return responseData;
      } else {
        log("[StripeScreen_BackendCall] Failed to fetch Payment Intent. Status: ${response.statusCode}, Reason: ${response.reasonPhrase}");
        return null;
      }
    } catch (e, s) {
      log("[StripeScreen_BackendCall] Error calling backend: $e\n$s", error: e, stackTrace: s);
      return null;
    }
  }

  Future<void> _presentPaymentSheet() async {
    if (!mounted) return;
    try {
      log("[StripeScreen] Presenting Payment Sheet to user...");
      await Stripe.instance.presentPaymentSheet();
      
      // If presentPaymentSheet completes without throwing an exception,
      // it means the user completed the flow (could be success or a Stripe-handled decline).
      // A true success needs backend verification.
      log("[StripeScreen] Payment Sheet flow completed by user.");

      // --- IMPORTANT ---
      // At this stage, for a production app, you MUST verify the payment status
      // on your backend by using the PaymentIntent ID to prevent fraud.
      // The client-side success is not enough.
      // -----------------
      
      String paymentIntentId = _paymentIntentData?['id'] ?? 
                               _paymentIntentData?['paymentIntent'] ?? // Some backends might return it as 'paymentIntent'
                               'N/A (clientSecret parse needed)';
      if(_paymentIntentData?['clientSecret'] != null && paymentIntentId.contains('clientSecret parse needed')) {
        paymentIntentId = _paymentIntentData!['clientSecret'].split('_secret_').first;
      }

      log("[StripeScreen] Frontend indication: Payment Succeeded. PaymentIntent ID: $paymentIntentId");
      log("[StripeScreen] NEXT STEP: Verify payment on backend, then send booking payload to your API:");
      log("[StripeScreen] Ticket Booking Payload: ${jsonEncode(widget.ticketBookingPayload)}");
      
      _showPaymentDialog(
        "Payment Successful (Client-side)", 
        "Payment processed. Next: Backend verification & booking finalization.\nPayment Intent: $paymentIntentId",
        isSuccess: true
      );
      // TODO: Call your API to create the ticket using widget.ticketBookingPayload AND paymentIntentId

    } on StripeException catch (e) {
      log("[StripeScreen] StripeException during presentPaymentSheet: Code: ${e.error.code}, Message: ${e.error.message}, Localized: ${e.error.localizedMessage}", error: e);
      if (e.error.code == FailureCode.Canceled) {
        log("[StripeScreen] Payment was canceled by the user.");
        _showPaymentDialog("Payment Canceled", "You canceled the payment process.");
      } else {
        _showPaymentDialog("Payment Failed", "Error: ${e.error.localizedMessage ?? e.error.message ?? 'Unknown Stripe error'}");
      }
    } catch (e, s) {
      log("[StripeScreen] Generic exception during presentPaymentSheet: $e\n$s", error: e, stackTrace: s);
      _showPaymentDialog("Payment Error", "An unexpected error occurred while processing your payment.");
    }
  }

  void _showPaymentDialog(String title, String message, {bool isSuccess = false}) {
     if (!mounted) return;
    showDialog(
      context: context,
      barrierDismissible: false, // User must tap button
      builder: (BuildContext dialogContext) => AlertDialog(
        title: Text(title),
        content: Text(message),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.of(dialogContext).pop(); // Close this dialog
              if (isSuccess) {
                // Navigate back past PaymentScreen and StripePaymentScreen
                int popCount = 0;
                Navigator.of(context).popUntil((_) => popCount++ >= 2); 
                // Consider navigating to a dedicated booking success/summary screen:
                // Navigator.of(context).pushAndRemoveUntil(
                //   MaterialPageRoute(builder: (context) => BookingConfirmationScreen(payload: widget.ticketBookingPayload)),
                //   (route) => route.isFirst, // Remove all previous routes
                // );
              } else if (title.contains("Setup Failed") || title.contains("Payment Error")) {
                 // Stay on StripeScreen to allow retry, or pop to let user fix details on PaymentScreen
                 // If setup failed, likely pop back to PaymentScreen
                 if (title.contains("Setup Failed")) Navigator.of(context).pop();

              }
            },
            child: Text("OK"),
          )
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Complete Your Payment"),
        leading: IconButton(
          icon: Icon(Icons.arrow_back),
          onPressed: _isLoading ? null : () { // Disable back if loading to prevent inconsistent state
            log("[StripeScreen] User tapped app bar back button.");
            Navigator.of(context).pop();
          },
        ),
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: _isLoading
              ? Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    CircularProgressIndicator(),
                    SizedBox(height: 20),
                    Text("Initializing payment gateway...\nPlease wait.", textAlign: TextAlign.center),
                  ],
                )
              : Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Icon(Icons.credit_card, size: 80, color: Theme.of(context).colorScheme.primary),
                    SizedBox(height: 24),
                    Text(
                      "Amount to Pay: ${currencyFormatter.format(widget.amount)}",
                      textAlign: TextAlign.center,
                      style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                    ),
                    SizedBox(height: 12),
                    Text(
                      _paymentIntentData != null && _paymentIntentData!['clientSecret'] != null
                          ? "Stripe payment sheet is ready or has been presented. Follow the prompts from Stripe."
                          : "There was an issue setting up the payment. Please try again.",
                      textAlign: TextAlign.center,
                      style: TextStyle(fontSize: 16, color: Colors.grey[700]),
                    ),
                    SizedBox(height: 30),
                    if (_paymentIntentData == null || _paymentIntentData!['clientSecret'] == null)
                      ElevatedButton.icon(
                        icon: Icon(Icons.refresh),
                        label: Text("Retry Payment Setup"),
                        onPressed: _initiatePaymentProcess,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.orange,
                          padding: EdgeInsets.symmetric(vertical: 12),
                          textStyle: TextStyle(fontSize: 16)
                        ),
                      ),
                    // The PaymentSheet is usually presented automatically by _initiatePaymentProcess.
                    // If you wanted a button to manually re-present it (e.g., if the user closed it accidentally
                    // and Stripe allows re-presentation without re-init), you could add it here.
                    // However, usually re-initializing the process is safer if the first attempt failed.
                  ],
                ),
        ),
      ),
    );
  }
}