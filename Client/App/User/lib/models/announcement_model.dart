import 'dart:convert';

// Hàm helper để parse danh sách JSON
List<Announcement> announcementFromJson(String str) => List<Announcement>.from(json.decode(str).map((x) => Announcement.fromJson(x)));

class Announcement {
    final String id;
    final String title;
    final String imageUrl;
    final String description; // Ánh xạ từ trường 'content' của API
    final String createdBy;
    final DateTime createdAt;

    Announcement({
        required this.id,
        required this.title,
        required this.imageUrl,
        required this.description,
        required this.createdBy,
        required this.createdAt,
    });

    // Factory constructor để parse JSON
    factory Announcement.fromJson(Map<String, dynamic> json) => Announcement(
        id: json["id"],
        title: json["title"],
        imageUrl: json["image_url"],
        description: json["content"], // Chú ý: 'content' từ API được ánh xạ vào 'description'
        createdBy: json["created_by"],
        createdAt: DateTime.parse(json["created_at"]),
    );
}