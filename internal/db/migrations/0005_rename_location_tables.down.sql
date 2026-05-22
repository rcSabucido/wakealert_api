BEGIN;

ALTER TABLE municipal_or_city RENAME CONSTRAINT municipal_or_city_pkey TO district_pkey;
ALTER TABLE province_or_huc RENAME CONSTRAINT province_or_huc_pkey TO city_pkey;
ALTER TABLE region RENAME CONSTRAINT region_pkey TO province_pkey;

ALTER TABLE province_or_huc RENAME CONSTRAINT fk_province_or_huc_region TO fk_city_province;
ALTER TABLE municipal_or_city RENAME CONSTRAINT fk_municipal_or_city_province_or_huc TO fk_district_city;
ALTER TABLE barangay RENAME CONSTRAINT fk_barangay_municipal_or_city TO fk_barangay_district;

ALTER TABLE municipal_or_city RENAME TO district;
ALTER TABLE province_or_huc RENAME TO city;
ALTER TABLE region RENAME TO province;

COMMIT;
