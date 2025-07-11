import React, { useState, useEffect } from "react";

// 1. Import các component động đã được làm lại
import KpiCards from "../components/admin/dashboard/KpiCards";
import TrendsChart from "../components/admin/dashboard/TrendsChart";
import RevenueChart from "../components/admin/dashboard/RevenueChart";
import TopSearchesDonut from "../components/admin/dashboard/TopSearchesDonut";

const DashboardPage = () => {
  // Logic lấy trạng thái dark mode của bạn được giữ nguyên
  const [darkMode, setDarkMode] = useState(false);
  useEffect(() => {
    const isDark = localStorage.getItem("darkMode") === "true";
    setDarkMode(isDark);

    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "class") {
          const isDarkNow = document.documentElement.classList.contains("dark");
          setDarkMode(isDarkNow);
        }
      });
    });
    observer.observe(document.documentElement, { attributes: true });
    return () => observer.disconnect();
  }, []);

  return (
    <>
      {/* 2. Sử dụng layout mới với các component động */}
      <div className="flex flex-col gap-4 md:gap-6 2xl:gap-7.5">
        {/* Hàng 1: Các thẻ chỉ số KPI chính */}
        {/* KpiCards đã có bộ lọc ngày riêng nên không cần truyền props */}
        <KpiCards />

        {/* Hàng 2: Biểu đồ xu hướng chính (Tìm kiếm & Doanh thu) */}
        {/* TrendsChart có bộ lọc ngày và tab riêng */}
        <TrendsChart darkMode={darkMode} />

        {/* Hàng 3: Các biểu đồ phụ, sắp xếp cạnh nhau */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6 xl:grid-cols-3 2xl:gap-7.5">
          {/* Biểu đồ doanh thu theo Ngày/Tuần/Tháng */}
          <div className="col-span-1 md:col-span-2">
            <RevenueChart darkMode={darkMode} />
          </div>

          {/* Biểu đồ tròn cho Top Tỉnh tìm kiếm */}
          <div className="col-span-1">
            <TopSearchesDonut darkMode={darkMode} />
          </div>
        </div>
      </div>
    </>
  );
};

export default DashboardPage;
