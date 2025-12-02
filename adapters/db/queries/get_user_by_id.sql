SELECT  
    id,
    title,
    description,
    start_time,
    end_time,
    created_at
FROM events
WHERE id = $1;