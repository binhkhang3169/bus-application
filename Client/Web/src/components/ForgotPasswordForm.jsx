/** @format */

import axios from 'axios';
import React, { useState } from 'react';
import { API_URL } from '../configs/env';
import { useNavigate } from 'react-router-dom';
import HomePromotion from './HomePromotion';

function ForgotPasswordForm() {
	const navigate = useNavigate();

	// Notìication
	const [message, setMessage] = useState('');

	// Data
	const [account, setAccount] = useState('');

	// Send POST request to forgot password
	const handleForgotPassword = () => {
		axios
			.post(`${API_URL}customer/forgot-password`, {
				account: account,
			})
			.then((res) => {
				navigate(`/forgot-password/verify?id=${res.data.customer_id}`);
			})
			.catch((err) => {
				setMessage(err.response.data.message);
			});
	};

	return (
		<div className="max-w-screen-lg mx-auto mb-20">
			<div
				// className="bg-slate-800 border border-slate-400 rounded-md p-8 shadow-lg backdrop-filter backdrop-blur-sm bg-opacity-30 relative"
				className="mx-auto w-full border-red-200 border-4 rounded-xl p-8 shadow-xl flex justify-center items-center"
			>
				<div className="w-full sm:basis-5/12 mx-4">
					<div className="mb-4">
						<h1 className="text-4xl font-bold text-center mb-3">Quên mật khẩu</h1>
						<p className="text-center">Vui lòng nhập thông tin của tài khoản bạn đã đăng ký</p>
					</div>
					<div>
						<div className="relative mb-2">
							<input
								type="text"
								className="block w-full py-2 px-0 text-sm bg-transparent border-0 border-b-2 border-gray-300 appearance-none focus:outline-none focus:ring-0  focus:border-red-600 peer"
								placeholder=""
								onChange={(e) => {
									setAccount(e.target.value);
									setMessage('');
								}}
								value={account}
								onKeyDown={(event) => {
									if (event.key === 'Enter') {
										handleForgotPassword();
									}
								}}
							/>
							<label
								htmlFor=""
								className="absolute text-md duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-red-500 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
							>
								Nhập số điện thoại/email
							</label>
						</div>
						{message && <div className="text-red-600 mt-2 text-sm font-medium">{message}</div>}

						<button
							type="button"
							className="transition-colors duration-300 w-full mb-4 text-md mt-6 rounded-full bg-red-500 text-white hover:bg-red-600 hover:text-white py-2"
							onClick={handleForgotPassword}
						>
							Gửi mã xác thực
						</button>
					</div>
				</div>
				<div className="hidden sm:block flex-grow mx-4">
					<img
						alt=""
						loading="lazy"
						decoding="async"
						data-nimg="fill"
						className="transition-all duration-200 relative hidden object-contain sm:block h-full w-full transparent"
						src="https://storage.googleapis.com/futa-busline-cms-dev/image_f922bef1bb/image_f922bef1bb.svg"
					></img>
					{/* <div className="w-full h-full rounded-xl bg-orange-200"></div> */}
				</div>
			</div>
			<HomePromotion />
		</div>
	);
}

export default ForgotPasswordForm;
