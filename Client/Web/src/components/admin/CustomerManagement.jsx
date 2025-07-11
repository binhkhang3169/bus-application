/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import { useNavigate } from "react-router-dom";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import CustomerForm from "./modal/CustomerForm";

const CustomerManagement = () => {
  const navigate = useNavigate();
  const [customerAll, setCustomerAll] = useState([]);
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [customerModal, setCustomerModal] = useState(false);
  const [message, setMessage] = useState("");
  const [customerId, setCustomerId] = useState("");

  useEffect(() => {
    getCustomerAll();
  }, []);

  const getCustomerAll = async () => {
    try {
      const res = await axios.get(API_URL + "employee/customer", {
        headers: { Authorization: "Bearer " + sessionStorage.getItem("token") },
      });
      setCustomerAll(res.data.customer);
    } catch (err) {
      if (err.response?.status === 401) navigate("/admin");
    }
  };

  const editBtn = (id) => {
    setCustomerId(id);
    openCustomerModal();
  };
  const deleteBtn = (id) => {
    setCustomerId(id);
    setDeleteModal(true);
  };
  const refresh = () => {
    setCustomerId("");
    getCustomerAll();
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setCustomerId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);
  const openCustomerModal = () => {
    setCustomerId("");
    setCustomerModal(true);
  };
  const closeCustomerModal = () => {
    setCustomerModal(false);
    setCustomerId("");
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <button
            className="mx-2 text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
            onClick={openCustomerModal}
          >
            Thêm khách hàng
          </button>
        </div>
        <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
          <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
              <tr>
                <th scope="col" className="px-2 py-3">
                  Mã
                </th>
                <th scope="col" className="px-2 py-3">
                  Họ
                </th>
                <th scope="col" className="px-2 py-3">
                  Tên
                </th>
                <th scope="col" className="px-2 py-3">
                  Email
                </th>
                <th scope="col" className="px-2 py-3">
                  Số điện thoại
                </th>
                <th scope="col" className="px-2 py-3">
                  Địa chỉ
                </th>
                <th scope="col" className="px-2 py-3">
                  Giới tính
                </th>
                <th scope="col" className="px-2 py-3">
                  Ngày sinh
                </th>
                <th scope="col" className="px-2 py-3">
                  Tùy chọn
                </th>
              </tr>
            </thead>
            <tbody>
              {customerAll &&
                customerAll.map((v, i) => (
                  <tr
                    key={i}
                    className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                  >
                    <td className="px-2 py-4">{v.id}</td>

                    <td className="px-2 py-4 font-medium text-gray-900 dark:text-white">
                      {v.first_name}
                    </td>
                    <td className="px-2 py-4">{v.last_name}</td>
                    <td className="px-2 py-4">{v.email}</td>
                    <td className="px-2 py-4">{v.phone_number}</td>
                    <td className="px-2 py-4">
                      <p className="line-clamp-1">{v.address}</p>
                    </td>
                    <td className="px-2 py-4">
                      {v.gender === 0 ? "Female" : "Male"}
                    </td>
                    <td className="px-2 py-4">{v.date_of_birth}</td>
                    <td className="px-2 py-4">
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
                        Hạn chế
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
          id={customerId}
          func={{
            refresh,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"customer"}
          action={"customer"}
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
      {customerModal && (
        <CustomerForm
          func={{
            closeModal: closeCustomerModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh,
          }}
          customerId={customerId}
        />
      )}
    </div>
  );
};

export default CustomerManagement;
