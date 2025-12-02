INSERT INTO events (id, title, description, start_time, end_time)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
)
RETURNING id;