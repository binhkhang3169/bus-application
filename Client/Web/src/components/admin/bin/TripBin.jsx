/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";
import WarningNotification from "../../Noti/WarningNotification";
import SuccessNotification from "../../Noti/SuccessNotification";
import FailureNotification from "../../Noti/FailureNotification";

const TripBin = () => {
  const navigate = useNavigate();

  // Modal states
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [message, setMessage] = useState("");
  const [tripId, setTripId] = useState("");

  // Data
  const [tripAll, setTripAll] = useState([]);

  // Get data for input
  useEffect(() => {
    getTripAll();
  }, []);

  // Send GET request to retrieve trips information
  const getTripAll = async () => {
    await axios
      .get(API_URL + "api/v1/trips")
      .then((res) => {
        const trips = res.data.data.filter(
          (trip) =>
            trip.status === 0 &&
            new Date(`${trip.departureDate}T${trip.departureTime}`) > new Date()
        );
        setTripAll(trips);
      })
      .catch((err) => {
        console.error("Error fetching trips:", err);
      });
  };

  // Restore trip function
  const restoreBtn = async (id) => {
    await axios
      .put(
        API_URL + `api/v1/trips/${id}`,
        { status: 1 },
        {
          headers: {
            Authorization: "Bearer " + sessionStorage.getItem("token"),
          },
        }
      )
      .then((res) => {
        setMessage("Trip restored successfully");
        openSuccessModal();
        refresh();
      })
      .catch((err) => {
        setMessage(err.response?.data?.message || "Failed to restore trip");
        openFailureModal();
      });
  };

  // Modal handlers
  const deleteBtn = (id) => {
    setTripId(id);
    setDeleteModal(true);
  };
  const refresh = () => {
    getTripAll();
    setTripId("");
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setTripId("");
  };
  const closeSuccessModal = () => {
    setSuccessModal(false);
  };
  const openSuccessModal = () => {
    setSuccessModal(true);
  };
  const closeFailureModal = () => {
    setFailureModal(false);
  };
  const openFailureModal = () => {
    setFailureModal(true);
  };

  // Format price to display with thousands separator
  const formatPrice = (price) => {
    return new Intl.NumberFormat("vi-VN").format(price);
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white">
            Trip Bin
          </h1>
          <button
            className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center"
            onClick={() => navigate("/admin/trip")}
          >
            Trip Management
          </button>
        </div>
        <div className="mt-4">
          <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
              <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-800 dark:text-gray-300">
                <tr>
                  <th scope="col" className="px-6 py-3">
                    Route
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Seat
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Departure
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Arrival
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Stations
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Price
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Vehicle
                  </th>
                  <th scope="col" className="px-6 py-3">
                    Action
                  </th>
                </tr>
              </thead>
              <tbody>
                {tripAll.length > 0 &&
                  tripAll.map((trip, i) => (
                    <tr
                      key={i}
                      className="odd:bg-white even:bg-gray-50 border-b dark:border-gray-700 dark:odd:bg-gray-900 dark:even:bg-gray-800/50"
                    >
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                      >
                        {`${trip.route.start.name} - ${trip.route.end.name}`}
                      </th>
                      <td className="px-6 py-4">{trip.vehicle.seat_number}</td>
                      <td className="px-6 py-4">{`${trip.departureDate} ${trip.departureTime}`}</td>
                      <td className="px-6 py-4">{`${trip.arrivalDate} ${trip.arrivalTime}`}</td>
                      <td className="px-6 py-4">{`${trip.route.start.name} - ${trip.route.end.name}`}</td>
                      <td className="px-6 py-4">
                        {formatPrice(trip.route.price)} đ
                        {trip.special && (
                          <span className="ml-1 text-green-600 dark:text-green-400">
                            (-{trip.special.percent}%)
                          </span>
                        )}
                      </td>
                      <td className="px-6 py-4">{trip.vehicle.license}</td>
                      <td className="px-6 py-4">
                        <button
                          onClick={() => restoreBtn(trip.id)}
                          className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                        >
                          Restore
                        </button>
                        <button
                          onClick={() => deleteBtn(trip.id)}
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
          id={tripId}
          func={{
            refresh,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"trip"}
          action={"trip"}
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

export default TripBin;
