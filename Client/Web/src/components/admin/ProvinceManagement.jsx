/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import ProvinceForm from "./modal/ProvinceForm";
import { useNavigate } from "react-router-dom";

const ProvinceManagement = () => {
  const navigate = useNavigate();
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [provinceFormModal, setProvinceFormModal] = useState(false);
  const [message, setMessage] = useState("");
  const [currentProvinceId, setCurrentProvinceId] = useState("");
  const [provinces, setProvinces] = useState([]);

  useEffect(() => {
    getAllProvinces();
  }, []);

  const getAllProvinces = async () => {
    try {
      const res = await axios.get(API_URL + "api/v1/provinces");
      setProvinces(res.data?.data || []);
    } catch (err) {
      console.error("Error fetching provinces:", err);
      setProvinces([]);
    }
  };

  const editBtn = (id) => {
    setCurrentProvinceId(id);
    setProvinceFormModal(true);
  };
  const deleteBtn = (id) => {
    setCurrentProvinceId(id);
    setDeleteModal(true);
  };
  const refreshData = () => {
    getAllProvinces();
    setCurrentProvinceId("");
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setCurrentProvinceId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);
  const openProvinceForm = () => {
    setCurrentProvinceId("");
    setProvinceFormModal(true);
  };
  const closeProvinceForm = () => {
    setProvinceFormModal(false);
    setCurrentProvinceId("");
  };

  const formatDate = (dateString) => {
    if (!dateString) return "N/A";
    const date = new Date(dateString);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0"); // Tháng bắt đầu từ 0
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <div className="flex items-center gap-2">
            <button
              className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
              onClick={openProvinceForm}
            >
              Thêm tỉnh/thành phố
            </button>
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
                  Tên
                </th>
                <th scope="col" className="px-6 py-3">
                  Trạng thái
                </th>
                <th scope="col" className="px-6 py-3">
                  Ngày tạo
                </th>
                <th scope="col" className="px-6 py-3">
                  Tùy chọn
                </th>
              </tr>
            </thead>
            <tbody>
              {provinces.length > 0 ? (
                provinces.map((province) => (
                  <tr
                    key={province.id}
                    className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                  >
                    <td className="px-6 py-4">{province.id}</td>
                    <th
                      scope="row"
                      className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                    >
                      {province.name}
                    </th>
                    <td className="px-6 py-4">
                      <span
                        className={`px-2 py-1 text-xs font-semibold rounded-full ${
                          province.status === 1 || province.status === true
                            ? "bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300"
                            : "bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-400"
                        }`}
                      >
                        {province.status === 1 || province.status === true
                          ? "Hoạt động"
                          : "Ngừng hoạt động"}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {formatDate(province.createdAt)}
                    </td>
                    <td className="px-6 py-4">
                      <button
                        onClick={() => editBtn(province.id)}
                        className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                      >
                        Sửa
                      </button>
                      <button
                        onClick={() => deleteBtn(province.id)}
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
                    colSpan="5"
                    className="text-center px-6 py-4 text-gray-500 dark:text-gray-400"
                  >
                    No provinces found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {deleteModal && (
        <WarningNotification
          id={currentProvinceId}
          func={{
            refresh: refreshData,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"province"}
          action={"provinces"}
          method={"DELETE"}
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
      {provinceFormModal && (
        <ProvinceForm
          func={{
            closeModal: closeProvinceForm,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh: refreshData,
          }}
          provinceId={currentProvinceId}
        />
      )}
    </div>
  );
};

export default ProvinceManagement;
