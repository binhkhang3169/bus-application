/** @format */

import "./App.css";
import {
  HashRouter,
  Outlet,
  Route,
  Routes,
  Navigate,
} from "react-router-dom";

// Layouts
import DashboardLayout from "./components/admin/layout/DashboardLayout";
import Navbar from "./components/Navbar"; // Navbar cho khách hàng
import Footer from "./components/Footer"; // Footer cho khách hàng

// --- Customer Pages ---
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import LookUpPage from "./pages/LookUpPage";
import TravelSchedulePage from "./pages/TravelSchedulePage";
import BookingTicketPage from "./pages/BookingTicketPage";
import ResultBookingPage from "./pages/ResultBookingPage";
import SignupPage from "./pages/SignupPage";
import EmailVerificationPage from "./pages/EmailVerificationPage";
import ProfileManagementPage from "./pages/ProfileManagementPage";
import ForgotPasswordPage from "./pages/ForgotPasswordPage";
import AiChatbox from "./components/AiChatbox";
import VerifyForgotPasswordPage from "./pages/VerifyForgotPassword";
import NewsPage from "./pages/NewsPage";
import NewsDetailPage from "./pages/NewsDetailPage";

// --- Admin Pages ---
import LoginPageAdmin from "./pages/admin/LoginPageAdmin";
import DashboardPage from "./pages/DashboardPage";
import TripManagement from "./components/admin/TripManagement";
import TicketManagement from "./components/admin/TicketManagement";
import BusManagerment from "./components/admin/BusManagement";
import BusStationManagement from "./components/admin/BusStationManagement";
import RouteManagement from "./components/admin/RouteManagement";
import EmployeeManagement from "./components/admin/EmployeeManagement";
import CustomerManagement from "./components/admin/CustomerManagement";
import AccountManagement from "./components/admin/AccountManagement";
import TripDriver from "./components/admin/TripDriver";
import PickupManagement from "./components/admin/PickupManagement";
import NewsManagement from "./components/admin/NewsManagement";
import ProvinceManagement from "./components/admin/ProvinceManagement";
import BusBin from "./components/admin/bin/BusBin";
import BusStationBin from "./components/admin/bin/BusStationBin";
import RouteBin from "./components/admin/bin/RouteBin";
import TripBin from "./components/admin/bin/TripBin";
import ShipmentManagement from "./components/admin/ShipmentManagement";
import InvoicePage from "./pages/admin/InvoicePage";
import AdminBookingPage from "./components/admin/AdminBookingPage";

// Component Layout cho các trang của khách hàng
const CustomerLayout = () => {
  return (
    <>
      <Navbar />
      <Outlet />
      <Footer />
      <AiChatbox />
    </>
  );
};

function App() {
  return (
    <HashRouter>
      <Routes>
        {/* === Customer Routes: Các route dành cho khách hàng === */}
        <Route path="/" element={<CustomerLayout />}>
          <Route index element={<HomePage />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="forgot-password" element={<ForgotPasswordPage />} />
          <Route
            path="forgot-password/verify"
            element={<VerifyForgotPasswordPage />}
          />
          <Route path="signup" element={<SignupPage />} />
          <Route path="email-verify" element={<EmailVerificationPage />} />
          <Route path="lookup-ticket" element={<LookUpPage />} />
          <Route path="search-trip" element={<TravelSchedulePage />} />
          <Route path="dat-ve" element={<BookingTicketPage />} />
          <Route path="ket-qua-dat-ve" element={<ResultBookingPage />} />
          <Route path="account" element={<ProfileManagementPage />} />
          <Route path="news" element={<NewsPage />} />
          <Route path="news/:id" element={<NewsDetailPage />} />
        </Route>

        {/* === Admin Routes: Các route dành cho nhân viên/quản trị === */}

        {/* 1. Route Đăng nhập: Nằm riêng biệt, không dùng chung layout */}
        <Route path="/admin" element={<LoginPageAdmin />} />

        {/* 2. Protected Routes: Tất cả các trang admin khác được bao bọc bởi DashboardLayout */}
        <Route element={<DashboardLayout />}>
          <Route path="/admin/home" element={<DashboardPage />} />
          <Route path="/admin/trip" element={<TripManagement />} />
          <Route path="/admin/trip/bin" element={<TripBin />} />
          <Route path="/admin/route" element={<RouteManagement />} />
          <Route path="/admin/route/bin" elemsent={<RouteBin />} />
          <Route path="/admin/pickup" element={<PickupManagement />} />
          <Route path="/admin/ticket" element={<TicketManagement />} />
          <Route path="/admin/bus" element={<BusManagerment />} />
          <Route path="/admin/booking" element={<AdminBookingPage />} /> 
          <Route path="/admin/bus/bin" element={<BusBin />} />
          <Route path="/admin/bus-station" element={<BusStationManagement />} />
          <Route path="/admin/bus-station/bin" element={<BusStationBin />} />
          <Route path="/admin/province" element={<ProvinceManagement />} />
          <Route path="/admin/employee" element={<EmployeeManagement />} />
          <Route path="/admin/customer" element={<CustomerManagement />} />
          <Route path="/admin/news" element={<NewsManagement />} />
          <Route path="/admin/shipment" element={<ShipmentManagement />} />
          <Route
            path="/admin/shipment/:shipmentId/invoice"
            element={<InvoicePage />}
          />
          <Route path="/admin/my-trip" element={<TripDriver />} />
          <Route path="/admin/account" element={<AccountManagement />} />
          {/* Lưu ý: Component ProtectRoute của bạn vẫn có thể được sử dụng ở đây
              để kiểm tra quyền hạn chi tiết (ví dụ: chỉ admin mới thấy trang Employee)
              <Route path="/admin/employee" element={<ProtectRoute allowedRoles={['ROLE_ADMIN']}><EmployeeManagement /></ProtectRoute>} />
          */}
        </Route>

        {/* Fallback Route: Chuyển hướng các đường dẫn không hợp lệ */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </HashRouter>
  );
}

export default App;
