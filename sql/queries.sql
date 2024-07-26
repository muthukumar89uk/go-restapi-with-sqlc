
-- name: CreateUser :one
INSERT INTO users (Username,Email,PhoneNumber,Password,Role)
VALUES ($1, $2,$3,$4,$5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByPhoneNumber :one
SELECT * FROM users
WHERE phonenumber = $1 LIMIT 1;

-- name: CreateCareer :one
INSERT INTO career (Company,Position,Jobtype,Description,StartDate,EndDate)
VALUES ($1, $2,$3,$4,$5,$6)
RETURNING *;

-- name: GetCareerByJobId :one
SELECT * FROM career
WHERE jobid = $1 LIMIT 1;

-- name: GetAllCareerDetails :many
SELECT * FROM career;

-- name: UpdateCareerByJobId :one
UPDATE career
SET company=$1,position=$2,jobtype=$3,description=$4
WHERE jobid = $5
RETURNING *;

-- name: DeleteCareerByJobId :one
DELETE
FROM career
WHERE jobid = $1
RETURNING *;


-- name: CreateProfile :one
INSERT INTO profile (FullName,Age,Gender,Address)
VALUES ($1, $2,$3,$4)
RETURNING *;

-- name: GetProfileByuserId :one
SELECT * FROM profile
WHERE userid = $1 LIMIT 1;

-- name: GetAllProfileDetails :many
SELECT * FROM profile;

-- name: UpdateProfileByuserId :one
UPDATE profile
SET FullName=$1,Age=$2,Gender=$3,Address=$4
WHERE userid = $5
RETURNING *;

-- name: DeleteProfileByUserId :one
DELETE
FROM profile
WHERE userid = $1
RETURNING *;

-- name: GetallusersEmail :many

SELECT email FROM users WHERE role = 'user';
