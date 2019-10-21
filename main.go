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

func Execute2(txn BonanzaTransaction) error {
	_, err := sqlPostgresHandler.Conn.Exec(`INSERT INTO portfolio_transactions (invest_tx_id,portfolio_id,portfolio_code,security_id,security_code,ref_security_id,ref_security_code,invest_tx_type_id,invest_tx_type_code,cash_type_id,cash_type_code,trade_date,settled_date,tradable_date,units,unit_cost,yield,cost_amount,accrued_int,commission_rate,commision_amount,wh_tax_rate,wh_tax_amount,vat_rate,vat_amount,settle_amount,net_amount,principal_amount,is_effect_cash,currency_id,currency_code,broker_id,broker_code,counter_party_id,counter_party_code,tax_payer_id,tax_payer_code) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37)`, txn.INVESTTXID, txn.PORTFOLIOID, txn.PORTFOLIOCODE, txn.SECURITYID, txn.SECURITYCODE, NewNullString(txn.REFSECURITYID), NewNullString(txn.REFSECURITYCODE), txn.INVESTTXTYPEID, txn.INVESTTXTYPECODE, NewNullString(txn.CASHTXTYPEID), NewNullString(txn.CASHTXTYPECODE), txn.TRADEDATE, txn.SETTLEDATE, NewNullString(txn.TRADABLEDATE), txn.UNIT, txn.UNITCOST, txn.YIELD, txn.COSTAMOUNT, txn.ACCRUEDINT, txn.COMMISSIONRATE, txn.COMMISSIONAMOUNT, txn.WHTAXRATE, txn.WHTAXAMT, NewNullString(txn.VATRATE), txn.VATAMOUNT, txn.SETTLEAMOUNT, txn.NETAMOUNT, txn.PRINCIPALAMOUNT, txn.ISEFFECTCASH, txn.CURRENCYID, txn.CURRENCYCODE, NewNullString(txn.BROKERID), NewNullString(txn.BROKERCODE), NewNullString(txn.COUNTERPARTYID), NewNullString(txn.COUNTERPARTYCODE), NewNullString(txn.TAXPAYERID), NewNullString(txn.TAXPAYERCODE))
	if err != nil {
		panic(err)
	}
	return err
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
				// fmt.Println(i, j, colCell)
				// rows ที่ i == 0 เป็น Header
				headerIndex = append(headerIndex, colCell)

				if i != 0 {
					switch header := headerIndex[j]; header {

					case "INVESTTXID":
						// post_transaction_no
						//iColCell, _ := strconv.Atoi(colCell)
						transaction.INVESTTXID = colCell
					case "PORTFOLIOID":
						transaction.PORTFOLIOID = colCell
					case "PORTFOLIOCODE":
						transaction.PORTFOLIOCODE = colCell
					case "SECURITYID":
						transaction.SECURITYID = colCell
					case "SECURITYCODE":
						// security_code
						transaction.SECURITYCODE = colCell
					case "REFSECURITYID":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.REFSECURITYID = colCell
					case "REFSECURITYCODE":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.REFSECURITYCODE = colCell
					case "INVESTTXTYPEID":
						transaction.INVESTTXTYPEID = colCell
					case "INVESTTXTYPECODE":
						transaction.INVESTTXTYPECODE = colCell
					case "CASHTXTYPEID":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.CASHTXTYPEID = colCell
					case "CASHTXTYPECODE":
						// ตรงนี้อยากแยก Managment fee กับ custodian fee ให้ได้
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.CASHTXTYPECODE = colCell
					case "TRADEDATE":
						// trade_date
						transaction.TRADEDATE = colCell
					case "SETTLEDATE":
						// settled_date
						transaction.SETTLEDATE = colCell
					case "TRADABLEDATE":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.TRADABLEDATE = colCell
					case "UNIT":
						// ไม่อยากให้มี unit เป็น 0 เลย
						if colCell == "0" || colCell == "" {
							colCell = "1"
						}
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
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
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
						if colCell == "Y" || colCell == "y" {
							colCell = "true"
						} else if colCell == "N" || colCell == "n" {
							colCell = "false"
						}
						transaction.ISEFFECTCASH = colCell
					case "CURRENCYID":
						transaction.CURRENCYID = colCell
					case "CURRENCYCODE":
						transaction.CURRENCYCODE = colCell
					case "BROKERID":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.BROKERID = colCell
					case "BROKERCODE":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.BROKERCODE = colCell
					case "COUNTERPARTYID":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.COUNTERPARTYID = colCell
					case "COUNTERPARTYCODE":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.COUNTERPARTYCODE = colCell
					case "TAXPAYERID":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
						transaction.TAXPAYERID = colCell
					case "TAXPAYERCODE":
						if colCell == "NULL" || colCell == "" {
							colCell = ""
						}
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
				}
			}
			if i != 0 {
				transactions = append(transactions, transaction)
			}
		}

		for _, txn := range transactions {
			fmt.Println(txn)
			Execute2(txn)
		}
	}

	// userId := 4 // mock ไปก่อน

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

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
