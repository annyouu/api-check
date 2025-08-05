package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// AppPay APIベースURL
// SwaggerドキュメントのURLに修正済み
const appPayAPIBaseURL = "https://tjufwmnunr.ap-northeast-1.awsapprunner.com/api/v1"
const appPayOrdersPath = "/orders"

// getOrdersData は AppPayの/ordersエンドポイントから全注文データを取得するハンドラ
func getOrdersData(w http.ResponseWriter, r *http.Request) {
    fullURL := appPayAPIBaseURL + appPayOrdersPath
    log.Printf("📡 Requesting AppPay API for all orders: %s\n", fullURL)

    response, err := http.Get(fullURL)
    if err != nil {
        log.Printf("Failed to fetch from AppPay API: %v", err)
        http.Error(w, fmt.Sprintf("AppPay API呼び出し失敗: %v", err), http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    log.Printf(" AppPay API responded with status: %d\n", response.StatusCode)

    if response.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(response.Body)
        log.Printf("⚠️ Error response body: %s", string(body))
        http.Error(w, fmt.Sprintf("AppPay APIステータス異常: %d\nBody: %s", response.StatusCode, string(body)), http.StatusInternalServerError)
        return
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Printf("Failed to read AppPay API response: %v", err)
        http.Error(w, fmt.Sprintf("AppPay APIレスポンス読み取り失敗: %v", err), http.StatusInternalServerError)
        return
    }

    log.Printf("Success: %d bytes retrieved from AppPay", len(body))

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(body)
}

// getSingleOrderData は AppPayの/orders/{orderId}エンドポイントから単一の注文データを取得するハンドラ
func getSingleOrderData(w http.ResponseWriter, r *http.Request) {
    // URLパスからorderIdを取得
    vars := mux.Vars(r)
    orderId := vars["orderId"]

    if orderId == "" {
        http.Error(w, "orderIdが指定されていません", http.StatusBadRequest)
        return
    }

    fullURL := fmt.Sprintf("%s%s/%s", appPayAPIBaseURL, appPayOrdersPath, orderId)
    log.Printf("Requesting AppPay API for single order: %s\n", fullURL)

    response, err := http.Get(fullURL)
    if err != nil {
        log.Printf("Failed to fetch from AppPay API: %v", err)
        http.Error(w, fmt.Sprintf("AppPay API呼び出し失敗: %v", err), http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    log.Printf("AppPay API responded with status: %d\n", response.StatusCode)

    if response.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(response.Body)
        log.Printf("Error response body: %s", string(body))
        http.Error(w, fmt.Sprintf("AppPay APIステータス異常: %d\nBody: %s", response.StatusCode, string(body)), http.StatusInternalServerError)
        return
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Printf("Failed to read AppPay API response: %v", err)
        http.Error(w, fmt.Sprintf("AppPay APIレスポンス読み取り失敗: %v", err), http.StatusInternalServerError)
        return
    }

    log.Printf("Success: %d bytes retrieved for orderId %s", len(body), orderId)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(body)
}

func main() {
	router := mux.NewRouter()

	// 全注文を取得するためのエンドポイント
	router.HandleFunc("/get-orders", getOrdersData).Methods("GET")
	// 単一の注文を取得するためのエンドポイント
	router.HandleFunc("/get-order/{orderId}", getSingleOrderData).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	port := "5001"
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
