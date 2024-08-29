// Package libs @Author lanpang
// @Date 2024/8/8 下午1:43:00
// @Desc
package libs

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/jackc/pgx/v5"
)

// Close 关闭所有数据库连接
func (dbClient *PGClient) Close() {
	for dbName, conn := range dbClient.Conn {
		if err := conn.Close(context.Background()); err != nil {
			log.Printf("Failed to close connection for database %s: %v", dbName, err)
		}
	}
}

func NewPGClient(conf DB) (*PGClient, error) {
	dbClient := &PGClient{
		Conn: make(map[string]*pgx.Conn),
	}
	databases := []string{"qv30", "users", "customer"}
	for _, dbName := range databases {
		connString := connStr(conf, dbName)
		conn, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			log.Printf("Failed to connect to database %s: %s\n", dbName, err)
			return nil, err
		}
		dbClient.Conn[dbName] = conn
	}

	return dbClient, nil
}

func connStr(conf DB, db string) string {
	scheme := map[bool]string{true: "require", false: "disable"}[conf.Sslmode]
	// 对密码进行 URL 编码
	encodedPassword := url.QueryEscape(conf.Password)
	str := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Username, encodedPassword, conf.Ip, conf.Port, db, scheme)
	return str
}
