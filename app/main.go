package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/get-token", GetTokenHandler).Methods("GET")

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	r.Handle("/status", NotImplemented).Methods("GET")
	// Добавляем прослойку к products и feedback роутам, все остальные
	// роуты у нас публичные
	r.Handle("/products",
		jwtMiddleware.Handler(ProductsHandler)).Methods("GET")
	r.Handle("/products/{slug}/feedback",
		jwtMiddleware.Handler(AddFeedbackHandler)).Methods("POST")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

var mySigningKey = []byte("secret")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter,
	r *http.Request) {
	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	// Устанавливаем набор параметров для токена
	claims["admin"] = true
	claims["name"] = "Ado Kukic"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims

	// Подписываем токен нашим секретным ключем
	tokenString, _ := token.SignedString(mySigningKey)

	// Отдаем токен клиенту
	w.Write([]byte(tokenString))
})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

type Product struct {
	Id          int
	Name        string
	Slug        string
	Description string
}

var products = []Product{
	Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters",
		Description: "Shoot your way to the top on 14 different hoverboards"},
	Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer",
		Description: "Explore the depths of the sea in this one of a kind"},
	Product{Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park",
		Description: "Go back 65 million years in the past and ride a T-Rex"},
	Product{Id: 4, Name: "Cars VR", Slug: "cars-vr",
		Description: "Get behind the wheel of the fastest cars in the world."},
	Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood",
		Description: "Pick up the bow and arrow and master the art of archery"},
	Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr",
		Description: "Explore the seven wonders of the world in VR"},
}

// Хендлер StatusHandler будет срабатывать в тот момент момент, когда
// пользователь обращается по роуту /status. Этот хендлер просто возвращает
// строку "API is up and running".
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

// ProductsHandler срабатывает в момент вызова роута /products
// Этот хендлер возвращает пользователю список продуктов для оценки.
var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Конвертируем наш слайс продуктов в json
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

// Этот хендлер позволяет добавить позитивный или негативный отзыв
// по конкретному продукту. Правильно было бы сохранить результат в базу
// данных, но для демо нам это не нужно.
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if product.Slug != "" {
		payload, _ := json.Marshal(product)
		w.Write([]byte(payload))

	} else {
		w.Write([]byte("Product Not Found"))
	}
})
