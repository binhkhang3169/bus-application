/** @format */

import React, { useState, useEffect } from "react";
import axios from "axios";
import { API_URL } from "../configs/env";

// --- HELPER FUNCTIONS ---

// Định dạng giá tiền
const formatPrice = (price) => {
  if (typeof price !== "number") return "N/A";
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
  }).format(price);
};

// Định dạng ngày giờ từ chuỗi ISO
const formatISODateTime = (isoString) => {
  if (!isoString) return "N/A";
  const options = {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  };
  return new Date(isoString).toLocaleString("vi-VN", options);
};

// Kết hợp ngày và giờ riêng lẻ và định dạng
const formatTripDateTime = (dateStr, timeStr) => {
  if (!dateStr || !timeStr) return "N/A";
  const isoString = `${dateStr}T${timeStr}`;
  return formatISODateTime(isoString);
};

// Diễn giải trạng thái (vé & thanh toán)
const getStatusBadge = (type, status) => {
  const statuses = {
    ticket: {
      1: <span className="font-semibold text-green-600">Đã xác nhận</span>,
      0: <span className="font-semibold text-yellow-600">Chờ xác nhận</span>,
    },
    payment: {
      1: <span className="font-semibold text-green-600">Đã thanh toán</span>,
      0: <span className="font-semibold text-red-600">Chưa thanh toán</span>,
    },
  };
  return (
    statuses[type]?.[status] || (
      <span className="font-semibold text-gray-600">Không xác định</span>
    )
  );
};

// --- TRIP DETAIL COMPONENT ---

// Component con để hiển thị thông tin chi tiết cho một chuyến đi
const TripDetails = ({ title, tripData }) => {
  if (!tripData) return null;

  const {
    route,
    departureDate,
    departureTime,
    arrivalDate,
    arrivalTime,
    vehicle,
    stock,
  } = tripData;

  return (
    <div className="mt-4 p-4 border rounded-lg bg-gray-50">
      <h5 className="text-md font-bold text-gray-700 mb-3">{title}</h5>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
        <p>
          <strong>Hành trình:</strong> {route?.start?.name} → {route?.end?.name}
        </p>
        <p>
          <strong>Khoảng cách:</strong> {route?.distance} (
          {route?.estimatedTime})
        </p>
        <p>
          <strong>Khởi hành:</strong>{" "}
          {formatTripDateTime(departureDate, departureTime)}
        </p>
        <p>
          <strong>Dự kiến đến:</strong>{" "}
          {formatTripDateTime(arrivalDate, arrivalTime)}
        </p>
        <p>
          <strong>Phương tiện:</strong> {vehicle?.type?.name} (
          {vehicle?.license})
        </p>
        <p>
          <strong>Số ghế trống:</strong> {stock}
        </p>
      </div>
    </div>
  );
};

// --- MAIN COMPONENT ---

function LookUpResult({ ticket }) {
  // State để lưu thông tin chi tiết của chuyến đi
  const [tripDetails, setTripDetails] = useState({ begin: null, end: null });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Hàm fetch dữ liệu cho một trip ID cụ thể
    const fetchTripById = async (tripId) => {
      if (!tripId) return null;
      // Giả sử API_URL đã được định nghĩa đúng
      const response = await axios.get(`${API_URL}api/v1/trips/${tripId}`);
      return response.data.data;
    };

    const fetchAllTripDetails = async () => {
      if (!ticket?.ticket_id) return;

      setIsLoading(true);
      setError(null);
      setTripDetails({ begin: null, end: null });

      try {
        // Tạo các promise để gọi API
        const beginTripPromise = fetchTripById(ticket.trip_id_begin);

        // *** ĐIỀU CHỈNH LOGIC TẠI ĐÂY ***
        // Chỉ cần trip_id_end hợp lệ là sẽ fetch, không cần kiểm tra 'type'
        const endTripPromise = ticket.trip_id_end?.Valid
          ? fetchTripById(ticket.trip_id_end.String)
          : Promise.resolve(null);

        // Chờ cả hai promise hoàn thành
        const [beginTripData, endTripData] = await Promise.all([
          beginTripPromise,
          endTripPromise,
        ]);

        setTripDetails({ begin: beginTripData, end: endTripData });
      } catch (err) {
        console.error("Failed to fetch trip details:", err);
        setError(
          "Không thể tải thông tin chi tiết hành trình. Vui lòng thử lại."
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchAllTripDetails();
  }, [ticket]);

  if (!ticket) return null;

  const {
    ticket_id,
    name,
    phone,
    email,
    price,
    booking_time,
    status,
    payment_status,
  } = ticket;

  // Xác định loại vé dựa trên dữ liệu thực tế
  const ticketTypeText = tripDetails.end ? "Khứ hồi" : "Một chiều";

  return (
    <div className="w-full border border-slate-300 rounded-xl flex flex-col my-10 shadow-lg">
      <div className="bg-gray-100 w-full text-center p-3 rounded-t-xl">
        <h3 className="text-xl font-bold text-gray-800">Kết Quả Tra Cứu Vé</h3>
      </div>

      {/* Phần thông tin vé & khách hàng */}
      <div className="p-6">
        <h4 className="text-lg font-semibold mb-4 text-blue-700">
          Thông tin chung
        </h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
          {/* Thông tin khách hàng */}
          <p>
            <strong>Họ và tên:</strong> {name?.String || "N/A"}
          </p>
          <p>
            <strong>Số điện thoại:</strong> {phone?.String || "N/A"}
          </p>
          <p>
            <strong>Email:</strong> {email?.String || "N/A"}
          </p>

          {/* Thông tin vé */}
          <p>
            <strong>Mã vé:</strong>{" "}
            <span className="font-mono bg-gray-200 px-2 py-1 rounded">
              {ticket_id}
            </span>
          </p>
          <p>
            <strong>Loại vé:</strong>{" "}
            {isLoading ? "Đang xác định..." : ticketTypeText}
          </p>
          <p>
            <strong>Tổng giá:</strong>{" "}
            <span className="font-bold text-lg text-red-500">
              {formatPrice(price)}
            </span>
          </p>
          <p>
            <strong>Thời gian đặt:</strong> {formatISODateTime(booking_time)}
          </p>
          <p>
            <strong>Trạng thái vé:</strong> {getStatusBadge("ticket", status)}
          </p>
          <p>
            <strong>Thanh toán:</strong>{" "}
            {getStatusBadge("payment", payment_status)}
          </p>
        </div>
      </div>

      <hr className="mx-6" />

      {/* Phần thông tin hành trình (lấy từ API trips) */}
      <div className="p-6">
        <h4 className="text-lg font-semibold text-blue-700">
          Thông tin hành trình
        </h4>
        {isLoading && (
          <p className="text-gray-500 mt-2">Đang tải chi tiết chuyến đi...</p>
        )}
        {error && <p className="text-red-500 mt-2">{error}</p>}

        {!isLoading && !error && (
          <>
            <TripDetails title="Chuyến đi" tripData={tripDetails.begin} />
            {/* Component TripDetails cho chuyến về sẽ chỉ render khi có dữ liệu */}
            <TripDetails title="Chuyến về" tripData={tripDetails.end} />
          </>
        )}
      </div>

      <div className="bg-green-100 rounded-b-xl p-4 text-center text-green-800 font-semibold mt-2">
        Vui lòng mang mã vé đến văn phòng để đổi vé lên xe trước giờ khởi hành
        ít nhất 60 phút.
      </div>
    </div>
  );
}

export default LookUpResult;
