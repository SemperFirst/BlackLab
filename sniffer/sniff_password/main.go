package main
import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device = "enp0s5"
	snapshotLength int32 = 1024
	promiscuous = false
	timeout = 5 * time.Second
	handle *pcap.Handle
	err error

	filter = "(tcp and dst port 21) or (tcp and dst port 80) or (tcp and dst port 25) or (tcp and dst port 110)"
	userList = []string{"user","username","login","login_user","manager","user_name","usr"}
	passList = []string{"pass","password","login_pass","pwd","passwd"}
)