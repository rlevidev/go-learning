-- name: CreateURL :one
INSERT INTO urls (url_original, short_code)
VALUES ($1, $2)
RETURNING *;

-- name: GetURLByShortCode :one
SELECT * FROM urls
WHERE short_code = $1 LIMIT 1;

-- name: IncrementAccessCount :exec
UPDATE urls
SET access_count = access_count + 1
WHERE short_code = $1;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE short_code = $1;

-- name: UpdateURL :one
UPDATE urls
SET url_original = $2
WHERE short_code = $1
RETURNING *;
