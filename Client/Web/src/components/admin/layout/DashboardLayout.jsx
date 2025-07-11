import React, { useState, useEffect, useCallback } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";
import Sidebar from "./Sidebar";
import Header from "./Header";
import Preloader from "./Preloader";

const ALLOWED_ADMIN_ROLES = [
  "ROLE_ADMIN",
  "ROLE_RECEPTION",
  "ROLE_DRIVER",
  "ROLE_OPERATOR",
];

const DashboardLayout = () => {
  const navigate = useNavigate();
  const [role, setRole] = useState(null);

  // Khởi tạo state dựa trên kích thước màn hình
  const [sidebarToggle, setSidebarToggle] = useState(window.innerWidth >= 1024);

  const [loading, setLoading] = useState(true);
  const [darkMode, setDarkMode] = useState(false);

  const clearAdminSessionAndRedirect = useCallback(() => {
    sessionStorage.removeItem("adminAccessToken");
    sessionStorage.removeItem("adminUsername");
    sessionStorage.removeItem("adminRole");
    setRole(null);
    window.dispatchEvent(
      new CustomEvent("authChangeAdmin", { detail: { loggedIn: false } })
    );
    navigate("/admin");
  }, [navigate]);

  const checkAdminAuthAndSetRole = useCallback(() => {
    const token = sessionStorage.getItem("adminAccessToken");
    const storedRole = sessionStorage.getItem("adminRole");
    if (token && storedRole) {
      try {
        const decodedToken = jwtDecode(token);
        const now = Math.floor(Date.now() / 1000);
        if (
          decodedToken.exp < now ||
          !ALLOWED_ADMIN_ROLES.includes(storedRole) ||
          decodedToken.role !== storedRole
        ) {
          clearAdminSessionAndRedirect();
        } else {
          setRole(storedRole);
        }
      } catch (err) {
        clearAdminSessionAndRedirect();
      }
    } else {
      clearAdminSessionAndRedirect();
    }
    setLoading(false);
  }, [clearAdminSessionAndRedirect]);

  useEffect(() => {
    const localDarkMode = localStorage.getItem("darkMode") === "true";
    setDarkMode(localDarkMode);
    if (localDarkMode) document.documentElement.classList.add("dark");

    checkAdminAuthAndSetRole();
    const handleAuthChange = () => {
      setLoading(true);
      checkAdminAuthAndSetRole();
    };
    window.addEventListener("authChangeAdmin", handleAuthChange);

    // Listener để thay đổi sidebar khi resize cửa sổ
    const handleResize = () => {
      if (window.innerWidth < 1024) {
        setSidebarToggle(false);
      } else {
        setSidebarToggle(true);
      }
    };
    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("authChangeAdmin", handleAuthChange);
      window.removeEventListener("resize", handleResize);
    };
  }, [checkAdminAuthAndSetRole]);

  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("darkMode", "true");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("darkMode", "false");
    }
  }, [darkMode]);

  if (loading) return <Preloader />;
  if (!role) return null;

  return (
    <div className="dark:bg-gray-900">
      <div className="flex h-screen overflow-hidden">
        {/* Render Sidebar có điều kiện */}
        {sidebarToggle && <Sidebar userRole={role} />}

        <div className="relative flex flex-1 flex-col overflow-y-auto overflow-x-hidden">
          <Header
            sidebarToggle={sidebarToggle}
            setSidebarToggle={setSidebarToggle}
            darkMode={darkMode}
            setDarkMode={setDarkMode}
            role={role}
            onLogout={clearAdminSessionAndRedirect}
          />
          <main>
            <div className="mx-auto max-w-screen-2xl p-4 md:p-6 2xl:p-10">
              <Outlet />
            </div>
          </main>
        </div>
      </div>
    </div>
  );
};

export default DashboardLayout;
