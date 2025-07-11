/** @format */

import React from "react";
import dayjs from "dayjs";

// Component con để hiển thị một hàng chi tiết
const DetailRow = ({ label, value }) => (
  <div className="grid grid-cols-3 gap-4 py-2 border-b border-gray-200 dark:border-gray-700 last:border-b-0">
    <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">{label}</dt>
    <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2 dark:text-white break-words">{value}</dd>
  </div>
);

const ShipmentDetailModal = ({ shipment, invoice, onClose }) => {
  if (!shipment) return null;

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex items-center justify-center">
      <div className="relative w-full max-w-2xl bg-white rounded-lg shadow-lg max-h-[90vh] overflow-y-auto dark:bg-gray-800">
        <div className="flex items-center justify-between p-4 border-b sticky top-0 bg-white z-10 dark:bg-gray-800 dark:border-gray-700">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
            Chi tiết Lô hàng #{shipment.id}
          </h3>
          <button
            onClick={onClose}
            type="button"
            className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
          >
            <svg className="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
              <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6" />
            </svg>
          </button>
        </div>

        <div className="p-6 space-y-6">
          {/* Thông tin lô hàng */}
          <div>
            <h4 className="text-lg font-semibold text-gray-800 dark:text-white mb-3">Thông tin Lô hàng</h4>
            <dl>
              <DetailRow label="Tên hàng" value={shipment.item_name} />
              <DetailRow label="Loại hàng" value={shipment.item_type} />
              <DetailRow label="Người gửi" value={shipment.sender_name} />
              <DetailRow label="Người nhận" value={shipment.receiver_name} />
              <DetailRow label="Kích thước (Dài x Rộng x Cao)" value={`${shipment.length} x ${shipment.width} x ${shipment.height} cm`} />
              <DetailRow label="Cân nặng" value={`${shipment.weight} kg`} />
              <DetailRow label="Thế tích" value={`${shipment.volume} cm³`} />
              <DetailRow label="Người trả phí" value={shipment.payer_type} />
              <DetailRow label="Ghi chú" value={shipment.note || "Không có"} />
              <DetailRow label="Ngày tạo" value={dayjs(shipment.created_at).format("DD/MM/YYYY HH:mm")} />
            </dl>
          </div>

          {/* Thông tin hóa đơn */}
          <div>
            <h4 className="text-lg font-semibold text-gray-800 dark:text-white mb-3">Thông tin Hóa đơn</h4>
            {invoice ? (
              <dl>
                <DetailRow label="ID Hóa đơn" value={invoice.id} />
                <DetailRow label="Số tiền" value={`${invoice.amount.toLocaleString("vi-VN")} VNĐ`} />
                <DetailRow label="Ngày xuất hóa đơn" value={dayjs(invoice.issued_at).format("DD/MM/YYYY HH:mm")} />
              </dl>
            ) : (
              <p className="text-center text-gray-500 dark:text-gray-400">Không tìm thấy hóa đơn cho lô hàng này.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ShipmentDetailModal;