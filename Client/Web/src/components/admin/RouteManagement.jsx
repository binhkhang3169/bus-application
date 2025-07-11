/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import RouteForm from "./modal/RouteForm";
import { useNavigate } from "react-router-dom";

const RouteManagement = () => {
  const navigate = useNavigate();
  const [busStationAll, setBusStationAll] = useState([]);
  const [routeAll, setRouteAll] = useState([]);
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [routeModal, setRouteModal] = useState(false);
  const [message, setMessage] = useState("");
  const [routeId, setRouteId] = useState("");

  useEffect(() => {
    getBusStationAll();
    getRouteAll();
  }, []);

  const getBusStationAll = async () => {
    try {
      const res = await axios.get(API_URL + "api/v1/provinces");
      setBusStationAll(res.data.data);
    } catch (err) {
      console.error("Error fetching bus stations:", err);
    }
  };

  const getRouteAll = async () => {
    try {
      const res = await axios.get(API_URL + "api/v1/routes");
      setRouteAll(res.data.data);
    } catch (err) {
      console.error("Error fetching routes:", err);
    }
  };

  const editBtn = (id) => {
    setRouteId(id);
    openRouteModal();
  };

  const deleteBtn = (id) => {
    setRouteId(id);
    setDeleteModal(true);
  };

  const refresh = () => {
    setRouteId("");
    getRouteAll();
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setRouteId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);
  const openRouteModal = () => setRouteModal(true);
  const closeRouteModal = () => {
    setRouteModal(false);
    setRouteId("");
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <div className="flex items-center gap-2">
            <button
              className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
              onClick={openRouteModal}
            >
              Thêm lộ trình
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
                  Điểm đi
                </th>
                <th scope="col" className="px-6 py-3">
                  Điểm đến
                </th>
                <th scope="col" className="px-6 py-3">
                  Giá
                </th>
                <th scope="col" className="px-6 py-3">
                  Thời gian đi
                </th>
                <th scope="col" className="px-6 py-3">
                  Tùy chọn
                </th>
              </tr>
            </thead>
            <tbody>
              {routeAll &&
                routeAll.map((v, i) => (
                  <tr
                    key={i}
                    className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                  >
                    <th
                      scope="row"
                      className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                    >
                      {v.id}
                    </th>
                    <td className="px-6 py-4">{v.start.name}</td>
                    <td className="px-6 py-4">{v.end.name}</td>
                    <td className="px-6 py-4">
                      {v.price.toLocaleString("vi-VN")} VND
                    </td>
                    <td className="px-6 py-4">{v.estimatedTime}</td>
                    <td className="px-6 py-4">
                      <button
                        onClick={() => editBtn(v.id)}
                        className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                      >
                        Sửa
                      </button>
                      <button
                        onClick={() => deleteBtn(v.id)}
                        className="font-medium text-red-600 dark:text-red-500 hover:underline"
                      >
                        Xóa
                      </button>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </div>
      {deleteModal && (
        <WarningNotification
          id={routeId}
          func={{
            refresh,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"route"}
          action={"route"}
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
      {routeModal && (
        <RouteForm
          func={{
            closeModal: closeRouteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh,
          }}
          routeId={routeId}
          busStationAll={busStationAll}
        />
      )}
    </div>
  );
};

export default RouteManagement;
