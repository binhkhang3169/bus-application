/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import { useNavigate } from "react-router-dom";

const parseJwt = (token) => {
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    return JSON.parse(jsonPayload);
  } catch (e) {
    console.error("Error parsing JWT token:", e);
    return null;
  }
};

const TripDriver = () => {
  const navigate = useNavigate();
  const [tripAll, setTripAll] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [loggedInDriverId, setLoggedInDriverId] = useState(null);

  useEffect(() => {
    const token =
      localStorage.getItem("adminAccessToken") ||
      sessionStorage.getItem("adminAccessToken");
    if (token) {
      const decodedToken = parseJwt(token);
      if (decodedToken && decodedToken.id) {
        setLoggedInDriverId(decodedToken.id);
      } else {
        setError("Unable to verify user. Please log in again.");
        setIsLoading(false);
      }
    } else {
      navigate("/admin");
    }
  }, [navigate]);

  useEffect(() => {
    if (loggedInDriverId) {
      getTripAll();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [loggedInDriverId]);

  const getTripAll = async () => {
    setIsLoading(true);
    setError(null);
    const token =
      localStorage.getItem("adminAccessToken") ||
      sessionStorage.getItem("adminAccessToken");
    if (!token) {
      navigate("/admin");
      setIsLoading(false);
      return;
    }

    try {
      const response = await axios.get(`${API_URL}api/v1/trips/driver`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (
        response.data &&
        response.data.code === 200 &&
        Array.isArray(response.data.data)
      ) {
        let trips = response.data.data.filter((v) => {
          const tripStartDateTime = new Date(
            `${v.departureDate}T${v.departureTime}`
          );
          const isFutureTrip = new Date() < tripStartDateTime;
          const isCorrectDriver = loggedInDriverId
            ? v.driverId === loggedInDriverId
            : false;
          return isFutureTrip && isCorrectDriver;
        });
        setTripAll(trips);
      } else {
        setError("Failed to retrieve trips. Invalid data format.");
        setTripAll([]);
      }
    } catch (err) {
      if (
        err.response &&
        (err.response.status === 401 || err.response.status === 403)
      ) {
        navigate("/admin");
      } else {
        setError(err.message || "An error occurred while fetching trips.");
      }
      setTripAll([]);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="w-full p-4 text-center dark:text-gray-300">
        Loading trips...
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full p-4 text-center text-red-500 dark:text-red-400">
        Error: {error}
      </div>
    );
  }

  return (
    <div className="w-full p-2">
      <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white mb-6">
        My Trip
      </h1>
      <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
        <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
          <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
            <tr>
              <th scope="col" className="px-6 py-3">
                Route
              </th>
              <th scope="col" className="px-6 py-3">
                Seats Available
              </th>
              <th scope="col" className="px-6 py-3">
                Time (Depart - Arrive)
              </th>
              <th scope="col" className="px-6 py-3">
                Date (Departure)
              </th>
              <th scope="col" className="px-6 py-3">
                Bus License
              </th>
            </tr>
          </thead>
          <tbody>
            {tripAll.length > 0 ? (
              tripAll.map((v) => (
                <tr
                  key={v.id}
                  className="odd:bg-white even:bg-gray-50 border-b hover:bg-gray-100 dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700/50"
                >
                  <td
                    scope="row"
                    className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                  >
                    {v.route?.start?.name && v.route?.end?.name
                      ? `${v.route.start.name} to ${v.route.end.name}`
                      : "N/A"}
                  </td>
                  <td className="px-6 py-4">
                    {typeof v.stock === "number" ? v.stock : "N/A"}
                  </td>
                  <td className="px-6 py-4">
                    {v.departureTime && v.arrivalTime
                      ? `${v.departureTime.substring(
                          0,
                          5
                        )} - ${v.arrivalTime.substring(0, 5)}`
                      : "N/A"}
                  </td>
                  <td className="px-6 py-4">{v.departureDate || "N/A"}</td>
                  <td className="px-6 py-4">{v.vehicle?.license || "N/A"}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td
                  colSpan="5"
                  className="px-6 py-4 text-center text-gray-500 dark:text-gray-400"
                >
                  No upcoming trips assigned to you.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default TripDriver;
