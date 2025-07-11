import re  # Add this import
import mysql.connector
import psycopg2
import os
from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher
from rasa_sdk.events import SlotSet
from typing import Any, Text, Dict, List

# Hàm dùng chung để lấy config DB từ biến môi trường
# Cấu hình riêng cho từng DB
# Cấu hình Postgres cho Ticket
PG_HOST = os.environ.get("PG_HOST", "localhost")
PG_USER = os.environ.get("PG_USER", "postgres")
PG_PASSWORD = os.environ.get("PG_PASSWORD", "root")
PG_NAME = os.environ.get("PG_NAME", "ticket_service")
PG_PORT = os.environ.get("PG_PORT", "5432")
PG_SSLMODE = os.environ.get("PG_SSLMODE", "prefer")

def get_pg_connection():
    return psycopg2.connect(
        host=PG_HOST,
        user=PG_USER,
        password=PG_PASSWORD,
        dbname=PG_NAME,
        port=PG_PORT,
        sslmode=PG_SSLMODE
    )

# Cấu hình MySQL cho các service khác
MYSQL_HOST = os.environ.get("MYSQL_HOST", "localhost")
MYSQL_USER = os.environ.get("MYSQL_USER", "root")
MYSQL_PASSWORD = os.environ.get("MYSQL_PASSWORD", "root")
MYSQL_NAME = os.environ.get("MYSQL_NAME", "trip_service")
MYSQL_PORT = os.environ.get("MYSQL_PORT", "3306")

def get_mysql_connection():
    return mysql.connector.connect(
        host=MYSQL_HOST,
        user=MYSQL_USER,
        password=MYSQL_PASSWORD,
        database=MYSQL_NAME,
        port=MYSQL_PORT
    )

############################################Serive ticket 
class ActionCheckTicketStatus(Action):
    def name(self) -> Text:
        return "action_check_ticket_status"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        phone_number = tracker.get_slot("phone_number")
        print(f"[DEBUG] phone_number: '{phone_number}'")

        if not phone_number:
            dispatcher.utter_message(text="Bạn vui lòng cung cấp số điện thoại để tôi kiểm tra vé nhé.")
            return []
        # Kết nối tới Postgres Ticket Service
        conn = get_pg_connection()
        cursor = conn.cursor()
        query = "SELECT status FROM Ticket WHERE phone LIKE %s LIMIT 1"
        cursor.execute(query, (f"%{phone_number}%",))
        result = cursor.fetchone()

        if result:
            dispatcher.utter_message(text=f"Vé của bạn hiện đang có trạng thái: {result[0]}")
        else:
            dispatcher.utter_message(text="Không tìm thấy vé với số điện thoại này.")

        cursor.close()
        conn.close()
        return []

class ActionAskTicketDetails(Action):
    def name(self) -> Text:
        return "action_ask_ticket_details"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        phone_number = tracker.get_slot("phone_number")
        entities = tracker.latest_message.get("entities", [])
        print(f"[DEBUG] phone_number slot: '{phone_number}'")
        print(f"[DEBUG] entities extracted: {entities}")

        if not phone_number:
            dispatcher.utter_message(response="utter_no_phone_number")
            return []

        conn = get_mysql_connection()
        cursor = conn.cursor()

        query = "SELECT route, time, seat_number, status FROM tickets WHERE phone = %s"
        cursor.execute(query, (phone_number,))
        result = cursor.fetchone()

        if result:
            dispatcher.utter_message(text=f"Thông tin vé: Tuyến {result[0]}, giờ {result[1]}, số ghế {result[2]}, trạng thái {result[3]}.")
        else:
            dispatcher.utter_message(text="Không tìm thấy thông tin vé với số điện thoại này.")

        cursor.close()
        conn.close()
        return []
############################################Serive ticket 


