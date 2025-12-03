package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"context"
	"log"

	"github.com/nxlak/go-pvz/internal/config"
	orderPostgres "github.com/nxlak/go-pvz/internal/repository/storage/postgres"
	"github.com/nxlak/go-pvz/internal/usecase/service"
	"github.com/nxlak/go-pvz/pkg/client/postgres"
)

func main() {
	cfg := config.GetConfig()

	client, err := postgres.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	orderRepo := orderPostgres.NewRepositoty(client)
	service := service.NewService(orderRepo)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите команду (или 'exit' для выхода):")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "accept-order":
			if len(parts) == 7 &&
				parts[1] == "--order-id" &&
				parts[3] == "--user-id" &&
				parts[5] == "--expires" {
				expires, err := time.Parse("2006-01-02", parts[6])
				if err != nil {
					fmt.Println("Ошибка даты. Используйте формат: yyyy-mm-dd")
					continue
				}
				fmt.Printf("Принятие заказа %s для пользователя %s до %s\n", parts[2], parts[4], parts[6])
				if err := service.AcceptOrder(parts[2], parts[4], expires); err != nil {
					fmt.Printf("Ошибка при принятии заказа: %v\n", err)
				} else {
					fmt.Println("Заказ успешно принят!")
				}
			} else {
				fmt.Println("Ошибка в команде. Проверьте правильность синтаксиса.")
				fmt.Println("Используйте: accept-order --order-id <id> --user-id <id> --expires <yyyy-mm-dd>")
			}

		case "return-order":
			if len(parts) == 3 && parts[1] == "--order-id" {
				fmt.Printf("Возврат заказа %s\n", parts[2])
				if err := service.ReturnOrder(parts[2]); err != nil {
					fmt.Printf("Ошибка при возврате заказа: %v\n", err)
				} else {
					fmt.Println("Заказ успешно возвращен!")
				}
			} else {
				fmt.Println("Ошибка в команде. Проверьте правильность синтаксиса.")
				fmt.Println("Используйте: return-order --order-id <id>")
			}

		case "help":
			fmt.Println("Доступные команды:")
			fmt.Println("  accept-order --order-id <id> --user-id <id> --expires <yyyy-mm-dd>")
			fmt.Println("  return-order --order-id <id>")
			fmt.Println("  help - показать справку")
			fmt.Println("  exit - выйти из программы")

		default:
			fmt.Println("Неизвестная команда. Введите 'help' для справки.")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения ввода: %v\n", err)
	}
}
