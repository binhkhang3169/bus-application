import React, { useState, useEffect } from 'react';
import ReactApexChart from 'react-apexcharts';
import DatePicker from 'react-datepicker';
import "react-datepicker/dist/react-datepicker.css";
import { getSearchesOverTime, getRevenueOverTime } from '../../../services/apiAdminService';

const TrendsChart = ({ darkMode }) => {
    const [selectedTab, setSelectedTab] = useState('Searches'); // 'Searches' or 'Revenue'
    const [startDate, setStartDate] = useState(new Date(new Date().setDate(new Date().getDate() - 30)));
    const [endDate, setEndDate] = useState(new Date());
    const [series, setSeries] = useState([{ name: '', data: [] }]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            if (!startDate || !endDate) return;
            setLoading(true);
            try {
                let response;
                let seriesName = '';
                let color = '#039855';

                if (selectedTab === 'Searches') {
                    response = await getSearchesOverTime(startDate, endDate);
                    seriesName = 'Lượt tìm kiếm';
                } else {
                    response = await getRevenueOverTime(startDate, endDate, 'day'); // Luôn lấy theo ngày cho biểu đồ này
                    seriesName = 'Doanh thu';
                    color = '#3C50E0';
                }
                
                const chartData = response.data || [];
                const formattedData = chartData.map(item => [new Date(item.date).getTime(), item.value]);

                setSeries([{ name: seriesName, data: formattedData, color: color }]);
            } catch (error) {
                console.error(`Không thể tải dữ liệu cho tab ${selectedTab}:`, error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [selectedTab, startDate, endDate]);

    const chartOptions = {
        chart: { type: 'area', height: 350, toolbar: { show: false }, foreColor: darkMode ? '#ccc' : '#333' },
        dataLabels: { enabled: false },
        stroke: { curve: 'smooth', width: 2 },
        xaxis: { type: 'datetime', labels: { style: { colors: darkMode ? '#ccc' : '#333' } } },
        yaxis: { labels: { style: { colors: darkMode ? '#ccc' : '#333' }, formatter: (val) => val.toLocaleString('vi-VN') } },
        grid: { borderColor: darkMode ? '#444' : '#e7e7e7' },
        tooltip: { theme: darkMode ? 'dark' : 'light', x: { format: 'dd/MM/yyyy' } },
        colors: [series.length > 0 ? series[0].color : '#039855']
    };

    return (
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white px-5 pt-5 dark:border-gray-800 dark:bg-gray-900 sm:px-6 sm:pt-6">
            <div className="flex flex-col gap-5 mb-6 sm:flex-row sm:justify-between">
                <div>
                    <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">Xu hướng theo thời gian</h3>
                </div>
                <div className="flex items-start w-full gap-3 sm:justify-end">
                    <div className="inline-flex w-fit items-center gap-0.5 rounded-lg bg-gray-100 p-0.5 dark:bg-gray-800">
                        <button onClick={() => setSelectedTab('Searches')} className={`px-3 py-2 font-medium rounded-md text-sm ${selectedTab === 'Searches' ? 'shadow-md text-gray-900 dark:text-white bg-white dark:bg-gray-700' : 'text-gray-500 hover:text-gray-900'}`}>Tìm kiếm</button>
                        <button onClick={() => setSelectedTab('Revenue')} className={`px-3 py-2 font-medium rounded-md text-sm ${selectedTab === 'Revenue' ? 'shadow-md text-gray-900 dark:text-white bg-white dark:bg-gray-700' : 'text-gray-500 hover:text-gray-900'}`}>Doanh thu</button>
                    </div>
                    <DatePicker
                        selected={startDate}
                        onChange={(dates) => { const [start, end] = dates; setStartDate(start); setEndDate(end); }}
                        startDate={startDate}
                        endDate={endDate}
                        selectsRange
                        className="h-10 w-full max-w-xs rounded-lg border border-gray-200 bg-white py-2.5 px-4 text-sm font-medium text-gray-700 shadow-sm focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
                    />
                </div>
            </div>
            <div className="max-w-full overflow-x-auto">
                 {loading ? <div className="h-[350px] flex items-center justify-center">Đang tải...</div> : 
                    <ReactApexChart options={chartOptions} series={series} type="area" height={350} />
                }
            </div>
        </div>
    );
};

export default TrendsChart;