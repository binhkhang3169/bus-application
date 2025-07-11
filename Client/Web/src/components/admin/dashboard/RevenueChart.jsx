import React, { useState, useEffect } from 'react';
import ReactApexChart from 'react-apexcharts';
import { getRevenueOverTime } from '../../../services/apiAdminService';

const RevenueChart = ({ darkMode }) => {
    const [series, setSeries] = useState([{ name: 'Doanh thu', data: [] }]);
    const [categories, setCategories] = useState([]);
    const [groupBy, setGroupBy] = useState('month'); // 'day', 'week', 'month'
    const [loading, setLoading] = useState(true);

    // Mặc định 1 năm gần nhất
    const defaultStartDate = new Date(new Date().setFullYear(new Date().getFullYear() - 1));
    const defaultEndDate = new Date();

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const response = await getRevenueOverTime(defaultStartDate, defaultEndDate, groupBy);
                const chartData = response.data || [];
                
                const newCategories = chartData.map(item => item.date);
                const newSeriesData = chartData.map(item => item.value);

                setCategories(newCategories);
                setSeries([{ name: 'Doanh thu', data: newSeriesData }]);
            } catch (error) {
                console.error("Không thể tải dữ liệu doanh thu:", error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [groupBy]);

    const chartOptions = {
        chart: { type: 'bar', height: 350, toolbar: { show: false }, foreColor: darkMode ? '#ccc' : '#333' },
        plotOptions: { bar: { horizontal: false, columnWidth: '55%', borderRadius: 4 } },
        dataLabels: { enabled: false },
        xaxis: { categories: categories, labels: { style: { colors: darkMode ? '#ccc' : '#333' } } },
        yaxis: { title: { text: 'Doanh thu (VND)', style: { color: darkMode ? '#ccc' : '#333' } }, labels: { style: { colors: darkMode ? '#ccc' : '#333' } } },
        grid: { borderColor: darkMode ? '#444' : '#e7e7e7' },
        tooltip: { theme: darkMode ? 'dark' : 'light', y: { formatter: val => new Intl.NumberFormat('vi-VN').format(val) + " VND" } },
        colors: ['#3C50E0']
    };

    return (
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-gray-900 sm:px-7.5 sm:pt-7.5">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">Doanh thu</h3>
                <div className="inline-flex items-center gap-1 rounded-lg bg-gray-100 p-1 dark:bg-gray-800">
                    <button onClick={() => setGroupBy('day')} className={`px-3 py-1 font-medium rounded text-sm ${groupBy === 'day' ? 'bg-white shadow dark:bg-gray-700' : ''}`}>Ngày</button>
                    <button onClick={() => setGroupBy('week')} className={`px-3 py-1 font-medium rounded text-sm ${groupBy === 'week' ? 'bg-white shadow dark:bg-gray-700' : ''}`}>Tuần</button>
                    <button onClick={() => setGroupBy('month')} className={`px-3 py-1 font-medium rounded text-sm ${groupBy === 'month' ? 'bg-white shadow dark:bg-gray-700' : ''}`}>Tháng</button>
                </div>
            </div>
             {loading ? <div className="h-[350px] flex items-center justify-center">Đang tải...</div> : 
                <ReactApexChart options={chartOptions} series={series} type="bar" height={350} />
            }
        </div>
    );
};

export default RevenueChart;