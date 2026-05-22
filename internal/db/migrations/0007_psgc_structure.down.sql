BEGIN;

ALTER TABLE address_line DROP CONSTRAINT IF EXISTS fk_address_barangay;
ALTER TABLE barangays DROP CONSTRAINT IF EXISTS fk_barangays_municipal_or_city;
ALTER TABLE municipal_or_city DROP CONSTRAINT IF EXISTS fk_municipal_or_city_province_or_huc;
ALTER TABLE province_or_huc DROP CONSTRAINT IF EXISTS fk_province_or_huc_region;

ALTER TABLE address_line DROP COLUMN IF EXISTS barangay_psgc;
ALTER TABLE address_line ADD COLUMN barangay_id bigint NOT NULL;

ALTER TABLE barangays DROP CONSTRAINT IF EXISTS barangays_pkey;
ALTER TABLE barangays DROP COLUMN IF EXISTS city_mun_psgc;
ALTER TABLE barangays RENAME COLUMN barangay_psgc TO barangay_id;
ALTER TABLE barangays ALTER COLUMN barangay_id TYPE bigint USING trim(barangay_id)::bigint;

CREATE SEQUENCE IF NOT EXISTS barangay_barangay_id_seq;
ALTER TABLE barangays ALTER COLUMN barangay_id SET DEFAULT nextval('barangay_barangay_id_seq'::regclass);
ALTER TABLE barangays ADD PRIMARY KEY (barangay_id);
ALTER TABLE barangays RENAME TO barangay;

ALTER TABLE municipal_or_city DROP CONSTRAINT IF EXISTS municipal_or_city_pkey;
ALTER TABLE municipal_or_city DROP COLUMN IF EXISTS province_or_huc_psgc;
ALTER TABLE municipal_or_city DROP COLUMN IF EXISTS type;
ALTER TABLE municipal_or_city RENAME COLUMN city_mun_name TO municipal_or_city_name;
ALTER TABLE municipal_or_city RENAME COLUMN city_mun_psgc TO municipal_or_city_id;
ALTER TABLE municipal_or_city ALTER COLUMN municipal_or_city_id TYPE bigint USING trim(municipal_or_city_id)::bigint;
ALTER TABLE municipal_or_city ADD COLUMN province_or_huc_id bigint NOT NULL;

CREATE SEQUENCE IF NOT EXISTS municipal_or_city_municipal_or_city_id_seq;
ALTER TABLE municipal_or_city ALTER COLUMN municipal_or_city_id
	SET DEFAULT nextval('municipal_or_city_municipal_or_city_id_seq'::regclass);
ALTER TABLE municipal_or_city ADD PRIMARY KEY (municipal_or_city_id);

ALTER TABLE province_or_huc DROP CONSTRAINT IF EXISTS province_or_huc_pkey;
ALTER TABLE province_or_huc DROP COLUMN IF EXISTS region_psgc;
ALTER TABLE province_or_huc RENAME COLUMN province_or_huc_psgc TO province_or_huc_id;
ALTER TABLE province_or_huc ALTER COLUMN province_or_huc_id TYPE bigint USING trim(province_or_huc_id)::bigint;
ALTER TABLE province_or_huc ADD COLUMN region_id bigint NOT NULL;

CREATE SEQUENCE IF NOT EXISTS province_or_huc_province_or_huc_id_seq;
ALTER TABLE province_or_huc ALTER COLUMN province_or_huc_id
	SET DEFAULT nextval('province_or_huc_province_or_huc_id_seq'::regclass);
ALTER TABLE province_or_huc ADD PRIMARY KEY (province_or_huc_id);

ALTER TABLE region DROP CONSTRAINT IF EXISTS region_pkey;
ALTER TABLE region RENAME COLUMN region_psgc TO region_id;
ALTER TABLE region ALTER COLUMN region_id TYPE bigint USING trim(region_id)::bigint;

CREATE SEQUENCE IF NOT EXISTS region_region_id_seq;
ALTER TABLE region ALTER COLUMN region_id SET DEFAULT nextval('region_region_id_seq'::regclass);
ALTER TABLE region ADD PRIMARY KEY (region_id);

ALTER TABLE province_or_huc
	ADD CONSTRAINT fk_province_or_huc_region
	FOREIGN KEY (region_id) REFERENCES region(region_id);
ALTER TABLE municipal_or_city
	ADD CONSTRAINT fk_municipal_or_city_province_or_huc
	FOREIGN KEY (province_or_huc_id) REFERENCES province_or_huc(province_or_huc_id);
ALTER TABLE barangay
	ADD CONSTRAINT fk_barangay_municipal_or_city
	FOREIGN KEY (municipal_or_city_id) REFERENCES municipal_or_city(municipal_or_city_id);
ALTER TABLE address_line
	ADD CONSTRAINT fk_address_barangay
	FOREIGN KEY (barangay_id) REFERENCES barangay(barangay_id);

COMMIT;
