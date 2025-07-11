// src/components/admin/dashboard/ChartThree.jsx
import React, { useState, useEffect } from 'react';
import ReactApexChart from 'react-apexcharts';
import DatePicker from 'react-datepicker';
import "react-datepicker/dist/react-datepicker.css";

// Mock data for different tabs
const overviewData = { name: "Tổng quan", data: [180, 210, 250, 190, 430, 310, 270, 400] };
const salesData = { name: "Đơn hàng", data: [40, 90, 120, 80, 150, 110, 95, 130] };
const revenueData = { name: "Doanh thu ($k)", data: [76, 85, 101, 98, 87, 105, 91, 114] };

const ChartThree = ({ darkMode }) => {
    const [selectedTab, setSelectedTab] = useState('Overview');
    const [startDate, setStartDate] = useState(new Date(new Date().setMonth(new Date().getMonth() - 1)));
    const [endDate, setEndDate] = useState(new Date());

    const [chartData, setChartData] = useState({
        series: [overviewData],
        options: { /* Initial options */ }
    });

    useEffect(() => {
        let newData;
        switch (selectedTab) {
            case 'Sales': newData = salesData; break;
            case 'Revenue': newData = revenueData; break;
            default: newData = overviewData;
        }

        setChartData({
            series: [newData],
            options: {
                chart: {
                    type: 'area',
                    height: 350,
                    toolbar: { show: false },
                    foreColor: darkMode ? '#ccc' : '#333'
                },
                dataLabels: { enabled: false },
                stroke: { curve: 'smooth', width: 2 },
                xaxis: {
                    type: 'datetime',
                    categories: ["2024-05-19T00:00:00.000Z", "2024-05-20T01:30:00.000Z", "2024-05-21T02:30:00.000Z", "2024-05-22T03:30:00.000Z", "2024-05-23T04:30:00.000Z", "2024-05-24T05:30:00.000Z", "2024-05-25T06:30:00.000Z", "2024-05-26T07:30:00.000Z"],
                    labels: { style: { colors: darkMode ? '#ccc' : '#333' } }
                },
                yaxis: { labels: { style: { colors: darkMode ? '#ccc' : '#333' } } },
                grid: { borderColor: darkMode ? '#444' : '#e7e7e7' },
                tooltip: {
                    theme: darkMode ? 'dark' : 'light',
                    x: { format: 'dd/MM/yy HH:mm' }
                },
                colors: ['#039855']
            }
        });
    }, [selectedTab, darkMode]);

    return (
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white px-5 pt-5 dark:border-gray-800 dark:bg-gray-900 sm:px-6 sm:pt-6">
            <div className="flex flex-col gap-5 mb-6 sm:flex-row sm:justify-between">
                <div className="w-full">
                    <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">Statistics</h3>
                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Thống kê chung trong khoảng thời gian đã chọn</p>
                </div>
                <div className="flex items-start w-full gap-3 sm:justify-end">
                    <div className="inline-flex w-fit items-center gap-0.5 rounded-lg bg-gray-100 p-0.5 dark:bg-gray-800">
                        <button onClick={() => setSelectedTab('Overview')} className={`px-3 py-2 font-medium rounded-md text-sm ${selectedTab === 'Overview' ? 'shadow-md text-gray-900 dark:text-white bg-white dark:bg-gray-700' : 'text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'}`}>Overview</button>
                        <button onClick={() => setSelectedTab('Sales')} className={`px-3 py-2 font-medium rounded-md text-sm ${selectedTab === 'Sales' ? 'shadow-md text-gray-900 dark:text-white bg-white dark:bg-gray-700' : 'text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'}`}>Sales</button>
                        <button onClick={() => setSelectedTab('Revenue')} className={`px-3 py-2 font-medium rounded-md text-sm ${selectedTab === 'Revenue' ? 'shadow-md text-gray-900 dark:text-white bg-white dark:bg-gray-700' : 'text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'}`}>Revenue</button>
                    </div>
                    <div className="relative">
                        <DatePicker
                            selected={startDate}
                            onChange={(dates) => {
                                const [start, end] = dates;
                                setStartDate(start);
                                setEndDate(end);
                            }}
                            startDate={startDate}
                            endDate={endDate}
                            selectsRange
                            className="h-10 w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-9 pr-4 text-sm font-medium text-gray-700 shadow-sm focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
                        />
                    </div>
                </div>
            </div>
            <div className="max-w-full overflow-x-auto">
                <div id="chartThree">
                    <ReactApexChart options={chartData.options} series={chartData.series} type="area" height={350} />
                </div>
            </div>
        </div>
    );
};

export default ChartThree;