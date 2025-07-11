import 'package:flutter/material.dart';
import 'chat_model.dart';
import 'chat_detail_page.dart';

class ChatListPage extends StatefulWidget {
  @override
  _ChatListPageState createState() => _ChatListPageState();
}

class _ChatListPageState extends State<ChatListPage> {
  final TextEditingController _searchController = TextEditingController();
  List<Chat> chats = [
    Chat(
      userName: "Nguyễn Cao Kỳ",
      lastMessage: "Chào! Có gì không?",
      time: "16:02",
      unreadCount: 0,
      messages: [
        Message(text: "Chào bạn!", isSentByMe: false, time: "10:00"),
        Message(text: "Chào! Có gì không?", isSentByMe: true, time: "10:01"),
      ],
    ),
    Chat(
      userName: "Nguyễn Minh Luân",
      lastMessage: "Hỏi Đáp: Ngoài lề một chút ạ",
      time: "14:27",
      unreadCount: 0,
      messages: [
        Message(
          text: "Hỏi Đáp: Ngoài lề một chút ạ",
          isSentByMe: false,
          time: "14:27",
        ),
      ],
    ),
  ];

  List<Chat> filteredChats = [];

  @override
  void initState() {
    super.initState();
    filteredChats = chats;
    _searchController.addListener(_filterChats);
  }

  void _filterChats() {
    setState(() {
      if (_searchController.text.isEmpty) {
        filteredChats = chats;
      } else {
        filteredChats =
            chats
                .where(
                  (chat) => chat.userName.toLowerCase().contains(
                    _searchController.text.toLowerCase(),
                  ),
                )
                .toList();
      }
    });
  }

  void _updateChat(
    String userName,
    String newMessage,
    List<Message> updatedMessages,
  ) {
    setState(() {
      final chatIndex = chats.indexWhere((chat) => chat.userName == userName);
      if (chatIndex != -1) {
        chats[chatIndex] = Chat(
          userName: chats[chatIndex].userName,
          lastMessage: newMessage,
          time: TimeOfDay.now().format(context),
          unreadCount: chats[chatIndex].unreadCount,
          messages: updatedMessages,
        );
        _filterChats();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: Text("Đoạn chat", style: TextStyle(color: Colors.white)),
        centerTitle: true,
        backgroundColor: Colors.blue,
        foregroundColor: Colors.white,
      ),
      body: Column(
        children: [
          Padding(
            padding: EdgeInsets.symmetric(horizontal: 10, vertical: 10),
            child: SizedBox(
              height: 35,
              child: TextField(
                controller: _searchController,
                style: TextStyle(fontSize: 12),
                decoration: InputDecoration(
                  hintText: "Tìm kiếm",
                  hintStyle: TextStyle(fontSize: 12),
                  prefixIcon: Icon(Icons.search, size: 18),
                  contentPadding: EdgeInsets.symmetric(
                    vertical: 0,
                    horizontal: 10,
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(25),
                    borderSide: BorderSide(color: Colors.grey),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(25),
                    borderSide: BorderSide(color: Colors.blue, width: 2),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(25),
                    borderSide: BorderSide(color: Colors.grey),
                  ),
                  filled: true,
                  fillColor: Colors.grey[200],
                ),
              ),
            ),
          ),
          Expanded(
            child: ListView.builder(
              itemCount: filteredChats.length,
              itemBuilder: (context, index) {
                final chat = filteredChats[index];
                return ListTile(
                  leading: CircleAvatar(child: Text(chat.userName[0])),
                  title: Text(chat.userName),
                  subtitle: Text(chat.lastMessage),
                  trailing: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(chat.time, style: TextStyle(color: Colors.grey)),
                      if (chat.unreadCount > 0)
                        Container(
                          margin: EdgeInsets.only(top: 5),
                          padding: EdgeInsets.all(5),
                          decoration: BoxDecoration(
                            color: Colors.red,
                            shape: BoxShape.circle,
                          ),
                          child: Text(
                            chat.unreadCount.toString(),
                            style: TextStyle(color: Colors.white, fontSize: 12),
                          ),
                        ),
                    ],
                  ),
                  onTap: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder:
                            (context) => ChatDetailPage(
                              userName: chat.userName,
                              messages: chat.messages,
                              onMessageSent: (newMessage, updatedMessages) {
                                _updateChat(
                                  chat.userName,
                                  newMessage,
                                  updatedMessages,
                                );
                              },
                            ),
                      ),
                    );
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }
}
