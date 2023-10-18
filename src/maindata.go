package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Meal struct {
	Name  string   `json:"name"`
	Foods []string `json:"foods"`
}

var meals = []Meal{
	{
		Name: "早饭",
		Foods: []string{
			"馄饨", "油条", "包子", "牛奶", "面条",
		},
	},
	{
		Name: "午饭",
		Foods: []string{
			"食堂", "超意兴", "肯德基", "麦当劳", "面条",
		},
	},
	{
		Name: "晚饭",
		Foods: []string{
			"食堂", "超意兴", "肯德基", "麦当劳", "面条",
		},
	},
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

var mutex = &sync.Mutex{}

func handleAddFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	newFood := r.URL.Query().Get("food")
	if newFood == "" {
		http.Error(w, "Food parameter is missing", http.StatusBadRequest)
		return
	}

	// 加锁
	mutex.Lock()
	defer mutex.Unlock() // 使用defer确保在函数退出时解锁

	// 添加新的食物到晚饭的列表中
	for i, meal := range meals {
		if meal.Name == "晚饭" {
			meals[i].Foods = append(meals[i].Foods, newFood)
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status": "success"}`)
}

func handleDeleteFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	foodToDelete := r.URL.Query().Get("food")
	if foodToDelete == "" {
		http.Error(w, "Food parameter is missing", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, meal := range meals {
		if meal.Name == "晚饭" {
			for j, food := range meal.Foods {
				if food == foodToDelete {
					meals[i].Foods = append(meal.Foods[:j], meal.Foods[j+1:]...)
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status": "success"}`)
}

func handleGetDinnerFoods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var dinnerFoods []string
	for _, meal := range meals {
		if meal.Name == "晚饭" {
			dinnerFoods = meal.Foods
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dinnerFoods)
}

func main() {
	http.HandleFunc("/api/getfood", handleRequest)
	http.HandleFunc("/api/addfood", handleAddFood)
	http.HandleFunc("/api/deletefood", handleDeleteFood) // 删除食物的接口
	http.HandleFunc("/api/getdinnerfoods", handleGetDinnerFoods)
	//http.Handle("/", http.FileServer(http.FS(staticFiles)))
	// 查询晚饭食物的接口
	http.Handle("/", http.FileServer(http.Dir("./static/chi")))
	http.ListenAndServe(":9091", nil)
}