# Ok
class ActionGetBusInfo(Action):
    def name(self) -> Text:
        return "action_get_bus_info"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        conn = get_mysql_connection()
        cursor = conn.cursor()

        query = """
        SELECT
        t.id AS `Chuyến`,
        s.name AS `Điểm đi`,
        e.name AS `Điểm đến`,
        t.arrival_date as `Ngày đi`,
        t.arrival_time as `Khởi hành`,
        t.departure_date as `Ngày đến`,
        t.departure_time as `Đến nơi`,
        t.stock as `Còn trống`,
        r.price as `Giá`,
        type.name as `Loại xe`
        FROM trip t
        JOIN route r ON t.route_id = r.id
        JOIN vehicle v ON t.vehicle_id = v.id
        JOIN type type ON type.id = v.type_id
        JOIN province s ON r.start = s.id
        JOIN province e ON r.end = e.id
        WHERE t.status = 1
        AND r.status = 1
        AND s.status = 1
        AND e.status = 1
        AND t.arrival_date >= CURRENT_DATE
        ORDER BY t.departure_date, t.departure_time
        LIMIT 5
        """

        cursor.execute(query)
        buses = cursor.fetchall()
        
        if buses:
            message = "Dưới đây là một số chuyến xe hiện có:\n"
            for b in buses:
                price = b[8]
                message += (f"- Chuyến: {b[0]}, Điểm đi: {b[1]}, Điểm đến: {b[2]}, "
                            f"Ngày đi: {b[3]} {b[4]}, Ngày đến: {b[5]} {b[6]}, "
                            f"Còn trống: {b[7]}, Giá: {b[8]:,} VND, Loại xe: {b[9]}\n")
            dispatcher.utter_message(text=message)
        else:
            dispatcher.utter_message(text="Hiện tại không có thông tin chuyến xe.")


        cursor.close()
        conn.close()
        return []

class ActionAskContact(Action):
    def name(self) -> Text:
        return "action_ask_contact"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:
        dispatcher.utter_message(text="Bạn có thể liên hệ với nhà xe qua số: 1900 1234 hoặc email: hotro@nhaxe.vn")
        return []



class ActionAskPrice(Action):
    def name(self) -> Text:
        return "action_ask_price"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        import re

        route = tracker.get_slot("route")
        entities = tracker.latest_message.get("entities", [])
        print(f"[DEBUG] route slot: '{route}'")
        print(f"[DEBUG] entities extracted: {entities}")

        route_entities = [entity['value'] for entity in entities if entity['entity'] == 'route']
        print(f"[DEBUG] route entities: {route_entities}")

        if len(route_entities) >= 2:
            start_point = route_entities[0]
            end_point = route_entities[1]
            print(f"[DEBUG] start_point: '{start_point}', end_point: '{end_point}' from entities")
        else:
            if route:
                parts = [s.strip() for s in route.split('-')]
                if len(parts) == 2:
                    start_point, end_point = parts
                    print(f"[DEBUG] start_point: '{start_point}', end_point: '{end_point}' from route slot")
                else:
                    input_text = tracker.latest_message.get("text", "")
                    route_pattern = r"([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)-([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)"
                    match = re.search(route_pattern, input_text)
                    if match:
                        start_point = match.group(1).strip()
                        end_point = match.group(2).strip()
                        print(f"[DEBUG] start_point: '{start_point}', end_point: '{end_point}' from input text")
                    else:
                        dispatcher.utter_message(response="utter_no_route_found")
                        return []
            else:
                dispatcher.utter_message(response="utter_ask_route")
                return []

        conn = get_mysql_connection()
        cursor = conn.cursor()

        query = """
        SELECT r.price, s.name, e.name
        FROM route r
        JOIN province s ON r.start = s.id
        JOIN province e ON r.end = e.id
        WHERE (s.name LIKE %s AND e.name LIKE %s) OR (s.name LIKE %s AND e.name LIKE %s)
        """
        cursor.execute(query, (f"%{start_point}%", f"%{end_point}%", f"%{end_point}%", f"%{start_point}%"))
        results = cursor.fetchall()

        if results:
            messages = []
            seen_pairs = set()  # Lưu các tuyến đã xử lý (theo cặp sắp xếp)
            for price, s_name, e_name in results:
                # Chuẩn hóa cặp tuyến: sorted để đảm bảo Hà Nội - Huế và Huế - Hà Nội là 1
                route_key = tuple(sorted([s_name.strip(), e_name.strip()]))
                if route_key not in seen_pairs:
                    seen_pairs.add(route_key)
                    price_formatted = f"{price:,} VND"
                    messages.append(f"Giá vé tuyến {s_name} - {e_name} là {price_formatted}")
            dispatcher.utter_message(text="\n".join(messages))
        else:
            dispatcher.utter_message(response="utter_no_route_found")

        cursor.close()
        conn.close()
        return []

