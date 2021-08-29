package main

import (
	"crypto.bot/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
	if l, err := log.ParseLevel("debug"); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
		log.SetOutput(os.Stdout)
	}

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

func pingAlive(conn *websocket.Conn) {
	const (
		maxMessageSize = 512
		writeWait      = 10 * time.Second
		pongWait       = 60 * time.Second
		pingIteration  = (pongWait * 9) / 10
	)

	conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(string) error {
		log.Debug("pong received")

		conn.SetReadLimit(maxMessageSize)
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive")); err != nil {
			return
		}

		log.Debug("ping iteration")
		time.Sleep(pingIteration)
	}
}

var channel chan []byte = make(chan []byte, 100000)
var quantityInsert = 0

func insertClick() {

	connect, err := sql.Open("clickhouse", "xxxxxxxx")
	connect.SetMaxOpenConns(50)
	if err != nil {
		log.Error(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Error(err)
		}
		return
	}

	type ResultJson struct {
		Type          string `json:"e"`
		EventTime     int64  `json:"E"`
		Symbol        string `json:"s"`
		TradeId       int    `json:"t"`
		Price         string `json:"p"`
		Quality       string `json:"q"`
		BuyerOrderId  int64  `json:"b"`
		SellerOrderId int64  `json:"a"`
		TradeTime     int64  `json:"T"`
		Temp1         bool   `json:"m"`
		Temp2         bool   `json:"M"`
	}

	var dataFromJson ResultJson
	var tx, _ = connect.Begin()
	stmt, err := tx.Prepare("INSERT INTO trade_buffer (date_add, symbol, price, quality, event_time, trade_time, trade_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
	}

	log.Info("Inserting start...")

	count := 0
	for {
		quantityInsert++

		message := <-channel
		json.Unmarshal(message, &dataFromJson)

		price, _ := strconv.ParseFloat(dataFromJson.Price, 32)
		quality, _ := strconv.ParseFloat(dataFromJson.Quality, 32)
		//log.Debugf("%+v", dataFromJson)

		if _, err := stmt.Exec(
			time.Now().Format("2006-01-02 15:04:05"),
			dataFromJson.Symbol,
			price,
			quality,
			dataFromJson.EventTime,
			dataFromJson.TradeTime,
			dataFromJson.TradeId,
		); err != nil {
			log.Error(err)
		}

		if count >= 1000 {
			//if time.Now().Second()%2 == 0 {
			if err := tx.Commit(); err != nil {
				log.Error(err)
			}

			tx, _ = connect.Begin()
			stmt, err = tx.Prepare("INSERT INTO trade_buffer (date_add, symbol, price, quality, event_time, trade_time, trade_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Error(err)
			}

			fmt.Print(".")

			count = 0
		}
		count++
	}
}

