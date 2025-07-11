/** @format */

import React, { useState, useEffect } from "react";
import dayjs from "dayjs";
import api from "../../services/apiService";
import { API_URL } from "../../configs/env";
import SuccessNotification from "../../components/Noti/SuccessNotification";
import FailureNotification from "../../components/Noti/FailureNotification";

// Reusable TripCard Component
const TripCard = ({ trip, onSelect }) => {
  const [showDetails, setShowDetails] = useState(false);
  return (
    <div className="bg-white rounded-xl shadow-lg hover:shadow-2xl transition-shadow duration-300 overflow-hidden border border-gray-200">
      <div className="p-5 flex flex-col md:flex-row gap-6">
        <div className="flex-grow">
          <div className="flex items-center gap-4 mb-4">
            <div>
              <p className="font-bold text-lg text-gray-800 capitalize">
                {trip.car_type}
              </p>
              <p className="text-sm text-gray-500">{trip.estimated_distance}</p>
            </div>
          </div>
          <div className="flex gap-4">
            <div className="flex flex-col items-center">
              <div className="w-5 h-5 rounded-full border-2 border-green-500 bg-white"></div>
              <div className="flex-1 w-px bg-gray-300 border-l-2 border-dotted border-gray-300 my-2"></div>
              <div className="w-5 h-5 rounded-full bg-green-500"></div>
            </div>
            <div className="flex-grow">
              <div className="flex items-baseline gap-3">
                <p className="text-xl font-bold text-gray-900">
                  {trip.start_time}
                </p>
                <p className="font-semibold text-gray-700">
                  {trip.start_address}
                </p>
              </div>
              <p className="pl-8 py-2 text-sm text-gray-500">
                {trip.estimated_time}
              </p>
              <div className="flex items-baseline gap-3">
                <p className="text-xl font-bold text-gray-900">
                  {trip.end_time}
                </p>
                <p className="font-semibold text-gray-700">
                  {trip.end_address}
                </p>
              </div>
            </div>
          </div>
        </div>
        <div className="md:w-56 flex-shrink-0 flex flex-col items-stretch md:items-end justify-between md:border-l md:pl-6 border-gray-200">
          <div className="text-left md:text-right">
            <p className="text-2xl font-bold text-orange-600">
              {trip.price.toLocaleString()} đ
            </p>
            <p className="text-green-600 font-medium mt-1">
              {trip.available_seats} chỗ trống
            </p>
          </div>
          <div className="w-full mt-4 flex flex-col gap-2">
            <button
              onClick={onSelect}
              className="w-full text-base font-semibold text-white hover:bg-blue-600 transition-all bg-blue-500 py-3 px-4 rounded-lg shadow-md hover:shadow-lg"
            >
              Chọn vé
            </button>
            <button
              onClick={() => setShowDetails(!showDetails)}
              className="w-full text-sm font-medium text-gray-600 hover:text-green-600"
            >
              {showDetails ? "Ẩn chi tiết" : "Xem chi tiết"}
            </button>
          </div>
        </div>
      </div>
      {showDetails && (
        <div className="bg-gray-50 p-4 border-t border-gray-200">
          <p className="font-semibold text-sm text-gray-700">Lộ trình:</p>
          <p className="text-sm text-gray-600 mt-1">
            {trip.full_route.split("→").join(" → ")}
          </p>
        </div>
      )}
    </div>
  );
};

