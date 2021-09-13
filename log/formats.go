package log

import (
	log "github.com/ron96G/log15"
)

func LogfmtFormat() log.Format {
	return log.LogfmtFormat()
}

func JsonFormat() log.Format {
	return log.JsonFormatEx(false, true)
}
