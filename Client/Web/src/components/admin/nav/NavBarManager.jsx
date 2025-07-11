/** @format */

import React, { useState } from "react";
import {
  Avatar,
  BusStation,
  Buses,
  Customer,
  Employee,
  Logout,
  Menu,
  Route,
  Pickup,
  Ticket,
  Trip,
} from "../../../svg/svg";
import axios from "axios";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";

const NavBarManager = () => {
  const navigate = useNavigate();

  const [nav, setNav] = useState(false);

  const logout = async () => {
    localStorage.removeItem("adminAccessToken");
    sessionStorage.removeItem("adminAccessToken");
    navigate("/admin");
    window.location.reload();
  };

  return (
    <>
      <div className="hidden lg:block">
        <div className="w-60"></div>
      </div>
      <div className="fixed hidden lg:block">
        <div className="bg-blue-500 w-60 h-screen flex flex-col">
          <h1 className="w-full text-white font-bold text-2xl mt-2 text-center">
            Manager
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
                onClick={() => navigate("province")}
              >
                <BusStation
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">
                  Province Management
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
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("ticket")}
              >
                <Ticket
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">
                  Ticket Management
                </div>
              </div>
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("employee")}
              >
                <Employee
                  className="ml-4"
                  style={{ width: "16", height: "16" }}
                />
                <div className="text-white  text-sm  m-2">
                  Employee Management
                </div>
              </div>
              <div
                className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                onClick={() => navigate("customer")}
              >
                <Customer
                  className="ml-4"
                  style={{ width: "14px", height: "14px" }}
                />
                <div className="text-white  text-sm  m-2">
                  Customer Management
                </div>
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
      <div
        className="flex justify-center items-center lg:hidden fixed top-2 left-2 w-12 h-12 rounded-xl 
							bg-blue-500 hover:bg-blue-600 cursor-pointer "
        onClick={() => setNav(!nav)}
      >
        <Menu className="h-8 w-8" />
      </div>
      {nav && (
        <div className="fixed z-50">
          <div className="bg-blue-500 w-screen h-screen flex flex-col">
            <div className="flex justify-center items-center">
              <Menu
                className="ml-4 mt-4 h-8 w-8"
                onClick={() => setNav(!nav)}
              />
              <h1 className="w-full text-white font-bold text-2xl mt-2 text-center">
                Manager
              </h1>
            </div>
            <div className="flex flex-col flex-grow justify-between my-4">
              <div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("trip");
                    setNav(!nav);
                  }}
                >
                  <Trip
                    className="ml-4"
                    style={{ width: "14px", height: "14px" }}
                  />
                  <div className="text-white  text-lg  m-2">
                    Trip Management
                  </div>
                </div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("route");
                    setNav(!nav);
                  }}
                >
                  <Route
                    className="ml-4"
                    style={{ width: "14px", height: "14px" }}
                  />
                  <div className="text-white  text-lg  m-2">
                    Route Management
                  </div>
                </div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("bus-station");
                    setNav(!nav);
                  }}
                >
                  <BusStation
                    className="ml-4"
                    style={{ width: "14px", height: "14px" }}
                  />
                  <div className="text-white  text-lg  m-2">
                    Bus Station Management
                  </div>
                </div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("bus");
                    setNav(!nav);
                  }}
                >
                  <Buses
                    className="ml-4"
                    style={{ width: "14px", height: "14px" }}
                  />
                  <div className="text-white  text-lg  m-2">Bus Management</div>
                </div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("ticket");
                    setNav(!nav);
                  }}
                >
                  <Ticket
                    className="ml-4"
                    style={{ width: "14px", height: "14px" }}
                  />
                  <div className="text-white  text-lg  m-2">
                    Ticket Management
                  </div>
                </div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("employee");
                    setNav(!nav);
                  }}
                >
                  <Employee
                    className="ml-4"
                    style={{ width: "16", height: "16" }}
                  />
                  <div className="text-white  text-lg  m-2">
                    Employee Management
                  </div>
                </div>
                <div
                  className="w-full flex items-center cursor-pointer py-2
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("customer");
                    setNav(!nav);
                  }}
                >
                  <Customer
                    className="ml-4"
                    style={{ width: "14px", height: "14px" }}
                  />
                  <div className="text-white  text-lg  m-2">
                    Customer Management
                  </div>
                </div>
              </div>
              <div>
                <div
                  className="w-full flex items-center cursor-pointer py-2 
                            hover:bg-blue-400 transition-all duration-100 ease-linear"
                  onClick={() => {
                    navigate("account");
                    setNav(!nav);
                  }}
                >
                  <Avatar
                    className="ml-4"
                    style={{ width: "16px", height: "16px" }}
                  />
                  <div className="text-white  text-lg  m-2">My Account</div>
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
                  <div className="text-white  text-lg  m-2">Logout</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default NavBarManager;
