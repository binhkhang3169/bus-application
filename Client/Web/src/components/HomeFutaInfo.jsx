/** @format */

import React from 'react';

function HomeFutaInfo() {
	return (
		<section className="flex flex-col items-center px-4 py-6 text-center sm:p-10 bg-white">
			<span className="home-title text-green-700 w-80 sm:w-full text-2xl font-bold">
				AnhPhung BUS LINES - CHẤT LƯỢNG LÀ DANH DỰ
			</span>
			<span className="home-title-content text-slate-500 my-3">Được khách hàng tin tưởng và lựa chọn</span>
			<hr className="w-4/5 mx-auto h-0.5 bg-gray-200 md:w-2/5 xl:w-1/5" />
			<div className="layout mt-8 sm:grid sm:grid-cols-2 sm:gap-16">
				<div>
					<div className="mb-6 flex items-center">
						<div className="">
							<img
								alt=""
								loading="lazy"
								width="96"
								height="96"
								decoding="async"
								data-nimg="1"
								className="transition-all duration-200 undefined transparent"
								src="https://storage.googleapis.com/futa-busline-cms-dev/Group_662c4422ba/Group_662c4422ba.svg"
							></img>
						</div>
						<div className="ml-4 flex flex-col text-left">
							<span className="text-2xl font-semibold text-black lg:text-3xl">
								Hơn 20 Triệu<span className="ml-3 text-base">Lượt khách</span>
							</span>
							<span className="text-gray">
								Phương Trang phục vụ hơn 20 triệu lượt khách bình quân 1 năm trên toàn quốc
							</span>
						</div>
					</div>
					<div className="mb-6 flex items-center">
						<div className="">
							<img
								alt=""
								loading="lazy"
								width="96"
								height="96"
								decoding="async"
								data-nimg="1"
								className="transition-all duration-200 undefined transparent"
								src="https://storage.googleapis.com/futa-busline-cms-dev/Store_55c0da8bd7/Store_55c0da8bd7.svg"
							></img>
						</div>
						<div className="ml-4 flex flex-col text-left">
							<span className="text-2xl font-semibold text-black lg:text-3xl">
								Hơn 350
								<span className="ml-3 text-base">Phòng vé - Bưu cục</span>
							</span>
							<span className="text-gray">
								Phương Trang có hơn 350 phòng vé, trạm trung chuyển, bến xe,... trên toàn hệ thống
							</span>
						</div>
					</div>
					<div className="mb-6 flex items-center">
						<div className="">
							<img
								alt=""
								loading="lazy"
								width="96"
								height="96"
								decoding="async"
								data-nimg="1"
								className="transition-all duration-200 undefined transparent"
								src="https://storage.googleapis.com/futa-busline-cms-dev/Group_2_75b5ed1dfd/Group_2_75b5ed1dfd.svg"
							></img>
						</div>
						<div className="ml-4 flex flex-col text-left">
							<span className="text-2xl font-semibold text-black lg:text-3xl">
								Hơn 1,000<span className="ml-3 text-base">Chuyến xe</span>
							</span>
							<span className="text-gray">
								Phương Trang phục vụ hơn 1,000 chuyến xe đường dài và liên tỉnh mỗi ngày
							</span>
						</div>
					</div>
				</div>
				<div className="relative hidden object-contain sm:block">
					<img
						alt=""
						loading="lazy"
						decoding="async"
						data-nimg="fill"
						className="transition-all duration-200 relative hidden object-contain sm:block h-full w-full transparent"
						src="https://storage.googleapis.com/futa-busline-cms-dev/image_f922bef1bb/image_f922bef1bb.svg"
					></img>
				</div>
			</div>
		</section>
	);
}

export default HomeFutaInfo;
