# -------------------------------------------------------
# 1) Build stage - PHẢI DÙNG JDK ĐỂ BUILD
# -------------------------------------------------------
    FROM eclipse-temurin:17-jdk AS build

    # Install Maven (Debian-based, so use apt)
    RUN apt-get update && apt-get install -y maven && rm -rf /var/lib/apt/lists/*
    
    WORKDIR /app
    
    # Copy pom.xml trước để cache dependencies
    COPY pom.xml .
    RUN mvn dependency:go-offline -B
    
    # Copy source code sau
    COPY src ./src
    RUN mvn clean package -DskipTests -B
    
    # -------------------------------------------------------
    # 2) Runtime stage  
    # -------------------------------------------------------
    FROM eclipse-temurin:17-jre
    
    WORKDIR /app
    
    # Tạo user non-root với Debian commands
    RUN groupadd -g 1001 appgroup && useradd -u 1001 -g appgroup -m appuser
    
    # Copy jar file
    COPY --from=build /app/target/*.jar user_service.jar
    
    # Chown cho user
    RUN chown appuser:appgroup user_service.jar
    
    USER appuser
    
    EXPOSE 8081
    
    ENTRYPOINT ["java", "-jar", "user_service.jar"]