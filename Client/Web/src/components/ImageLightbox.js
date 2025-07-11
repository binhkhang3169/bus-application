/** @format */

import React from "react";

function ImageLightbox({ imageUrl, onClose }) {
  if (!imageUrl) return null;

  // Ngăn sự kiện click từ ảnh lan ra nền (để không bị đóng khi click vào ảnh)
  const handleImageClick = (e) => {
    e.stopPropagation();
  };

  return (
    // Lớp phủ nền
    <div
      className="fixed inset-0 bg-black bg-opacity-75 flex justify-center items-center z-[100]"
      onClick={onClose} // Đóng khi click vào nền
    >
      {/* Nút đóng ở góc trên phải */}
      <button
        onClick={onClose}
        className="absolute top-4 right-4 text-white text-4xl font-bold hover:text-gray-300"
        aria-label="Đóng"
      >
        &times;
      </button>

      {/* Container cho ảnh */}
      <div className="relative p-4">
        <img
          src={imageUrl}
          alt="Mã QR được phóng to"
          className="max-w-[90vw] max-h-[90vh] object-contain"
          onClick={handleImageClick}
        />
      </div>
    </div>
  );
}

export default ImageLightbox;