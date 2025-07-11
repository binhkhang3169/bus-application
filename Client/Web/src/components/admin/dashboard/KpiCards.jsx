import React, { useState, useEffect } from 'react';
import { getKpis } from '../../../services/apiAdminService';
import DatePicker from 'react-datepicker';
import "react-datepicker/dist/react-datepicker.css";

const KpiCards = () => {
    const [kpiData, setKpiData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    
    // Mặc định 30 ngày gần nhất
    const [startDate, setStartDate] = useState(new Date(new Date().setDate(new Date().getDate() - 30)));
    const [endDate, setEndDate] = useState(new Date());

    useEffect(() => {
        const fetchKpis = async () => {
            try {
                setLoading(true);
                const response = await getKpis(startDate, endDate);
                setKpiData(response.data);
                setError(null);
            } catch (err) {
                setError('Không thể tải dữ liệu KPI.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        if (startDate && endDate) {
            fetchKpis();
        }
    }, [startDate, endDate]);

    const formatCurrency = (value) => {
        return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value);
    }

    if (loading) return <div className="text-center p-10">Đang tải dữ liệu...</div>;
    if (error) return <div className="text-center text-red-500 p-10">{error}</div>;

    return (
        <div>
            <div className="flex justify-end mb-4">
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
                    className="h-10 w-full max-w-xs rounded-lg border border-gray-200 bg-white py-2.5 px-4 text-sm font-medium text-gray-700 shadow-sm focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
                />
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 md:gap-6">
                <KpiCard title="Tổng Doanh thu" value={formatCurrency(kpiData?.total_revenue || 0)} />
                <KpiCard title="Tổng Vé đã bán" value={kpiData?.total_tickets?.toLocaleString('vi-VN') || '0'} />
                <KpiCard title="Tổng Hóa đơn" value={kpiData?.total_invoices?.toLocaleString('vi-VN') || '0'} />
            </div>
        </div>
    );
};

const KpiCard = ({ title, value }) => (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-gray-900 md:p-6">
        <div className="mt-2">
            <span className="text-sm text-gray-500 dark:text-gray-400">{title}</span>
            <h4 className="mt-2 text-2xl font-bold text-gray-800 dark:text-white/90">{value}</h4>
        </div>
    </div>
);

export default KpiCards;