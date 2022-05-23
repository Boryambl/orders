# orders

Simple rest app with users, products and orders.

To run this app you need golang and docker-compose.

To setup postgres run "docker-compose -f postgres.yml up -d".

DB structure in sql/up.sql .

After "make" the binary will appear in the folder bin/ .

To run app use ./bin/orders server --port {port} --pg postgres://user:password@127.0.0.1:5439/db (if you start db from postgres.yml use postgres://admin:nimda@127.0.0.1:5439/orders).

You can create user by post request to /api/users with struct like {
surname,firstname,middlename(all string), age(int),sex(bool)
} in body.
To update user use put request to /api/users/:id with same struct.
To delete user use delete request to same url.
To get users use get request to /api/users, also you can get one user by get request to /api/users?id={id}.

You can create product by post request to /api/products with struct like {
description(string), price({currency(string), value(int)}, left_in_stock(int))
} in body.
To update the stock quantity use put request to /api/products/:id with struct like {count(int)}.
To delete product use delete request to same url.
To get products use get request to /api/products, also you can get one products by get request to /api/products?id={id}.

You can create order by post request to /api/orders with struct like {
user_id(string), products[{product_id(string), count(int)}]
} in body.
To add more products to order use put request to /api/orders/:id with struct like {products[{product_id(string), count(int)}]}.
To delete order use delete request to same url.
To get orders use get request to /api/products, also you can get all products from user by get request to /api/products?user_id={id}.
