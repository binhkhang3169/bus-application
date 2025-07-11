/** @format */

import axios from "axios";
import React, { useEffect, useState } from "react";
import { API_URL } from "../../configs/env";
import { useNavigate } from "react-router-dom";
import SuccessNotification from "../Noti/SuccessNotification";
import FailureNotification from "../Noti/FailureNotification";
import { imageDB } from "../../configs/firebase";
import {
  deleteObject,
  getDownloadURL,
  ref,
  uploadBytes,
} from "firebase/storage";
import { v4 } from "uuid";
import ChangePassword from "./modal/ChangePassword";
import UpdateAccount from "./modal/UpdateAccount";

const AccountManagement = () => {
  const navigate = useNavigate();
  const [account, setAccount] = useState({});
  const [avatar, setAvatar] = useState(
    "https://images.unsplash.com/photo-1618500299034-abce7ed0e8df?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
  );
  const [successModal, setSuccessModal] = useState(false);
  const [failureModal, setFailureModal] = useState(false);
  const [changePasswordModal, setChangePasswordModal] = useState(false);
  const [modelUpdate, setUpdateModal] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    getEmployee();
  }, []);

  const getEmployee = async () => {
    const token = sessionStorage.getItem("token");
    if (token) {
      axios
        .get(API_URL + "employee/me", {
          headers: { Authorization: `Bearer ${token}` },
        })
        .then(async (res) => {
          await renderAvatar(res.data.employee);
          setAccount(res.data.employee);
        })
        .catch((err) => {
          navigate("/admin");
        });
    } else {
      navigate("/admin");
    }
  };

  const renderAvatar = async (employee) => {
    if (employee.avatar) {
      try {
        const url = await getDownloadURL(ref(imageDB, employee.avatar));
        setAvatar(url);
      } catch (error) {
        console.error("Error fetching avatar, using default.", error);
      }
    }
  };

  const closeUpdateModal = () => setUpdateModal(false);
  const closeChangePasswordModal = () => setChangePasswordModal(false);
  const closeSuccessModal = () => setSuccessModal(false);
  const closeFailureModal = () => setFailureModal(false);
  const openSuccessModal = () => setSuccessModal(true);
  const openFailureModal = () => setFailureModal(true);

  const handleChangeImage = async (file) => {
    if (!file) return;
    const url = `avatars/${v4()}`;
    const imageRef = ref(imageDB, url);
    await uploadBytes(imageRef, file);

    let data = { avatar: url };
    const token = sessionStorage.getItem("token");
    if (token) {
      axios
        .put(API_URL + "employee/change-avatar", data, {
          headers: { Authorization: `Bearer ${token}` },
        })
        .then(async (res) => {
          setMessage(res.data.message);
          openSuccessModal();
          if (account.avatar) {
            try {
              const oldImageRef = ref(imageDB, account.avatar);
              await deleteObject(oldImageRef);
            } catch (deleteError) {
              console.error(
                "Failed to delete old avatar, it might not exist.",
                deleteError
              );
            }
          }
          getEmployee();
        })
        .catch((err) => {
          if (err.response.status === 401) navigate("/admin");
          setMessage(err.response.data.message);
          openFailureModal();
        });
    } else {
      navigate("/admin");
    }
  };

  return (
    <>
      <div className="w-full px-4">
        <h3 className="ml-12 lg:ml-0 h-14 text-2xl font-semibold flex items-center dark:text-white">
          My Account
        </h3>
        <div className="w-full mt-2 border border-gray-200 rounded-xl p-3 flex flex-col md:flex-row dark:border-gray-700 dark:bg-gray-900">
          <div className="basis-1/3 flex flex-col p-2">
            <div className="flex justify-center mx-auto md:mx-0">
              <img
                src={avatar}
                alt="avatar"
                className="aspect-square w-full max-w-[200px] rounded-full object-cover border-4 border-white dark:border-gray-800 shadow-lg"
              />
            </div>
            <div className="mb-5 mt-8 text-center">
              <label
                htmlFor="image"
                className="cursor-pointer bg-gray-200 px-4 py-2 rounded-full hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-200"
              >
                Choose picture
              </label>
              <input
                type="file"
                accept=".jpg,.png"
                id="image"
                onChange={(e) => handleChangeImage(e.target.files[0])}
                className="hidden"
              />
            </div>
            <div className="text-center text-slate-500 dark:text-gray-400 text-xs">
              The maximum file size is 1 MB, and the only accepted formats are
              JPEG and PNG.
            </div>
          </div>
          <div className="basis-2/3 w-full flex flex-col p-3 md:p-5">
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                Last Name
              </div>
              <div className="basis-2/3 text-gray-800 dark:text-gray-200">
                : {account.last_name}
              </div>
            </div>
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                First Name
              </div>
              <div className="basis-2/3 text-gray-800 dark:text-gray-200">
                : {account.first_name}
              </div>
            </div>
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                Phone
              </div>
              <div className="basis-2/3 text-gray-800 dark:text-gray-200">
                : {account.phone_number}
              </div>
            </div>
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                Gender
              </div>
              <div className="basis-2/3 text-gray-800 dark:text-gray-200">
                : {account.gender === 0 ? "Female" : "Male"}
              </div>
            </div>
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                Email
              </div>
              <div className="basis-2/3 text-gray-800 dark:text-gray-200">
                : {account.email}
              </div>
            </div>
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                Date of Birth
              </div>
              <div className="basis-2/3 text-gray-800 dark:text-gray-200">
                : {account.date_of_birth}
              </div>
            </div>
            <div className="flex flex-row mb-3 items-center">
              <div className="basis-1/3 text-slate-500 dark:text-gray-400">
                Address
              </div>
              <div className="basis-2/3 line-clamp-1 text-gray-800 dark:text-gray-200">
                : {account.address}
              </div>
            </div>
            <div className="flex justify-center md:justify-start gap-4 mt-4">
              <button
                onClick={() => setUpdateModal(true)}
                className="bg-blue-600 px-8 py-2 rounded-full text-white hover:bg-blue-700 transition-colors cursor-pointer"
              >
                Update
              </button>
              <button
                onClick={() => setChangePasswordModal(true)}
                className="bg-blue-600 px-8 py-2 rounded-full text-white hover:bg-blue-700 transition-colors cursor-pointer"
              >
                Change password
              </button>
            </div>
          </div>
        </div>
      </div>

      {changePasswordModal && (
        <ChangePassword
          closeModal={closeChangePasswordModal}
          refresh={getEmployee}
          setMessage={setMessage}
          openFailureModal={openFailureModal}
          openSuccessModal={openSuccessModal}
        />
      )}
      {modelUpdate && (
        <UpdateAccount
          closeModal={closeUpdateModal}
          refresh={getEmployee}
          setMessage={setMessage}
          openFailureModal={openFailureModal}
          openSuccessModal={openSuccessModal}
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
    </>
  );
};

export default AccountManagement;
