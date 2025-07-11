/** @format */

import React from 'react';

// =================================================================
// TripBookingSection Component
// =================================================================
const TripBookingSection = ({ 
    tripInfo, 
    tripType, 
    onSeatSelect, 
    selectedSeats, 
    onLocationChange,
    selectedPickupId,
    selectedDropoffId 
}) => {
    if (!tripInfo) return null;

    const { availableSeats, seatMapping, pickupLocations, dropoffLocations } = tripInfo.details;
    const title = tripType === 'departure' ? 'Thông tin lượt đi' : 'Thông tin lượt về';

    const isSeatAvailable = (seat) => availableSeats.includes(seat);
    const isSeatSelected = (seat) => selectedSeats.includes(seat);
    
    const getSeatImage = (seat) => {
        if (isSeatSelected(seat)) return "https://futabus.vn/images/icons/seat_selecting.svg";
        if (isSeatAvailable(seat)) return "https://futabus.vn/images/icons/seat_active.svg";
        return "https://futabus.vn/images/icons/seat_disabled.svg";
    };

    const getSeatStatusClass = (seat) => {
      if (isSeatSelected(seat)) return "text-red-400";
      if (isSeatAvailable(seat)) return "text-blue-400";
      return "text-gray-400";
    };
    
    const formatSeatDisplay = (seat) => {
        if (!seat || seat.length < 2) return seat;
        const row = seat.charAt(0);
        const num = parseInt(seat.substring(1));
        return `${row}${num < 10 ? "0" : ""}${num}`;
    };
    
    // Giả định xe có 2 dãy A, B và mỗi dãy 15 ghế
    const generateAllSeats = (row) => Array.from({ length: 15 }, (_, i) => `${row}${i + 1}`);

    return (
        <div className="bg-white rounded-lg p-6 border border-slate-200 shadow-sm">
            <h3 className="text-xl font-medium">{title}</h3>
            <p className="mt-2 text-gray-600">{tripInfo.departureStation} → {tripInfo.arrivalStation} | {tripInfo.departureDate}</p>
            <div className="flex flex-col mt-5 gap-8">
                <div className="flex flex-row gap-8 justify-center">
                    {['A', 'B'].map(row => (
                        <div key={row} className="grid grid-cols-3 gap-x-10 gap-y-2">
                            {generateAllSeats(row).map((seat) => (
                                <div key={seat}
                                    className={`mt-1 text-center relative flex justify-center ${!isSeatAvailable(seat) ? "cursor-not-allowed" : "cursor-pointer"}`}
                                    onClick={() => isSeatAvailable(seat) && onSeatSelect(seat, seatMapping[seat])}>
                                    <img width="32" src={getSeatImage(seat)} alt="seat icon" />
                                    <span className={`absolute text-sm font-semibold sm:text-[12px] ${getSeatStatusClass(seat)} top-1`}>{formatSeatDisplay(seat)}</span>
                                </div>
                            ))}
                        </div>
                    ))}
                </div>
                <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Điểm đón</label>
                        <select 
                            value={selectedPickupId}
                            onChange={(e) => onLocationChange('pickup', e.target.value)} 
                            className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
                            <option value="">Chọn điểm đón</option>
                            {pickupLocations.map((loc) => <option key={loc.id} value={loc.id}>{loc.name}</option>)}
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Điểm trả</label>
                        <select 
                            value={selectedDropoffId}
                            onChange={(e) => onLocationChange('dropoff', e.target.value)} 
                            className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
                            <option value="">Chọn điểm trả</option>
                            {dropoffLocations.map((loc) => <option key={loc.id} value={loc.id}>{loc.name}</option>)}
                        </select>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default TripBookingSection;