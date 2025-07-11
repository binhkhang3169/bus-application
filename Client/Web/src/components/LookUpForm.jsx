/** @format */

import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import LookUpResult from "./LookUpResult";
import axios from "axios";
import { API_URL } from "../configs/env";
import FailureNotification from "./Noti/FailureNotification";

function LookUpForm() {
  // Data
  const [ticket, setTicket] = useState(null);

  // UI State
  const [isLoading, setIsLoading] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [message, setMessage] = useState("");

  // Input
  const [phone_number, setPhoneNumber] = useState("");
  const [ticket_id, setTicketId] = useState("");

  const navigate = useNavigate();

  const getTicket = async () => {
    if (!phone_number.trim() || !ticket_id.trim()) {
      setMessage("Vui lòng nhập đầy đủ Số điện thoại và Mã vé.");
      openFailureModal();
      return;
    }

    setIsLoading(true);
    setTicket(null);

    const params = {
      phone: phone_number,
      ticket_id: ticket_id,
    };

    try {
      const res = await axios.post(API_URL + "api/v1/ticket-by-phone", params);

      // Cập nhật đường dẫn để truy cập đúng đối tượng ticket
      if (res.data && res.data.data && res.data.data.ticket) {
        setTicket(res.data.data.ticket);
      } else {
        setMessage("Không tìm thấy thông tin vé phù hợp.");
        openFailureModal();
        setTicket(null); // Đảm bảo không còn dữ liệu cũ
      }
    } catch (err) {
      if (err.response) {
        if (err.response.status === 401) {
          navigate("/admin");
        } else {
          setMessage(
            err.response.data.message ||
              err.response.data.error ||
              "Tra cứu vé thất bại."
          );
          openFailureModal();
        }
      } else if (err.request) {
        setMessage("Lỗi kết nối mạng hoặc máy chủ không phản hồi.");
        openFailureModal();
      } else {
        setMessage("Đã xảy ra lỗi không xác định.");
        openFailureModal();
      }
      setTicket(null); // Đảm bảo không còn dữ liệu cũ khi có lỗi
    } finally {
      setIsLoading(false);
    }
  };

  const closeFailureModal = () => setFailureModal(false);
  const openFailureModal = () => setFailureModal(true);

  return (
    <>
      <div className="lookupform flex flex-col max-w-screen-xl my-10 mx-auto">
        <h1 className="text-green-700 text-2xl text-center font-bold mb-5">
          TRA CỨU THÔNG TIN ĐẶT VÉ
        </h1>
        <hr className="max-w-screen-sm w-full mx-auto h-0.5 bg-gray-200" />
        <div className="max-w-screen-sm w-full flex flex-col gap-y-8 p-5 mx-auto mt-5">
          <input
            onChange={(e) => setPhoneNumber(e.target.value)}
            value={phone_number}
            type="text"
            className="border border-slate-400 rounded-xl p-3"
            placeholder="Vui lòng nhập số điện thoại"
            disabled={isLoading}
          />
          <input
            onChange={(e) => setTicketId(e.target.value)}
            value={ticket_id}
            type="text"
            className="border border-slate-400 rounded-xl p-3"
            placeholder="Vui lòng nhập mã vé"
            onKeyDown={(event) => {
              if (event.key === "Enter" && !isLoading) {
                getTicket();
              }
            }}
            disabled={isLoading}
          />
          <button
            className="font-semibold text-white hover:bg-blue-600 transition-all mx-auto border border-transparent bg-blue-500 px-10 py-3 rounded-full disabled:bg-gray-400"
            onClick={getTicket}
            disabled={isLoading}
          >
            {isLoading ? "Đang tra cứu..." : "Tra cứu vé"}
          </button>
        </div>
      </div>

      <div className="flex flex-col max-w-screen-lg my-10 mx-auto">
        {isLoading && (
          <div className="text-center text-lg">Đang tải dữ liệu...</div>
        )}
        {ticket && <LookUpResult ticket={ticket} />}
      </div>

      {failureModal && (
        <FailureNotification
          func={{ closeModal: closeFailureModal }}
          message={message}
        />
      )}
    </>
  );
}

export default LookUpForm;
