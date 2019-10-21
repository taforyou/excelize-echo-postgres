package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

type BonanzaTransaction struct {
	INVESTTXID       string
	PORTFOLIOID      string
	PORTFOLIOCODE    string
	SECURITYID       string
	SECURITYCODE     string
	REFSECURITYID    string
	REFSECURITYCODE  string
	INVESTTXTYPEID   string
	INVESTTXTYPECODE string
	CASHTXTYPEID     string
	CASHTXTYPECODE   string
	TRADEDATE        string
	SETTLEDATE       string
	TRADABLEDATE     string
	UNIT             string
	UNITCOST         string
	YIELD            string
	COSTAMOUNT       string
	ACCRUEDINT       string
	COMMISSIONRATE   string
	COMMISSIONAMOUNT string
	WHTAXRATE        string
	WHTAXAMT         string
	VATRATE          string
	VATAMOUNT        string
	SETTLEAMOUNT     string
	NETAMOUNT        string
	PRINCIPALAMOUNT  string
	ISEFFECTCASH     string
	CURRENCYID       string
	CURRENCYCODE     string
	BROKERID         string
	BROKERCODE       string
	COUNTERPARTYID   string
	COUNTERPARTYCODE string
	TAXPAYERID       string
	TAXPAYERCODE     string
	ISCONFIRMED      string
	ISPOSTED         string
	ISCLOSED         string
	POSTTIME         string
	INFORMATION      string
	CASHGLID         string
	CASHGLCODE       string
	REMARK           string
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

	file := "./02012019.xlsx"
	f, err := excelize.OpenFile(file)

	if err != nil {
		fmt.Println(err)
	} else {

		transactions := []BonanzaTransaction{}
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
		}
		var headerIndex []string

		for i, row := range rows {

			transaction := BonanzaTransaction{}
			for j, colCell := range row {
				// rows ที่ i == 0 เป็น Header
				if i == 0 {
					headerIndex = append(headerIndex, colCell)
				} else {
					//fmt.Println("cell name : ", headerIndex[j], " => ", colCell)

					switch header := headerIndex[j]; header {
					case "INVESTTXID":
						transaction.INVESTTXID = colCell
					case "PORTFOLIOID":
						transaction.PORTFOLIOID = colCell
					case "PORTFOLIOCODE":
						transaction.PORTFOLIOCODE = colCell
					case "SECURITYID":
						transaction.SECURITYID = colCell
					case "SECURITYCODE":
						transaction.SECURITYCODE = colCell
					case "REFSECURITYID":
						transaction.REFSECURITYID = colCell
					case "REFSECURITYCODE":
						transaction.REFSECURITYCODE = colCell
					case "INVESTTXTYPEID":
						transaction.INVESTTXTYPEID = colCell
					case "INVESTTXTYPECODE":
						transaction.INVESTTXTYPECODE = colCell
					case "CASHTXTYPEID":
						transaction.CASHTXTYPEID = colCell
					case "CASHTXTYPECODE":
						transaction.CASHTXTYPECODE = colCell
					case "TRADEDATE":
						transaction.TRADEDATE = colCell
					case "SETTLEDATE":
						transaction.SETTLEDATE = colCell
					case "TRADABLEDATE":
						transaction.TRADABLEDATE = colCell
					case "UNIT":
						transaction.UNIT = colCell
					case "UNITCOST":
						transaction.UNITCOST = colCell
					case "YIELD":
						transaction.YIELD = colCell
					case "COSTAMOUNT":
						transaction.COSTAMOUNT = colCell
					case "ACCRUEDINT":
						transaction.ACCRUEDINT = colCell
					case "COMMISSIONRATE":
						transaction.COMMISSIONRATE = colCell
					case "COMMISSIONAMOUNT":
						transaction.COMMISSIONAMOUNT = colCell
					case "WHTAXRATE":
						transaction.WHTAXRATE = colCell
					case "WHTAXAMT":
						transaction.WHTAXAMT = colCell
					case "VATRATE":
						transaction.VATRATE = colCell
					case "VATAMOUNT":
						transaction.VATAMOUNT = colCell
					case "SETTLEAMOUNT":
						transaction.SETTLEAMOUNT = colCell
					case "NETAMOUNT":
						transaction.NETAMOUNT = colCell
					case "PRINCIPALAMOUNT":
						transaction.PRINCIPALAMOUNT = colCell
					case "ISEFFECTCASH":
						transaction.ISEFFECTCASH = colCell
					case "CURRENCYID":
						transaction.CURRENCYID = colCell
					case "CURRENCYCODE":
						transaction.CURRENCYCODE = colCell
					case "BROKERID":
						transaction.BROKERID = colCell
					case "BROKERCODE":
						transaction.BROKERCODE = colCell
					case "COUNTERPARTYID":
						transaction.COUNTERPARTYID = colCell
					case "COUNTERPARTYCODE":
						transaction.COUNTERPARTYCODE = colCell
					case "TAXPAYERID":
						transaction.TAXPAYERID = colCell
					case "TAXPAYERCODE":
						transaction.TAXPAYERCODE = colCell
					case "ISCONFIRMED":
						transaction.ISCONFIRMED = colCell
					case "ISPOSTED":
						transaction.ISPOSTED = colCell
					case "ISCLOSED":
						transaction.ISCLOSED = colCell
					case "POSTTIME":
						transaction.POSTTIME = colCell
					case "INFORMATION":
						transaction.INFORMATION = colCell
					case "CASHGLID":
						transaction.CASHGLID = colCell
					case "CASHGLCODE":
						transaction.CASHGLCODE = colCell
					case "REMARK":
						transaction.REMARK = colCell
					default:
						fmt.Println("!!! OUT OF HEADER !!! ==> ", header)
					}

					transactions = append(transactions, transaction)
					// switch {
					// case t.Hour() < 12:
					// 	fmt.Println("It's before noon")
					// default:
					// 	fmt.Println("It's after noon")
					// }

					// switch os := runtime.GOOS; os {
					// case "darwin":
					// 	fmt.Println("OS X.")
					// case "linux":
					// 	fmt.Println("Linux.")
					// default:
					// 	// freebsd, openbsd,
					// 	// plan9, windows...
					// 	fmt.Printf("%s.\n", os)
					// }

					//fmt.Println("headerIndex[j] ", headerIndex[j])
					// test := fmt.Sprintf("transaction.%v", headerIndex[j])
					// test = 11
					//transaction.headerIndex[j] = colCell
				}

			}

		}

		// แปลง transactions เพื่อ Inset เข้า DB
		//params["sql"] = fmt.Sprintf("INSERT INTO portfolio_transactions (user_id,set,trade_date,settled_date,security_code,class,security_type,investment,order_no,pre_transaction_no,post_transaction_no,linked_transaction_no,transaction_type,issuer,counter_party_broker,cash_type,yield,units,transaction_currency,cost_per_unit_exclude_commission,commission_percentage,commission_amount,vat,vat_amount,commission_and_vat,settlement_amount,net_amount,cash_security_code,status,posttime,remark,is_import_order,portfolio_name,institute,action,maturity_date,cost_amount,portfolio_currency,security_id,security_code_ref,portfolio_id) VALUES ('%v','%s','%s')", jsonbSetStocks, jsonbMaiStocks)
		// params := make(map[interface{}]interface{})
		// Execute(params)

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
