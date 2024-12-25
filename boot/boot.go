package boot

import (
	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/signailUtils"
	"time"
)

func Boot() {
	signailUtils.WaitExit(1 * time.Second)
}
