/** @format */

import React from "react";

// Thêm prop onImageClick để xử lý việc phóng to
function SeatQrList({ title, seats, qrCodeUrls, qrCodeErrors, onImageClick }) {
  if (!seats || seats.length === 0) {
    return null;
  }

  return (
    <div className="mb-6">
      <h4 className="text-lg font-semibold mb-3 text-gray-800">{title}</h4>
      <div className="flex space-x-4 overflow-x-auto pb-4">
        {seats.map((seat, index) => {
          const qrImageUrl = qrCodeUrls[seat.seat_id];
          const hasError = qrCodeErrors[seat.seat_id];

          return (
            <div
              key={seat.id || seat.seat_id || index}
              className="p-3 border rounded-lg text-center bg-white shadow-md flex-shrink-0 w-48"
            >
              <p className="text-base font-medium text-gray-800">
                Ghế {seat.seat_id}
              </p>
              <p className="text-sm mb-2 text-gray-600">
                {seat.status === 1 ? "Đã xác nhận" : "Chờ"}
              </p>

              {/* Thêm position: relative để định vị nút download */}
              <div className="relative w-36 h-36 mx-auto mt-2 flex items-center justify-center bg-gray-50 rounded">
                {qrImageUrl ? (
                  <>
                    <img
                      src={qrImageUrl}
                      alt={`Mã QR cho ghế ${seat.seat_id}`}
                      className="w-full h-full object-contain cursor-pointer"
                      onClick={() => onImageClick(qrImageUrl)} // <-- GỌI HÀM PHÓNG TO
                    />
                    {/* NÚT TẢI VỀ */}
                    <a
                      href={qrImageUrl}
                      download={`QR_Seat_${seat.seat_id}.png`}
                      className="absolute top-1 right-1 bg-white bg-opacity-70 p-1 rounded-full text-gray-700 hover:bg-opacity-100 hover:text-black transition-opacity"
                      title="Tải mã QR"
                      onClick={(e) => e.stopPropagation()} // Ngăn sự kiện click lan ra ảnh
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    </a>
                  </>
                ) : hasError ? (
                  <p className="text-xs text-red-500 italic">
                    Lỗi tải mã QR.
                  </p>
                ) : (
                  <p className="text-xs text-gray-400 italic">Đang tải...</p>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default SeatQrList;