class ActionAskBusType(Action):
    def name(self) -> Text:
        return "action_ask_bus_type"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        import re
        route = tracker.get_slot("route")
        entities = tracker.latest_message.get("entities", [])
        print(f"[DEBUG] route slot: '{route}'")
        print(f"[DEBUG] entities extracted: {entities}")

        route_entities = [entity['value'] for entity in entities if entity['entity'] == 'route']
        print(f"[DEBUG] route entities: {route_entities}")

        start_point = end_point = None

        if len(route_entities) >= 2:
            start_point = route_entities[0]
            end_point = route_entities[1]
            print(f"[DEBUG] start_point: '{start_point}', end_point: '{end_point}' from entities")
        else:
            if route:
                parts = [s.strip() for s in route.split('-')]
                if len(parts) == 2:
                    start_point, end_point = parts
                    print(f"[DEBUG] start_point: '{start_point}', end_point: '{end_point}' from route slot")
                else:
                    input_text = tracker.latest_message.get("text", "")
                    route_pattern = r"([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)-([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)"
                    match = re.search(route_pattern, input_text)
                    if match:
                        start_point = match.group(1).strip()
                        end_point = match.group(2).strip()
                        print(f"[DEBUG] start_point: '{start_point}', end_point: '{end_point}' from input text")
                    else:
                        dispatcher.utter_message(response="utter_no_trip_found")
                        return []
            else:
                dispatcher.utter_message(text=f"Hiện tại không có chuyến xe nào giữa {start_point} và {end_point}. Vui lòng thử tuyến khác.")
                return []

        conn = get_mysql_connection()
        cursor = conn.cursor()

        # Truy vấn tất cả tuyến phù hợp theo 2 chiều
        query = """
        SELECT v.license,type.name, s.name, e.name
        FROM trip t
        JOIN route r ON t.route_id = r.id
        JOIN vehicle v ON t.vehicle_id = v.id
        JOIN type  ON v.type_id = type.id
        JOIN province s ON r.start = s.id
        JOIN province e ON r.end = e.id
        WHERE (s.name LIKE %s AND e.name LIKE %s)
           OR (e.name LIKE %s AND s.name LIKE %s)
        """
        cursor.execute(query, (f"%{start_point}%", f"%{end_point}%", f"%{end_point}%", f"%{start_point}%"))
        results = cursor.fetchall()

        if results:
            messages = []
            seen_routes = set()
            for bus_type, amenities, s_name, e_name in results:
                route_key = tuple(sorted([s_name.strip(), e_name.strip()]))
                if route_key not in seen_routes:
                    seen_routes.add(route_key)
                    messages.append(
                        f"Tuyến {s_name} - {e_name}: Xe khởi hành {bus_type} - {amenities}."
                    )
            dispatcher.utter_message(text="\n".join(messages))
        else:
            dispatcher.utter_message(response="utter_no_route_found")

        cursor.close()
        conn.close()
        return []

class ActionAskPickupDropoff(Action):
    def name(self) -> Text:
        return "action_ask_pickup_dropoff"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        import re
        import mysql.connector

        route = tracker.get_slot("route")
        entities = tracker.latest_message.get("entities", [])
        print(f"[DEBUG] route slot: '{route}'")
        print(f"[DEBUG] entities extracted: {entities}")

        route_entities = [entity['value'] for entity in entities if entity['entity'] == 'route']
        print(f"[DEBUG] route entities: {route_entities}")

        from_province = to_province = None

        # Ưu tiên lấy từ entity
        if len(route_entities) >= 2:
            from_province = route_entities[0]
            to_province = route_entities[1]
            print(f"[DEBUG] from_province: '{from_province}', to_province: '{to_province}' from entities")
        else:
            # Tiếp theo lấy từ slot
            if route:
                parts = [s.strip() for s in route.split('-')]
                if len(parts) == 2:
                    from_province, to_province = parts
                    print(f"[DEBUG] from_province: '{from_province}', to_province: '{to_province}' from route slot")
                else:
                    # Cuối cùng dùng regex từ văn bản người dùng
                    input_text = tracker.latest_message.get("text", "")
                    route_pattern = r"([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)-([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)"
                    match = re.search(route_pattern, input_text)
                    if match:
                        from_province = match.group(1).strip()
                        to_province = match.group(2).strip()
                        print(f"[DEBUG] from_province: '{from_province}', to_province: '{to_province}' from input text")
                    else:
                        dispatcher.utter_message(response="utter_no_route_found")
                        return []

        # Kiểm tra nếu không có đủ dữ liệu
        if not from_province or not to_province:
            dispatcher.utter_message(text="Vui lòng cung cấp đầy đủ thông tin về tuyến đường (ví dụ: Hà Nội - Huế).")
            return []

        # Kết nối DB và truy vấn
        conn = get_mysql_connection()
        cursor = conn.cursor()

        query = """
        SELECT 
            s_start.name AS departureStation,
            s_end.name AS arrivalStation,
            (
                SELECT GROUP_CONCAT(CONCAT(st.name, ' (', st.address, ')') ORDER BY p.id SEPARATOR ' → ')
                FROM pickup p
                JOIN station st ON p.station_id = st.id
                WHERE p.route_id = r.id
                  AND p.path_id = p_start.path_id
                GROUP BY p.path_id
            ) AS fullRoute
        FROM route r
        JOIN pickup p_start ON p_start.route_id = r.id AND p_start.self_id = '-1'
        JOIN station s_start ON p_start.station_id = s_start.id
        JOIN pickup p_end ON p_end.route_id = r.id AND p_end.self_id = '-2' AND p_end.path_id = p_start.path_id
        JOIN station s_end ON p_end.station_id = s_end.id
        JOIN province p_from ON r.start = p_from.id
        JOIN province p_to ON r.end = p_to.id
        WHERE LOWER(p_from.name) = LOWER(%s)
          AND LOWER(p_to.name) = LOWER(%s)
          AND r.status = 1;
        """

        cursor.execute(query, (from_province, to_province))
        results = cursor.fetchall()

        if results:
            message = f"🔍 Tìm thấy {len(results)} tuyến từ {from_province} đến {to_province}:\n\n"
            for i, (departure, arrival, full_route) in enumerate(results, 1):
                message += f"🚌 Tuyến {i}:\n"
                message += f"• Điểm đón: {departure}\n"
                message += f"• Điểm trả: {arrival}\n"
                message += "• Lộ trình:\n"
                for point in full_route.split(" → "):
                    message += f"  - {point}\n"
                message += "\n"

            dispatcher.utter_message(text=message.strip())
        else:
            dispatcher.utter_message(response="utter_no_route_found")

        cursor.close()
        conn.close()
        return []

    
