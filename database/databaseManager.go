package database

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS User (
			id INT AUTO_INCREMENT PRIMARY KEY,
			firstname VARCHAR (20) NOT NULL,
			lastname VARCHAR (20) NOT NULL,
			phone INT NOT NULL,
			email VARCHAR (20) NOT NULL,
			password VARCHAR (255) NOT NULL,
			role ENUM ('admin', 'operator') NOT NULL,
			status ENUM ('work', 'leave') DEFAULT 'work',
			created_at DATETIME
			)`,
		`CREATE TABLE IF NOT EXISTS Vehicle(
			id INT AUTO_INCREMENT PRIMARY KEY,
			register_number VARCHAR (10) NOT NULL,
			name VARCHAR (20) NOT NULL,
			model VARCHAR (20) NOT NULL,
			type ENUM('truck', 'van', 'car', 'drone') NOT NULL,
			type_charge ENUM('electric', 'fuel') NOT NULL,
			currrent_charge FLOAT,
			charge_capacity FLOAT,
			current_distance FLOAT,
			current_position VARCHAR(255),
			status ENUM('available', 'use', 'maintenance') DEFAULT 'available',
			connection_key VARCHAR (20),
			created_at DATETIME,
			created_by INT,
			FOREIGN KEY (created_by) REFERENCES User (id)
			
		)`,
		`CREATE TABLE IF NOT EXISTS Driver(
			id INT AUTO_INCREMENT PRIMARY KEY,
			firstname VARCHAR (20) NOT NULL,
			lastname VARCHAR (20) NOT NULL,
			birthday DATE,
			phone INT,
			email VARCHAR (20) NOT NULL,
			password VARCHAR (255) NOT NULL,
			class ENUM('a','b', 'c', 'd'),
			status ENUM('work', 'leave') DEFAULT 'work',
			created_at DATETIME,
			created_by INT,
			FOREIGN KEY (created_by) REFERENCES User (id)
		)
		`,
		`CREATE TABLE IF NOT EXISTS Route(
			id INT AUTO_INCREMENT PRIMARY KEY,
			status ENUM ('progress', 'canceled', 'completed') NOT NULL,
			departure_date DATE NOT NULL,
			arrival_date DATE NOT NULL,
			driver_id INT,
			vehicle_id INT,
			created_at DATETIME,
			created_by INT,
			FOREIGN KEY (driver_id) REFERENCES Driver (id),
			FOREIGN KEY (vehicle_id) REFERENCES Vehicle (id),
			FOREIGN KEY (created_by) REFERENCES User (id)
			)`,
		`CREATE TABLE IF NOT EXISTS Route_logs(
			id INT AUTO_INCREMENT PRIMARY KEY,
			vehicle_id INT,
			lat FLOAT NOT NULL,
			lng FLOAT NOT NULL,
			datetime DATETIME,
			FOREIGN KEY (vehicle_id) REFERENCES Vehicle (id)
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
