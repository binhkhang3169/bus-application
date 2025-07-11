/** @format */
import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';
import { API_URL } from '../configs/env';

function ResultBookingPage() {
	const navigate = useNavigate();
	const location = useLocation();
	const [status, setStatus] = useState(null);
	const [paymentDetails, setPaymentDetails] = useState(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState(null);
	function formatVnpPayDate(dateStr) {
  if (!dateStr || dateStr.length !== 14) return dateStr;
  const year = dateStr.substring(0, 4);
  const month = dateStr.substring(4, 6);
  const day = dateStr.substring(6, 8);
  const hour = dateStr.substring(8, 10);
  const minute = dateStr.substring(10, 12);
  const second = dateStr.substring(12, 14);

  return `${hour}:${minute}:${second} ngày ${day}/${month}/${year}`;
}
	useEffect(() => {
		// Get all query parameters from the URL
		const queryParams = new URLSearchParams(location.search);
		
		// Check if we have query parameters (came from VNPay redirect)
		if (queryParams.toString()) {
			verifyPayment(queryParams);
		} else {
			// No query parameters, likely direct access to the page
			setLoading(false);
			setError('Không tìm thấy thông tin thanh toán');
		}
	}, [location.search]);

	// Verify payment with backend
	const verifyPayment = async (queryParams) => {
		try {
			// // Get token from sessionStorage
			// const token = localStorage.getItem('accessToken') || sessionStorage.getItem('accessToken');
			// if (!token) {
			// 	setError('Phiên đăng nhập đã hết hạn');
			// 	setLoading(false);
			// 	setStatus('failure');
			// 	return;
			// }

			// Send all query parameters to backend for verification
			const response = await axios.get(
				`${API_URL}api/v1/vnpay/return?${queryParams.toString()}`
			);
			console.log(response);


			// Check response
			if (response.data && response.data.data.is_valid_signature && response.data.data.vnp_response_code === '00') {
				setStatus('success');
				setPaymentDetails(response.data.data);
			} else {
				setStatus('failure');
				setError('Xác thực thanh toán thất bại');
			}
		} catch (error) {
			console.log(error);
			console.error('Error verifying payment:', error);
			setError('Lỗi xác thực thanh toán: ' + 
				(error.response ? error.response.data.message || JSON.stringify(error.response.data) : error.message));
			setStatus('failure');
		} finally {
			setLoading(false);
		}
	};

	if (loading) {
		return (
			<div className="flex flex-col items-center justify-center h-screen">
				<div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				<p className="mt-4 text-lg">Đang xác thực thanh toán...</p>
			</div>
		);
	}

	return (
		<div className="flex flex-col items-center justify-center min-h-screen py-12 px-4 bg-slate-50">
			<div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
				{status === 'success' ? (
					<div className="text-center">
						<div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100">
							<svg 
								className="h-10 w-10 text-green-600" 
								fill="none" 
								stroke="currentColor" 
								viewBox="0 0 24 24" 
								xmlns="http://www.w3.org/2000/svg"
							>
								<path 
									strokeLinecap="round" 
									strokeLinejoin="round" 
									strokeWidth="2" 
									d="M5 13l4 4L19 7"
								></path>
							</svg>
						</div>
						<h2 className="mt-4 text-2xl font-bold text-gray-900">Thanh toán thành công</h2>
						
						{paymentDetails && (
							<div className="mt-6 text-left">
								<div className="border-t border-b py-4">
									<div className="flex justify-between items-center mb-2">
										<span className="text-gray-600">Mã giao dịch:</span>
										<span className="font-medium">{paymentDetails.vnp_transaction_no}</span>
									</div>
									<div className="flex justify-between items-center mb-2">
										<span className="text-gray-600">Số tiền:</span>
										<span className="font-medium">{new Intl.NumberFormat('vi-VN').format(paymentDetails.amount)} đ</span>
									</div>
									<div className="flex justify-between items-center mb-2">
										<span className="text-gray-600">Ngân hàng:</span>
										<span className="font-medium">{paymentDetails.vnp_bank_code}</span>
									</div>
									<div className="flex justify-between items-center">
										<span className="text-gray-600">Thời gian:</span>
										<span className="font-medium">{formatVnpPayDate(paymentDetails.vnp_pay_date)}</span>
									</div>
								</div>
								<p className="mt-4 text-sm text-center text-gray-500">
									Chi tiết đặt vé đã được gửi đến email của bạn
								</p>
							</div>
						)}
					</div>
				) : (
					<div className="text-center">
						<div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-red-100">
							<svg 
								className="h-10 w-10 text-red-600" 
								fill="none" 
								stroke="currentColor" 
								viewBox="0 0 24 24" 
								xmlns="http://www.w3.org/2000/svg"
							>
								<path 
									strokeLinecap="round" 
									strokeLinejoin="round" 
									strokeWidth="2" 
									d="M6 18L18 6M6 6l12 12"
								></path>
							</svg>
						</div>
						<h2 className="mt-4 text-2xl font-bold text-gray-900">Thanh toán thất bại</h2>
						{error && <p className="mt-2 text-red-600">{error}</p>}
					</div>
				)}
				
				<div className="mt-8 text-center">
					<button
						onClick={() => navigate('/')}
						className="inline-flex justify-center items-center px-6 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
					>
						Trở về trang chủ
					</button>
				</div>
			</div>
		</div>
	); 	
}

export default ResultBookingPage;