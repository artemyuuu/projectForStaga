package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Получение списка заказов из аргументов командной строки
	orders := os.Args[1:]
	if len(orders) == 0 {
		fmt.Println("No order numbers provided.")
		os.Exit(1)
	}

	// Формирование запроса к базе данных
	query := fmt.Sprintf(`
	    SELECT p.name, o.count, o.order_id, p.product_id, m.main, s.shelf_name, s.shelf_id
	    FROM Orders o
	    JOIN products p ON o.product_id = p.product_id
	    JOIN shelf_product m ON o.product_id = m.product_id
		JOIN shelves s ON m.shelf_id = s.shelf_id
	    WHERE o.order_id IN (%s)
	    ORDER BY m.main DESC, m.shelf_id, p.name;
	`, strings.Join(orders, ","))

	// query, err := ioutil.ReadFile("./sql/products.sql")

	// Выполнение запроса к базе данных
	rows, err := db.Query(string(query))
	if err != nil {
		fmt.Println("Error executing query:", err)
		os.Exit(1)
	}
	defer rows.Close()

	// Обработка результатов запроса
	var (
		name       string
		count      int
		order_id   int
		product_id int
		main       bool
		shelf_name string
		shelf_id   int
	)

	fmt.Println("=+=+=+=")
	fmt.Println("Страница сборки заказов", strings.Join(orders, ","))
	var currentShelf = ""
	for rows.Next() {
		err := rows.Scan(&name, &count, &order_id, &product_id, &main, &shelf_name, &shelf_id)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		if main {
			if currentShelf != shelf_name {
				currentShelf = shelf_name
				fmt.Printf("\n===Стеллаж %s", shelf_name)
			}

			fmt.Printf("\n%s (id=%d)\nзаказ %d, %d шт", name, product_id, order_id, count)
			fmt.Println()
			additionalShelfs := getAdditionalShelvesF(product_id,order_id,db)
			if len(additionalShelfs) > 0 {
				fmt.Println("доп стеллаж: " + additionalShelfs[0])
			}
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
	}
}

func getAdditionalShelvesF(product_id, order_id int, db *sql.DB) []string {
	var additionalShelves []string
	query := fmt.Sprintf(`
        SELECT m.shelf_id, s.shelf_name
        FROM orders o
        JOIN products p ON o.product_id = p.product_id
        JOIN shelf_product m ON o.product_id = m.product_id
        JOIN shelves s ON m.shelf_id = s.shelf_id
        WHERE o.order_id = %d AND m.main != true AND p.product_id = %d;
    `, order_id, product_id)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return additionalShelves
	}
	defer rows.Close()
	for rows.Next() {
		var shelf_id int
		var name string
		err := rows.Scan(&shelf_id, &name)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		additionalShelves = append(additionalShelves, fmt.Sprintf("%s", name))
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
	}
	return additionalShelves
}
