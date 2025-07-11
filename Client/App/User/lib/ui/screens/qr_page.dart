import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart'; 
import 'package:image_picker/image_picker.dart';

class QRPage extends StatefulWidget {
  @override
  _QRPageState createState() => _QRPageState();
}

class _QRPageState extends State<QRPage> {
  final MobileScannerController _scannerController = MobileScannerController();
  final ImagePicker _picker = ImagePicker();
  bool isFlashOn = false;

  void _onQRDetected(BarcodeCapture barcodeCapture) {
    final String code = barcodeCapture.barcodes.first.rawValue ?? ''; 
    if (code.isNotEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text("Mã QR: $code")),
      );
    }
  }

  Future<void> _pickImage() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      print("Đã chọn ảnh: ${image.path}");
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text("Đã chọn ảnh để quét QR (giả lập)")),
      );
    }
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: PreferredSize(
        preferredSize: Size.fromHeight(70), 
        child: Container(
          height: 70, 
          decoration: BoxDecoration(
            image: DecorationImage(
              image: AssetImage('assets/images/background1.jpg'),
              fit: BoxFit.cover,
            ),
          ),
          child: AppBar(
            backgroundColor: Colors.transparent,
            elevation: 0,
            leading: IconButton(
              icon: Icon(Icons.arrow_back, color: Colors.white),
              onPressed: () {
                Navigator.pop(context);
              },
            ),
            title: Text(
              'Quét mã QR',
              style: TextStyle(color: Colors.white, fontSize: 20),
            ),
            centerTitle: true,
          ),
        ),
      ),
      body: Container(
        color: Colors.blue, 
        child: Padding(
          padding: EdgeInsets.all(16.0), 
          child: Stack(
            children: [
              MobileScanner(
                controller: _scannerController,
                onDetect: _onQRDetected, 
              ),
              Positioned(
                top: 50,
                left: 0,
                right: 0,
                child: Center(
                  child: Text(
                    "Di chuyển camera đến vùng chứa mã QR",
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.w400,
                    ),
                  ),
                ),
              ),
              Positioned(
                bottom: 50,
                left: 0,
                right: 0,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Column(
                      children: [
                        IconButton(
                          icon: Icon(
                            isFlashOn ? Icons.flash_on : Icons.flash_off,
                            color: Colors.white,
                            size: 30,
                          ),
                          onPressed: () {
                            setState(() {
                              isFlashOn = !isFlashOn;
                            });
                            _scannerController.toggleTorch(); 
                          },
                        ),
                        Text(
                          "Bật đèn pin",
                          style: TextStyle(color: Colors.white, fontSize: 14),
                        ),
                      ],
                    ),
                    SizedBox(width: 50),
                    Column(
                      children: [
                        IconButton(
                          icon: Icon(Icons.photo_library, color: Colors.white, size: 30),
                          onPressed: _pickImage,
                        ),
                        Text(
                          "Thư viện ảnh",
                          style: TextStyle(color: Colors.white, fontSize: 14),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
