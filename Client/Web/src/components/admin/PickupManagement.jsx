/** @format */

import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { API_URL } from "../../configs/env";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import PickupForm from "./modal/PickupForm";

const PickupManagement = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [pickupModal, setPickupModal] = useState(false);
  const [message, setMessage] = useState("");
  const [pickups, setPickups] = useState([]);
  const [stations, setStations] = useState([]);
  const [routes, setRoutes] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        const token = sessionStorage.getItem("adminAccessToken");
        const [pickupRes, stationRes, routeRes] = await Promise.all([
          axios.get(`${API_URL}api/v1/pickups`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
          axios.get(`${API_URL}api/v1/stations`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
          axios.get(`${API_URL}api/v1/routes`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);
        setPickups(pickupRes.data.data || []);
        setStations(stationRes.data.data || []);
        setRoutes(routeRes.data.data || []);
      } catch (error) {
        console.error("Failed to fetch data:", error);
        setMessage("Failed to load data.");
        setFailureModal(true);
      }
      setIsLoading(false);
    };
    fetchData();
  }, []);

  const refreshData = async () => {
    setIsLoading(true);
    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const pickupRes = await axios.get(`${API_URL}api/v1/pickups`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setPickups(pickupRes.data.data || []);
    } catch (error) {
      console.error("Failed to refresh pickups:", error);
    }
    setIsLoading(false);
  };

  const openPickupModal = () => setPickupModal(true);
  const closePickupModal = () => setPickupModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const closeSuccessModal = () => setSuccessModal(false);
  const openFailureModal = () => setFailureModal(true);
  const closeFailureModal = () => setFailureModal(false);

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-2 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <button
            className="mr-2 text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
            onClick={openPickupModal}
          >
            Thêm điểm đón trả
          </button>
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
                    Lộ trình
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Trạm
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Thời gian có mặt (Ước lượng)
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Đoạn đường
                  </th>
                </tr>
              </thead>
              <tbody>
                {pickups.length > 0 ? (
                  pickups.map((pickup) => (
                    <tr
                      key={pickup.id}
                      className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                    >
                      <td className="px-6 py-4 font-medium text-gray-900 dark:text-white">
                        {pickup.id}
                      </td>
                      <td className="px-6 py-4">{`${pickup.route.start.name} - ${pickup.route.end.name}`}</td>
                      <td className="px-6 py-4">{pickup.station.name}</td>
                      <td className="px-6 py-4">{pickup.time}</td>
                      <td className="px-6 py-4">
                        {pickup.selfId === "-1"
                          ? "Điểm đầu"
                          : pickup.selfId === "-2"
                          ? "Điểm cuối"
                          : "Giữa đường"}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td
                      colSpan="5"
                      className="text-center py-4 text-gray-500 dark:text-gray-400"
                    >
                      Không tìm thấy điểm đón/trả.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {pickupModal && (
        <PickupForm
          func={{
            closeModal: closePickupModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh: refreshData,
          }}
          stations={stations}
          routes={routes}
          pickups={pickups}
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
    </div>
  );
};

export default PickupManagement;
