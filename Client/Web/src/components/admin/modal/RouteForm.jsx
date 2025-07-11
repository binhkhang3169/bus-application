/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";

const RouteForm = ({ func, routeId, busStationAll }) => {
  const navigate = useNavigate();

  // Input
  const [startAddress, setStartAddress] = useState("");
  const [endAddress, setEndAddress] = useState("");
  const [price, setPrice] = useState("");
  const [estimatedTime, setEstimatedTime] = useState("");
  const [distance, setDistance] = useState("");

  useEffect(() => {
    if (routeId !== "") {
      axios
        .get(API_URL + `api/v1/routes/${routeId}`)
        .then((res) => {
          setInput(res.data.data);
        })
        .catch((err) => {
          if (err.response && err.response.status === 401) {
            navigate("/admin");
          }
        });
    } else {
      setStartAddress("");
      setEndAddress("");
      setPrice("");
      setEstimatedTime("");
      setDistance("");
    }
  }, [routeId, navigate]);

  const sendRequestCreateRoute = async () => {
    let data = {
      start: { id: startAddress },
      end: { id: endAddress },
      price,
      estimatedTime,
      distance,
    };

    await axios
      .post(API_URL + "api/v1/routes", data, {
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      })
      .then((res) => {
        func.closeModal();
        func.setMessage(res.data.message);
        func.openSuccessModal();
        func.refresh();
      })
      .catch((err) => {
        func.closeModal();
        func.setMessage(
          err.response?.data?.message || "Failed to create route."
        );
        func.openFailureModal();
      });
  };

  const sendRequestUpdateRoute = async () => {
    let data = {
      start: { id: startAddress },
      end: { id: endAddress },
      price,
      estimatedTime,
      distance,
    };
    await axios
      .put(API_URL + `api/v1/routes/${routeId}`, data, {
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      })
      .then((res) => {
        func.closeModal();
        func.setMessage(res.data.message);
        func.openSuccessModal();
        func.refresh();
      })
      .catch((err) => {
        func.closeModal();
        func.setMessage(
          err.response?.data?.message || "Failed to update route."
        );
        func.openFailureModal();
      });
  };

  const validateInput = () => {
    if (!startAddress) {
      alert("Vui lòng chọn điểm đi.");
      return false;
    }
    if (!endAddress) {
      alert("Vui lòng chọn điểm đến.");
      return false;
    }
    if (startAddress === endAddress) {
      alert("Điểm đi và điểm đến không được trùng nhau.");
      return false;
    }
    if (!price || parseFloat(price) <= 0) {
      alert("Vui lòng nhập giá hợp lệ.");
      return false;
    }
    if (!estimatedTime || estimatedTime.trim() === "") {
      alert("Vui lòng nhập thời gian ước lượng.");
      return false;
    }
    if (!distance || parseFloat(distance) <= 0) {
      alert("Vui lòng nhập khoảng cách hợp lệ.");
      return false;
    }
    return true;
  };

  const setInput = (data) => {
    setStartAddress(data.start_address?.id?.toString() || "");
    setEndAddress(data.end_address?.id?.toString() || "");
    setPrice(data.price || "");
    setEstimatedTime(data.estimated_time || "");
    setDistance(data.distance || "");
  };

  const formatVND = (value) => {
    if (!value) return "";
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
      minimumFractionDigits: 0,
    }).format(value);
  };

  const labelClasses =
    "block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300";
  const inputClasses =
    "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white";

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex justify-center items-center">
      <div className="relative p-4 w-full max-w-md max-h-full">
        <div className="relative bg-white rounded-lg shadow dark:bg-gray-800">
          <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-700">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              {routeId === "" ? "Tạo lộ trình mới" : "Cập nhật lộ trình"}
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
            </button>
          </div>
          <div className="p-4 md:p-5">
            <div className="space-y-4">
              <div>
                <label htmlFor="start-address" className={labelClasses}>
                  Điểm đi
                </label>
                <select
                  onChange={(e) => setStartAddress(e.target.value)}
                  value={startAddress}
                  id="start-address"
                  className={inputClasses}
                >
                  <option value="">Chọn điểm đi</option>
                  {busStationAll &&
                    busStationAll.map((v) => (
                      <option key={v.id} value={v.id.toString()}>
                        {v.name}
                      </option>
                    ))}
                </select>
              </div>
              <div>
                <label htmlFor="end-address" className={labelClasses}>
                  Điểm đến
                </label>
                <select
                  onChange={(e) => setEndAddress(e.target.value)}
                  value={endAddress}
                  id="end-address"
                  className={inputClasses}
                >
                  <option value="">Chọn điểm đến</option>
                  {busStationAll &&
                    busStationAll.map((v) => (
                      <option key={v.id} value={v.id.toString()}>
                        {v.name}
                      </option>
                    ))}
                </select>
              </div>
              <div>
                <label htmlFor="price" className={labelClasses}>
                  Giá
                </label>
                <input
                  onChange={(e) => setPrice(e.target.value)}
                  value={price}
                  type="number"
                  id="price"
                  className={inputClasses}
                  placeholder="Giá vé"
                  required
                />
                {price && (
                  <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
                    {formatVND(price)}
                  </p>
                )}
              </div>

              <div>
                <label htmlFor="estimated-time" className={labelClasses}>
                  Thời gian ước lượng
                </label>
                <input
                  onChange={(e) => setEstimatedTime(e.target.value)}
                  value={estimatedTime}
                  type="text"
                  id="estimated-time"
                  className={inputClasses}
                  placeholder="Ví dụ: 2h 30m or 150m"
                />
              </div>
              <div>
                <label htmlFor="distance" className={labelClasses}>
                  Khoảng cách (km)
                </label>
                <input
                  onChange={(e) => setDistance(e.target.value)}
                  value={distance}
                  type="number"
                  id="distance"
                  className={inputClasses}
                  placeholder="Khoảng cách lộ trình"
                />
              </div>
              <button
                className="w-full text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
                onClick={() => {
                  if (!validateInput()) return;
                  if (routeId === "") {
                    sendRequestCreateRoute();
                  } else {
                    sendRequestUpdateRoute();
                  }
                }}
              >
                {routeId === "" ? "Thêm" : "Cập nhật"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RouteForm;
