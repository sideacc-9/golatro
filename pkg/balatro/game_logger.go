package balatro

import (
	"bytes"
	"fmt"
	"time"
)

type LogType int

const (
	ROUND_STARTED LogType = iota
	HAND_PLAYED
	HAND_SCORED
	HAND_DISCARDED
	JOKER_ACTIVATED
	JOKER_REMOVED
	ENHANCEMENT_ACTIVATED
	EDITION_ACTIVATED
	PACK_OPENED
	UPGRADE_APPLIED
	MONEY_GAINED
	MONEY_SPENT
	CARD_ADDED
	CARD_DESTROYED
)

func (t LogType) String() string {
	switch t {
	case ROUND_STARTED:
		return "ROUND_STARTED"
	case HAND_PLAYED:
		return "HAND_PLAYED"
	case HAND_SCORED:
		return "HAND_SCORED"
	case HAND_DISCARDED:
		return "HAND_DISCARDED"
	case JOKER_ACTIVATED:
		return "JOKER_ACTIVATED"
	case JOKER_REMOVED:
		return "JOKER_REMOVED"
	case ENHANCEMENT_ACTIVATED:
		return "ENHANCEMENT_ACTIVATED"
	case EDITION_ACTIVATED:
		return "EDITION_ACTIVATED"
	case PACK_OPENED:
		return "PACK_OPENED"
	case UPGRADE_APPLIED:
		return "UPGRADE_APPLIED"
	case MONEY_GAINED:
		return "MONEY_GAINED"
	case MONEY_SPENT:
		return "MONEY_SPENT"
	case CARD_ADDED:
		return "CARD_ADDED"
	case CARD_DESTROYED:
		return "CARD_DESTROYED"
	}
	return ""
}

type GameLog struct {
	time    time.Time
	Type    LogType
	Details string
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}

func (log GameLog) String() string {
	return insertNth(fmt.Sprintf("[%v] %v - %v\n", log.time.Format("2006-01-02 15:04:05"), log.Type, log.Details), 80)
}

type GameLogger struct {
	log []GameLog
}

func (l *GameLogger) Add(t LogType, details string) {
	l.log = append(l.log, GameLog{time: time.Now(), Type: t, Details: details})
}

func (l *GameLogger) All() string {
	str := ""
	for _, log := range l.log {
		str += fmt.Sprintf("%v\n", log)
	}
	return str
}
