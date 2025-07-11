/** @format */

import React, { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import api from "../services/apiService"; // Use the centralized API service
import CustomerAccountUpdateModal from "./CustomerAccountUpdateModal";
import SuccessNotification from "./Noti/SuccessNotification";
import FailureNotification from "./Noti/FailureNotification";
import CustomerChangePasswordModal from "./CustomerChangePasswordModal";
import TicketDetailModal from "./TicketDetailModal"; // --- IMPORT TicketDetailModal ---

// Firebase imports - ensure firebase is configured in 'src/configs/firebase.js'
import { imageDB } from "../configs/firebase";
import {
  deleteObject,
  getDownloadURL,
  ref,
  uploadBytes,
} from "firebase/storage";
import { v4 as uuidv4 } from "uuid";

function ProfileManagementForm() {
  const navigate = useNavigate();

  // Data
  const [customer, setCustomer] = useState(null);
  const [tickets, setTickets] = useState([]);
  const [loadingTickets, setLoadingTickets] = useState(false);

  // Modal states
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [showFailureModal, setShowFailureModal] = useState(false);
  const [showChangePasswordModal, setShowChangePasswordModal] = useState(false);
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [modalMessage, setModalMessage] = useState("");

  // --- NEW STATE for Ticket Detail Modal ---
  const [showTicketDetailModal, setShowTicketDetailModal] = useState(false);
  const [selectedTicketDetail, setSelectedTicketDetail] = useState(null);
  const [rootApiUrl, setRootApiUrl] = useState("");

  // Avatar state
  const [avatarUrl, setAvatarUrl] = useState(
    "https://images.unsplash.com/photo-1618500299034-abce7ed0e8df?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
  );

  const [activeSection, setActiveSection] = useState("account");

  // --- Function to get Root API URL for QR Codes ---
  useEffect(() => {
    if (api && api.defaults && api.defaults.baseURL) {
      try {
        const url = new URL(api.defaults.baseURL);
        setRootApiUrl(url.origin); // e.g., "http://localhost:8080"
      } catch (error) {
        console.error("Error parsing API base URL:", error);
        // Fallback or default if parsing fails, adjust as needed
        setRootApiUrl("http://localhost:8080");
      }
    } else {
      // Fallback if api.defaults.baseURL is not available
      setRootApiUrl("http://localhost:8080");
    }
  }, []);

  const fetchCustomerInfo = useCallback(async () => {
    const accessToken =
      localStorage.getItem("accessToken") ||
      sessionStorage.getItem("accessToken");
    if (!accessToken) {
      console.log("ProfileManagement: No token found, navigating to login.");
      navigate("/login");
      return;
    }
    try {
      console.log("ProfileManagement: Fetching user-info...");
      const response = await api.get("customer/info");
      if (response.data && response.status === 200) {
        setCustomer(response.data.data);
        if (
          response.data.data.avatar &&
          response.data.data.avatar.startsWith("avatars/")
        ) {
          console.log(
            "ProfileManagement: Attempting to render avatar from Firebase path:",
            response.data.data.avatar
          );
          try {
            const firebaseUrl = await getDownloadURL(
              ref(imageDB, response.data.data.avatar)
            );
            setAvatarUrl(firebaseUrl);
          } catch (fbError) {
            console.error("Error fetching avatar URL from Firebase:", fbError);
            setAvatarUrl(
              "https://images.unsplash.com/photo-1618500299034-abce7ed0e8df?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
            );
          }
        } else if (response.data.data.avatar) {
          setAvatarUrl(response.data.data.avatar);
        } else {
          console.log(
            "ProfileManagement: No Firebase avatar path found or invalid format."
          );
          setAvatarUrl(
            "https://images.unsplash.com/photo-1618500299034-abce7ed0e8df?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
          );
        }
      } else {
        console.warn(
          "ProfileManagement: User info response not as expected.",
          response
        );
        navigate("/login");
      }
    } catch (error) {
      console.error(
        "ProfileManagement: Error fetching customer info:",
        error.response?.data?.message || error.message
      );
      if (error.response?.status !== 401) {
        navigate("/login");
      }
    }
  }, [navigate]);

  const fetchTickets = useCallback(async () => {
    const accessToken =
      localStorage.getItem("accessToken") ||
      sessionStorage.getItem("accessToken");
    if (!accessToken) {
      return;
    }
    setLoadingTickets(true);
    try {
      console.log("ProfileManagement: Fetching tickets...");
      const response = await api.get("tickets");
      if (response.data && response.status === 200) {
        setTickets(response.data.data || []);
        console.log("ProfileManagement: Tickets fetched: ", response.data.data);
      } else {
        console.warn(
          "ProfileManagement: Tickets response not as expected.",
          response
        );
        setTickets([]);
        setModalMessage(
          response.data?.message || "Không thể tải danh sách vé."
        );
        setShowFailureModal(true);
      }
    } catch (error) {
      console.error(
        "ProfileManagement: Error fetching tickets:",
        error.response?.data?.message || error.message
      );
      setTickets([]);
      if (error.response?.status !== 401) {
        setModalMessage(
          error.response?.data?.message || "Lỗi khi tải danh sách vé."
        );
        setShowFailureModal(true);
      }
    } finally {
      setLoadingTickets(false);
    }
  }, []);

  useEffect(() => {
    fetchCustomerInfo();
    fetchTickets();
  }, [fetchCustomerInfo, fetchTickets]);

  const handleLogout = async () => {
    try {
      console.log("ProfileManagement: Attempting logout via API.");
      await api.post("auth/logout");
      setModalMessage("Đăng xuất thành công!");
    } catch (error) {
      console.error(
        "ProfileManagement: Logout API call failed:",
        error.response?.data?.message || error.message
      );
      setModalMessage("Đăng xuất thất bại hoặc có lỗi xảy ra.");
    } finally {
      localStorage.removeItem("accessToken");
      sessionStorage.removeItem("accessToken");
      localStorage.removeItem("username");
      sessionStorage.removeItem("username");
      localStorage.removeItem("userRole");
      sessionStorage.removeItem("userRole");
      localStorage.removeItem("userFullName");
      sessionStorage.removeItem("userFullName");
      localStorage.removeItem("refreshToken");
      sessionStorage.removeItem("refreshToken");
      window.dispatchEvent(
        new CustomEvent("authChange", { detail: { loggedIn: false } })
      );
      console.log("ProfileManagement: Navigating to home after logout.");
      navigate("/");
    }
  };

  const openChangePasswordModalHandler = () => setShowChangePasswordModal(true);

  const handleChangeAvatar = async (file) => {
    if (!file) return;

    const oldAvatarPath = customer?.avatar;
    const newAvatarFirebasePath = `avatars/${uuidv4()}-${file.name}`;
    const imageRef = ref(imageDB, newAvatarFirebasePath);

    try {
      console.log(
        `ProfileManagement: Uploading new avatar to Firebase: ${newAvatarFirebasePath}`
      );
      await uploadBytes(imageRef, file);
      const newFirebaseUrl = await getDownloadURL(imageRef);

      try {
        console.log(
          `ProfileManagement: Calling backend to update avatar path to: ${newAvatarFirebasePath}`
        );
        await api.put("customer/change-avatar", {
          avatar: newAvatarFirebasePath,
        });

        setModalMessage("Thay đổi ảnh đại diện thành công!");
        setShowSuccessModal(true);
        setAvatarUrl(newFirebaseUrl);
        setCustomer((prev) => ({ ...prev, avatar: newAvatarFirebasePath }));

        if (
          oldAvatarPath &&
          oldAvatarPath !== newAvatarFirebasePath &&
          oldAvatarPath.startsWith("avatars/")
        ) {
          console.log(
            `ProfileManagement: Deleting old avatar from Firebase: ${oldAvatarPath}`
          );
          const oldImageRef = ref(imageDB, oldAvatarPath);
          await deleteObject(oldImageRef).catch((delError) =>
            console.error("Error deleting old Firebase avatar:", delError)
          );
        }
      } catch (backendError) {
        console.error(
          "ProfileManagement: Backend failed to update avatar:",
          backendError.response?.data?.message || backendError.message
        );
        setModalMessage(
          backendError.response?.data?.message ||
            "Lỗi khi cập nhật ảnh đại diện ở server."
        );
        setShowFailureModal(true);
        await deleteObject(imageRef).catch((delError) =>
          console.error(
            "Error deleting new Firebase avatar after backend failure:",
            delError
          )
        );
      }
    } catch (error) {
      console.error(
        "ProfileManagement: Error during avatar change process:",
        error
      );
      setModalMessage("Lỗi khi tải lên hoặc xử lý ảnh đại diện.");
      setShowFailureModal(true);
    }
  };

  const handlePayment = async (ticket) => {
    const token =
      localStorage.getItem("accessToken") ||
      sessionStorage.getItem("accessToken");
    if (!token) {
      setModalMessage("Phiên đăng nhập hết hạn. Vui lòng đăng nhập lại.");
      setShowFailureModal(true);
      navigate("/login");
      return;
    }

    if (
      !ticket ||
      typeof ticket.Price === "undefined" ||
      !ticket.TicketID ||
      !ticket.trip_id
    ) {
      console.error("Invalid ticket data for payment:", ticket);
      setModalMessage("Dữ liệu vé không hợp lệ để thanh toán.");
      setShowFailureModal(true);
      return;
    }

    const customerIdForPayment = ticket.CustomerID?.Int32?.toString();
    if (!customerIdForPayment) {
      console.error("Missing CustomerID in ticket data:", ticket.CustomerID);
      setModalMessage("Thiếu thông tin ID khách hàng trong vé.");
      setShowFailureModal(true);
      return;
    }

    const paymentGatewayPayload = {
      amount: ticket.Price,
      bank_code: "NCB",
      language: "vn",
      invoice_type: "bus_ticket",
      customer_id: customerIdForPayment,
      ticket_id: ticket.TicketID,
      notes: `Thanh toán cho vé ${ticket.TicketID} của chuyến ${ticket.trip_id}`,
    };

    console.log("VNPay Payload:", paymentGatewayPayload);

    try {
      const paymentGatewayResponse = await api.post(
        "vnpay/create-payment",
        paymentGatewayPayload
      );
      console.log("Payment Gateway Response:", paymentGatewayResponse.data);

      if (
        paymentGatewayResponse.data &&
        paymentGatewayResponse.data.data?.payment_url
      ) {
        setModalMessage("Đang chuyển hướng đến cổng thanh toán...");
        setShowSuccessModal(true);
        window.location.href = paymentGatewayResponse.data.data.payment_url;
      } else {
        console.error(
          "Payment initiation failed:",
          paymentGatewayResponse.data
        );
        setModalMessage(
          paymentGatewayResponse.data?.message ||
            "Không thể khởi tạo thanh toán. Vui lòng thử lại."
        );
        setShowFailureModal(true);
      }
    } catch (error) {
      console.error(
        "Error calling payment gateway:",
        error.response?.data || error.message
      );
      setModalMessage(
        error.response?.data?.message || "Lỗi kết nối đến cổng thanh toán."
      );
      setShowFailureModal(true);
    }
  };

  // Modal close handlers
  const closeUpdateModal = () => setShowUpdateModal(false);
  const closeChangePasswordModal = () => setShowChangePasswordModal(false);
  const closeSuccessModal = () => setShowSuccessModal(false);
  const closeFailureModal = () => setShowFailureModal(false);

  // --- Handlers for Ticket Detail Modal ---
  const openTicketDetailModalHandler = (ticket) => {
    setSelectedTicketDetail(ticket);
    setShowTicketDetailModal(true);
  };
  const closeTicketDetailModalHandler = () => {
    setShowTicketDetailModal(false);
    setSelectedTicketDetail(null);
  };

  const openUpdateModalHandler = () => {
    if (customer) {
      setShowUpdateModal(true);
    } else {
      console.log("Cannot open update modal, customer data not loaded.");
    }
  };

  if (!customer) {
    return (
      <div className="flex justify-center items-center h-screen">
        <p>Đang tải thông tin người dùng...</p>
      </div>
    );
  }

  const formatPrice = (price) => {
    console.log("Price:", price);
    if (typeof price !== "number") return "N/A";
    return price.toLocaleString("vi-VN", {
      style: "currency",
      currency: "VND",
    });
  };

  const formatDate = (dateString) => {
    if (!dateString) return "N/A";
    try {
      return new Date(dateString).toLocaleString("vi-VN", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit", // Added seconds for more detail if needed
      });
    } catch (e) {
      return "Invalid Date";
    }
  };

  return (
    <div className="flex-1 p-5">
      <div className="flex flex-col-reverse md:flex-row w-full lg:w-[80%] 2xl:w-[70%] mx-auto gap-y-8 md:gap-x-8 my-8">
        <div className="flex basis-full md:basis-1/4 border border-slate-300 p-2 flex-col rounded-xl h-fit">
          <div
            className={`flex flex-row p-2 mb-2 items-center hover:bg-slate-100 cursor-pointer rounded ${
              activeSection === "account" ? "bg-slate-200 font-semibold" : ""
            }`}
            onClick={() => setActiveSection("account")}
          >
            <div className="basis-1/4">
              <img
                src="https://cdn-icons-png.flaticon.com/512/1177/1177568.png"
                alt="info"
                className="w-6 h-6"
              />
            </div>
            <div className="basis-3/4 ml-3">Thông tin tài khoản</div>
          </div>
          <div
            className={`flex flex-row p-2 mb-2 items-center hover:bg-slate-100 cursor-pointer rounded ${
              activeSection === "tickets" ? "bg-slate-200 font-semibold" : ""
            }`}
            onClick={() => setActiveSection("tickets")}
          >
            <div className="basis-1/4">
              <img
                src="https://cdn-icons-png.flaticon.com/512/2990/2990693.png"
                alt="tickets"
                className="w-6 h-6"
              />
            </div>
            <div className="basis-3/4 ml-3">Vé của tôi</div>
          </div>
          <div
            className="flex flex-row p-2 mb-2 items-center hover:bg-slate-100 cursor-pointer rounded"
            onClick={openChangePasswordModalHandler}
          >
            <div className="basis-1/4">
              <img
                src="https://futabus.vn/images/header/profile/Password.svg"
                alt="password"
              />
            </div>
            <div className="basis-3/4 ml-3">Đặt lại mật khẩu</div>
          </div>
          <div
            className="flex flex-row p-2 mb-2 items-center hover:bg-slate-100 cursor-pointer rounded"
            onClick={handleLogout}
          >
            <div className="basis-1/4">
              <img
                src="https://futabus.vn/images/header/profile/Logout.svg"
                alt="logout"
              />
            </div>
            <div className="basis-3/4 ml-3">Đăng xuất</div>
          </div>
        </div>

        <div className="md:basis-3/4 w-full">
          {activeSection === "account" && (
            <div id="account-info" className="mb-12">
              <h3 className="text-2xl font-semibold mb-2">
                Thông tin tài khoản
              </h3>
              <p className="text-sm text-slate-500 mb-5">
                Quản lý thông tin hồ sơ để bảo mật tài khoản
              </p>
              <div className="w-full border border-slate-300 rounded-xl p-3 flex flex-col md:flex-row">
                <div className="basis-1/3 flex flex-col p-2 items-center md:items-start">
                  <div className="flex justify-center">
                    <img
                      src={avatarUrl}
                      alt="avatar"
                      className="aspect-square w-[150px] h-[150px] md:w-[200px] md:h-[200px] rounded-full object-cover border border-slate-200"
                    />
                  </div>
                  <div className="mb-5 mt-8 text-center w-full">
                    <label
                      htmlFor="imageUpload"
                      className="cursor-pointer bg-slate-200 px-4 py-2 rounded-full hover:bg-slate-300 text-sm"
                    >
                      Chọn ảnh
                    </label>
                    <input
                      type="file"
                      accept=".jpg,.jpeg,.png"
                      id="imageUpload"
                      onChange={(e) => handleChangeAvatar(e.target.files[0])}
                      className="hidden"
                    />
                  </div>
                  <div className="text-center text-slate-500 text-xs w-full">
                    Dung lượng file tối đa 1 MB
                    <br />
                    Định dạng: .JPEG, .PNG
                  </div>
                </div>

                <div className="basis-2/3 w-full flex flex-col p-3 md:p-5">
                  <div className="flex flex-row mb-3 items-center">
                    <div className="basis-1/3 text-slate-500">Tên</div>
                    <div className="basis-2/3 text-gray-800">
                      :{" "}
                      <span className="ml-3">{customer.fullName || "N/A"}</span>
                    </div>
                  </div>
                  <div className="flex flex-row mb-3 items-center">
                    <div className="basis-1/3 text-slate-500">
                      Số điện thoại
                    </div>
                    <div className="basis-2/3 text-gray-800">
                      :{" "}
                      <span className="ml-3">
                        {customer.phoneNumber || "N/A"}
                      </span>
                    </div>
                  </div>
                  <div className="flex flex-row mb-3 items-center">
                    <div className="basis-1/3 text-slate-500">Giới tính</div>
                    <div className="basis-2/3 text-gray-800">
                      : <span className="ml-3">{customer.gender || "N/A"}</span>
                    </div>
                  </div>
                  <div className="flex flex-row mb-3 items-center">
                    <div className="basis-1/3 text-slate-500">
                      Tên đăng nhập
                    </div>
                    <div className="basis-2/3 text-gray-800">
                      :{" "}
                      <span className="ml-3">{customer.username || "N/A"}</span>
                    </div>
                  </div>
                  <div className="flex flex-row mb-3 items-center">
                    <div className="basis-1/3 text-slate-500">Địa chỉ</div>
                    <div className="basis-2/3 text-gray-800">
                      :{" "}
                      <span className="ml-3">{customer.address || "N/A"}</span>
                    </div>
                  </div>
                  <div
                    className="flex justify-center items-center bg-blue-500 mx-auto px-8 py-2 rounded-full text-white mt-3 hover:bg-blue-600 transition-colors cursor-pointer w-fit"
                    onClick={openUpdateModalHandler}
                  >
                    Cập nhật thông tin
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeSection === "tickets" && (
            <div id="my-tickets" className="w-full">
              <h3 className="text-2xl font-semibold mb-2">Vé của tôi</h3>
              <p className="text-sm text-slate-500 mb-5">
                Danh sách các vé đã đặt
              </p>
              {loadingTickets ? (
                <p>Đang tải danh sách vé...</p>
              ) : tickets && tickets.length > 0 ? (
                <div className="overflow-x-auto border border-slate-300 rounded-xl">
                  <table className="min-w-full divide-y divide-slate-200">
                    <thead className="bg-slate-50">
                      <tr>
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          Mã vé
                        </th>
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          Tên người đặt
                        </th>
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          SĐT
                        </th>
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          Giá vé
                        </th>
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          Thời gian đặt
                        </th>
                        
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          Trạng thái TT
                        </th>
                        <th
                          scope="col"
                          className="px-4 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider"
                        >
                          Hành động
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-slate-200">
                      {tickets.map((ticket) => (
                        <tr key={ticket.TicketID}>
                          <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">
                            {ticket.ticket_id}
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">
                            {ticket.name?.String || "N/A"}
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">
                            {ticket.phone?.String || "N/A"}
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">
                            {formatPrice(ticket.price)}
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">
                            {formatDate(ticket.booking_time)}
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap text-sm">
                            {ticket.payment_status === 1 ? (
                              <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                Đã thanh toán
                              </span>
                            ) : ticket.payment_status === 0 ? (
                              <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800">
                                Chưa thanh toán
                              </span>
                            ) : (
                              <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
                                Đã hủy
                              </span>
                            )}
                          </td>
                          {/* --- UPDATED Actions Column --- */}
                          <td className="px-4 py-3 whitespace-nowrap text-sm font-medium space-x-2">
                            {ticket.PaymentStatus === 0 &&
                              ticket.Status !== 2 && (
                                <button
                                  onClick={() => handlePayment(ticket)}
                                  className="text-indigo-600 hover:text-indigo-900 bg-indigo-100 hover:bg-indigo-200 px-3 py-1 rounded-md text-xs"
                                >
                                  Thanh toán
                                </button>
                              )}
                            {ticket.payment_status === 1 &&
                              ticket.status !== 2 && (
                                <span className="text-green-600 px-3 py-1 text-xs">
                                  Hoàn tất
                                </span>
                              )}
                            {ticket.status === 2 && (
                              <span className="text-red-600 px-3 py-1 text-xs">
                                Đã hủy
                              </span>
                            )}
                            <button
                              onClick={() =>
                                openTicketDetailModalHandler(ticket)
                              }
                              className="text-sky-600 hover:text-sky-900 bg-sky-100 hover:bg-sky-200 px-3 py-1 rounded-md text-xs"
                            >
                              Chi tiết
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <p className="text-slate-600">Không tìm thấy vé nào.</p>
              )}
            </div>
          )}
        </div>
      </div>

      {showChangePasswordModal && (
        <CustomerChangePasswordModal
          closeModal={closeChangePasswordModal}
          setMessage={setModalMessage}
          openFailureModal={() => setShowFailureModal(true)}
          openSuccessModal={() => setShowSuccessModal(true)}
        />
      )}
      {showUpdateModal && customer && (
        <CustomerAccountUpdateModal
          closeModal={closeUpdateModal}
          refresh={() => {
            fetchCustomerInfo();
          }}
          setMessage={setModalMessage}
          openFailureModal={() => setShowFailureModal(true)}
          openSuccessModal={() => setShowSuccessModal(true)}
          currentUserData={customer}
        />
      )}
      {showSuccessModal && (
        <SuccessNotification
          func={{ closeModal: closeSuccessModal }}
          message={modalMessage}
        />
      )}
      {showFailureModal && (
        <FailureNotification
          func={{ closeModal: closeFailureModal }}
          message={modalMessage}
        />
      )}
      {/* --- RENDER TicketDetailModal --- */}
      {showTicketDetailModal && selectedTicketDetail && rootApiUrl && (
        <TicketDetailModal
          isOpen={showTicketDetailModal}
          onClose={closeTicketDetailModalHandler}
          ticket={selectedTicketDetail}
          formatPrice={formatPrice}
          formatDate={formatDate}
          rootApiUrl={rootApiUrl}
        />
      )}
    </div>
  );
}

export default ProfileManagementForm;
