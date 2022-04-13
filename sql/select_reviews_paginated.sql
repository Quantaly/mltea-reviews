-- args: limit, offset
WITH reviews AS (SELECT count(1) FROM review)
SELECT reviews.count, review.reviewer, review.rating, tea.name, tea.caffeinated, review.comment
FROM review
         JOIN tea ON review.tea_id = tea.id
         CROSS JOIN reviews
ORDER BY review.id DESC
LIMIT $1 OFFSET $2;
