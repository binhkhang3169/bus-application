/** @format */

import React, { useState, useEffect, useRef } from "react";
import "./AiChatbox.css"; // Giữ nguyên file CSS của bạn
import { API_URL } from '../configs/env';
// Chú ý: Đảm bảo rằng API_URL đã được định nghĩa trong file env.js
// --- Hằng số cho API và Local Storage ---
// thay apuurl vào dưới
const RASA_API_URL = `${API_URL}webhooks/rest/webhook`; // Đường dẫn API của RASA
const LOCAL_STORAGE_KEY = "ai_chat_messages";
const SENDER_ID = "user1"; // ID người dùng, có thể thay đổi nếu cần

const AiChatbox = () => {
  // State chính: Quản lý việc chat đang mở hay đóng
  const [isChatOpen, setIsChatOpen] = useState(false);

  // State quản lý tin nhắn, khởi tạo từ localStorage hoặc tin nhắn mặc định
  const [messages, setMessages] = useState(() => {
    try {
      const savedMessages = localStorage.getItem(LOCAL_STORAGE_KEY);
      // Nếu có tin nhắn đã lưu, parse nó, ngược lại trả về mảng chào mừng
      return savedMessages
        ? JSON.parse(savedMessages)
        : [{ id: 1, text: "Xin chào! Tôi có thể giúp gì cho bạn?", sender: "ai" }];
    } catch (error) {
      console.error("Failed to parse messages from localStorage", error);
      return [{ id: 1, text: "Xin chào! Tôi có thể giúp gì cho bạn?", sender: "ai" }];
    }
  });

  const [inputValue, setInputValue] = useState("");
  const [isTyping, setIsTyping] = useState(false); // State cho chỉ báo "đang nhập"
  const chatWindowRef = useRef(null);

  // Effect #1: Tự động cuộn xuống tin nhắn mới nhất
  useEffect(() => {
    if (isChatOpen && chatWindowRef.current) {
      chatWindowRef.current.scrollTop = chatWindowRef.current.scrollHeight;
    }
  }, [messages, isChatOpen, isTyping]);

  // Effect #2: Lưu tin nhắn vào localStorage mỗi khi có thay đổi
  useEffect(() => {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(messages));
  }, [messages]);

  // Hàm xử lý gửi tin nhắn (đã tích hợp API)
  const handleSendMessage = async (e) => {
    // Chỉ hoạt động khi nhấn Enter và có nội dung
    if (e.key !== "Enter" || inputValue.trim() === "") return;

    const userMessageText = inputValue.trim();

    // 1. Thêm tin nhắn của người dùng vào giao diện ngay lập tức
    const newUserMessage = {
      id: Date.now(),
      text: userMessageText,
      sender: "user",
    };
    setMessages((prevMessages) => [...prevMessages, newUserMessage]);
    setInputValue(""); // Xóa input
    setIsTyping(true); // Hiển thị "Bot đang nhập..."

    // 2. Gửi yêu cầu đến API
    try {
      const response = await fetch(RASA_API_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sender: SENDER_ID,
          message: userMessageText,
        }),
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.statusText}`);
      }

      const aiData = await response.json();

      // 3. Xử lý phản hồi từ AI
      if (aiData && aiData.length > 0) {
        const aiResponses = aiData.map((msg, index) => ({
          id: Date.now() + index + 1, // Đảm bảo id là duy nhất
          text: msg.text,
          sender: "ai",
        }));
        setMessages((prevMessages) => [...prevMessages, ...aiResponses]);
      } else {
         // Trường hợp API trả về mảng rỗng
         const fallbackResponse = {
            id: Date.now() + 1,
            text: "Xin lỗi, tôi chưa hiểu ý bạn. Bạn có thể diễn đạt khác được không?",
            sender: "ai",
        };
        setMessages((prevMessages) => [...prevMessages, fallbackResponse]);
      }

    } catch (error) {
      console.error("Failed to send message:", error);
      // Hiển thị tin nhắn lỗi trên giao diện
      const errorMessage = {
        id: Date.now() + 1,
        text: "Đã có lỗi xảy ra. Vui lòng thử lại sau.",
        sender: "ai",
      };
      setMessages((prevMessages) => [...prevMessages, errorMessage]);
    } finally {
      setIsTyping(false); // Ẩn "Bot đang nhập..."
    }
  };

  // --- PHẦN RENDERING DỰA TRÊN STATE ---

  // TRẠNG THÁI BAN ĐẦU: Hiển thị nút chat đơn giản
  if (!isChatOpen) {
    return (
      <button
        onClick={() => setIsChatOpen(true)}
        className="chat-opener-button"
        aria-label="Open chat"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-8 w-8"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
          />
        </svg>
      </button>
    );
  }

  // TRẠNG THÁI KHI CHAT MỞ: Hiển thị giao diện chat đầy đủ
  return (
    <div className="chat-container active">
      {/* Cửa sổ hiển thị tin nhắn */}
      <div ref={chatWindowRef} className="chat-window">
        {messages.map((msg) => (
          <div key={msg.id} className={`chat-message ${msg.sender}`}>
            <p>{msg.text}</p>
          </div>
        ))}
        {/* Chỉ báo "Bot đang nhập..." */}
        {isTyping && (
          <div className="chat-message ai">
            <p>
              <i>Bot đang nhập...</i>
            </p>
          </div>
        )}
      </div>

      {/* Thanh Chatbar */}
      <div className="chat-bar">
        <div className="chat-bar__input-wrapper">
          <input
            className="chat-bar__input"
            type="text"
            placeholder="Nhắn tin..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleSendMessage}
            autoFocus
            disabled={isTyping} // Vô hiệu hóa input khi bot đang trả lời
          />
        </div>

        {/* Nút đóng chat */}
        <div className="chat-bar__close" onClick={() => setIsChatOpen(false)}>
          <i className="material-icons">close</i>
        </div>
      </div>
    </div>
  );
};

export default AiChatbox;