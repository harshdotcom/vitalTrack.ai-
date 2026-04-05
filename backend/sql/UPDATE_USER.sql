UPDATE users
SET name = $1,
dob = $2,
gender = $3,
profile_pic = $4,
is_verified = $5
WHERE user_id = $6