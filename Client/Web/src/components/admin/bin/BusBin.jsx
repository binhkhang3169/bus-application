/** @format */

import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { API_URL } from '../../../configs/env';
import WarningNotification from '../../Noti/WarningNotification';
import SuccessNotification from '../../Noti/SuccessNotification';
import FailureNotification from '../../Noti/FailureNotification';
import { useNavigate } from 'react-router-dom';

const BusBin = () => {
	const navigate = useNavigate();
	// Data
	const [busAll, setBusAll] = useState([]);

	// Modal
	const [deleteModal, setDeleteModal] = useState(false);
	const [successModal, setSuccessModal] = useState(false);
	const [failureModal, setFailureModal] = useState(false);
	const [message, setMessage] = useState('');
	const [busId, setBusId] = useState('');

	useEffect(() => {
		getBusAll();
	}, []);

	// Send GET request to retrieve buses information
	const getBusAll = async () => {
		await axios
			.get(API_URL + 'employee/bus')
			.then((res) => {
				setBusAll(res.data.data.filter((v) => v.status === 0));
			})
			.catch((err) => {
				if (err.response.status === 401) {
					navigate('/admin');
				}
			});
	};

	// Open bus edit modal
	const restoreBtn = async (id) => {
		await axios
			.put(
				API_URL + `employee/bus/${id}`,
				{ status: 1 },
				{
					headers: { Authorization: 'Bearer' + sessionStorage.getItem('token') },
				}
			)
			.then((res) => {
				setMessage(res.data.message);
				openSuccessModal();
				refresh();
			})
			.catch((err) => {
				setMessage(err.response.data.message);
				openFailureModal();
			});
	};

	// Open delete modal
	const deleteBtn = async (id) => {
		setBusId(id);
		setDeleteModal(true);
	};

	// Refresh page
	const refresh = () => {
		getBusAll();
		setBusId('');
	};

	// Close delete modal
	const closeDeleteModal = () => {
		setDeleteModal(false);
		setBusId('');
	};

	// Close Success Modal
	const closeSuccessModal = () => {
		setSuccessModal(false);
	};

	// Close Failure Modal
	const closeFailureModal = () => {
		setFailureModal(false);
	};

	// Open Success Modal
	const openSuccessModal = () => {
		setSuccessModal(true);
	};

	// Open Failure Modal
	const openFailureModal = () => {
		setFailureModal(true);
	};
	return (
		<div className="w-full p-2">
			<div className="flex justify-between items-center mb-8">
				<h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white">Bus Bin</h1>
				<button
					className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center"
					onClick={() => navigate('/admin/bus')}
				>
					Bus Management
				</button>
			</div>
			
			<div className="flex flex-col xl:flex-row gap-4">
				{/* Table 1 */}
				<div className="flex-1 relative overflow-x-auto shadow-md sm:rounded-lg bg-white dark:bg-gray-900">
					<table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
						<thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-800 dark:text-gray-300">
							<tr>
								<th scope="col" className="px-6 py-3">License</th>
								<th scope="col" className="px-6 py-3">Action</th>
							</tr>
						</thead>
						<tbody>
							{busAll.slice(0, 10).map((v, i) => (
								<tr key={i} className="border-b dark:border-gray-700 odd:bg-white even:bg-gray-50 dark:odd:bg-gray-900 dark:even:bg-gray-800/50">
									<td className="px-6 py-4 text-gray-800 dark:text-gray-200">{v.license}</td>
									<td className="px-6 py-4">
										<button onClick={() => restoreBtn(v.id)} className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline">
											Restore
										</button>
										<button onClick={() => deleteBtn(v.id)} className="font-medium text-red-600 dark:text-red-500 hover:underline">
											Delete
										</button>
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>

				{/* Table 2 for more than 10 items */}
				{busAll.length > 10 && (
					<div className="flex-1 relative overflow-x-auto shadow-md sm:rounded-lg bg-white dark:bg-gray-900">
						<table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
							<thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-800 dark:text-gray-300">
								<tr>
									<th scope="col" className="px-6 py-3">License</th>
									<th scope="col" className="px-6 py-3">Action</th>
								</tr>
							</thead>
							<tbody>
								{busAll.slice(10).map((v, i) => (
									<tr key={i + 10} className="border-b dark:border-gray-700 odd:bg-white even:bg-gray-50 dark:odd:bg-gray-900 dark:even:bg-gray-800/50">
										<td className="px-6 py-4 text-gray-800 dark:text-gray-200">{v.license}</td>
										<td className="px-6 py-4">
											<button onClick={() => restoreBtn(v.id)} className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline">
												Restore
											</button>
											<button onClick={() => deleteBtn(v.id)} className="font-medium text-red-600 dark:text-red-500 hover:underline">
												Delete
											</button>
										</td>
									</tr>
								))}
							</tbody>
						</table>
					</div>
				)}
			</div>

			{deleteModal && (
				<WarningNotification
					id={busId}
					func={{ refresh, closeModal: closeDeleteModal, openSuccessModal, openFailureModal, setMessage }}
					type={'bus'}
					action={'bus'}
					method={'delete'}
				/>
			)}
			{successModal && <SuccessNotification func={{ closeModal: closeSuccessModal }} message={message} />}
			{failureModal && <FailureNotification func={{ closeModal: closeFailureModal }} message={message} />}
		</div>
	);
};

export default BusBin;