/** @format */

import React, { useEffect, useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import api from "../services/apiService"; // Import the updated API service

function LoginForm() {
  const navigate = useNavigate();

  // Notification
  const [message, setMessage] = useState("");

  // Input
  const [usernameInput, setUsernameInput] = useState(""); // Renamed to avoid conflict
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);

  // Check logged in status
  useEffect(() => {
    const token =
      localStorage.getItem("accessToken") ||
      sessionStorage.getItem("accessToken");
    if (token) {
      // Optional: Add a token validation call here if necessary
      // For now, assume if token exists, user is "logged in" client-side
      console.log("User already has a token, navigating to home.");
      navigate("/");
    }
  }, [navigate]);

  // Login function
  const submitLogin = async () => {
    if (!usernameInput || !password) {
      setMessage("Vui lòng nhập đầy đủ thông tin.");
      return;
    }

    const loginData = {
      username: usernameInput,
      password,
      rememberMe, // Include rememberMe in the request payload
    };

    try {
      // Ensure the endpoint is correct: 'auth/login' relative to baseURL
      const res = await api.post("auth/login", loginData);

      console.log("Login response:", res);

      if (res.status === 200 && res.data.code === 200) {
        const { accessToken, username, role, refreshToken } = res.data.data; // Destructure accessToken

        if (rememberMe) {
          localStorage.setItem("accessToken", accessToken);
          localStorage.setItem("username", username); // Store username
          localStorage.setItem("userRole", role); // Store role
          localStorage.setItem("refreshToken", refreshToken);
        } else {
          sessionStorage.setItem("accessToken", accessToken);
          sessionStorage.setItem("username", username); // Store username
          sessionStorage.setItem("userRole", role); // Store role
          sessionStorage.setItem("refreshToken", refreshToken);
        }

        // Dispatch an event to notify other components (e.g., Navbar) about auth change
        window.dispatchEvent(new Event("authChange"));
        setMessage("Đăng nhập thành công!");
        navigate("/");
      } else {
        // Handle cases where status is 200 but backend indicates an error via 'code'
        setMessage(
          res.data.message || "Đăng nhập không thành công. Vui lòng thử lại."
        );
      }
    } catch (err) {
      console.error("Login error:", err);
      if (err.response) {
        // The response interceptor in apiService should handle 401s for token expiry
        // For login specific errors (e.g. wrong credentials), they will come here.
        setMessage(
          err.response.data.message || "Đã xảy ra lỗi. Vui lòng thử lại."
        );
      } else if (err.request) {
        setMessage(
          "Không thể kết nối đến máy chủ. Vui lòng kiểm tra kết nối mạng."
        );
      } else {
        setMessage("Đã có lỗi xảy ra trong quá trình đăng nhập.");
      }
    }
  };

  return (
    <div className="max-w-screen-lg mx-auto mb-20 px-2">
      <div className="mx-auto w-full border-red-200 border-4 rounded-xl p-8 shadow-xl flex flex-col sm:flex-row justify-center items-center">
        <div className="w-full sm:basis-5/12 mx-4">
          <h1 className="text-4xl font-bold text-center mb-10">Đăng nhập</h1>
          <div>
            <div className="relative mb-8">
              <input
                type="text"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=" " // Important for label animation
                autoComplete="username"
                value={usernameInput}
                onChange={(e) => {
                  setUsernameInput(e.target.value);
                  setMessage("");
                }}
              />
              <label
                htmlFor="username" // Should match input id if you add one
                className="absolute text-sm duration-300 transform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Tên đăng nhập
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
                    d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"
                  />
                </svg>
              </span>
            </div>
            <div className="relative mb-4">
              <input
                type="password"
                className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 focus:border-red-600 peer"
                placeholder=" " // Important for label animation
                autoComplete="current-password"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                  setMessage("");
                }}
                onKeyDown={(event) => {
                  if (event.key === "Enter") {
                    submitLogin();
                  }
                }}
              />
              <label
                htmlFor="password" // Should match input id if you add one
                className="absolute text-sm duration-300 transform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Nhập mật khẩu
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

              {message && (
                <div className="text-red-600 mt-2 text-sm font-medium">
                  {message}
                </div>
              )}
            </div>
            <div className="flex justify-between items-center mb-4">
              <div className="flex gap-2 items-center">
                <input
                  type="checkbox"
                  id="rememberMe"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                  className="h-4 w-4 text-red-600 border-gray-300 rounded focus:ring-red-500"
                />
                <label htmlFor="rememberMe" className="text-sm">
                  Ghi nhớ tài khoản
                </label>
              </div>
              <span className="text-sm text-blue-700 hover:underline">
                <NavLink to="/forgot-password">Quên mật khẩu?</NavLink>
              </span>
            </div>
            <button
              type="button"
              className="transition-colors duration-300 w-full mb-4 text-[18px] mt-6 rounded-full bg-blue-500 text-white hover:bg-blue-600 py-2.5 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
              onClick={submitLogin}
            >
              Đăng nhập
            </button>
            <div>
              <span className="m-4 text-sm">
                Chưa có tài khoản?{" "}
                <NavLink to="/signup" className="text-blue-700 hover:underline">
                  Đăng ký ngay
                </NavLink>
              </span>
            </div>
          </div>
        </div>
        <div className="hidden sm:block flex-grow mx-4">
          <img
            alt="Login illustration"
            className="transition-all duration-200 relative hidden object-contain sm:block h-full w-full"
            src="https://storage.googleapis.com/futa-busline-cms-dev/image_f922bef1bb/image_f922bef1bb.svg"
          />
        </div>
      </div>
      {/* Assuming HomePromotion component exists and is correctly imported */}
      {/* <HomePromotion /> */}
    </div>
  );
}

export default LoginForm;
