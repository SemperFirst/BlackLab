package models 

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"

	"SafeDP/scanner/password_crack/logger"
	"SafeDP/scanner/password_crack/util/hash"
	"SafeDP/scanner/password_crack/vars"

	"github.com/patrickmn/go-cache"
)

func init() {
	gob.Register(Service{})
	gob.Register(ScanResult{})
}