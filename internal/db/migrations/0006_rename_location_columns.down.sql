BEGIN;

ALTER SEQUENCE municipal_or_city_municipal_or_city_id_seq RENAME TO district_district_id_seq;
ALTER TABLE municipal_or_city ALTER COLUMN municipal_or_city_id SET DEFAULT nextval('district_district_id_seq');
ALTER SEQUENCE district_district_id_seq OWNED BY municipal_or_city.municipal_or_city_id;

ALTER SEQUENCE province_or_huc_province_or_huc_id_seq RENAME TO city_city_id_seq;
ALTER TABLE province_or_huc ALTER COLUMN province_or_huc_id SET DEFAULT nextval('city_city_id_seq');
ALTER SEQUENCE city_city_id_seq OWNED BY province_or_huc.province_or_huc_id;

ALTER SEQUENCE region_region_id_seq RENAME TO province_province_id_seq;
ALTER TABLE region ALTER COLUMN region_id SET DEFAULT nextval('province_province_id_seq');
ALTER SEQUENCE province_province_id_seq OWNED BY region.region_id;

ALTER TABLE barangay RENAME COLUMN municipal_or_city_id TO district_id;

ALTER TABLE municipal_or_city RENAME COLUMN province_or_huc_id TO city_id;
ALTER TABLE municipal_or_city RENAME COLUMN municipal_or_city_name TO district_name;
ALTER TABLE municipal_or_city RENAME COLUMN municipal_or_city_id TO district_id;

ALTER TABLE province_or_huc RENAME COLUMN region_id TO province_id;
ALTER TABLE province_or_huc RENAME COLUMN province_or_huc_name TO city_name;
ALTER TABLE province_or_huc RENAME COLUMN province_or_huc_id TO city_id;

ALTER TABLE region RENAME COLUMN region_name TO province_name;
ALTER TABLE region RENAME COLUMN region_id TO province_id;

COMMIT;
