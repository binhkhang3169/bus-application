/** @format */

import { useEffect, useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import api from "../services/apiService"; // Assuming this path is correct

function Navbar() {
  const [customerInfo, setCustomerInfo] = useState(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  // const navigate = useNavigate(); // Not strictly needed here if not navigating from Navbar directly

  const fetchAndUpdateUserInfo = async () => {
    const accessToken =
      localStorage.getItem("accessToken") ||
      sessionStorage.getItem("accessToken");

    if (!accessToken) {
      clearAuthClientState();
      return;
    }

    try {
      console.log("Navbar: Attempting to fetch user-info");
      const response = await api.get("customer/info"); // Endpoint relative to baseURL in apiService
      if (response.data && response.status === 200) {
        console.log("Navbar: User info received", response.data);
        setCustomerInfo(response.data);
        setIsLoggedIn(true);

        // Optionally update localStorage/sessionStorage if other components read from there directly
        // This helps keep stored info fresh if user info changes during a session
        const storage = localStorage.getItem("accessToken")
          ? localStorage
          : sessionStorage;
        storage.setItem("phoneNumber", response.data.data.phoneNumber);
        storage.setItem("username", response.data.data.username);
        storage.setItem("userFullName", response.data.data.fullName);
        // You might also store other details like role if your app uses them globally
      } else {
        // This case might be an unexpected response structure
        console.warn(
          "Navbar: User info response not as expected or error code returned.",
          response
        );
        clearAuthClientState(); // Clear client state if response is not valid
      }
    } catch (error) {
      console.error(
        "Navbar: Failed to fetch user info:",
        error.response?.data?.message || error.message
      );
      // The apiService interceptor should handle 401s (token refresh or calling clearAuthDataAndRedirect).
      // If it's not a 401 or if refresh failed, the interceptor might have already cleared tokens.
      // We ensure client state is also cleared here as a fallback or if interceptor didn't fully handle UI.
      clearAuthClientState();
    }
  };

  const clearAuthClientState = () => {
    console.log("Navbar: Clearing client authentication state.");
    setCustomerInfo(null);
    setIsLoggedIn(false);
    // No need to remove tokens here as apiService's clearAuthDataAndRedirect or logout function should handle it.
    // This function primarily updates the Navbar's own state.
  };

  useEffect(() => {
    fetchAndUpdateUserInfo(); // Check initially

    const handleAuthChange = (event) => {
      console.log("Navbar: authChange event detected.", event.detail);
      // event.detail can carry information, e.g., { loggedIn: true/false }
      // For simplicity, we'll just re-fetch user info or clear state
      const tokenExists = !!(
        localStorage.getItem("accessToken") ||
        sessionStorage.getItem("accessToken")
      );
      if (tokenExists) {
        fetchAndUpdateUserInfo();
      } else {
        clearAuthClientState();
      }
    };

    window.addEventListener("authChange", handleAuthChange);
    return () => window.removeEventListener("authChange", handleAuthChange);
  }, []);

  const userDisplayFullName =
    customerInfo?.fullName ||
    localStorage.getItem("userFullName") ||
    sessionStorage.getItem("userFullName") ||
    "Account";

  return (
    <nav className="bg-white w-full z-20 top-0 start-0 border-b border-gray-200">
      <div className="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4">
        <NavLink
          to="/"
          className="flex items-center space-x-3 rtl:space-x-reverse"
        >
          <img
            src="https://cdn.iconscout.com/icon/free/png-256/free-bus-1817190-1538058.png"
            className="h-12"
            alt="Anh Phụng BUS LINE Logo"
          />
          <span className="self-center text-2xl font-semibold whitespace-nowrap">
            Anh Phụng BUS LINE
          </span>
        </NavLink>
        <div className="flex lg:order-2 space-x-3 lg:space-x-0 rtl:space-x-reverse">
          {/* Mobile View */}
          <div className="lg:hidden">
            {!isLoggedIn ? (
              <NavLink
                to="/login"
                className="text-white bg-red-500 hover:bg-red-600 font-medium rounded-lg text-sm px-4 py-2 text-center"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                  className="w-6 h-6"
                >
                  <path
                    fillRule="evenodd"
                    d="M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z"
                    clipRule="evenodd"
                  />
                </svg>
              </NavLink>
            ) : (
              <NavLink
                to="/account" // Assuming '/account' is the route for profile management
                className="text-white bg-red-500 hover:bg-red-600 font-medium rounded-lg text-sm px-4 py-2 text-center"
              >
                <div className="flex items-center justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    className="w-6 h-6"
                  >
                    <path
                      fillRule="evenodd"
                      d="M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z"
                      clipRule="evenodd"
                    />
                  </svg>
                  {/* Optionally show name on mobile if space allows, or keep it icon-only */}
                  {/* <span className="ml-2 hidden sm:block">{userDisplayFullName}</span> */}
                </div>
              </NavLink>
            )}
          </div>

          {/* Desktop View */}
          <div className="hidden lg:flex">
            {!isLoggedIn ? (
              <NavLink
                to="/login"
                className="cursor-pointer text-white bg-blue-500 hover:bg-blue-600 font-medium rounded-lg text-sm px-4 py-2 text-center"
              >
                Đăng nhập/Đăng ký
              </NavLink>
            ) : (
              <NavLink
                to="/account"
                className="cursor-pointer flex items-center text-white bg-blue-500 hover:bg-blue-600 font-medium rounded-lg text-sm px-4 py-2 text-center"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                  className="w-6 h-6 mr-2"
                >
                  <path
                    fillRule="evenodd"
                    d="M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z"
                    clipRule="evenodd"
                  />
                </svg>
                <span className="hidden sm:block">{userDisplayFullName}</span>
              </NavLink>
            )}
          </div>

          <button
            data-collapse-toggle="navbar-sticky"
            type="button"
            className="inline-flex items-center p-2 w-10 h-10 justify-center text-sm text-gray-500 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200"
            aria-controls="navbar-sticky"
            aria-expanded="false" // This should be dynamic if you implement the collapse toggle
          >
            <span className="sr-only">Open main menu</span>
            <svg
              className="w-5 h-5"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 17 14"
            >
              <path
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M1 1h15M1 7h15M1 13h15"
              />
            </svg>
          </button>
        </div>
        <div
          className="items-center justify-between hidden w-full lg:flex lg:w-auto lg:order-1"
          id="navbar-sticky"
        >
          <ul className="flex flex-col p-4 lg:p-0 mt-4 font-medium border border-gray-100 rounded-lg bg-gray-50 lg:bg-white lg:space-x-4 rtl:space-x-reverse lg:flex-row lg:mt-0 lg:border-0">
            <li>
              <NavLink
                to="/"
                className={({ isActive }) =>
                  isActive
                    ? "block py-2 px-3 text-white bg-red-500 rounded lg:bg-transparent lg:text-red-700 lg:p-0"
                    : "block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 lg:hover:bg-transparent lg:hover:text-red-700 lg:p-0"
                }
                aria-current="page"
              >
                Trang chủ
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/news" // Assuming a route for news
                className={({ isActive }) =>
                  isActive
                    ? "block py-2 px-3 text-white bg-red-500 rounded lg:bg-transparent lg:text-red-700 lg:p-0"
                    : "block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 lg:hover:bg-transparent lg:hover:text-red-700 lg:p-0"
                }
              >
                Tin tức
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/search-trip"
                className={({ isActive }) =>
                  isActive
                    ? "block py-2 px-3 text-white bg-red-500 rounded lg:bg-transparent lg:text-red-700 lg:p-0"
                    : "block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 lg:hover:bg-transparent lg:hover:text-red-700 lg:p-0"
                }
              >
                Tìm chuyến
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/lookup-ticket"
                className={({ isActive }) =>
                  isActive
                    ? "block py-2 px-3 text-white bg-red-500 rounded lg:bg-transparent lg:text-red-700 lg:p-0"
                    : "block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 lg:hover:bg-transparent lg:hover:text-red-700 lg:p-0"
                }
              >
                Tra cứu vé
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/contact" // Assuming a route for contact
                className={({ isActive }) =>
                  isActive
                    ? "block py-2 px-3 text-white bg-red-500 rounded lg:bg-transparent lg:text-red-700 lg:p-0"
                    : "block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 lg:hover:bg-transparent lg:hover:text-red-700 lg:p-0"
                }
              >
                Liên hệ
              </NavLink>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
