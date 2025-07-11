/** @format */

import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import { API_URL } from '../configs/env'; // Giả sử bạn có file config này

function NewsList() {
    const [news, setNews] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchNews = async () => {
            setIsLoading(true);
            try {
                // Endpoint lấy danh sách tin tức từ backend
                // Lưu ý: Endpoint có thể là /api/v1/news tùy vào cấu trúc API Gateway của bạn
                const response = await axios.get(`${API_URL}api/v1/news?limit=20&offset=0`);
                if (response.data && Array.isArray(response.data)) {
                    setNews(response.data);
                } else {
                    setNews([]);
                }
            } catch (err) {
                setError('Không thể tải danh sách tin tức. Vui lòng thử lại sau.');
                console.error('Failed to fetch news:', err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchNews();
    }, []); // Mảng rỗng đảm bảo useEffect chỉ chạy một lần khi component được mount

    const formatDate = (dateString) => {
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        return new Date(dateString).toLocaleDateString('vi-VN', options);
    };

    if (isLoading) {
        return <div className="text-center text-lg my-10">Đang tải dữ liệu...</div>;
    }

    if (error) {
        return <div className="text-center text-red-500 text-lg my-10">{error}</div>;
    }

    return (
        <div className="max-w-screen-xl mx-auto my-10 p-5">
            <h1 className="text-blue-700 text-3xl text-center font-bold mb-8">Tin Tức & Thông Báo</h1>
            <hr className="max-w-screen-sm w-full mx-auto h-0.5 bg-gray-200 mb-10" />
            
            {news.length === 0 ? (
                <p className="text-center text-gray-500">Hiện tại chưa có tin tức nào.</p>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    {news.map((article) => (
                        <Link to={`/news/${article.id}`} key={article.id} className="block group border rounded-lg overflow-hidden shadow-md hover:shadow-xl transition-shadow duration-300">
                            <div className="overflow-hidden">
                                <img 
                                    src={article.image_url || 'https://via.placeholder.com/400x250?text=No+Image'} 
                                    alt={article.title}
                                    className="w-full h-56 object-cover group-hover:scale-105 transition-transform duration-300"
                                />
                            </div>
                            <div className="p-5">
                                <h2 className="text-xl font-bold text-gray-800 mb-2 group-hover:text-green-600 transition-colors duration-300">{article.title}</h2>
                                <p className="text-gray-500 text-sm mb-4">
                                    Ngày đăng: {formatDate(article.created_at)}
                                </p>
                                <p className="text-gray-600 line-clamp-3">
                                    {article.content}
                                </p>
                            </div>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
}

export default NewsList;