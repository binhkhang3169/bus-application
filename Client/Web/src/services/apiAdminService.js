import axios from "axios";
import { format } from "date-fns";
import { API_URL } from "../configs/env";

const apiClient = axios.create({
  baseURL: API_URL + "api/v1", // Địa chỉ backend Go của bạn
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    Authorization: `Bearer ${sessionStorage.getItem("adminAccessToken") || ""}`, // Lấy token từ sessionStorage
  },
});
// Hàm định dạng ngày tháng sang 'YYYY-MM-DD'
const formatDateForApi = (date) => format(new Date(date), "yyyy-MM-dd");

// Lấy các chỉ số KPI
export const getKpis = (startDate, endDate) => {
  const params = {
    start_date: formatDateForApi(startDate),
    end_date: formatDateForApi(endDate),
  };
  return apiClient.get("/kpis", { params });
};

// Lấy dữ liệu doanh thu theo thời gian
export const getRevenueOverTime = (startDate, endDate, groupBy = "day") => {
  const params = {
    start_date: formatDateForApi(startDate),
    end_date: formatDateForApi(endDate),
    group_by: groupBy,
  };
  return apiClient.get("/charts/revenue-over-time", { params });
};

// Lấy dữ liệu phân bổ vé
export const getTicketDistribution = (startDate, endDate) => {
  const params = {
    start_date: formatDateForApi(startDate),
    end_date: formatDateForApi(endDate),
  };
  return apiClient.get("/charts/ticket-distribution", { params });
};

// Lấy dữ liệu xu hướng tìm kiếm theo thời gian
export const getSearchesOverTime = (startDate, endDate) => {
  const params = {
    start_date: formatDateForApi(startDate),
    end_date: formatDateForApi(endDate),
  };
  return apiClient.get("/analytics/searches/over-time", { params });
};

// Lấy top các tỉnh được tìm kiếm
export const getTopProvinces = (startDate, endDate, limit = 5) => {
  const params = {
    start_date: formatDateForApi(startDate),
    end_date: formatDateForApi(endDate),
    limit,
  };
  return apiClient.get("/analytics/searches/top-provinces", { params });
};
