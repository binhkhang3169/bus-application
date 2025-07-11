/** @format */

import React, { useState } from "react";
import { API_URL, REACT_URL } from "../configs/env";
import axios from "axios";
import HomePromotion from "./HomePromotion";
import SuccessNotification from "./Noti/SuccessNotification";
import FailureNotification from "./Noti/FailureNotification";

function SignupForm() {
  // Notification
  const [message, setMessage] = useState("");
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);

  // Input
  const [username, setUsername] = useState("");
  const [fullName, setFullName] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [address, setAddress] = useState("");
  const [gender, setGender] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  // Close Success Modal
  const closeSuccessModal = () => {
    setSuccessModal(false);
  };

  // Open Success Modal
  const openSuccessModal = () => {
    setSuccessModal(true);
  };

  // Close Failure Modal
  const closeFailureModal = () => {
    setFailureModal(false);
  };

  // Open Failure Modal
  const openFailureModal = () => {
    setFailureModal(true);
  };

  // Send POST request for initial signup
  const signupBtn = async () => {
    // Validate inputs
    if (
      !username ||
      !address ||
      !gender ||
      !password ||
      !fullName ||
      !phoneNumber
    ) {
      setMessage("Vui lòng điền đầy đủ thông tin bắt buộc");
      openFailureModal();
      return;
    }

    const data = {
      username,
      fullName,
      phoneNumber,
      password,
      address,
      gender,
    };
    setIsLoading(true);

    await axios
      .post(API_URL + "api/v1/signup", data)
      .then((res) => {
        if (res.status === 200) {
          setMessage(res.data.message);
          openSuccessModal();

          // Save user data to localStorage
          localStorage.setItem("signupData", JSON.stringify(data));

          // Redirect to verify page
          window.location.href = REACT_URL + "email-verify?email=" + username;
        }
      })
      .catch((err) => {
        setMessage(err.response?.data?.message || "Có lỗi xảy ra khi đăng ký");
        openFailureModal();
      });
  };

  return (
    <div className="max-w-screen-lg mx-auto mb-20 px-2">
      <div className="w-full -mx-4 border-red-200 border-4 rounded-xl p-8 shadow-xl flex justify-center items-center">
        <div className="w-full sm:basis-5/12 mx-4">
          <h1 className="text-4xl font-bold text-center mb-10">
            Đăng ký tài khoản
          </h1>
          <div>
            <div className="relative mb-5">
              <input
                type="text"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=""
                onChange={(e) => setUsername(e.target.value)}
                value={username}
              />
              <label
                htmlFor=""
                className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Email
              </label>
            </div>
            <div className="relative mb-5">
              <input
                type="text"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=""
                onChange={(e) => setFullName(e.target.value)}
                value={fullName}
              />
              <label
                htmlFor=""
                className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Họ và tên
              </label>
            </div>
            <div className="relative mb-5">
              <input
                type="text"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=""
                onChange={(e) => setPhoneNumber(e.target.value)}
                value={phoneNumber}
              />
              <label
                htmlFor=""
                className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Số điện thoại
              </label>
              <span className="absolute top-1 right-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M10.5 1.5H8.25A2.25 2.25 0 0 0 6 3.75v16.5a2.25 2.25 0 0 0 2.25 2.25h7.5A2.25 2.25 0 0 0 18 20.25V3.75a2.25 2.25 0 0 0-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3"
                  />
                </svg>
              </span>
            </div>

            <div className="relative mb-5">
              <input
                type="text"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=""
                onChange={(e) => setAddress(e.target.value)}
                value={address}
              />
              <label
                htmlFor=""
                className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Địa chỉ
              </label>
            </div>
            <div className="relative mb-5">
              <select
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                onChange={(e) => setGender(e.target.value)}
                value={gender}
              >
                <option value="">Chọn giới tính</option>
                <option value="Nam">Nam</option>
                <option value="Nữ">Nữ</option>
                <option value="Khác">Khác</option>
              </select>
              <label
                htmlFor=""
                className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Giới tính
              </label>
            </div>
            <div className="relative mb-5">
              <input
                type="password"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=""
                onChange={(e) => setPassword(e.target.value)}
                value={password}
                onKeyDown={(event) => {
                  if (event.key === "Enter") {
                    signupBtn();
                  }
                }}
              />
              <label
                htmlFor=""
                className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Mật khẩu
              </label>
              <span className="absolute top-1 right-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
                  />
                </svg>
              </span>
            </div>

            {message && (
              <div className="text-red-600 text-sm font-medium">{message}</div>
            )}

            <button
              type="submit"
              className="transition-colors duration-300 w-full mb-4 text-[18px] mt-6 rounded-full bg-blue-500 text-white hover:bg-blue-600 py-2 flex justify-center items-center disabled:bg-gray-400"
              onClick={signupBtn}
              disabled={isLoading}
            >
              {isLoading ? <div className="loader"></div> : "Đăng ký"}
            </button>

            <div>
              <span className="m-4">
                Đã có tài khoản?{" "}
                <a href="/login" className="text-blue-500">
                  Đăng nhập
                </a>
              </span>
            </div>
          </div>
        </div>
        <div className="hidden sm:block flex-grow mx-4">
          <img
            alt=""
            loading="lazy"
            decoding="async"
            data-nimg="fill"
            className="transition-all duration-200 relative hidden object-contain sm:block h-full w-full transparent"
            src="https://storage.googleapis.com/futa-busline-cms-dev/image_f922bef1bb/image_f922bef1bb.svg"
          ></img>
        </div>
      </div>
      <HomePromotion />

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
}

export default SignupForm;
