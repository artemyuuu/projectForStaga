CREATE TABLE Products (
    product_id INTEGER PRIMARY KEY,
    name VARCHAR(50)
);

INSERT INTO Products (product_id, name) VALUES
(1, 'Ноутбук'),
(2, 'Телевизор'),
(3, 'Телефон'),
(4, 'Системный блок'),
(5, 'Часы'),
(6, 'Микрофон');

CREATE TABLE Shelves (
    shelf_id INTEGER PRIMARY KEY,
    shelf_name VARCHAR(5)
);

INSERT INTO Shelves (shelf_id, shelf_name) VALUES
(1, 'А'),
(2, 'Б'),
(3, 'В'),
(4, 'З');
(5, 'Ж')

CREATE TABLE Orders (
   order_id INTEGER,
   product_id INTEGER,
   count INTEGER,
   FOREIGN KEY (product_id) REFERENCES Products (product_id)
);
insert into Orders (order_id, product_id, count) values 
(10,1,2),
(10,3,1),
(10,6,1),
(11,2,3),
(14,1,3),
(14,4,4),
(15,5,1);


CREATE TABLE Shelf_Product (
    shelf_id integer,
    product_id integer, 
    main boolean,
    FOREIGN KEY (shelf_id) REFERENCES Shelves (shelf_id),
    FOREIGN KEY (product_id) REFERENCES Products (product_id)
);

INSERT INTO Shelf_Product (shelf_id, product_id, main) VALUES
(1, 1, true),
(1, 2, true),
(1, 5,, false)
(2, 3, true),
(3, 3, false),
(4, 3, false),
(5, 4, true),
(5, 5, true),
(5, 6, true),



