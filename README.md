# Dự án Microservices - Hướng dẫn Triển khai Toàn diện

Tài liệu này cung cấp hướng dẫn đầy đủ để chạy và triển khai hệ thống microservices từ môi trường phát triển cục bộ đến production trên Azure Kubernetes Service (AKS).

## 📋 Mục lục

- [Yêu cầu hệ thống](#yêu-cầu-hệ-thống)
- [1. Chạy dự án với Docker Compose](#1-chạy-dự-án-với-docker-compose)
- [2. Chạy ứng dụng Frontend](#2-chạy-ứng-dụng-frontend)
- [3. Deploy lên Azure Kubernetes Service (AKS)](#3-deploy-lên-azure-kubernetes-service-aks)
- [4. ETL Pipeline với Apache Airflow](#4-etl-pipeline-với-apache-airflow)

## 🔧 Yêu cầu hệ thống

### Môi trường phát triển cục bộ:
- **Docker** và **Docker Compose** (phiên bản mới nhất)
- **Node.js** (v16 trở lên) và **npm/yarn**
- **Flutter SDK** (phiên bản ổn định)
- **Git**

### Môi trường production:
- **Azure CLI**
- **kubectl**
- **Docker** (để build images)

---

## 1. Chạy dự án với Docker Compose

Docker Compose cung cấp cách dễ dàng nhất để chạy toàn bộ hệ thống microservices trong môi trường phát triển.

### 🚀 Khởi động nhanh

```bash
# Clone repository
git clone <repository-url>
cd <project-directory>
cd Server

# Tạo file environment variables
cp .env.example .env

# Chỉnh sửa file .env với thông tin cấu hình của bạn
nano .env
```

### 📝 Cấu hình Environment Variables

Tạo file `.env` trong thư mục gốc với các biến sau:

```env
# Database Configuration
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=password123
POSTGRES_DB=ticketdb

MYSQL_HOST=mysql
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password123
MYSQL_DATABASE=invoicedb

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis123

# Kafka Configuration
KAFKA_BOOTSTRAP_SERVERS=kafka:9092
KAFKA_TOPIC_TRIP_SEARCH=trip_search

# JWT & Security
JWT_SECRET_KEY=your-super-secret-jwt-key-here
VNPAY_HASH_SECRET=your-vnpay-hash-secret
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_stripe_webhook_secret

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-gmail-app-password

# File Upload
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
```

### 🐳 Chạy với Docker Compose

```bash
# Khởi động tất cả services
docker-compose up -d

# Xem logs của tất cả services
docker-compose logs -f

# Xem logs của một service cụ thể
docker-compose logs -f api-gateway

# Kiểm tra trạng thái các services
docker-compose ps

# Dừng tất cả services
docker-compose down

# Dừng và xóa volumes (dữ liệu sẽ bị mất)
docker-compose down -v
```

### 🔍 Kiểm tra các services

Sau khi khởi động thành công, các services sẽ có sẵn tại:

- **API Gateway**: http://localhost:8080
- **Payment Service**: http://localhost:8083
- **Ticket Service**: http://localhost:8081
- **User Service**: http://localhost:8082
- **Notification Service**: http://localhost:8084
- **Chatbot Service**: http://localhost:5005
- **PostgreSQL**: localhost:5432
- **MySQL**: localhost:3306
- **Redis**: localhost:6379
- **Kafka UI**: http://localhost:8090

### 🛠️ Troubleshooting Docker Compose

**Lỗi port đã được sử dụng:**
```bash
# Tìm và dừng process sử dụng port
lsof -ti:8080 | xargs kill -9

# Hoặc thay đổi port trong docker-compose.yml
```

**Lỗi kết nối database:**
```bash
# Kiểm tra logs database
docker-compose logs postgres
docker-compose logs mysql

# Restart database services
docker-compose restart postgres mysql
```

**Rebuild services sau khi thay đổi code:**
```bash
# Rebuild một service cụ thể
docker-compose build api-gateway

# Rebuild tất cả services
docker-compose build

# Restart với image mới
docker-compose up -d --build
```

---

## 2. Chạy ứng dụng Frontend

Dự án bao gồm hai ứng dụng frontend: **React Web App** và **Flutter Mobile App**.

### 🌐 React Web Application

#### Cài đặt và chạy:

```bash
# Di chuyển vào thư mục frontend
cd client
cd web

# Cài đặt dependencies
npm install

npm start
```

**Ứng dụng sẽ chạy tại**: http://localhost:3000

### 📱 Flutter Mobile Application

#### Cài đặt và chạy:

```bash
# Di chuyển vào thư mục mobile
cd client
cd user

# Cài đặt dependencies
flutter pub get

# Kiểm tra thiết bị kết nối
flutter devices

# Chạy ứng dụng
flutter run


cd driver


# Cài đặt dependencies
flutter pub get

# Kiểm tra thiết bị kết nối
flutter devices

# Chạy ứng dụng
flutter run
```


### 🔗 Tích hợp Frontend - Backend

##### Kiểm tra kết nối API:

**Từ React App:**
```javascript
// Test API connection
fetch('http://localhost:8080/api/health')
  .then(response => response.json())
  .then(data => console.log('API Health:', data))
  .catch(error => console.error('API Error:', error));
```

**Từ Flutter App:**
```dart
// Test API connection
import 'package:http/http.dart' as http;

Future<void> testApiConnection() async {
  try {
    final response = await http.get(
      Uri.parse('http://localhost:8080/api/health'),
    );
    print('API Status: ${response.statusCode}');
    print('API Response: ${response.body}');
  } catch (e) {
    print('API Error: $e');
  }
}
```
---

## 3. Deploy lên Azure Kubernetes Service (AKS)

### 🏗️ Kiến trúc Hệ thống trên Kubernetes

Trước khi triển khai, hãy hiểu rõ kiến trúc mục tiêu trên K8S:

#### Cổng vào (Entrypoint)
Một Nginx Reverse Proxy sẽ đóng vai trò là cổng vào duy nhất cho toàn bộ hệ thống. Nó được expose ra internet thông qua một Service loại LoadBalancer của Kubernetes. Nginx sẽ chịu trách nhiệm xử lý CORS, nén Gzip, và sau này là cả SSL/HTTPS.

#### Luồng Request
Internet → Nginx LoadBalancer (Public IP) → Pod Nginx → Service API Gateway → Pod API Gateway → Các Microservice nội bộ khác.

#### Giao tiếp nội bộ
Tất cả các service bên trong cluster (API Gateway, payment, ticket, v.v.) sẽ giao tiếp với nhau thông qua tên Service của Kubernetes (ví dụ: `http://payment-service:8083`), hoạt động như một hệ thống DNS nội bộ.

#### Dịch vụ Nền tảng (Backend Services)
Để đảm bảo tính ổn định và dễ quản lý, chúng ta không chạy database hay message broker trong cluster. Thay vào đó, chúng ta kết nối tới các dịch vụ bên ngoài (Managed Services) như NeonDB (PostgreSQL), Railway (MySQL), Azure Cache for Redis, và Confluent Cloud (Kafka).

#### Quản lý Cấu hình

- **ConfigMap**: Dùng để lưu trữ các cấu hình không nhạy cảm (tên host, tên database, URL dịch vụ).
- **Secret**: Dùng để lưu trữ TOÀN BỘ thông tin nhạy cảm (mật khẩu, API keys, token, connection strings chứa mật khẩu).

### 📋 Các giai đoạn chuyển đổi chính

#### Giai đoạn 1: Chuẩn bị Hạ tầng & Công cụ

##### Đăng nhập và Cấu hình Azure:

```bash
# Đăng nhập vào Azure
az login --use-device-code

# Đăng ký các nhà cung cấp tài nguyên (chỉ chạy một lần cho mỗi subscription)
az provider register --namespace Microsoft.ContainerRegistry
az provider register --namespace Microsoft.ContainerService
```

##### Tạo các Tài nguyên cần thiết:

```bash
# Tạo Resource Group
az group create --name MyResourceGroup --location southeastasia

# Tạo Azure Container Registry (ACR) để lưu trữ images
az acr create --resource-group MyResourceGroup --name duancntt --sku Basic

# Đăng nhập Docker vào ACR
az acr login --name duancntt

# Tạo cụm Azure Kubernetes Service (AKS)
az aks create --resource-group MyResourceGroup --name MyAKSCluster --node-count 2 --node-vm-size Standard_B2s --generate-ssh-keys

# Kết nối kubectl tới cụm AKS
az aks get-credentials --resource-group MyResourceGroup --name MyAKSCluster

# Cấp quyền cho AKS để kéo image từ ACR
az aks update -n MyAKSCluster -g MyResourceGroup --attach-acr duancntt
```

###### Chuẩn bị các Dịch vụ Bên ngoài:

Hãy đảm bảo bạn đã tạo và có đầy đủ thông tin kết nối (host, port, user, password, database names) cho NeonDB (PostgreSQL) và Railway (MySQL).

Đặc biệt, hãy vào mục Firewall / IP Allowlist của các dịch vụ này và thêm địa chỉ IP đầu ra của cụm AKS vào danh sách cho phép để tránh lỗi timeout.

```bash
# Lệnh tìm IP đầu ra của AKS
NODE_RG=$(az aks show -g MyResourceGroup -n MyAKSCluster --query nodeResourceGroup -o tsv)
az network public-ip list -g $NODE_RG --query "[].ipAddress" -o tsv
```

#### Giai đoạn 2: Container hóa & Đẩy Image

##### Viết Dockerfile
Mỗi service cần có một Dockerfile được tối ưu hóa.

**Best Practice**: Luôn sử dụng một tag phiên bản duy nhất cho mỗi lần build (ví dụ: `v1.0`, `v1.1`, `v2.0`) để Kubernetes biết khi nào cần cập nhật image. Tránh dùng lại tag cũ như `:v1`.

##### Build & Push:

```bash
# Ví dụ cho một service
docker build -t duancntt.azurecr.io/payment-service:v1.0 .
docker push duancntt.azurecr.io/payment-service:v1.0
```

Lặp lại cho tất cả các service. Đối với Rasa Chatbot, bạn cần build 2 image riêng biệt: `rasa-server` và `rasa-action-server`.

#### Giai đoạn 3: Viết Manifests Kubernetes

Đây là nơi bạn "dạy" cho Kubernetes cách chạy ứng dụng của mình bằng các file `.yaml`.

##### Tạo Secret (Bước quan trọng nhất):
Tập hợp tất cả các thông tin nhạy cảm và tạo một Secret duy nhất.

```bash
kubectl create secret generic platform-secrets \
  --from-literal=POSTGRES_PASSWORD='YOUR_NEONDB_PASSWORD' \
  --from-literal=MYSQL_PASSWORD='YOUR_RAILWAY_MYSQL_PASSWORD' \
  --from-literal=REDIS_URL='redis://:YOUR_REDIS_ACCESS_KEY@...' \
  --from-literal=KAFKA_SASL_USER='YOUR_KAFKA_USER' \
  --from-literal=KAFKA_SASL_PASS='YOUR_KAFKA_PASSWORD' \
  --from-literal=JWT_SECRET_KEY='YOUR_STRONG_JWT_SECRET' \
  --from-literal=VNPAY_HASH_SECRET='YOUR_VNPAY_HASH_SECRET' \
  --from-literal=STRIPE_SECRET_KEY='sk_test_YOUR_STRIPE_SECRET' \
  --from-literal=STRIPE_WEBHOOK_SECRET='whsec_YOUR_STRIPE_WEBHOOK_SECRET' \
  --from-literal=SMTP_USERNAME='your-email@gmail.com' \
  --from-literal=SMTP_PASSWORD='your-gmail-app-password' \
  --from-literal=CLOUDINARY_URL='cloudinary://api_key:api_secret@cloud_name'
```

##### Tạo các file YAML:
Tổ chức code vào các file riêng lẻ (`00-platform-config.yaml`, `nginx-proxy.yaml`, `12-api-gateway.yaml`, v.v.). Các mẫu quan trọng bao gồm:

- **ConfigMap** (`00-platform-config.yaml`): Chứa các cấu hình không nhạy cảm như tên host DB, tên các topic Kafka, URL nội bộ của các service.

- **Nginx Proxy** (`nginx-proxy.yaml`): Bao gồm ConfigMap cho `nginx.conf`, Deployment cho Nginx, và một Service loại LoadBalancer để expose Nginx ra ngoài.

- **API Gateway** (`12-api-gateway.yaml`): Bao gồm Deployment và một Service loại ClusterIP để chỉ Nginx mới có thể gọi vào.

- **Các Backend Service**: Mỗi service sẽ có một file YAML chứa Deployment và Service loại ClusterIP. Deployment sẽ lấy cấu hình từ `platform-config` và `platform-secrets`.

- **Rasa Chatbot** (`11-chat-bot-service.yaml`): File này đặc biệt vì Deployment của nó sẽ chứa 2 container: `rasa-server` và `rasa-action-server`. Service sẽ chỉ expose port của rasa-server (5005).

#### Giai đoạn 4: Triển khai & Vận hành

#### Chuẩn bị
Đảm bảo bạn đang ở trong thư mục chứa tất cả các file `.yaml`.

##### Triển khai:

```bash
# Áp dụng tất cả cấu hình
kubectl apply -f .
```

##### Lấy IP Công khai:

```bash
# Chờ IP xuất hiện ở cột EXTERNAL-IP
kubectl get service nginx-loadbalancer -w
```

#### Kiểm tra
Dùng IP của `nginx-loadbalancer` để gửi request API từ Postman hoặc trình duyệt.

#### Giai đoạn 5: Gỡ lỗi Thường gặp

##### ImagePullBackOff
Sai tên/tag image hoặc chưa cấp quyền cho AKS.

**Fix**: Kiểm tra lại tên image trong file YAML. Chạy lại `az aks update --attach-acr ....`

##### CreateContainerConfigError
Thiếu ConfigMap, Secret hoặc một key trong đó.

**Fix**: Dùng `kubectl describe pod <tên_pod>` để xem Events ở cuối, nó sẽ chỉ rõ key nào đang bị thiếu. Tạo lại Secret hoặc sửa ConfigMap cho đúng.

##### CrashLoopBackOff
Lỗi ứng dụng bên trong container.

**Fix**: Dùng `kubectl logs <tên_pod>` để xem log ứng dụng.

##### password authentication failed
Sai mật khẩu, cập nhật lại Secret.

##### database "..." does not exist
Chưa tạo database đó trên server. Hãy kết nối tới DB server và chạy `CREATE DATABASE ...;`.

##### context canceled khi kết nối ra ngoài
Lỗi Firewall/IP Allowlist.

##### 504 Gateway Time-out (từ Nginx)
Nginx không kết nối được tới service phía sau (API Gateway).

**Fix**: Kiểm tra xem Service của API Gateway (`api-gateway-service`) có tồn tại và có Endpoints không bằng lệnh `kubectl describe service api-gateway-service`.

### 🔄 Cập nhật Frontend sau khi Deploy

Sau khi triển khai thành công lên AKS và có được IP công khai, bạn cần cập nhật cấu hình frontend:

#### Cập nhật React App:

```env
# Cập nhật .env.production
REACT_APP_API_BASE_URL=http://YOUR_AKS_PUBLIC_IP/api
REACT_APP_WEBSOCKET_URL=ws://YOUR_AKS_PUBLIC_IP/ws
REACT_APP_ENVIRONMENT=production
```

#### Cập nhật Flutter App:

```dart
// Cập nhật lib/config/api_config.dart
class ApiConfig {
  static const String PROD_BASE_URL = 'http://YOUR_AKS_PUBLIC_IP/api';
  static const String PROD_WEBSOCKET_URL = 'ws://YOUR_AKS_PUBLIC_IP/ws';
  
  // Chuyển sang production mode
  static const bool IS_DEVELOPMENT = false;
}
```

#### Build và Deploy Frontend:

```bash
# Build React for production
cd frontend-react
npm run build

# Deploy lên hosting service (Netlify, Vercel, etc.)
# hoặc serve từ Nginx static files

# Build Flutter APK mới
cd mobile-flutter
flutter build apk --release
```

Bằng cách tuân theo quy trình trên, bạn đã xây dựng thành công một hệ thống microservices mạnh mẽ và chuyên nghiệp trên nền tảng Kubernetes.

---

## 4. ETL Pipeline với Apache Airflow

### 📊 Tổng quan

Sau khi hệ thống microservices đã hoạt động, một yêu cầu phổ biến là trích xuất, chuyển đổi và tải dữ liệu (ETL) để phục vụ cho việc phân tích. Phần này mô tả các bước cần thiết để triển khai và thực thi hai quy trình (pipeline) ETL được quản lý bởi Apache Airflow, sử dụng dữ liệu được tạo ra bởi hệ thống microservices của chúng ta.

- **Pipeline 1** (`kafka_to_bigquery_pipeline_final`): Lắng nghe và tiêu thụ dữ liệu sự kiện tìm kiếm chuyến đi (trip search) từ một topic trên Kafka, sau đó tải dữ liệu này lên Google Cloud Storage (GCS) và cuối cùng nạp vào BigQuery.

- **Pipeline 2** (`etl_postgres_gcs_bigquery_daily_v2`): Trích xuất dữ liệu định kỳ hàng ngày từ hai bảng (ticket và invoices) trong hai cơ sở dữ liệu PostgreSQL khác nhau, lưu trữ dưới dạng file CSV lên GCS, và sau đó nạp vào các bảng tương ứng trong BigQuery.

### 6.1. Yêu cầu và Điều kiện tiên quyết

Trước khi bắt đầu, hãy đảm bảo môi trường Airflow của bạn đã đáp ứng các yêu cầu sau:

#### a. Cài đặt các thư viện Python cần thiết:
Bạn cần cài đặt các Airflow providers và thư viện kafka-python. Cách tốt nhất là thêm các gói này vào file `requirements.txt` của môi trường Airflow.

```
apache-airflow-providers-google
apache-airflow-providers-postgres
kafka-python-ng  # Hoặc kafka-python, tùy thuộc vào môi trường của bạn
pandas
```

Sau đó, chạy lệnh cài đặt:
```bash
pip install -r requirements.txt
```

#### b. Truy cập và Thông tin xác thực:

**Kafka (Confluent Cloud/Redpanda):**
- Bootstrap servers (host)
- Username (login)
- Password

**PostgreSQL:**
- Connection string hoặc thông tin chi tiết (host, port, username, password, database name) cho cả hai database `ticket_db` và `invoice_db`.

**Google Cloud Platform (GCP):**
- Một Service Account của Google Cloud có quyền truy cập vào:
  - **Google Cloud Storage (GCS)**: Quyền đọc/ghi (`roles/storage.objectAdmin`) trên bucket được chỉ định.
  - **BigQuery**: Quyền tạo và ghi dữ liệu vào bảng (`roles/bigquery.dataEditor`) và quyền tạo job (`roles/bigquery.jobUser`) trong project của bạn.
- File JSON key của Service Account này.

### 6.2. Cấu hình trong Airflow UI

Sau khi đã có đủ thông tin xác thực, bạn cần tạo các "Connections" trong giao diện người dùng của Airflow. Truy cập vào **Admin -> Connections**.

#### a. Cấu hình cho Kafka

- **Connection ID**: `kafka_redpanda_cloud`
- **Connection Type**: Kafka
- **Host**: Dán danh sách các bootstrap servers của bạn (ví dụ: `your-cluster.redpanda.cloud:9092`).
- **Login**: Username để truy cập Kafka.
- **Password**: Password để truy cập Kafka.
- **Extra**: (Không bắt buộc, nhưng khuyến khích)

```json
{
    "security_protocol": "SASL_SSL",
    "sasl_mechanism": "SCRAM-SHA-256"
}
```

#### b. Cấu hình cho Google Cloud Platform

- **Connection ID**: `google_cloud_default`
- **Connection Type**: Google Cloud
- **Keyfile JSON**: Dán toàn bộ nội dung của file JSON key từ Service Account của bạn vào đây.

#### c. Cấu hình cho PostgreSQL (Cơ sở dữ liệu ticket)

- **Connection ID**: `postgres_ticket_db`
- **Connection Type**: Postgres
- **Host**: Địa chỉ host của database ticket.
- **Schema**: Tên database (ví dụ: `ticketdb`).
- **Login**: Username của database.
- **Password**: Password của database.
- **Port**: Cổng kết nối (ví dụ: `5432`).

#### d. Cấu hình cho PostgreSQL (Cơ sở dữ liệu invoice)

- **Connection ID**: `postgres_invoice_db`
- **Connection Type**: Postgres
- **Host**: Địa chỉ host của database invoice.
- **Schema**: Tên database (ví dụ: `invoicedb`).
- **Login**: Username của database.
- **Password**: Password của database.
- **Port**: Cổng kết nối (ví dụ: `5432`).

### 6.3. Kiểm tra và Tùy chỉnh các biến trong file DAG

Trước khi đưa các file Python vào thư mục dags, hãy mở chúng ra và kiểm tra lại các biến cấu hình ở đầu mỗi file để đảm bảo chúng khớp với môi trường của bạn.

#### Trong file `search.py`:

```python
# =============================================================================
# CẤU HÌNH
# =============================================================================
KAFKA_CONN_ID = "kafka_redpanda_cloud"
KAFKA_TOPIC = "trip_search" # <-- Tên topic Kafka bạn muốn lắng nghe

GCP_CONN_ID = "google_cloud_default"
GCS_BUCKET_NAME = "asia-southeast1-invoice-dua-e3630dcb-bucket" # <-- Tên GCS bucket của bạn

BIGQUERY_PROJECT_ID = "dacntt-dfabb" # <-- Project ID của bạn
BIGQUERY_DATASET_NAME = "duancntt"  # <-- Tên dataset trong BigQuery
BIGQUERY_TABLE_NAME = "trip_searches_raw" # <-- Tên bảng đích
```

#### Trong file `invoice.py`:

```python
# =============================================================================
# Cấu hình chung
# =============================================================================
GCP_CONN_ID = "google_cloud_default"
GCS_BUCKET_NAME = "asia-southeast1-invoice-dua-e3630dcb-bucket" # <-- Tên GCS bucket của bạn
POSTGRES_TICKET_CONN_ID = "postgres_ticket_db"
POSTGRES_INVOICE_CONN_ID = "postgres_invoice_db"

# Cấu hình cho BigQuery
BIGQUERY_PROJECT_ID = "dacntt-dfabb"        # <-- Project ID của bạn
BIGQUERY_DATASET_NAME = "duancntt"           # <-- Tên dataset trong BigQuery
```

### 6.4. Triển khai và Kích hoạt DAGs

#### Sao chép file DAG:
Đặt cả hai file `search.py` và `invoice.py` vào thư mục `dags` của môi trường Apache Airflow.

#### Kích hoạt DAGs trong Airflow UI:
Mở giao diện người dùng Airflow. Sau một vài phút, bạn sẽ thấy hai DAG mới xuất hiện với tên:

- `kafka_to_bigquery_pipeline_final`
- `etl_postgres_gcs_bigquery_daily_v2`

Mặc định, các DAG mới sẽ ở trạng thái "paused". Hãy bật chúng lên bằng cách gạt nút toggle bên cạnh tên DAG.

### 6.5. Lịch trình và Thực thi

#### `kafka_to_bigquery_pipeline_final`:

- **Lịch chạy**: Được cấu hình để chạy mỗi 2 giờ (`schedule_interval='0 */2 * * *'`).
- **Hoạt động**: Mỗi lần chạy, DAG sẽ kết nối đến Kafka, thu thập các message trong vòng 60 giây (dựa trên `consumer_timeout_ms`), đẩy chúng thành một file JSON xuống GCS, và sau đó nạp vào BigQuery. Nếu không có message nào, task GCS-to-BigQuery sẽ được bỏ qua.

#### `etl_postgres_gcs_bigquery_daily_v2`:

- **Lịch chạy**: Được cấu hình để chạy vào lúc 20:00 (8 PM) mỗi ngày (`schedule_interval='0 20 * * *'`).
- **Hoạt động**: DAG sẽ trích xuất dữ liệu từ ngày hôm trước (dựa trên `data_interval_start` và `data_interval_end`) từ hai bảng `ticket` và `invoices` trong PostgreSQL, lưu thành hai file CSV riêng biệt trên GCS, và cuối cùng nạp chúng vào hai bảng tương ứng trong BigQuery.

Để thực thi ngay lập tức mà không cần chờ lịch trình, bạn có thể nhấn vào nút "Play" (Trigger DAG) bên phải tên mỗi DAG.