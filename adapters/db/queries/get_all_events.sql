SELECT  
    id,
    title,
    description,
    status,
    start_time,
    end_time,
    created_at
FROM events
ORDER BY start_time ASC;