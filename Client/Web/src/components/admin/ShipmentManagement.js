/** @format */

import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { API_URL } from "../../configs/env"; // Đảm bảo bạn đã cấu hình API_URL
import ShipmentForm from "./modal/ShipmentForm";
import ShipmentDetailModal from "./modal/ShipmentDetailModal";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";

// Hàm trợ giúp để xử lý thông báo lỗi an toàn
const getSafeErrorMessage = (error, defaultMessage) => {
  const apiError = error.response?.data?.error || error.response?.data?.message;
  if (typeof apiError === 'string') {
    return apiError;
  }
  if (typeof apiError === 'object' && apiError !== null) {
    return JSON.stringify(apiError); // Chuyển đối tượng lỗi thành chuỗi
  }
  return defaultMessage;
};

const ShipmentManagement = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [shipments, setShipments] = useState([]);
  const [tripIdFilter, setTripIdFilter] = useState("");

  // States for modals
  const [shipmentFormModal, setShipmentFormModal] = useState(false);
  const [detailModal, setDetailModal] = useState({
    isOpen: false,
    shipment: null,
    invoice: null,
  });
  const [notification, setNotification] = useState({
    success: false,
    failure: false,
    message: "",
  });

  // Hàm gọi API để lấy danh sách lô hàng
  const fetchShipments = useCallback(async (tripId) => {
    setIsLoading(true);
    const token = sessionStorage.getItem("adminAccessToken");
    if (!token) {
      console.error("Authorization token not found.");
      setIsLoading(false);
      return;
    }

    const url = tripId
      ? `${API_URL}api/v1/tripsshipments/${tripId}/shipments`
      : `${API_URL}api/v1/shipments`;

    try {
      const response = await axios.get(url, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setShipments(response.data || []);
    } catch (error) {
      console.error("Error fetching shipments:", error);
      setShipments([]);
      setNotification({
        failure: true,
        message: getSafeErrorMessage(error, "Không thể tải danh sách lô hàng."),
      });
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Tải tất cả lô hàng khi component được mount
  useEffect(() => {
    fetchShipments();
  }, [fetchShipments]);

  const handleFilterByTripId = () => {
    if (tripIdFilter.trim()) {
      fetchShipments(tripIdFilter.trim());
    }
  };

  const handleClearFilter = () => {
    setTripIdFilter("");
    fetchShipments();
  };

  const handleViewDetails = async (shipment) => {
    setIsLoading(true);
    const token = sessionStorage.getItem("adminAccessToken");
    if (!token) {
        setIsLoading(false);
        return;
    }

    try {
      const shipmentUrl = `${API_URL}api/v1/shipments/${shipment.id}`;
      const invoiceUrl = `${API_URL}api/v1/shipments/${shipment.id}/invoice`;

      const [shipmentRes, invoiceRes] = await Promise.all([
          axios.get(shipmentUrl, { headers: { Authorization: `Bearer ${token}` } }),
          axios.get(invoiceUrl, { headers: { Authorization: `Bearer ${token}` } })
      ]);

      setDetailModal({
        isOpen: true,
        shipment: shipmentRes.data,
        invoice: invoiceRes.data,
      });
    } catch (error) {
      console.error("Error fetching shipment details:", error);
      setNotification({
        failure: true,
        message: getSafeErrorMessage(error, "Không thể lấy thông tin chi tiết."),
      });
    } finally {
      setIsLoading(false);
    }
  };

  const closeDetailModal = () => setDetailModal({ isOpen: false, shipment: null, invoice: null });
  const openShipmentFormModal = () => setShipmentFormModal(true);
  const closeShipmentFormModal = () => setShipmentFormModal(false);
  const closeNotification = () => setNotification({ success: false, failure: false, message: "" });
  const refreshData = () => fetchShipments(tripIdFilter.trim() || null);

  return (
    <div className="w-full p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="font-bold text-2xl text-gray-800 dark:text-white">
          Quản lý Lô hàng
        </h1>
        <button
          onClick={openShipmentFormModal}
          className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center"
        >
          Thêm Lô hàng
        </button>
      </div>

      <div className="mb-4 flex items-center gap-2">
        <input
          type="number"
          value={tripIdFilter}
          onChange={(e) => setTripIdFilter(e.target.value)}
          placeholder="Nhập Trip ID để lọc"
          className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg p-2.5 w-64 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
        />
        <button onClick={handleFilterByTripId} className="px-4 py-2.5 text-sm font-medium text-white bg-green-600 rounded-lg hover:bg-green-700">
          Tìm kiếm
        </button>
        <button onClick={handleClearFilter} className="px-4 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700">
          Xóa bộ lọc
        </button>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      ) : (
        <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
          <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
              <tr>
                <th scope="col" className="px-6 py-3">ID Lô hàng</th>
                <th scope="col" className="px-6 py-3">ID Chuyến đi</th>
                <th scope="col" className="px-6 py-3">Tên hàng</th>
                <th scope="col" className="px-6 py-3">Người gửi</th>
                <th scope="col" className="px-6 py-3">Người nhận</th>
                <th scope="col" className="px-6 py-3 text-center">Hành động</th>
              </tr>
            </thead>
            <tbody>
              {shipments.length === 0 ? (
                <tr>
                  <td colSpan="6" className="px-6 py-4 text-center">
                    Không có lô hàng nào.
                  </td>
                </tr>
              ) : (
                shipments.map((shipment) => (
                  <tr key={shipment.id} className="bg-white border-b hover:bg-gray-50 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-600">
                    <td className="px-6 py-4 font-medium text-gray-900 dark:text-white">{shipment.id}</td>
                    <td className="px-6 py-4">{shipment.trip_id}</td>
                    <td className="px-6 py-4">{shipment.item_name}</td>
                    <td className="px-6 py-4">{shipment.sender_name}</td>
                    <td className="px-6 py-4">{shipment.receiver_name}</td>
                    <td className="px-6 py-4 text-center">
                      <button
                        onClick={() => handleViewDetails(shipment)}
                        className="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                      >
                        Xem chi tiết
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}

      {shipmentFormModal && (
        <ShipmentForm
          func={{
            closeModal: closeShipmentFormModal,
            refresh: refreshData,
            setNotification,
          }}
        />
      )}

      {detailModal.isOpen && (
        <ShipmentDetailModal
          shipment={detailModal.shipment}
          invoice={detailModal.invoice}
          onClose={closeDetailModal}
        />
      )}

      {notification.success && (
        <SuccessNotification func={{ closeModal: closeNotification }} message={notification.message} />
      )}
      {notification.failure && (
        <FailureNotification func={{ closeModal: closeNotification }} message={notification.message} />
      )}
    </div>
  );
};

export default ShipmentManagement;