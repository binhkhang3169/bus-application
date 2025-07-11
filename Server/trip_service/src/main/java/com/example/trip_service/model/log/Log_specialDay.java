package com.example.trip_service.model.log;

import jakarta.persistence.*;
import lombok.Data;

import java.util.Date;


@Entity
@Table(name = "log_specialDay")
@Data
public class Log_specialDay {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    private Date updatedAt;

    private Integer updatedBy;

    private Integer specialDayId;

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public Date getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(Date updatedAt) {
        this.updatedAt = updatedAt;
    }

    public Integer getUpdatedBy() {
        return updatedBy;
    }

    public void setUpdatedBy(Integer updatedBy) {
        this.updatedBy = updatedBy;
    }

    public Integer getSpecialDayId() {
        return specialDayId;
    }

    public void setSpecialDayId(Integer specialDayId) {
        this.specialDayId = specialDayId;
    }
}