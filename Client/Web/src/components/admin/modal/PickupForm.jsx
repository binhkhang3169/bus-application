/** @format */

import React, { useState } from "react";
import axios from "axios";
import { API_URL } from "../../../configs/env";

const PickupForm = ({ func, stations, routes, pickups }) => {
  const [routeId, setRouteId] = useState("");
  const [pickupPoints, setPickupPoints] = useState([
    { station: "", time: "" },
    { station: "", time: "" },
  ]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handlePointChange = (index, field, value) => {
    const updatedPoints = [...pickupPoints];
    updatedPoints[index][field] = value;
    setPickupPoints(updatedPoints);
  };

  const addPoint = () => {
    setPickupPoints([...pickupPoints, { station: "", time: "" }]);
  };

  const removePoint = (index) => {
    if (pickupPoints.length > 2) {
      const updatedPoints = pickupPoints.filter((_, i) => i !== index);
      setPickupPoints(updatedPoints);
    }
  };
  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validate route
    if (!routeId) {
      func.setMessage("Vui lòng chọn lộ trình.");
      func.openFailureModal();
      return;
    }

    // Validate từng điểm đón
    for (let i = 0; i < pickupPoints.length; i++) {
      const point = pickupPoints[i];
      if (!point.station || !point.time) {
        func.setMessage(
          `Vui lòng điền đầy đủ trạm và thời gian cho điểm ${
            i === 0
              ? "xuất phát"
              : i === pickupPoints.length - 1
              ? "cuối"
              : `dừng ${i}`
          }.`
        );
        func.openFailureModal();
        return;
      }

      if (parseInt(point.time) <= 0 || isNaN(point.time)) {
        func.setMessage(`Thời gian điểm ${i} không hợp lệ (phải > 0).`);
        func.openFailureModal();
        return;
      }
    }

    // Validate trùng lặp trạm
    const stationIds = pickupPoints.map((p) => p.station);
    const uniqueStations = new Set(stationIds);
    if (uniqueStations.size < stationIds.length) {
      func.setMessage("Các trạm không được trùng nhau.");
      func.openFailureModal();
      return;
    }

    // Nếu qua hết validate
    setIsSubmitting(true);
    const token = sessionStorage.getItem("adminAccessToken");
    let previousPickupId = null;

    try {
      for (let i = 0; i < pickupPoints.length; i++) {
        const point = pickupPoints[i];
        let selfId;

        if (i === 0) {
          selfId = "-1";
        } else if (i === pickupPoints.length - 1) {
          selfId = "-2";
        } else {
          selfId = previousPickupId;
        }

        const payload = {
          station: { id: parseInt(point.station, 10) },
          route: { id: parseInt(routeId, 10) },
          time: point.time,
          selfId: selfId,
          status: 1,
          id: (pickups.length + i + 1).toString(),
        };

        const response = await axios.post(`${API_URL}api/v1/pickups`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (i < pickupPoints.length - 1) {
          previousPickupId = response.data.data.id;
        }
      }

      func.setMessage("Tạo điểm đón/trả thành công!");
      func.openSuccessModal();
      func.refresh();
      func.closeModal();
    } catch (error) {
      console.error("Failed to create pickup path:", error);
      func.setMessage(error.response?.data?.message || "Đã xảy ra lỗi.");
      func.openFailureModal();
    } finally {
      setIsSubmitting(false);
    }
  };

  const labelClasses =
    "block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300";
  const inputClasses =
    "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white";

  return (
    <div className="fixed top-0 right-0 left-0 z-50 flex justify-center items-center w-full h-full bg-black bg-opacity-50">
      <div className="relative p-4 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="relative bg-white rounded-lg shadow dark:bg-gray-800">
          <div className="flex items-center justify-between p-4 border-b rounded-t dark:border-gray-700">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              Thêm điểm mới
            </h3>
            <button
              onClick={func.closeModal}
              className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
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
          <form onSubmit={handleSubmit} className="p-4 md:p-5">
            <div className="grid gap-4 mb-4 grid-cols-1">
              <div>
                <label htmlFor="route" className={labelClasses}>
                  Lộ trình
                </label>
                <select
                  id="route"
                  value={routeId}
                  onChange={(e) => setRouteId(e.target.value)}
                  className={inputClasses}
                >
                  <option value="">Chọn lộ trình</option>
                  {routes.map((r) => (
                    <option
                      key={r.id}
                      value={r.id}
                    >{`${r.start.name} - ${r.end.name}`}</option>
                  ))}
                </select>
              </div>

              {pickupPoints.map((point, index) => (
                <div
                  key={index}
                  className="p-3 border rounded-lg grid grid-cols-1 md:grid-cols-3 gap-4 items-center dark:border-gray-700"
                >
                  <div className="md:col-span-2 grid grid-cols-2 gap-4">
                    <div>
                      <label className={labelClasses}>
                        {index === 0
                          ? "Điểm xuất phát"
                          : index === 1
                          ? "Điểm dừng cuối"
                          : `Stop ${index - 1}`}
                      </label>
                      <select
                        value={point.station}
                        onChange={(e) =>
                          handlePointChange(index, "station", e.target.value)
                        }
                        className={inputClasses}
                      >
                        <option value="">Chọn trạm</option>
                        {stations.map((s) => (
                          <option key={s.id} value={s.id}>
                            {s.name}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div>
                      <label className={labelClasses}>
                        Thời gian có mặt (phút)
                      </label>
                      <input
                        type="number"
                        value={point.time}
                        onChange={(e) =>
                          handlePointChange(index, "time", e.target.value)
                        }
                        className={inputClasses}
                        placeholder="Ví dụ: 20, 30, ..."
                      />
                    </div>
                  </div>
                  <div className="flex items-end h-full">
                    {index > 1 && (
                      <button
                        type="button"
                        onClick={() => removePoint(index)}
                        className="text-red-600 dark:text-red-500 hover:text-red-700 dark:hover:text-red-400 font-medium"
                      >
                        Xóa
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
            <div className="flex items-center gap-4">
              <button
                type="button"
                onClick={addPoint}
                className="text-white bg-green-700 hover:bg-green-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-green-600 dark:hover:bg-green-700"
              >
                Thêm trạm dừng giữa đường
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 disabled:opacity-50"
              >
                {isSubmitting ? "Đang tạo..." : "Hoàn tất"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default PickupForm;
