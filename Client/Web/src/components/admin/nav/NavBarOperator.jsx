/** @format */

import React, { useEffect } from "react";
import {
  Avatar,
  BusStation,
  Buses,
  Logout,
  Pickup,
  Route,
  Trip,
} from "../../../svg/svg";
import axios from "axios";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";

const NavBarOperator = () => {
  const navigate = useNavigate();

  useEffect(() => {
    authEmployee();
  });

  const logout = async () => {
    localStorage.removeItem("adminAccessToken");
    sessionStorage.removeItem("adminAccessToken");
    navigate("/admin");
    window.location.reload();
  };

  const authEmployee = async () => {
    // get customer
    const token = sessionStorage.getItem("token");
    if (token) {
      axios
        .get(API_URL + "employee/me", {
          headers: { Authorization: `Bearer ${token}` },
        })
        .then((res) => {
          if (res.data.employee.role !== "operator") {
            checkRole(res.data.employee.role);
          }
        })
        .catch((err) => {
          navigate("/admin");
        });
    } else {
      navigate("/admin");
    }
  };

  const checkRole = (role) => {
    switch (role) {
      case "manager":
        navigate("/admin/quan-ly");
        break;
      case "accountant":
        navigate("/admin/ke-toan");
        break;
      case "customer_service":
        navigate("/admin/cham-soc-khach-hang");
        break;
      case "driver":
        navigate("/admin/tai-xe");
        break;
      case "operator":
        navigate("/admin/van-hanh");
        break;
      default:
        break;
    }
  };

  return (
    <div>
      <div className="w-60"></div>
      <div className="fixed">
        <div className="bg-blue-500 w-60 h-screen flex flex-col">
          <h1 className="w-full text-white font-bold text-2xl mt-2 text-center">
            Operator
          </h1>
          <div className="flex flex-col flex-grow justify-between my-4">
            <div>
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("trip")}
              >
                <Trip
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">Trip Management</div>
              </div>
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("route")}
              >
                <Route
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">Route Management</div>
              </div>
              <div
                className="w-full flex items-center cursor-pointer py-2 hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("pickup")}
              >
                <Pickup
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white text-sm m-2">Pickup Management</div>
              </div>
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("bus-station")}
              >
                <BusStation
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">
                  Bus Station Management
                </div>
              </div>
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("bus")}
              >
                <Buses
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">Bus Management</div>
              </div>
            </div>
            <div>
              <div
                className="w-full flex items-center cursor-pointer py-2 
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("account")}
              >
                <Avatar
                  className="ml-4"
                  style={{ width: "16px", height: "16px" }}
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
                  style={{ width: "16px", height: "16px" }}
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

export default NavBarOperator;
