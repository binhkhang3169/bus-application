# Giai đoạn 1: Builder
FROM golang:1.24-alpine AS builder

# Cài đặt tzdata để hỗ trợ timezone và các chứng chỉ CA
RUN apk add --no-cache tzdata ca-certificates

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép go.mod và go.sum trước để tận dụng Docker cache
COPY go.mod go.sum ./

# Tải dependencies
RUN go mod download

# Sao chép toàn bộ mã nguồn
COPY . .

# Build ứng dụng Go thành một file thực thi tĩnh
# Đặt tên file thực thi là 'emailservice'
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /app/emailservice main.go

# Giai đoạn 2: Final image từ scratch
FROM scratch

# Sao chép timezone data từ builder stage
# Điều này cần thiết để ENV TZ=Asia/Ho_Chi_Minh hoạt động chính xác trong scratch image
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Thiết lập timezone Việt Nam
ENV TZ=Asia/Ho_Chi_Minh

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép file thực thi đã được build từ builder stage
COPY --from=builder /app/emailservice /app/emailservice
COPY --from=builder /app/templates /app/templates

# Expose port mà ứng dụng sẽ lắng nghe
EXPOSE 8085

# Entrypoint để chạy ứng dụng
# File thực thi đã được đặt tên là 'emailservice'
ENTRYPOINT ["/app/emailservice"]
