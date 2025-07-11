/** @format */

import React, { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import axios from "axios";
import { API_URL } from "../configs/env";

function NewsDetail() {
  const { id } = useParams(); // Lấy ID từ URL
  const [article, setArticle] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchArticle = async () => {
      if (!id) return;
      setIsLoading(true);
      try {
        // Endpoint lấy chi tiết tin tức từ backend
        const response = await axios.get(`${API_URL}api/v1/news/${id}`);
        setArticle(response.data);
      } catch (err) {
        if (err.response && err.response.status === 404) {
          setError("Không tìm thấy bài viết này.");
        } else {
          setError("Không thể tải dữ liệu bài viết. Vui lòng thử lại sau.");
        }
        console.error("Failed to fetch article:", err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchArticle();
  }, [id]); // useEffect sẽ chạy lại nếu ID trên URL thay đổi

  const formatDate = (dateString) => {
    const options = {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    };
    return new Date(dateString).toLocaleDateString("vi-VN", options);
  };

  if (isLoading) {
    return (
      <div className="text-center text-lg my-10">Đang tải bài viết...</div>
    );
  }

  if (error) {
    return (
      <div className="text-center text-red-500 text-lg my-10">
        <p>{error}</p>
        <Link to="/news" className="text-blue-500 hover:underline mt-4 block">
          Quay lại danh sách tin tức
        </Link>
      </div>
    );
  }

  if (!article) {
    return null; // Hoặc một thông báo không tìm thấy khác
  }

  return (
    <div className="max-w-screen-lg mx-auto my-10 p-5 bg-white shadow-lg rounded-lg">
      <h1 className="text-4xl font-extrabold text-gray-900 mb-4">
        {article.title}
      </h1>
      <div className="text-gray-500 text-md mb-6 border-b pb-4">
        <span>
          Đăng bởi: <strong>{article.created_by}</strong>
        </span>
        <span className="mx-2">|</span>
        <span>Ngày đăng: {formatDate(article.created_at)}</span>
      </div>

      {article.image_url && (
        <img
          src={article.image_url}
          alt={article.title}
          className="w-full h-auto max-h-96 object-cover rounded-md mb-8"
        />
      )}

      {/* Sử dụng div với white-space để giữ lại định dạng xuống dòng từ content */}
      <div
        className="prose lg:prose-xl max-w-none text-gray-800"
        style={{ whiteSpace: "pre-line" }}
      >
        {article.content}
      </div>

      <div className="mt-10 pt-5 border-t">
        <Link
          to="/news"
          className="text-green-600 hover:text-green-800 font-semibold transition-colors duration-300"
        >
          &larr; Quay lại danh sách tin tức
        </Link>
      </div>
    </div>
  );
}

export default NewsDetail;
