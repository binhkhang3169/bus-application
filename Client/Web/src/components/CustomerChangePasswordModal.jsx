/** @format */

import axios from "axios";
import React, { useState } from "react";
import { API_URL } from "../configs/env";
import { useNavigate } from "react-router-dom";

const CustomerChangePasswordModal = ({
  closeModal,
  setMessage,
  openFailureModal,
  openSuccessModal,
}) => {
  const navigate = useNavigate();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newPasswordConfirmation, setPasswordConfirmation] = useState("");

  const changePassword = async () => {
    let data = {
      current_password: currentPassword,
      new_password: newPassword,
      new_password_confirmation: newPasswordConfirmation,
    };

    const token = sessionStorage.getItem("token");
    if (token) {
      axios
        .put(API_URL + "customer/change-password", data, {
          headers: { Authorization: `Bearer ${token}` },
        })
        .then((res) => {
          console.log(res.data);
          setMessage(res.data.message);
          openSuccessModal();
          closeModal();

          resetInput();
        })
        .catch((err) => {
          if (err.response.status === 401) {
            navigate("/login");
          }
          setMessage(err.response.data.message);
          openFailureModal();
        });
    } else {
      navigate("/login");
    }
  };

  // Reset input
  const resetInput = () => {
    setCurrentPassword("");
    setNewPassword("");
    setPasswordConfirmation("");
  };
  return (
    <div className="fixed z-50 top-0 left-0 bg-black/20 w-full h-full">
      <div
        id="authentication-modalCustomerChangePasswordModal"
        className=" overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full"
      >
        <div className="relative p-4 w-full max-w-md max-h-full mx-auto">
          {/* <!-- Modal content --> */}
          <div className="relative bg-white rounded-lg shadow ">
            {/* <!-- Modal header --> */}
            <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t ">
              <h3 className="text-xl font-semibold text-gray-900">
                Change password
              </h3>
              <button
                onClick={() => closeModal()}
                type="button"
                className="end-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center "
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
                <span className="sr-only">
                  Close modalCustomerChangePasswordModal
                </span>
              </button>
            </div>
            {/* <!-- Modal body --> */}
            <div className="p-4 md:p-5">
              <div className="space-y-4">
                <div>
                  <label
                    htmlFor="password_current"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Mật khẩu hiện tại
                  </label>
                  <input
                    type="password"
                    name="password_current"
                    id="password_current"
                    className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-red-500 focus:border-red-500 block w-full p-2.5 "
                    placeholder="••••••••"
                    value={currentPassword}
                    onChange={(e) => setCurrentPassword(e.target.value)}
                    required
                  />
                </div>
                <div>
                  <label
                    htmlFor="password"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Mật khẩu mới
                  </label>
                  <input
                    type="password"
                    name="password"
                    id="password"
                    placeholder="••••••••"
                    className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-red-500 focus:border-red-500 block w-full p-2.5 "
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    required
                  />
                </div>
                <div>
                  <label
                    htmlFor="password_confirmation"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Nhập lại mật khẩu
                  </label>
                  <input
                    type="password"
                    name="password_confirmation"
                    id="password_confirmation"
                    placeholder="••••••••"
                    className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-red-500 focus:border-red-500 block w-full p-2.5 "
                    value={newPasswordConfirmation}
                    onChange={(e) => setPasswordConfirmation(e.target.value)}
                    required
                  />
                </div>
                <button
                  type="submit"
                  className="w-full text-white bg-blue-500 hover:bg-blue-600 font-medium rounded-lg text-sm px-5 py-2.5 text-center"
                  onClick={changePassword}
                >
                  Đổi mật khẩu
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CustomerChangePasswordModal;
