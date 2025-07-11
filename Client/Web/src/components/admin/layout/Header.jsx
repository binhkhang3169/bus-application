import React, { useState } from "react";
import { NavLink } from "react-router-dom";

const SunIcon = () => (
  <svg className="fill-current" width="16" height="16" viewBox="0 0 16 16">
    <path
      d="M8 11.5C9.933 11.5 11.5 9.933 11.5 8C11.5 6.067 9.933 4.5 8 4.5C6.067 4.5 4.5 6.067 4.5 8C4.5 9.933 6.067 11.5 8 11.5ZM8 12.5C5.51472 12.5 3.5 10.4853 3.5 8C3.5 5.51472 5.51472 3.5 8 3.5C10.4853 3.5 12.5 5.51472 12.5 8C12.5 10.4853 10.4853 12.5 8 12.5Z"
      fill="currentColor"
    />
    <path
      d="M8 1.5C8.41421 1.5 8.75 1.83579 8.75 2.25V3.75C8.75 4.16421 8.41421 4.5 8 4.5C7.58579 4.5 7.25 4.16421 7.25 3.75V2.25C7.25 1.83579 7.58579 1.5 8 1.5Z"
      fill="currentColor"
    />
  </svg>
);

const MoonIcon = () => (
  <svg className="fill-current" width="16" height="16" viewBox="0 0 16 16">
    <path
      d="M10.9745 12.213C10.8404 12.3551 10.6385 12.5 10.4219 12.5C9.28943 12.5 8.24393 12.067 7.44688 11.27C6.64983 10.473 6.2168 9.42749 6.2168 8.29499C6.2168 7.07099 6.68531 6.01074 7.51443 5.18162C7.58438 5.10662 7.68063 5.05662 7.78188 5.04912C7.88313 5.04162 7.98563 5.07662 8.06063 5.14662L8.79813 5.81912C8.95438 5.96037 8.96563 6.19787 8.82438 6.35412C8.20313 7.04287 7.88188 7.97037 7.88188 8.94162C7.88188 10.6979 9.30254 12.1185 11.0588 12.1185C11.5338 12.1185 11.9903 12.0223 12.4053 11.8473C12.5616 11.7811 12.7328 11.8123 12.8578 11.9373L13.5678 12.6473C13.6703 12.7448 13.7145 12.8854 13.6845 13.02C13.6545 13.1545 13.5538 13.2659 13.4193 13.3084C12.6478 13.5679 11.8195 13.7024 10.9745 13.7024V12.213Z"
      fill="currentColor"
    />
  </svg>
);

const Header = ({
  sidebarToggle,
  setSidebarToggle,
  darkMode,
  setDarkMode,
  role,
  onLogout,
}) => {
  const [notificationOpen, setNotificationOpen] = useState(false);
  const [profileOpen, setProfileOpen] = useState(false);
  const [notifying, setNotifying] = useState(true);

  return (
    <header className="sticky top-0 z-50 flex w-full bg-white drop-shadow-sm dark:bg-gray-800 dark:drop-shadow-none">
      <div className="flex flex-grow items-center justify-between py-4 px-4 shadow-sm md:px-6 2xl:px-11">
        <div className="flex items-center gap-2 sm:gap-4">
          {/* --- Nút Hamburger cho Mobile --- */}
          <button
            aria-controls="sidebar"
            onClick={(e) => {
              e.stopPropagation();
              setSidebarToggle(!sidebarToggle);
            }}
            className="z-99999 block rounded-sm border border-gray-200 bg-white p-1.5 shadow-sm dark:border-gray-700 dark:bg-gray-900 lg:hidden"
          >
            <span className="relative block h-5.5 w-5.5 cursor-pointer">
              {/* Hamburger Lines */}
            </span>
          </button>

          {/* --- Nút Toggle Sidebar cho Desktop --- */}
          <button
            onClick={() => setSidebarToggle(!sidebarToggle)}
            className="z-99999 hidden h-9 w-9 items-center justify-center rounded-lg border border-gray-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-900 lg:flex"
          >
            <svg
              className="w-6 h-6 dark:text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 6h16M4 12h16M4 18h16"
              ></path>
            </svg>
          </button>
        </div>

        <div className="flex items-center gap-3 2xsm:gap-7">
          <ul className="flex items-center gap-2 2xsm:gap-4">
            {/* Dark Mode Toggler */}
            <li>
              <label
                className={`relative m-0 block h-8 w-14 rounded-full ${
                  darkMode ? "bg-blue-600" : "bg-gray-200 dark:bg-gray-700"
                }`}
              >
                <input
                  type="checkbox"
                  onChange={() => setDarkMode(!darkMode)}
                  checked={darkMode}
                  className="absolute top-0 z-50 m-0 h-full w-full cursor-pointer opacity-0"
                />
                <span
                  className={`absolute top-1/2 left-1 flex h-6 w-6 -translate-y-1/2 translate-x-0 items-center justify-center rounded-full bg-white shadow-switcher duration-200 ease-linear ${
                    darkMode && "!right-1 !translate-x-full"
                  }`}
                >
                  <span className="dark:hidden">
                    <SunIcon />
                  </span>
                  <span className="hidden dark:inline">
                    <MoonIcon />
                  </span>
                </span>
              </label>
            </li>
          </ul>

          {/* User Profile Dropdown */}
          <div className="relative">
            <button
              onClick={() => setProfileOpen(!profileOpen)}
              className="flex items-center gap-4"
            >
              <span className="hidden text-right lg:block">
                <span className="block text-sm font-medium text-black dark:text-white">
                  {sessionStorage.getItem("adminUsername") || "Admin"}
                </span>
                {/* <span className="block text-xs dark:text-gray-400">{role}</span> */}
              </span>
            </button>
            {profileOpen && (
              <div className="absolute right-0 mt-4 flex w-62.5 flex-col rounded-lg border border-gray-200 bg-white shadow-default dark:border-gray-700 dark:bg-gray-800">
                <ul className="flex flex-col gap-5 border-b border-gray-200 px-6 py-7.5 dark:border-gray-700">
                  <li>
                    <NavLink
                      to="/admin/account"
                      onClick={() => setProfileOpen(false)}
                      className="flex items-center gap-3.5 text-sm font-medium duration-300 ease-in-out hover:text-blue-600 dark:text-gray-300 dark:hover:text-blue-500 lg:text-base"
                    >
                      Thông tin
                    </NavLink>
                  </li>
                </ul>
                <button
                  onClick={onLogout}
                  className="flex items-center gap-3.5 py-4 px-6 text-sm font-medium duration-300 ease-in-out hover:text-blue-600 dark:text-gray-300 dark:hover:text-blue-500 lg:text-base"
                >
                  Đăng xuất
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
