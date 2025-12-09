package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
)

const (
	serverURL      = "http://localhost:8080"
	defaultOrderId = "3"
)

func main() {
	rootCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := order_v1.NewClient(serverURL)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании клиента: %v", err)
	}

	log.Println("=== Тестирование API (параллельно) ===")
	log.Println()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		ctx, ctxCancel := context.WithTimeout(rootCtx, 10*time.Second)
		defer ctxCancel()

		log.Printf("[GET] Получение данных заказа для ID %s\n", defaultOrderId)
		log.Println("===================================================")

		orderResp, err := client.GetOrderById(ctx, order_v1.GetOrderByIdParams{
			ID: defaultOrderId,
		})
		if err != nil {
			var errResp *order_v1.AppErrorStatusCode
			if errors.As(err, &errResp) && errResp.StatusCode == 404 {
				log.Printf("ℹ️ [GET] Данные не найдены (404)\n")
				return
			}

			log.Printf("❌ [GET] Ошибка при получении заказа: %v\n", err)
			return
		}

		log.Printf("✅ [GET] Данные о заказе для ID %s: %+v\n", defaultOrderId, orderResp)
	}()

	go func() {
		defer wg.Done()

		ctx, ctxCancel := context.WithTimeout(rootCtx, 10*time.Second)
		defer ctxCancel()

		go func() {
			time.Sleep(500 * time.Millisecond)

			log.Printf("🗑️ [DELETE] Return ордера для ID %s\n", defaultOrderId)
			log.Println("===========================================================")

			orderReturnResp, err := client.ReturnOrder(ctx, order_v1.ReturnOrderParams{
				ID: defaultOrderId,
			})
			if err != nil {
				log.Printf("❌ [DELETE] Ошибка при return ордера: %v\n", err)
				return
			}

			log.Printf("✅ [DELETE] Ордер успешно return'нут: %+v\n", orderReturnResp)
		}()

		log.Printf("[GET] Получение данных заказа для ID %s\n", defaultOrderId)
		log.Println("===================================================")

		orderResp, err := client.GetOrderById(ctx, order_v1.GetOrderByIdParams{
			ID: defaultOrderId,
		})
		if err != nil {
			var errResp *order_v1.AppErrorStatusCode
			if errors.As(err, &errResp) && errResp.StatusCode == 404 {
				log.Printf("ℹ️ [GET] Данные не найдены (404)\n")
				return
			}

			log.Printf("❌ [GET] Ошибка при получении заказа: %v\n", err)
			return
		}

		log.Printf("✅ [GET] Данные о заказе для ID %s: %+v\n", defaultOrderId, orderResp)
	}()

	wg.Wait()

	if err := rootCtx.Err(); err != nil {
		log.Printf("⚠️ root context завершился с ошибкой: %v\n", err)
	} else {
		log.Println("🎉 Тестирование завершено успешно!")
	}
}
