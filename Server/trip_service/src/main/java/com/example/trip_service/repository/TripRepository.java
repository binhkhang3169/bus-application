package com.example.trip_service.repository;

import java.time.LocalDate;
import java.time.LocalTime;
import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import com.example.trip_service.dto.TripInfoProjection;
import com.example.trip_service.model.Trip;

public interface TripRepository extends JpaRepository<Trip, Integer> {

  @Query(value = """
      SELECT t.* FROM trip t
      WHERE t.driver_id = :driverId
        AND t.status != 0
        AND (:tripId IS NULL OR t.id != :tripId) -- Exclude current trip if updating
        AND (
          -- Check if existing trip's departure is before the new trip's sensitive range ends
          (t.departure_date < :checkRangeEndDate OR (t.departure_date = :checkRangeEndDate AND t.departure_time < :checkRangeEndTime))
          AND
          -- Check if existing trip's arrival is after the new trip's sensitive range starts
          (t.arrival_date > :checkRangeStartDate OR (t.arrival_date = :checkRangeStartDate AND t.arrival_time > :checkRangeStartTime))
        )
      LIMIT 1
      """, nativeQuery = true)
  Trip findConflictingTripForDriver(
      @Param("driverId") Integer driverId,
      @Param("checkRangeStartDate") LocalDate checkRangeStartDate,
      @Param("checkRangeStartTime") LocalTime checkRangeStartTime,
      @Param("checkRangeEndDate") LocalDate checkRangeEndDate,
      @Param("checkRangeEndTime") LocalTime checkRangeEndTime,
      @Param("tripId") Integer tripId);

  @Query(value = """
      SELECT
          t.id AS tripId,
          CAST(t.vehicle_id AS CHAR) AS vehicleId,
          v.license AS license,
          ty.name AS vehicleType,
          t.status AS status,
          CAST(t.departure_date AS CHAR) AS departureDate,
          CAST(t.departure_time AS CHAR) AS departureTime,
          CAST(t.arrival_date AS CHAR) AS arrivalDate,
          CAST(t.arrival_time AS CHAR) AS arrivalTime,
          t.stock AS stock,
          r.price AS price,
          CAST(r.distance AS CHAR) AS estimatedDistance,
          r.estimated_time AS estimatedTime,
          s_start.name AS departureStation,
          s_end.name AS arrivalStation,
          (
              SELECT GROUP_CONCAT(st.name ORDER BY p.id SEPARATOR ' → ')
              FROM pickup p
              JOIN station st ON p.station_id = st.id
              WHERE p.route_id = r.id
                AND p.path_id = p_start.path_id
              GROUP BY p.path_id
          ) AS fullRoute
      FROM trip t
      JOIN route r
        ON t.route_id = r.id
      JOIN pickup p_start
        ON p_start.route_id = r.id
       AND p_start.self_id = '-1'
      JOIN station s_start
        ON p_start.station_id = s_start.id
      JOIN pickup p_end
        ON p_end.route_id = r.id
       AND p_end.self_id = '-2'
       AND p_end.path_id = p_start.path_id
      JOIN station s_end
        ON p_end.station_id = s_end.id
      JOIN vehicle v ON t.vehicle_id = v.id
      JOIN type ty ON v.type_id = ty.id
      WHERE r.id = :routeId
        AND t.departure_date = :departureDate
        AND t.stock >= :quantity
      """, nativeQuery = true)
  List<TripInfoProjection> findTripsByRouteAndDateRaw(
      @Param("routeId") Integer routeId,
      @Param("departureDate") String departureDate,
      @Param("quantity") Integer quantity);

