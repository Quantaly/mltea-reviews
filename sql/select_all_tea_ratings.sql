WITH tea_info AS (SELECT tea.id, avg(review.rating) AS rating, count(1) AS rating_count
                  FROM tea
                           JOIN review ON tea.id = review.tea_id
                  GROUP BY tea.id)
SELECT tea.name, tea.caffeinated, coalesce(tea_info.rating, -1) AS rating, coalesce(tea_info.rating_count, 0) AS rating_count
FROM tea
         LEFT JOIN tea_info ON tea.id = tea_info.id
ORDER BY tea_info.rating DESC NULLS LAST, tea.name;