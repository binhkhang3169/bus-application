/** @format */

import React, { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import FailureNotification from "./Noti/FailureNotification";
import BookingStatus from "./BookingStatus"; // Component màn hình chờ và kết quả
import api from "../services/apiService";
import { API_URL } from "../configs/env";
import TripBookingSection from "./TripBookingSection"; // Giả định component này đã tồn tại

// =================================================================
// BookingTicketForm Component (Logic hoàn chỉnh)
// =================================================================
function BookingTicketForm() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const departureTripId = searchParams.get("departureTripId");
    const returnTripId = searchParams.get("returnTripId");

    const [departureTripInfo, setDepartureTripInfo] = useState(null);
    const [returnTripInfo, setReturnTripInfo] = useState(null);
    
    const [bookingDetails, setBookingDetails] = useState({
        departure: { seats: [], seatIds: [], pickupId: '', dropoffId: '' },
        return: { seats: [], seatIds: [], pickupId: '', dropoffId: '' }
    });

    // Giá trị mặc định để test, tương tự file index.html
    const [customer, setCustomer] = useState({ name: "", phone_number: "", email: "" });
    const [initialLoading, setInitialLoading] = useState(true);
    const [termsAccepted, setTermsAccepted] = useState(false);
    const [failureModal, setFailureModal] = useState(false);
    const [message, setMessage] = useState("");

    // --- STATE MỚI CHO QUY TRÌNH ĐẶT VÉ ---
    const [bookingState, setBookingState] = useState('IDLE'); // IDLE, PROCESSING
    const [bookingId, setBookingId] = useState(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    
    // --- STATE CHO GIAO DIỆN TABS (Lượt đi/Lượt về) ---
    const [activeTripView, setActiveTripView] = useState('departure');

    useEffect(() => {
        if (!departureTripId) {
            setMessage("Không tìm thấy thông tin chuyến đi. Vui lòng chọn lại.");
            openFailureModal();
            setTimeout(() => navigate("/"), 3000);
            return;
        }

        const fetchAllTripData = async () => {
            setInitialLoading(true);
            try {
                const departureData = await fetchTripAndSeatDetails(departureTripId);
                setDepartureTripInfo(departureData);

                if (returnTripId) {
                    const returnData = await fetchTripAndSeatDetails(returnTripId);
                    setReturnTripInfo(returnData);
                }
                loadCustomer();
            } catch (error) {
                setMessage(error.message || "Đã xảy ra lỗi khi tải dữ liệu chuyến đi.");
                openFailureModal();
            } finally {
                setInitialLoading(false);
            }
        };

        fetchAllTripData();
    }, [departureTripId, returnTripId, navigate]);

    const loadCustomer = () => {
       const token = localStorage.getItem("accessToken") || sessionStorage.getItem("accessToken");
       if (token) {
           setCustomer({
               name: sessionStorage.getItem("userFullName") || localStorage.getItem("userFullName") || "Khách Hàng Web",
               phone_number: sessionStorage.getItem("phoneNumber") || localStorage.getItem("phoneNumber") || "0987654321",
               email: sessionStorage.getItem("username") || localStorage.getItem("username") || "customer.web@example.com",
           });
       }
    };
    
    const parseRouteStations = (fullRouteString) => {
        if (!fullRouteString) return { pickup: [], dropoff: [] };
        const stationRegex = /([\w\sÀ-ỹ]+)\s*\((\d+)\)/g;
        const stations = [];
        let match;
        while ((match = stationRegex.exec(fullRouteString)) !== null) {
            stations.push({ name: match[1].trim(), id: parseInt(match[2], 10) });
        }
        if (stations.length === 0) {
            const parts = fullRouteString.split(" → ");
            return {
                pickup: parts.slice(0, -1).map((name, i) => ({ name, id: i + 1 })),
                dropoff: parts.slice(1).map((name, i) => ({ name, id: i + 2 })),
            };
        }
        return { pickup: stations.slice(0, -1), dropoff: stations.slice(1) };
    };

    const fetchTripAndSeatDetails = async (tripID) => {
        const [tripDetailsResponse, availableSeatsResponse] = await Promise.all([
            api.get(`${API_URL}api/v1/trips/${tripID}/seats`),
            api.get(`${API_URL}api/v1/tickets-available/${tripID}`),
        ]);

        if (!tripDetailsResponse.data?.data) throw new Error(`Không tìm thấy chi tiết chuyến đi ${tripID}.`);
        if (!availableSeatsResponse.data?.data?.seats) throw new Error(`Không thể tải ghế trống cho chuyến ${tripID}.`);
        
        const tripData = tripDetailsResponse.data.data;
        const seatsData = availableSeatsResponse.data.data.seats;
        
        const availableSeatNames = seatsData.map((seat) => seat.name);
        const seatNameToIdMapping = seatsData.reduce((acc, seat) => {
            acc[seat.name] = seat.id;
            return acc;
        }, {});
        
        const { pickup, dropoff } = parseRouteStations(tripData.fullRoute);

        return {
            ...tripData,
            details: {
                availableSeats: availableSeatNames,
                seatMapping: seatNameToIdMapping,
                pickupLocations: pickup,
                dropoffLocations: dropoff,
            }
        };
    };

    const handleSeatSelect = (tripType, seatName, seatId) => {
        setBookingDetails(prev => {
            const currentTrip = prev[tripType];
            const isSelected = currentTrip.seats.includes(seatName);
            const newSeats = isSelected ? currentTrip.seats.filter(s => s !== seatName) : [...currentTrip.seats, seatName];
            const newSeatIds = isSelected ? currentTrip.seatIds.filter(id => id !== seatId) : [...currentTrip.seatIds, seatId];
            return { ...prev, [tripType]: { ...currentTrip, seats: newSeats, seatIds: newSeatIds }};
        });
    };
    
    const handleLocationChange = (tripType, locationType, locationId) => {
       setBookingDetails(prev => ({
           ...prev,
           [tripType]: { ...prev[tripType], [`${locationType}Id`]: locationId }
       }));
    };

    const handleCustomerChange = (e) => {
        setCustomer({ ...customer, [e.target.name]: e.target.value });
    };
    
    const resetBookingProcess = () => {
        setBookingState('IDLE');
        setBookingId(null);
        setIsSubmitting(false);
        // Tùy chọn: reset cả ghế đã chọn
        setBookingDetails({
            departure: { seats: [], seatIds: [], pickupId: '', dropoffId: '' },
            return: { seats: [], seatIds: [], pickupId: '', dropoffId: '' }
        });
    };

    // Hàm `initiateBooking` được cập nhật để theo luồng bất đồng bộ
    const initiateBooking = async () => {
        // --- VALIDATION ---
        if (bookingDetails.departure.seatIds.length === 0 || (returnTripId && bookingDetails.return.seatIds.length === 0)) {
            setMessage("Bạn chưa chọn đủ ghế cho tất cả các chặng."); openFailureModal(); return;
        }
        if (!bookingDetails.departure.pickupId || !bookingDetails.departure.dropoffId || (returnTripId && (!bookingDetails.return.pickupId || !bookingDetails.return.dropoffId))) {
            setMessage("Vui lòng chọn đủ điểm đón/trả cho các chặng."); openFailureModal(); return;
        }
        if (!termsAccepted) {
            setMessage("Vui lòng chấp nhận điều khoản và điều kiện."); openFailureModal(); return;
        }

        setIsSubmitting(true);

        // --- Dữ liệu gửi đi theo format của main.py ---
        const payload = {
            ticket_type: returnTripId ? 1 : 0, // 0: một chiều, 1: khứ hồi
            price: departureTripInfo.price, // Giá của một vé
            booking_channel: 2, // 2 for Web client
            policy_id: 1, // Chính sách mặc định
            name: customer.name,
            phone: customer.phone_number,
            email: customer.email,
            booked_by: "customer_web",
            trip_id_begin: departureTripId,
            seat_id_begin: bookingDetails.departure.seatIds.map(id => parseInt(id, 10)),
            pickup_location_begin: parseInt(bookingDetails.departure.pickupId),
            dropoff_location_begin: parseInt(bookingDetails.departure.dropoffId),
        };

        if (returnTripId) {
            payload.trip_id_end = returnTripId;
            payload.seat_id_end = bookingDetails.return.seatIds.map(id => parseInt(id, 10));
            payload.pickup_location_end = parseInt(bookingDetails.return.pickupId);
            payload.dropoff_location_end = parseInt(bookingDetails.return.dropoffId);
        }
        
        try {
            // 1. Gửi yêu cầu đến server để khởi tạo booking
            const response = await api.post(`${API_URL}api/v1/initiate-booking`, payload);

            // 2. Server trả về 202, nhận bookingId và chuyển sang màn hình chờ
            if (response.status === 202 && response.data?.data?.bookingId) {
                const id = response.data.data.bookingId;
                setBookingId(id);
                setBookingState('PROCESSING'); 
            } else {
                throw new Error(response.data?.message || "Server không phản hồi đúng định dạng.");
            }

        } catch (error) {
            const errorMessage = error.response?.data?.detail || error.message || "Lỗi khi khởi tạo đặt vé.";
            setMessage(errorMessage);
            openFailureModal();
            setIsSubmitting(false);
        }
    };
    
    const closeFailureModal = () => setFailureModal(false);
    const openFailureModal = () => setFailureModal(true);
    
    // --- RENDER LOGIC ---

    if (initialLoading) {
        return <div className="flex justify-center items-center h-screen"><p>Đang tải thông tin chuyến đi...</p></div>;
    }

    // Nếu không tải được thông tin chuyến đi chính
    if (!departureTripInfo) {
        return (
            <div className="flex flex-col justify-center items-center h-screen">
                <p className="text-red-500">Không thể tải thông tin chuyến đi hoặc chuyến đi không tồn tại.</p>
                {failureModal && ( <FailureNotification func={{ closeModal: closeFailureModal }} message={message} /> )}
            </div>
        );
    }
    
    // --- HIỂN THỊ MÀN HÌNH CHỜ VÀ KẾT QUẢ ---
    if (bookingState === 'PROCESSING') {
        return (
            <div className="w-full bg-slate-50 lg:p-10 flex justify-center items-center min-h-screen">
                <div className="max-w-screen-lg w-full">
                    <BookingStatus bookingId={bookingId} onReset={resetBookingProcess} />
                </div>
            </div>
        );
    }

    const totalCost = (departureTripInfo.price * bookingDetails.departure.seats.length) +
                      (returnTripInfo ? (returnTripInfo.price * bookingDetails.return.seats.length) : 0);

    const activeTripData = activeTripView === 'departure' ? departureTripInfo : returnTripInfo;
    
    // --- GIAO DIỆN FORM ĐẶT VÉ (Khi state là 'IDLE') ---
    return (
        <div className="w-full bg-slate-50 lg:p-10">
            <div className="flex max-w-screen-lg flex-col md:flex-row mx-auto gap-8">
                <div className="basis-2/3 w-full flex flex-col gap-5">

                    {/* --- TABS LƯỢT ĐI/VỀ --- */}
                    {returnTripInfo && (
                        <div className="flex bg-white p-1 rounded-lg border border-slate-200 shadow-sm w-full">
                            <button
                                onClick={() => setActiveTripView('departure')}
                                className={`w-1/2 p-3 text-center font-semibold rounded-md transition-colors duration-300 ${activeTripView === 'departure' ? 'bg-red-500 text-white' : 'bg-transparent text-gray-700 hover:bg-red-50'}`}
                            >
                                Thông tin lượt đi
                            </button>
                            <button
                                onClick={() => setActiveTripView('return')}
                                className={`w-1/2 p-3 text-center font-semibold rounded-md transition-colors duration-300 ${activeTripView === 'return' ? 'bg-red-500 text-white' : 'bg-transparent text-gray-700 hover:bg-red-50'}`}
                            >
                                Thông tin lượt về
                            </button>
                        </div>
                    )}
                    
                    {/* --- HIỂN THỊ THÔNG TIN CHUYẾN ĐI THEO TAB --- */}
                    {activeTripData && (
                        <TripBookingSection 
                            tripInfo={activeTripData}
                            tripType={activeTripView}
                            selectedSeats={bookingDetails[activeTripView].seats}
                            onSeatSelect={(seatName, seatId) => handleSeatSelect(activeTripView, seatName, seatId)}
                            onLocationChange={(type, id) => handleLocationChange(activeTripView, type, id)}
                            selectedPickupId={bookingDetails[activeTripView].pickupId}
                            selectedDropoffId={bookingDetails[activeTripView].dropoffId}
                        />
                    )}
                    
                    {/* --- THÔNG TIN KHÁCH HÀNG --- */}
                    <div className="user-info-section bg-white border border-slate-200 p-5 rounded-lg shadow-sm">
                        <div className="flex flex-col md:flex-row gap-8">
                            <div className="basis-full md:basis-1/2 w-full">
                                <h3 className="text-xl font-semibold">Thông tin khách hàng</h3>
                                <form action="" className="mt-5">
                                    <div className="mb-5">
                                        <p className="text-sm mb-1">Họ và tên <span className="text-red-500">*</span></p>
                                        <input type="text" name="name" className="rounded-xl border border-slate-300 w-full p-2 focus:border-indigo-500 focus:ring-indigo-500 shadow-sm" value={customer.name} onChange={handleCustomerChange} />
                                    </div>
                                    <div className="mb-5">
                                        <p className="text-sm mb-1">Số điện thoại <span className="text-red-500">*</span></p>
                                        <input type="text" name="phone_number" className="rounded-xl border border-slate-300 w-full p-2 focus:border-indigo-500 focus:ring-indigo-500 shadow-sm" value={customer.phone_number} onChange={handleCustomerChange} />
                                    </div>
                                    <div>
                                        <p className="text-sm mb-1">Email <span className="text-red-500">*</span></p>
                                        <input type="email" name="email" className="rounded-xl border border-slate-300 w-full p-2 focus:border-indigo-500 focus:ring-indigo-500 shadow-sm" value={customer.email} onChange={handleCustomerChange} />
                                    </div>
                                </form>
                            </div>
                            <div className="basis-full md:basis-1/2 w-full mt-8 md:mt-0">
                                <h3 className="text-red-500 text-center mb-5 font-semibold">ĐIỀU KHOẢN & LƯU Ý</h3>
                                <p className="text-[15px] text-justify mb-3 font-[500] leading-6">
                                    (*) Quý khách vui lòng có mặt tại bến xuất phát của xe trước ít nhất 30 phút giờ xe khởi hành, mang theo thông báo đã thanh toán vé thành công có chứa mã vé. Vui lòng liên hệ Trung tâm tổng đài <span className="text-red-500">1900 6067</span> để được hỗ trợ.
                                </p>
                            </div>
                        </div>
                        <div className="mt-5 text-sm font-[400]">
                            <input type="checkbox" id="termsAccepted" checked={termsAccepted} onChange={(e) => setTermsAccepted(e.target.checked)} className="mr-2 align-middle"/>
                            <label htmlFor="termsAccepted" className="cursor-pointer">
                                <span className="text-red-500 underline">Chấp nhận điều khoản</span> đặt vé & chính sách bảo mật thông tin của FUTABusline
                            </label>
                        </div>
                    </div>
                </div>

                {/* --- CHI TIẾT GIÁ --- */}
                <div className="basis-1/3 w-full flex flex-col gap-y-5">
                    <div className="bg-white border border-slate-200 p-5 rounded-xl shadow-sm">
                        <h3 className="text-xl font-medium">Chi tiết giá</h3>
                        <div className="mt-4 flex items-center justify-between">
                            <span className="text-slate-500">Giá vé lượt đi ({bookingDetails.departure.seats.length} ghế)</span>
                            <span className="text-red-500 font-medium">{(departureTripInfo.price * bookingDetails.departure.seats.length).toLocaleString()} đ</span>
                        </div>
                        {bookingDetails.departure.seats.length > 0 && (
                            <div className="text-right text-sm text-slate-500 mt-1 pr-1 font-semibold">
                                {bookingDetails.departure.seats.join(', ')}
                            </div>
                        )}

                        {returnTripInfo && (
                           <>
                             <div className="mt-2 flex items-center justify-between">
                                 <span className="text-slate-500">Giá vé lượt về ({bookingDetails.return.seats.length} ghế)</span>
                                 <span className="text-red-500 font-medium">{(returnTripInfo.price * bookingDetails.return.seats.length).toLocaleString()} đ</span>
                             </div>
                             {bookingDetails.return.seats.length > 0 && (
                                 <div className="text-right text-sm text-slate-500 mt-1 pr-1 font-semibold">
                                     {bookingDetails.return.seats.join(', ')}
                                 </div>
                             )}
                           </>
                        )}
                        <hr className="my-3" />
                        <div className="flex items-center justify-between">
                            <span className="text-slate-500 text-lg">Tổng tiền</span>
                            <span className="text-red-500 font-bold text-xl">{totalCost.toLocaleString()} đ</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* --- THANH TOÁN --- */}
            <div className="mt-5 flex max-w-screen-lg flex-col md:flex-row mx-auto gap-8 ">
                <div className="md:basis-2/3 w-full gap-5 payment-section bg-white border border-slate-200 p-5 rounded-b-xl flex items-center shadow-sm">
                    <div className="flex flex-col">
                        <span className="text-slate-500 text-sm">Tổng tiền</span>
                        <span className="mt-1 text-2xl font-medium text-black">{totalCost.toLocaleString()} đ</span>
                    </div>
                    <div className="flex flex-auto items-center justify-end">
                        <button
                            type="button"
                            className="px-5 py-3 text-white rounded-lg mr-6 bg-red-500 hover:bg-red-600 transition-all disabled:bg-gray-400 disabled:cursor-not-allowed"
                            onClick={initiateBooking}
                            disabled={isSubmitting || !termsAccepted || bookingDetails.departure.seatIds.length === 0}
                        >
                            {isSubmitting ? "Đang gửi yêu cầu..." : "Khởi tạo Đặt vé"}
                        </button>
                    </div>
                </div>
                <div className="basis-1/3 gap-5"></div>
            </div>

            {failureModal && ( <FailureNotification func={{ closeModal: closeFailureModal }} message={message} /> )}
        </div>
    );
}

export default BookingTicketForm;