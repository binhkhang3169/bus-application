/** @format */
import axios from "axios";
import React, { useState, useRef, useEffect } from "react";
import { API_URL } from "../../configs/env";

const TicketManagement = () => {
  const printRef = useRef();
  const [isLoading, setIsLoading] = useState(false);
  const [ticket, setTicket] = useState({});
  const [ticket_id, setTicketId] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  // State cho bảng danh sách vé
  const [allTickets, setAllTickets] = useState([]);
  const [isTableLoading, setIsTableLoading] = useState(false);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
  });

  // State cho việc in và đếm ảnh đã tải, giúp tránh lỗi mất QR
  const [ticketToPrint, setTicketToPrint] = useState(null);
  const [loadedImages, setLoadedImages] = useState(0);

  // --- API & Data Fetching ---

  // Lấy vé theo mã ID, hiển thị chi tiết ngay lập tức
  const getTicket = async () => {
    if (!ticket_id.trim()) {
      setErrorMessage("Vui lòng nhập mã vé.");
      return;
    }
    setIsLoading(true);
    setTicket({}); // Reset vé cũ khi tìm kiếm vé mới
    setErrorMessage("");
    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const response = await axios.get(
        `${API_URL}api/v1/public/ticket/${ticket_id}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      if (response.data?.data?.ticket) {
        setTicket(response.data.data.ticket);
      } else {
        setTicket({});
        setErrorMessage("Không tìm thấy vé.");
      }
    } catch (error) {
      setTicket({});
      setErrorMessage(
        error.response?.data?.message || "Không thể lấy thông tin vé."
      );
    } finally {
      setIsLoading(false);
    }
  };

  // Lấy toàn bộ vé cho bảng
  const fetchAllTickets = async (page = 1) => {
    setIsTableLoading(true);
    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const response = await axios.get(
        `${API_URL}api/v1/tickets/all?page=${page}&limit=${pagination.limit}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      const { tickets, total } = response.data.data;
      setAllTickets(tickets || []);
      setPagination((prev) => ({ ...prev, page, total }));
    } catch (error) {
      console.error("Failed to fetch all tickets:", error);
      setErrorMessage("Không thể tải danh sách vé.");
    } finally {
      setIsTableLoading(false);
    }
  };

  const generateQrCodesForSeats = async (seats, ticketId) => {
    if (!seats || seats.length === 0) return { urls: {} };
    const qrPromises = seats.map(async (seat) => {
      if (!seat?.seat_id) return null;
      const content = `TICKET:${ticketId}-SEAT:${seat.seat_id}`;
      const qrApiUrl = `${API_URL}api/v1/qr/image?content=${encodeURIComponent(
        content
      )}`;
      try {
        const response = await fetch(qrApiUrl);
        if (!response.ok) throw new Error("QR fetch failed");
        const blob = await response.blob();
        return { seatId: seat.seat_id, url: URL.createObjectURL(blob) };
      } catch {
        return { seatId: seat.seat_id, url: null };
      }
    });
    const results = await Promise.all(qrPromises);
    return {
      urls: results.filter(Boolean).reduce((acc, res) => {
        acc[res.seatId] = res.url;
        return acc;
      }, {}),
    };
  };

  // --- Logic In Vé (đã sửa lỗi) ---

  useEffect(() => {
    if (
      ticketToPrint &&
      loadedImages > 0 &&
      loadedImages === ticketToPrint.SeatTicketsBegin.length
    ) {
      window.print();
      setTicketToPrint(null);
      setLoadedImages(0);
    }
  }, [ticketToPrint, loadedImages]);

  const handlePrint = async (ticketData) => {
    setLoadedImages(0);
    const { urls } = await generateQrCodesForSeats(
      ticketData.SeatTicketsBegin,
      ticketData.ticket_id
    );
    setTicketToPrint({ ...ticketData, qrUrls: urls });
  };

  // --- Component Lifecycle & Helpers ---

  useEffect(() => {
    fetchAllTickets(1);
  }, []);

  const refresh = () => {
    setTicketId("");
    setTicket({});
    setErrorMessage("");
  };

  const handlePageChange = (newPage) => {
    if (newPage > 0 && newPage <= Math.ceil(pagination.total / pagination.limit)) {
      fetchAllTickets(newPage);
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 0: return "Đã đặt";
      case 1: return "Đã xác nhận";
      case 2: return "Đã hủy";
      case 4: return "Hoàn thành";
      default: return "Không xác định";
    }
  };

  return (
    <div className="w-full p-6 dark:bg-gray-900 min-h-screen">
      <style>{`
        @media print { body * { visibility: hidden; } .print-area, .print-area * { visibility: visible; } .print-area { position: absolute; left: 0; top: 0; width: 100%; font-family: 'Inter', 'Roboto', sans-serif; } .hidden-print { display: none !important; } .page-break { page-break-after: always; } .ticket-card { width: 85mm; height: auto; margin: 5mm auto; border: none; box-shadow: none; } .qr-code { width: 60mm; height: 60mm; } }
        .input-field { transition: all 0.3s ease; border-color: #d1d5db; }
        .input-field:focus { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1); outline: none; }
        .btn-primary { background-color: #2563eb; transition: background-color 0.3s ease; }
        .btn-primary:hover { background-color: #1d4ed8; }
        .btn-secondary { background-color: #6b7280; transition: background-color 0.3s ease; }
        .btn-secondary:hover { background-color: #4b5563; }
        .ticket-card { width: 85mm; max-width: 100%; border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden; background: white; box-shadow: 0 4px 12px rgba(0,0,0,0.1); margin: 16px auto; }
        .ticket-header { background: linear-gradient(to right, #2563eb, #60a5fa); color: white; padding: 12px; text-align: center; border-bottom: 4px solid #ffffff; }
        .ticket-body { padding: 16px; font-size: 13px; color: #1f2937; }
        .ticket-footer { border-top: 2px dashed #d1d5db; padding: 12px; text-align: center; background: #f9fafb; }
        .qr-code { border: 2px solid #e5e7eb; padding: 8px; background: white; margin: 0 auto; width: 60mm; height: 60mm; }
        .ticket-details { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; line-height: 1.5; }
        .ticket-details strong { color: #111827; }
        .pagination-btn { margin: 0 4px; padding: 8px 12px; border-radius: 8px; background-color: #f3f4f6; color: #374151; font-weight: 500;}
        .pagination-btn.active { background-color: #2563eb; color: white; }
        .pagination-btn:disabled { opacity: 0.5; cursor: not-allowed; }
      `}</style>
      
      {/* --- Phần tìm kiếm vé theo mã --- */}
      <div className="hidden-print">
        <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white tracking-tight">Quản Lý Vé</h1>
        <div className="max-w-md mx-auto bg-white dark:bg-gray-800 p-6 rounded-xl shadow-lg">
          <div className="mb-5">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">Mã vé</label>
            <input value={ticket_id} onChange={(e) => setTicketId(e.target.value)} type="text" className="input-field w-full p-3 border rounded-lg bg-gray-50 dark:bg-gray-700 dark:border-gray-600 dark:text-white text-sm" placeholder="Nhập mã vé của bạn" />
          </div>
          <div className="flex gap-3 mb-4">
            <button onClick={getTicket} className="btn-primary flex-1 py-2.5 text-white rounded-lg font-medium text-sm" disabled={isLoading}>{isLoading ? "Đang tìm..." : "Xem vé"}</button>
            <button onClick={refresh} className="btn-secondary flex-1 py-2.5 text-white rounded-lg font-medium text-sm">Làm mới</button>
          </div>
          {errorMessage && <p className="text-red-500 text-sm bg-red-50 dark:bg-red-900/30 p-3 rounded-lg">{errorMessage}</p>}
        </div>
      </div>

      {/* --- Chi tiết vé tìm được (hiển thị ngay) --- */}
      {ticket?.ticket_id && (
        <div className="max-w-md mx-auto mt-6 hidden-print">
          <div className="flex gap-3 mb-4">
            <button onClick={() => handlePrint(ticket)} className="btn-primary flex-1 py-2.5 text-white rounded-lg font-medium text-sm">In vé</button>
          </div>
          {/* Đây là phần xem trước, phần để in thực sự nằm ở khu vực print-area ở cuối file */}
          {ticket.SeatTicketsBegin?.map((seat) => (
            <div key={seat.seat_id} className="ticket-card">
              <div className="ticket-header"><h2 className="text-lg font-bold">Vé xe khách</h2><p className="text-xs opacity-80">Comfort Travel Co.</p></div>
              <div className="ticket-body">
                <div className="ticket-details">
                  <p><strong>Mã vé:</strong> {ticket.ticket_id}</p>
                  <p><strong>Ghế:</strong> {seat.seat_name?.String || "N/A"}</p>
                  <p><strong>Họ tên:</strong> {ticket.name?.String || "N/A"}</p>
                  <p><strong>SĐT:</strong> {ticket.phone?.String || "N/A"}</p>
                  <p><strong>Giá vé:</strong> {ticket.price?.toLocaleString()} VND</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* --- Bảng danh sách toàn bộ vé --- */}
      <div className="mt-12 hidden-print">
        <h2 className="text-2xl font-bold mb-4 text-gray-800 dark:text-white">Toàn bộ vé</h2>
        <div className="bg-white dark:bg-gray-800 p-4 rounded-xl shadow-lg overflow-x-auto">
          {isTableLoading ? (
            <p className="text-center p-8">Đang tải dữ liệu...</p>
          ) : (
            <>
              <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                  <tr>
                    <th scope="col" className="px-6 py-3">Mã Vé</th>
                    <th scope="col" className="px-6 py-3">Khách Hàng</th>
                    <th scope="col" className="px-6 py-3">SĐT</th>
                    <th scope="col" className="px-6 py-3">Thời Gian Đặt</th>
                    <th scope="col" className="px-6 py-3">Trạng Thái</th>
                    <th scope="col" className="px-6 py-3">Thao Tác</th>
                  </tr>
                </thead>
                <tbody>
                  {allTickets.map((t) => (
                    <tr key={t.ticket_id} className="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                      <td className="px-6 py-4 font-medium text-gray-900 dark:text-white">{t.ticket_id}</td>
                      <td className="px-6 py-4">{t.name?.String || "N/A"}</td>
                      <td className="px-6 py-4">{t.phone?.String || "N/A"}</td>
                      <td className="px-6 py-4">{new Date(t.booking_time).toLocaleString("vi-VN")}</td>
                      <td className="px-6 py-4">{getStatusText(t.status)}</td>
                      <td className="px-6 py-4">
                        <button onClick={() => handlePrint(t)} className="font-medium text-blue-600 dark:text-blue-500 hover:underline">In Vé</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              <div className="flex justify-center items-center mt-6">
                <button onClick={() => handlePageChange(pagination.page - 1)} disabled={pagination.page <= 1} className="pagination-btn">Previous</button>
                <span className="p-2">Trang {pagination.page} / {Math.ceil(pagination.total / pagination.limit)}</span>
                <button onClick={() => handlePageChange(pagination.page + 1)} disabled={pagination.page >= Math.ceil(pagination.total / pagination.limit)} className="pagination-btn">Next</button>
              </div>
            </>
          )}
        </div>
      </div>

      {/* --- Khu vực riêng để in (ẩn) --- */}
      <div ref={printRef} className="print-area">
        {ticketToPrint?.SeatTicketsBegin?.map((seat) => (
          <div key={seat.seat_id} className="ticket-card page-break">
            <div className="ticket-header">
              <h2 className="text-lg font-bold">Vé xe khách</h2>
              <p className="text-xs opacity-80">Comfort Travel Co.</p>
            </div>
            <div className="ticket-body">
              <div className="ticket-details">
                <p><strong>Mã vé:</strong> {ticketToPrint.ticket_id}</p>
                <p><strong>Ghế:</strong> {seat.seat_name?.String || "N/A"}</p>
                <p><strong>Họ tên:</strong> {ticketToPrint.name?.String || "N/A"}</p>
                <p><strong>Số điện thoại:</strong> {ticketToPrint.phone?.String || "N/A"}</p>
                <p><strong>Giá vé:</strong> {ticketToPrint.price?.toLocaleString() || "N/A"} VND</p>
                <p><strong>Thời gian đặt:</strong> {new Date(ticketToPrint.booking_time).toLocaleString("vi-VN")}</p>
              </div>
            </div>
            <div className="ticket-footer">
              {ticketToPrint.qrUrls?.[seat.seat_id] ? (
                <img
                  src={ticketToPrint.qrUrls[seat.seat_id]}
                  alt="QR Code"
                  className="qr-code"
                  onLoad={() => setLoadedImages((prev) => prev + 1)}
                />
              ) : (
                <p className="text-red-500 text-xs">Không thể tạo mã QR</p>
              )}
              <p className="text-xs text-gray-500 mt-2">
                Quét mã để xác nhận vé
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TicketManagement;