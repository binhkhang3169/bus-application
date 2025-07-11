import React from "react";
import { NavLink } from "react-router-dom";

const navLinksConfig = {
  ROLE_ADMIN: [
    { to: "/admin/trip", text: "Quản lý chuyến đi" },
    { to: "/admin/route", text: "Quản lý lộ trình" },
    { to: "/admin/pickup", text: "Quản lý các điểm đón - trả" },
    { to: "/admin/bus-station", text: "Quản lý các trạm xe" },
    { to: "/admin/province", text: "Quản lý các tỉnh" },
    { to: "/admin/bus", text: "Quản lý xe" },
    { to: "/admin/ticket", text: "Quản lý vé xe" },
    { to: "/admin/employee", text: "Quản lý nhân viên" },
    { to: "/admin/news", text: "Quản lý tin tức" },
    { to: "/admin/shipment", text: "Hàng hóa" },
    { to: "/admin/booking", text: "Đặt vé" },
    { to: "/admin/customer", text: "Quản lý khách hàng" },
  ],
  ROLE_OPERATOR: [
    { to: "/admin/trip", text: "Quản lý chuyến đi" },
    { to: "/admin/route", text: "Quản lý lộ trình" },
    { to: "/admin/pickup", text: "Quản lý các điểm đón - trả" },
    { to: "/admin/bus-station", text: "Quản lý các trạm xe" },
    { to: "/admin/bus", text: "Quản lý xe" },
    { to: "/admin/booking", text: "Đặt vé" },
    { to: "/admin/shipment", text: "Hàng hóa" },
  ],
  ROLE_RECEPTION: [
    { to: "/admin/ticket", text: "Quản lý vé xe" },
    { to: "/admin/shipment", text: "Hàng hóa" },
    { to: "/admin/booking", text: "Đặt vé" },
    { to: "/admin/customer", text: "Quản lý khách hàng" },
  ],
  ROLE_DRIVER: [{ to: "/admin/my-trip", text: "My Trip" }],
};

const Sidebar = ({ userRole }) => {
  // onLogout không còn cần thiết ở đây
  const roleNameMap = {
    ROLE_ADMIN: "Quản trị viên",
    ROLE_OPERATOR: "Điều phối viên",
    ROLE_RECEPTION: "Lễ tân",
    ROLE_DRIVER: "Tài xế",
  };

  const navLinks = navLinksConfig[userRole] || [];
  const roleName = roleNameMap[userRole] || "Admin";

  const baseLinkClasses =
    "relative flex items-center rounded-lg py-3 px-4 font-medium duration-300 ease-in-out hover:bg-gray-100 dark:hover:bg-gray-800";
  const activeLinkClasses = "bg-blue-500 text-white dark:bg-blue-600";
  const inactiveLinkClasses = "text-gray-800 dark:text-gray-200";

  return (
    // ClassName được đơn giản hóa, vì việc ẩn hiện đã do DashboardLayout quản lý
    <aside className="static flex h-screen w-72 flex-col overflow-y-hidden bg-white duration-300 ease-linear ">
      {/* <div className="flex items-center justify-between px-4 py-4 ">
        <div className="px-4 py-4  ">
          <span className="text-xl font-semibold text-gray-800 dark:text-white">
            {roleName}
          </span>
        </div>
      </div> */}

      <div className="no-scrollbar flex flex-col overflow-y-auto duration-300 ease-linear">
        <nav className=" py-4 px-4 lg:px-6">
          <div>
            {/* <h3 className="mb-4 ml-4 text-sm font-semibold text-gray-500 dark:text-gray-400">
              MENU
            </h3> */}
            <ul className="mb-6 flex flex-col gap-1.5">
              <li>
                <NavLink
                  to="/admin/home"
                  end
                  className={({ isActive }) =>
                    `${baseLinkClasses} ${
                      isActive ? activeLinkClasses : inactiveLinkClasses
                    }`
                  }
                >
                  {roleName}
                </NavLink>
              </li>
              {navLinks.map((link) => (
                <li key={link.to}>
                  <NavLink
                    to={link.to}
                    className={({ isActive }) =>
                      `${baseLinkClasses} ${
                        isActive ? activeLinkClasses : inactiveLinkClasses
                      }`
                    }
                  >
                    {link.text}
                  </NavLink>
                </li>
              ))}
            </ul>
          </div>
        </nav>
      </div>
    </aside>
  );
};

export default Sidebar;
