import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'chat_model.dart';

class ChatDetailPage extends StatefulWidget {
  final String userName;
  final List<Message> messages;
  final Function(String, List<Message>) onMessageSent;

  const ChatDetailPage({
    Key? key,
    required this.userName,
    required this.messages,
    required this.onMessageSent,
  }) : super(key: key);

  @override
  _ChatDetailPageState createState() => _ChatDetailPageState();
}

class _ChatDetailPageState extends State<ChatDetailPage> {
  final TextEditingController _messageController = TextEditingController();
  late List<Message> messages;
  bool _isTyping = false;
  bool _showEmojiPicker = false;
  bool _isRecording = false;
  final ImagePicker _picker = ImagePicker();

  final List<String> emojis = [
    "😊",
    "😂",
    "😍",
    "😢",
    "😡",
    "👍",
    "👏",
    "🙌",
    "🎉",
    "💖",
  ];

  @override
  void initState() {
    super.initState();
    messages = List.from(widget.messages);
    _messageController.addListener(() {
      setState(() {
        _isTyping = _messageController.text.trim().isNotEmpty;
      });
    });
  }

  void _sendMessage() {
    if (_messageController.text.trim().isEmpty) return;

    setState(() {
      messages.add(
        Message(
          text: _messageController.text.trim(),
          isSentByMe: true,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent(_messageController.text.trim(), messages);
      messages.add(
        Message(
          text: "Đã nhận tin nhắn: ${_messageController.text.trim()}",
          isSentByMe: false,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent(
        "Đã nhận tin nhắn: ${_messageController.text.trim()}",
        messages,
      );
      _messageController.clear();
    });
  }

  void _sendLike() {
    setState(() {
      messages.add(
        Message(
          text: "👍",
          isSentByMe: true,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent("👍", messages);
      messages.add(
        Message(
          text: "Đã nhận tin nhắn: 👍",
          isSentByMe: false,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent("Đã nhận tin nhắn: 👍", messages);
    });
  }

  void _sendEmoji(String emoji) {
    setState(() {
      messages.add(
        Message(
          text: emoji,
          isSentByMe: true,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent(emoji, messages);
      messages.add(
        Message(
          text: "Đã nhận tin nhắn: $emoji",
          isSentByMe: false,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent("Đã nhận tin nhắn: $emoji", messages);
      _showEmojiPicker = false;
    });
  }

  void _pickImage() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      setState(() {
        messages.add(
          Message(
            text: "[Hình ảnh]",
            isSentByMe: true,
            time: TimeOfDay.now().format(context),
          ),
        );
        widget.onMessageSent("[Hình ảnh]", messages);
        messages.add(
          Message(
            text: "Đã nhận tin nhắn: [Hình ảnh]",
            isSentByMe: false,
            time: TimeOfDay.now().format(context),
          ),
        );
        widget.onMessageSent("Đã nhận tin nhắn: [Hình ảnh]", messages);
      });
    }
  }

  void _startRecording() {
    setState(() {
      _isRecording = true;
    });
  }

  void _stopRecording() {
    setState(() {
      _isRecording = false;
      messages.add(
        Message(
          text: "[Ghi âm]",
          isSentByMe: true,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent("[Ghi âm]", messages);
      messages.add(
        Message(
          text: "Đã nhận tin nhắn: [Ghi âm]",
          isSentByMe: false,
          time: TimeOfDay.now().format(context),
        ),
      );
      widget.onMessageSent("Đã nhận tin nhắn: [Ghi âm]", messages);
    });
  }

  @override
  Widget build(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    final screenHeight = MediaQuery.of(context).size.height;

    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.white,

        leading: IconButton(
          icon: Icon(
            Icons.arrow_back,
            color: Colors.black,
            size: screenWidth * 0.06,
          ),
          onPressed: () => Navigator.pop(context),
        ),
        title: Row(
          children: [
            CircleAvatar(
              radius: screenWidth * 0.05,
              child: Text(
                widget.userName[0],
                style: TextStyle(fontSize: screenWidth * 0.04),
              ),
            ),
            SizedBox(width: screenWidth * 0.03),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.userName,
                  style: TextStyle(
                    color: Colors.black,
                    fontSize: screenWidth * 0.04,
                  ),
                ),
                Row(
                  children: [
                    Container(
                      width: screenWidth * 0.025,
                      height: screenWidth * 0.025,
                      decoration: BoxDecoration(
                        color: Colors.green,
                        shape: BoxShape.circle,
                      ),
                    ),
                    SizedBox(width: screenWidth * 0.015),
                    Text(
                      "Đang hoạt động",
                      style: TextStyle(
                        color: Colors.grey,
                        fontSize: screenWidth * 0.03,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: Icon(
              Icons.phone,
              color: Colors.blue,
              size: screenWidth * 0.06,
            ),
            onPressed: () {},
          ),
          IconButton(
            icon: Icon(
              Icons.videocam,
              color: Colors.blue,
              size: screenWidth * 0.06,
            ),
            onPressed: () {},
          ),
        ],
      ),
      body: Container(
        color: const Color.fromARGB(255, 248, 248, 248),
        child: Column(
          children: [
            Expanded(
              child: ListView.builder(
                padding: EdgeInsets.all(screenWidth * 0.03),
                itemCount: messages.length,
                itemBuilder: (context, index) {
                  final message = messages[index];
                  return Align(
                    alignment:
                        message.isSentByMe
                            ? Alignment.centerRight
                            : Alignment.centerLeft,
                    child: Container(
                      margin: EdgeInsets.symmetric(
                        vertical: screenHeight * 0.005,
                      ),
                      padding: EdgeInsets.symmetric(
                        horizontal: screenWidth * 0.04,
                        vertical: screenHeight * 0.01,
                      ),
                      decoration: BoxDecoration(
                        color:
                            message.isSentByMe
                                ? Colors.blue[100]
                                : Colors.grey[200],
                        borderRadius: BorderRadius.circular(screenWidth * 0.04),
                      ),
                      child: Column(
                        crossAxisAlignment:
                            message.isSentByMe
                                ? CrossAxisAlignment.end
                                : CrossAxisAlignment.start,
                        children: [
                          Text(
                            message.text,
                            style: TextStyle(fontSize: screenWidth * 0.04),
                          ),
                          SizedBox(height: screenHeight * 0.005),
                          Text(
                            message.time,
                            style: TextStyle(
                              fontSize: screenWidth * 0.025,
                              color: Colors.grey,
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                },
              ),
            ),
            if (_showEmojiPicker)
              Container(
                height: screenHeight * 0.2,
                padding: EdgeInsets.all(screenWidth * 0.02),
                child: GridView.builder(
                  gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 5,
                    crossAxisSpacing: screenWidth * 0.02,
                    mainAxisSpacing: screenHeight * 0.01,
                  ),
                  itemCount: emojis.length,
                  itemBuilder: (context, index) {
                    return GestureDetector(
                      onTap: () => _sendEmoji(emojis[index]),
                      child: Text(
                        emojis[index],
                        style: TextStyle(fontSize: screenWidth * 0.06),
                      ),
                    );
                  },
                ),
              ),
            Padding(
              padding: EdgeInsets.all(screenWidth * 0.02),
              child: Row(
                children: [
                  IconButton(
                    icon: Icon(
                      Icons.camera_alt,
                      color: Colors.blue,
                      size: screenWidth * 0.06,
                    ),
                    onPressed: () {},
                  ),
                  IconButton(
                    icon: Icon(
                      Icons.image,
                      color: Colors.blue,
                      size: screenWidth * 0.06,
                    ),
                    onPressed: _pickImage,
                  ),
                  IconButton(
                    icon: Icon(
                      _isRecording ? Icons.stop : Icons.mic,
                      color: Colors.blue,
                      size: screenWidth * 0.06,
                    ),
                    onPressed: _isRecording ? _stopRecording : _startRecording,
                  ),
                  Expanded(
                    child:
                        _isRecording
                            ? Container(
                              padding: EdgeInsets.symmetric(
                                horizontal: screenWidth * 0.03,
                                vertical: screenHeight * 0.01,
                              ),
                              decoration: BoxDecoration(
                                color: Colors.grey[200],
                                borderRadius: BorderRadius.circular(
                                  screenWidth * 0.05,
                                ),
                              ),
                              child: Row(
                                children: [
                                  Icon(Icons.mic, color: Colors.red),
                                  SizedBox(width: screenWidth * 0.02),
                                  Text(
                                    "Đang ghi âm...",
                                    style: TextStyle(
                                      fontSize: screenWidth * 0.035,
                                    ),
                                  ),
                                ],
                              ),
                            )
                            : TextField(
                              controller: _messageController,
                              decoration: InputDecoration(
                                hintText: "Aa",
                                hintStyle: TextStyle(
                                  fontSize: screenWidth * 0.035,
                                ),
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(
                                    screenWidth * 0.05,
                                  ),
                                  borderSide: BorderSide(color: Colors.grey),
                                ),
                                focusedBorder: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(
                                    screenWidth * 0.05,
                                  ),
                                  borderSide: BorderSide(
                                    color: Colors.blue,
                                    width: 2,
                                  ),
                                ),
                                enabledBorder: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(
                                    screenWidth * 0.05,
                                  ),
                                  borderSide: BorderSide(color: Colors.grey),
                                ),
                                contentPadding: EdgeInsets.symmetric(
                                  horizontal: screenWidth * 0.03,
                                  vertical: screenHeight * 0.01,
                                ),
                                suffixIcon: IconButton(
                                  icon: Icon(
                                    Icons.emoji_emotions,
                                    color: Colors.blue,
                                    size: screenWidth * 0.06,
                                  ),
                                  onPressed: () {
                                    setState(() {
                                      _showEmojiPicker = !_showEmojiPicker;
                                    });
                                  },
                                ),
                              ),
                              style: TextStyle(fontSize: screenWidth * 0.035),
                            ),
                  ),
                  IconButton(
                    icon: Icon(
                      _isTyping ? Icons.send : Icons.thumb_up,
                      color: Colors.blue,
                      size: screenWidth * 0.06,
                    ),
                    onPressed: _isTyping ? _sendMessage : _sendLike,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _messageController.dispose();
    super.dispose();
  }
}
