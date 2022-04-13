WITH tea_info AS (SELECT tea.id, avg(review.rating) AS rating, count(1) AS review_count
                  FROM tea
                           JOIN review ON tea.id = review.tea_id
                  GROUP BY tea.id)
SELECT tea.name, tea.caffeinated, tea_info.rating, tea_info.review_count
FROM tea
         JOIN tea_info ON tea.id = tea_info.id
ORDER BY tea_info.rating DESC
LIMIT 10;
