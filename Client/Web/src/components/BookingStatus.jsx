/** @format */

import React, { useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";
import api from "../services/apiService";
import { API_URL } from "../configs/env";

const REAL_SERVER_WS_URL = "ws://172.188.164.9";

const BookingStatus = ({ bookingId, onReset }) => {
  const [status, setStatus] = useState("WAITING_RESULT");
  const [resultPayload, setResultPayload] = useState(null);
  const [errorInfo, setErrorInfo] = useState({ title: "", message: "" });

  const wsUrl = `${REAL_SERVER_WS_URL}/api/v1/ws/track/${bookingId}`;

  const { lastJsonMessage, sendJsonMessage } = useWebSocket(wsUrl, {
    onOpen: () => console.log("WebSocket connection established."),
    onClose: () => console.log("WebSocket connection closed."),
    onError: (error) => {
      console.error("WebSocket error:", error);
      setErrorInfo({
        title: "⚠️ Lỗi Kết Nối",
        message: "Không thể kết nối tới WebSocket server.",
      });
      setStatus("FINAL_STATE");
    },
    shouldReconnect: () => false,
  });

  const redirectToPayment = async (ticketData) => {
    const token =
      localStorage.getItem("accessToken") ||
      sessionStorage.getItem("accessToken");

    // Destructure các trường cần thiết từ payload
    const { ticket_id, price, customer_id, phone } = ticketData;

    if (!ticket_id || price === undefined) {
      setErrorInfo({
        title: "❌ Lỗi Dữ Liệu",
        message:
          "Phản hồi từ server không đủ thông tin để thanh toán (thiếu ticket_id hoặc price).",
      });
      setStatus("FINAL_STATE");
      return;
    }

    setStatus("REDIRECTING");

    try {
      // --- LOGIC MỚI: KIỂM TRA VÀ CHỌN customer_id HOẶC SĐT ---
      // Kiểm tra nếu customer_id có tồn tại, hợp lệ và lớn hơn 0.
      // Nếu không, sử dụng số điện thoại để thay thế.
      const customerIdentifier =
        customer_id && customer_id.Valid && customer_id.Int32 > 0
          ? customer_id.Int32
          : phone.String;

      const paymentGatewayPayload = {
        amount: price,
        bank_code: "",
        language: "vn",
        ticket_id: ticket_id,
        notes: `Thanh toán vé ${ticket_id}`,
        // Sử dụng giá trị đã được xác định ở trên
        customer_id: String(customerIdentifier || ""),
      };

      const response = await api.post(
        `${API_URL}api/v1/vnpay/create-payment`,
        paymentGatewayPayload,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      if (response.data?.data?.payment_url) {
        window.location.href = response.data.data.payment_url;
      } else {
        throw new Error("Không tạo được URL thanh toán từ server.");
      }
    } catch (error) {
      const errorMessage =
        error.response?.data?.message ||
        error.message ||
        "Lỗi trong quá trình tạo yêu cầu thanh toán.";
      setErrorInfo({ title: "❌ Lỗi Thanh Toán", message: errorMessage });
      setStatus("FINAL_STATE");
    }
  };

  useEffect(() => {
    if (lastJsonMessage) {
      console.log("Nhận được tin nhắn từ server:", lastJsonMessage);

      sendJsonMessage({ type: "ack" });

      if (
        lastJsonMessage.type === "result" &&
        lastJsonMessage.payload?.ticket_id
      ) {
        setResultPayload(lastJsonMessage.payload);
        redirectToPayment(lastJsonMessage.payload);
      } else {
        setErrorInfo({
          title: "❌ Đặt vé thất bại",
          message:
            lastJsonMessage.payload?.message ||
            "Server trả về lỗi không xác định.",
        });
        setResultPayload(lastJsonMessage.payload);
        setStatus("FINAL_STATE");
      }
    }
  }, [lastJsonMessage]);

  const renderContent = () => {
    if (status === "WAITING_RESULT") {
      return (
        <div className="mx-auto">
          <div className="spinner w-9 h-9 border-4 border-t-blue-600 border-solid rounded-full animate-spin mx-auto"></div>
          <p className="mt-2 font-semibold text-blue-600">
            Đang xử lý... Vui lòng chờ trong giây lát.
          </p>
          <p className="text-gray-500 text-sm">
            Booking ID: <span className="font-mono">{bookingId}</span>
          </p>
        </div>
      );
    }

    if (status === "REDIRECTING") {
      return (
        <div className="mx-auto">
          <div className="spinner w-9 h-9 border-4 border-t-green-600 border-solid rounded-full animate-spin mx-auto"></div>
          <p className="mt-2 font-semibold text-green-600">
            Đặt vé thành công! Đang chuyển hướng đến trang thanh toán...
          </p>
        </div>
      );
    }

    if (status === "FINAL_STATE") {
      return (
        <>
          <h3 className="text-2xl font-bold text-red-600">{errorInfo.title}</h3>
          <pre className="mt-2 text-left bg-gray-100 p-4 rounded-lg text-sm font-mono whitespace-pre-wrap">
            {JSON.stringify(
              resultPayload || { message: errorInfo.message },
              null,
              2
            )}
          </pre>
          <button
            onClick={onReset}
            className="mt-4 bg-gray-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-gray-600 transition duration-300"
          >
            Thực hiện lại
          </button>
        </>
      );
    }

    return null;
  };

  return (
    <div className="w-full bg-white rounded-xl shadow-lg p-8 space-y-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-800">Trạng Thái Đặt Vé</h1>
        <p className="text-gray-500 mt-2">
          Theo dõi kết quả từ server real-time
        </p>
      </div>
      <div className="text-center p-6 border-2 border-dashed rounded-lg space-y-4">
        {renderContent()}
      </div>
    </div>
  );
};

export default BookingStatus;
