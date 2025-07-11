/** @format */

import React, { useState, useEffect, useCallback } from "react";
import api from "../../../services/apiService"; // Giả sử apiService đã được cấu hình
import axios from "axios";
import dayjs from "dayjs";
import { API_URL } from "../../../configs/env"; // Giả sử bạn đã cấu hình API_URL

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

const ShipmentForm = ({ func }) => {
  const [isLoading, setIsLoading] = useState(false);

  // Step 1: Trip Search
  const [locations, setLocations] = useState([]);
  const [startLocationId, setStartLocationId] = useState("");
  const [endLocationId, setEndLocationId] = useState("");
  const [searchDate, setSearchDate] = useState(dayjs().format("YYYY-MM-DD"));
  const [foundTrips, setFoundTrips] = useState([]);
  const [isSearching, setIsSearching] = useState(false);

  // Step 2: Shipment Details
  const [selectedTrip, setSelectedTrip] = useState(null);

  // Form fields
  const [senderName, setSenderName] = useState("");
  const [receiverName, setReceiverName] = useState("");
  const [itemName, setItemName] = useState("");
  const [itemType, setItemType] = useState("document");
  const [weight, setWeight] = useState("");
  const [length, setLength] = useState("");
  const [width, setWidth] = useState("");
  const [height, setHeight] = useState("");
  const [price, setPrice] = useState(0); // Giá sẽ được tính tự động
  const [payerType, setPayerType] = useState("sender");
  const [note, setNote] = useState("");
  const [formErrors, setFormErrors] = useState({});
  
  const calculatePrice = useCallback(() => {
    const numWeight = parseFloat(weight);
    const numLength = parseFloat(length);
    const numWidth = parseFloat(width);
    const numHeight = parseFloat(height);

    if (numWeight > 0 && numLength > 0 && numWidth > 0 && numHeight > 0) {
      const BASE_PRICE = 15000;
      const PRICE_PER_KG = 2500;
      const DIMENSIONAL_FACTOR = 5000;
      const volumetricWeight = (numLength * numWidth * numHeight) / DIMENSIONAL_FACTOR;
      const chargeableWeight = Math.max(numWeight, volumetricWeight);
      let itemTypeMultiplier = 1.0;
      if (itemType === "electronics") itemTypeMultiplier = 1.5;
      else if (itemType === "furniture") itemTypeMultiplier = 1.3;
      const calculated = (BASE_PRICE + (chargeableWeight * PRICE_PER_KG)) * itemTypeMultiplier;
      setPrice(Math.round(calculated / 1000) * 1000);
    } else {
      setPrice(0);
    }
  }, [weight, length, width, height, itemType]);
  
  useEffect(() => {
    calculatePrice();
  }, [calculatePrice]);
  
  useEffect(() => {
    const fetchLocations = async () => {
      try {
        const res = await api.get(`/provinces`);
        const activeLocations = res.data?.data.filter((loc) => loc.status === 1) || [];
        setLocations(activeLocations);
      } catch (error) {
        func.setNotification({ failure: true, message: getSafeErrorMessage(error, "Không thể tải danh sách địa điểm.") });
      }
    };
    fetchLocations();
  }, [func]);

  const handleSearchTrips = async () => {
    if (!startLocationId || !endLocationId || !searchDate) {
      func.setNotification({ failure: true, message: "Vui lòng chọn điểm đi, điểm đến và ngày đi." });
      return;
    }
    setIsSearching(true);
    setFoundTrips([]);
    try {
      const startLocationName = locations.find((l) => l.id == startLocationId)?.name;
      const endLocationName = locations.find((l) => l.id == endLocationId)?.name;
      const response = await api.get(
        `trips/search?from=${encodeURIComponent(startLocationName)}&to=${encodeURIComponent(endLocationName)}&fromId=${startLocationId}&toId=${endLocationId}&fromTime=${searchDate}`
      );
      setFoundTrips(response.data?.data || []);
    } catch (error) {
      func.setNotification({ failure: true, message: getSafeErrorMessage(error, "Không tìm thấy chuyến đi nào.") });
    } finally {
      setIsSearching(false);
    }
  };

  const validateForm = () => {
    const errors = {};
    if (!senderName.trim()) errors.senderName = "Tên người gửi là bắt buộc.";
    if (!receiverName.trim()) errors.receiverName = "Tên người nhận là bắt buộc.";
    if (!itemName.trim()) errors.itemName = "Tên hàng hóa là bắt buộc.";
    if (parseFloat(weight) <= 0 || !weight) errors.weight = "Cân nặng phải là số dương.";
    if (parseFloat(length) <= 0 || !length) errors.length = "Chiều dài phải là số dương.";
    if (parseFloat(width) <= 0 || !width) errors.width = "Chiều rộng phải là số dương.";
    if (parseFloat(height) <= 0 || !height) errors.height = "Chiều cao phải là số dương.";
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm() || !selectedTrip) return;

    setIsLoading(true);
    const payload = {
      trip_id: selectedTrip.tripId,
      sender_name: senderName,
      receiver_name: receiverName,
      item_name: itemName,
      item_type: itemType,
      weight: parseFloat(weight),
      dimensions: {
        length: parseFloat(length),
        width: parseFloat(width),
        height: parseFloat(height),
      },
      price: price,
      payer_type: payerType,
      note: note,
    };

    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const url = `${API_URL}api/v1/tripsshipments/${selectedTrip.tripId}/shipments`;
      await axios.post(url, payload, { headers: { Authorization: `Bearer ${token}` } });
      func.setNotification({ success: true, message: "Tạo lô hàng thành công!" });
      func.refresh();
      func.closeModal();
    } catch (error) {
      func.setNotification({ failure: true, message: getSafeErrorMessage(error, "Tạo lô hàng thất bại.") });
    } finally {
      setIsLoading(false);
    }
  };

  const inputClasses = "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white";
  const labelClasses = "block mb-2 text-sm font-medium text-gray-900 dark:text-white";

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex items-center justify-center">
      <div className="relative w-full max-w-2xl bg-white rounded-lg shadow-lg max-h-[90vh] overflow-y-auto dark:bg-gray-800">
        <div className="flex items-center justify-between p-4 border-b sticky top-0 bg-white z-10 dark:bg-gray-800 dark:border-gray-700">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
            {selectedTrip ? `Tạo đơn hàng cho chuyến #${selectedTrip.tripId}` : "Bước 1: Chọn chuyến xe"}
          </h3>
          <button onClick={func.closeModal} type="button" className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white">
             <svg className="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14"><path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6" /></svg>
          </button>
        </div>

        <div className="p-5">
          {!selectedTrip ? (
            <div className="space-y-4">
               <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                 <div>
                    <label className={labelClasses}>Điểm đi</label>
                    <select value={startLocationId} onChange={(e) => setStartLocationId(e.target.value)} className={inputClasses}>
                        <option value="">Chọn địa điểm</option>
                        {locations.map((loc) => (<option key={loc.id} value={loc.id}>{loc.name}</option>))}
                    </select>
                 </div>
                 <div>
                    <label className={labelClasses}>Điểm đến</label>
                    <select value={endLocationId} onChange={(e) => setEndLocationId(e.target.value)} className={inputClasses}>
                        <option value="">Chọn địa điểm</option>
                        {locations.map((loc) => (<option key={loc.id} value={loc.id}>{loc.name}</option>))}
                    </select>
                 </div>
                 <div>
                    <label className={labelClasses}>Ngày</label>
                    <input type="date" value={searchDate} onChange={(e) => setSearchDate(e.target.value)} min={dayjs().format("YYYY-MM-DD")} className={inputClasses} />
                 </div>
               </div>
               <div className="flex justify-center">
                  <button onClick={handleSearchTrips} disabled={isSearching} className="text-white bg-green-600 hover:bg-green-700 font-medium rounded-lg text-sm px-5 py-2.5 text-center">
                    {isSearching ? "Đang tìm..." : "Tìm kiếm"}
                  </button>
               </div>
               {foundTrips.length > 0 && (
                 <div className="mt-4 border-t pt-4 dark:border-gray-700">
                    <h4 className="font-semibold mb-2 dark:text-white">Các chuyến xe phù hợp:</h4>
                    <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                        {foundTrips.map((trip) => (
                            <li key={trip.tripId} className="py-3 flex justify-between items-center">
                                <div>
                                    <p className="font-medium text-gray-900 dark:text-white">{`${trip.departureStation} → ${trip.arrivalStation}`}</p>
                                    <p className="text-sm text-gray-500 dark:text-gray-400">{`Giờ đi: ${trip.departureTime.substring(0,5)} - Xe: ${trip.license}`}</p>
                                </div>
                                <button onClick={() => setSelectedTrip(trip)} className="text-white bg-blue-600 hover:bg-blue-700 font-medium rounded-lg text-sm px-4 py-2">Chọn</button>
                            </li>
                        ))}
                    </ul>
                 </div>
               )}
            </div>
          ) : (
            <form className="space-y-4">
              <div className="p-3 bg-blue-50 border border-blue-200 rounded-lg dark:bg-gray-700 dark:border-blue-500">
                  <p className="font-semibold text-gray-900 dark:text-white">Chuyến xe đã chọn:</p>
                  <p className="text-sm text-blue-800 dark:text-blue-300">{`${selectedTrip.departureStation} → ${selectedTrip.arrivalStation}`}</p>
                  <p className="text-sm text-blue-800 dark:text-blue-300">{`Ngày: ${selectedTrip.departureDate}, Giờ: ${selectedTrip.departureTime.substring(0, 5)}`}</p>
                  <button onClick={() => setSelectedTrip(null)} className="text-xs text-red-600 hover:underline mt-1 dark:text-red-400">Đổi chuyến khác</button>
              </div>

              <div className="grid grid-cols-2 gap-4">
                  <div>
                      <label className={labelClasses}>Tên người gửi <span className="text-red-500">*</span></label>
                      <input type="text" value={senderName} onChange={(e) => setSenderName(e.target.value)} className={`${inputClasses} ${formErrors.senderName ? "border-red-500" : ""}`} />
                      {formErrors.senderName && <p className="text-red-500 text-xs mt-1">{formErrors.senderName}</p>}
                  </div>
                  <div>
                      <label className={labelClasses}>Tên người nhận <span className="text-red-500">*</span></label>
                      <input type="text" value={receiverName} onChange={(e) => setReceiverName(e.target.value)} className={`${inputClasses} ${formErrors.receiverName ? "border-red-500" : ""}`} />
                      {formErrors.receiverName && <p className="text-red-500 text-xs mt-1">{formErrors.receiverName}</p>}
                  </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                  <div>
                      <label className={labelClasses}>Tên hàng <span className="text-red-500">*</span></label>
                      <input type="text" value={itemName} onChange={(e) => setItemName(e.target.value)} className={`${inputClasses} ${formErrors.itemName ? "border-red-500" : ""}`} />
                      {formErrors.itemName && <p className="text-red-500 text-xs mt-1">{formErrors.itemName}</p>}
                  </div>
                  <div>
                      <label className={labelClasses}>Loại hàng <span className="text-red-500">*</span></label>
                      <select value={itemType} onChange={(e) => setItemType(e.target.value)} className={inputClasses}>
                          <option value="document">Tài liệu</option>
                          <option value="electronics">Đồ điện tử</option>
                          <option value="furniture">Nội thất</option>
                      </select>
                  </div>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div>
                      <label className={labelClasses}>Cân nặng (kg) <span className="text-red-500">*</span></label>
                      <input type="number" value={weight} onChange={(e) => setWeight(e.target.value)} className={`${inputClasses} ${formErrors.weight ? "border-red-500" : ""}`} />
                      {formErrors.weight && <p className="text-red-500 text-xs mt-1">{formErrors.weight}</p>}
                  </div>
                  <div>
                      <label className={labelClasses}>Dài (cm) <span className="text-red-500">*</span></label>
                      <input type="number" value={length} onChange={(e) => setLength(e.target.value)} className={`${inputClasses} ${formErrors.length ? "border-red-500" : ""}`} />
                      {formErrors.length && <p className="text-red-500 text-xs mt-1">{formErrors.length}</p>}
                  </div>
                  <div>
                      <label className={labelClasses}>Rộng (cm) <span className="text-red-500">*</span></label>
                      <input type="number" value={width} onChange={(e) => setWidth(e.target.value)} className={`${inputClasses} ${formErrors.width ? "border-red-500" : ""}`} />
                      {formErrors.width && <p className="text-red-500 text-xs mt-1">{formErrors.width}</p>}
                  </div>
                  <div>
                      <label className={labelClasses}>Cao (cm) <span className="text-red-500">*</span></label>
                      <input type="number" value={height} onChange={(e) => setHeight(e.target.value)} className={`${inputClasses} ${formErrors.height ? "border-red-500" : ""}`} />
                      {formErrors.height && <p className="text-red-500 text-xs mt-1">{formErrors.height}</p>}
                  </div>
              </div>
              
               <div className="grid grid-cols-2 gap-4">
                  <div>
                      <label className={labelClasses}>Giá cước ước tính (VNĐ)</label>
                      <input
                          type="text"
                          value={price.toLocaleString('vi-VN')}
                          readOnly
                          className="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5 cursor-not-allowed dark:bg-gray-600 dark:border-gray-500"
                      />
                  </div>
                  <div>
                      <label className={labelClasses}>Người trả phí <span className="text-red-500">*</span></label>
                      <select value={payerType} onChange={(e) => setPayerType(e.target.value)} className={inputClasses}>
                          <option value="sender">Người gửi</option>
                          <option value="receiver">Người nhận</option>
                      </select>
                  </div>
              </div>

              <div>
                  <label className={labelClasses}>Ghi chú</label>
                  <textarea value={note} onChange={(e) => setNote(e.target.value)} rows="2" className={inputClasses}></textarea>
              </div>
            </form>
          )}
        </div>

        {selectedTrip && (
            <div className="flex justify-end items-center p-4 border-t sticky bottom-0 bg-white z-10 dark:bg-gray-800 dark:border-gray-700">
                <button onClick={func.closeModal} className="mr-2 px-4 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-500 dark:hover:bg-gray-600">Hủy</button>
                <button type="button" onClick={handleSubmit} disabled={isLoading} className="px-4 py-2 text-sm font-medium text-white bg-blue-700 rounded-lg hover:bg-blue-800 flex items-center dark:bg-blue-600 dark:hover:bg-blue-700">
                    {isLoading && <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white mr-2"></div>}
                    Tạo Lô hàng
                </button>
            </div>
        )}
      </div>
    </div>
  );
};

export default ShipmentForm;