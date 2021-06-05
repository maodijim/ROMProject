package esClient

import (
	Cmd "ROMProject/Cmds"
	"fmt"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"time"
)

type ExchangeTemplate struct {
	ServerId     uint32           `json:"server_id"`
	ItemId       uint32           `json:"item_id"`
	ItemName     string           `json:"item_name"`
	ItemPrice    uint64           `json:"item_price"`
	ItemCategory uint32           `json:"item_category"`
	ItemEnhance  *Cmd.EnchantData `json:"item_enhance"`
	ItemRefineLv uint32           `json:"item_refine_lv"`
	IsPub        bool             `json:"is_pub"`
	ExpireTime   uint32           `json:"expire_time"`
	Count        uint32           `json:"count"`
	BuyerCount   uint32           `json:"buyer_count"`
	TimeStamp    time.Time        `json:"timestamp"`
	TradeType    Cmd.ETradeType   `json:"trade_type"`
	IsDamage     bool             `json:"is_damage"`
	Guid         string           `json:"guid"`
}

func (e *ExchangeTemplate) GetIndexName() string {
	return fmt.Sprintf("rom-exchange-%s", time.Now().Format("2006-01-02"))
}

func NewEsClient(urls []string) *elastic.Client {
	log.Infof("Establishing connection to elastic search: %v", urls)
	client, err := elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		log.Errorf("failed to connect to elasticsearch: %s", err)
		log.Infof("Exiting")
		log.Exit(1)
	}
	return client
}
