/** @format */

import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate, useSearchParams } from "react-router-dom";
import dayjs from "dayjs";
import { API_URL } from "../configs/env";

// You would replace this with your actual API URL from env or config

function SearchResult() {
  const [searchParams] = useSearchParams();
  const [trips, setTrips] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const from = searchParams.get("from");
  const fromId = searchParams.get("fromId");
  const to = searchParams.get("to");
  const toId = searchParams.get("toId");
  const fromTime = searchParams.get("fromTime");

  useEffect(() => {
    if (from && fromId && to && toId && fromTime) {
      fetchTrips();
    } else {
      console.log("Missing required search parameters:", {
        from,
        fromId,
        to,
        toId,
        fromTime,
      });
    }
  }, [from, fromId, to, toId, fromTime]);

  const fetchTrips = async () => {
    setLoading(true);

    // Log the request parameters
    console.log("Search Request Parameters:", {
      from,
      fromId,
      to,
      toId,
      fromTime,
    });

    const requestUrl = `${API_URL}/v1/trips/search?from=${encodeURIComponent(
      from
    )}&fromId=${fromId}&fromTime=${fromTime}&to=${encodeURIComponent(
      to
    )}&toId=${toId}`;
    console.log("API Request URL:", requestUrl);

    try {
      console.log("Sending API request...");
      const response = await axios.get(requestUrl);
      console.log("API Response Data:", response.data);
      console.log("API Response Status:", response.status);

      setTrips(response.data);
      setLoading(false);
    } catch (err) {
      console.error("Error fetching trips:", err);
      console.log("Error details:", {
        message: err.message,
        response: err.response
          ? {
              status: err.response.status,
              data: err.response.data,
            }
          : "No response data",
      });

      setError("Có lỗi xảy ra khi tìm kiếm chuyến xe. Vui lòng thử lại sau.");
      setLoading(false);

      // For development/demo purposes, use mock data
      const mockData = [
        {
          status: "AVAILABLE",
          price: 900000,
          estimatedTime: "28h",
          stock: 10,
          departureDate: fromTime,
          departureTime: "08:00:00",
          arrivalDate: dayjs(fromTime).add(1, "day").format("YYYY-MM-DD"),
          arrivalTime: "12:00:00",
          license: "51A-00003",
          tripId: 11,
          vehicleId: "3",
          vehicleType: "Limousine",
          estimatedDistance: "1700 km",
          departureStation: "BX Mỹ Đình",
          arrivalStation: "BX Miền Đông",
          fullRoute: "BX Mỹ Đình → BX Vinh → BX Miền Đông",
        },
      ];
      console.log("Using mock data:", mockData);
      setTrips(mockData);
      setLoading(false);
      setError(null); // Clear error when using mock data
    }
  };

  const chooseTrip = (tripId) => {
    console.log("Selected trip ID:", tripId);
    navigate(`/bookingticket?id=${tripId}`);
  };

  // Format currency number to VND
  const formatPrice = (price) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  // Combine date and time for display
  const formatDateTime = (date, time) => {
    return `${dayjs(date).format("DD/MM/YYYY")} ${time}`;
  };

  return (
    <div className="container mx-auto py-8 px-4">
      {from && to && (
        <>
          <h3 className="text-center text-2xl font-bold text-blue-700 mb-5">
            KẾT QUẢ TÌM KIẾM ({trips.length}) : {from} -&gt; {to} (
            {dayjs(fromTime).format("DD/MM/YYYY")})
          </h3>
          <hr className="w-4/5 mx-auto h-0.5 bg-gray-200 md:w-2/5 xl:w-1/5 mb-8" />
        </>
      )}

      {loading && (
        <div className="text-center py-10">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-solid border-green-500 border-r-transparent"></div>
          <p className="mt-4">Đang tìm chuyến xe...</p>
        </div>
      )}

      {error && !loading && (
        <div className="text-center py-10 text-red-600">
          <p>{error}</p>
        </div>
      )}

      {!loading && trips.length === 0 && !error && (
        <div className="text-center py-10">
          <p>
            Không tìm thấy chuyến xe phù hợp. Vui lòng thử lại với tuyến đường
            khác hoặc ngày khác.
          </p>
        </div>
      )}

      {!loading && trips.length > 0 && (
        <div className="space-y-6">
          {trips.map((trip, index) => (
            <div
              key={index}
              className="border border-slate-300 hover:border-orange-500 rounded-xl p-5 hover:shadow-2xl transition-all"
            >
              <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                {/* Departure Information */}
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-lg font-semibold">Thời gian đi</span>
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth="1.5"
                      stroke="currentColor"
                      className="w-5 h-5 hidden md:block"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                      />
                    </svg>
                  </div>
                  <p className="text-xl">
                    {formatDateTime(trip.departureDate, trip.departureTime)}
                  </p>
                  <p className="text-gray-600">{trip.departureStation}</p>
                </div>

                {/* Arrival Information */}
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-lg font-semibold">Thời gian đến</span>
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth="1.5"
                      stroke="currentColor"
                      className="w-5 h-5 hidden md:block"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
                      />
                    </svg>
                  </div>
                  <p className="text-xl">
                    {formatDateTime(trip.arrivalDate, trip.arrivalTime)}
                  </p>
                  <p className="text-gray-600">{trip.arrivalStation}</p>
                </div>

                {/* Trip Details */}
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-lg font-semibold">Thông tin</span>
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth="1.5"
                      stroke="currentColor"
                      className="w-5 h-5 hidden md:block"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z"
                      />
                    </svg>
                  </div>
                  <div className="space-y-1">
                    <p>
                      <span className="font-medium">Xe:</span>{" "}
                      {trip.vehicleType}
                    </p>
                    <p>
                      <span className="font-medium">Thời gian:</span>{" "}
                      {trip.estimatedTime}
                    </p>
                    <p>
                      <span className="font-medium">Khoảng cách:</span>{" "}
                      {trip.estimatedDistance}
                    </p>
                    <p>
                      <span className="font-medium">Ghế trống:</span>{" "}
                      {trip.stock}
                    </p>
                  </div>
                </div>

                {/* Price and Book Button */}
                <div className="flex flex-col items-center justify-center">
                  <p className="text-orange-500 text-xl font-semibold mb-3">
                    {formatPrice(trip.price)}
                  </p>
                  <button
                    className="font-semibold text-white hover:bg-red-600 transition-all border border-transparent bg-red-500 px-6 py-3 rounded-full"
                    onClick={() => chooseTrip(trip.tripId)}
                  >
                    Chọn chuyến
                  </button>
                </div>
              </div>

              {/* Full Route */}
              <div className="mt-4 pt-4 border-t border-gray-200">
                <p className="text-gray-700">
                  <span className="font-medium">Lộ trình: </span>
                  {trip.fullRoute}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default SearchResult;
