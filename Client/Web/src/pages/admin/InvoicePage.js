import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/apiService';

const InvoicePage = () => {
    const { shipmentId } = useParams();
    const navigate = useNavigate();
    const [invoice, setInvoice] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchInvoice = async () => {
            try {
                const response = await api.get(`/api/v1/shipments/${shipmentId}/invoice`);
                setInvoice(response.data?.data || response.data);
            } catch (err) {
                if (err.response?.status === 404) {
                    setError('Invoice not found for this shipment.');
                } else {
                    setError('Failed to load invoice details.');
                }
                console.error("Fetch invoice error:", err);
            } finally {
                setLoading(false);
            }
        };

        fetchInvoice();
    }, [shipmentId]);

    if (loading) {
        return <div className="flex justify-center items-center h-screen"><div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-blue-500"></div></div>;
    }

    if (error) {
        return (
            <div className="text-center mt-10 p-4">
                <h2 className="text-2xl font-bold text-red-600">Error</h2>
                <p className="text-gray-600 mt-2">{error}</p>
                <button onClick={() => navigate(-1)} className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">Go Back</button>
            </div>
        );
    }
    
    if (!invoice) {
        return <div className="text-center mt-10 p-4">No invoice data.</div>
    }

    return (
        <div className="max-w-4xl mx-auto p-4 sm:p-6 lg:p-8">
            <div className="bg-white shadow-lg rounded-lg p-8">
                <div className="flex justify-between items-center border-b pb-4 mb-6">
                    <h1 className="text-3xl font-bold text-gray-800">Invoice</h1>
                    <div className="text-right">
                        <p className="text-lg font-semibold">Invoice #{invoice.id}</p>
                        <p className="text-sm text-gray-500">Date: {new Date(invoice.created_at).toLocaleDateString()}</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
                    <div>
                        <h2 className="text-xl font-semibold mb-2 text-gray-700">Shipment Details:</h2>
                        <p><strong>Shipment ID:</strong> {invoice.shipment_id}</p>
                        <p><strong>Item:</strong> {invoice.shipment?.item_name}</p>
                        <p><strong>Sender:</strong> {invoice.shipment?.sender_name}</p>
                        <p><strong>Receiver:</strong> {invoice.shipment?.receiver_name}</p>
                    </div>
                    <div className="text-right">
                         <h2 className="text-xl font-semibold mb-2 text-gray-700">Trip Info:</h2>
                         <p><strong>Route:</strong> {invoice.shipment?.trip ? `${invoice.shipment.trip.route.start.name} → ${invoice.shipment.trip.route.end.name}` : 'N/A'}</p>
                         <p><strong>Date:</strong> {invoice.shipment?.trip?.departureDate}</p>
                    </div>
                </div>

                <div className="mb-8">
                    <h3 className="text-xl font-semibold border-b pb-2 mb-4">Invoice Summary</h3>
                    <div className="flex justify-between items-center mb-2">
                        <span className="text-gray-600">Base Fare:</span>
                        <span className="font-medium">{new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(invoice.base_fare)}</span>
                    </div>
                    <div className="flex justify-between items-center mb-2">
                        <span className="text-gray-600">Surcharge:</span>
                        <span className="font-medium">{new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(invoice.surcharge)}</span>
                    </div>
                    <div className="flex justify-between items-center text-xl font-bold border-t pt-4 mt-4">
                        <span>Total Amount:</span>
                        <span>{new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(invoice.total_amount)}</span>
                    </div>
                </div>
                 <div className="flex justify-end">
                    <button onClick={() => navigate(-1)} className="mt-4 px-6 py-2 bg-gray-600 text-white rounded hover:bg-gray-700">
                        Back to Management
                    </button>
                 </div>
            </div>
        </div>
    );
};

export default InvoicePage;