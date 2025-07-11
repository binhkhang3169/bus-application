/** @format */

import React, { useEffect, useState } from "react";
import axios from "axios";
import { API_URL } from "../../../configs/env";
import { useNavigate } from "react-router-dom";

const NewsForm = ({ func, newsId }) => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  // Form fields
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [createdBy, setCreatedBy] = useState("Admin");

  const [formErrors, setFormErrors] = useState({});

  useEffect(() => {
    if (newsId) {
      setIsLoading(true);
      const fetchNewsDetails = async () => {
        try {
          const token = sessionStorage.getItem("adminAccessToken");
          const response = await axios.get(`${API_URL}api/v1/news/${newsId}`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          const newsData = response.data;
          setTitle(newsData.title);
          setContent(newsData.content);
          setImageUrl(newsData.image_url || newsData.imageUrl);
          setCreatedBy(newsData.created_by);
        } catch (error) {
          func.setMessage("Failed to load news details.");
          func.openFailureModal();
        } finally {
          setIsLoading(false);
        }
      };
      fetchNewsDetails();
    }
  }, [newsId, func]);

  const handleImageSelectAndUpload = async (file) => {
    if (!file) return;

    if (file.size > 1024 * 1024) {
      // 1MB limit
      func.setMessage("Dung lượng file tối đa là 1 MB.");
      func.openFailureModal();
      return;
    }

    const formData = new FormData();
    formData.append("image", file);

    setIsUploading(true);
    try {
      const token = sessionStorage.getItem("adminAccessToken");
      const response = await axios.post(
        `${API_URL}api/v1/upload/image`,
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.data.success) {
        setImageUrl(response.data.cloudinary_url);
      } else {
        throw new Error(
          response.data.Error || "Lỗi không xác định khi tải ảnh."
        );
      }
    } catch (error) {
      const errorMsg =
        error.response?.data?.Error || error.message || "Tải ảnh thất bại.";
      func.setMessage(errorMsg);
      func.openFailureModal();
    } finally {
      setIsUploading(false);
    }
  };

  const validateForm = () => {
    const errors = {};
    if (!title.trim()) errors.title = "Tiêu đề là bắt buộc.";
    if (!content.trim()) errors.content = "Nội dung là bắt buộc.";
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;

    setIsLoading(true);
    const payload = {
      title,
      content,
      image_url: imageUrl,
      created_by: createdBy,
    };

    try {
      const token = sessionStorage.getItem("adminAccessToken");
      let response;
      if (newsId) {
        response = await axios.put(`${API_URL}api/v1/news/${newsId}`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
      } else {
        response = await axios.post(`${API_URL}api/v1/news`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
      }

      func.setMessage(`Đã ${newsId ? "cập nhật" : "tạo"} tin tức thành công!`);
      func.openSuccessModal();
      func.refresh();
      func.closeModal();
    } catch (error) {
      const errorMsg = error.response?.data?.error || `Thao tác thất bại.`;
      func.setMessage(errorMsg);
      func.openFailureModal();
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed z-50 top-0 left-0 bg-black/50 w-full h-full flex items-center justify-center">
      {/* DARK MODE: Thêm background cho modal panel */}
      <div className="relative w-full max-w-3xl bg-white rounded-lg shadow-lg max-h-[90vh] overflow-y-auto dark:bg-gray-800">
        {/* DARK MODE: Thêm background và border cho header */}
        <div className="flex items-center justify-between p-4 border-b sticky top-0 bg-white z-10 dark:bg-gray-800 dark:border-gray-700">
          {/* DARK MODE: Thêm text color cho title */}
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
            {newsId ? "Chỉnh sửa Tin tức" : "Tạo Tin tức mới"}
          </h3>
          <button
            onClick={func.closeModal}
            type="button"
            className="text-gray-400 bg-transparent hover:bg-gray-200 rounded-lg text-sm w-8 h-8 inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
          >
            <svg className="w-3 h-3" fill="none" viewBox="0 0 14 14">
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

        <div className="p-5">
          <form className="space-y-4" onSubmit={(e) => e.preventDefault()}>
            <div className="flex flex-col md:flex-row gap-6">
              {/* Left side: Image Upload */}
              <div className="basis-1/3 flex flex-col p-2 items-center">
                <div className="relative flex justify-center w-[200px] h-[200px]">
                  <img
                    src={
                      imageUrl ||
                      "https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg"
                    }
                    alt="Preview"
                    // DARK MODE: Thêm border color
                    className="aspect-square w-[200px] h-[200px] rounded-full object-cover border-2 border-slate-200 dark:border-gray-600"
                    key={imageUrl}
                  />
                  {isUploading && (
                    <div className="absolute inset-0 bg-black/50 flex items-center justify-center rounded-full">
                      <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-white"></div>
                    </div>
                  )}
                </div>
                <div className="mb-5 mt-8 text-center w-full">
                  {/* DARK MODE: Thêm background và hover state */}
                  <label
                    htmlFor="imageUpload"
                    className="cursor-pointer bg-slate-200 px-4 py-2 rounded-full hover:bg-slate-300 text-sm dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-300"
                  >
                    Chọn ảnh
                  </label>
                  <input
                    type="file"
                    accept=".jpg,.jpeg,.png"
                    id="imageUpload"
                    onChange={(e) =>
                      handleImageSelectAndUpload(e.target.files[0])
                    }
                    className="hidden"
                    disabled={isUploading}
                  />
                </div>
                {/* DARK MODE: Thêm text color */}
                <div className="text-center text-slate-500 text-xs w-full dark:text-gray-400">
                  Dung lượng file tối đa 1 MB
                  <br />
                  Định dạng: .JPEG, .PNG
                </div>
              </div>

              {/* Right side: Form Fields */}
              <div className="basis-2/3 flex flex-col space-y-4">
                <div>
                  {/* DARK MODE: Thêm text color */}
                  <label
                    htmlFor="title"
                    className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                  >
                    Tiêu đề <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    id="title"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    // DARK MODE: Thêm các class cho input field
                    className={`bg-gray-50 border ${
                      formErrors.title ? "border-red-500" : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white`}
                  />
                  {formErrors.title && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.title}
                    </p>
                  )}
                </div>
                <div>
                  {/* DARK MODE: Thêm text color */}
                  <label
                    htmlFor="content"
                    className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                  >
                    Nội dung <span className="text-red-500">*</span>
                  </label>
                  <textarea
                    id="content"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    rows="12"
                    // DARK MODE: Thêm các class cho textarea
                    className={`bg-gray-50 border ${
                      formErrors.content ? "border-red-500" : "border-gray-300"
                    } text-gray-900 text-sm rounded-lg block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white`}
                  ></textarea>
                  {formErrors.content && (
                    <p className="text-red-500 text-xs mt-1">
                      {formErrors.content}
                    </p>
                  )}
                </div>
              </div>
            </div>

            {/* DARK MODE: Thêm border cho footer */}
            <div className="flex justify-end pt-4 border-t dark:border-gray-700">
              {/* DARK MODE: Thêm các class cho nút Hủy */}
              <button
                type="button"
                onClick={func.closeModal}
                className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700"
              >
                Hủy
              </button>
              <button
                type="button"
                onClick={handleSubmit}
                disabled={isLoading || isUploading}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-700 rounded-lg hover:bg-blue-800 flex items-center disabled:bg-blue-400"
              >
                {(isLoading || isUploading) && (
                  <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-white mr-2"></div>
                )}
                {newsId ? "Cập nhật" : "Tạo mới"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default NewsForm;
