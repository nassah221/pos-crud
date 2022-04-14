-- name: CreateOrder :execresult
INSERT INTO orders (
    cashier_id, payment_id, total_price, price_paid, total_return, receipt_id
) VALUES (
    ?,?,?,?,?,?
);

-- name: CreateOrderDetails :exec
INSERT INTO order_details (
    order_id,product_id, product_name, discount_id, qty, price, total_final_price, total_normal_price
) VALUES (
    ?,?,?,?,?,?,?,?
);

-- name: GetOrderDetails :many
select o.id as orderId, o.payment_id as paymentTypesId, o.total_price,o.price_paid,o.total_return,o.receipt_id,o.created_at,
c.id as cashiersId, c.name as cashierName,
p.id as paymentTypeId, p.name as paymentName, p.logo, p.type paymentType,
od.product_id,od.product_name,od.discount_id as discountsId,od.qty as productQty,od.price,od.total_final_price,od.total_normal_price,
d.id,d.qty as discountQty,d.type as discountType,d.result,d.expired_at,d.expired_at_format,d.string_format
from orders as o
left join cashiers as c on c.id = o.cashier_id 
left join payments as p on p.id = o.payment_id
inner join order_details as od on od.order_id = o.id
left join discounts as d on od.discount_id = d.id
where o.id=?;

-- name: ListAllOrderDetails :many
select o.id as orderId, o.payment_id as paymentTypesId, o.total_price,o.price_paid,o.total_return,o.receipt_id,o.created_at,
c.id as cashiersId, c.name as cashierName,
p.id as paymentTypeId, p.name as paymentName, p.logo, p.type as paymentType,
od.product_id,od.product_name,od.price,od.discount_id as discountsId,od.qty as productQty,od.total_final_price,od.total_normal_price,
d.id,d.qty as discountQty,d.type as discountType,d.result,d.expired_at,d.expired_at_format,d.string_format
from orders as o
left join cashiers as c on c.id = o.cashier_id 
left join payments as p on p.id = o.payment_id
inner join order_details as od on od.order_id = o.id
left join discounts as d on od.discount_id = d.id
ORDER BY o.id LIMIT ? OFFSET ?;
