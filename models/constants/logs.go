package constants

import "github.com/rs/zerolog"

const (
	LogFileName      = "fileName"
	LogCorrelationID = "correlationID"
	LogAnkamaID      = "ankamaID"
	LogQueryID       = "queryID"
	LogQueryType     = "queryType"
	LogItemType      = "itemType"
	LogEntityCount   = "entityCount"
	LogKey           = "key"

	LogLevelFallback = zerolog.InfoLevel
)
