/** @format */

import React, { useState, useEffect } from "react";
import { API_URL } from "../configs/env";
import SeatQrList from "./SeatQrList";
import ImageLightbox from "./ImageLightbox";

// --- CÁC HÀM HELPER ---

// Helper lấy giá trị chuỗi an toàn từ object lồng nhau của Go
const getString = (obj, defaultValue = "N/A") => {
  if (obj && obj.Valid && typeof obj.String === "string" && obj.String !== "") {
    return obj.String;
  }
  return defaultValue;
};

// Helper lấy giá trị số nguyên an toàn từ object lồng nhau của Go
const getInt = (obj, defaultValue = "N/A") => {
  if (obj && obj.Valid && typeof obj.Int32 === "number") {
    return obj.Int32;
  }
  return defaultValue;
};

// Helper hiển thị trạng thái thanh toán
const getPaymentStatusText = (status) => {
  if (status === 1)
    return (
      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
        Đã thanh toán
      </span>
    );
  if (status === 0)
    return (
      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800">
        Chưa thanh toán
      </span>
    );
  return (
    <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800">
      Không xác định
    </span>
  );
};

// Helper hiển thị trạng thái vé
const getTicketStatusText = (status, paymentStatus) => {
  if (status === 2)
    return (
      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
        Đã hủy
      </span>
    );
  if (status === 1 && paymentStatus === 1)
    return (
      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
        Hoàn thành
      </span>
    );
  if (status === 0 && paymentStatus === 0)
    return (
      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800">
        Chờ xử lý
      </span>
    );
  return (
    <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
      Đang hoạt động
    </span>
  );
};

// Helper hiển thị kênh đặt vé
const getBookingChannelText = (channel) => {
  switch (channel) {
    case 0:
      return "Trang web";
    case 1:
      return "App";
    case 2:
      return "Tại quầy";
    default:
      return "N/A";
  }
};

// Helper để fetch QR Code song song cho một danh sách ghế
const generateQrCodesForSeats = async (seats, ticketId) => {
  if (!seats || seats.length === 0) {
    return { urls: {}, errors: {} };
  }
  const qrPromises = seats.map(async (seat) => {
    if (!seat || !seat.seat_id) return null;
    const content = `TICKET:${ticketId}-SEAT:${seat.seat_id}`;
    const qrApiUrl = `${API_URL}api/v1/qr/image?content=${encodeURIComponent(
      content
    )}`;
    try {
      const response = await fetch(qrApiUrl);
      if (!response.ok) throw new Error("QR fetch failed");
      const blob = await response.blob();
      return {
        seatId: seat.seat_id,
        url: URL.createObjectURL(blob),
        error: false,
      };
    } catch (error) {
      console.error(`Error fetching QR for seat ${seat.seat_id}:`, error);
      return { seatId: seat.seat_id, url: null, error: true };
    }
  });
  const results = await Promise.all(qrPromises);
  const urls = {},
    errors = {};
  results
    .filter((r) => r)
    .forEach((res) => {
      urls[res.seatId] = res.url;
      errors[res.seatId] = res.error;
    });
  return { urls, errors };
};

// --- COMPONENT CHÍNH ---

