CREATE TABLE `cashiers` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT now(),
  `updated_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `products` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `sku` varchar(255) UNIQUE NOT NULL,
  `name` varchar(255) NOT NULL,
  `stock` int NOT NULL,
  `price` int NOT NULL,
  `image_url` text NOT NULL,
  `category_id` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT now(),
  `updated_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `product_discount` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `discount_id` int NOT NULL,
  `product_id` int NOT NULL
);

CREATE TABLE `categories` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT now(),
  `updated_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `discounts` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `qty` int NOT NULL,
  `type` varchar(255) NOT NULL,
  `result` int NOT NULL,
  `expired_at` timestamp NOT NULL,
  `expired_at_format` varchar(255) NOT NULL,
  `string_format` text NOT NULL
);

CREATE TABLE `payments` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `logo` text,
  `created_at` timestamp NOT NULL DEFAULT now(),
  `updated_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `orders` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `cashier_id` int NOT NULL,
  `payment_id` int NOT NULL,
  `total_price` int NOT NULL,
  `price_paid` int NOT NULL,
  `total_return` int NOT NULL,
  `receipt_id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT now(),
  `updated_at` timestamp NOT NULL DEFAULT now()
);

CREATE TABLE `sold` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `product_id` int NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `total_qty` int NOT NULL,
  `total_amount` bigint NOT NULL
);

CREATE TABLE `order_products` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `order_id` int,
  `product_id` int
);

CREATE TABLE `order_details` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `product_id` int NOT NULL,
  `order_id` int NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `discount_id` int,
  `qty` int NOT NULL,
  `price` int NOT NULL,
  `total_final_price` bigint NOT NULL,
  `total_normal_price` bigint NOT NULL
);

ALTER TABLE `sold` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `order_products` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_products` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `order_details` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_details` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `order_details` ADD FOREIGN KEY (`discount_id`) REFERENCES `discounts` (`id`);

ALTER TABLE `products` ADD FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);

ALTER TABLE `product_discount` ADD FOREIGN KEY (`discount_id`) REFERENCES `discounts` (`id`);

ALTER TABLE `product_discount` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`cashier_id`) REFERENCES `cashiers` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`);
