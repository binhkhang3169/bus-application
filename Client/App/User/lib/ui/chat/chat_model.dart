class Chat {
  final String userName;
  String lastMessage;
  String time;
  int unreadCount;
  List<Message> messages;

  Chat({
    required this.userName,
    required this.lastMessage,
    required this.time,
    required this.unreadCount,
    required this.messages,
  });
}

class Message {
  final String text;
  final bool isSentByMe;
  final String time;

  Message({required this.text, required this.isSentByMe, required this.time});
}
