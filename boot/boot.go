package boot

import (
	"time"

	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/signailUtils"
)

func Boot() {
	signailUtils.WaitExit(1 * time.Second)
}
