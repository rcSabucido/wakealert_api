-- name: GetRegion :one
SELECT * FROM region WHERE region_psgc = $1;

-- name: ListRegions :many
SELECT * FROM region ORDER BY region_name;

-- name: GetProvinceOrHucsByRegion :many
SELECT * FROM province_or_huc WHERE region_psgc = $1 ORDER BY province_or_huc_name;

-- name: GetMunicipalOrCitiesByProvinceOrHuc :many
SELECT * FROM municipal_or_city WHERE province_or_huc_psgc = $1 ORDER BY city_mun_name;

-- name: GetBarangaysByMunicipalOrCity :many
SELECT * FROM barangays WHERE city_mun_psgc = $1 ORDER BY barangay_name;

-- name: GetAddressLine :one
SELECT * FROM address_line WHERE address_id = $1;

-- name: CreateAddressLine :one
INSERT INTO address_line (barangay_psgc, address_line)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateAddressLine :one
UPDATE address_line
SET barangay_psgc = $1,
	address_line = $2
WHERE address_id = $3
RETURNING *;

-- name: GetAddressDetails :one
SELECT
	address_line.address_id,
	address_line.address_line,
	barangays.barangay_name,
	municipal_or_city.city_mun_name,
	province_or_huc.province_or_huc_name,
	region.region_name
FROM address_line
JOIN barangays ON barangay.barangay_psgc = address_line.barangay_psgc
JOIN municipal_or_city ON municipal_or_city.city_mun_psgc = barangay.city_mun_psgc
JOIN province_or_huc ON province_or_huc.province_or_huc_psgc = municipal_or_city.province_or_huc_psgc
JOIN region ON region.region_psgc = province_or_huc.region_psgc
WHERE address_line.address_id = $1;