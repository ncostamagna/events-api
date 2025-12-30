INSERT INTO events (id, title, description, status, start_time, end_time)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id;