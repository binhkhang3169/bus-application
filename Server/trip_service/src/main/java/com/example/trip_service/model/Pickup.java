package com.example.trip_service.model;

import jakarta.persistence.*;
import lombok.Data;

import java.util.Date;

@Entity
@Table(name = "pickup")
@Data
public class Pickup {

    @Id
    private String id;

    @ManyToOne
    @JoinColumn(name = "station_id")
    private Station station;

    private String selfId;
    private String time;

    private int path_id;

    private Integer status=1;

    private Date createdAt;

    private Integer createdBy;

    public Date getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Date createdAt) {
        this.createdAt = createdAt;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Station getStation() {
        return station;
    }

    public void setStation(Station station) {
        this.station = station;
    }

    public String getSelfId() {
        return selfId;
    }

    public void setSelfId(String selfId) {
        this.selfId = selfId;
    }

    public String getTime() {
        return time;
    }

    public void setTime(String time) {
        this.time = time;
    }

    public int getPath_id() {
        return path_id;
    }

    public void setPath_id(int path_id) {
        this.path_id = path_id;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public Route getRoute() {
        return route;
    }

    public void setRoute(Route route) {
        this.route = route;
    }

    public Integer getCreatedBy() {
        return createdBy;
    }

    public void setCreatedBy(Integer createdBy) {
        this.createdBy = createdBy;
    }


    @ManyToOne
    @JoinColumn(name = "route_id")
    private Route route;

}
