/** @format */

import axios from 'axios';
import React, { useEffect } from 'react';
import { API_URL } from '../../configs/env';
import { useNavigate } from 'react-router-dom';

const ProtectRoute = ({ children, allowedRoles }) => {
	const navigate = useNavigate();

	useEffect(() => {
		authEmployee();
	});

	const authEmployee = async () => {
		const token = sessionStorage.getItem('token');
		if (token) {
			await axios
				.get(`${API_URL}employee/me`, {
					headers: {
						Authorization: `Bearer ${token}`,
					},
				})
				.then((res) => {
					checkRole(res.data.employee.role);
				})
				.catch((err) => {
					navigate('/admin');
				});
		}
	};

	const checkRole = (role) => {
		!allowedRoles.includes(role) && navigate('/admin/home');
	};

	return children;
};

export default ProtectRoute;
