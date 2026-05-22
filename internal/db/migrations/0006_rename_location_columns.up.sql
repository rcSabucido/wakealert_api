BEGIN;

ALTER TABLE region RENAME COLUMN province_id TO region_id;
ALTER TABLE region RENAME COLUMN province_name TO region_name;

ALTER TABLE province_or_huc RENAME COLUMN city_id TO province_or_huc_id;
ALTER TABLE province_or_huc RENAME COLUMN city_name TO province_or_huc_name;
ALTER TABLE province_or_huc RENAME COLUMN province_id TO region_id;

ALTER TABLE municipal_or_city RENAME COLUMN district_id TO municipal_or_city_id;
ALTER TABLE municipal_or_city RENAME COLUMN district_name TO municipal_or_city_name;
ALTER TABLE municipal_or_city RENAME COLUMN city_id TO province_or_huc_id;

ALTER TABLE barangay RENAME COLUMN district_id TO municipal_or_city_id;

ALTER SEQUENCE province_province_id_seq RENAME TO region_region_id_seq;
ALTER TABLE region ALTER COLUMN region_id SET DEFAULT nextval('region_region_id_seq');
ALTER SEQUENCE region_region_id_seq OWNED BY region.region_id;

ALTER SEQUENCE city_city_id_seq RENAME TO province_or_huc_province_or_huc_id_seq;
ALTER TABLE province_or_huc ALTER COLUMN province_or_huc_id SET DEFAULT nextval('province_or_huc_province_or_huc_id_seq');
ALTER SEQUENCE province_or_huc_province_or_huc_id_seq OWNED BY province_or_huc.province_or_huc_id;

ALTER SEQUENCE district_district_id_seq RENAME TO municipal_or_city_municipal_or_city_id_seq;
ALTER TABLE municipal_or_city ALTER COLUMN municipal_or_city_id SET DEFAULT nextval('municipal_or_city_municipal_or_city_id_seq');
ALTER SEQUENCE municipal_or_city_municipal_or_city_id_seq OWNED BY municipal_or_city.municipal_or_city_id;

COMMIT;
