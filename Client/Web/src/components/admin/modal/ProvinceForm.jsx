/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../../configs/env";

const ProvinceForm = ({ func, provinceId }) => {
  const [name, setName] = useState("");
  const [status, setStatus] = useState(1);

  const setFormInputs = (data) => {
    setName(data.name || "");
    setStatus(data.status === 1 || data.status === true ? 1 : 0);
  };

  const resetForm = () => {
    setName("");
    setStatus(1);
  };

  useEffect(() => {
    if (provinceId && provinceId !== "") {
      axios
        .get(API_URL + `api/v1/provinces/${provinceId}`, {
          headers: {
            Authorization:
              "Bearer " + sessionStorage.getItem("adminAccessToken"),
          },
        })
        .then((res) => {
          if (res.data && res.data.data) {
            setFormInputs(res.data.data);
          } else {
            func.setMessage("Could not load province data for editing.");
            func.openFailureModal();
            func.closeModal();
          }
        })
        .catch((err) => {
          func.setMessage("Failed to fetch province details.");
          func.openFailureModal();
          func.closeModal();
        });
    } else {
      resetForm();
    }
  }, [provinceId]);

  const handleSubmit = async (isUpdate) => {
    if (!name.trim()) {
      func.setMessage("Province name is required.");
      func.openFailureModal();
      return;
    }

    let payload = {
      name: name.trim(),
      status: parseInt(status, 10),
    };

    const endpoint = isUpdate
      ? API_URL + `api/v1/provinces/${provinceId}`
      : API_URL + "api/v1/provinces";
    const method = isUpdate ? "put" : "post";

    try {
      const res = await axios({
        method: method,
        url: endpoint,
        data: payload,
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      });
      func.closeModal();
      func.setMessage(
        res.data.message ||
          `Province ${isUpdate ? "updated" : "created"} successfully!`
      );
      func.openSuccessModal();
      func.refresh();
    } catch (err) {
      func.setMessage(
        err.response?.data?.message ||
          `Failed to ${isUpdate ? "update" : "create"} province.`
      );
      func.openFailureModal();
    }
  };

  const handleCreate = () => handleSubmit(false);
  const handleUpdate = () => handleSubmit(true);

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
              {provinceId === "" ? "Thêm tỉnh/thành phố" : "Sửa thông tin"}
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
                  Tên tỉnh/thành phố <span className="text-red-500">*</span>
                </label>
                <input
                  onChange={(e) => setName(e.target.value)}
                  value={name}
                  type="text"
                  id="name"
                  className={inputClasses}
                  placeholder="Ví dụ: TP. Hồ Chí Minh"
                  required
                />
              </div>
              {provinceId !== "" && (
                <div>
                  <label htmlFor="status" className={labelClasses}>
                    Trạng thái <span className="text-red-500">*</span>
                  </label>
                  <select
                    id="status"
                    value={status}
                    onChange={(e) => setStatus(parseInt(e.target.value, 10))}
                    className={inputClasses}
                  >
                    <option value={1}>Hoạt động</option>
                    <option value={0}>Ngừng hoạt động</option>
                  </select>
                </div>
              )}

              <button
                type="button"
                className="w-full text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-50 dark:bg-blue-600 dark:hover:bg-blue-700"
                onClick={provinceId === "" ? handleCreate : handleUpdate}
                disabled={!name.trim()}
              >
                {provinceId === "" ? "Add Province" : "Update Province"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProvinceForm;
