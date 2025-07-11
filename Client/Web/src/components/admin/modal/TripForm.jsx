import React, { useEffect, useState } from "react";
import axios from "axios";
import { API_URL } from "../../../configs/env"; // Assuming this configuration is correct
import { useNavigate } from "react-router-dom";

const TripForm = ({
  func,
  tripId,
  routeAll,
  busAll,
  specialAll,
  driverAll,
}) => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [isPickupsLoading, setIsPickupsLoading] = useState(false); // Loading state for pickup points

  // Sample data, which you can pass via props if needed
  const vehicleTypesAll = [
    { id: 1, name: "Ghế ngồi" },
    { id: 2, name: "Giường nằm" },
    { id: 3, name: "Limousine" },
  ];

  const statusTextToValue = (textStatus) => {
    switch (
      String(textStatus).toUpperCase() // Convert to uppercase to ensure a match
    ) {
      case "AVAILABLE":
        return "1";
      case "UNAVAILABLE":
        return "0";
      case "CANCELLED":
        return "2";
      case "COMPLETED":
        return "3";
      default:
        // If the backend returns a number (as a number or string), keep it
        if (!isNaN(textStatus) && textStatus !== null && textStatus !== "") {
          return String(textStatus);
        }
        return "1"; // Default value if no match, e.g., "AVAILABLE"
    }
  };

  // Form input states
  const [status, setStatus] = useState(statusTextToValue("AVAILABLE"));
  const [route, setRoute] = useState("");
  const [selectedVehicleType, setSelectedVehicleType] = useState("");
  const [vehicle, setVehicle] = useState("");
  const [special, setSpecial] = useState("");
  const [departureDate, setDepartureDate] = useState("");
  const [departureTime, setDepartureTime] = useState("");
  const [arrivalDate, setArrivalDate] = useState("");
  const [arrivalTime, setArrivalTime] = useState("");
  const [total, setTotal] = useState(0);
  const [stock, setStock] = useState(0);
  const [driver, setDriver] = useState("");

  // States for pickup points
  const [pickupPoints, setPickupPoints] = useState([]);
  const [selectedPickupPoint, setSelectedPickupPoint] = useState(""); // Stores the ID of the selected pickup point

  // Selected route and vehicle details
  const [selectedRoute, setSelectedRoute] = useState(null);
  const [selectedVehicle, setSelectedVehicle] = useState(null);
  const [filteredVehicles, setFilteredVehicles] = useState([]);

  // Form validation
  const [formErrors, setFormErrors] = useState({});

  // Get trip data if tripId is provided (edit mode)
  useEffect(() => {
    if (tripId) {
      fetchTripDetails();
    } else {
      // Set default dates and times for new trips
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      setDepartureDate(tomorrow.toISOString().split("T")[0]);
      setDepartureTime("08:00");
      setStatus("AVAILABLE");
      setSelectedVehicleType(""); // Ensure type is reset for new form
      setVehicle(""); // Ensure vehicle is reset
      setFilteredVehicles([]); // Ensure filtered list is reset
      setTotal(0); // Reset total seats
      setStock(0); // Reset stock
    }
  }, [tripId]);

  // Fetch trip details for editing
  const fetchTripDetails = async () => {
    setIsLoading(true);
    try {
      const token = sessionStorage.getItem("token");
      const response = await axios.get(`${API_URL}api/v1/trips/${tripId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (response.data) {
        const tripData = response.data.data;
        setFormValues(tripData);
      }
    } catch (error) {
      if (error.response && error.response.status === 401) {
        navigate("/admin");
      } else {
        func.setMessage("Failed to load trip details");
        func.openFailureModal();
      }
    } finally {
      setIsLoading(false);
    }
  };

  // Set form values from trip data
  const setFormValues = (tripData) => {
    console.log(tripData.status);
    setStatus(
      statusTextToValue(
        tripData.status !== undefined && tripData.status !== null
          ? tripData.status
          : "AVAILABLE"
      )
    );
    setRoute(tripData.route?.id || "");

    // Set the value for the pickup point when editing. The useEffect for the route will fetch the list,
    // and this value will be used to select the correct option.
    if (tripData.pickupId) {
      setSelectedPickupPoint(tripData.pickupId);
    }

    if (tripData.vehicle?.type?.id) {
      setSelectedVehicleType(String(tripData.vehicle.type.id));
    } else {
      setSelectedVehicleType("");
    }
    setVehicle(tripData.vehicle?.id || "");

    setSpecial(tripData.special?.id || "");
    setDepartureDate(tripData.departureDate || "");
    setArrivalDate(tripData.arrivalDate || "");
    setTotal(tripData.total || 0);
    setStock(tripData.stock || 0);

    if (tripData.departureTime) {
      const timeValue = tripData.departureTime.substring(0, 5);
      setDepartureTime(timeValue);
    }

    if (tripData.arrivalTime) {
      const timeValue = tripData.arrivalTime.substring(0, 5);
      setArrivalTime(timeValue);
    }

    if (tripData.route?.id) {
      const foundRoute = routeAll.find((r) => r.id === tripData.route.id);
      if (foundRoute) setSelectedRoute(foundRoute);
    }
  };

  // When route is selected, update selectedRoute and fetch pickups
  useEffect(() => {
    const fetchPickups = async (routeId) => {
      setIsPickupsLoading(true);
      setPickupPoints([]);
      setSelectedPickupPoint(""); // Reset pickup point selection when changing route
      setFormErrors((prev) => ({ ...prev, selectedPickupPoint: undefined }));

      try {
        // Temporarily hold the old pickupId to re-set it after fetching if in edit mode
        const currentPickupId = selectedPickupPoint;

        const token = sessionStorage.getItem("adminAccessToken"); // Or your preferred token storage
        const response = await axios.get(
          `${API_URL}api/v1/pickups/byRoute/${routeId}`,
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );

        if (response.data && response.data.data) {
          const fetchedPickups = response.data.data;
          setPickupPoints(fetchedPickups);
          // If in edit mode, check if the old pickupId exists in the new list
          if (
            tripId &&
            fetchedPickups.some((p) => p.id.toString() === currentPickupId)
          ) {
            setSelectedPickupPoint(currentPickupId);
          }
        }
      } catch (error) {
        console.error("Failed to fetch pickup points:", error);
        func.setMessage("Không thể tải danh sách điểm đón.");
        func.openFailureModal();
      } finally {
        setIsPickupsLoading(false);
      }
    };

    if (route) {
      const foundRoute = routeAll.find((r) => r.id === parseInt(route));
      setSelectedRoute(foundRoute || null);
      fetchPickups(route);
    } else {
      setSelectedRoute(null);
      setPickupPoints([]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [route, routeAll, tripId]); // Add tripId to handle edit mode logic

  // Filter vehicles when selectedVehicleType or busAll changes
  useEffect(() => {
    if (selectedVehicleType && busAll && busAll.length > 0) {
      const filtered = busAll.filter(
        (bus) => bus.type && bus.type.id === parseInt(selectedVehicleType)
      );
      setFilteredVehicles(filtered);
    } else {
      setFilteredVehicles([]);
    }
  }, [selectedVehicleType, busAll]);

  // When vehicle (ID) is selected or changes, update total seats
  useEffect(() => {
    if (vehicle) {
      const foundVehicle = busAll.find((b) => b.id === parseInt(vehicle));
      setSelectedVehicle(foundVehicle || null);
      if (foundVehicle && typeof foundVehicle.seat_number !== "undefined") {
        setTotal(foundVehicle.seat_number);
      }
    } else {
      setSelectedVehicle(null);
      setTotal(0);
    }
  }, [vehicle, busAll]);

  // Handler for vehicle type change
  const handleVehicleTypeChange = (e) => {
    const newTypeId = e.target.value;
    setSelectedVehicleType(newTypeId);
    setVehicle("");
    setFormErrors((prev) => ({ ...prev, vehicle: undefined }));
  };

  // Auto-calculate arrival time
  useEffect(() => {
    if (departureDate && departureTime && selectedRoute?.estimatedTime) {
      const departureDateTime = new Date(
        `${departureDate}T${departureTime}:00`
      );
      let hoursToAdd = 0;
      let minutesToAdd = 0;
      const estimatedTime = selectedRoute.estimatedTime;

      if (estimatedTime.includes("h")) {
        const parts = estimatedTime.split("h");
        hoursToAdd = parseInt(parts[0].trim(), 10) || 0;
        if (parts.length > 1 && parts[1].includes("m")) {
          minutesToAdd = parseInt(parts[1].replace("m", "").trim(), 10) || 0;
        }
      } else if (estimatedTime.includes("m")) {
        minutesToAdd = parseInt(estimatedTime.replace("m", "").trim(), 10) || 0;
      }

      const arrivalDateTime = new Date(departureDateTime);
      arrivalDateTime.setHours(arrivalDateTime.getHours() + hoursToAdd);
      arrivalDateTime.setMinutes(arrivalDateTime.getMinutes() + minutesToAdd);

      setArrivalDate(arrivalDateTime.toISOString().split("T")[0]);
      setArrivalTime(
        `${String(arrivalDateTime.getHours()).padStart(2, "0")}:${String(
          arrivalDateTime.getMinutes()
        ).padStart(2, "0")}`
      );
    }
  }, [departureDate, departureTime, selectedRoute]);

  // Auto-update available seats for new trips
  useEffect(() => {
    if (!tripId) {
      setStock(total);
    }
  }, [total, tripId]);

  // Validate form
  const validateForm = () => {
    const errors = {};

    if (!status) errors.status = "Trạng thái là bắt buộc";
    if (!route) errors.route = "Tuyến đường là bắt buộc";

    // Validate pickup point if there is a list to choose from
    if (pickupPoints.length > 0 && !selectedPickupPoint) {
      errors.selectedPickupPoint = "Vui lòng chọn điểm đón.";
    }

    if (!selectedVehicleType && vehicleTypesAll && vehicleTypesAll.length > 0) {
      errors.selectedVehicleType = "Loại xe là bắt buộc";
    }
    if (!vehicle) errors.vehicle = "Xe là bắt buộc";
    if (driverAll && driverAll.length > 0 && !driver) {
      errors.driver = "Tài xế là bắt buộc.";
    }
    if (!departureDate) errors.departureDate = "Ngày đi là bắt buộc";
    if (!departureTime) errors.departureTime = "Giờ đi là bắt buộc";

    if (total > 0 && stock > total)
      errors.stock = "Số ghế trống không thể lớn hơn tổng số ghế";

    // ... (other validations can be added here)

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (isUpdate) => {
    if (!validateForm()) return;

    setIsLoading(true);
    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const tripDataPayload = {
        status: statusTextToValue(status),
        departureDate: departureDate,
        departureTime: `${departureTime}:00`,
        arrivalDate: arrivalDate,
        arrivalTime: `${arrivalTime}:00`,
        total: parseInt(total),
        stock: parseInt(stock),
        route: { id: parseInt(route) },
        vehicle: { id: parseInt(vehicle) },
        special: special ? { id: parseInt(special) } : null,
        driverId: driver ? parseInt(driver) : null,
        pickupId: selectedPickupPoint || null, // Add pickupId to payload
      };

      console.log("Submitting payload:", tripDataPayload);

      let response;
      if (isUpdate) {
        response = await axios.put(
          `${API_URL}api/v1/trips/${tripId}`,
          tripDataPayload,
          { headers: { Authorization: `Bearer ${token}` } }
        );
      } else {
        response = await axios.post(`${API_URL}api/v1/trips`, tripDataPayload, {
          headers: { Authorization: `Bearer ${token}` },
        });
      }

      func.setMessage(
        response.data.message ||
          `Chuyến đi đã được ${isUpdate ? "cập nhật" : "tạo"} thành công`
      );
      func.openSuccessModal();
      func.refresh();
      func.closeModal();
    } catch (error) {
      if (error.response && error.response.status === 401) {
        navigate("/admin");
      } else {
        const errorMsg =
          error.response?.data?.message ||
          `Thất bại khi ${isUpdate ? "cập nhật" : "tạo"} chuyến đi`;
        func.setMessage(errorMsg);
        func.openFailureModal();
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateTrip = () => handleSubmit(false);
  const handleUpdateTrip = () => handleSubmit(true);

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex items-center justify-center">
      <div className="relative w-full max-w-lg bg-white rounded-lg shadow-lg max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t sticky top-0 bg-white z-10">
          <h3 className="text-xl font-semibold text-gray-900">
            {tripId ? "Cập nhật chuyến đi" : "Tạo chuyến đi mới"}
          </h3>
          <button
            onClick={func.closeModal}
            type="button"
            className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 inline-flex justify-center items-center"
            aria-label="Close"
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
        <div className="p-4 md:p-5">
          {isLoading && !tripId ? (
            <div className="flex justify-center py-4">
              <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500"></div>
            </div>
          ) : (
            <form className="space-y-4">
              {/* Status selection */}
              <div>
                <label
                  htmlFor="status"
                  className="block mb-2 text-sm font-medium text-gray-900"
                >
                  Trạng thái <span className="text-red-500">*</span>
                </label>
                <select
                  id="status"
                  value={status}
                  onChange={(e) => setStatus(e.target.value)}
                  className={`bg-gray-50 border ${
                    formErrors.status ? "border-red-500" : "border-gray-300"
                  } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                >
                  {tripId ? (
                    <>
                      <option value="1">Có sẵn (AVAILABLE)</option>
                      <option value="0">Không có sẵn (UNAVAILABLE)</option>
                      <option value="2">Đã hủy (CANCELLED)</option>
                      <option value="3">Đã hoàn thành (COMPLETED)</option>
                    </>
                  ) : (
                    <option value="1">Có sẵn (AVAILABLE)</option>
                  )}
                </select>

                {formErrors.status && (
                  <p className="text-red-500 text-xs mt-1">
                    {formErrors.status}
                  </p>
                )}
              </div>

              {/* Route selection */}
              <div>
                <label
                  htmlFor="route"
                  className="block mb-2 text-sm font-medium text-gray-900"
                >
                  Tuyến đường <span className="text-red-500">*</span>
                </label>
                <select
                  id="route"
                  value={route}
                  onChange={(e) => setRoute(e.target.value)}
                  className={`bg-gray-50 border ${
                    formErrors.route ? "border-red-500" : "border-gray-300"
                  } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                >
                  <option value="">Chọn một tuyến đường</option>
                  {routeAll &&
                    routeAll.map((routeItem) => (
                      <option key={routeItem.id} value={routeItem.id}>
                        {routeItem.name
                          ? `${routeItem.name} | ${
                              routeItem.start_address?.name ||
                              routeItem.start?.name
                            } → ${
                              routeItem.end_address?.name || routeItem.end?.name
                            }`
                          : `${
                              routeItem.start_address?.name ||
                              routeItem.start?.name
                            } → ${
                              routeItem.end_address?.name || routeItem.end?.name
                            }`}
                      </option>
                    ))}
                </select>
                {formErrors.route && (
                  <p className="text-red-500 text-xs mt-1">
                    {formErrors.route}
                  </p>
                )}
                {selectedRoute && (
                  <p className="text-blue-600 text-xs mt-1">
                    Khoảng cách: {selectedRoute.distance} | Thời gian dự kiến:{" "}
                    {selectedRoute.estimatedTime}
                  </p>
                )}
              </div>

              {/* Pickup Point Selection */}
              <div>
                <label
                  htmlFor="pickupPoint"
                  className="block mb-2 text-sm font-medium text-gray-900"
                >
                  Điểm đón{" "}
                  {pickupPoints.length > 0 && (
                    <span className="text-red-500">*</span>
                  )}
                </label>
                <select
                  id="pickupPoint"
                  value={selectedPickupPoint}
                  onChange={(e) => setSelectedPickupPoint(e.target.value)}
                  disabled={
                    !route || isPickupsLoading || pickupPoints.length === 0
                  }
                  className={`bg-gray-50 border ${
                    formErrors.selectedPickupPoint
                      ? "border-red-500"
                      : "border-gray-300"
                  } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 ${
                    !route || isPickupsLoading
                      ? "opacity-50 cursor-not-allowed"
                      : ""
                  }`}
                >
                  <option value="">
                    {isPickupsLoading
                      ? "Đang tải điểm đón..."
                      : !route
                      ? "Vui lòng chọn tuyến đường trước"
                      : pickupPoints.length === 0
                      ? "Không có điểm đón cho tuyến này"
                      : "Chọn một điểm đón"}
                  </option>
                  {pickupPoints.map((pickup) => (
                    <option key={pickup.id} value={pickup.id}>
                      {pickup.station.name} ({pickup.station.address})
                    </option>
                  ))}
                </select>
                {formErrors.selectedPickupPoint && (
                  <p className="text-red-500 text-xs mt-1">
                    {formErrors.selectedPickupPoint}
                  </p>
                )}
              </div>

              {/* Vehicle Type selection */}
              {vehicleTypesAll && vehicleTypesAll.length > 0 && (
                <div>
                  <label
                    htmlFor="vehicleType"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Loại xe <span className="text-red-500">*</span>
                  </label>
                  <select
                    id="vehicleType"
                    value={selectedVehicleType}
                    onChange={handleVehicleTypeChange}
                    className={`bg-gray-50 border ${
                      formErrors.selectedVehicleType
                        ? "border-red-500"
                        : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                  >
                    <option value="">Chọn một loại xe</option>
                    {vehicleTypesAll.map((type) => (
                      <option key={type.id} value={type.id}>
                        {type.name}
                      </option>
                    ))}
                  </select>
                  {formErrors.selectedVehicleType && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.selectedVehicleType}
                    </p>
                  )}
                </div>
              )}

              {/* Vehicle selection */}
              <div>
                <label
                  htmlFor="vehicle"
                  className="block mb-2 text-sm font-medium text-gray-900"
                >
                  Xe <span className="text-red-500">*</span>
                </label>
                <select
                  id="vehicle"
                  value={vehicle}
                  onChange={(e) => setVehicle(e.target.value)}
                  disabled={
                    !selectedVehicleType || filteredVehicles.length === 0
                  }
                  className={`bg-gray-50 border ${
                    formErrors.vehicle ? "border-red-500" : "border-gray-300"
                  } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 ${
                    !selectedVehicleType || filteredVehicles.length === 0
                      ? "opacity-50 cursor-not-allowed"
                      : ""
                  }`}
                >
                  <option value="">
                    {!selectedVehicleType
                      ? "Chọn loại xe trước"
                      : filteredVehicles.length === 0
                      ? "Không có xe cho loại này"
                      : "Chọn một xe"}
                  </option>
                  {filteredVehicles.map((busItem) => (
                    <option key={busItem.id} value={busItem.id}>
                      {busItem.license} - {busItem.type?.name} (
                      {busItem.seat_number} ghế)
                    </option>
                  ))}
                </select>
                {formErrors.vehicle && (
                  <p className="text-red-500 text-xs mt-1">
                    {formErrors.vehicle}
                  </p>
                )}
              </div>

              {/* Special selection */}
              <div>
                <label
                  htmlFor="special"
                  className="block mb-2 text-sm font-medium text-gray-900"
                >
                  Ưu đãi đặc biệt
                </label>
                <select
                  id="special"
                  value={special}
                  onChange={(e) => setSpecial(e.target.value)}
                  className={`bg-gray-50 border ${
                    formErrors.special ? "border-red-500" : "border-gray-300"
                  } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                >
                  <option value="">Không có</option>
                  {specialAll &&
                    specialAll.map((specialItem) => (
                      <option key={specialItem.id} value={specialItem.id}>
                        {specialItem.name ||
                          specialItem.type ||
                          "Special " + specialItem.id}
                      </option>
                    ))}
                </select>
                {formErrors.special && (
                  <p className="text-red-500 text-xs mt-1">
                    {formErrors.special}
                  </p>
                )}
              </div>

              {/* Driver Selection */}
              {driverAll && driverAll.length > 0 && (
                <div>
                  <label
                    htmlFor="driver"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Tài xế <span className="text-red-500">*</span>
                  </label>
                  <select
                    id="driver"
                    value={driver}
                    onChange={(e) => setDriver(e.target.value)}
                    className={`bg-gray-50 border ${
                      formErrors.driver ? "border-red-500" : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                  >
                    <option value="">Chọn một tài xế</option>
                    {driverAll.map((driverItem) => (
                      <option key={driverItem.id} value={driverItem.id}>
                        {driverItem.details?.fullName ||
                          driverItem.username ||
                          `Driver ID: ${driverItem.id}`}
                      </option>
                    ))}
                  </select>
                  {formErrors.driver && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.driver}
                    </p>
                  )}
                </div>
              )}

              {/* Total seats and Available seats */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label
                    htmlFor="total"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Tổng số ghế <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    id="total"
                    value={total}
                    readOnly
                    className="bg-gray-100 border border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5"
                  />
                  {formErrors.total && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.total}
                    </p>
                  )}
                </div>
                {/* <div>
                  <label
                    htmlFor="stock"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Ghế trống <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    id="stock"
                    value={stock}
                    onChange={(e) =>
                      setStock(Math.max(0, parseInt(e.target.value) || 0))
                    }
                    min="0"
                    max={total > 0 ? total : undefined}
                    className={`bg-gray-50 border ${
                      formErrors.stock ? "border-red-500" : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                  />
                  {formErrors.stock && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.stock}
                    </p>
                  )}
                </div> */}
              </div>

              {/* Departure Date and Time */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label
                    htmlFor="departureDate"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Ngày đi <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="date"
                    id="departureDate"
                    value={departureDate}
                    onChange={(e) => setDepartureDate(e.target.value)}
                    className={`bg-gray-50 border ${
                      formErrors.departureDate
                        ? "border-red-500"
                        : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                    min={new Date().toISOString().split("T")[0]}
                  />
                  {formErrors.departureDate && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.departureDate}
                    </p>
                  )}
                </div>
                <div>
                  <label
                    htmlFor="departureTime"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Giờ đi <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="time"
                    id="departureTime"
                    value={departureTime}
                    onChange={(e) => setDepartureTime(e.target.value)}
                    className={`bg-gray-50 border ${
                      formErrors.departureTime
                        ? "border-red-500"
                        : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5`}
                  />
                  {formErrors.departureTime && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.departureTime}
                    </p>
                  )}
                </div>
              </div>

              {/* Arrival Date and Time */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label
                    htmlFor="arrivalDate"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Ngày đến <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="date"
                    id="arrivalDate"
                    value={arrivalDate}
                    readOnly
                    className="bg-gray-100 border border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5"
                  />
                </div>
                <div>
                  <label
                    htmlFor="arrivalTime"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Giờ đến <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="time"
                    id="arrivalTime"
                    value={arrivalTime}
                    readOnly
                    className="bg-gray-100 border border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5"
                  />
                </div>
              </div>
              <p className="text-gray-500 text-xs -mt-2">
                Tự động tính toán dựa trên giờ đi và thời gian dự kiến của
                tuyến.
              </p>

              {/* Submit button */}
              <div className="flex justify-end pt-2 sticky bottom-0 bg-white z-10 pb-4">
                <button
                  type="button"
                  onClick={() => func.closeModal()}
                  className="mr-2 px-4 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:ring-4 focus:ring-gray-200"
                >
                  Hủy
                </button>
                <button
                  type="button"
                  onClick={tripId ? handleUpdateTrip : handleCreateTrip}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-700 rounded-lg hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 flex items-center"
                  disabled={isLoading}
                >
                  {isLoading && (
                    <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white mr-2"></div>
                  )}
                  {tripId ? "Cập nhật" : "Tạo mới"}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};

export default TripForm;
