/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";

const BusForm = ({ func, busId }) => {
  const navigate = useNavigate();

  // Form Inputs
  const [license, setLicense] = useState("");
  const [seatNumber, setSeatNumber] = useState("");
  const [status, setStatus] = useState(1); // Default to active (1)
  const [typeId, setTypeId] = useState("");
  const [vehicleTypes, setVehicleTypes] = useState([]);
  const [formError, setFormError] = useState(""); // For displaying form-specific errors

  const labelClasses =
    "block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300";
  const inputClasses =
    "shadow-sm bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white";

  useEffect(() => {
    // Fetch vehicle types for the dropdown
    axios
      .get(API_URL + "api/v1/types", {
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      })
      .then((res) => {
        setVehicleTypes(res.data.data || []);
      })
      .catch((err) => {
        console.error("Error fetching vehicle types:", err);
      });

    if (busId && busId !== "") {
      axios
        .get(API_URL + `api/v1/vehicles/${busId}`, {
          headers: {
            Authorization:
              "Bearer " + sessionStorage.getItem("adminAccessToken"),
          },
        })
        .then((res) => {
          const busData = res.data.data;
          setLicense(busData.license);
          setSeatNumber(busData.seat_number.toString());
          setStatus(busData.status);
          setTypeId(busData.type ? busData.type.id.toString() : "");
          setFormError("");
        })
        .catch((err) => {
          console.error("Error fetching bus details:", err);
          if (err.response && err.response.status === 401) {
            navigate("/admin");
          } else {
            setFormError(
              "Failed to load bus details. Please close and try again."
            );
          }
        });
    } else {
      setLicense("");
      setSeatNumber("");
      setStatus(1);
      setTypeId("");
      setFormError("");
    }
  }, [busId, navigate]);

  const handleSubmit = async () => {
    setFormError("");
    if (!license || !seatNumber || !typeId) {
      setFormError(
        "Please fill in all required fields: License, Seat Number, and Type."
      );
      return;
    }

    const data = {
      license,
      seat_number: parseInt(seatNumber, 10),
      status: parseInt(status, 10),
      type: { id: parseInt(typeId, 10) },
    };

    const request = busId
      ? axios.put(API_URL + `api/v1/vehicles/${busId}`, data, {
          headers: {
            Authorization:
              "Bearer " + sessionStorage.getItem("adminAccessToken"),
          },
        })
      : axios.post(API_URL + "api/v1/vehicles", data, {
          headers: {
            Authorization:
              "Bearer " + sessionStorage.getItem("adminAccessToken"),
          },
        });

    request
      .then((res) => {
        func.closeModal();
        func.setMessage(
          res.data.message ||
            (busId ? "Bus updated successfully!" : "Bus created successfully!")
        );
        func.openSuccessModal();
        func.refresh();
      })
      .catch((err) => {
        console.error("Error submitting form:", err);
        if (err.response && err.response.status === 401) {
          navigate("/admin");
        } else {
          func.closeModal();
          func.setMessage(
            err.response?.data?.message ||
              "An error occurred. Please try again."
          );
          func.openFailureModal();
        }
      });
  };

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex justify-center items-center">
      <div className="relative p-4 w-full max-w-md max-h-full">
        <div className="relative bg-white rounded-lg shadow dark:bg-gray-800">
          <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-700">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              {busId ? "Sửa xe" : "Thêm xe mới"}
            </h3>
            <button
              onClick={() => func.closeModal()}
              type="button"
              className="end-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
            >
              <svg
                className="w-3 h-3"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 14 14"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
                />
              </svg>
              <span className="sr-only">Close modal</span>
            </button>
          </div>
          <div className="p-4 md:p-5">
            {formError && (
              <div className="mb-4 text-sm text-red-700 bg-red-100 p-3 rounded-md dark:bg-red-900/30 dark:text-red-400">
                {formError}
              </div>
            )}
            <div className="space-y-4">
              <div>
                <label htmlFor="license" className={labelClasses}>
                  Biển số xe
                </label>
                <input
                  onChange={(e) => setLicense(e.target.value)}
                  value={license}
                  type="text"
                  id="license"
                  className={inputClasses}
                  placeholder="Ví dụ: 51A-12345"
                  required
                />
              </div>
              <div>
                <label htmlFor="seat_number" className={labelClasses}>
                  Số lượng ghế
                </label>
                <input
                  onChange={(e) => setSeatNumber(e.target.value)}
                  value={seatNumber}
                  type="number"
                  id="seat_number"
                  className={inputClasses}
                  placeholder="Ví dụ: 17"
                  required
                  min="1"
                />
              </div>
              <div>
                <label htmlFor="type" className={labelClasses}>
                  Loại xe
                </label>
                <select
                  id="type"
                  value={typeId}
                  onChange={(e) => setTypeId(e.target.value)}
                  className={inputClasses}
                  required
                >
                  <option value="">-- Chọn loại xe --</option>
                  {vehicleTypes.map((type) => (
                    <option key={type.id} value={type.id.toString()}>
                      {type.name}
                    </option>
                  ))}
                </select>
              </div>
              {busId !== "" && (
                <div>
                  <label htmlFor="status" className={labelClasses}>
                    Trạng thái
                  </label>
                  <select
                    id="status"
                    value={status.toString()}
                    onChange={(e) => setStatus(parseInt(e.target.value, 10))}
                    className={inputClasses}
                  >
                    <option value="1">Hoạt động</option>
                    <option value="0">Ngừng hoạt động</option>
                  </select>
                </div>
              )}

              <button
                type="button"
                onClick={handleSubmit}
                className="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
              >
                {busId ? "Xác nhận sửa" : "Thêm xe"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BusForm;