function TicketDetailModal({
  isOpen,
  onClose,
  ticket,
  formatPrice,
  formatDate,
}) {
  const [qrCodeUrls, setQrCodeUrls] = useState({});
  const [qrCodeErrors, setQrCodeErrors] = useState({});
  const [locationNames, setLocationNames] = useState({});
  const [isLoading, setIsLoading] = useState(false);
  const [zoomedImageUrl, setZoomedImageUrl] = useState(null);

  useEffect(() => {
    if (!isOpen || !ticket) {
      setQrCodeUrls({});
      setQrCodeErrors({});
      setLocationNames({});
      setIsLoading(false);
      return;
    }

    const fetchAllData = async () => {
      setIsLoading(true);
      const [locationData, qrData] = await Promise.all([
        (async () => {
          try {
            const response = await fetch(`${API_URL}api/v1/stations/status/1`);
            if (!response.ok) throw new Error("Failed to fetch stations");
            const result = await response.json();
            if (result && result.data) {
              return result.data.reduce((acc, station) => {
                acc[station.id] = station.name;
                return acc;
              }, {});
            }
          } catch (error) {
            console.error("Error fetching station names:", error);
          }
          return {};
        })(),
        (async () => {
          const [beginResult, endResult] = await Promise.all([
            generateQrCodesForSeats(ticket.SeatTicketsBegin, ticket.ticket_id),
            generateQrCodesForSeats(ticket.SeatTicketsEnd, ticket.ticket_id),
          ]);
          return {
            urls: { ...beginResult.urls, ...endResult.urls },
            errors: { ...beginResult.errors, ...endResult.errors },
          };
        })(),
      ]);

      setLocationNames(locationData);
      setQrCodeUrls(qrData.urls);
      setQrCodeErrors(qrData.errors);
      setIsLoading(false);
    };

    fetchAllData();

    return () => {
      Object.values(qrCodeUrls).forEach((url) => {
        if (url) URL.revokeObjectURL(url);
      });
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, ticket]);

  const handleImageZoom = (url) => {
    setZoomedImageUrl(url);
  };

  const handleCloseZoom = () => {
    setZoomedImageUrl(null);
  };

  if (!isOpen || !ticket) {
    return null;
  }

  const detail =
    ticket.Details && ticket.Details.length > 0 ? ticket.Details[0] : {};
  const pickupIdBegin = getInt(detail.pickup_location_begin);
  const dropoffIdBegin = getInt(detail.dropoff_location_begin);
  const pickupIdEnd = getInt(detail.pickup_location_end);
  const dropoffIdEnd = getInt(detail.dropoff_location_end);

  return (
    <>
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex justify-center items-center p-4">
        <div className="relative mx-auto p-6 border w-full max-w-4xl shadow-lg rounded-lg bg-gray-50">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-2xl font-semibold text-gray-900">
              Chi Tiết Vé: {ticket.ticket_id}
            </h3>
            <button
              onClick={onClose}
              className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center"
            >
              <svg
                className="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-3 gap-x-6 gap-y-4 text-sm mb-6 p-4 bg-white rounded-lg shadow">
            <div>
              <strong>Mã vé:</strong> {ticket.ticket_id}
            </div>
            <div>
              <strong>Tên người đặt:</strong> {getString(ticket.name)}
            </div>
            <div>
              <strong>Số điện thoại:</strong> {getString(ticket.phone)}
            </div>
            <div>
              <strong>Email:</strong> {getString(ticket.email)}
            </div>
            <div>
              <strong>Giá vé:</strong> {formatPrice(ticket.price)}
            </div>
            <div>
              <strong>Thời gian đặt:</strong> {formatDate(ticket.booking_time)}
            </div>
            <div>
              <strong>Trạng thái vé:</strong>{" "}
              {getTicketStatusText(ticket.status, ticket.payment_status)}
            </div>
            <div>
              <strong>Trạng thái thanh toán:</strong>{" "}
              {getPaymentStatusText(ticket.payment_status)}
            </div>
            <div>
              <strong>Kênh đặt vé:</strong>{" "}
              {getBookingChannelText(ticket.booking_channel)}
            </div>
          </div>

          {/* <div className="p-4 bg-white rounded-lg shadow mb-6">
            <h4 className="text-lg font-semibold mb-3 text-gray-800">
              Chi tiết hành trình
            </h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {(pickupIdBegin !== "N/A" || dropoffIdBegin !== "N/A") && (
                <div className="p-3 bg-gray-50 rounded-md border border-gray-200">
                  <p className="font-semibold text-indigo-700">
                    Lượt đi (Mã chuyến: {ticket.trip_id_begin || "N/A"})
                  </p>
                  <p>
                    <strong>Điểm đón:</strong>{" "}
                    {isLoading
                      ? "Đang tải..."
                      : locationNames[pickupIdBegin] || `ID: ${pickupIdBegin}`}
                  </p>
                  <p>
                    <strong>Điểm trả:</strong>{" "}
                    {isLoading
                      ? "Đang tải..."
                      : locationNames[dropoffIdBegin] ||
                        `ID: ${dropoffIdBegin}`}
                  </p>
                </div>
              )}
              {(pickupIdEnd !== "N/A" || dropoffIdEnd !== "N/A") && (
                <div className="p-3 bg-gray-50 rounded-md border border-gray-200">
                  <p className="font-semibold text-green-700">
                    Lượt về (Mã chuyến: {getString(ticket.trip_id_end, "Không có")})
                  </p>
                  <p>
                    <strong>Điểm đón:</strong>{" "}
                    {isLoading
                      ? "Đang tải..."
                      : locationNames[pickupIdEnd] || `ID: ${pickupIdEnd}`}
                  </p>
                  <p>
                    <strong>Điểm trả:</strong>{" "}
                    {isLoading
                      ? "Đang tải..."
                      : locationNames[dropoffIdEnd] || `ID: ${dropoffIdEnd}`}
                  </p>
                </div>
              )}
            </div>
          </div> */}

          <SeatQrList
            title="Vé Lượt Đi"
            seats={ticket.SeatTicketsBegin}
            qrCodeUrls={qrCodeUrls}
            qrCodeErrors={qrCodeErrors}
            onImageClick={handleImageZoom}
          />
          <SeatQrList
            title="Vé Lượt Về"
            seats={ticket.SeatTicketsEnd}
            qrCodeUrls={qrCodeUrls}
            qrCodeErrors={qrCodeErrors}
            onImageClick={handleImageZoom}
          />

          <div className="mt-8 text-right">
            <button
              onClick={onClose}
              className="px-6 py-2 bg-gray-500 text-white text-base font-medium rounded-md shadow-sm hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-400"
            >
              Đóng
            </button>
          </div>
        </div>
      </div>

      <ImageLightbox imageUrl={zoomedImageUrl} onClose={handleCloseZoom} />
    </>
  );
}

export default TicketDetailModal;
