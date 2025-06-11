package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igor-fedko/go_multiply_app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection
var logsCollection *mongo.Collection

// подключение к MongoDB
func connectDB() (*mongo.Client, error) {
	// Создаем URI для подключения к MongoDB с учетными данными
	uri := "mongodb://admin:password@mongodb:27017"

	// Настраиваем опции подключения
	clientOptions := options.Client().ApplyURI(uri)

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Успешное подключение к MongoDB!")

	return client, nil
}

// Функция логирования операций
func logOperation(c *gin.Context, operation string, input string, result float64) {
	logEntry := models.LogEntry{
		Operation: operation,
		Input:     input,
		Result:    fmt.Sprintf("%f", result),
		UserIP:    c.ClientIP(),
		Timestamp: time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := logsCollection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Ошибка при логировании операции: %v", err)
	}
}

// Обработчик главной страницы
func indexHandler(c *gin.Context) {
	// Получаем все результаты из базы данных
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при получении результатов: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var results []models.Result
	if err = cursor.All(ctx, &results); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при обработке результатов: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Results": results,
	})
}

// Обработчик умножения чисел
func multiplyHandler(c *gin.Context) {
	// Получаем данные из формы
	number1Str := c.PostForm("number1")
	number2Str := c.PostForm("number2")

	// Преобразуем строки в числа
	number1, err := strconv.ParseFloat(number1Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат первого числа",
		})
		return
	}

	number2, err := strconv.ParseFloat(number2Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат второго числа",
		})
		return
	}

	// Вычисляем произведение
	product := number1 * number2

	// Создаем новый результат
	result := models.Result{
		Number1:   number1,
		Number2:   number2,
		Result:    product,
		Operation: "multiply",
		CreatedAt: time.Now().UTC(),
	}

	// Сохраняем результат в MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, result)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при сохранении результата: " + err.Error(),
		})
		return
	}

	// Логируем операцию
	logOperation(c, "multiply", fmt.Sprintf("%f * %f", number1, number2), product)

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusSeeOther, "/")
}

// Обработчик деления чисел
func divideHandler(c *gin.Context) {
	// Получаем данные из формы
	number1Str := c.PostForm("number1")
	number2Str := c.PostForm("number2")

	// Преобразуем строки в числа
	number1, err := strconv.ParseFloat(number1Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат первого числа",
		})
		return
	}

	number2, err := strconv.ParseFloat(number2Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат второго числа",
		})
		return
	}

	// Проверяем деление на ноль
	if number2 == 0 {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Деление на ноль невозможно",
		})
		return
	}

	// Вычисляем частное
	quotient := number1 / number2

	// Создаем новый результат
	result := models.Result{
		Number1:   number1,
		Number2:   number2,
		Result:    quotient,
		Operation: "divide",
		CreatedAt: time.Now().UTC(),
	}

	// Сохраняем результат в MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, result)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при сохранении результата: " + err.Error(),
		})
		return
	}

	// Логируем операцию
	logOperation(c, "divide", fmt.Sprintf("%f / %f", number1, number2), quotient)

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusSeeOther, "/")
}

// Обработчик сложения чисел
func addHandler(c *gin.Context) {
	// Получаем данные из формы
	number1Str := c.PostForm("number1")
	number2Str := c.PostForm("number2")

	// Преобразуем строки в числа
	number1, err := strconv.ParseFloat(number1Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат первого числа",
		})
		return
	}

	number2, err := strconv.ParseFloat(number2Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат второго числа",
		})
		return
	}

	// Вычисляем сумму
	sum := number1 + number2

	// Создаем новый результат
	result := models.Result{
		Number1:   number1,
		Number2:   number2,
		Result:    sum,
		Operation: "add",
		CreatedAt: time.Now().UTC(),
	}

	// Сохраняем результат в MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, result)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при сохранении результата: " + err.Error(),
		})
		return
	}

	// Логируем операцию
	logOperation(c, "add", fmt.Sprintf("%f + %f", number1, number2), sum)

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusSeeOther, "/")
}

// Обработчик вычитания чисел
func subtractHandler(c *gin.Context) {
	// Получаем данные из формы
	number1Str := c.PostForm("number1")
	number2Str := c.PostForm("number2")

	// Преобразуем строки в числа
	number1, err := strconv.ParseFloat(number1Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат первого числа",
		})
		return
	}

	number2, err := strconv.ParseFloat(number2Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат второго числа",
		})
		return
	}

	// Вычисляем разность
	difference := number1 - number2

	// Создаем новый результат
	result := models.Result{
		Number1:   number1,
		Number2:   number2,
		Result:    difference,
		Operation: "subtract",
		CreatedAt: time.Now().UTC(),
	}

	// Сохраняем результат в MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, result)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при сохранении результата: " + err.Error(),
		})
		return
	}

	// Логируем операцию
	logOperation(c, "subtract", fmt.Sprintf("%f - %f", number1, number2), difference)

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusSeeOther, "/")
}

// Обработчик возведения в квадрат
func squareHandler(c *gin.Context) {
	// Получаем данные из формы
	number1Str := c.PostForm("number1")

	// Преобразуем строку в число
	number1, err := strconv.ParseFloat(number1Str, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Неверный формат числа",
		})
		return
	}

	// Вычисляем квадрат
	squared := number1 * number1

	// Создаем новый результат
	result := models.Result{
		Number1:   number1,
		Number2:   0, // Для операции возведения в квадрат второе число не используется
		Result:    squared,
		Operation: "square",
		CreatedAt: time.Now().UTC(),
	}

	// Сохраняем результат в MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, result)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Ошибка при сохранении результата: " + err.Error(),
		})
		return
	}

	// Логируем операцию
	logOperation(c, "square", fmt.Sprintf("%f²", number1), squared)

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusSeeOther, "/")
}

func main() {
	// Подключаемся к MongoDB
	var err error
	client, err = connectDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к MongoDB: %v", err)
	}

	// Закрываем соединение при завершении
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Ошибка при отключении от MongoDB: %v", err)
		}
	}()

	// Получаем коллекцию для работы с результатами
	collection = client.Database("multiply_app").Collection("results")

	// Получаем коллекцию для логов
	logsCollection = client.Database("multiply_app").Collection("logs")

	// Создаем Gin роутер
	router := gin.Default()

	// Загружаем HTML шаблоны
	router.SetHTMLTemplate(template.Must(template.ParseFiles("templates/index.html")))

	// Определяем маршруты
	router.GET("/", indexHandler)
	router.POST("/multiply", multiplyHandler)
	router.POST("/divide", divideHandler)
	router.POST("/add", addHandler)
	router.POST("/subtract", subtractHandler)
	router.POST("/square", squareHandler)

	// Запускаем сервер
	log.Println("Сервер запущен на http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
