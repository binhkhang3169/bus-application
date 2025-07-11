package com.example.trip_service.dto;

import java.util.Date;

// DTO chứa thông tin cho sự kiện tìm kiếm chuyến đi
public class TripSearchEvent {

    private Integer fromProvinceId;
    private Integer toProvinceId;
    private String departureDate;
    private Date searchTimestamp;
    // === THAY ĐỔI BẮT ĐẦU ===
    private Integer quantity; // Số lượng ghế tìm kiếm
    private Integer userId; // ID của người dùng thực hiện tìm kiếm (có thể null)
    // === THAY ĐỔI KẾT THÚC ===

    // === CONSTRUCTOR ĐÃ CẬP NHẬT ===
    public TripSearchEvent(Integer fromProvinceId, Integer toProvinceId, String departureDate, Date searchTimestamp,
            Integer quantity, Integer userId) {
        this.fromProvinceId = fromProvinceId;
        this.toProvinceId = toProvinceId;
        this.departureDate = departureDate;
        this.searchTimestamp = searchTimestamp;
        this.quantity = quantity;
        this.userId = userId;
    }

    // Getters and Setters
    public Integer getFromProvinceId() {
        return fromProvinceId;
    }

    public void setFromProvinceId(Integer fromProvinceId) {
        this.fromProvinceId = fromProvinceId;
    }

    public Integer getToProvinceId() {
        return toProvinceId;
    }

    public void setToProvinceId(Integer toProvinceId) {
        this.toProvinceId = toProvinceId;
    }

    public String getDepartureDate() {
        return departureDate;
    }

    public void setDepartureDate(String departureDate) {
        this.departureDate = departureDate;
    }

    public Date getSearchTimestamp() {
        return searchTimestamp;
    }

    public void setSearchTimestamp(Date searchTimestamp) {
        this.searchTimestamp = searchTimestamp;
    }

    // === GETTERS AND SETTERS MỚI ===
    public Integer getQuantity() {
        return quantity;
    }

    public void setQuantity(Integer quantity) {
        this.quantity = quantity;
    }

    public Integer getUserId() {
        return userId;
    }

    public void setUserId(Integer userId) {
        this.userId = userId;
    }
}