func main() {
	utils.Test()

	go func() {
		return
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":3434", nil))
	}()

	go func() {
		for {
			qc := len(channel)
			if qc > 10000 {
				fmt.Print("\n")
				log.Info("queue: ", qc)
			}

			time.Sleep(5 * time.Second)
			//time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for {

			log.Info("quantity insert: ", quantityInsert)

			quantityInsert = 0

			time.Sleep(5 * time.Second)
			//time.Sleep(100 * time.Millisecond)
		}
	}()

	go insertClick()

	go checking()

	go clearClick()

	for {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com/ws/ethbtc@trade/ltcbtc@trade/bnbbtc@trade/neobtc@trade/bccbtc@trade/gasbtc@trade/hsrbtc@trade/mcobtc@trade/wtcbtc@trade/lrcbtc@trade/qtumbtc@trade/yoyobtc@trade/omgbtc@trade/zrxbtc@trade/stratbtc@trade/snglsbtc@trade/bqxbtc@trade/kncbtc@trade/funbtc@trade/snmbtc@trade/iotabtc@trade/linkbtc@trade/xvgbtc@trade/saltbtc@trade/mdabtc@trade/mtlbtc@trade/subbtc@trade/eosbtc@trade/sntbtc@trade/etcbtc@trade/mthbtc@trade/engbtc@trade/dntbtc@trade/zecbtc@trade/bntbtc@trade/astbtc@trade/dashbtc@trade/oaxbtc@trade/icnbtc@trade/btgbtc@trade/evxbtc@trade/reqbtc@trade/vibbtc@trade/trxbtc@trade/powrbtc@trade/arkbtc@trade/xrpbtc@trade/modbtc@trade/enjbtc@trade/storjbtc@trade/venbtc@trade/kmdbtc@trade/rcnbtc@trade/nulsbtc@trade/rdnbtc@trade/xmrbtc@trade/dltbtc@trade/ambbtc@trade/batbtc@trade/bcptbtc@trade/arnbtc@trade/gvtbtc@trade/cdtbtc@trade/gxsbtc@trade/poebtc@trade/qspbtc@trade/btsbtc@trade/xzcbtc@trade/lskbtc@trade/tntbtc@trade/fuelbtc@trade/manabtc@trade/bcdbtc@trade/dgdbtc@trade/adxbtc@trade/adabtc@trade/pptbtc@trade/cmtbtc@trade/xlmbtc@trade/cndbtc@trade/lendbtc@trade/wabibtc@trade/tnbbtc@trade/wavesbtc@trade/gtobtc@trade/icxbtc@trade/ostbtc@trade/elfbtc@trade/aionbtc@trade/neblbtc@trade/brdbtc@trade/edobtc@trade/wingsbtc@trade/navbtc@trade/lunbtc@trade/trigbtc@trade/appcbtc@trade/vibebtc@trade/rlcbtc@trade/insbtc@trade/pivxbtc@trade/iostbtc@trade/chatbtc@trade/steembtc@trade/nanobtc@trade/viabtc@trade/blzbtc@trade/aebtc@trade/rpxbtc@trade/ncashbtc@trade/poabtc@trade/zilbtc@trade/ontbtc@trade/stormbtc@trade/xembtc@trade/wanbtc@trade/wprbtc@trade/qlcbtc@trade/sysbtc@trade/grsbtc@trade/cloakbtc@trade/gntbtc@trade/loombtc@trade/bcnbtc@trade/repbtc@trade/tusdbtc@trade/zenbtc@trade/skybtc@trade/cvcbtc@trade/thetabtc@trade/iotxbtc@trade/qkcbtc@trade/agibtc@trade/nxsbtc@trade/databtc@trade/scbtc@trade/npxsbtc@trade/keybtc@trade/nasbtc@trade/mftbtc@trade/dentbtc@trade/ardrbtc@trade/hotbtc@trade/vetbtc@trade/dockbtc@trade/polybtc@trade/phxbtc@trade/hcbtc@trade/gobtc@trade/paxbtc@trade/rvnbtc@trade/dcrbtc@trade/mithbtc@trade/bchabcbtc@trade/bchsvbtc@trade/renbtc@trade/bttbtc@trade/ongbtc@trade/fetbtc@trade/celrbtc@trade/maticbtc@trade/atombtc@trade/phbbtc@trade/tfuelbtc@trade/onebtc@trade/ftmbtc@trade/btcbbtc@trade/algobtc@trade/erdbtc@trade/dogebtc@trade/duskbtc@trade/ankrbtc@trade/winbtc@trade/cosbtc@trade/cocosbtc@trade/tomobtc@trade/perlbtc@trade/chzbtc@trade/bandbtc@trade/beambtc@trade/xtzbtc@trade/hbarbtc@trade/nknbtc@trade/stxbtc@trade/kavabtc@trade/arpabtc@trade/ctxcbtc@trade/bchbtc@trade/troybtc@trade/vitebtc@trade/fttbtc@trade/ognbtc@trade/drepbtc@trade/tctbtc@trade/wrxbtc@trade/ltobtc@trade/mblbtc@trade/cotibtc@trade/stptbtc@trade/solbtc@trade/ctsibtc@trade/hivebtc@trade/chrbtc@trade/mdtbtc@trade/stmxbtc@trade/pntbtc@trade/dgbbtc@trade/compbtc@trade/sxpbtc@trade/snxbtc@trade/irisbtc@trade/mkrbtc@trade/daibtc@trade/runebtc@trade/fiobtc@trade/avabtc@trade/balbtc@trade/yfibtc@trade/jstbtc@trade/srmbtc@trade/antbtc@trade/crvbtc@trade/sandbtc@trade/oceanbtc@trade/nmrbtc@trade/dotbtc@trade/lunabtc@trade/idexbtc@trade/rsrbtc@trade/paxgbtc@trade/wnxmbtc@trade/trbbtc@trade/bzrxbtc@trade/wbtcbtc@trade/sushibtc@trade/yfiibtc@trade/ksmbtc@trade/egldbtc@trade/diabtc@trade/umabtc@trade/belbtc@trade/wingbtc@trade/unibtc@trade/nbsbtc@trade/oxtbtc@trade/sunbtc@trade/avaxbtc@trade/hntbtc@trade/flmbtc@trade/scrtbtc@trade/ornbtc@trade/utkbtc@trade/xvsbtc@trade/alphabtc@trade/vidtbtc@trade/aavebtc@trade/nearbtc@trade/filbtc@trade/injbtc@trade/aergobtc@trade/audiobtc@trade/ctkbtc@trade/botbtc@trade/akrobtc@trade/axsbtc@trade/hardbtc@trade/renbtcbtc@trade/straxbtc@trade/forbtc@trade/unfibtc@trade/rosebtc@trade/sklbtc@trade/susdbtc@trade/glmbtc@trade/grtbtc@trade/juvbtc@trade/psgbtc@trade/1inchbtc@trade/reefbtc@trade/ogbtc@trade/atmbtc@trade/asrbtc@trade/celobtc@trade/rifbtc@trade/btcstbtc@trade/trubtc@trade/ckbbtc@trade/twtbtc@trade/firobtc@trade/litbtc@trade/sfpbtc@trade/fxsbtc@trade/dodobtc@trade/frontbtc@trade/easybtc@trade/cakebtc@trade/acmbtc@trade/auctionbtc@trade/phabtc@trade/tvkbtc@trade/badgerbtc@trade/fisbtc@trade/ombtc@trade/pondbtc@trade/degobtc@trade/alicebtc@trade/linabtc@trade/perpbtc@trade/rampbtc@trade/superbtc@trade/cfxbtc@trade/epsbtc@trade/autobtc@trade/tkobtc@trade/tlmbtc@trade/mirbtc@trade/barbtc@trade/forthbtc@trade/ezbtc@trade/icpbtc@trade/arbtc@trade/polsbtc@trade/mdxbtc@trade/lptbtc@trade/agixbtc@trade/nubtc@trade/atabtc@trade/gtcbtc@trade/tornbtc@trade/bakebtc@trade/keepbtc@trade/klaybtc@trade/bondbtc@trade/mlnbtc@trade/btcusdt@trade/ethusdt@trade/bnbusdt@trade/bccusdt@trade/neousdt@trade/ltcusdt@trade/qtumusdt@trade/adausdt@trade/xrpusdt@trade/eosusdt@trade/tusdusdt@trade/iotausdt@trade/xlmusdt@trade/ontusdt@trade/trxusdt@trade/etcusdt@trade/icxusdt@trade/venusdt@trade/nulsusdt@trade/vetusdt@trade/paxusdt@trade/bchabcusdt@trade/bchsvusdt@trade/usdcusdt@trade/linkusdt@trade/wavesusdt@trade/bttusdt@trade/usdsusdt@trade/ongusdt@trade/hotusdt@trade/zilusdt@trade/zrxusdt@trade/fetusdt@trade/batusdt@trade/xmrusdt@trade/zecusdt@trade/iostusdt@trade/celrusdt@trade/dashusdt@trade/nanousdt@trade/omgusdt@trade/thetausdt@trade/enjusdt@trade/mithusdt@trade/maticusdt@trade/atomusdt@trade/tfuelusdt@trade/oneusdt@trade/ftmusdt@trade/algousdt@trade/usdsbusdt@trade/gtousdt@trade/erdusdt@trade/dogeusdt@trade/duskusdt@trade/ankrusdt@trade/winusdt@trade/cosusdt@trade/npxsusdt@trade/cocosusdt@trade/mtlusdt@trade/tomousdt@trade/perlusdt@trade/dentusdt@trade/mftusdt@trade/keyusdt@trade/stormusdt@trade/dockusdt@trade/wanusdt@trade/funusdt@trade/cvcusdt@trade/chzusdt@trade/bandusdt@trade/busdusdt@trade/beamusdt@trade/xtzusdt@trade/renusdt@trade/rvnusdt@trade/hcusdt@trade/hbarusdt@trade/nknusdt@trade/stxusdt@trade/kavausdt@trade/arpausdt@trade/iotxusdt@trade/rlcusdt@trade/mcousdt@trade/ctxcusdt@trade/bchusdt@trade/troyusdt@trade/viteusdt@trade/fttusdt@trade/eurusdt@trade/ognusdt@trade/drepusdt@trade/bullusdt@trade/bearusdt@trade/ethbullusdt@trade/ethbearusdt@trade/tctusdt@trade/wrxusdt@trade/btsusdt@trade/lskusdt@trade/bntusdt@trade/ltousdt@trade/eosbullusdt@trade/eosbearusdt@trade/xrpbullusdt@trade/xrpbearusdt@trade/stratusdt@trade/aionusdt@trade/mblusdt@trade/cotiusdt@trade/bnbbullusdt@trade/bnbbearusdt@trade/stptusdt@trade/wtcusdt@trade/datausdt@trade/xzcusdt@trade/solusdt@trade/ctsiusdt@trade/hiveusdt@trade/chrusdt@trade/btcupusdt@trade/btcdownusdt@trade/gxsusdt@trade/ardrusdt@trade/lendusdt@trade/mdtusdt@trade/stmxusdt@trade/kncusdt@trade/repusdt@trade/lrcusdt@trade/pntusdt@trade/compusdt@trade/bkrwusdt@trade/scusdt@trade/zenusdt@trade/snxusdt@trade/ethupusdt@trade/ethdownusdt@trade/adaupusdt@trade/adadownusdt@trade/linkupusdt@trade/linkdownusdt@trade/vthousdt@trade/dgbusdt@trade/gbpusdt@trade/sxpusdt@trade/mkrusdt@trade/daiusdt@trade/dcrusdt@trade/storjusdt@trade/bnbupusdt@trade/bnbdownusdt@trade/xtzupusdt@trade/xtzdownusdt@trade/manausdt@trade/audusdt@trade/yfiusdt@trade/balusdt@trade/blzusdt@trade/irisusdt@trade/kmdusdt@trade/jstusdt@trade/srmusdt@trade/antusdt@trade/crvusdt@trade/sandusdt@trade/oceanusdt@trade/nmrusdt@trade/dotusdt@trade/lunausdt@trade/rsrusdt@trade/paxgusdt@trade/wnxmusdt@trade/trbusdt@trade/bzrxusdt@trade/sushiusdt@trade/yfiiusdt@trade/ksmusdt@trade/egldusdt@trade/diausdt@trade/runeusdt@trade/fiousdt@trade/umausdt@trade/eosupusdt@trade/eosdownusdt@trade/trxupusdt@trade/trxdownusdt@trade/xrpupusdt@trade/xrpdownusdt@trade/dotupusdt@trade/dotdownusdt@trade/belusdt@trade/wingusdt@trade/ltcupusdt@trade/ltcdownusdt@trade/uniusdt@trade/nbsusdt@trade/oxtusdt@trade/sunusdt@trade/avaxusdt@trade/hntusdt@trade/flmusdt@trade/uniupusdt@trade/unidownusdt@trade/ornusdt@trade/utkusdt@trade/xvsusdt@trade/alphausdt@trade/aaveusdt@trade/nearusdt@trade/sxpupusdt@trade/sxpdownusdt@trade/filusdt@trade/filupusdt@trade/fildownusdt@trade/yfiupusdt@trade/yfidownusdt@trade/injusdt@trade/audiousdt@trade/ctkusdt@trade/bchupusdt@trade/bchdownusdt@trade/akrousdt@trade/axsusdt@trade/hardusdt@trade/dntusdt@trade/straxusdt@trade/unfiusdt@trade/roseusdt@trade/avausdt@trade/xemusdt@trade/aaveupusdt@trade/aavedownusdt@trade/sklusdt@trade/susdusdt@trade/sushiupusdt@trade/sushidownusdt@trade/xlmupusdt@trade/xlmdownusdt@trade/grtusdt@trade/juvusdt@trade/psgusdt@trade/1inchusdt@trade/reefusdt@trade/ogusdt@trade/atmusdt@trade/asrusdt@trade/celousdt@trade/rifusdt@trade/btcstusdt@trade/truusdt@trade/ckbusdt@trade/twtusdt@trade/firousdt@trade/litusdt@trade/sfpusdt@trade/dodousdt@trade/cakeusdt@trade/acmusdt@trade/badgerusdt@trade/fisusdt@trade/omusdt@trade/pondusdt@trade/degousdt@trade/aliceusdt@trade/linausdt@trade/perpusdt@trade/rampusdt@trade/superusdt@trade/cfxusdt@trade/epsusdt@trade/autousdt@trade/tkousdt@trade/pundixusdt@trade/tlmusdt@trade/1inchupusdt@trade/1inchdownusdt@trade/btgusdt@trade/mirusdt@trade/barusdt@trade/forthusdt@trade/bakeusdt@trade/burgerusdt@trade/slpusdt@trade/shibusdt@trade/icpusdt@trade/arusdt@trade/polsusdt@trade/mdxusdt@trade/maskusdt@trade/lptusdt@trade/nuusdt@trade/xvgusdt@trade/atausdt@trade/gtcusdt@trade/tornusdt@trade/keepusdt@trade/ernusdt@trade/klayusdt@trade/phausdt@trade/bondusdt@trade/mlnusdt@trade/qtumeth@trade/eoseth@trade/snteth@trade/bnteth@trade/bnbeth@trade/oaxeth@trade/dnteth@trade/mcoeth@trade/icneth@trade/wtceth@trade/lrceth@trade/omgeth@trade/zrxeth@trade/strateth@trade/snglseth@trade/bqxeth@trade/knceth@trade/funeth@trade/snmeth@trade/neoeth@trade/iotaeth@trade/linketh@trade/xvgeth@trade/salteth@trade/mdaeth@trade/mtleth@trade/subeth@trade/etceth@trade/mtheth@trade/engeth@trade/zeceth@trade/asteth@trade/dasheth@trade/btgeth@trade/evxeth@trade/reqeth@trade/vibeth@trade/hsreth@trade/trxeth@trade/powreth@trade/arketh@trade/yoyoeth@trade/xrpeth@trade/modeth@trade/enjeth@trade/storjeth@trade/veneth@trade/kmdeth@trade/rcneth@trade/nulseth@trade/rdneth@trade/xmreth@trade/dlteth@trade/ambeth@trade/bcceth@trade/bateth@trade/bcpteth@trade/arneth@trade/gvteth@trade/cdteth@trade/gxseth@trade/poeeth@trade/qspeth@trade/btseth@trade/xzceth@trade/lsketh@trade/tnteth@trade/fueleth@trade/manaeth@trade/bcdeth@trade/dgdeth@trade/adxeth@trade/adaeth@trade/ppteth@trade/cmteth@trade/xlmeth@trade/cndeth@trade/lendeth@trade/wabieth@trade/ltceth@trade/tnbeth@trade/waveseth@trade/gtoeth@trade/icxeth@trade/osteth@trade/elfeth@trade/aioneth@trade/nebleth@trade/brdeth@trade/edoeth@trade/wingseth@trade/naveth@trade/luneth@trade/trigeth@trade/appceth@trade/vibeeth@trade/rlceth@trade/inseth@trade/pivxeth@trade/iosteth@trade/chateth@trade/steemeth@trade/nanoeth@trade/viaeth@trade/blzeth@trade/aeeth@trade/rpxeth@trade/ncasheth@trade/poaeth@trade/zileth@trade/onteth@trade/stormeth@trade/xemeth@trade/waneth@trade/wpreth@trade/qlceth@trade/syseth@trade/grseth@trade/cloaketh@trade/gnteth@trade/loometh@trade/bcneth@trade/repeth@trade/tusdeth@trade/zeneth@trade/skyeth@trade/cvceth@trade/thetaeth@trade/iotxeth@trade/qkceth@trade/agieth@trade/nxseth@trade/dataeth@trade/sceth@trade/npxseth@trade/keyeth@trade/naseth@trade/mfteth@trade/denteth@trade/ardreth@trade/hoteth@trade/veteth@trade/docketh@trade/phxeth@trade/hceth@trade/paxeth@trade/stmxeth@trade/wbtceth@trade/scrteth@trade/aaveeth@trade/easyeth@trade/renbtceth@trade/slpeth@trade/cvpeth@trade/straxeth@trade/fronteth@trade/hegiceth@trade/susdeth@trade/covereth@trade/glmeth@trade/ghsteth@trade/dfeth@trade/grteth@trade/dexeeth@trade/firoeth@trade/betheth@trade/proseth@trade/ufteth@trade/pundixeth@trade/ezeth@trade"), nil)
		if err != nil {
			log.Error(err)
			time.Sleep(1 * time.Second)

			continue
		} else {
			go pingAlive(conn)

			for {
				_, message, readErr := conn.ReadMessage()
				//fmt.Printf("%+v", message)
				//return

				if readErr != nil {
					log.Error(readErr)
					time.Sleep(1 * time.Second)
					break
				} else {
					channel <- message
				}
			}
		}
	}
}

