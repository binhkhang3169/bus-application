/** @format */

import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import dayjs from "dayjs";
import { API_URL } from "../configs/env";
import api from "../services/apiService";

// ===================================================================================
// COMPONENT: TripCard - Giao diện được thiết kế lại cho mỗi chuyến đi
// ===================================================================================
const TripCard = ({ trip, onSelect, isRoundTrip, selectionStep }) => {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <div className="bg-white rounded-xl shadow-lg hover:shadow-2xl transition-shadow duration-300 overflow-hidden border border-gray-200">
      <div className="p-5 flex flex-col md:flex-row gap-6">
        {/* Cột trái: Thông tin hành trình */}
        <div className="flex-grow">
          <div className="flex items-center gap-4 mb-4">
            {/* <img src="https://storage.googleapis.com/futa-bus-bucket/icon_new/bus_1.png" alt="Vehicle" className="h-10 w-10"/> */}
            <div>
              <p className="font-bold text-lg text-gray-800 capitalize">
                {trip.car_type}
              </p>
              <p className="text-sm text-gray-500">{trip.estimated_distance}</p>
            </div>
          </div>

          <div className="flex gap-4">
            {/* Timeline */}
            <div className="flex flex-col items-center">
              <div className="w-5 h-5 rounded-full border-2 border-green-500 bg-white"></div>
              <div className="flex-1 w-px bg-gray-300 border-l-2 border-dotted border-gray-300 my-2"></div>
              <div className="w-5 h-5 rounded-full bg-green-500"></div>
            </div>
            {/* Thông tin thời gian và địa điểm */}
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

        {/* Cột phải: Giá và hành động */}
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
              {isRoundTrip && selectionStep === "departure"
                ? "Chọn chuyến đi"
                : "Chọn vé"}
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

      {/* Phần chi tiết lộ trình có thể mở rộng */}
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

// ===================================================================================
// COMPONENT: TravelScheduleInfo - Component chính của trang
// ===================================================================================
function TravelScheduleInfo() {
  const location = useLocation();
  const navigate = useNavigate();

  const [departureTrips, setDepartureTrips] = useState([]);
  const [returnTrips, setReturnTrips] = useState([]);
  const [filteredTrips, setFilteredTrips] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchInfo, setSearchInfo] = useState({
    from: "",
    to: "",
    fromId: "",
    toId: "",
    date: "",
    returnDate: "",
  });

  const [isRoundTrip, setIsRoundTrip] = useState(false);
  const [selectionStep, setSelectionStep] = useState("departure"); // 'departure' hoặc 'return'
  const [selectedDepartureTrip, setSelectedDepartureTrip] = useState(null);

  const [filters, setFilters] = useState({ timeOfDay: "", carType: "" });

  useEffect(() => {
    const query = new URLSearchParams(location.search);
    const from = decodeURIComponent(query.get("from") || "");
    const fromId = query.get("fromId") || "";
    const to = decodeURIComponent(query.get("to") || "");
    const toId = query.get("toId") || "";
    const date = query.get("fromTime") || dayjs().format("YYYY-MM-DD");
    const returnDate = query.get("returnTime");

    const newSearchInfo = { from, to, fromId, toId, date, returnDate };
    setSearchInfo(newSearchInfo);

    setIsRoundTrip(!!returnDate);
    setSelectionStep("departure");
    setSelectedDepartureTrip(null);
    setFilters({ timeOfDay: "", carType: "" });

    if (fromId && toId && date) {
      fetchTrips(fromId, toId, from, to, date, setDepartureTrips, "departure");
    } else {
      setError("Thông tin tìm kiếm không hợp lệ. Vui lòng quay lại trang chủ.");
      setIsLoading(false);
    }
  }, [location.search]);

  useEffect(() => {
    const tripsToFilter =
      selectionStep === "departure" ? departureTrips : returnTrips;
    applyFilters(tripsToFilter);
  }, [filters, departureTrips, returnTrips, selectionStep]);

  const fetchTrips = async (fromId, toId, from, to, date, tripSetter) => {
    setIsLoading(true);
    setError(null);
    tripSetter([]);

    try {
      const response = await api.get(
        `${API_URL}api/v1/trips/search?from=${from}&to=${to}&fromId=${fromId}&toId=${toId}&fromTime=${date}`
      );

      if (response.data && response.data.data) {
        const formattedTrips = response.data.data.map((trip) => ({
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
          estimated_time: trip.estimatedTime,
          estimated_distance: trip.estimatedDistance,
        }));
        tripSetter(formattedTrips);
      } else {
        tripSetter([]);
      }
    } catch (err) {
      console.error("Failed to fetch trips:", err);
      setError("Không thể tải dữ liệu chuyến đi. Vui lòng thử lại sau.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSelectTrip = (trip) => {
    if (isRoundTrip && selectionStep === "departure") {
      setSelectedDepartureTrip(trip);
      setSelectionStep("return");
      setFilters({ timeOfDay: "", carType: "" }); // Reset filter cho chuyến về
      fetchTrips(
        searchInfo.toId,
        searchInfo.fromId,
        searchInfo.to,
        searchInfo.from,
        searchInfo.returnDate,
        setReturnTrips
      );
    } else {
      const departureId = isRoundTrip ? selectedDepartureTrip.id : trip.id;
      const returnId = isRoundTrip ? trip.id : null;

      let navUrl = `/dat-ve?departureTripId=${departureId}`;
      if (returnId) {
        navUrl += `&returnTripId=${returnId}`;
      }
      navigate(navUrl);
    }
  };

  const applyFilters = (tripsToFilter) => {
    let filtered = [...tripsToFilter];
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
    // Sắp xếp các chuyến đi theo giờ khởi hành
    setFilteredTrips(
      filtered.sort((a, b) => a.start_time.localeCompare(b.start_time))
    );
  };

  const handleFilterChange = (filterType, value) => {
    setFilters((prev) => ({
      ...prev,
      [filterType]: prev[filterType] === value ? "" : value,
    }));
  };

  const handleResetDeparture = () => {
    setSelectedDepartureTrip(null);
    setSelectionStep("departure");
    setReturnTrips([]);
  };

  const getTitle = () => {
    const formatDate = (dateStr) => dayjs(dateStr).format("DD/MM/YYYY");
    if (selectionStep === "return") {
      return `CHỌN CHUYẾN VỀ: ${searchInfo.to} → ${
        searchInfo.from
      } (${formatDate(searchInfo.returnDate)})`;
    }
    return `Kết quả các chuyến xe đi từ: ${searchInfo.from} → ${
      searchInfo.to
    } (${formatDate(searchInfo.date)})`;
  };

  return (
    <div className="bg-gray-50 min-h-screen">
      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8 mt-10 mb-32 mx-auto max-w-screen-xl px-4">
        {/* Cột bộ lọc */}
        <div className="col-span-1 bg-white rounded-xl shadow-md p-5 border border-gray-200 h-fit">
          <h3 className="text-xl font-bold text-gray-800 mb-5">Bộ lọc</h3>
          <div className="space-y-6">
            <div>
              <label className="block mb-3 font-semibold text-gray-700">
                Khung giờ đi
              </label>
              <div className="grid grid-cols-2 gap-2">
                {["", "sang", "chieu", "toi"].map((val) => {
                  const labels = {
                    "": "Tất cả",
                    sang: "Sáng",
                    chieu: "Chiều",
                    toi: "Tối",
                  };
                  const isSelected = filters.timeOfDay === val;
                  return (
                    <button
                      key={val}
                      onClick={() => handleFilterChange("timeOfDay", val)}
                      className={`py-2 px-3 text-sm rounded-lg border transition-colors ${
                        isSelected
                          ? "bg-blue-600 border-blue-600 text-white font-semibold"
                          : "border-gray-300 hover:bg-gray-100"
                      }`}
                    >
                      {labels[val]}
                    </button>
                  );
                })}
              </div>
            </div>
            <div>
              <label className="block mb-3 font-semibold text-gray-700">
                Loại xe
              </label>
              <div className="flex flex-col space-y-2">
                {["", "ghế ngồi", "giường nằm", "limousine"].map((val) => {
                  const labels = {
                    "": "Tất cả",
                    "ghế ngồi": "Ghế ngồi",
                    "giường nằm": "Giường nằm",
                    limousine: "Limousine",
                  };
                  const isSelected = filters.carType === val;
                  return (
                    <button
                      key={val}
                      onClick={() => handleFilterChange("carType", val)}
                      className={`py-2 px-3 text-sm rounded-lg border transition-colors text-left ${
                        isSelected
                          ? "bg-blue-600 border-blue-600 text-white font-semibold"
                          : "border-gray-300 hover:bg-gray-100"
                      }`}
                    >
                      {labels[val]}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </div>

        {/* Cột kết quả */}
        <div className="col-span-1 lg:col-span-3">
          <h3 className="text-center text-2xl font-bold text-blue-700 mb-5">
            {getTitle()} ({filteredTrips.length})
          </h3>

          {isRoundTrip && selectedDepartureTrip && (
            <div className="bg-blue-50 border-l-4 border-blue-500 text-blue-800 p-4 rounded-r-lg mb-6 shadow-md flex justify-between items-center">
              <div>
                <h4 className="font-bold">Chuyến đi đã chọn:</h4>
                <p className="text-sm">
                  {selectedDepartureTrip.start_address} →{" "}
                  {selectedDepartureTrip.end_address}
                </p>
                <p className="text-sm">
                  Thời gian: {selectedDepartureTrip.start_time} -{" "}
                  {dayjs(selectedDepartureTrip.date).format("DD/MM/YYYY")}
                </p>
              </div>
              <button
                onClick={handleResetDeparture}
                className="font-semibold text-blue-600 hover:text-blue-800 underline"
              >
                Thay đổi
              </button>
            </div>
          )}

          {isLoading ? (
            <div className="text-center py-10">
              <div className="inline-block h-10 w-10 animate-spin rounded-full border-4 border-solid border-green-500 border-r-transparent"></div>
              <p className="mt-3 text-gray-600 font-medium">
                Đang tìm kiếm chuyến đi...
              </p>
            </div>
          ) : error ? (
            <div className="text-center py-10 bg-red-50 text-red-600 rounded-lg shadow">
              {error}
            </div>
          ) : filteredTrips.length === 0 ? (
            <div className="text-center py-10 bg-gray-100 rounded-lg shadow">
              <p className="font-semibold text-gray-700">
                Không tìm thấy chuyến xe phù hợp.
              </p>
              <p className="text-sm text-gray-500 mt-1">
                Vui lòng thử lại với bộ lọc hoặc ngày khác.
              </p>
            </div>
          ) : (
            <div className="space-y-6">
              {filteredTrips.map((trip) => (
                <TripCard
                  key={trip.id}
                  trip={trip}
                  onSelect={() => handleSelectTrip(trip)}
                  isRoundTrip={isRoundTrip}
                  selectionStep={selectionStep}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default TravelScheduleInfo;
