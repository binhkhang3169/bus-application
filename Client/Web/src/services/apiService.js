/** @format */

import axios from 'axios';
import { API_URL } from '../configs/env';

// Create an axios instance with default settings
const api = axios.create({
  baseURL: `${API_URL}api/v1/`, //api/v1/', // Base URL for your API
  withCredentials: true, // Crucial for sending cookies (like the refresh token)
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add auth token (access token) to requests if available
api.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem('accessToken') || sessionStorage.getItem('accessToken');
    const refreshToken = localStorage.getItem('refreshToken') || sessionStorage.getItem('refreshToken');

    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }

    if (refreshToken) {
      config.headers['X-Refresh-Token'] = refreshToken;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
// Variable to prevent multiple concurrent token refresh attempts
let isRefreshing = false;
// Array to hold requests that are waiting for a new token
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// Response interceptor - handle common errors, including token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Check if the error is 401 (Unauthorized) and it's not a retry request
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // If already refreshing, add the request to the queue
        return new Promise(function(resolve, reject) {
          failedQueue.push({resolve, reject});
        }).then(token => {
          originalRequest.headers['Authorization'] = 'Bearer ' + token;
          return axios(originalRequest); // Use global axios for retry to avoid interceptor loop
        }).catch(err => {
          return Promise.reject(err);
        });
      }

      originalRequest._retry = true; // Mark as a retry
      isRefreshing = true;

      try {
        console.log('Attempting to refresh token...');
        // The refresh token is in an HttpOnly cookie, so no need to send it in the body
        const rs = await api.post('auth/refresh-token');

        if (rs.status === 200 && rs.data.data.accessToken) {
          const newAccessToken = rs.data.data.accessToken;
          console.log('Token refreshed successfully:', newAccessToken);

          // Update token in storage (check if it was in localStorage or sessionStorage)
          if (localStorage.getItem('accessToken')) {
            localStorage.setItem('accessToken', newAccessToken);
          } else if (sessionStorage.getItem('accessToken')) {
            sessionStorage.setItem('accessToken', newAccessToken);
          } else {
            // If token wasn't stored, default to sessionStorage or handle as error
            sessionStorage.setItem('accessToken', newAccessToken);
          }

          api.defaults.headers.common['Authorization'] = `Bearer ${newAccessToken}`;
          originalRequest.headers['Authorization'] = `Bearer ${newAccessToken}`;
          processQueue(null, newAccessToken);
          return api(originalRequest); // Retry original request with new token
        } else {
           // If refresh fails with a valid response but no token (should not happen with your backend logic)
          console.error('Refresh token response was not 200 or did not contain accessToken.');
          clearAuthDataAndRedirect();
          processQueue(new Error('Refresh token failed, no new access token.'), null);
          return Promise.reject(error);
        }
      } catch (_error) {
        console.error('Error during token refresh:', _error.response?.data || _error.message);
        processQueue(_error, null);
        clearAuthDataAndRedirect();
        return Promise.reject(_error);
      } finally {
        isRefreshing = false;
      }
    }
    // For other errors, or if it's not a 401, just reject
    return Promise.reject(error);
  }
);

function clearAuthDataAndRedirect() {
  localStorage.removeItem('accessToken');
  sessionStorage.removeItem('accessToken');
  localStorage.removeItem('userRole'); // If you store role
  localStorage.removeItem('username'); // If you store username
  // Dispatch a custom event to notify other parts of the app about logout
  window.dispatchEvent(new Event('authChange'));
  // Optional: redirect to login
  // if (window.location.pathname !== '/login') {
  //   window.location.href = '/login';
  // }
  console.log('Cleared auth data and redirected (or would redirect).');
}

export default api;