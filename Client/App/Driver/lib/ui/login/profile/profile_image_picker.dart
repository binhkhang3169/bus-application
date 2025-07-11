import 'dart:io';
import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:image_picker/image_picker.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:taixe/services/api_user_service.dart';

class ProfileImagePicker extends StatefulWidget {
  final String? imageUrl;

  const ProfileImagePicker({super.key, this.imageUrl});
  @override
  _ProfileImagePickerState createState() => _ProfileImagePickerState();
}

class _ProfileImagePickerState extends State<ProfileImagePicker> {
  File? _selectedImage;
  String? _imageUrl;
  bool _isUploading = false;
  late ApiUserService apiService;
  String token = "";
  String imageUrl = "";

  @override
  void initState() {
    super.initState();
    apiService = ApiUserService(Dio());
    _loadUserData();
  }

  Future<void> _loadUserData() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    setState(() {
      token = prefs.getString('accessToken') ?? "";
    });
  }

  Future<void> _pickImage() async {
    final picker = ImagePicker();
    final pickedFile = await picker.pickImage(source: ImageSource.gallery);

    if (pickedFile != null) {
      setState(() {
        _selectedImage = File(pickedFile.path);
      });

      await uploadImageToCloudinary();
    }
  }

  Future<void> uploadImageToCloudinary() async {
    if (_selectedImage == null) return;

    setState(() => _isUploading = true);

    try {
      String cloudinaryUrl = "https://api.cloudinary.com/v1_1/dwskd7iqr/upload";
      String uploadPreset = "flutter";

      var request = http.MultipartRequest('POST', Uri.parse(cloudinaryUrl));
      request.fields['upload_preset'] = uploadPreset;
      request.files.add(
        await http.MultipartFile.fromPath('file', _selectedImage!.path),
      );

      var response = await request.send();
      var responseData = await response.stream.bytesToString();
      var jsonResponse = json.decode(responseData);

      if (response.statusCode == 200) {
        imageUrl = jsonResponse['secure_url'];

        setState(() {
          _imageUrl = imageUrl;
        });

        print("✅ Upload thành công: $imageUrl");

        if (token != null && token.isNotEmpty) {
          await apiService.changeImage(token, imageUrl);
          print("✅ Gửi ảnh về API thành công");
        } else {
          print("⚠️ Token không tồn tại");
        }
      } else {
        print("❌ Lỗi khi upload: ${jsonResponse['error']['message']}");
      }
    } catch (e) {
      print("❌ Lỗi upload: $e");
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('image', imageUrl);
    }

    setState(() => _isUploading = false);
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: (token != null) ? _pickImage : null,
      child: Container(
        width: 80,
        height: 80,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: Colors.white,
          image:
              (_imageUrl != null && _imageUrl!.isNotEmpty)
                  ? DecorationImage(
                    image: NetworkImage(_imageUrl!),
                    fit: BoxFit.cover,
                  )
                  : (widget.imageUrl != null && widget.imageUrl!.isNotEmpty)
                  ? DecorationImage(
                    image: NetworkImage(widget.imageUrl!),
                    fit: BoxFit.cover,
                  )
                  : null,
        ),
        child:
            (_imageUrl != null && _imageUrl!.isNotEmpty) ||
                    (widget.imageUrl != null && widget.imageUrl!.isNotEmpty)
                ? null
                : (_isUploading
                    ? const CircularProgressIndicator()
                    : const Icon(Icons.person, size: 40, color: Colors.black)),
      ),
    );
  }
}
