/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";
import { Arrow } from "../../../svg/svg";

const BusStationForm = ({ func, busStationId, provinces }) => {
  const navigate = useNavigate();
  const [citySearchModal, setCitySearchModal] = useState(false);
  const [citySearch, setCitySearch] = useState("");
  const [name, setName] = useState("");
  const [city, setCity] = useState("");
  const [selectedProvinceId, setSelectedProvinceId] = useState(null);
  const [address, setAddress] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");

  const setFormInputs = (data) => {
    setName(data.name || "");
    setAddress(data.address || "");
    if (data.province && typeof data.province === "object") {
      setCity(data.province.name || "");
      setSelectedProvinceId(data.province.id || null);
    } else {
      setCity("");
      setSelectedProvinceId(null);
      if (data.province?.id) {
        setSelectedProvinceId(data.province.id);
        const foundProvince = provinces.find((p) => p.id === data.province.id);
        setCity(
          foundProvince ? foundProvince.name : "Error: Province not found"
        );
      }
    }
    setPhoneNumber(data.phone_number || "");
  };

  const resetForm = () => {
    setName("");
    setAddress("");
    setCity("");
    setSelectedProvinceId(null);
    setPhoneNumber("");
    setCitySearch("");
  };

  useEffect(() => {
    if (busStationId && busStationId !== "") {
      axios
        .get(API_URL + `api/v1/stations/${busStationId}`)
        .then((res) => {
          if (res.data && res.data.data) {
            setFormInputs(res.data.data);
          }
        })
        .catch((err) => {
          if (err.response?.status === 401) navigate("/admin");
        });
    } else {
      resetForm();
    }
  }, [busStationId, provinces, navigate]);

  const handleSubmit = async (isUpdate) => {
    if (!selectedProvinceId || !name.trim() || !address.trim()) {
      func.setMessage("Name, Address, and City/Province are required.");
      func.openFailureModal();
      return;
    }

    let payload = {
      name: name.trim(),
      address: address.trim(),
      province: { id: selectedProvinceId },
    };
    if (phoneNumber.trim()) payload.phone_number = phoneNumber.trim();

    const endpoint = isUpdate
      ? API_URL + `api/v1/stations/${busStationId}`
      : API_URL + "api/v1/stations";
    const method = isUpdate ? "put" : "post";

    try {
      const res = await axios({
        method,
        url: endpoint,
        data: payload,
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      });
      func.closeModal();
      func.setMessage(
        res.data.message || `${isUpdate ? "Sửa" : "Thêm"} thành công!`
      );
      func.openSuccessModal();
      func.refresh();
    } catch (err) {
      if (err.response?.status === 401) navigate("/admin");
      func.setMessage(
        err.response?.data?.message ||
          `Failed to ${isUpdate ? "update" : "create"} station.`
      );
      func.openFailureModal();
    }
  };

  const sendRequestCreateBusStation = () => handleSubmit(false);
  const sendRequestUpdateBusStation = () => handleSubmit(true);

  const labelClasses =
    "block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300";
  const inputClasses =
    "shadow-sm bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white";

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex justify-center items-center">
      <div className="relative p-4 w-full max-w-md max-h-full">
        <div className="relative bg-white rounded-lg shadow dark:bg-gray-800">
          <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-700">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
              {busStationId === "" ? "Thêm trạm xe mới" : "Sửa trạm xe"}
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
                <label htmlFor="name" className={labelClasses}>
                  Tên trạm <span className="text-red-500">*</span>
                </label>
                <input
                  onChange={(e) => setName(e.target.value)}
                  value={name}
                  type="text"
                  id="name"
                  className={inputClasses}
                  placeholder="Ví dụ: Bến xe Miền Đông, ..."
                  required
                />
              </div>
              <div>
                <label htmlFor="address" className={labelClasses}>
                  Địa chỉ <span className="text-red-500">*</span>
                </label>
                <input
                  onChange={(e) => setAddress(e.target.value)}
                  value={address}
                  type="text"
                  id="address"
                  className={inputClasses}
                  placeholder="Ví dụ: 292 Đinh Bộ Lĩnh, P.26, Bình Thạnh, ..."
                  required
                />
              </div>
              <div>
                <label htmlFor="city" className={labelClasses}>
                  Tỉnh/thành phố <span className="text-red-500">*</span>
                </label>
                <div className="text-sm rounded-lg relative h-10 border border-gray-300 bg-gray-50 text-gray-900 cursor-pointer dark:bg-gray-700 dark:border-gray-600 dark:text-white">
                  <div
                    className="w-full h-full flex items-center justify-between px-2.5"
                    onClick={() => {
                      setCitySearchModal(!citySearchModal);
                      if (citySearchModal) setCitySearch("");
                    }}
                  >
                    <span
                      className={city ? "dark:text-white" : "text-gray-400"}
                    >
                      {city || "Chọn tỉnh/thành phố"}
                    </span>
                    <Arrow
                      className={`h-4 w-4 fill-gray-500 transition-transform duration-200 dark:fill-gray-300 ${
                        citySearchModal ? "rotate-0" : "rotate-180"
                      }`}
                    />
                  </div>
                  {citySearchModal && (
                    <div className="absolute z-40 top-full mt-1 left-0 w-full bg-white rounded-md shadow-lg border border-gray-200 dark:bg-gray-900 dark:border-gray-700">
                      <input
                        value={citySearch}
                        className="h-10 w-full px-3 py-2 border-b border-gray-200 focus:outline-none text-sm dark:bg-gray-900 dark:text-white dark:border-gray-700"
                        type="text"
                        onChange={(e) => setCitySearch(e.target.value)}
                        placeholder="Search province..."
                        autoFocus
                      />
                      <ul className="w-full overflow-y-auto max-h-40 text-sm">
                        {provinces && provinces.length > 0 ? (
                          provinces
                            .filter((p) =>
                              p.name
                                .toLowerCase()
                                .includes(citySearch.toLowerCase())
                            )
                            .map((province, i) => (
                              <li
                                key={province.id || i}
                                className="px-3 py-2 text-gray-800 hover:bg-gray-100 cursor-pointer dark:text-gray-300 dark:hover:bg-gray-700"
                                onClick={() => {
                                  setCity(province.name);
                                  setSelectedProvinceId(province.id);
                                  setCitySearchModal(false);
                                  setCitySearch("");
                                }}
                              >
                                {province.name}
                              </li>
                            ))
                        ) : (
                          <li className="px-3 py-2 text-gray-500 dark:text-gray-400">
                            Không có tỉnh nào
                          </li>
                        )}
                        {provinces &&
                          provinces.filter((p) =>
                            p.name
                              .toLowerCase()
                              .includes(citySearch.toLowerCase())
                          ).length === 0 &&
                          citySearch !== "" && (
                            <li className="px-3 py-2 text-gray-500 dark:text-gray-400">
                              No matching provinces
                            </li>
                          )}
                      </ul>
                    </div>
                  )}
                </div>
              </div>
              {/* <div>
                <label htmlFor="phoneNumber" className={labelClasses}>
                  Phone
                </label>
                <input
                  onChange={(e) => setPhoneNumber(e.target.value)}
                  value={phoneNumber}
                  type="tel"
                  id="phoneNumber"
                  className={inputClasses}
                  placeholder="Ví dụ: 02838994056"
                />
              </div> */}
              <button
                type="button"
                className="w-full text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-50 dark:bg-blue-600 dark:hover:bg-blue-700"
                onClick={
                  busStationId === ""
                    ? sendRequestCreateBusStation
                    : sendRequestUpdateBusStation
                }
                disabled={
                  !name.trim() || !address.trim() || !selectedProvinceId
                }
              >
                {busStationId === "" ? "Thêm trạm" : "Xác nhận sửa"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
export default BusStationForm;
