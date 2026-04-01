CREATE OR REPLACE VIEW user_usage AS
    SELECT uploaded_by AS user_id,
    SUM(file_size) AS total_storage_used
    FROM files
    GROUP BY uploaded_by;