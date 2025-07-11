-- phpMyAdmin SQL Dump
-- version 5.2.2
-- https://www.phpmyadmin.net/
--
-- Host: mysql
-- Generation Time: May 25, 2025 at 04:17 AM
-- Server version: 9.2.0
-- PHP Version: 8.2.27

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `trip_service`
--

-- --------------------------------------------------------

--
-- Table structure for table `log_type`
--

CREATE TABLE `log_type` (
  `id` int NOT NULL,
  `type_id` int DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `updated_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `log_type`
--

INSERT INTO `log_type` (`id`, `type_id`, `updated_at`, `updated_by`) VALUES
(1, 4, '2025-05-25 10:42:05.313000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `log_vehicle`
--

CREATE TABLE `log_vehicle` (
  `id` int NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `updated_by` int DEFAULT NULL,
  `vehicle_id` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `log_vehicle`
--

INSERT INTO `log_vehicle` (`id`, `updated_at`, `updated_by`, `vehicle_id`) VALUES
(1, '2025-05-25 10:29:44.112000', 1, 1),
(2, '2025-05-25 10:31:15.429000', 1, 1),
(3, '2025-05-25 10:31:32.643000', 1, 1),
(4, '2025-05-25 10:35:17.386000', 1, 2),
(5, '2025-05-25 10:35:38.295000', 1, 2);

-- --------------------------------------------------------

--
-- Table structure for table `pickup`
--

CREATE TABLE `pickup` (
  `id` varchar(255) NOT NULL,
  `self_id` varchar(255) DEFAULT NULL,
  `route_id` int DEFAULT NULL,
  `station_id` int DEFAULT NULL,
  `path_id` int NOT NULL,
  `time` varchar(255) DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `pickup`
--

INSERT INTO `pickup` (`id`, `self_id`, `route_id`, `station_id`, `path_id`, `time`, `status`, `created_at`, `created_by`) VALUES
('P1', '-1', 1, 1, 1, '0', 1, '2025-05-25 02:59:11.000000', 1),
('P10', '-1', 4, 3, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P11', 'P10', 4, 6, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P12', '-2', 4, 1, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P13', '-1', 5, 5, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P14', '-2', 5, 7, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P15', 'P13', 5, 2, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P16', '-1', 1, 2, 2, '60', 1, '2025-05-25 02:59:11.000000', 1),
('P17', '-2', 1, 4, 2, NULL, 1, '2025-05-25 02:59:11.000000', 1),
('P18', 'P2', 1, 10, 1, '200', 1, '2025-05-25 02:59:11.000000', 1),
('P19', 'P16', 1, 12, 2, '120', 1, '2025-05-25 02:59:11.000000', 1),
('P2', 'P1', 1, 7, 1, '0', 1, '2025-05-25 02:59:11.000000', 1),
('P3', '-2', 1, 3, 1, NULL, 1, '2025-05-25 02:59:11.000000', 1),
('P4', '-1', 2, 1, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P5', 'P4', 2, 6, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P6', 'P5', 2, 10, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P7', '-2', 2, 5, 0, '0', 1, '2025-05-25 02:59:11.000000', 1),
('P8', '-1', 3, 3, 0, '30', 1, '2025-05-25 02:59:11.000000', 1),
('P9', '-2', 3, 1, 0, '30', 1, '2025-05-25 02:59:11.000000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `province`
--

CREATE TABLE `province` (
  `id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `province`
--

INSERT INTO `province` (`id`, `name`, `status`, `created_at`, `created_by`) VALUES
(1, 'Hà Nội', 0, '2025-05-25 02:59:11.000000', 1),
(2, 'Hồ Chí Minh', 1, '2025-05-25 02:59:11.000000', 1),
(3, 'Đà Nẵng', 1, '2025-05-25 02:59:11.000000', 1),
(4, 'Hải Phòng', 1, '2025-05-25 02:59:11.000000', 1),
(5, 'Nghệ An', 1, '2025-05-25 02:59:11.000000', 1),
(6, 'Quảng Ninh', 1, '2025-05-25 02:59:11.000000', 1),
(7, 'Huế', 1, '2025-05-25 02:59:11.000000', 1),
(8, 'Cần Thơ', 1, '2025-05-25 02:59:11.000000', 1),
(9, 'Lâm Đồng', 1, '2025-05-25 02:59:11.000000', 1),
(10, 'Bình Dương', 1, '2025-05-25 02:59:11.000000', 1),
(11, 'Bình Thuận', 1, '2025-05-25 02:59:11.000000', 1),
(12, 'Quảng Nam', 1, '2025-05-25 02:59:11.000000', 1),
(13, 'Nam Định', 1, '2025-05-25 02:59:11.000000', 1),
(14, 'Thái Bình', 1, '2025-05-25 02:59:11.000000', 1),
(15, 'Thanh Hóa', 1, '2025-05-25 02:59:11.000000', 1),
(16, 'Quảng Trị', 1, '2025-05-25 02:59:11.000000', 1),
(17, 'Ninh Bình', 1, '2025-05-25 02:59:11.000000', 1),
(18, 'Hòa Bình', 1, '2025-05-25 02:59:11.000000', 1),
(19, 'Vĩnh Phúc', 1, '2025-05-25 02:59:11.000000', 1),
(20, 'Phú Thọ', 1, '2025-05-25 02:59:11.000000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `route`
--

CREATE TABLE `route` (
  `id` int NOT NULL,
  `distance` varchar(255) DEFAULT NULL,
  `estimated_time` varchar(255) DEFAULT NULL,
  `price` int DEFAULT NULL,
  `end` int DEFAULT NULL,
  `start` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `route`
--

INSERT INTO `route` (`id`, `distance`, `estimated_time`, `price`, `end`, `start`, `status`, `created_at`, `created_by`) VALUES
(1, '1700 km', '28h', 900000, 2, 1, 1, '2025-05-25 02:59:11.000000', 1),
(2, '950 km', '17h', 700000, 3, 1, 1, '2025-05-25 02:59:11.000000', 1),
(3, '1700 km', '28h', 900000, 1, 2, 1, '2025-05-25 02:59:11.000000', 1),
(4, '940 km', '17h', 690000, 3, 2, 1, '2025-05-25 02:59:11.000000', 1),
(5, '750 km', '14h', 600000, 5, 3, 1, '2025-05-25 02:59:11.000000', 1),
(6, '750 km', '14h', 600000, 3, 5, 1, '2025-05-25 02:59:11.000000', 1),
(7, '120 km', '3h', 200000, 6, 1, 1, '2025-05-25 02:59:11.000000', 1),
(8, '120 km', '3h', 200000, 1, 6, 1, '2025-05-25 02:59:11.000000', 1),
(9, '300 km', '6h', 350000, 7, 1, 1, '2025-05-25 02:59:11.000000', 1),
(10, '300 km', '6h', 350000, 1, 7, 1, '2025-05-25 02:59:11.000000', 1),
(11, '300 km', '7h', 400000, 12, 2, 1, '2025-05-25 02:59:11.000000', 1),
(12, '300 km', '7h', 400000, 2, 12, 1, '2025-05-25 02:59:11.000000', 1),
(13, '20 km', '1h', 50000, 4, 3, 1, '2025-05-25 02:59:11.000000', 1),
(14, '20 km', '1h', 50000, 3, 4, 1, '2025-05-25 02:59:11.000000', 1),
(15, '700 km', '13h', 550000, 10, 2, 1, '2025-05-25 02:59:11.000000', 1),
(16, '700 km', '13h', 550000, 2, 10, 1, '2025-05-25 02:59:11.000000', 1),
(17, '170 km', '4h', 220000, 11, 3, 1, '2025-05-25 02:59:11.000000', 1),
(18, '170 km', '4h', 220000, 3, 11, 1, '2025-05-25 02:59:11.000000', 1),
(19, '20 km', '30m', 30000, 13, 3, 1, '2025-05-25 02:59:11.000000', 1),
(20, '20 km', '30m', 30000, 3, 13, 1, '2025-05-25 02:59:11.000000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `special_day`
--

CREATE TABLE `special_day` (
  `id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `percent` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `special_day`
--

INSERT INTO `special_day` (`id`, `name`, `percent`, `status`, `created_at`, `created_by`) VALUES
(1, 'Tết Dương Lịch', 20, 1, '2025-05-25 02:59:11.000000', 1),
(2, 'Tết Nguyên Đán', 50, 1, '2025-05-25 02:59:11.000000', 1),
(3, 'Giỗ Tổ Hùng Vương', 15, 1, '2025-05-25 02:59:11.000000', 1),
(4, '30/4 - 1/5', 25, 1, '2025-05-25 02:59:11.000000', 1),
(5, '2/9', 25, 1, '2025-05-25 02:59:11.000000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `station`
--

CREATE TABLE `station` (
  `id` int NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `province_id` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `station`
--

INSERT INTO `station` (`id`, `address`, `name`, `province_id`, `status`, `created_at`, `created_by`) VALUES
(1, 'Nam Từ Liêm, Hà Nội', 'BX Mỹ Đình', 1, 1, '2025-05-25 02:59:10.000000', 1),
(2, 'Hoàng Mai, Hà Nội', 'BX Giáp Bát', 1, 1, '2025-05-25 02:59:10.000000', 1),
(3, 'Bình Thạnh, TP.HCM', 'BX Miền Đông', 2, 1, '2025-05-25 02:59:10.000000', 1),
(4, 'Bình Tân, TP.HCM', 'BX Miền Tây', 2, 1, '2025-05-25 02:59:10.000000', 1),
(5, 'Hải Châu, Đà Nẵng', 'BX Đà Nẵng', 3, 1, '2025-05-25 02:59:10.000000', 1),
(6, 'Lê Chân, Hải Phòng', 'BX Cầu Rào', 4, 1, '2025-05-25 02:59:10.000000', 1),
(7, 'TP Vinh, Nghệ An', 'BX Vinh', 5, 1, '2025-05-25 02:59:10.000000', 1),
(8, 'TP Hạ Long, Quảng Ninh', 'BX Bãi Cháy', 6, 1, '2025-05-25 02:59:10.000000', 1),
(9, 'TP Đông Hà, Quảng Trị', 'BX Đông Hà', 16, 1, '2025-05-25 02:59:10.000000', 1),
(10, 'TP Huế, Thừa Thiên Huế', 'BX Huế', 7, 1, '2025-05-25 02:59:10.000000', 1),
(11, 'Ninh Kiều, Cần Thơ', 'BX Cần Thơ', 8, 1, '2025-05-25 02:59:10.000000', 1),
(12, 'TP Đà Lạt, Lâm Đồng', 'BX Đà Lạt', 9, 1, '2025-05-25 02:59:10.000000', 1),
(13, 'Dĩ An, Bình Dương', 'BX Dĩ An', 10, 1, '2025-05-25 02:59:10.000000', 1),
(14, 'TP Phan Thiết, Bình Thuận', 'BX Phan Thiết', 11, 1, '2025-05-25 02:59:10.000000', 1),
(15, 'TP Tam Kỳ, Quảng Nam', 'BX Tam Kỳ', 12, 1, '2025-05-25 02:59:10.000000', 1),
(16, 'TP Nam Định, Nam Định', 'BX Nam Định', 13, 1, '2025-05-25 02:59:10.000000', 1),
(17, 'TP Thái Bình, Thái Bình', 'BX Thái Bình', 14, 1, '2025-05-25 02:59:10.000000', 1),
(18, 'TP Thanh Hóa, Thanh Hóa', 'BX Thanh Hóa', 15, 1, '2025-05-25 02:59:10.000000', 1),
(19, 'TP Ninh Bình, Ninh Bình', 'BX Ninh Bình', 17, 1, '2025-05-25 02:59:10.000000', 1),
(20, 'TP Việt Trì, Phú Thọ', 'BX Việt Trì', 20, 1, '2025-05-25 02:59:10.000000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `trip`
--

CREATE TABLE `trip` (
  `id` int NOT NULL,
  `arrival_date` date DEFAULT NULL,
  `arrival_time` time(6) DEFAULT NULL,
  `departure_date` date DEFAULT NULL,
  `departure_time` time(6) DEFAULT NULL,
  `total` int DEFAULT NULL,
  `route_id` int DEFAULT NULL,
  `special_id` int DEFAULT NULL,
  `stock` int DEFAULT NULL,
  `vehicle_id` int DEFAULT NULL,
  `pickup_id` varchar(255) DEFAULT NULL,
  `driver_id` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `trip`
--

INSERT INTO `trip` (`id`, `arrival_date`, `arrival_time`, `departure_date`, `departure_time`, `total`, `route_id`, `special_id`, `stock`, `vehicle_id`, `pickup_id`, `driver_id`, `status`, `created_at`, `created_by`) VALUES
(1, '2025-05-15', '12:00:00.000000', '2025-05-15', '08:00:00.000000', 45, 1, 1, 3, 1, 'P1', NULL, 1, '2025-05-25 02:58:05.000000', 1),
(2, '2025-05-16', '12:00:00.000000', '2025-05-15', '08:00:00.000000', 45, 1, 2, 12, 3, 'P16', NULL, 1, '2025-05-25 02:58:05.000000', 1),
(3, '2025-05-15', '12:00:00.000000', '2025-05-16', '20:00:00.000000', 45, 3, NULL, 10, 2, 'P8', NULL, 1, '2025-05-25 02:58:05.000000', 1),
(4, '2025-05-17', '17:00:00.000000', '2025-05-17', '08:00:00.000000', 45, 4, NULL, 3, 5, 'P10', NULL, 1, '2025-05-25 02:58:05.000000', 1),
(5, '2025-05-18', '23:00:00.000000', '2025-05-18', '10:00:00.000000', 45, 5, 3, 14, 4, 'P13', NULL, 1, '2025-05-25 02:58:05.000000', 1),
(6, '2025-05-21', '07:30:00.000000', '2025-05-20', '08:30:00.000000', 100, 3, NULL, 100, 1, 'P8', NULL, 1, '2025-05-25 02:58:05.000000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `type`
--

CREATE TABLE `type` (
  `id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `type`
--

INSERT INTO `type` (`id`, `name`, `status`, `created_at`, `created_by`) VALUES
(1, 'Ghế ngồi', 1, '2025-05-25 02:57:53.000000', 1),
(2, 'Giường nằm', 1, '2025-05-25 02:57:53.000000', 1),
(3, 'Limousine', 1, '2025-05-25 02:57:53.000000', 1),
(4, 'Ghế nhựa', 0, '2025-05-25 10:41:24.084000', 1);

-- --------------------------------------------------------

--
-- Table structure for table `vehicle`
--

CREATE TABLE `vehicle` (
  `id` int NOT NULL,
  `license` varchar(255) DEFAULT NULL,
  `seat_number` int DEFAULT NULL,
  `type_id` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `created_by` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `vehicle`
--

INSERT INTO `vehicle` (`id`, `license`, `seat_number`, `type_id`, `status`, `created_at`, `created_by`) VALUES
(1, '51A-00001', 16, 1, 1, '2025-05-16 09:56:40.000000', 1),
(2, '51A-00002', 29, 2, 1, '2025-05-16 09:56:40.000000', 1),
(3, '51A-00003', 45, 3, 1, '2025-05-16 09:56:40.000000', 1),
(4, '51A-00004', 30, 2, 1, '2025-05-16 09:56:40.000000', 1),
(5, '51A-00005', 16, 1, 1, '2025-05-16 09:56:40.000000', 1),
(6, '51A-00006', 29, 3, 1, '2025-05-16 09:56:40.000000', 1),
(7, '51A-00007', 45, 1, 1, '2025-05-16 09:56:40.000000', 1),
(8, '51A-00008', 30, 3, 1, '2025-05-16 09:56:40.000000', 1),
(9, '51A-00009', 16, 3, 1, '2025-05-16 09:56:40.000000', 1),
(10, '51A-00010', 29, 2, 1, '2025-05-16 09:56:40.000000', 1);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `log_type`
--
ALTER TABLE `log_type`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `log_vehicle`
--
ALTER TABLE `log_vehicle`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `pickup`
--
ALTER TABLE `pickup`
  ADD PRIMARY KEY (`id`),
  ADD KEY `FKjeea3rj5roemgumid0fct7in5` (`route_id`),
  ADD KEY `FK2opa7yb8aanv570ubyb02q9oe` (`station_id`);

--
-- Indexes for table `province`
--
ALTER TABLE `province`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `route`
--
ALTER TABLE `route`
  ADD PRIMARY KEY (`id`),
  ADD KEY `FKl3owc56nquwqu75528odjdf1q` (`end`),
  ADD KEY `FKk4cindug6pvrerd6d9ha8ng0s` (`start`);

--
-- Indexes for table `special_day`
--
ALTER TABLE `special_day`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `station`
--
ALTER TABLE `station`
  ADD PRIMARY KEY (`id`),
  ADD KEY `FKad8w80o8c90l9es4lo2alh6yp` (`province_id`);

--
-- Indexes for table `trip`
--
ALTER TABLE `trip`
  ADD PRIMARY KEY (`id`),
  ADD KEY `FKeva4adpyk6glllffnw5ypj20j` (`route_id`),
  ADD KEY `FKrrpr36y1k0upu67h1f3xc0mxv` (`special_id`),
  ADD KEY `FKrji8htecrp06ao6s7nfubswnr` (`vehicle_id`),
  ADD KEY `FKma43gpn91rhelbhgtovr57m5m` (`pickup_id`);

--
-- Indexes for table `type`
--
ALTER TABLE `type`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `vehicle`
--
ALTER TABLE `vehicle`
  ADD PRIMARY KEY (`id`),
  ADD KEY `FKspo7hdc8fxcdccgnr2lbje930` (`type_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `log_type`
--
ALTER TABLE `log_type`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `log_vehicle`
--
ALTER TABLE `log_vehicle`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `province`
--
ALTER TABLE `province`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT for table `route`
--
ALTER TABLE `route`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT for table `special_day`
--
ALTER TABLE `special_day`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `station`
--
ALTER TABLE `station`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT for table `trip`
--
ALTER TABLE `trip`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `type`
--
ALTER TABLE `type`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `vehicle`
--
ALTER TABLE `vehicle`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `pickup`
--
ALTER TABLE `pickup`
  ADD CONSTRAINT `FK2opa7yb8aanv570ubyb02q9oe` FOREIGN KEY (`station_id`) REFERENCES `station` (`id`),
  ADD CONSTRAINT `FKjeea3rj5roemgumid0fct7in5` FOREIGN KEY (`route_id`) REFERENCES `route` (`id`);

--
-- Constraints for table `route`
--
ALTER TABLE `route`
  ADD CONSTRAINT `FKk4cindug6pvrerd6d9ha8ng0s` FOREIGN KEY (`start`) REFERENCES `province` (`id`),
  ADD CONSTRAINT `FKl3owc56nquwqu75528odjdf1q` FOREIGN KEY (`end`) REFERENCES `province` (`id`);

--
-- Constraints for table `station`
--
ALTER TABLE `station`
  ADD CONSTRAINT `FKad8w80o8c90l9es4lo2alh6yp` FOREIGN KEY (`province_id`) REFERENCES `province` (`id`);

--
-- Constraints for table `trip`
--
ALTER TABLE `trip`
  ADD CONSTRAINT `FKeva4adpyk6glllffnw5ypj20j` FOREIGN KEY (`route_id`) REFERENCES `route` (`id`),
  ADD CONSTRAINT `FKma43gpn91rhelbhgtovr57m5m` FOREIGN KEY (`pickup_id`) REFERENCES `pickup` (`id`),
  ADD CONSTRAINT `FKrji8htecrp06ao6s7nfubswnr` FOREIGN KEY (`vehicle_id`) REFERENCES `vehicle` (`id`),
  ADD CONSTRAINT `FKrrpr36y1k0upu67h1f3xc0mxv` FOREIGN KEY (`special_id`) REFERENCES `special_day` (`id`);

--
-- Constraints for table `vehicle`
--
ALTER TABLE `vehicle`
  ADD CONSTRAINT `FKspo7hdc8fxcdccgnr2lbje930` FOREIGN KEY (`type_id`) REFERENCES `type` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
