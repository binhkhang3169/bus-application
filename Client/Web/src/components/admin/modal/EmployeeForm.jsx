/** @format */

import axios from "axios";
import React, { useEffect, useState, useCallback } from "react";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";

const EmployeeForm = ({ func, employeeId, roleForNewEmployee }) => {
  const navigate = useNavigate();

  // Common Fields
  const [username, setUsername] = useState("");
  const [fullName, setFullName] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [address, setAddress] = useState("");
  const [gender, setGender] = useState("");
  const [dateOfBirth, setDateOfBirth] = useState("");
  const [identityNumber, setIdentityNumber] = useState("");
  const [employeeType, setEmployeeType] = useState("");

  // Driver Specific Fields
  const [licenseNumber, setLicenseNumber] = useState("");
  const [licenseClass, setLicenseClass] = useState("");
  const [licenseIssuedDate, setLicenseIssuedDate] = useState("");
  const [licenseExpiryDate, setLicenseExpiryDate] = useState("");
  const [vehicleType, setVehicleType] = useState("");

  // Error state for validation
  const [errors, setErrors] = useState({});

  // Helper to format date string YYYY-MM-DD from ISO or null
  const formatDateForInput = (isoDateString) => {
    if (!isoDateString) return "";
    try {
      return new Date(isoDateString).toISOString().split("T")[0];
    } catch {
      return "";
    }
  };

  // Helper to format date to ISO string or null
  const formatDateForAPI = (dateString) => {
    if (!dateString) return null;
    try {
      return new Date(dateString).toISOString();
    } catch {
      return null;
    }
  };

  const resetFormFields = useCallback(() => {
    setUsername("");
    setFullName("");
    setPhoneNumber("");
    setAddress("");
    setGender("");
    setDateOfBirth("");
    setIdentityNumber("");
    setLicenseNumber("");
    setLicenseClass("");
    setLicenseIssuedDate("");
    setLicenseExpiryDate("");
    setVehicleType("");
    setErrors({});
  }, []);

  // Set input fields from fetched data
  const setInput = useCallback(
    (data) => {
      resetFormFields();

      let determinedType = "";
      if (data.roles && data.roles.length > 0) {
        determinedType = data.roles[0];
      } else if (data.employeeType) {
        determinedType = data.employeeType;
      } else if (data.role) {
        if (data.role === "driver") determinedType = "ROLE_DRIVER";
        else if (
          ["manager", "operator", "accountant", "customer_service"].includes(
            data.role
          )
        )
          determinedType = "ROLE_RECEPTION";
      }
      setEmployeeType(determinedType);

      setUsername(data.username || data.email || "");
      setFullName(
        data.details?.fullName ||
          `${data.first_name || ""} ${data.last_name || ""}`.trim() ||
          ""
      );
      setPhoneNumber(data.details?.phoneNumber || data.phone_number || "");
      setAddress(data.details?.address || data.address || "");

      let apiGender = data.details?.gender || data.gender;
      if (typeof apiGender === "number") {
        setGender(apiGender === 1 ? "MALE" : apiGender === 0 ? "FEMALE" : "");
      } else {
        setGender(apiGender || "");
      }

      setDateOfBirth(
        formatDateForInput(data.details?.dateOfBirth || data.date_of_birth)
      );
      setIdentityNumber(
        data.details?.identityNumber || data.identityNumber || ""
      );

      if (determinedType === "ROLE_DRIVER") {
        setLicenseNumber(
          data.details?.licenseNumber || data.licenseNumber || ""
        );
        setLicenseClass(data.details?.licenseClass || data.licenseClass || "");
        setLicenseIssuedDate(
          formatDateForInput(
            data.details?.licenseIssuedDate || data.licenseIssuedDate
          )
        );
        setLicenseExpiryDate(
          formatDateForInput(
            data.details?.licenseExpiryDate || data.licenseExpiryDate
          )
        );
        setVehicleType(data.details?.vehicleType || data.vehicleType || "");
      }
    },
    [resetFormFields]
  );

  useEffect(() => {
    if (employeeId) {
      axios
        .get(API_URL + `api/v1/create/${employeeId}`, {
          headers: {
            Authorization:
              "Bearer " + sessionStorage.getItem("adminAccessToken"),
          },
        })
        .then((res) => {
          setInput(res.data.employee);
        })
        .catch((err) => {
          if (err.response?.status === 401) {
            navigate("/admin");
          } else {
            console.error("Error fetching employee for edit:", err);
            func.setMessage(
              err.response?.data?.message || "Không thể lấy dữ liệu nhân viên."
            );
            func.openFailureModal();
            func.closeModal();
          }
        });
    } else {
      resetFormFields();
      if (roleForNewEmployee) {
        setEmployeeType(roleForNewEmployee);
      } else {
        setEmployeeType("ROLE_RECEPTION");
      }
    }
  }, [
    employeeId,
    roleForNewEmployee,
    navigate,
    func,
    setInput,
    resetFormFields,
  ]);

  // Validation function
  const validateForm = () => {
    const newErrors = {};

    // Common fields validation
    if (!username.trim()) {
      newErrors.username = "Email là bắt buộc.";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(username)) {
      newErrors.username = "Email không hợp lệ.";
    }

    if (!fullName.trim()) {
      newErrors.fullName = "Họ tên là bắt buộc.";
    }

    if (!phoneNumber.trim()) {
      newErrors.phoneNumber = "Số điện thoại là bắt buộc.";
    } else if (!/^\d{10,11}$/.test(phoneNumber)) {
      newErrors.phoneNumber = "Số điện thoại phải có 10-11 số.";
    }

    if (!address.trim()) {
      newErrors.address = "Địa chỉ là bắt buộc.";
    }

    if (!gender) {
      newErrors.gender = "Vui lòng chọn giới tính.";
    }

    if (!dateOfBirth) {
      newErrors.dateOfBirth = "Ngày sinh là bắt buộc.";
    } else {
      const dob = new Date(dateOfBirth);
      const today = new Date();
      if (isNaN(dob.getTime())) {
        newErrors.dateOfBirth = "Ngày sinh không hợp lệ.";
      } else if (dob > today) {
        newErrors.dateOfBirth = "Ngày sinh không được ở tương lai.";
      }
    }

    if (!identityNumber.trim()) {
      newErrors.identityNumber = "Số CCCD là bắt buộc.";
    } else if (!/^\d{12}$/.test(identityNumber)) {
      newErrors.identityNumber = "Số CCCD phải có 12 số.";
    }

    // Driver specific fields validation
    if (employeeType === "ROLE_DRIVER") {
      if (!licenseNumber.trim()) {
        newErrors.licenseNumber = "Số bằng lái là bắt buộc.";
      }

      if (!licenseClass.trim()) {
        newErrors.licenseClass = "Cấp độ bằng lái là bắt buộc.";
      }

      if (!licenseIssuedDate) {
        newErrors.licenseIssuedDate = "Ngày cấp bằng là bắt buộc.";
      } else {
        const issuedDate = new Date(licenseIssuedDate);
        if (isNaN(issuedDate.getTime())) {
          newErrors.licenseIssuedDate = "Ngày cấp bằng không hợp lệ.";
        }
      }

      if (!licenseExpiryDate) {
        newErrors.licenseExpiryDate = "Ngày hết hạn bằng là bắt buộc.";
      } else {
        const expiryDate = new Date(licenseExpiryDate);
        const issuedDate = new Date(licenseIssuedDate);
        if (isNaN(expiryDate.getTime())) {
          newErrors.licenseExpiryDate = "Ngày hết hạn bằng không hợp lệ.";
        } else if (licenseIssuedDate && expiryDate <= issuedDate) {
          newErrors.licenseExpiryDate = "Ngày hết hạn phải sau ngày cấp bằng.";
        }
      }

      if (!vehicleType.trim()) {
        newErrors.vehicleType = "Loại xe là bắt buộc.";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) {
      func.setMessage("Vui lòng kiểm tra và điền đầy đủ các trường bắt buộc.");
      func.openFailureModal();
      return;
    }

    let payload = {
      username: username,
      fullName: fullName,
      phoneNumber: phoneNumber,
      address: address,
      gender: gender,
      dateOfBirth: formatDateForAPI(dateOfBirth),
      identityNumber: identityNumber,
      employeeType: employeeType,
    };

    if (employeeType === "ROLE_DRIVER") {
      payload.licenseNumber = licenseNumber;
      payload.licenseClass = licenseClass;
      payload.licenseIssuedDate = formatDateForAPI(licenseIssuedDate);
      payload.licenseExpiryDate = formatDateForAPI(licenseExpiryDate);
      payload.vehicleType = vehicleType;
    }

    try {
      let response;
      if (employeeId) {
        response = await axios.put(
          API_URL + `api/v1/create/${employeeId}`,
          payload,
          {
            headers: {
              Authorization:
                "Bearer " + sessionStorage.getItem("adminAccessToken"),
            },
          }
        );
      } else {
        response = await axios.post(API_URL + "api/v1/create", payload, {
          headers: {
            Authorization:
              "Bearer " + sessionStorage.getItem("adminAccessToken"),
          },
        });
      }
      func.closeModal();
      func.setMessage(response.data.message);
      func.openSuccessModal();
      func.refresh();
    } catch (err) {
      if (err.response?.status === 401) {
        navigate("/admin");
      } else {
        func.closeModal();
        func.setMessage(err.response?.data?.message || "Đã xảy ra lỗi.");
        func.openFailureModal();
      }
    }
  };

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/20 w-full h-full">
      <div className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
        <div className="relative p-4 w-full max-w-3xl max-h-full mx-auto mt-10">
          <div className="relative bg-white rounded-lg shadow">
            <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t">
              <h3 className="text-xl font-semibold text-gray-900">
                {employeeId
                  ? "Sửa nhân viên"
                  : `Thêm nhân viên ${
                      employeeType === "ROLE_DRIVER" ? "Tài xế" : "Lễ tân"
                    }`}
              </h3>
              <button
                onClick={() => func.closeModal()}
                type="button"
                className="end-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center"
              >
                <svg
                  className="w-3 h-3"
                  aria-hidden="true"
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
              </button>
            </div>
            <div
              className="p-4 md:p-5"
              style={{ maxHeight: "70vh", overflowY: "auto" }}
            >
              <div className="space-y-4">
                {/* Row 1 */}
                <div className="flex flex-col sm:flex-row sm:space-x-4">
                  <div className="basis-1/2">
                    <label
                      htmlFor="fullName"
                      className="block mb-2 text-sm font-medium text-gray-900"
                    >
                      Họ tên
                    </label>
                    <input
                      onChange={(e) => setFullName(e.target.value)}
                      value={fullName}
                      type="text"
                      id="fullName"
                      className={`input-field ${
                        errors.fullName ? "border-red-500" : ""
                      }`}
                      placeholder="Nhập đầy đủ họ tên"
                    />
                    {errors.fullName && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.fullName}
                      </p>
                    )}
                  </div>
                  <div className="basis-1/2 mt-4 sm:mt-0">
                    <label
                      htmlFor="username"
                      className="block mb-2 text-sm font-medium text-gray-900"
                    >
                      Email
                    </label>
                    <input
                      onChange={(e) => setUsername(e.target.value)}
                      value={username}
                      type="email"
                      id="username"
                      className={`input-field ${
                        errors.username ? "border-red-500" : ""
                      }`}
                      placeholder="employee@example.com"
                    />
                    {errors.username && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.username}
                      </p>
                    )}
                  </div>
                </div>
                {/* Row 2 */}
                <div className="flex flex-col sm:flex-row sm:space-x-4">
                  <div className="basis-1/2">
                    <label
                      htmlFor="phoneNumber"
                      className="block mb-2 text-sm font-medium text-gray-900"
                    >
                      Số điện thoại
                    </label>
                    <input
                      onChange={(e) => setPhoneNumber(e.target.value)}
                      value={phoneNumber}
                      type="text"
                      id="phoneNumber"
                      className={`input-field ${
                        errors.phoneNumber ? "border-red-500" : ""
                      }`}
                      placeholder="000x00123121"
                    />
                    {errors.phoneNumber && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.phoneNumber}
                      </p>
                    )}
                  </div>
                  <div className="basis-1/2 mt-4 sm:mt-0">
                    <label
                      htmlFor="identityNumber"
                      className="block mb-2 text-sm font-medium text-gray-900"
                    >
                      Số CCCD
                    </label>
                    <input
                      onChange={(e) => setIdentityNumber(e.target.value)}
                      value={identityNumber}
                      type="text"
                      id="identityNumber"
                      className={`input-field ${
                        errors.identityNumber ? "border-red-500" : ""
                      }`}
                      placeholder="Nhập tại đây"
                    />
                    {errors.identityNumber && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.identityNumber}
                      </p>
                    )}
                  </div>
                </div>
                {/* Row 3 */}
                <div>
                  <label
                    htmlFor="address"
                    className="block mb-2 text-sm font-medium text-gray-900"
                  >
                    Địa chỉ
                  </label>
                  <input
                    onChange={(e) => setAddress(e.target.value)}
                    value={address}
                    type="text"
                    id="address"
                    className={`input-field ${
                      errors.address ? "border-red-500" : ""
                    }`}
                    placeholder="Nhập địa chỉ"
                  />
                  {errors.address && (
                    <p className="text-red-500 text-xs mt-1">
                      {errors.address}
                    </p>
                  )}
                </div>
                {/* Row 4 */}
                <div className="flex flex-col sm:flex-row sm:space-x-4">
                  <div className="basis-1/2">
                    <label
                      htmlFor="gender"
                      className="block mb-2 text-sm font-medium text-gray-900"
                    >
                      Giới tính
                    </label>
                    <select
                      onChange={(e) => setGender(e.target.value)}
                      value={gender}
                      id="gender"
                      className={`select-field ${
                        errors.gender ? "border-red-500" : ""
                      }`}
                    >
                      <option value="">Lựa chọn</option>
                      <option value="MALE">Nam</option>
                      <option value="FEMALE">Nữ</option>
                    </select>
                    {errors.gender && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.gender}
                      </p>
                    )}
                  </div>
                  <div className="basis-1/2 mt-4 sm:mt-0">
                    <label
                      htmlFor="dateOfBirth"
                      className="block mb-2 text-sm font-medium text-gray-900"
                    >
                      Ngày sinh
                    </label>
                    <input
                      onChange={(e) => setDateOfBirth(e.target.value)}
                      value={dateOfBirth}
                      type="date"
                      id="dateOfBirth"
                      className={`input-field ${
                        errors.dateOfBirth ? "border-red-500" : ""
                      }`}
                    />
                    {errors.dateOfBirth && (
                      <p className="text-red-500 text-xs mt-1">
                        {errors.dateOfBirth}
                      </p>
                    )}
                  </div>
                </div>

                {/* Driver Specific Fields */}
                {employeeType === "ROLE_DRIVER" && (
                  <>
                    <hr className="my-4" />
                    <h4 className="text-lg font-semibold text-gray-700 mb-2">
                      Thông tin chi tiết tài xế
                    </h4>
                    {/* Row 5 - Driver */}
                    <div className="flex flex-col sm:flex-row sm:space-x-4">
                      <div className="basis-1/2">
                        <label
                          htmlFor="licenseNumber"
                          className="block mb-2 text-sm font-medium text-gray-900"
                        >
                          Bằng lái xe
                        </label>
                        <input
                          onChange={(e) => setLicenseNumber(e.target.value)}
                          value={licenseNumber}
                          type="text"
                          id="licenseNumber"
                          className={`input-field ${
                            errors.licenseNumber ? "border-red-500" : ""
                          }`}
                          placeholder="Nhập tại đây"
                        />
                        {errors.licenseNumber && (
                          <p className="text-red-500 text-xs mt-1">
                            {errors.licenseNumber}
                          </p>
                        )}
                      </div>
                      <div className="basis-1/2 mt-4 sm:mt-0">
                        <label
                          htmlFor="licenseClass"
                          className="block mb-2 text-sm font-medium text-gray-900"
                        >
                          Cấp độ
                        </label>
                        <input
                          onChange={(e) => setLicenseClass(e.target.value)}
                          value={licenseClass}
                          type="text"
                          id="licenseClass"
                          className={`input-field ${
                            errors.licenseClass ? "border-red-500" : ""
                          }`}
                          placeholder="Ví dụ: B2, C"
                        />
                        {errors.licenseClass && (
                          <p className="text-red-500 text-xs mt-1">
                            {errors.licenseClass}
                          </p>
                        )}
                      </div>
                    </div>
                    {/* Row 6 - Driver */}
                    <div className="flex flex-col sm:flex-row sm:space-x-4">
                      <div className="basis-1/2">
                        <label
                          htmlFor="licenseIssuedDate"
                          className="block mb-2 text-sm font-medium text-gray-900"
                        >
                          Ngày có bằng lái
                        </label>
                        <input
                          onChange={(e) => setLicenseIssuedDate(e.target.value)}
                          value={licenseIssuedDate}
                          type="date"
                          id="licenseIssuedDate"
                          className={`input-field ${
                            errors.licenseIssuedDate ? "border-red-500" : ""
                          }`}
                        />
                        {errors.licenseIssuedDate && (
                          <p className="text-red-500 text-xs mt-1">
                            {errors.licenseIssuedDate}
                          </p>
                        )}
                      </div>
                      <div className="basis-1/2 mt-4 sm:mt-0">
                        <label
                          htmlFor="licenseExpiryDate"
                          className="block mb-2 text-sm font-medium text-gray-900"
                        >
                          Ngày hết hạn
                        </label>
                        <input
                          onChange={(e) => setLicenseExpiryDate(e.target.value)}
                          value={licenseExpiryDate}
                          type="date"
                          id="licenseExpiryDate"
                          className={`input-field ${
                            errors.licenseExpiryDate ? "border-red-500" : ""
                          }`}
                        />
                        {errors.licenseExpiryDate && (
                          <p className="text-red-500 text-xs mt-1">
                            {errors.licenseExpiryDate}
                          </p>
                        )}
                      </div>
                    </div>
                    {/* Row 7 - Driver */}
                    <div>
                      <label
                        htmlFor="vehicleType"
                        className="block mb-2 text-sm font-medium text-gray-900"
                      >
                        Loại xe
                      </label>
                      <input
                        onChange={(e) => setVehicleType(e.target.value)}
                        value={vehicleType}
                        type="text"
                        id="vehicleType"
                        className={`input-field ${
                          errors.vehicleType ? "border-red-500" : ""
                        }`}
                        placeholder="VD: xe máy, xe ô tô, ..."
                      />
                      {errors.vehicleType && (
                        <p className="text-red-500 text-xs mt-1">
                          {errors.vehicleType}
                        </p>
                      )}
                    </div>
                  </>
                )}

                <button
                  onClick={handleSubmit}
                  className="mt-6 text-white bg-blue-700 hover:bg-blue-800 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center"
                >
                  {employeeId ? "Cập nhật" : "Hoàn tất"}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <style jsx>{`
        .input-field {
          background-color: #f9fafb;
          border: 1px solid #d1d5db;
          color: #111827;
          font-size: 0.875rem;
          border-radius: 0.5rem;
          display: block;
          width: 100%;
          padding: 0.625rem;
        }
        .input-field:focus {
          outline: none;
          border-color: #3b82f6;
          box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.25);
        }
        .input-field.border-red-500 {
          border-color: #ef4444;
        }
        .select-field {
          background-color: #f9fafb;
          border: 1px solid #d1d5db;
          color: #111827;
          font-size: 0.875rem;
          border-radius: 0.5rem;
          display: block;
          width: 100%;
          padding: 0.625rem;
        }
        .select-field:focus {
          outline: none;
          border-color: #3b82f6;
          box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.25);
        }
        .select-field.border-red-500 {
          border-color: #ef4444;
        }
      `}</style>
    </div>
  );
};

export default EmployeeForm;
