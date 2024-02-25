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
	orders := []string{"10", "11", "14", "15"}
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

	// Выполнение запроса к базе данных
	rows, err := db.Query(query)
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
		shelf_id int
	)
	// currentShelf := ""
	fmt.Println("=+=+=+=")
	fmt.Println("Страница сборки заказов", strings.Join(orders, ","))
	for rows.Next() {
		err := rows.Scan(&name, &count, &order_id, &product_id, &main, &shelf_name, &shelf_id)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		if main  {
			fmt.Printf("\n===Стеллаж %s\n", shelf_name)
			fmt.Printf("%s (id=%d)\nзаказ %d, %d шт", name, product_id, order_id, count)
		}

		// if currentShelf != "" && currentShelf != fmt.Sprintf("%c", 65+product_id-1) {
		// 	fmt.Println()
		// }
		// if currentShelf != fmt.Sprintf("%c", 65+product_id-1) {
		// 	currentShelf = fmt.Sprintf("%c", 65+product_id-1)
		// 	fmt.Printf("\n===Стеллаж %s\n", currentShelf)
		// }
 
		if !main {
			fmt.Print("\nдоп стеллаж: ")
			//Получение дополнительных стеллажей
			additionalShelves := getAdditionalShelves(db, order_id, name, shelf_id)
			fmt.Print(additionalShelves)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
	}
}

// Функция для получения дополнительных стеллажей
func getAdditionalShelves(db *sql.DB, order_id int, name string, shelf_id int) []string {
	var additionalShelves []string
	query := fmt.Sprintf(`
        SELECT m.shelf_id, s.shelf_name
        FROM Orders o
        JOIN Products p ON o.product_id = p.product_id
        JOIN shelf_product m ON o.product_id = m.product_id
        JOIN Shelves s ON m.shelf_id = s.shelf_id
        WHERE o.order_id = %d AND p.name = '%s' AND m.shelf_id != %d;
    `, order_id, name, shelf_id)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return additionalShelves
	}
	defer rows.Close()
	for rows.Next() {
		var shelf_id int
		var shelf_name string
		err := rows.Scan(&shelf_id, &shelf_name)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		additionalShelves = append(additionalShelves, fmt.Sprintf("%s (id=%d)", shelf_name, shelf_id))
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
	}
	return additionalShelves
}
