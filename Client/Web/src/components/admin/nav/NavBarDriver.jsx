/** @format */

import React from 'react';
import { Avatar, Logout, Trip } from '../../../svg/svg';
import axios from 'axios';
import { API_URL } from '../../../configs/env';
import { useNavigate } from 'react-router-dom';

const NavBarDriver = () => {
	const navigate = useNavigate();

	const logout = async () => {
		localStorage.removeItem('adminAccessToken');
		sessionStorage.removeItem('adminAccessToken');
		navigate('/admin');
		window.location.reload();
	};

	return (
		<div>
			<div className="w-60"></div>
			<div className="fixed">
				<div className="bg-blue-500 w-60 h-screen flex flex-col">
					<h1 className="w-full text-white font-bold text-2xl mt-2 text-center">Driver</h1>
					<div className="flex flex-col flex-grow justify-between my-4">
						<div>
							<div
								className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
								onClick={() => navigate('my-trip')}
							>
								<Trip
									className="ml-4"
									style={{ width: '14px', height: '14px' }}
								/>
								<div className="text-white  text-sm  m-2">My Trip</div>
							</div>
						</div>
						<div>
							<div
								className="w-full flex items-center cursor-pointer py-2 
								hover:bg-blue-400 transition-all duration-100 ease-linear"
								onClick={() => navigate('account')}
							>
								<Avatar
									className="ml-4"
									style={{ width: '16px', height: '16px' }}
								/>
								<div className="text-white  text-sm  m-2">My Account</div>
							</div>
							<div
								className="w-full flex items-center cursor-pointer py-2 
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
								onClick={logout}
							>
								<Logout
									className="ml-4"
									style={{ width: '16px', height: '16px' }}
								/>
								<div className="text-white  text-sm  m-2">Logout</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default NavBarDriver;
