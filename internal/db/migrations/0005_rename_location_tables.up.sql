BEGIN;  

ALTER TABLE province RENAME TO region;
ALTER TABLE city RENAME TO province_or_huc;
ALTER TABLE district RENAME TO municipal_or_city;

ALTER TABLE region RENAME CONSTRAINT province_pkey TO region_pkey;
ALTER TABLE province_or_huc RENAME CONSTRAINT city_pkey TO province_or_huc_pkey;
ALTER TABLE municipal_or_city RENAME CONSTRAINT district_pkey TO municipal_or_city_pkey;

ALTER TABLE province_or_huc RENAME CONSTRAINT fk_city_province TO fk_province_or_huc_region;
ALTER TABLE municipal_or_city RENAME CONSTRAINT fk_district_city TO fk_municipal_or_city_province_or_huc;
ALTER TABLE barangay RENAME CONSTRAINT fk_barangay_district TO fk_barangay_municipal_or_city;

COMMIT;
