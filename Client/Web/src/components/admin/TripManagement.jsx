import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { API_URL } from "../../configs/env";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import TripForm from "./modal/TripForm";

const TripManagement = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [tripModal, setTripModal] = useState(false);
  const [message, setMessage] = useState("");
  const [tripId, setTripId] = useState("");
  const [routeAll, setRouteAll] = useState([]);
  const [busAll, setBusAll] = useState([]);
  const [driverAll, setDriverAll] = useState([]);
  const [tripAll, setTripAll] = useState([]);
  const [specialDays, setSpecialDays] = useState([]);

  useEffect(() => {
    const fetchAllData = async () => {
      setIsLoading(true);
      try {
        const token = sessionStorage.getItem("adminAccessToken");
        await Promise.all([
          getRouteAll(token),
          getBusAll(token),
          getDriverAll(token),
          getTripAll(token),
          getSpecialDays(token),
        ]);
      } catch (error) {
        console.error("Error fetching initial data:", error);
        if (error.response?.status === 401) navigate("/admin");
      } finally {
        setIsLoading(false);
      }
    };
    fetchAllData();
  }, [navigate]);

  const getRouteAll = async (token) => {
    try {
      const res = await axios.get(`${API_URL}api/v1/routes`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setRouteAll(res.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch routes", err);
    }
  };

  const getBusAll = async (token) => {
    try {
      const res = await axios.get(`${API_URL}api/v1/vehicles`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setBusAll(res.data?.data.filter((bus) => bus.status === 1) || []);
    } catch (err) {
      console.error("Failed to fetch buses", err);
    }
  };

  const getDriverAll = async (token) => {
    try {
      const res = await axios.get(
        `${API_URL}api/v1/users/by-role?roleName=ROLE_DRIVER`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setDriverAll(
        res.data?.data.filter((driver) => driver.active === 1) || []
      );
    } catch (err) {
      console.error("Failed to fetch drivers", err);
    }
  };

  const getTripAll = async (token) => {
    try {
      const res = await axios.get(`${API_URL}api/v1/trips`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const futureTrips = (res.data?.data || []).filter(
        (trip) =>
          new Date(`${trip.arrivalDate}T${trip.arrivalTime}`) > new Date()
      );
      setTripAll(futureTrips);
    } catch (err) {
      console.error("Failed to fetch trips", err);
    }
  };

  const getSpecialDays = async (token) => {
    try {
      const res = await axios.get(`${API_URL}api/v1/special-days`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setSpecialDays(res.data?.data || []);
    } catch (error) {
      console.error("Failed to fetch special days:", error);
    }
  };

  const handleEdit = (id) => {
    setTripId(id);
    setTripModal(true);
  };
  const handleDelete = (id) => {
    setTripId(id);
    setDeleteModal(true);
  };
  const refreshData = () => {
    setIsLoading(true);
    const token = sessionStorage.getItem("adminAccessToken");
    Promise.all([
      getRouteAll(token),
      getBusAll(token),
      getDriverAll(token),
      getTripAll(token),
      getSpecialDays(token),
    ]).finally(() => {
      setIsLoading(false);
      setTripId("");
    });
  };

  const closeDeleteModal = () => {
    setDeleteModal(false);
    setTripId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const closeFailureModal = () => setFailureModal(false);
  const openFailureModal = () => setFailureModal(true);
  const openTripModal = () => {
    setTripId("");
    setTripModal(true);
  };
  const closeTripModal = () => {
    setTripModal(false);
    setTripId("");
  };

  const formatDate = (dateString) => {
    if (!dateString) return "N/A";
    const date = new Date(dateString);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0"); // Tháng bắt đầu từ 0
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  };

  const formatPrice = (price) =>
    new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price || 0);
  const formatTime = (timeString) =>
    timeString ? String(timeString).substring(0, 5) : "N/A";

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-2 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <div className="flex items-center gap-2">
            <button
              onClick={openTripModal}
              className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
            >
              Thêm chuyến
            </button>
            {/* <button
              onClick={() => navigate("bin")}
              className="text-white bg-yellow-500 hover:bg-yellow-600 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-yellow-600 dark:hover:bg-yellow-700"
            >
              View Bin
            </button> */}
          </div>
        </div>
        {isLoading ? (
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        ) : (
          <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
              <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
                <tr>
                  <th scope="col" className="px-6 py-3">
                    Mã
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Chuyến đi
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Ghế
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Giờ đi
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Ngày đi
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Xe
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Tùy chọn
                  </th>
                </tr>
              </thead>
              <tbody>
                {tripAll.length === 0 ? (
                  <tr>
                    <td
                      colSpan="6"
                      className="px-6 py-4 text-center text-gray-500 dark:text-gray-400"
                    >
                      Chưa có chuyến xe
                    </td>
                  </tr>
                ) : (
                  tripAll.map((trip) => (
                    <tr
                      key={trip.id}
                      className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                    >
                      <td className="px-6 py-4">{trip.id}</td>
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                      >
                        {trip.route?.start?.name && trip.route?.end?.name
                          ? `${trip.route.start.name} - ${trip.route.end.name}`
                          : "N/A"}
                      </th>
                      <td className="px-6 py-4">
                          {trip.stock}/ {trip.total ?? trip.vehicle?.seat_number ?? "N/A"}
                      </td>
                      <td className="px-6 py-4">{`${formatTime(
                        trip.departureTime
                      )} - ${formatTime(trip.arrivalTime)}`}</td>
                      <td className="px-6 py-4">
                        {formatDate(trip.departureDate)}
                      </td>
                      <td className="px-6 py-4">
                        {trip.vehicle?.license || "N/A"}
                      </td>
                      <td className="px-6 py-4 flex space-x-2">
                        <button
                          onClick={() => handleEdit(trip.id)}
                          className="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                        >
                          Sửa
                        </button>
                        <button
                          onClick={() => handleDelete(trip.id)}
                          className="font-medium text-red-600 dark:text-red-500 hover:underline"
                        >
                          Xóa
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
      {deleteModal && (
        <WarningNotification
          id={tripId}
          func={{
            refresh: refreshData,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"trip"}
          action={"trips"}
          method={"put"}
        />
      )}
      {successModal && (
        <SuccessNotification
          func={{ closeModal: closeSuccessModal }}
          message={message}
        />
      )}
      {failureModal && (
        <FailureNotification
          func={{ closeModal: closeFailureModal }}
          message={message}
        />
      )}
      {tripModal && (
        <TripForm
          func={{
            closeModal: closeTripModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh: refreshData,
          }}
          tripId={tripId}
          routeAll={routeAll}
          busAll={busAll}
          driverAll={driverAll}
          specialAll={specialDays}
        />
      )}
    </div>
  );
};
export default TripManagement;
