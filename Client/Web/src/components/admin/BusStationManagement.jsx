/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import BusStationForm from "./modal/BusStationForm";
import { useNavigate } from "react-router-dom";

const BusStationManagement = () => {
  const navigate = useNavigate();
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [busStationModal, setBusStationModal] = useState(false);
  const [message, setMessage] = useState("");
  const [stationId, setStationId] = useState("");
  const [stations, setStations] = useState([]);
  const [provinces, setProvinces] = useState([]);

  useEffect(() => {
    getStations();
    getProvinces();
  }, []);

  const getProvinces = async () => {
    try {
      const res = await axios.get(API_URL + "api/v1/provinces");
      setProvinces(res.data?.data || []);
    } catch (err) {
      console.error("Error fetching provinces:", err);
      setProvinces([]);
    }
  };

  const getStations = async () => {
    try {
      const res = await axios.get(API_URL + "api/v1/stations");
      setStations(res.data?.data || []);
    } catch (err) {
      console.error("Error fetching stations:", err);
      setStations([]);
    }
  };

  const editBtn = (id) => {
    setStationId(id);
    openBusStationModal();
  };

  const deleteBtn = (id) => {
    setStationId(id);
    setDeleteModal(true);
  };

  const refresh = () => {
    getStations();
    setStationId("");
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setStationId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);
  const openBusStationModal = () => setBusStationModal(true);
  const closeBusStationModal = () => {
    setBusStationModal(false);
    setStationId("");
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <div className="flex items-center gap-2">
            <button
              className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
              onClick={() => {
                setStationId("");
                openBusStationModal();
              }}
            >
              Thêm trạm xe
            </button>
            {/* <button
              className="text-white bg-yellow-500 hover:bg-yellow-600 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-yellow-600 dark:hover:bg-yellow-700"
              onClick={() => navigate("bin")}
            >
              Bin
            </button> */}
          </div>
        </div>
        <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
          <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
              <tr>
                <th scope="col" className="px-6 py-3">
                  Mã
                </th>
                <th scope="col" className="px-6 py-3">
                  Tên trạm
                </th>
                <th scope="col" className="px-6 py-3">
                  Địa chỉ cụ thể
                </th>
                <th scope="col" className="px-6 py-3">
                  Tỉnh/thành phố
                </th>
                <th scope="col" className="px-6 py-3">
                  Tùy chọn
                </th>
              </tr>
            </thead>
            <tbody>
              {stations.length > 0 ? (
                stations.map((station, i) => (
                  <tr
                    key={station.id || i}
                    className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                  >
                    <td className="px-6 py-4">{station.id}</td>

                    <th
                      scope="row"
                      className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                    >
                      {station.name}
                    </th>
                    <td className="px-6 py-4">{station.address}</td>
                    <td className="px-6 py-4">
                      {station.province?.name || "N/A"}
                    </td>
                    <td className="px-6 py-4">
                      <button
                        onClick={() => editBtn(station.id)}
                        className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                      >
                        Sửa
                      </button>
                      <button
                        onClick={() => deleteBtn(station.id)}
                        className="font-medium text-red-600 dark:text-red-500 hover:underline"
                      >
                        Xóa
                      </button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan="4"
                    className="text-center px-6 py-4 text-gray-500 dark:text-gray-400"
                  >
                    No stations found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
      {deleteModal && (
        <WarningNotification
          id={stationId}
          func={{
            refresh,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"bus station"}
          action={"stations"}
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
      {busStationModal && (
        <BusStationForm
          func={{
            closeModal: closeBusStationModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh,
          }}
          busStationId={stationId}
          provinces={provinces}
        />
      )}
    </div>
  );
};

export default BusStationManagement;
