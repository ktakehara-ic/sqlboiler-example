package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ktakehara-icd/sqlboiler-example/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	db, err := sql.Open("mysql",
		"user:pass@tcp(127.0.0.1:33060)/items?parseTime=true&charset=utf8mb4")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// トランザクションの中での処理
	{
		// NOTE: トランザクションの作成は context.Background() (timeout しないようにする)
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			panic(err)
		}

		time.Sleep(5 * time.Second) // TIMEOUT

		item := models.Item{
			Name:  "item1",
			Price: 100,
		}
		// NOTE: Insert はタイムアウトする context を渡す
		if err := item.Insert(ctx, tx, boil.Infer()); err != nil {
			fmt.Printf("INSERT ERROR: %v\n", err)
		}

		if err := tx.Commit(); err != nil {
		  // NOTE: トランザクションはタイムアウトしないので、ここには落ちない
			fmt.Printf("COMMIT ERROR: %v\n", err)
		}
	}

	items, err := models.Items().All(context.Background(), db)
	if err != nil {
		panic(err)
	}

	if len(items) == 0 {
	  // NOTE: INSERT がタイムアウトで失敗しているのでアイテムなし
		fmt.Println("No items")
		return
	}
	for _, item := range items {
		fmt.Printf("ID: %d, Name: %s\n", item.ID, item.Name)
	}
}
