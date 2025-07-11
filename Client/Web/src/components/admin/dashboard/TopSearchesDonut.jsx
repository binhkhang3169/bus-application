import React, { useState, useEffect } from "react";
import ReactApexChart from "react-apexcharts";
import { getTopProvinces } from "../../../services/apiAdminService";
import { API_URL } from "../../../configs/env";

const TopSearchesDonut = ({ darkMode }) => {
  const [series, setSeries] = useState([]);
  const [labels, setLabels] = useState([]);
  const [loading, setLoading] = useState(true);
  const [provinceMap, setProvinceMap] = useState({});

  // Function to fetch provinces and create province ID to name mapping
  const fetchProvinceMapping = async () => {
    try {
      const response = await fetch(API_URL + "api/v1/provinces");
      const data = await response.json();

      if (data.code === 200 && data.data) {
        const mapping = {};

        // Create mapping from province ID to name
        data.data.forEach((province) => {
          mapping[province.id] = province.name;
        });

        setProvinceMap(mapping);
        return mapping;
      }
    } catch (error) {
      console.error("Không thể tải dữ liệu mapping tỉnh:", error);
      return {};
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        // First fetch province mapping
        const mapping = await fetchProvinceMapping();

        // Then fetch top provinces data
        const startDate = new Date(
          new Date().setDate(new Date().getDate() - 30)
        );
        const endDate = new Date();
        const response = await getTopProvinces(startDate, endDate, 5);
        const topOrigins = response.data.top_origins || [];

        setSeries(topOrigins.map((item) => item.search_count));
        console.log("Top Origins Data:", mapping);
        // Use province names instead of IDs
        setLabels(
          topOrigins.map((item) => {
            const provinceName = mapping[item.province_id];
            return provinceName || `Tỉnh ${item.province_id}`;
          })
        );
      } catch (error) {
        console.error("Không thể tải dữ liệu top tỉnh tìm kiếm:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const chartOptions = {
    chart: {
      type: "donut",
      foreColor: darkMode ? "#ccc" : "#333",
      background: "transparent",
    },
    labels: labels,
    dataLabels: {
      enabled: true,
      formatter: (val, opts) => {
        return opts.w.config.series[opts.seriesIndex];
      },
    },
    plotOptions: {
      pie: {
        donut: {
          labels: {
            show: true,
            total: {
              show: true,
              label: "Tổng lượt",
              fontSize: "16px",
              fontWeight: 600,
              color: darkMode ? "#fff" : "#333",
            },
          },
        },
      },
    },
    legend: {
      position: "bottom",
      fontSize: "14px",
      fontWeight: 400,
      labels: {
        colors: darkMode ? "#ccc" : "#333",
      },
    },
    tooltip: {
      theme: darkMode ? "dark" : "light",
      y: {
        formatter: function (val) {
          return val + " lượt tìm kiếm";
        },
      },
    },
    colors: ["#3C50E0", "#6577F3", "#8092FF", "#A8B3FF", "#D0D7FF"],
    responsive: [
      {
        breakpoint: 480,
        options: {
          chart: {
            width: 200,
          },
          legend: {
            position: "bottom",
          },
        },
      },
    ],
  };

  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-gray-900">
      <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">
        Top 5 Tỉnh xuất phát
      </h3>
      <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
        Các tỉnh được tìm kiếm làm điểm đi nhiều nhất (30 ngày qua)
      </p>
      <div className="my-4 h-[280px]">
        {loading ? (
          <div className="h-full flex items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span className="ml-2 text-gray-600 dark:text-gray-400">
              Đang tải...
            </span>
          </div>
        ) : series.length > 0 ? (
          <ReactApexChart
            options={chartOptions}
            series={series}
            type="donut"
            height={280}
          />
        ) : (
          <div className="h-full flex items-center justify-center">
            <p className="text-gray-500 dark:text-gray-400">
              Không có dữ liệu để hiển thị
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default TopSearchesDonut;
