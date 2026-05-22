BEGIN;

ALTER TABLE address_line DROP CONSTRAINT IF EXISTS fk_address_barangay;
ALTER TABLE barangay DROP CONSTRAINT IF EXISTS fk_barangay_municipal_or_city;
ALTER TABLE municipal_or_city DROP CONSTRAINT IF EXISTS fk_municipal_or_city_province_or_huc;
ALTER TABLE province_or_huc DROP CONSTRAINT IF EXISTS fk_province_or_huc_region;

ALTER TABLE address_line DROP COLUMN IF EXISTS barangay_id;
ALTER TABLE barangay DROP COLUMN IF EXISTS municipal_or_city_id;
ALTER TABLE municipal_or_city DROP COLUMN IF EXISTS province_or_huc_id;
ALTER TABLE province_or_huc DROP COLUMN IF EXISTS region_id;

ALTER TABLE region DROP CONSTRAINT IF EXISTS region_pkey;
ALTER TABLE region RENAME COLUMN region_id TO region_psgc;
ALTER TABLE region ALTER COLUMN region_psgc TYPE CHAR(2);
ALTER TABLE region ALTER COLUMN region_psgc DROP DEFAULT;
ALTER TABLE region ADD PRIMARY KEY (region_psgc);

ALTER TABLE province_or_huc DROP CONSTRAINT IF EXISTS province_or_huc_pkey;
ALTER TABLE province_or_huc RENAME COLUMN province_or_huc_id TO province_or_huc_psgc;
ALTER TABLE province_or_huc ALTER COLUMN province_or_huc_psgc TYPE CHAR(5);
ALTER TABLE province_or_huc ALTER COLUMN province_or_huc_psgc DROP DEFAULT;
ALTER TABLE province_or_huc ADD COLUMN region_psgc CHAR(2) NOT NULL;
ALTER TABLE province_or_huc ADD PRIMARY KEY (province_or_huc_psgc);
ALTER TABLE province_or_huc
	ADD CONSTRAINT fk_province_or_huc_region
	FOREIGN KEY (region_psgc) REFERENCES region(region_psgc);

ALTER TABLE municipal_or_city DROP CONSTRAINT IF EXISTS municipal_or_city_pkey;
ALTER TABLE municipal_or_city RENAME COLUMN municipal_or_city_id TO city_mun_psgc;
ALTER TABLE municipal_or_city ALTER COLUMN city_mun_psgc TYPE CHAR(7);
ALTER TABLE municipal_or_city ALTER COLUMN city_mun_psgc DROP DEFAULT;
ALTER TABLE municipal_or_city RENAME COLUMN municipal_or_city_name TO city_mun_name;
ALTER TABLE municipal_or_city ADD COLUMN province_or_huc_psgc CHAR(5);
ALTER TABLE municipal_or_city ADD COLUMN type VARCHAR(20) CHECK (type IN ('City', 'Municipality'));
ALTER TABLE municipal_or_city ADD PRIMARY KEY (city_mun_psgc);
ALTER TABLE municipal_or_city
	ADD CONSTRAINT fk_municipal_or_city_province_or_huc
	FOREIGN KEY (province_or_huc_psgc) REFERENCES province_or_huc(province_or_huc_psgc);

ALTER TABLE barangay RENAME TO barangays;
ALTER TABLE barangays DROP CONSTRAINT IF EXISTS barangay_pkey;
ALTER TABLE barangays RENAME COLUMN barangay_id TO barangay_psgc;
ALTER TABLE barangays ALTER COLUMN barangay_psgc TYPE CHAR(10);
ALTER TABLE barangays ALTER COLUMN barangay_psgc DROP DEFAULT;
ALTER TABLE barangays ADD COLUMN city_mun_psgc CHAR(7) NOT NULL;
ALTER TABLE barangays ADD PRIMARY KEY (barangay_psgc);
ALTER TABLE barangays
	ADD CONSTRAINT fk_barangays_municipal_or_city
	FOREIGN KEY (city_mun_psgc) REFERENCES municipal_or_city(city_mun_psgc);

ALTER TABLE address_line ADD COLUMN barangay_psgc CHAR(10) NOT NULL;
ALTER TABLE address_line
	ADD CONSTRAINT fk_address_barangay
	FOREIGN KEY (barangay_psgc) REFERENCES barangays(barangay_psgc);

DROP SEQUENCE IF EXISTS region_region_id_seq;
DROP SEQUENCE IF EXISTS province_or_huc_province_or_huc_id_seq;
DROP SEQUENCE IF EXISTS municipal_or_city_municipal_or_city_id_seq;
DROP SEQUENCE IF EXISTS barangay_barangay_id_seq;

COMMIT;
