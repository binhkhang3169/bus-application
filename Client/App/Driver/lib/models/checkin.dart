class Checkin {
  final String ticketId;
  final String tripId;
  final String? seatName; // Có thể null
  final DateTime? checkedInAt; // Có thể null
  final String checkinNote;
  final String message;

  Checkin({
    required this.ticketId,
    required this.tripId,
    this.seatName,
    this.checkedInAt,
    required this.checkinNote,
    required this.message,
  });

  factory Checkin.fromJson(Map<String, dynamic> json) {
    return Checkin(
      ticketId: json['ticket_id'],
      tripId: json['trip_id'],
      seatName: json['seat_name'],
      // Chuyển đổi chuỗi thời gian thành DateTime, xử lý trường hợp null
      checkedInAt: json['checked_in_at'] != null
          ? DateTime.parse(json['checked_in_at'])
          : null,
      checkinNote: json['checkin_note'] ?? '',
      message: json['message'] ?? '',
    );
  }
}