  // @Query(value = """
  // SELECT
  // t.id AS tripId,
  // CAST(t.vehicle_id AS CHAR) AS vehicleId,
  // v.license AS license,
  // ty.name AS vehicleType,
  // t.status AS status,
  // CAST(t.departure_date AS CHAR) AS departureDate,
  // CAST(t.departure_time AS CHAR) AS departureTime,
  // CAST(t.arrival_date AS CHAR) AS arrivalDate,
  // CAST(t.arrival_time AS CHAR) AS arrivalTime,
  // t.stock AS stock,
  // r.price AS price,
  // CAST(r.distance AS CHAR) AS estimatedDistance,
  // r.estimated_time AS estimatedTime,
  // s_start.name AS departureStation,
  // s_end.name AS arrivalStation,
  // (
  // SELECT GROUP_CONCAT(st.name ORDER BY p.id SEPARATOR ' → ')
  // FROM pickup p
  // JOIN station st ON p.station_id = st.id
  // WHERE p.route_id = r.id
  // AND p.path_id = p_start.path_id
  // GROUP BY p.path_id
  // ) AS fullRoute
  // FROM trip t
  // JOIN route r
  // ON t.route_id = r.id
  // JOIN pickup p_start
  // ON p_start.route_id = r.id
  // AND p_start.self_id = '-1'
  // AND p_start.id = t.pickup_id
  // JOIN station s_start
  // ON p_start.station_id = s_start.id
  // JOIN pickup p_end
  // ON p_end.route_id = r.id
  // AND p_end.self_id = '-2'
  // AND p_end.path_id = p_start.path_id
  // JOIN station s_end
  // ON p_end.station_id = s_end.id
  // JOIN vehicle v ON t.vehicle_id = v.id
  // JOIN type ty ON v.type_id = ty.id
  // JOIN province p_from ON r.start = p_from.id
  // JOIN province p_to ON r.end = p_to.id
  // WHERE p_from.id = :fromProvinceId
  // AND p_to.id = :toProvinceId
  // AND t.departure_date = :departureDate
  // AND t.stock >= :quantity
  // AND t.status = 1
  // """, nativeQuery = true)
  // List<TripInfoProjection> findTripsByLocationsAndDate(
  // @Param("fromProvinceId") Integer fromProvinceId,
  // @Param("toProvinceId") Integer toProvinceId,
  // @Param("departureDate") String departureDate,
  // @Param("quantity") Integer quantity);
  @Query(value = """
      WITH RECURSIVE route_cte AS (
          SELECT
              p.route_id,
              p.path_id,
              st.name,
              p.id,
              p.self_id,
              1 AS level
          FROM pickup p
          JOIN station st ON st.id = p.station_id
          WHERE p.self_id = '-1' COLLATE utf8mb4_unicode_ci
          UNION ALL
          SELECT
              p2.route_id,
              p2.path_id,
              st2.name,
              p2.id,
              p2.self_id,
              cte.level + 1
          FROM pickup p2
          JOIN station st2 ON st2.id = p2.station_id
          JOIN route_cte cte ON p2.self_id = CAST(cte.id AS CHAR) COLLATE utf8mb4_unicode_ci
          WHERE p2.route_id = cte.route_id
            AND p2.path_id = cte.path_id
      )
      SELECT
          t.id AS tripId,
          CAST(t.vehicle_id AS CHAR) AS vehicleId,
          v.license AS license,
          ty.name AS vehicleType,
          t.status AS status,
          CAST(t.departure_date AS CHAR) AS departureDate,
          CAST(t.departure_time AS CHAR) AS departureTime,
          CAST(t.arrival_date AS CHAR) AS arrivalDate,
          CAST(t.arrival_time AS CHAR) AS arrivalTime,
          t.stock AS stock,
          r.price AS price,
          CAST(r.distance AS CHAR) AS estimatedDistance,
          r.estimated_time AS estimatedTime,
          s_start.name AS departureStation,
          s_end.name AS arrivalStation,
          (
              SELECT GROUP_CONCAT(name ORDER BY level SEPARATOR ' → ')
              FROM (
                  SELECT name, level
                  FROM route_cte
                  WHERE route_id = r.id AND path_id = p_start.path_id
                  UNION ALL
                  SELECT s_end.name, 9999
              ) AS full
          ) AS fullRoute
      FROM trip t
      JOIN route r ON t.route_id = r.id
      JOIN pickup p_start ON p_start.route_id = r.id
          AND p_start.self_id = '-1' COLLATE utf8mb4_unicode_ci
          AND p_start.id = t.pickup_id
      JOIN station s_start ON p_start.station_id = s_start.id
      JOIN pickup p_end ON p_end.route_id = r.id
          AND p_end.self_id = '-2' COLLATE utf8mb4_unicode_ci
          AND p_end.path_id = p_start.path_id
      JOIN station s_end ON p_end.station_id = s_end.id
      JOIN vehicle v ON t.vehicle_id = v.id
      JOIN type ty ON v.type_id = ty.id
      JOIN province p_from ON r.start = p_from.id
      JOIN province p_to ON r.end = p_to.id
      WHERE p_from.id = :fromProvinceId
        AND p_to.id = :toProvinceId
        AND t.departure_date = :departureDate
        AND t.stock >= :quantity
        AND t.status = 1
      """, nativeQuery = true)
  List<TripInfoProjection> findTripsByLocationsAndDate(
      @Param("fromProvinceId") Integer fromProvinceId,
      @Param("toProvinceId") Integer toProvinceId,
      @Param("departureDate") String departureDate,
      @Param("quantity") Integer quantity);

  @Query(value = """
      SELECT
          t.id AS tripId,
          CAST(t.vehicle_id AS CHAR) AS vehicleId,
          v.license AS license,
          ty.name AS vehicleType,
          t.status AS status,
          CAST(t.departure_date AS CHAR) AS departureDate,
          CAST(t.departure_time AS CHAR) AS departureTime,
          CAST(t.arrival_date AS CHAR) AS arrivalDate,
          CAST(t.arrival_time AS CHAR) AS arrivalTime,
          t.stock AS stock,
          r.price AS price,
          CAST(r.distance AS CHAR) AS estimatedDistance,
          r.estimated_time AS estimatedTime,
          s_start.name AS departureStation,
          s_end.name AS arrivalStation,
          (
              SELECT GROUP_CONCAT(CONCAT(st.name, ' (', st.id, ')') ORDER BY p.id SEPARATOR ' → ')
              FROM pickup p
              JOIN station st ON p.station_id = st.id
              WHERE p.route_id = r.id
                AND p.path_id = p_start.path_id
              GROUP BY p.path_id
          ) AS fullRoute
      FROM trip t
      JOIN route r
        ON t.route_id = r.id
      JOIN pickup p_start
        ON p_start.route_id = r.id
       AND p_start.self_id = '-1'
       AND p_start.id = t.pickup_id
      JOIN station s_start
        ON p_start.station_id = s_start.id
      JOIN pickup p_end
        ON p_end.route_id = r.id
       AND p_end.self_id = '-2'
       AND p_end.path_id = p_start.path_id
      JOIN station s_end
        ON p_end.station_id = s_end.id
      JOIN vehicle v ON t.vehicle_id = v.id
      JOIN type ty ON v.type_id = ty.id
      WHERE t.id = :tripId
      """, nativeQuery = true)
  List<TripInfoProjection> findTripInfoById(@Param("tripId") Integer tripId);

  List<Trip> findByStatusIn(List<Integer> statuses);

  // New method to find trips by driver ID
  List<Trip> findByDriverId(Integer driverId);
}