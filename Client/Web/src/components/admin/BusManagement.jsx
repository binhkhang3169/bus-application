/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import { useNavigate } from "react-router-dom";
import BusForm from "./modal/BusForm";

const BusManagement = () => {
  const navigate = useNavigate();
  const [busAll, setBusAll] = useState([]);
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [createModal, setCreateModal] = useState(false);
  const [message, setMessage] = useState("");
  const [busId, setBusId] = useState("");

  useEffect(() => {
    getBusAll();
  }, []);

  const getBusAll = async () => {
    try {
      const res = await axios.get(API_URL + "api/v1/vehicles", {
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      });
      setBusAll(res.data.data.filter((v) => v.status === 1));
    } catch (err) {
      console.error("Error fetching buses:", err);
      if (err.response?.status === 401) navigate("/admin");
      else {
        setMessage("Failed to load bus data. Please try again.");
        setFailureModal(true);
      }
    }
  };

  const editBtn = (id) => {
    setBusId(id);
    setCreateModal(true);
  };

  const deleteBtn = (id) => {
    setBusId(id);
    setDeleteModal(true);
  };

  const refresh = () => {
    getBusAll();
    setBusId("");
  };

  const closeCreateModal = () => {
    setCreateModal(false);
    setBusId("");
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setBusId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);

  return (
    <div className="w-full p-2">
      <div className="flex justify-between items-center mb-4">
        <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
        <div className="flex gap-2">
          <button
            className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
            onClick={() => {
              setBusId("");
              setCreateModal(true);
            }}
          >
            Thêm xe mới
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
                Biển số xe
              </th>
              <th scope="col" className="px-6 py-3">
                Số ghế ngồi
              </th>
              <th scope="col" className="px-6 py-3">
                Loại xe
              </th>
              <th scope="col" className="px-6 py-3">
                Tùy chọn
              </th>
            </tr>
          </thead>
          <tbody>
            {busAll.length > 0 ? (
              busAll.map((v) => (
                <tr
                  key={v.id}
                  className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                >
                  <td className="px-6 py-4">{v.id}</td>

                  <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                    {v.license}
                  </td>
                  <td className="px-6 py-4">{v.seat_number}</td>
                  <td className="px-6 py-4">{v.type ? v.type.name : "N/A"}</td>
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
              ))
            ) : (
              <tr>
                <td
                  colSpan="4"
                  className="px-6 py-4 text-center text-gray-500 dark:text-gray-400"
                >
                  No active buses found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {deleteModal && (
        <WarningNotification
          id={busId}
          func={{
            refresh,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"bus"}
          action={"delete"}
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
      {createModal && (
        <BusForm
          func={{
            closeModal: closeCreateModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh,
          }}
          busId={busId}
        />
      )}
    </div>
  );
};

export default BusManagement;
