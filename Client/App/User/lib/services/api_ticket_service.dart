import 'dart:math';

import 'package:caoky/models/api_response.dart';
import 'package:caoky/models/api_response1.dart';
import 'package:caoky/models/trip/address_info.dart';
import 'package:caoky/models/trip/station.dart';
import 'package:caoky/models/trip/trip.dart';
import 'package:caoky/models/trip/trip_info.dart';
// lib/services/api_trip_service.dart

// Note: Your original file had many commented-out imports.
// I'm only including those relevant to the defined methods.
import 'package:caoky/models/api_response.dart';
import 'package:caoky/models/trip/station.dart'; // User provided Station model
import 'package:caoky/models/trip/address_info.dart'; // From your provided service
// import 'package:caoky/models/trip/trip_info.dart'; // Not directly used as return type here
// import 'package:caoky/models/trip/trip.dart'; // Not directly used as return type here

import 'package:dio/dio.dart';
import 'package:retrofit/http.dart'; // Changed from retrofit.dart to http.dart as per common practice

part 'api_ticket_service.g.dart'; // Ensure you run build_runner

@RestApi(baseUrl: "http://57.155.76.74/api/v1") // Your ngrok URL
abstract class ApiTicketService {
  factory ApiTicketService(Dio dio, {String baseUrl}) = _ApiTicketService;

  @GET("/tickets-available/{id}")
  Future<ApiResponse1<List<TripInfo>>> getTicketsAvailable(@Path("id") int id);

  // Other commented methods from your original file can be added here if needed.
}

// To generate/update the .g.dart file, run in your terminal:
// flutter pub run build_runner build --delete-conflicting-outputs
