/** @format */

import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { API_URL } from "../../configs/env";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import NewsForm from "./modal/NewsForm";

const NewsManagement = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);

  // Modal states
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [newsFormModal, setNewsFormModal] = useState(false);

  // Data and state
  const [message, setMessage] = useState("");
  const [newsId, setNewsId] = useState("");
  const [newsList, setNewsList] = useState([]);

  useEffect(() => {
    fetchNewsList();
  }, []);

  const fetchNewsList = async () => {
    setIsLoading(true);
    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const res = await axios.get(`${API_URL}api/v1/news`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setNewsList(res.data || []);
    } catch (error) {
      console.error("Error fetching news list:", error);
      if (error.response?.status === 401) navigate("/admin");
    } finally {
      setIsLoading(false);
    }
  };

  const handleEdit = (id) => {
    setNewsId(id);
    setNewsFormModal(true);
  };

  const handleDelete = (id) => {
    setNewsId(id);
    setDeleteModal(true);
  };

  const refreshData = () => {
    fetchNewsList();
    setNewsId("");
  };

  const closeDeleteModal = () => setDeleteModal(false);
  const closeSuccessModal = () => setSuccessModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const closeFailureModal = () => setFailureModal(false);
  const openFailureModal = () => setFailureModal(true);

  const openNewsFormModal = () => {
    setNewsId("");
    setNewsFormModal(true);
  };
  const closeNewsFormModal = () => setNewsFormModal(false);

  const formatDate = (dateString) =>
    new Date(dateString).toLocaleDateString("vi-VN");

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          {/* DARK MODE: Thêm text color */}
          <h1 className="ml-2 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <button
            onClick={openNewsFormModal}
            className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center"
          >
            Thêm tin tức
          </button>
        </div>

        {isLoading ? (
          <div className="text-center py-8 dark:text-gray-300">Loading...</div>
        ) : (
          <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
            {/* DARK MODE: Thêm text color */}
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
              {/* DARK MODE: Thêm background và text color */}
              <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
                <tr>
                  <th scope="col" className="px-6 py-3">
                    Tiêu đề
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Người tạo
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Thời gian
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Tùy chọn
                  </th>
                </tr>
              </thead>
              <tbody>
                {newsList.length === 0 ? (
                  <tr>
                    <td colSpan="4" className="px-6 py-4 text-center">
                      Chưa có tin tức nào
                    </td>
                  </tr>
                ) : (
                  newsList.map((newsItem) => (
                    // DARK MODE: Thêm background và border color
                    <tr
                      key={newsItem.id}
                      className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                    >
                      {/* <td className="px-6 py-4">{newsItem.id}</td> */}

                      {/* DARK MODE: Thêm text color */}
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                      >
                        {newsItem.title}
                      </th>

                      <td className="px-6 py-4">{newsItem.created_by}</td>
                      <td className="px-6 py-4">
                        {formatDate(newsItem.created_at)}
                      </td>
                      <td className="px-6 py-4 flex space-x-2">
                        <button
                          onClick={() => handleEdit(newsItem.id)}
                          className="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                        >
                          Sửa
                        </button>
                        <button
                          onClick={() => handleDelete(newsItem.id)}
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

      {/* Các modals không cần thay đổi vì chúng là các component riêng đã hỗ trợ dark mode */}
      {deleteModal && (
        <WarningNotification
          id={newsId}
          func={{
            refresh: refreshData,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"news"}
          action={"news"}
          method={"delete"}
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
      {newsFormModal && (
        <NewsForm
          func={{
            closeModal: closeNewsFormModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh: refreshData,
          }}
          newsId={newsId}
        />
      )}
    </div>
  );
};

export default NewsManagement;
