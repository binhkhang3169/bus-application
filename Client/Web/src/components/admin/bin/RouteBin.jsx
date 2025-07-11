/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../../configs/env";
import WarningNotification from "../../Noti/WarningNotification";
import SuccessNotification from "../../Noti/SuccessNotification";
import FailureNotification from "../../Noti/FailureNotification";
import { useNavigate } from "react-router-dom";

const RouteBin = () => {
  const navigate = useNavigate();

  // Data
  const [routeAll, setRouteAll] = useState([]);

  // Modal
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [message, setMessage] = useState("");
  const [routeId, setRouteId] = useState("");

  useEffect(() => {
    getRouteAll();
  }, []);

  // Send GET request to retrieve trip
  const getRouteAll = async () => {
    await axios
      .get(API_URL + "api/v1/routes")
      .then((res) => {
        setRouteAll(res.data.route.filter((v) => v.status === 0));
      })
      .catch((err) => {});
  };

  // Open route edit modal
  const restoreBtn = async (id) => {
    await axios
      .put(
        API_URL + `v1/routes/${id}`,
        { status: 1 },
        {
          headers: {
            Authorization: "Bearer" + sessionStorage.getItem("token"),
          },
        }
      )
      .then((res) => {
        setMessage(res.data.message);
        openSuccessModal();
        refresh();
      })
      .catch((err) => {
        setMessage(err.response.data.message);
        openFailureModal();
      });
  };

  // Open delete confirm modal
  const deleteBtn = async (id) => {
    setRouteId(id);
    setDeleteModal(true);
  };

  // Refresh page
  const refresh = () => {
    setRouteId("");
    getRouteAll();
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setRouteId("");
  };
  const closeSuccessModal = () => {
    setSuccessModal(false);
  };
  const closeFailureModal = () => {
    setFailureModal(false);
  };
  const openSuccessModal = () => {
    setSuccessModal(true);
  };
  const openFailureModal = () => {
    setFailureModal(true);
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white">
            Route Bin
          </h1>
          <button
            className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center"
            onClick={() => navigate("/admin/route")}
          >
            Route Management
          </button>
        </div>
        <div className="mt-4">
          <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
              <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-800 dark:text-gray-300">
                <tr>
                  <th scope="col" className="px-6 py-3">
                    Name
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Start Address
                  </th>
                  <th scope="col" className="px-6 py-3">
                    End Address
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Price
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Time
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Action
                  </th>
                </tr>
              </thead>
              <tbody>
                {routeAll &&
                  routeAll.map((v, i) => (
                    <tr
                      key={i}
                      className="odd:bg-white even:bg-gray-50 border-b dark:border-gray-700 dark:odd:bg-gray-900 dark:even:bg-gray-800/50"
                    >
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                      >
                        {v.name}
                      </th>
                      <td className="px-6 py-4">{v.start_address.name}</td>
                      <td className="px-6 py-4">{v.end_address.name}</td>
                      <td className="px-6 py-4">{v.price}</td>
                      <td className="px-6 py-4">{v.time}</td>
                      <td className="px-6 py-4">
                        <button
                          onClick={() => restoreBtn(v.id)}
                          className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                        >
                          Restore
                        </button>
                        <button
                          onClick={() => deleteBtn(v.id)}
                          className="font-medium text-red-600 dark:text-red-500 hover:underline"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
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

export default RouteBin;
