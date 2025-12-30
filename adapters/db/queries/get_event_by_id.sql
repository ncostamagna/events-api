SELECT  
    id,
    title,
    description,
    status,
    start_time,
    end_time,
    created_at
FROM events
WHERE id = $1;