/** @format */

import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { API_URL } from '../configs/env';
import FailureNotification from './Noti/FailureNotification';

function BookingTicketForm() {
	const [searchParams, setSearchParams] = useSearchParams();
	const navigate = useNavigate();

	// Trip ID
	let tripID = searchParams.get('id');

	// Modal
	const [failureModal, setFailureModal] = useState(false);
	const [message, setMessage] = useState('');

	// Input
	const [seat1, setSeat1] = useState([]);
	const [seat2, setSeat2] = useState([]);
	const [numberOfSeats, setNumberOfSeatsSelected] = useState(0);
	const [price, setPrice] = useState(0);
	const [start_time, setStartTime] = useState('');
	const [route, setRoute] = useState('');
	const [customer, setCustomer] = useState('');

	useEffect(() => {
		if (tripID) {
			getTrip();
			authCustomer();
		} else {
			navigate('/');
		}
	}, []);

	// Send POST request to payment
	const payment = async () => {
		if (!numberOfSeats) {
			setMessage('Bạn chưa chọn ghế');
			openFailureModal();
			return;
		}
		let body = {
			trip_id: tripID,
			customer_id: customer.id,
			seat: seatID(),
			discount: '',
			price: price,
			quantity: numberOfSeats,
			amount: price * numberOfSeats,
			language: 'vn',
			bankCode: '',
		};
		authCustomer();

		const token = sessionStorage.getItem('token');
		await axios
			.post(API_URL + `customer/payment`, body, { headers: { Authorization: `Bearer ${token}` } })
			.then((res) => {
				if (res.status === 200) {
					window.location.href = res.data;
				}
			})
			.catch((err) => {
				setMessage(err.response.data.message);
				openFailureModal();
			});
	};

	// Send GET request to retrieve trip information by id
	const getTrip = async () => {
		await axios.get(API_URL + `employee/trip/${tripID}`).then((res) => {
			setPrice(res.data.trip.price);
			console.log(res.data.trip.date);
			setStartTime(
				res.data.trip.start_time.substring(0, res.data.trip.start_time.length - 3) + ' ' + res.data.trip.date
			);
			setRoute(res.data.trip.route.name);

			let newSeat1 = new Array(18).fill('0');
			let newSeat2 = new Array(18).fill('0');

			for (let i = 1; i < 37; i++) {
				if (res.data.trip.seat.includes(i.toString())) {
					if (parseInt(i) > 18) {
						newSeat2[i - 19] = '2';
					} else {
						newSeat1[i - 1] = '2';
					}
				}
			}

			setSeat1(newSeat1);
			setSeat2(newSeat2);
		});
	};

	// User authentication
	const authCustomer = async () => {
		const token = sessionStorage.getItem('token');
		if (token) {
			axios
				.get(API_URL + 'customer/me', { headers: { Authorization: `Bearer ${token}` } })
				.then((res) => {
					setCustomer(res.data.customer);
				})
				.catch((err) => {
					navigate('/login');
				});
		} else {
			navigate('/login');
		}
	};

	// Choose seats
	const chooseSeat = (seat) => {
		console.log(start_time);
		let newSeat = [];
		if (seat > 18) {
			newSeat = [...seat2];
			if (newSeat[seat - 19] === '0') {
				newSeat[seat - 19] = '1';
				setNumberOfSeatsSelected(numberOfSeats + 1);
			} else if (newSeat[seat - 19] === '1') {
				newSeat[seat - 19] = '0';
				setNumberOfSeatsSelected(numberOfSeats - 1);
			}
			setSeat2(newSeat);
		} else {
			newSeat = [...seat1];
			if (newSeat[seat - 1] === '0') {
				newSeat[seat - 1] = '1';
				setNumberOfSeatsSelected(numberOfSeats + 1);
			} else if (newSeat[seat - 1] === '1') {
				newSeat[seat - 1] = '0';
				setNumberOfSeatsSelected(numberOfSeats - 1);
			}
			setSeat1(newSeat);
		}
	};

	// Return seat id to
	const seatID = () => {
		let res = '';
		seat1.map((v, i) => {
			if (v === '1') {
				res += i + 1 + ',';
			}
		});
		seat2.map((v, i) => {
			if (v === '1') {
				res += i + 19 + ',';
			}
		});
		return res.substring(0, res.length - 1);
	};

	// Close Failure Modal
	const closeFailureModal = () => {
		setFailureModal(false);
	};

	// Open Failure Modal
	const openFailureModal = () => {
		setFailureModal(true);
	};

	return (
		<div className="w-full bg-slate-50 lg:p-10">
			<div className="flex max-w-screen-lg flex-col md:flex-row mx-auto gap-8">
				<div className="basis-2/3 w-full flex flex-col gap-5">
					<div className="seat-section bg-white rounded-t-xl p-5 border border-slate-200">
						<h3 className="text-xl font-medium">Chọn ghế</h3>
						<div className="flex flex-row mt-5 gap-8">
							<div className="grid grid-cols-3 gap-2 cursor-pointer">
								{seat1.map((v, i) => {
									return (
										<div
											key={i}
											className="mt-1 text-center relative flex justify-center"
											onClick={() => chooseSeat(i + 1)}
										>
											<img
												width="32"
												src={
													v === '0'
														? 'https://futabus.vn/images/icons/seat_active.svg'
														: v === '2'
														? 'https://futabus.vn/images/icons/seat_disabled.svg'
														: 'https://futabus.vn/images/icons/seat_selecting.svg'
												}
												alt="seat icon"
											/>
											<span
												className={`absolute text-sm font-semibold sm:text-[12px] ${
													v === '2'
														? 'text-gray-400'
														: v === '0'
														? 'text-blue-400'
														: 'text-red-400'
												} top-1`}
											>
												{i + 1}
											</span>
										</div>
									);
								})}
							</div>

							<div className="grid grid-cols-3 gap-2 cursor-pointer">
								{seat2.map((v, i) => {
									return (
										<div
											key={i}
											className="mt-1 text-center relative flex justify-center"
											onClick={() => chooseSeat(i + 19)}
										>
											<img
												width="32"
												src={
													v === '0'
														? 'https://futabus.vn/images/icons/seat_active.svg'
														: v === '2'
														? 'https://futabus.vn/images/icons/seat_disabled.svg'
														: 'https://futabus.vn/images/icons/seat_selecting.svg'
												}
												alt="seat icon"
											/>
											<span
												className={`absolute text-sm font-semibold sm:text-[12px] ${
													v === '2'
														? 'text-gray-400'
														: v === '0'
														? 'text-blue-400'
														: 'text-red-400'
												} top-1`}
											>
												{i + 19}
											</span>
										</div>
									);
								})}
							</div>

							<div className="flex flex-col gap-3 text-sm">
								<span className="flex items-center">
									<div className="mr-2 h-4 w-4 rounded bg-[#D5D9DD] border-[#C0C6CC]"></div>
									Đã bán
								</span>
								<span className="flex items-center">
									<div className="mr-2 h-4 w-4 rounded bg-[#DEF3FF] border-[#96C5E7]"></div>
									Còn trống
								</span>
								<span className="flex items-center">
									<div className="mr-2 h-4 w-4 rounded bg-[#FDEDE8] border-[#F8BEAB]"></div>
									Đang chọn
								</span>
							</div>
						</div>
					</div>
					<div className="user-info-section bg-white border border-slate-200 p-5">
						<div className="flex flex-row gap-8">
							<div className="basis-1/2 w-full">
								<h3 className="text-xl font-semibold">Thông tin khách hàng</h3>
								<form
									action=""
									className="mt-5"
								>
									<div>
										<div className="mb-5">
											<p className="text-sm mb-1">
												Họ và tên <span className="text-red-500">*</span>
											</p>
											<input
												type="text"
												name="fullname"
												className="rounded-xl border border-slate-200 w-full bg-gray-200 text-gray-500"
												value={`${customer.last_name} ${customer.first_name}`}
												disabled
											/>
										</div>
										<div className="mb-5">
											<p className="text-sm mb-1">
												Số điện thoại <span className="text-red-500">*</span>
											</p>
											<input
												disabled
												type="text"
												name="fullname"
												className="rounded-xl border border-slate-200 w-full bg-gray-200 text-gray-500"
												value={customer.phone_number}
											/>
										</div>
										<div>
											<p className="text-sm mb-1">
												Email <span className="text-red-500">*</span>
											</p>
											<input
												disabled
												type="text"
												name="fullname"
												className="rounded-xl border border-slate-200 w-full bg-gray-200 text-gray-500"
												value={customer.email}
											/>
										</div>
									</div>
								</form>
							</div>
							<div className="basis-1/2 w-full">
								<h3 className="text-red-500 text-center mb-5 font-semibold">ĐIỀU KHOẢN & LƯU Ý</h3>
								<p className="text-[15px] text-justify mb-3 font-[500] leading-6">
									(*) Quý khách vui lòng có mặt tại bến xuất phát của xe trước ít nhất 30 phút giờ xe
									khởi hành, mang theo thông báo đã thanh toán vé thành công có chứa mã vé được gửi từ
									hệ thống FUTA BUS LINE. Vui lòng liên hệ Trung tâm tổng đài{' '}
									<span className="text-red-500">1900 6067</span> để được hỗ trợ.
								</p>
								<p className="text-[15px] text-justify font-[500] leading-6">
									(*) Nếu quý khách có nhu cầu trung chuyển, vui lòng liên hệ Tổng đài trung chuyển{' '}
									<span className="text-red-500">1900 6918</span> trước khi đặt vé. Chúng tôi không
									đón/trung chuyển tại những điểm xe trung chuyển không thể tới được.
								</p>
							</div>
						</div>
						<div className="mt-5 text-sm font-[400]">
							<input type="checkbox"></input>
							<span>
								<span className="cursor-pointer text-red-500 underline ml-3">Chấp nhận điều khoản</span>{' '}
								đặt vé &amp; chính sách bảo mật thông tin của FUTABusline
							</span>
						</div>
					</div>
					{/* <div className="address-info-section bg-white border border-slate-200 p-5">
						<div className="flex flex-row items-center gap-x-3 mb-5">
							<h3 className="text-xl font-semibold">Thông tin đón trả</h3>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								strokeWidth="1.5"
								stroke="currentColor"
								className="w-7 h-7 text-red-500"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
								/>
							</svg>
						</div>
						<div className="flex flex-col lg:flex-row gap-8">
							<div className="basis-1/2 w-full">
								<h4 className="text-medium font-[500] mb-5">ĐIỂM ĐÓN</h4>
								<div className="flex flex-row items-center gap-x-5">
									<div>
										<input
											disabled
											checked
											id="disabled-radio-1"
											type="radio"
											value=""
											name="disabled-radio1"
											className="w-4 h-4 text-blue-500 bg-gray-100 border-gray-300 focus:ring-blue-500"
										></input>
										<label
											htmlFor="disabled-radio-1"
											className="ms-2 text-sm font-medium text-blue-500 "
										>
											Điểm đón
										</label>
									</div>
									<div>
										{' '}
										<input
											disabled
											id="disabled-radio-2"
											type="radio"
											value=""
											name="disabled-radio2"
											className="w-4 h-4 text-blue-500 bg-gray-100 border-gray-300 focus:ring-blue-500"
										></input>
										<label
											htmlFor="disabled-radio-2"
											className="ms-2 text-sm font-medium text-gray-400"
										>
											Trung chuyển
										</label>
									</div>
									<div>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											strokeWidth="1.5"
											stroke="currentColor"
											className="w-7 h-7 text-red-500"
										>
											<path
												strokeLinecap="round"
												strokeLinejoin="round"
												d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
											/>
										</svg>
									</div>
								</div>
								<div className="my-5">
									<input
										type="text"
										value="Bến Xe Miền Tây"
										disabled
										className="border border-slate-300 rounded-xl w-full"
									/>
								</div>
								<div className="flex flex-wrap gap-1 text-justify font-[420] text-[15px] leading-6">
									<span>
										Quý khách vui lòng có mặt tại Bến xe/Văn Phòng <span className="font-bold">BX Miền Tây</span>
										<span className="font-bold text-red-500"> Trước 19:15 04/05/2024</span> để được trung chuyển hoặc kiểm tra thông tin
										trước khi lên xe.
									</span>
								</div>
							</div>
							<div className="basis-1/2 w-full">
								<h4 className="text-medium font-[500] mb-5">ĐIỂM TRẢ</h4>
								<div className="flex flex-row items-center gap-x-5">
									<div>
										<input
											disabled
											checked
											id="disabled-radio-3"
											type="radio"
											value=""
											name="disabled-radio-3"
											className="w-4 h-4 text-blue-500 bg-gray-100 border-gray-300 focus:ring-blue-500"
										></input>
										<label
											htmlFor="disabled-radio-3"
											className="ms-2 text-sm font-medium text-blue-500 "
										>
											Điểm đón
										</label>
									</div>
									<div>
										{' '}
										<input
											disabled
											id="disabled-radio-4"
											type="radio"
											value=""
											name="disabled-radio-4"
											className="w-4 h-4 text-blue-500 bg-gray-100 border-gray-300 focus:ring-blue-500"
										></input>
										<label
											htmlFor="disabled-radio-4"
											className="ms-2 text-sm font-medium text-gray-400"
										>
											Trung chuyển
										</label>
									</div>
									<div>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											strokeWidth="1.5"
											stroke="currentColor"
											className="w-7 h-7 text-red-500"
										>
											<path
												strokeLinecap="round"
												strokeLinejoin="round"
												d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
											/>
										</svg>
									</div>
								</div>
								<div className="my-5">
									<input
										type="text"
										value="Vũng Tàu"
										disabled
										className="border border-slate-300 rounded-xl w-full"
									/>
								</div>
							</div>
						</div>
						<div className="payment-section"></div>
					</div> */}
				</div>
				<div className="basis-1/3 w-full flex flex-col gap-y-5">
					<div className="bg-white border border-slate-200 p-5 rounded-xl">
						<h3 className="text-xl font-medium">Thông tin lượt đi</h3>
						<div className="mt-4 flex justify-between">
							<span className="text-slate-500">Tuyến xe</span>
							<span className="text-right font-medium">{route}</span>
						</div>
						<div className="mt-1 flex items-center justify-between">
							<span className="text-slate-500">Thời gian xuất bến</span>
							<span className="text-green-500 font-medium">{start_time}</span>
						</div>
						<div className="mt-1 flex items-center justify-between">
							<span className="text-slate-500">Số lượng ghế</span>
							<span className="font-semibold">{numberOfSeats} Ghế</span>
						</div>
						<div className="mt-1 flex items-center justify-between">
							<span className="text-slate-500">Số ghế</span>
							<span className="text-green-500 font-medium">{seatID()}</span>
						</div>
						<div className="mt-1 flex items-center justify-between">
							<span className="text-slate-500">Tổng tiền lượt đi</span>
							<span className="text-green-500 font-medium">
								{(price * numberOfSeats) / 1000 + '.000'} đ
							</span>
						</div>
					</div>
					<div className="bg-white border border-slate-200 p-5 rounded-xl">
						<div className="flex flex-row items-center gap-x-3">
							{' '}
							<h3 className="text-xl font-medium">Chi tiết giá</h3>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								strokeWidth="1.5"
								stroke="currentColor"
								className="w-7 h-7 text-red-500"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
								/>
							</svg>
						</div>
						<div className="mt-4 flex items-center justify-between">
							<span className="text-slate-500">Giá vé lượt đi</span>
							<span className="text-red-500 font-medium">
								{(price * numberOfSeats) / 1000 + '.000'} đ
							</span>
						</div>
						<div className="mt-1 flex items-center justify-between">
							<span className="text-slate-500">Phí thanh toán</span>
							<span className="font-medium">0đ</span>
						</div>
						<hr className="my-3" />
						<div className="flex items-center justify-between">
							<span className="text-slate-500">Tổng tiền</span>
							<span className="text-red-500 font-medium">
								{(price * numberOfSeats) / 1000 + '.000'} đ
							</span>
						</div>
					</div>
				</div>
			</div>
			<div className="mt-5 flex max-w-screen-lg flex-col md:flex-row mx-auto gap-8 ">
				<div className="md:basis-2/3 w-full gap-5 payment-section bg-white border border-slate-200 p-5 rounded-b-xl flex items-center">
					<div className="flex flex-col">
						{/* <span className="w-16 rounded-xl bg-blue-600 py-1 text-center text-xs text-white">VNPAY</span> */}
						<span className="mt-2 text-2xl font-medium text-black">
							{(price * numberOfSeats) / 1000 + '.000'} đ
						</span>
					</div>
					<div className="flex flex-auto items-center justify-end">
						{/* <button
								type="button"
								className="px-10 py-2 border border-slate-200 text-blue-500 rounded-full mr-6 hover:bg-blue-500 hover:text-white transition-all"
							>
								<span className="text-sm font-medium">Hủy</span>
							</button> */}
						<button
							type="button"
							className="px-5 py-2 text-white rounded-full mr-6 bg-red-500 hover:bg-red-500 transition-all"
							onClick={payment}
						>
							<span className="text-sm font-medium">Thanh toán</span>
						</button>
					</div>
				</div>
				<div className="basis-1/3 gap-5"></div>
			</div>
			{failureModal && (
				<FailureNotification
					func={{ closeModal: closeFailureModal }}
					message={message}
				/>
			)}
		</div>
	);
}

export default BookingTicketForm;
