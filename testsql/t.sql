select * from 
    stats_query_summary where 
    query_name like 'rt-%' or 
    query_name like 'archive-%' or 
    query_name like 'live-%' or 
    query_name IN ('trending-repos', 'events-total') limit 3;
