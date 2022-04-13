SELECT review.reviewer, review.rating, tea.name, tea.caffeinated, review.comment
FROM review
         JOIN tea ON review.tea_id = tea.id
ORDER BY review.id DESC
LIMIT 5;
