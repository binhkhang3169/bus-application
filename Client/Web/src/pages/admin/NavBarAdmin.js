/** @format */

import React, { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import NavBarCustomerService from '../../components/admin/nav/NavBarCustomerService';
import NavBarOperator from '../../components/admin/nav/NavBarOperator';
import NavBarManager from '../../components/admin/nav/NavBarManager';
import NavBarDriver from '../../components/admin/nav/NavBarDriver';
import { jwtDecode } from 'jwt-decode'; // Using the imported library

// Define allowed admin roles (consistency with LoginPageAdmin)
const ALLOWED_ADMIN_ROLES = ['ROLE_ADMIN', 'ROLE_RECEPTION', 'ROLE_DRIVER', 'ROLE_OPERATOR'];

const NavBarAdmin = () => {
    const navigate = useNavigate();
    const [role, setRole] = useState('');

    const clearAdminSessionAndRedirect = useCallback(() => {
        sessionStorage.removeItem('adminAccessToken');
        sessionStorage.removeItem('adminUsername');
        sessionStorage.removeItem('adminRole');
        setRole(''); // Clear local role state
        // Dispatch an event so other admin components can react if needed
        window.dispatchEvent(new CustomEvent('authChangeAdmin', { detail: { loggedIn: false } }));
        navigate('/admin'); // Redirect to admin login
    }, [navigate]);

    const checkAdminAuthAndSetRole = useCallback(() => {
        const token = sessionStorage.getItem('adminAccessToken');
        const storedRole = sessionStorage.getItem('adminRole'); // Get role stored at login

        if (token && storedRole) {
            try {
                const decodedToken = jwtDecode(token); // Decode JWT

                const now = Math.floor(Date.now() / 1000); // Current time in seconds
                const isExpired = decodedToken.exp < now;

                // Check if the role from token and stored role are allowed
                const isTokenRoleAllowed = ALLOWED_ADMIN_ROLES.includes(decodedToken.role);
                const isStoredRoleAllowed = ALLOWED_ADMIN_ROLES.includes(storedRole);
                // Also ensure the role in token matches the one stored during login for consistency
                const rolesMatch = decodedToken.role === storedRole;

                if (isExpired || !isTokenRoleAllowed || !isStoredRoleAllowed || !rolesMatch) {
                    console.warn('Admin auth check failed: Expired, role mismatch, or not allowed.', {
                        isExpired, isTokenRoleAllowed, isStoredRoleAllowed, rolesMatch,
                        tokenRole: decodedToken.role, storedRole
                    });
                    clearAdminSessionAndRedirect();
                } else {
                    // console.log('Admin auth check success. Role:', storedRole);
                    setRole(storedRole); // Set role from the initially stored role after validation
                }
            } catch (err) {
                console.error('Invalid admin token or decoding error:', err);
                clearAdminSessionAndRedirect();
            }
        } else {
            // No token or role found in sessionStorage
            // console.log('No admin token or role found in session storage.');
            // If role is already empty, no need to redirect, otherwise clear and redirect
            if (role) setRole(''); // Clear local state if it was somehow set
            // No automatic redirect here if already on a public admin page or login page.
            // Protected routes should handle redirection.
            // However, if this component is part of a layout for protected admin pages,
            // then a redirect might be appropriate if role is not set.
            // For now, it will simply not render a navbar if no valid role.
        }
    }, [clearAdminSessionAndRedirect, role]);

    useEffect(() => {
        checkAdminAuthAndSetRole();

        const handleAuthChange = (event) => {
            console.log('NavBarAdmin: authChangeAdmin event detected.', event.detail);
            checkAdminAuthAndSetRole(); // Re-check auth status
        };

        window.addEventListener('authChangeAdmin', handleAuthChange);
        return () => {
            window.removeEventListener('authChangeAdmin', handleAuthChange);
        };
    }, [checkAdminAuthAndSetRole]);

    const renderNavbarByRole = () => {
        if (!role) return null; // Don't render if role is not set or invalid

        switch (role) {
            case 'ROLE_ADMIN': // This usually implies manager or superuser
                return <NavBarManager />;
            // case 'accountant': // Example for a role not in current ALLOWED_ADMIN_ROLES
            //     // return <NavBarAccountant />; // If you had this component
            //     break;
            case 'ROLE_RECEPTION':
                return <NavBarCustomerService />;
            case 'ROLE_DRIVER':
                return <NavBarDriver />;
            case 'ROLE_OPERATOR':
                return <NavBarOperator />;
            default:
                console.warn(`No specific navbar defined for role: ${role}`);
                // Optionally, redirect or show a generic "access denied" or empty navbar
                // For now, clearing session and redirecting if role is unknown but token existed.
                // This default case should ideally not be hit if ALLOWED_ADMIN_ROLES is comprehensive
                // and login ensures only these roles get a token.
                // If we reach here, it might imply an unexpected role in a valid token.
                // clearAdminSessionAndRedirect(); // Consider if this is the desired behavior
                return null;
        }
    };

    // Only render if a valid role is set.
    return <>{role ? renderNavbarByRole() : null}</>;
};

export default NavBarAdmin;