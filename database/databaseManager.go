package database

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS USER (
			id INT AUTO_INCREMENT PRIMARY KEY,
			firstname VARCHAR (20) NOT NULL,
			lastname VARCHAR (20) NOT NULL,		
			phone VARCHAR(10),
			email VARCHAR (20) NOT NULL,
			password VARCHAR (255) NOT NULL,
			role ENUM ('admin', 'operator') NOT NULL,
			status ENUM ('work', 'leave') DEFAULT 'work',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)`,
		`CREATE TABLE IF NOT EXISTS VEHICLE(
			id INT AUTO_INCREMENT PRIMARY KEY,
			register_number VARCHAR (10) NOT NULL,
			name VARCHAR (30) NOT NULL,
			model VARCHAR (30) NOT NULL,
			type ENUM('truck', 'van', 'car', 'drone') NOT NULL,
			type_charge ENUM('electric', 'fuel') NOT NULL,
			current_charge FLOAT,
			charge_capacity FLOAT,
			current_distance FLOAT,
			current_position VARCHAR(255),
			status ENUM('available', 'use', 'maintenance') DEFAULT 'available',
			connection_key VARCHAR (20),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			created_by INT,
			FOREIGN KEY (created_by) REFERENCES USER(id)
		)`,
		`CREATE TABLE IF NOT EXISTS DRIVER(
			id INT AUTO_INCREMENT PRIMARY KEY,
			register_number VARCHAR (10),
			firstname VARCHAR (20) NOT NULL,
			lastname VARCHAR (20) NOT NULL,
			date_of_birth DATE,
			phone VARCHAR(10),
			email VARCHAR (20) NOT NULL,
			password VARCHAR (255) NOT NULL,
			class ENUM('a','b', 'c', 'd'),
			status ENUM('work', 'leave') DEFAULT 'work',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			created_by INT,
			FOREIGN KEY (created_by) REFERENCES USER(id)
		)
		`,
		`CREATE TABLE IF NOT EXISTS ROUTE(
			id INT AUTO_INCREMENT PRIMARY KEY,
			status ENUM ('progress', 'canceled', 'completed') DEFAULT 'progress',
			departure_date DATE NOT NULL,
			arrival_date DATE NOT NULL,
			lat_start VARCHAR(255),
			lng_start VARCHAR(255),
			lat_end	VARCHAR(255),
			lng_end VARCHAR(255),
			departure_city VARCHAR(30),
			arrival_city VARCHAR(30),
			driver_id INT,
			vehicle_id INT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			created_by INT,
			FOREIGN KEY (driver_id) REFERENCES DRIVER (id),
			FOREIGN KEY (vehicle_id) REFERENCES VEHICLE (id),
			FOREIGN KEY (created_by) REFERENCES USER (id)
			)`,
		`CREATE TABLE IF NOT EXISTS ROUTE_LOGS(
			id INT AUTO_INCREMENT PRIMARY KEY,
			vehicle_id INT,
			lat FLOAT NOT NULL,
			lng FLOAT NOT NULL,
			datetime DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (vehicle_id) REFERENCES VEHICLE (id)
		)`,
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil

}
