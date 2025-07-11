/** @format */

import React from 'react';

function HomeNews() {
	return (
		<div className="my-10">
			<h3 className="text-center text-2xl font-bold text-green-700 mb-5">TIN TỨC MỚI</h3>
			<hr className="w-4/5 mx-auto h-0.5 bg-gray-200 md:w-2/5 xl:w-1/5" />
			<div className="max-w-screen-lg mx-auto mt-10 flex flex-col md:flex-row justify-center items-start gap-8">
				<div className="basis-1/3 shawdow-2xl">
					<img
						loading="lazy"
						src="https://storage.googleapis.com/futa-busline-web-cms-prod/599_x_337_px_min_f4d6bf9dec/599_x_337_px_min_f4d6bf9dec.png"
						alt="news-1"
						className="rounded-xl shadow-2xl cursor-pointer"
					/>
					<div className="py-3">
						<p className="font-semibold">PHƯƠNG TRANG - LIÊN TIẾP ĐẠT TOP 1 CHẤT LƯỢNG ASEAN 2024</p>
						<div className="flex justify-between items-center py-3">
							<span className="text-sm text-slate-500">8/12/2023</span>
							<span className="text-red-500 cursor-pointer text-sm font-semibold">Chi tiết &gt;</span>
						</div>
					</div>
				</div>
				<div className="basis-1/3 shawdow-2xl">
					<img
						loading="lazy"
						src="https://storage.googleapis.com/futa-busline-web-cms-prod/599_x_337_px_0882f6f487/599_x_337_px_0882f6f487.png"
						alt="news-2"
						className="rounded-xl shadow-2xl cursor-pointer"
					/>
					<div className="py-3">
						<p className="font-semibold">🚌🚌 CÙNG PHƯƠNG TRANG ĐI AN GIANG, NÚI SẬP 🚌🚌</p>
						<div className="flex justify-between items-center py-3">
							<span className="text-sm text-slate-500">8/12/2023</span>
							<span className="text-red-500 cursor-pointer text-sm font-semibold">Chi tiết &gt;</span>
						</div>
					</div>
				</div>
				<div className="basis-1/3 shawdow-2xl">
					<img
						loading="lazy"
						src="https://storage.googleapis.com/futa-busline-web-cms-prod/599_x_337_px_7aab7fa492/599_x_337_px_7aab7fa492.png"
						alt="news-3"
						className="rounded-xl shadow-2xl cursor-pointer"
					/>
					<div className="mt-3">
						<p className="font-semibold">CÔNG TY PHƯƠNG TRANG KHAI TRƯƠNG 2 TUYẾN XE BUÝT Ở VĨNH LONG</p>
						<div className="flex justify-between items-center mt-3">
							<span className="text-sm text-slate-500">8/12/2023</span>
							<span className="text-red-500 cursor-pointer text-sm font-semibold">Chi tiết &gt;</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}

export default HomeNews;
