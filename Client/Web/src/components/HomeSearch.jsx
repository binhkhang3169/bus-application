/** @format */
import React, { useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import FailureNotification from "./Noti/FailureNotification";
import api from "../services/apiService";
import { API_URL } from "../configs/env";
import dayjs from "dayjs";
import isSameOrAfter from "dayjs/plugin/isSameOrAfter"; // SỬA LỖI: Import plugin

dayjs.extend(isSameOrAfter); // SỬA LỖI: Kích hoạt plugin

function HomeSearch() {
  const navigate = useNavigate();
  const location = useLocation();

  const [locations, setLocations] = useState([]);
  const [failureModal, setFailureModal] = useState(false);
  const [startAddressSearchModal, setStartAddressSearchModal] = useState(false);
  const [endAddressSearchModal, setEndAddressSearchModal] = useState(false);
  const [message, setMessage] = useState("");
  const [startAddress, setStartAddress] = useState("");
  const [endAddress, setEndAddress] = useState("");
  const [startAddressId, setStartAddressId] = useState("");
  const [endAddressId, setEndAddressId] = useState("");
  const [startAddressSearch, setStartAddressSearch] = useState("");
  const [endAddressSearch, setEndAddressSearch] = useState("");
  const [date, setDate] = useState(dayjs().format("YYYY-MM-DD"));
  const [history, setHistory] = useState([]);
  const [isRoundTrip, setIsRoundTrip] = useState(false); // Mặc định là khứ hồi
  const [returnDate, setReturnDate] = useState(
    dayjs().add(1, "day").format("YYYY-MM-DD")
  );

  useEffect(() => {
    getLocations();
    parseUrlParams();
    loadHistory();
  }, []);

  // Thêm useEffect để đóng modal khi click ra ngoài
  useEffect(() => {
    const handleClickOutside = (event) => {
      // Kiểm tra xem click có nằm ngoài các phần tử modal không
      if (
        startAddressSearchModal &&
        !event.target.closest(".start-address-modal-container")
      ) {
        setStartAddressSearchModal(false);
      }
      if (
        endAddressSearchModal &&
        !event.target.closest(".end-address-modal-container")
      ) {
        setEndAddressSearchModal(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [startAddressSearchModal, endAddressSearchModal]);

  const getLocations = async () => {
    try {
      const res = await api.get(`${API_URL}api/v1/provinces`);
      if (res.data && res.data.code === 200) {
        const activeLocations = res.data.data.filter((loc) => loc.status === 1);
        setLocations(activeLocations);
      } else {
        throw new Error("Failed to fetch locations from API");
      }
    } catch (err) {
      console.error("Failed to load locations from API:", err);
      setMessage("Không thể tải danh sách địa điểm. Vui lòng thử lại.");
      setFailureModal(true);
    }
  };

  const parseUrlParams = () => {
    const query = new URLSearchParams(location.search);
    const start = query.get("from") || "";
    const startId = query.get("fromId") || "";
    const end = query.get("to") || "";
    const endId = query.get("toId") || "";
    const dateParam = query.get("fromTime") || dayjs().format("YYYY-MM-DD");
    const returnDateParam = query.get("returnTime");

    setStartAddress(decodeURIComponent(start).trim());
    setStartAddressId(startId);
    setEndAddress(decodeURIComponent(end).trim());
    setEndAddressId(endId);
    setDate(dateParam);

    if (returnDateParam) {
      setIsRoundTrip(true);
      setReturnDate(returnDateParam);
    }
  };

  const saveToHistory = (start, startId, end, endId) => {
    let stored = JSON.parse(localStorage.getItem("trip_history")) || [];
    const pair = { start, startId, end, endId };
    stored = stored.filter((p) => !(p.start === start && p.end === end));
    stored.unshift(pair);
    if (stored.length > 3) stored = stored.slice(0, 3);
    localStorage.setItem("trip_history", JSON.stringify(stored));
    setHistory(stored);
  };

  const loadHistory = () => {
    const stored = JSON.parse(localStorage.getItem("trip_history")) || [];
    setHistory(stored);
  };

  const deleteHistoryItem = (index) => {
    const newHistory = [...history];
    newHistory.splice(index, 1);
    setHistory(newHistory);
    localStorage.setItem("trip_history", JSON.stringify(newHistory));
  };

  const handleSearch = () => {
    if (startAddress && startAddressId && endAddress && endAddressId && date) {
      if (isRoundTrip && !returnDate) {
        setMessage("Vui lòng chọn ngày về cho chuyến khứ hồi!");
        openFailureModal();
        return;
      }
      if (isRoundTrip && dayjs(returnDate).isBefore(dayjs(date))) {
        setMessage("Ngày về phải sau ngày đi!");
        openFailureModal();
        return;
      }

      saveToHistory(startAddress, startAddressId, endAddress, endAddressId);

      let searchString = `/search-trip?from=${encodeURIComponent(
        startAddress
      )}&fromId=${startAddressId}&fromTime=${date}&to=${encodeURIComponent(
        endAddress
      )}&toId=${endAddressId}`;
      if (isRoundTrip) {
        searchString += `&returnTime=${returnDate}`;
      }
      navigate(searchString);
    } else {
      setMessage("Vui lòng điền đầy đủ thông tin!");
      openFailureModal();
    }
  };

  const swapLocations = () => {
    const tempAddress = startAddress;
    const tempAddressId = startAddressId;

    setStartAddress(endAddress);
    setStartAddressId(endAddressId);

    setEndAddress(tempAddress);
    setEndAddressId(tempAddressId);
  };

  const openFailureModal = () => setFailureModal(true);
  const closeFailureModal = () => setFailureModal(false);

  const handleDateChange = (e) => {
    const newDate = dayjs(e.target.value).format("YYYY-MM-DD");
    setDate(newDate);
    // Tự động cập nhật ngày về nếu ngày đi mới sau ngày về hiện tại
    if (isRoundTrip && dayjs(newDate).isSameOrAfter(dayjs(returnDate))) {
      setReturnDate(dayjs(newDate).add(1, "day").format("YYYY-MM-DD"));
    }
  };

  const handleReturnDateChange = (e) => {
    const rawValue = e.target.value;
    const formatted = dayjs(rawValue).format("YYYY-MM-DD");
    setReturnDate(formatted);
  };

  const handleLocationSelection = (item, isStart) => {
    if (isStart) {
      setStartAddress(item.name);
      setStartAddressId(item.id);
      setStartAddressSearchModal(false);
    } else {
      setEndAddress(item.name);
      setEndAddressId(item.id);
      setEndAddressSearchModal(false);
    }
  };

  const filterLocations = (search, items) => {
    if (!search) return items;
    return items.filter((item) =>
      item.name.toLowerCase().includes(search.toLowerCase())
    );
  };

  return (
    <>
      <div className="mt-10 mb-20">
        <h3 className="text-center text-2xl font-bold text-blue-700 mb-5">
          LỰA CHỌN CHUYẾN ĐI
        </h3>

        <hr className="w-4/5 mx-auto h-0.5 bg-gray-200 md:w-2/5 xl:w-1/5" />
        <div className="max-w-screen-lg bg-white text-black mx-auto mt-10 p-5 border border-orange-400 shadow-2xl rounded-lg flex flex-col gap-4">
          <div className="flex items-center justify-end gap-4">
            <label className="font-semibold">Một chiều</label>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={isRoundTrip}
                onChange={() => setIsRoundTrip(!isRoundTrip)}
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 rounded-full peer peer-focus:ring-4 peer-focus:ring-green-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-green-600"></div>
            </label>
            <label className="font-semibold">Khứ hồi</label>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-4 items-end gap-4">
            {/* Điểm đi */}
            <div className="input-search flex flex-col w-full">
              <label className="mb-3 font-semibold">Điểm đi</label>
              <div className="rounded-lg relative h-10 border border-slate-300 cursor-pointer start-address-modal-container">
                <div
                  className="w-full h-full flex items-center"
                  onClick={() => setStartAddressSearchModal(true)}
                >
                  <span className="ml-2 truncate">
                    {startAddress || "Chọn điểm đi"}
                  </span>
                </div>
                {startAddressSearchModal && (
                  <div className="absolute top-full mt-1 left-0 w-full bg-white rounded-xl shadow-lg z-50 border">
                    <input
                      value={startAddressSearch}
                      className="h-10 w-full rounded-t-lg border-b border-slate-300 px-2 focus:outline-none"
                      type="text"
                      placeholder="Tìm điểm đi..."
                      onChange={(e) => setStartAddressSearch(e.target.value)}
                      autoFocus
                    />
                    <ul className="w-full overflow-y-auto max-h-52">
                      {filterLocations(startAddressSearch, locations).map(
                        (item) => (
                          <li
                            key={item.id}
                            className="py-2 px-4 hover:bg-gray-100 cursor-pointer"
                            onClick={() => handleLocationSelection(item, true)}
                          >
                            {item.name}
                          </li>
                        )
                      )}
                    </ul>
                  </div>
                )}
              </div>
            </div>

            {/* Nút Đảo chiều */}
            <button
              onClick={swapLocations}
              className="hidden md:flex items-center justify-center h-10 w-10 mx-auto bg-gray-200 rounded-full hover:bg-gray-300 transition-colors"
              title="Đảo chiều"
            >
              <svg
                className="w-6 h-6 text-gray-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                ></path>
              </svg>
            </button>

            {/* Điểm đến */}
            <div className="input-search flex flex-col w-full md:col-start-3">
              <label className="mb-3 font-semibold">Điểm đến</label>
              <div className="rounded-lg relative h-10 border border-slate-300 cursor-pointer end-address-modal-container">
                <div
                  className="w-full h-full flex items-center"
                  onClick={() => setEndAddressSearchModal(true)}
                >
                  <span className="ml-2 truncate">
                    {endAddress || "Chọn điểm đến"}
                  </span>
                </div>
                {endAddressSearchModal && (
                  <div className="absolute top-full mt-1 left-0 w-full bg-white rounded-xl shadow-lg z-50 border">
                    <input
                      value={endAddressSearch}
                      className="h-10 w-full rounded-t-lg border-b border-slate-300 px-2 focus:outline-none"
                      type="text"
                      placeholder="Tìm điểm đến..."
                      onChange={(e) => setEndAddressSearch(e.target.value)}
                      autoFocus
                    />
                    <ul className="w-full overflow-y-auto max-h-52">
                      {filterLocations(endAddressSearch, locations).map(
                        (item) => (
                          <li
                            key={item.id}
                            className="py-2 px-4 hover:bg-gray-100 cursor-pointer"
                            onClick={() => handleLocationSelection(item, false)}
                          >
                            {item.name}
                          </li>
                        )
                      )}
                    </ul>
                  </div>
                )}
              </div>
            </div>

            {/* Ngày đi */}
            <div className="input-search flex flex-col w-full md:col-span-2">
              <label className="mb-3 font-semibold">Ngày đi</label>
              <input
                value={date}
                className="rounded-lg border-slate-300 h-10 px-2 border focus:outline-orange-400"
                type="date"
                min={dayjs().format("YYYY-MM-DD")}
                onChange={handleDateChange}
              />
            </div>

            {/* --- CẬP NHẬT: Ô NGÀY VỀ LUÔN HIỂN THỊ --- */}
            <div
              className={`input-search flex flex-col w-full md:col-span-2 transition-opacity duration-300 ${
                !isRoundTrip ? "opacity-50" : "opacity-100"
              }`}
            >
              <label className="mb-3 font-semibold">Ngày về</label>
              <input
                value={returnDate}
                className={`rounded-lg border-slate-300 h-10 px-2 border focus:outline-orange-400 ${
                  !isRoundTrip ? "bg-gray-100 cursor-not-allowed" : ""
                }`}
                type="date"
                min={dayjs(date).add(0, "day").format("YYYY-MM-DD")}
                onChange={handleReturnDateChange}
                disabled={!isRoundTrip} // Thuộc tính disabled được điều khiển bởi state isRoundTrip
              />
            </div>

            {/* Nút tìm kiếm */}
            <div className="search-button flex items-center md:col-span-4 justify-center mt-4">
              <button
                className="font-semibold text-white hover:bg-blue-600 transition-all w-full md:w-auto border border-transparent bg-blue-500 px-10 py-3 rounded-full"
                onClick={handleSearch}
              >
                Tìm chuyến xe
              </button>
            </div>
          </div>

          {history.length > 0 && (
            <div className="text-left mt-4 border-t pt-4">
              <h4 className="font-semibold mb-2 text-gray-600">
                Các tuyến đã tìm gần đây:
              </h4>
              <div className="flex justify-start gap-2 flex-wrap">
                {history.map((item, index) => (
                  <div
                    key={index}
                    className="flex items-center bg-gray-100 border border-gray-300 rounded-lg pl-3 pr-2 py-2 text-sm"
                  >
                    <button
                      className="hover:text-orange-500"
                      onClick={() => {
                        setStartAddress(item.start);
                        setStartAddressId(item.startId);
                        setEndAddress(item.end);
                        setEndAddressId(item.endId);
                      }}
                    >
                      {item.start} → {item.end}
                    </button>
                    <button
                      className="ml-2 text-gray-400 font-bold hover:text-red-500"
                      onClick={() => deleteHistoryItem(index)}
                      title="Xóa"
                    >
                      <svg
                        className="w-4 h-4"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                          clipRule="evenodd"
                        ></path>
                      </svg>
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
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

export default HomeSearch;
