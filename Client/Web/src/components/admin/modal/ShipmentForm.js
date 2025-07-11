/** @format */

import React, { useState, useEffect } from "react";
import api from "../../../services/apiService"; // Giả sử apiService đã được cấu hình đúng
import dayjs from "dayjs";

const ShipmentForm = ({ func }) => {
  const [isLoading, setIsLoading] = useState(false);

  // --- States for Step 1: Trip Search ---
  const [locations, setLocations] = useState([]);
  const [startLocationId, setStartLocationId] = useState("");
  const [endLocationId, setEndLocationId] = useState("");
  const [searchDate, setSearchDate] = useState(dayjs().format("YYYY-MM-DD"));
  const [foundTrips, setFoundTrips] = useState([]);
  const [isSearching, setIsSearching] = useState(false);

  // --- States for Step 2: Shipment Details ---
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
  const [note, setNote] = useState("");
  const [formErrors, setFormErrors] = useState({});

  useEffect(() => {
    fetchLocations();
  }, []);

  const fetchLocations = async () => {
    try {
      // Thay đổi URL API cho đúng với cấu trúc của bạn
      const res = await api.get(`/provinces`);
      const activeLocations =
        res.data?.data.filter((loc) => loc.status === 1) || [];
      setLocations(activeLocations);
    } catch (error) {
      console.error("Failed to load locations:", error);
      func.setMessage("Cannot load locations data.");
      func.openFailureModal();
    }
  };

  const handleSearchTrips = async () => {
    if (!startLocationId || !endLocationId || !searchDate) {
      func.setMessage("Please select start, end locations and a date.");
      func.openFailureModal();
      return;
    }
    setIsSearching(true);
    setFoundTrips([]);
    try {
      // URL tìm kiếm cần được điều chỉnh cho đúng với backend
      const response = await api.get(
        `/trips/search?fromId=${startLocationId}&toId=${endLocationId}&date=${searchDate}`
      );
      setFoundTrips(response.data?.data || []);
    } catch (error) {
      console.error("Failed to search for trips:", error);
      func.setMessage(error.response?.data?.message || "Failed to find trips.");
      func.openFailureModal();
    } finally {
      setIsSearching(false);
    }
  };

  const handleSelectTrip = (trip) => {
    setSelectedTrip(trip);
  };

  const validateForm = () => {
    const errors = {};
    if (!senderName) errors.senderName = "Sender name is required.";
    if (!receiverName) errors.receiverName = "Receiver name is required.";
    if (!itemName) errors.itemName = "Item name is required.";
    if (!itemType) errors.itemType = "Item type is required.";
    if (parseFloat(weight) <= 0 || !weight)
      errors.weight = "Weight must be greater than 0.";
    if (parseFloat(length) <= 0 || !length)
      errors.length = "Length must be greater than 0.";
    if (parseFloat(width) <= 0 || !width)
      errors.width = "Width must be greater than 0.";
    if (parseFloat(height) <= 0 || !height)
      errors.height = "Height must be greater than 0.";
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm() || !selectedTrip) return;

    setIsLoading(true);
    const payload = {
      trip_id: selectedTrip.id,
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
      note: note,
    };

    try {
      // API để tạo shipment
      const response = await api.post(
        `/api/v1/trips/${selectedTrip.id}/shipments`,
        payload
      );
      func.setMessage(
        response.data.message || "Shipment created successfully!"
      );
      func.openSuccessModal();
      func.refresh();
      func.closeModal();
    } catch (error) {
      const errorMsg =
        error.response?.data?.error || "Failed to create shipment.";
      func.setMessage(errorMsg);
      func.openFailureModal();
    } finally {
      setIsLoading(false);
    }
  };

  // Định nghĩa các lớp CSS dùng chung để dễ quản lý
  const inputFieldClasses =
    "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500";
  const labelClasses =
    "block mb-2 text-sm font-medium text-gray-900 dark:text-white";

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex items-center justify-center">
      <div className="relative w-full max-w-2xl bg-white rounded-lg shadow-lg max-h-[90vh] overflow-y-auto dark:bg-gray-800">
        {/* Header with Dark Mode */}
        <div className="flex items-center justify-between p-4 border-b sticky top-0 bg-white z-10 dark:bg-gray-800 dark:border-gray-700">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
            {selectedTrip
              ? `Create Shipment for Trip #${selectedTrip.id}`
              : "Step 1: Find a Suitable Trip"}
          </h3>
          <button
            onClick={func.closeModal}
            type="button"
            className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
          >
            <svg
              className="w-3 h-3"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 14 14"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
              />
            </svg>
          </button>
        </div>

        {/* Body with Dark Mode */}
        <div className="p-5">
          {!selectedTrip ? (
            // --- STEP 1: TRIP SEARCH ---
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                <div>
                  <label className={labelClasses}>From</label>
                  <select
                    value={startLocationId}
                    onChange={(e) => setStartLocationId(e.target.value)}
                    className={inputFieldClasses}
                  >
                    <option value="">Select location</option>
                    {locations.map((loc) => (
                      <option key={loc.id} value={loc.id}>
                        {loc.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className={labelClasses}>To</label>
                  <select
                    value={endLocationId}
                    onChange={(e) => setEndLocationId(e.target.value)}
                    className={inputFieldClasses}
                  >
                    <option value="">Select location</option>
                    {locations.map((loc) => (
                      <option key={loc.id} value={loc.id}>
                        {loc.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className={labelClasses}>Date</label>
                  <input
                    type="date"
                    value={searchDate}
                    onChange={(e) => setSearchDate(e.target.value)}
                    min={dayjs().format("YYYY-MM-DD")}
                    className={inputFieldClasses}
                  />
                </div>
              </div>
              <div className="flex justify-center">
                <button
                  onClick={handleSearchTrips}
                  disabled={isSearching}
                  className="text-white bg-green-600 hover:bg-green-700 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-green-600 dark:hover:bg-green-700"
                >
                  {isSearching ? "Searching..." : "Find Trips"}
                </button>
              </div>

              {/* Trip Results with Dark Mode */}
              {foundTrips.length > 0 && (
                <div className="mt-4 border-t pt-4 dark:border-gray-700">
                  <h4 className="font-semibold mb-2 dark:text-white">
                    Available Trips:
                  </h4>
                  <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                    {foundTrips.map((trip) => (
                      <li
                        key={trip.id}
                        className="py-3 flex justify-between items-center"
                      >
                        <div>
                          <p className="font-medium text-gray-900 dark:text-white">{`${trip.route.start.name} → ${trip.route.end.name}`}</p>
                          <p className="text-sm text-gray-500 dark:text-gray-400">{`Time: ${trip.departureTime.substring(
                            0,
                            5
                          )} - Vehicle: ${trip.vehicle.license}`}</p>
                        </div>
                        <button
                          onClick={() => handleSelectTrip(trip)}
                          className="text-white bg-blue-600 hover:bg-blue-700 font-medium rounded-lg text-sm px-4 py-2"
                        >
                          Select
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ) : (
            // --- STEP 2: SHIPMENT DETAILS ---
            <form className="space-y-4">
              <div className="p-3 bg-blue-50 border border-blue-200 rounded-lg dark:bg-gray-700 dark:border-blue-500">
                <p className="font-semibold text-gray-900 dark:text-white">
                  Selected Trip:
                </p>
                <p className="text-sm text-blue-800 dark:text-blue-300">{`${selectedTrip.route.start.name} → ${selectedTrip.route.end.name}`}</p>
                <p className="text-sm text-blue-800 dark:text-blue-300">{`Date: ${
                  selectedTrip.departureDate
                }, Time: ${selectedTrip.departureTime.substring(0, 5)}`}</p>
                <button
                  onClick={() => setSelectedTrip(null)}
                  className="text-xs text-red-600 hover:underline mt-1 dark:text-red-400"
                >
                  Change trip
                </button>
              </div>

              {/* Sender & Receiver */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className={labelClasses}>
                    Sender Name <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={senderName}
                    onChange={(e) => setSenderName(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.senderName
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.senderName && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.senderName}
                    </p>
                  )}
                </div>
                <div>
                  <label className={labelClasses}>
                    Receiver Name <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={receiverName}
                    onChange={(e) => setReceiverName(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.receiverName
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.receiverName && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.receiverName}
                    </p>
                  )}
                </div>
              </div>

              {/* Item Details */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className={labelClasses}>
                    Item Name <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={itemName}
                    onChange={(e) => setItemName(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.itemName
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.itemName && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.itemName}
                    </p>
                  )}
                </div>
                <div>
                  <label className={labelClasses}>
                    Item Type <span className="text-red-500">*</span>
                  </label>
                  <select
                    value={itemType}
                    onChange={(e) => setItemType(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.itemType
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  >
                    <option value="document">Document</option>
                    <option value="electronics">Electronics</option>
                    <option value="furniture">Furniture</option>
                  </select>
                </div>
              </div>

              {/* Weight & Dimensions */}
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div>
                  <label className={labelClasses}>
                    Weight (kg) <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    value={weight}
                    onChange={(e) => setWeight(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.weight
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.weight && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.weight}
                    </p>
                  )}
                </div>
                <div>
                  <label className={labelClasses}>
                    Length (cm) <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    value={length}
                    onChange={(e) => setLength(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.length
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.length && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.length}
                    </p>
                  )}
                </div>
                <div>
                  <label className={labelClasses}>
                    Width (cm) <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    value={width}
                    onChange={(e) => setWidth(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.width
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.width && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.width}
                    </p>
                  )}
                </div>
                <div>
                  <label className={labelClasses}>
                    Height (cm) <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    value={height}
                    onChange={(e) => setHeight(e.target.value)}
                    className={`${inputFieldClasses} ${
                      formErrors.height
                        ? "border-red-500 dark:border-red-500"
                        : ""
                    }`}
                  />
                  {formErrors.height && (
                    <p className="text-red-500 text-xs mt-1 dark:text-red-400">
                      {formErrors.height}
                    </p>
                  )}
                </div>
              </div>

              {/* Note */}
              <div>
                <label className={labelClasses}>Note</label>
                <textarea
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  rows="2"
                  className={inputFieldClasses}
                ></textarea>
              </div>
            </form>
          )}
        </div>

        {/* Footer with Dark Mode */}
        {selectedTrip && (
          <div className="flex justify-end items-center p-4 border-t sticky bottom-0 bg-white z-10 dark:bg-gray-800 dark:border-gray-700">
            <button
              onClick={func.closeModal}
              className="mr-2 px-4 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-500 dark:hover:bg-gray-600"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleSubmit}
              disabled={isLoading}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-700 rounded-lg hover:bg-blue-800 flex items-center dark:bg-blue-600 dark:hover:bg-blue-700"
            >
              {isLoading && (
                <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white mr-2"></div>
              )}
              Create Shipment
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default ShipmentForm;
