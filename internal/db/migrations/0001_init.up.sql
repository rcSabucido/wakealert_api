CREATE TABLE mobile_user (
    mobile_user_id BIGSERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE receiver_user (
    receiver_user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE relationship_type (    
    relationship_type_id BIGSERIAL PRIMARY KEY,
    relationship_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE pregnancy_status (
    pregnancy_status_id BIGSERIAL PRIMARY KEY,
    pregnancy_status_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE organ_donor_status (
    donor_status_id BIGSERIAL PRIMARY KEY,
    donor_status_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE blood_type (
    blood_type_id BIGSERIAL PRIMARY KEY,
    blood_type_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE medical_information (
    medical_info_id BIGSERIAL PRIMARY KEY,
    allergies VARCHAR(500),
    medication VARCHAR(500),
    medical_notes VARCHAR(500),
    pregnancy_status_id BIGINT NOT NULL,
    donor_status_id BIGINT NOT NULL,
    blood_type_id BIGINT NOT NULL,
    last_diagnosis_date DATE,
    last_diagnosis_hospital_name VARCHAR(100),
    CONSTRAINT fk_medinfo_pregnancy_status FOREIGN KEY (pregnancy_status_id) REFERENCES pregnancy_status(pregnancy_status_id),
    CONSTRAINT fk_medinfo_donor_status FOREIGN KEY (donor_status_id) REFERENCES organ_donor_status(donor_status_id),
    CONSTRAINT fk_medinfo_blood_type FOREIGN KEY (blood_type_id) REFERENCES blood_type(blood_type_id)
);

CREATE TABLE province (
    province_id BIGSERIAL PRIMARY KEY,
    province_name VARCHAR(100) NOT NULL
);

CREATE TABLE city (
    city_id BIGSERIAL PRIMARY KEY,
    city_name VARCHAR(100) NOT NULL,
    province_id BIGINT NOT NULL,
    CONSTRAINT fk_city_province FOREIGN KEY (province_id) REFERENCES province(province_id)
);

CREATE TABLE district (
    district_id BIGSERIAL PRIMARY KEY,
    district_name VARCHAR(100) NOT NULL,
    city_id BIGINT NOT NULL,
    CONSTRAINT fk_district_city FOREIGN KEY (city_id) REFERENCES city(city_id)
);

CREATE TABLE barangay (
    barangay_id BIGSERIAL PRIMARY KEY,
    barangay_name VARCHAR(100) NOT NULL,
    district_id BIGINT NOT NULL,
    CONSTRAINT fk_barangay_district FOREIGN KEY (district_id) REFERENCES district(district_id)
);

CREATE TABLE address_line (
    address_id BIGSERIAL PRIMARY KEY,
    barangay_id BIGINT NOT NULL,
    address_line VARCHAR(100) NOT NULL,
    CONSTRAINT fk_address_barangay FOREIGN KEY (barangay_id) REFERENCES barangay(barangay_id)
);

CREATE TABLE victim (
    victim_id BIGSERIAL PRIMARY KEY,
    mobile_user_id BIGINT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    birth_date DATE,
    address_id BIGINT NOT NULL,
    medical_info_id BIGINT NOT NULL UNIQUE,
    CONSTRAINT fk_victim_mobile_user FOREIGN KEY (mobile_user_id) REFERENCES mobile_user(mobile_user_id),
    CONSTRAINT fk_victim_address FOREIGN KEY (address_id) REFERENCES address_line(address_id),
    CONSTRAINT fk_victim_medical_info FOREIGN KEY (medical_info_id) REFERENCES medical_information(medical_info_id)
);

CREATE TABLE contact (
    contact_id BIGSERIAL PRIMARY KEY,
    mobile_user_id BIGINT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    relationship_type_id BIGINT NOT NULL,
    is_primary BOOLEAN NOT NULL,
    CONSTRAINT fk_contact_mobile_user FOREIGN KEY (mobile_user_id) REFERENCES mobile_user(mobile_user_id),
    CONSTRAINT fk_contact_relationship_type FOREIGN KEY (relationship_type_id) REFERENCES relationship_type(relationship_type_id)
);

CREATE TABLE alert (
    alert_id BIGSERIAL PRIMARY KEY,
    latitude NUMERIC(9, 6) NOT NULL,
    longitude NUMERIC(10, 6) NOT NULL,
    alert_time TIMESTAMP NOT NULL DEFAULT NOW(),
    victim_id BIGINT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_alert_victim FOREIGN KEY (victim_id) REFERENCES victim(victim_id)
);

CREATE TABLE medical_history_entry (
    history_entry_id BIGSERIAL PRIMARY KEY,
    medical_info_id BIGINT NOT NULL,
    diagnosis VARCHAR(500) NOT NULL,
    is_recent BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_history_medical_info FOREIGN KEY (medical_info_id) REFERENCES medical_information(medical_info_id)
);