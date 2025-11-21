-- Create database
CREATE DATABASE IF NOT EXISTS weather_label_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE weather_label_db;

-- Monitoring stations table
CREATE TABLE IF NOT EXISTS stations (
    id VARCHAR(255) KEY,
    name VARCHAR(255) NOT NULL,
    longitude DECIMAL(10, 5) NOT NULL,
    latitude DECIMAL(10, 5) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_coordinates (longitude, latitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Images table
CREATE TABLE IF NOT EXISTS images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    filepath VARCHAR(512) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    annotated BOOLEAN DEFAULT FALSE,
    is_standard BOOLEAN DEFAULT NULL COMMENT 'NULL=未处理, TRUE=标准图片(有时间和地点), FALSE=非标准图片',
    ocr_time VARCHAR(255) DEFAULT NULL COMMENT 'OCR识别的时间',
    ocr_location VARCHAR(255) DEFAULT NULL COMMENT 'OCR识别的地点',
    INDEX idx_annotated (annotated),
    INDEX idx_is_standard (is_standard)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Annotations table
CREATE TABLE IF NOT EXISTS annotations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    image_id INT NOT NULL,
    category ENUM('积涝', '大雾', '结冰') NOT NULL,
    severity ENUM('无', '轻度', '中度', '重度') NOT NULL,
    observation_time DATETIME NOT NULL,
    location VARCHAR(255) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    latitude DECIMAL(10, 7) NOT NULL,
    station_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (station_id) REFERENCES stations(id),
    UNIQUE KEY unique_image_annotation (image_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;