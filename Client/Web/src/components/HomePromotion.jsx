/** @format */

import React from 'react';

function HomePromotion() {
	return (
		<div className="my-10">
			<h3 className="text-center text-2xl font-bold text-green-700 mb-5">KHUYẾN MÃI NỔI BẬT</h3>
			<hr className="w-4/5 mx-auto h-0.5 bg-gray-200 md:w-2/5 xl:w-1/5" />
			<div className="max-w-screen-lg mx-auto mt-10 flex flex-col md:flex-row justify-center items-center gap-8">
				<div className="basis-1/3 shawdow-2xl">
					<img
						loading="lazy"
						src="https://storage.googleapis.com/futa-busline-cms-dev/343x184_4x_29d182ce55/343x184_4x_29d182ce55.png"
						alt="promotion-img1"
						className="rounded-xl shadow-2xl cursor-pointer"
					/>
				</div>
				<div className="basis-1/3 shawdow-2xl">
					<img
						loading="lazy"
						src="https://storage.googleapis.com/futa-busline-web-cms-prod/Zalo_11b66ecb81/Zalo_11b66ecb81.png"
						alt="promotion-img2"
						className="rounded-xl shadow-2xl cursor-pointer"
					/>
				</div>
				<div className="basis-1/3 shawdow-2xl">
					<img
						loading="lazy"
						src="https://storage.googleapis.com/futa-busline-web-cms-prod/343_x_184_px_0b1588190d/343_x_184_px_0b1588190d.png"
						alt="promotion-img3"
						className="rounded-xl shadow-2xl cursor-pointer"
					/>
				</div>
			</div>
		</div>
	);
}

export default HomePromotion;