class ActionAskTravelTime(Action):
    def name(self) -> Text:
        return "action_ask_travel_time"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        import re
        import mysql.connector

        route = tracker.get_slot("route")
        entities = tracker.latest_message.get("entities", [])
        print(f"[DEBUG] route slot: '{route}'")
        print(f"[DEBUG] entities extracted: {entities}")

        route_entities = [entity['value'] for entity in entities if entity['entity'] == 'route']
        print(f"[DEBUG] route entities: {route_entities}")

        from_province = to_province = None

        # Ưu tiên lấy từ entity
        if len(route_entities) >= 2:
            from_province = route_entities[0]
            to_province = route_entities[1]
            print(f"[DEBUG] from_province: '{from_province}', to_province: '{to_province}' from entities")
        else:
            # Nếu có slot route thì xử lý từ slot
            if route:
                parts = [s.strip() for s in route.split('-')]
                if len(parts) == 2:
                    from_province, to_province = parts
                    print(f"[DEBUG] from_province: '{from_province}', to_province: '{to_province}' from route slot")
                else:
                    # Cuối cùng dùng regex từ input text
                    input_text = tracker.latest_message.get("text", "")
                    route_pattern = r"([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)-([A-Za-zÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝĂĐĨŨƠƯ\s]+)"
                    match = re.search(route_pattern, input_text)
                    if match:
                        from_province = match.group(1).strip()
                        to_province = match.group(2).strip()
                        print(f"[DEBUG] from_province: '{from_province}', to_province: '{to_province}' from input text")
                    else:
                        dispatcher.utter_message(response="utter_no_route_found")
                        return []

        if not from_province or not to_province:
            dispatcher.utter_message(text="Vui lòng cung cấp tuyến đường hợp lệ (ví dụ: Hà Nội - Đà Nẵng).")
            return []

        # Kết nối cơ sở dữ liệu
        conn = get_mysql_connection()
        cursor = conn.cursor()

        query = """
        SELECT r.estimated_time
        FROM route r
        JOIN province p_from ON r.start = p_from.id
        JOIN province p_to ON r.end = p_to.id
        WHERE LOWER(p_from.name) = LOWER(%s)
          AND LOWER(p_to.name) = LOWER(%s)
          AND r.status = 1
        LIMIT 1;
        """

        cursor.execute(query, (from_province, to_province))
        result = cursor.fetchone()

        if result:
            dispatcher.utter_message(
                text=f"⏱️ Thời gian di chuyển từ {from_province} đến {to_province} là khoảng {result[0]}."
            )
        else:
            dispatcher.utter_message(response="utter_no_route_found")

        cursor.close()
        conn.close()
        return []