func clearClick() {
	for {
		connect, err := sql.Open("clickhouse", "xxxxxxxx")
		connect.SetMaxOpenConns(50)
		if err != nil {
			log.Error(err)
		}
		if err := connect.Ping(); err != nil {
			if exception, ok := err.(*clickhouse.Exception); ok {
				log.Info("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			} else {
				log.Error(err)
			}
			return
		}

		connect.Exec("ALTER TABLE trade DELETE WHERE date_index BETWEEN '2021-01-01' AND NOW() - INTERVAL 4 DAY")
		connect.Close()

		log.Info("Cleaning clickhouse...")

		time.Sleep(3600 * time.Second)
	}
}

type rate struct {
	couple string
	way    string
	price  float64
	chatId int32
}

func checking() {
	connect, err := sql.Open("clickhouse", "xxxxxxxx")
	connect.SetMaxOpenConns(50)
	if err != nil {
		log.Error(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Info("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Error(err)
		}
		return
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "/app/crypto_bot.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for {
		rows, err := db.Query("select couple, way, price, chat_id from rates")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		rates := []rate{}

		for rows.Next() {
			r := rate{}
			err := rows.Scan(&r.couple, &r.way, &r.price, &r.chatId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			rates = append(rates, r)
		}
		for _, r := range rates {
			//log.Print(r.couple, r.way, r.price, r.chat_id)

			rows, err := connect.Query("SELECT date_add, price FROM trade WHERE date_index BETWEEN NOW() - INTERVAL 24 HOUR AND NOW() AND symbol = '" + strings.ToUpper(reg.ReplaceAllString(r.couple, "")) + "' ORDER BY date_add DESC LIMIT 1")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var (
					dateAdd  time.Time
					floatOut float64
				)
				if err := rows.Scan(&dateAdd, &floatOut); err != nil {
					log.Error(err)
				}

				log.Printf("date_add: %s, price: %f, find price: %f, symbol: %s, chat id: %d", dateAdd, floatOut, r.price, r.couple, r.chatId)

				if r.way == "up" {
					if floatOut >= r.price {
						send(r.chatId, "Binance notify: "+r.couple+" up to "+fmt.Sprintf("%f", r.price))

						log.Info("send")
						return
					}
				} else if r.way == "down" {
					if floatOut <= r.price {
						send(r.chatId, "Binance notify: "+r.couple+" up to "+fmt.Sprintf("%f", r.price))

						log.Info("send")
						return
					}
				}
			}

			if err := rows.Err(); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func send(chatId int32, text string) {
	resp, _ := http.Get("https://api.telegram.org/xxxxxxxx/sendMessage?chat_id=" + strconv.FormatInt(int64(chatId), 10) + "&text=" + text)
	fmt.Printf("%s", resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
	stringSlice := strings.Split(r.URL.Path[1:], "/")

	go runObserver(stringSlice)

	log.Info(w, "ok")
}

func runObserver(data []string) {

	log.Info(data)

	connect, err := sql.Open("clickhouse", "xxxxxxxx")
	connect.SetMaxOpenConns(50)
	if err != nil {
		log.Error(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Info("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Error(err)
		}
		return
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	for {
		rows, err := connect.Query("SELECT date_add, price FROM trade WHERE date_index BETWEEN NOW() - INTERVAL 24 HOUR AND NOW() AND symbol = '" + strings.ToUpper(reg.ReplaceAllString(data[0], "")) + "' ORDER BY date_add DESC LIMIT 1")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				dateAdd  time.Time
				floatOut float64
			)
			if err := rows.Scan(&dateAdd, &floatOut); err != nil {
				log.Error(err)
			}

			floatUser, _ := strconv.ParseFloat(data[2], 42)

			log.Printf("date_add: %s, price: %d, find price: %d, symbol: %s, chat id: %s", dateAdd, floatOut, floatUser, data[0], data[3])

			if data[1] == "up" {
				if floatOut >= floatUser {
					//send(data[3], "Binance notify: "+data[0]+" up to "+data[2])

					log.Info("send")
					return
				}
			} else if data[1] == "down" {
				if floatOut <= floatUser {
					//send(data[3], "Binance notify: "+data[0]+" up to "+data[2])

					log.Info("send")
					return
				}
			}
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)
	}
}
