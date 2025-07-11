/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import { useNavigate } from "react-router-dom";
import WarningNotification from "../Noti/WarningNotification";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import EmployeeForm from "./modal/EmployeeForm";

const EmployeeManagement = () => {
  const navigate = useNavigate();
  const [employees, setEmployees] = useState([]);
  const [currentViewRole, setCurrentViewRole] = useState("ROLE_RECEPTION");
  const [deleteModal, setDeleteModal] = useState(false);
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [employeeModal, setEmployeeModal] = useState(false);
  const [message, setMessage] = useState("");
  const [employeeId, setEmployeeId] = useState("");

  useEffect(() => {
    getEmployeesByRole();
  }, [currentViewRole]);

  const getEmployeesByRole = async () => {
    const path = `api/v1/users/by-role?roleName=${currentViewRole}`;
    try {
      const res = await axios.get(API_URL + path, {
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("adminAccessToken"),
        },
      });
      setEmployees(res.data.data || []);
    } catch (err) {
      console.error("Error fetching employees:", err);
      if (err.response?.status === 401) navigate("/admin");
      setEmployees([]);
    }
  };

  const renderRoleForDisplay = (rolesArray) => {
    if (!rolesArray || rolesArray.length === 0) return "N/A";
    const role = rolesArray[0];
    switch (role) {
      case "ROLE_DRIVER":
        return "Driver";
      case "ROLE_RECEPTION":
        return "Receptionist";
      default:
        return role.replace("ROLE_", "").replace("_", " ");
    }
  };

  const editBtn = (id) => {
    setEmployeeId(id);
    openEmployeeModal();
  };
  const deleteBtn = (id) => {
    setEmployeeId(id);
    setDeleteModal(true);
  };
  const refresh = () => {
    getEmployeesByRole();
    setEmployeeId("");
  };
  const closeDeleteModal = () => {
    setDeleteModal(false);
    setEmployeeId("");
  };
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);
  const openEmployeeModal = () => setEmployeeModal(true);
  const closeEmployeeModal = () => {
    setEmployeeModal(false);
    setEmployeeId("");
  };
  const handleAddNewEmployee = () => {
    setEmployeeId("");
    openEmployeeModal();
  };

  return (
    <div className="w-full p-2">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="ml-16 lg:ml-0 font-bold text-2xl text-gray-800 dark:text-white"></h1>
          <div className="flex items-center gap-2">
            <button
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                currentViewRole === "ROLE_RECEPTION"
                  ? "bg-blue-700 text-white dark:bg-blue-600"
                  : "bg-gray-200 text-gray-700 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
              }`}
              onClick={() => setCurrentViewRole("ROLE_RECEPTION")}
            >
              Lễ tân
            </button>
            <button
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                currentViewRole === "ROLE_DRIVER"
                  ? "bg-blue-700 text-white dark:bg-blue-600"
                  : "bg-gray-200 text-gray-700 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
              }`}
              onClick={() => setCurrentViewRole("ROLE_DRIVER")}
            >
              Tài xế
            </button>
          </div>
          <button
            className="text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
            onClick={handleAddNewEmployee}
          >
            {currentViewRole === "ROLE_RECEPTION"
              ? "Thêm lễ tân"
              : "Thêm tài xế"}
          </button>
        </div>
        <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
          <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-300">
              <tr>
                <th scope="col" className="px-2 py-3">
                  Mã nhân viên
                </th>
                <th scope="col" className="px-2 py-3">
                  Họ tên
                </th>
                <th scope="col" className="px-2 py-3">
                  Email
                </th>
                <th scope="col" className="px-2 py-3">
                  Số điện thoại
                </th>
                <th scope="col" className="px-2 py-3">
                  Địa chỉ
                </th>
                <th scope="col" className="px-2 py-3">
                  Giới tính
                </th>
                <th scope="col" className="px-2 py-3">
                  Ngày sinh
                </th>

                <th scope="col" className="px-2 py-3">
                  Tùy chọn
                </th>
              </tr>
            </thead>
            <tbody>
              {employees && employees.length > 0 ? (
                employees.map((employee, i) => (
                  <tr
                    key={employee.id || i}
                    className="odd:bg-white even:bg-gray-50 border-b dark:odd:bg-gray-900 dark:even:bg-gray-800 dark:border-gray-700"
                  >
                    <td className="px-2 py-4">{employee.id}</td>

                    <td className="px-2 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                      {employee.details?.fullName || "N/A"}
                    </td>
                    <td className="px-2 py-4">{employee.username || "N/A"}</td>
                    <td className="px-2 py-4">
                      {employee.details?.phoneNumber || "N/A"}
                    </td>
                    <td className="px-2 py-4">
                      <p className="line-clamp-1">
                        {employee.details?.address || "N/A"}
                      </p>
                    </td>
                    <td className="px-2 py-4">
                      {employee.details?.gender === "MALE"
                        ? "Nam"
                        : employee.details?.gender === "FEMALE"
                        ? "Nữ"
                        : "N/A"}
                    </td>
                    <td className="px-2 py-4">
                      {employee.details?.dateOfBirth
                        ? new Date(
                            employee.details.dateOfBirth
                          ).toLocaleDateString()
                        : "N/A"}
                    </td>

                    <td className="px-2 py-4">
                      <button
                        onClick={() => editBtn(employee.id)}
                        className="mr-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                      >
                        Sửa
                      </button>
                      <button
                        onClick={() => deleteBtn(employee.id)}
                        className="font-medium text-red-600 dark:text-red-500 hover:underline"
                      >
                        Cấm
                      </button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan="8"
                    className="text-center py-4 text-gray-500 dark:text-gray-400"
                  >
                    No{" "}
                    {currentViewRole === "ROLE_RECEPTION"
                      ? "receptionists"
                      : "drivers"}{" "}
                    found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
      {deleteModal && (
        <WarningNotification
          id={employeeId}
          func={{
            refresh,
            closeModal: closeDeleteModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
          }}
          type={"employee"}
          action={"employee"}
        />
      )}
      {successModal && (
        <SuccessNotification
          func={{ closeModal: closeSuccessModal }}
          message={message}
        />
      )}
      {failureModal && (
        <FailureNotification
          func={{ closeModal: closeFailureModal }}
          message={message}
        />
      )}
      {employeeModal && (
        <EmployeeForm
          func={{
            closeModal: closeEmployeeModal,
            openSuccessModal,
            openFailureModal,
            setMessage,
            refresh,
          }}
          employeeId={employeeId}
          roleForNewEmployee={employeeId === "" ? currentViewRole : null}
        />
      )}
    </div>
  );
};

export default EmployeeManagement;
