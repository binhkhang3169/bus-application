/** @format */

import React from 'react';
import { WarningIcon } from '../../svg/svg';
import axios from 'axios';
import { API_URL } from '../../configs/env';

const WarningNotification = ({ type, action, func, id, method }) => {
	const handleDelete = async () => {
		if (method === 'put') {
			await axios
				.put(
					API_URL + `api/v1/${action}/${id}`,
					{ status: 0 },
					{
						headers: { Authorization: 'Bearer ' + sessionStorage.getItem('adminAccessToken') },
					}
				)
				.then((res) => {
					console.log(res.data);
					func.closeModal();
					func.setMessage(res.data.message);
					func.openSuccessModal();
					func.refresh();
				})
				.catch((err) => {
					console.log(err.response);
					func.setMessage(err.response.data.message);
					func.openFailureModal();
					func.closeModal();
				});
		} else {
			await axios
				.delete(API_URL + `api/v1/${action}/${id}`, {
					headers: { Authorization: 'Bearer ' + sessionStorage.getItem('adminAccessToken') },
				})
				.then((res) => {
					func.closeModal();
					func.setMessage(res.data.message);
					func.openSuccessModal();
					func.refresh();
				})
				.catch((err) => {
					func.setMessage(err.response.data.message);
					func.openFailureModal();
					func.closeModal();
				});
		}
	};
	return (
		<div className="fixed z-50 top-0 left-0 w-full h-screen justify-center bg-black/20 ">
			<div className="overflow-y-auto overflow-x-hidden top-0 right-0 left-0 z-50 flex justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
				<div className="relative p-4 w-full max-w-md max-h-full">
					<div className="relative bg-white rounded-lg shadow ">
						<button
							type="button"
							className="absolute top-3 end-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center"
							onClick={func.closeModal}
						>
							<svg
								className="w-3 h-3"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 14 14"
							>
								<path
									stroke="currentColor"
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth="2"
									d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
								/>
							</svg>
							<span className="sr-only">Close modal</span>
						</button>
						<div className="p-4 md:p-5 text-center">
							<div className="flex justify-center">
								<WarningIcon className="w-[80px] aspect-square" />
							</div>

							<h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">{`Are you sure you want to delete this ${type}?`}</h3>
							<button
								type="button"
								className="text-white bg-red-600 hover:bg-red-800 font-medium rounded-lg text-sm inline-flex items-center px-5 py-2.5 text-center"
								onClick={handleDelete}
							>
								Yes, I'm sure
							</button>
							<button
								type="button"
								className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10  focus:ring-gray-100 "
								onClick={func.closeModal}
							>
								No, cancel
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default WarningNotification;