// Main Component: AdminBookingPage
const AdminBookingPage = () => {
  const [view, setView] = useState("SEARCH");
  const [provinces, setProvinces] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [searchParams, setSearchParams] = useState({
    fromId: "",
    toId: "",
    fromTime: dayjs().format("YYYY-MM-DD"),
  });
  const [searchError, setSearchError] = useState("");
  const [foundTrips, setFoundTrips] = useState([]);
  const [filteredTrips, setFilteredTrips] = useState([]);
  const [filters, setFilters] = useState({ timeOfDay: "", carType: "" });
  const [selectedTrip, setSelectedTrip] = useState(null);
  const [bookingDetails, setBookingDetails] = useState({
    seats: [],
    seatIds: [],
    pickupId: "",
    dropoffId: "",
  });
  const [customerInfo, setCustomerInfo] = useState({
    name: "",
    phone: "",
    email: "",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [message, setMessage] = useState("");
  const [errors, setErrors] = useState({}); // State for validation errors

  // Load provinces
  useEffect(() => {
    const fetchProvinces = async () => {
      try {
        const res = await api.get(API_URL + "api/v1/provinces");
        setProvinces(res.data?.data || []);
      } catch (err) {
        console.error("Error fetching provinces:", err);
      }
    };
    fetchProvinces();
  }, []);

  // Apply filters
  useEffect(() => {
    let filtered = [...foundTrips];
    if (filters.timeOfDay) {
      filtered = filtered.filter((trip) => {
        const hour = parseInt(trip.start_time.split(":")[0]);
        if (filters.timeOfDay === "sang") return hour >= 5 && hour < 12;
        if (filters.timeOfDay === "chieu") return hour >= 12 && hour < 18;
        if (filters.timeOfDay === "toi") return hour >= 18 || hour < 5;
        return true;
      });
    }
    if (filters.carType) {
      filtered = filtered.filter((trip) => trip.car_type === filters.carType);
    }
    setFilteredTrips(
      filtered.sort((a, b) => a.start_time.localeCompare(b.start_time))
    );
  }, [filters, foundTrips]);

  // Handle Trip Search
  const handleSearch = async (e) => {
    e.preventDefault();
    if (!searchParams.fromId || !searchParams.toId || !searchParams.fromTime) {
      setSearchError("Vui lòng chọn điểm đi, điểm đến và ngày khởi hành.");
      return;
    }
    setIsLoading(true);
    setSearchError("");
    setFoundTrips([]);

    const fromProvince = provinces.find(
      (p) => p.id === parseInt(searchParams.fromId, 10)
    );
    const toProvince = provinces.find(
      (p) => p.id === parseInt(searchParams.toId, 10)
    );

    try {
      const response = await api.get(
        `${API_URL}api/v1/trips/search?from=${fromProvince.name}&to=${toProvince.name}&fromId=${searchParams.fromId}&toId=${searchParams.toId}&fromTime=${searchParams.fromTime}`
      );
      const tripsData = response.data?.data || [];
      setFoundTrips(
        tripsData.map((trip) => ({
          ...trip,
          id: trip.tripId,
          start_address: trip.departureStation,
          end_address: trip.arrivalStation,
          date: trip.departureDate,
          start_time: trip.departureTime.substring(0, 5),
          end_time: trip.arrivalTime.substring(0, 5),
          price: trip.price,
          car_type: trip.vehicleType
            ? trip.vehicleType.toLowerCase()
            : "không xác định",
          available_seats: trip.stock || 0,
          full_route: trip.fullRoute,
        }))
      );
      setView("RESULTS");
    } catch (err) {
      setSearchError("Không thể tìm thấy chuyến xe. Vui lòng thử lại.");
    } finally {
      setIsLoading(false);
    }
  };

  // Handle Trip Selection
  const handleSelectTrip = async (trip) => {
    setIsLoading(true);
    try {
      const [tripDetailsResponse, availableSeatsResponse] = await Promise.all([
        api.get(`${API_URL}api/v1/trips/${trip.id}/seats`),
        api.get(`${API_URL}api/v1/tickets-available/${trip.id}`),
      ]);

      const tripData = tripDetailsResponse.data.data;
      const seatsData = availableSeatsResponse.data.data.seats;

      const availableSeatNames = seatsData.map((seat) => seat.name);
      const seatNameToIdMapping = seatsData.reduce((acc, seat) => {
        acc[seat.name] = seat.id;
        return acc;
      }, {});

      const parseRouteStations = (fullRouteString) => {
        if (!fullRouteString) return { pickup: [], dropoff: [] };
        const stationRegex = /([\w\sÀ-ỹ]+)\s*\((\d+)\)/g;
        const stations = [];
        let match;
        while ((match = stationRegex.exec(fullRouteString)) !== null) {
          stations.push({ name: match[1].trim(), id: parseInt(match[2], 10) });
        }
        return { pickup: stations.slice(0, -1), dropoff: stations.slice(1) };
      };
      const { pickup, dropoff } = parseRouteStations(tripData.fullRoute);

      setSelectedTrip({
        ...tripData,
        details: {
          availableSeats: availableSeatNames,
          seatMapping: seatNameToIdMapping,
          pickupLocations: pickup,
          dropoffLocations: dropoff,
        },
      });
      setView("BOOKING");
    } catch (error) {
      setMessage(error.message || "Không thể tải chi tiết chuyến xe.");
      setFailureModal(true);
    } finally {
      setIsLoading(false);
    }
  };

  // Validation function
  const validateBooking = () => {
    const newErrors = {};

    // Validate booking details
    if (bookingDetails.seatIds.length === 0) {
      newErrors.seats = "Vui lòng chọn ít nhất một ghế.";
    }
    if (!bookingDetails.pickupId) {
      newErrors.pickupId = "Vui lòng chọn điểm đón.";
    }
    if (!bookingDetails.dropoffId) {
      newErrors.dropoffId = "Vui lòng chọn điểm xuống.";
    }

    // Validate customer info
    if (!customerInfo.name.trim()) {
      newErrors.name = "Họ và tên là bắt buộc.";
    }
    if (!customerInfo.phone.trim()) {
      newErrors.phone = "Số điện thoại là bắt buộc.";
    } else if (!/^\d{10,11}$/.test(customerInfo.phone)) {
      newErrors.phone = "Số điện thoại phải có 10-11 số.";
    }
    if (
      customerInfo.email.trim() &&
      !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(customerInfo.email)
    ) {
      newErrors.email = "Email không hợp lệ.";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle Final Booking Confirmation
  const handleConfirmBooking = async () => {
    if (!validateBooking()) {
      setMessage("Vui lòng kiểm tra và điền đầy đủ các thông tin bắt buộc.");
      setFailureModal(true);
      return;
    }

    setIsSubmitting(true);
    const staffID =
      localStorage.getItem("userID") || sessionStorage.getItem("userID");

    const payload = {
      ticket_type: 0,
      price: selectedTrip.price,
      status: 1,
      payment_status: 1,
      booking_channel: 1,
      policy_id: 1,
      name: customerInfo.name,
      phone: customerInfo.phone,
      email: customerInfo.email,
      booked_by: `staff_${staffID}`,
      trip_id_begin: String(selectedTrip.tripId),
      seat_id_begin: bookingDetails.seatIds.map((id) => parseInt(id, 10)),
      pickup_location_begin: parseInt(bookingDetails.pickupId),
      dropoff_location_begin: parseInt(bookingDetails.dropoffId),
    };

    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const head = {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      };

      await api.post(`${API_URL}api/v1/staff/tickets`, payload, head);
      setMessage("Đặt vé thành công!");
      setSuccessModal(true);
      resetToSearch();
    } catch (error) {
      setMessage(error.response?.data?.message || "Không thể đặt vé.");
      setFailureModal(true);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Helper to reset the entire page
  const resetToSearch = () => {
    setView("SEARCH");
    setFoundTrips([]);
    setSelectedTrip(null);
    setBookingDetails({ seats: [], seatIds: [], pickupId: "", dropoffId: "" });
    setCustomerInfo({ name: "", phone: "", email: "" });
    setFilters({ timeOfDay: "", carType: "" });
    setErrors({});
  };

  // Seat map helper functions
  const isSeatAvailable = (seat) =>
    selectedTrip?.details.availableSeats.includes(seat);
  const isSeatSelected = (seat) => bookingDetails.seats.includes(seat);
  const getSeatImage = (seat) => {
    if (isSeatSelected(seat))
      return "https://futabus.vn/images/icons/seat_selecting.svg";
    if (isSeatAvailable(seat))
      return "https://futabus.vn/images/icons/seat_active.svg";
    return "https://futabus.vn/images/icons/seat_disabled.svg";
  };
  const generateAllSeats = (row) =>
    Array.from({ length: 15 }, (_, i) => `${row}${i + 1}`);

  // Render
  return (
    <div className="w-full p-4 md:p-8 min-h-full">
      <h1 className="font-bold text-2xl text-gray-800 mb-6">
        Đặt vé cho khách hàng
      </h1>

      {successModal && (
        <SuccessNotification
          func={{ closeModal: () => setSuccessModal(false) }}
          message={message}
        />
      )}
      {failureModal && (
        <FailureNotification
          func={{ closeModal: () => setFailureModal(false) }}
          message={message}
        />
      )}

      {/* Section 1: Search Form */}
      <div
        className={`bg-white p-6 rounded-lg shadow-md mb-8 ${
          view === "BOOKING" ? "hidden" : ""
        }`}
      >
        <h2 className="text-xl font-semibold mb-4">Bước 1: Tìm chuyến</h2>
        <form
          onSubmit={handleSearch}
          className="grid grid-cols-1 md:grid-cols-4 gap-4 items-end"
        >
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Điểm đi
            </label>
            <select
              className={`mt-1 block w-full p-2 border rounded-md ${
                searchError.includes("điểm đi")
                  ? "border-red-500"
                  : "border-gray-300"
              }`}
              value={searchParams.fromId}
              onChange={(e) =>
                setSearchParams({ ...searchParams, fromId: e.target.value })
              }
            >
              <option value="">Chọn điểm đi</option>
              {provinces.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Điểm đến
            </label>
            <select
              className={`mt-1 block w-full p-2 border rounded-md ${
                searchError.includes("điểm đến")
                  ? "border-red-500"
                  : "border-gray-300"
              }`}
              value={searchParams.toId}
              onChange={(e) =>
                setSearchParams({ ...searchParams, toId: e.target.value })
              }
            >
              <option value="">Chọn điểm đến</option>
              {provinces.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Ngày
            </label>
            <input
              type="date"
              className={`mt-1 block w-full p-2 border rounded-md ${
                searchError.includes("ngày")
                  ? "border-red-500"
                  : "border-gray-300"
              }`}
              value={searchParams.fromTime}
              onChange={(e) =>
                setSearchParams({ ...searchParams, fromTime: e.target.value })
              }
            />
          </div>
          <button
            type="submit"
            disabled={isLoading}
            className="bg-blue-600 text-white p-2.5 rounded-md hover:bg-blue-700 disabled:bg-gray-400"
          >
            {isLoading ? "Đang tìm..." : "Tìm kiếm"}
          </button>
        </form>
        {searchError && <p className="text-red-500 mt-2">{searchError}</p>}
      </div>

      {/* Section 2: Results and Filters */}
      {view === "RESULTS" && (
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          <div className="col-span-1 bg-white rounded-xl shadow-md p-5 h-fit">
            <h3 className="text-xl font-bold text-gray-800 mb-5">Bộ lọc</h3>
            <div className="space-y-6">
              <div>
                <label className="block mb-3 font-semibold text-gray-700">
                  Thời gian trong ngày
                </label>
                <div className="grid grid-cols-2 gap-2">
                  {["", "sang", "chieu", "toi"].map((val) => (
                    <button
                      key={val}
                      onClick={() =>
                        setFilters((prev) => ({ ...prev, timeOfDay: val }))
                      }
                      className={`py-2 px-3 text-sm rounded-lg border ${
                        filters.timeOfDay === val
                          ? "bg-blue-600 text-white"
                          : "hover:bg-gray-100"
                      }`}
                    >
                      {
                        {
                          "": "Tất cả",
                          sang: "Buổi sáng",
                          chieu: "Buổi chiều",
                          toi: "Buổi tối",
                        }[val]
                      }
                    </button>
                  ))}
                </div>
              </div>
              <div>
                <label className="block mb-3 font-semibold text-gray-700">
                  Loại xe
                </label>
                <div className="flex flex-col space-y-2">
                  {["", "ghế ngồi", "giường nằm", "limousine"].map((val) => (
                    <button
                      key={val}
                      onClick={() =>
                        setFilters((prev) => ({ ...prev, carType: val }))
                      }
                      className={`py-2 px-3 text-sm rounded-lg border text-left ${
                        filters.carType === val
                          ? "bg-blue-600 text-white"
                          : "hover:bg-gray-100"
                      }`}
                    >
                      {
                        {
                          "": "Tất cả",
                          "ghế ngồi": "Ghế ngồi",
                          "giường nằm": "Giường nằm",
                          limousine: "Limousine",
                        }[val]
                      }
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>
          <div className="col-span-1 lg:col-span-3">
            <h2 className="text-xl font-semibold mb-4">
              Bước 2: Kết quả ({filteredTrips.length} chuyến)
            </h2>
            {isLoading && <div className="text-center p-10">Loading...</div>}
            {!isLoading && filteredTrips.length > 0 ? (
              <div className="space-y-4">
                {filteredTrips.map((trip) => (
                  <TripCard
                    key={trip.id}
                    trip={trip}
                    onSelect={() => handleSelectTrip(trip)}
                  />
                ))}
              </div>
            ) : (
              !isLoading && (
                <p>Không tìm thấy chuyến xe nào cho tuyến và ngày đã chọn.</p>
              )
            )}
          </div>
        </div>
      )}

      {/* Section 3: Booking Details */}
      {view === "BOOKING" && selectedTrip && (
        <div className="bg-white p-6 rounded-lg shadow-md">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Bước 3: Đặt vé</h2>
            <button
              onClick={resetToSearch}
              className="text-sm font-medium text-blue-600 hover:underline"
            >
              ← Quay về
            </button>
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Seat Map & Location */}
            <div className="bg-white rounded-lg p-6 border border-slate-200 shadow-sm">
              <h3 className="text-xl font-medium">
                {selectedTrip.departureStation} → {selectedTrip.arrivalStation}
              </h3>
              <div className="flex flex-row gap-8 justify-center mt-5">
                {["A", "B"].map((row) => (
                  <div key={row} className="grid grid-cols-3 gap-x-10 gap-y-2">
                    {generateAllSeats(row).map((seat) => (
                      <div
                        key={seat}
                        className={`mt-1 text-center relative flex justify-center ${
                          !isSeatAvailable(seat)
                            ? "cursor-not-allowed"
                            : "cursor-pointer"
                        } ${
                          errors.seats && !isSeatSelected(seat)
                            ? "border-2 border-red-500 rounded"
                            : ""
                        }`}
                        onClick={() =>
                          isSeatAvailable(seat) &&
                          setBookingDetails((prev) => {
                            const isSelected = prev.seats.includes(seat);
                            const newSeats = isSelected
                              ? prev.seats.filter((s) => s !== seat)
                              : [...prev.seats, seat];
                            const newSeatIds = isSelected
                              ? prev.seatIds.filter(
                                  (id) =>
                                    id !==
                                    selectedTrip.details.seatMapping[seat]
                                )
                              : [
                                  ...prev.seatIds,
                                  selectedTrip.details.seatMapping[seat],
                                ];
                            return {
                              ...prev,
                              seats: newSeats,
                              seatIds: newSeatIds,
                            };
                          })
                        }
                      >
                        <img
                          width="32"
                          src={getSeatImage(seat)}
                          alt="seat icon"
                        />
                        <span
                          className={`absolute text-sm font-semibold top-1 ${
                            isSeatSelected(seat)
                              ? "text-red-400"
                              : "text-blue-400"
                          }`}
                        >
                          {seat}
                        </span>
                      </div>
                    ))}
                  </div>
                ))}
              </div>
              {errors.seats && (
                <p className="text-red-500 text-xs mt-2">{errors.seats}</p>
              )}
              <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Điểm đón
                  </label>
                  <select
                    value={bookingDetails.pickupId}
                    onChange={(e) =>
                      setBookingDetails((p) => ({
                        ...p,
                        pickupId: e.target.value,
                      }))
                    }
                    className={`mt-1 block w-full p-2 border rounded-md ${
                      errors.pickupId ? "border-red-500" : "border-gray-300"
                    }`}
                  >
                    <option value="">Chọn điểm</option>
                    {selectedTrip.details.pickupLocations.map((loc) => (
                      <option key={loc.id} value={loc.id}>
                        {loc.name}
                      </option>
                    ))}
                  </select>
                  {errors.pickupId && (
                    <p className="text-red-500 text-xs mt-1">
                      {errors.pickupId}
                    </p>
                  )}
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Điểm xuống
                  </label>
                  <select
                    value={bookingDetails.dropoffId}
                    onChange={(e) =>
                      setBookingDetails((p) => ({
                        ...p,
                        dropoffId: e.target.value,
                      }))
                    }
                    className={`mt-1 block w-full p-2 border rounded-md ${
                      errors.dropoffId ? "border-red-500" : "border-gray-300"
                    }`}
                  >
                    <option value="">Chọn điểm</option>
                    {selectedTrip.details.dropoffLocations.map((loc) => (
                      <option key={loc.id} value={loc.id}>
                        {loc.name}
                      </option>
                    ))}
                  </select>
                  {errors.dropoffId && (
                    <p className="text-red-500 text-xs mt-1">
                      {errors.dropoffId}
                    </p>
                  )}
                </div>
              </div>
            </div>
            {/* Customer Info & Confirmation */}
            <div className="space-y-6">
              <div>
                <h3 className="text-lg font-semibold text-gray-800 mb-3">
                  Thông tin khách hàng
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Họ và tên
                    </label>
                    <input
                      type="text"
                      value={customerInfo.name}
                      onChange={(e) =>
                        setCustomerInfo({
                          ...customerInfo,
                          name: e.target.value,
                        })
                      }
                      className={`mt-1 block w-full p-2 border rounded-md ${
                        errors.name ? "border-red-500" : "border-gray-300"
                      }`}
                    />
                    {errors.name && (
                      <p className="text-red-500 text-xs mt-1">{errors.name}</p>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Số điện thoại
                    </label>
                    <input
                      type="text"
                      value={customerInfo.phone}
                      onChange={(e) =>
                        setCustomerInfo({
                          ...customerInfo,
                          phone: e.target.value,
                        })
                      }
                      className={`mt-1 block w-full p-2 border rounded-md ${
                        errors.phone ? "border-red-500" : "border-gray-300"
                      }`}
                    />
                    {errors.phone && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.phone}
                      </p>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Email (không bắt buộc)
                    </label>
                    <input
                      type="email"
                      value={customerInfo.email}
                      onChange={(e) =>
                        setCustomerInfo({
                          ...customerInfo,
                          email: e.target.value,
                        })
                      }
                      className={`mt-1 block w-full p-2 border rounded-md ${
                        errors.email ? "border-red-500" : "border-gray-300"
                      }`}
                    />
                    {errors.email && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.email}
                      </p>
                    )}
                  </div>
                </div>
              </div>
              <div className="p-4 bg-gray-100 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-800 mb-3">
                  Tóm tắt
                </h3>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span>Chỗ ngồi:</span>{" "}
                    <span className="font-medium">
                      {bookingDetails.seats.join(", ") || "Chưa chọn"}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span>Tổng tiền:</span>{" "}
                    <span className="font-bold text-xl text-red-600">
                      {(
                        selectedTrip.price * bookingDetails.seats.length
                      ).toLocaleString()}{" "}
                      đ
                    </span>
                  </div>
                </div>
              </div>
              <button
                onClick={handleConfirmBooking}
                disabled={isSubmitting}
                className="w-full bg-green-600 text-white font-bold p-3 rounded-md hover:bg-green-700 disabled:bg-gray-400"
              >
                {isSubmitting ? "Đang xử lý..." : "Xác nhận & Đặt vé"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminBookingPage;
