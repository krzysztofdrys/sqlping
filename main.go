package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/krzysztofdrys/sqlping/pinger"
)

type noopLogger struct{}

func (n noopLogger) Print(v ...interface{}) {}

func main() {
	if err := mysql.SetLogger(&noopLogger{}); err != nil {
		panic(err)
	}

	flag.Parse()
	dsn := flag.Arg(0)
	if dsn == "" {
		fmt.Fprintf(os.Stderr, "DSN of the database is required")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	err = db.PingContext(ctx)
	cancel()

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	p := pinger.State{StartedAt: time.Now().UTC(), Error: errStr}
	lastPrint := time.Time{}
	for range time.Tick(100 * time.Millisecond) {
		ctx := context.Background()
		ctx, cancel = context.WithTimeout(ctx, 200*time.Millisecond)
		err = db.PingContext(ctx)
		cancel()

		next := p.OnPing(err)
		if next != p {
			p.EndedAt = time.Now()
			log.Println(p)
			if p.Error != "" && err == nil {
				log.Println("Connection is now up")
			} else if p.Error == "" && err != nil {
				log.Printf("Connection is now down, the latest error: %v", err)
			}
			p = next
		} else if time.Since(lastPrint) > 30*time.Second {
			log.Println(p)
			lastPrint = time.Now().UTC()
		}
	}

}
