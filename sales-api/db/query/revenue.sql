-- name: GetRevenue :many
select o.payment_id, p.*, sum(total_price) as total_amount
from orders o join payments p on p.id = o.payment_id 
where date(o.created_at) = curdate() 
group by payment_id;
