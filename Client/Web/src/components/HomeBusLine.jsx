/** @format */

import React from 'react';

function HomeBusLine() {
	return (
		<div className="py-20">
			<h3 className="text-center text-2xl font-bold text-green-700 mb-3">TUYẾN PHỔ BIẾN</h3>
			<p className="text-center text-slate-500 mb-5">Được khách hàng tin tưởng và lựa chọn</p>
			<hr className="w-4/5 mx-auto h-0.5 bg-gray-200 md:w-2/5 xl:w-1/5" />
			<div className="max-w-screen-lg mx-auto mt-10 flex flex-col md:flex-row justify-center items-center gap-8">
				<div className="basis-1/3 bus-line-card border border-slate-300 rounded-2xl flex flex-col shadow-2xl cursor-pointer">
					<div className="img relative">
						<img
							loading="lazy"
							src="https://storage.googleapis.com/futa-busline-cms-dev/Rectangle_23_2_8bf6ed1d78/Rectangle_23_2_8bf6ed1d78.png"
							alt="busline-1"
							className="rounded-xl"
						/>
						<div className="txt-inside-img absolute bottom-2 left-4 text-white ">
							Tuyến xe từ <p className="text-xl font-bold">Tp Hồ Chí Minh</p>
						</div>
					</div>
					<div className="bus-line-info p-3 border-b border-slate-300">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Đà Lạt</h5>
							<p className="font-semibold">290.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">305km - 8 giờ - 1/5/2024</p>
					</div>
					<div className="bus-line-info p-3 border-b border-slate-300">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Cần Thơ</h5>
							<p className="font-semibold">165.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">366km - 3 giờ 12 phút - 1/5/2024</p>
					</div>
					<div className="bus-line-info p-3">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Long Xuyên</h5>
							<p className="font-semibold">190.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">186km - 5 giờ - 1/5/2024</p>
					</div>
				</div>
				<div className="basis-1/3 bus-line-card border border-slate-300 rounded-2xl flex flex-col shadow-2xl cursor-pointer">
					<div className="img relative">
						<img
							loading="lazy"
							src="https://storage.googleapis.com/futa-busline-cms-dev/Rectangle_23_3_2d8ce855bc/Rectangle_23_3_2d8ce855bc.png"
							alt="busline-2"
							className="rounded-xl"
						/>
						<div className="txt-inside-img absolute bottom-2 left-4 text-white ">
							Tuyến xe từ <p className="text-xl font-bold">Đà Lạt</p>
						</div>
					</div>
					<div className="bus-line-info p-3 border-b border-slate-300">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">TP. Hồ Chí Minh</h5>
							<p className="font-semibold">290.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">310km - 8 giờ - 2/5/2024</p>
					</div>
					<div className="bus-line-info p-3 border-b border-slate-300">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Đà Nẵng</h5>
							<p className="font-semibold">410.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">757km - 17 giờ - 1/5/2024</p>
					</div>
					<div className="bus-line-info p-3">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Cần Thơ</h5>
							<p className="font-semibold">435.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">457km - 11 giờ - 1/5/2024</p>
					</div>
				</div>
				<div className="basis-1/3 bus-line-card border border-slate-300 rounded-2xl flex flex-col shadow-2xl cursor-pointer">
					<div className="img relative">
						<img
							loading="lazy"
							src="https://storage.googleapis.com/futa-busline-cms-dev/Rectangle_23_4_061f4249f6/Rectangle_23_4_061f4249f6.png"
							alt="busline-3"
							className="rounded-xl"
						/>
						<div className="txt-inside-img absolute bottom-2 left-4 text-white ">
							Tuyến xe từ <p className="text-xl font-bold">Đà Nẵng</p>
						</div>
					</div>
					<div className="bus-line-info p-3 border-b border-slate-300">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Đà Lạt</h5>
							<p className="font-semibold">410.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">666km - 17 giờ - 1/5/2024</p>
					</div>
					<div className="bus-line-info p-3 border-b border-slate-300">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">BX An Sương</h5>
							<p className="font-semibold">410.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">966km - 20 giờ - 1/5/2024</p>
					</div>
					<div className="bus-line-info p-3">
						<div className="flex justify-between">
							<h5 className="text-green-800 font-semibold text-lg">Nha Trang</h5>
							<p className="font-semibold">300.000đ</p>
						</div>
						<p className="text-slate-500 text-[16px]">528km - 9 giờ 25 phút - 1/5/2024</p>
					</div>
				</div>
			</div>
		</div>
	);
}

export default HomeBusLine;
