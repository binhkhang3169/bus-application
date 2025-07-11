/** @format */

import React, { useEffect, useState } from 'react';
import { useNavigate }  from 'react-router-dom';
// Adjust this path to correctly point to your apiService.js
// Assuming it's in a 'services' folder, one level up from 'pages' (if this is in a pages folder)
// For example: import api from '../../services/apiService';
import api from '../../services/apiService'; // <--- !!! ADJUST THIS IMPORT PATH !!!

const ALLOWED_ADMIN_ROLES = ['ROLE_ADMIN', 'ROLE_RECEPTION', 'ROLE_DRIVER', 'ROLE_OPERATOR'];

// Helper function to decode JWT - use cautiously and prefer server-side validation
const decodeJwt = (token) => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(function (c) {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        })
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (e) {
    console.error("Failed to decode JWT:", e);
    return null;
  }
};


function LoginPageAdmin() {
  const navigate = useNavigate();

  const [message, setMessage] = useState('');
  const [username, setUsername] = useState(''); // Renamed from setUserName for consistency
  const [password, setPassword] = useState('');

  useEffect(() => {
    const adminToken = sessionStorage.getItem('adminAccessToken');
    const adminRole = sessionStorage.getItem('adminRole');

    if (adminToken && adminRole) {
      const payload = decodeJwt(adminToken);
      if (payload) {
        const now = Math.floor(Date.now() / 1000);
        const isExpired = payload.exp < now;
        const isRoleAllowed = ALLOWED_ADMIN_ROLES.includes(adminRole) && ALLOWED_ADMIN_ROLES.includes(payload.role); // Double check role

        if (isExpired || !isRoleAllowed) {
          sessionStorage.removeItem('adminAccessToken');
          sessionStorage.removeItem('adminUsername');
          sessionStorage.removeItem('adminRole');
          // No navigation here, let it stay on login page or redirect if this page itself requires no token
        } else {
          console.log('Admin token and role valid, navigating to admin home.');
          navigate('/admin/home'); // Or your main admin dashboard route
        }
      } else {
        // Invalid token format
        sessionStorage.removeItem('adminAccessToken');
        sessionStorage.removeItem('adminUsername');
        sessionStorage.removeItem('adminRole');
      }
    }
  }, [navigate]);

  const submitLogin = async () => {
    if (!username || !password) {
      setMessage('Vui lòng nhập đầy đủ thông tin.');
      return;
    }

    const loginData = {
      username,
      password,
      rememberMe: false, // Admin login typically doesn't use "remember me" for long sessions
    };

    try {
      setMessage(''); // Clear previous messages
      // The baseURL in apiService is 'http://localhost:8080/v1/'
      // The login endpoint is 'auth/login'
      const response = await api.post('auth/login', loginData);

      console.log('Admin Login Response:', response);

      if (response.status === 200 && response.data.code === 200) {
        const accessToken = response.data.data.accessToken

        const {username, role } = response.data.data.user;

        if (ALLOWED_ADMIN_ROLES.includes(role)) {
          sessionStorage.setItem('adminAccessToken', accessToken);
          sessionStorage.setItem('adminUsername', username);
          sessionStorage.setItem('adminRole', role);

          // Dispatch an authChange event if other parts of the admin UI need to react
          window.dispatchEvent(new CustomEvent('authChangeAdmin', { detail: { loggedIn: true } }));
          
          navigate('/admin/home'); // Or your main admin dashboard route
        } else {
          setMessage('Vai trò của bạn không được phép truy cập trang quản trị.');
        }
      } else {
        // Handle cases where status is 200 but backend indicates an error via 'code' or 'message'
        setMessage(response.data.message || 'Đăng nhập không thành công. Vui lòng thử lại.');
      }
    } catch (err) {
      console.error('Admin Login error:', err);
      if (err.response && err.response.data && err.response.data.message) {
        setMessage(err.response.data.message);
      } else if (err.message) {
        setMessage(err.message);
      } else {
        setMessage('Đã xảy ra lỗi không xác định. Vui lòng thử lại.');
      }
    }
  };

  return (
    <div className="text-white h-[100vh] flex flex-row justify-center items-center bg-center bg-cover bg-[url('https://images.unsplash.com/photo-1557223563-8db8e5e7d90b?q=80&w=2071&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')]">
      <div className="basis-10/12 sm:basis-8/12 md:basis-6/12 lg:basis-5/12 xl:basis-4/12 2xl:basis-3/12 mb-40">
        <div className="bg-slate-800 border border-slate-400 rounded-md p-8 shadow-lg backdrop-filter backdrop-blur-sm bg-opacity-30 relative">
          <h1 className="text-4xl text-white font-bold text-center mb-2">LOGIN</h1>
          <h2 className="text-md text-white text-center mb-4">EMPLOYEE</h2>

          <div>
            <div className="relative mb-8">
              <input
                type="text"
                className="block w-full py-2 px-0 text-sm text-white bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0 peer"
                placeholder=" " // Important for label animation
                autoComplete="username"
                value={username}
                onChange={(e) => {
                  setUsername(e.target.value);
                  setMessage('');
                }}
              />
              <label
                htmlFor="adminUsername" // Good practice to have 'for' attribute
                className="absolute text-md text-white duration-300 transform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-blue-400 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Tên đăng nhập / Số điện thoại
              </label>
              <span className="absolute top-1 right-1 ">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="w-6 h-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 1.5H8.25A2.25 2.25 0 0 0 6 3.75v16.5a2.25 2.25 0 0 0 2.25 2.25h7.5A2.25 2.25 0 0 0 18 20.25V3.75a2.25 2.25 0 0 0-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3" />
                </svg>
              </span>
            </div>
            <div className="relative mb-2">
              <input
                type="password"
                className="block w-full py-2 px-0 text-sm text-white bg-transparent border-0 border-b-2 border-gray-300 appearance-none peer focus:outline-none focus:ring-0 focus:border-blue-400"
                placeholder=" " // Important for label animation
                autoComplete="current-password"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                  setMessage('');
                }}
                onKeyDown={(event) => {
                  if (event.key === 'Enter') {
                    submitLogin();
                  }
                }}
              />
              <label
                htmlFor="adminPassword" // Good practice to have 'for' attribute
                className="absolute text-md text-white duration-300 transform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-blue-400 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
              >
                Mật khẩu
              </label>
              <span className="absolute top-1 right-1 ">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="w-6 h-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z" />
                </svg>
              </span>
            </div>
            {message && <div className="text-red-400 mt-4 text-sm font-medium text-center p-2 bg-red-900 bg-opacity-50 rounded">{message}</div>}

            <button
              type="button" // Changed from submit to button
              className="transition-colors duration-300 w-full mb-4 text-[18px] mt-8 rounded-full bg-white text-slate-800 hover:bg-cyan-500 hover:text-white py-2 font-semibold focus:outline-none focus:ring-2 focus:ring-cyan-400 focus:ring-opacity-75"
              onClick={submitLogin}
            >
              Đăng nhập
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default LoginPageAdmin;