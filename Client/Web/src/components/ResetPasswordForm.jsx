/** @format */

import React from 'react';

function ResetPasswordForm() {
	return (
		<div className="basis-10/12 sm:basis-8/12 md:basis-6/12 lg:basis-5/12 xl:basis-4/12 2xl:basis-3/12 mb-40">
			<div className="bg-slate-800 border border-slate-400 rounded-md p-8 shadow-lg backdrop-filter backdrop-blur-sm bg-opacity-30 relative">
				<h1 className="text-4xl text-white font-bold text-center mb-10">Đặt lại mật khẩu</h1>
				<form action="">
					<div className="relative mb-8">
						<input
							type="text"
							className="block w-full py-2 px-0 text-sm text-white bg-transparent border-0 border-b-2 border-gray-300 appearance-none peer"
							placeholder=""
						/>
						<label
							htmlFor=""
							className="absolute text-md text-white duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-blue-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
						>
							Nhập mật khẩu mới
						</label>
						<span className="absolute top-1 right-1 ">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								strokeWidth="1.5"
								stroke="currentColor"
								className="w-6 h-6"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
								/>
							</svg>
						</span>
					</div>
					<div className="relative mb-8">
						<input
							type="password"
							className="block w-full py-2 px-0 text-sm text-white bg-transparent border-0 border-b-2 border-gray-300 appearance-none peer"
							placeholder=""
						/>
						<label
							htmlFor=""
							className="absolute text-md text-white duration-300 tranform -translate-y-7 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-blue-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-7"
						>
							Xác thực mật khẩu
						</label>
						<span className="absolute top-1 right-1 ">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								strokeWidth="1.5"
								stroke="currentColor"
								className="w-6 h-6"
							>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
								/>
							</svg>
						</span>
					</div>

					<button
						type="submit"
						className="transition-colors duration-300 w-full mb-4 text-[18px] mt-6 rounded-full bg-white text-black hover:bg-cyan-500 hover:text-white py-2"
					>
						Đặt lại mật khẩu
					</button>
					<div></div>
				</form>
			</div>
		</div>
	);
}

export default ResetPasswordForm;
