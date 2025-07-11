// src/components/admin/dashboard/ChartOne.jsx
import React, { useState } from 'react';
import ReactApexChart from 'react-apexcharts';

const ChartOne = ({ darkMode }) => {
    const [chartData] = useState({
        series: [{
            name: 'Doanh thu',
            data: [2300, 1100, 2200, 2700, 1300, 2200, 3700, 2100, 4400, 2200, 3000, 4500]
        }],
        options: {
            chart: {
                type: 'bar',
                height: 350,
                toolbar: { show: false },
                foreColor: darkMode ? '#ccc' : '#333'
            },
            plotOptions: {
                bar: {
                    horizontal: false,
                    columnWidth: '55%',
                    endingShape: 'rounded'
                },
            },
            dataLabels: { enabled: false },
            stroke: {
                show: true,
                width: 2,
                colors: ['transparent']
            },
            xaxis: {
                categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
                labels: {
                    style: {
                        colors: darkMode ? '#ccc' : '#333'
                    }
                }
            },
            yaxis: {
                title: {
                    text: '$ (thousands)',
                    style: {
                        color: darkMode ? '#ccc' : '#333'
                    }
                },
                labels: {
                    style: {
                        colors: darkMode ? '#ccc' : '#333'
                    }
                }
            },
            grid: {
                borderColor: darkMode ? '#444' : '#e7e7e7',
            },
            fill: { opacity: 1 },
            tooltip: {
                theme: darkMode ? 'dark' : 'light',
                y: {
                    formatter: function (val) {
                        return "$ " + val + " thousands"
                    }
                }
            },
            colors: ['#3C50E0']
        }
    });

    return (
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-gray-900 sm:px-7.5 sm:pt-7.5">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">Doanh thu hàng tháng</h3>
            </div>
            <div id="chartOne">
                <ReactApexChart options={chartData.options} series={chartData.series} type="bar" height={350} />
            </div>
        </div>
    );
};

export default ChartOne;