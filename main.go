package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "msg"

	k := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		GroupID: "clear-old-cache",
		GroupTopics: []string{
			"debezium.chii.bangumi.chii_subject_interests",
			"debezium.chii.bangumi.chii_episodes",
		},
	})
	mc := memcache.New(os.Getenv("MEMCACHED"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		msg, err := k.ReadMessage(ctx)

		if err != nil {
			log.Err(err).Msg("failed to read kafka message")
			continue
		}

		if len(msg.Value) == 0 {
			continue
		}

		if strings.HasSuffix(msg.Topic, ".chii_episodes") {
			var m KafkaMessageValue[ChiiEpisode]
			err = json.Unmarshal(msg.Value, &m)
			if err != nil {
				log.Err(err).Bytes("value", msg.Value).Msg("failed to decode kafka message as json")
				continue
			}

			if (m.Op != OpUpdate) && (m.Op != OpCreate) {
				continue
			}

			_ = mc.Delete(fmt.Sprintf("subject_eps_%d", m.After.SubjectID))
			continue
		}

		if strings.HasSuffix(msg.Topic, ".chii_subject_interests") {
			var m KafkaMessageValue[ChiiInterest]
			err = json.Unmarshal(msg.Value, &m)
			if err != nil {
				log.Err(err).Bytes("value", msg.Value).Msg("failed to decode kafka message as json")
				continue
			}

			if m.Op != OpUpdate {
				continue
			}

			_ = mc.Delete(fmt.Sprintf("prg_ep_status_%d", m.After.InterestUID))
			_ = mc.Delete(fmt.Sprintf("prg_watching_v3_%d", m.After.InterestUID))
			continue
		}
	}
}

const (
	OpCreate = "c"
	OpUpdate = "u"
	OpDelete = "d"
)

type Source struct {
	TsMs int64 `json:"ts_ms"`
}

func (s Source) Timestamp() time.Time {
	return time.UnixMicro(s.TsMs)
}

type KafkaMessageValue[T any] struct {
	Before *T     `json:"before"`
	After  *T     `json:"after"`
	Op     string `json:"op"`
	Source Source `json:"source"`
}

type ChiiEpisode struct {
	SubjectID uint64 `json:"ep_subject_id"`
}

type ChiiInterest struct {
	InterestID              int    `json:"interest_id"`
	InterestUID             int    `json:"interest_uid"`
	InterestSubjectID       int    `json:"interest_subject_id"`
	InterestSubjectType     int    `json:"interest_subject_type"`
	InterestRate            int    `json:"interest_rate"`
	InterestType            int    `json:"interest_type"`
	InterestHasComment      int    `json:"interest_has_comment"`
	InterestComment         string `json:"interest_comment"`
	InterestTag             string `json:"interest_tag"`
	InterestEpStatus        int    `json:"interest_ep_status"`
	InterestVolStatus       int    `json:"interest_vol_status"`
	InterestWishDateline    int    `json:"interest_wish_dateline"`
	InterestDoingDateline   int    `json:"interest_doing_dateline"`
	InterestCollectDateline int    `json:"interest_collect_dateline"`
	InterestOnHoldDateline  int    `json:"interest_on_hold_dateline"`
	InterestDroppedDateline int    `json:"interest_dropped_dateline"`
	InterestCreateIP        string `json:"interest_create_ip"`
	InterestLasttouchIP     string `json:"interest_lasttouch_ip"`
	InterestLasttouch       int    `json:"interest_lasttouch"`
	InterestPrivate         int    `json:"interest_private"`
}
