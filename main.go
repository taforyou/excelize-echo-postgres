package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "level11_v2"
)

type SqlHandler struct {
	Conn *sql.DB
}

var sqlPostgresHandler = new(SqlHandler)

func Execute(params map[interface{}]interface{}) {

	sql := params["sql"].(string)
	fmt.Println(sql)
	_, err := sqlPostgresHandler.Conn.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(result)
	// ถ้าไม่เรียกนี้แปลว่า Connection อาจจะต่อตลอดเวลา
	//defer sqlPostgresHandler.Conn.Close()
}

type StockPrice struct {
	Bid     string   `json:"bid"`
	Name    string   `json:"name"`
	Open    string   `json:"open"`
	Marks   []string `json:"marks"`
	Offer   string   `json:"offer"`
	Value   string   `json:"value"`
	Change  string   `json:"change"`
	Lowest  string   `json:"lowest"`
	Volume  string   `json:"volume"`
	Highest string   `json:"highest"`
	Lastest float64  `json:"lastest"`
	Pchange string   `json:"pchange"`
}

func main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	e.GET("/cal", calStock)

	e.Logger.Fatal(e.Start(":8084"))

}

func calStock(c echo.Context) (err error) {

	db := ConnDB()

	// วนไปเรื่อยๆ ทั้งเดือน

	sday := "30"
	month := "08"
	year := "2018"

	iday := 0

	// ทำแค่วันเฉยๆ เดือนไม่เอา
	for i := 0; i < 1; i++ {

		iday, _ = strconv.Atoi(sday)
		iday++
		// พอเป็น หลักหน่วยล่ะมีปัญหา 01 => 1
		if iday < 10 {
			sday = "0" + strconv.Itoa(iday)
		} else {
			sday = strconv.Itoa(iday)
		}

		file := "./" + sday + month + year + ".xlsx"
		f, err := excelize.OpenFile(file)

		if err != nil {
			fmt.Println(err)
		} else {
			var maiStocks []StockPrice
			var setStocks []StockPrice

			stockPrice := StockPrice{}
			stockPrice.Name = ""
			// Get value from cell by given worksheet name and axis.

			stockCount := 0
			for i := 4; i < 601; i++ {
				_name := fmt.Sprintf("A%v", i)

				name := f.GetCellValue("Sheet1", _name)

				_lastest := fmt.Sprintf("E%v", i)
				lastest := f.GetCellValue("Sheet1", _lastest)

				// fmt.Println(name)
				// fmt.Println(lastest)
				if strings.TrimSpace(name) == "SPA" {
					// MAI
					stockPrice.Name = strings.TrimSpace(name)
					floatvalue, _ := strconv.ParseFloat(lastest, 64)
					stockPrice.Lastest = floatvalue
					stockPrice.Marks = append(stockPrice.Marks, "")
					maiStocks = append(maiStocks, stockPrice)
				} else if strings.TrimSpace(name) == "AU" {
					// MAI
					fmt.Println("AU")
				} else {
					// SET
					stockPrice.Name = strings.TrimSpace(name)
					floatvalue, _ := strconv.ParseFloat(lastest, 64)
					stockPrice.Lastest = floatvalue
					stockPrice.Marks = append(stockPrice.Marks, "")
					setStocks = append(setStocks, stockPrice)
				}

				stockPrice.Marks = nil
				stockCount++
			}

			// เอา วันที่จากใน Excel มาแปลงเวลาเพื่อเตรียมบรรจุลง DB
			date := f.GetCellValue("Sheet1", "D1")
			layout := "01-02-06"
			t, _ := time.Parse(layout, date)

			jsonbMaiStocks, _ := json.Marshal(maiStocks)
			jsonbSetStocks, _ := json.Marshal(setStocks)
			// อย่าทำแบบนี้เพราะแบบนี้ SQL INJECTION ได้ไว้ก่อนเดี๋ยวมาแก้
			params := make(map[interface{}]interface{})
			params["sql"] = fmt.Sprintf("INSERT INTO stock_price (date,set,mai) VALUES ('%v','%s','%s')", t.Format(time.RFC3339), jsonbSetStocks, jsonbMaiStocks)
			Execute(params)

			fmt.Println("Total Stock SET & MAI", stockCount)
		}
	}

	defer db.Close()

	return c.JSON(http.StatusOK, "OK")

}

func ConnDB() *sql.DB {

	psqlInfo := fmt.Sprintf("port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected! DB")

	sqlPostgresHandler.Conn = db

	return db

}
