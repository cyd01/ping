package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-ping/ping"
)

var usagePingTxt = `
Usage:

    ping [-count count] [-interval interval] [-timeout timeout] [--debug] host

Examples:

    # ping google continuously
    ping www.google.com

    # ping google 5 times
    ping -count 5 www.google.com

    # ping google 5 times at 500ms intervals
    ping -count 5 -interval 500ms www.google.com

    # ping google for 10 seconds
    ping -timeout 10s www.google.com

    # Send a privileged raw ICMP ping
    sudo ping --debug www.google.com
`

func usagePing() {
	fmt.Print(usagePingTxt)
}

func mainPing( count int, interval time.Duration, timeout time.Duration, debug bool, host string ) {
	/*timeout := flag.Duration("t", time.Second*100000, "")
	interval := flag.Duration("i", time.Second, "")
	count := flag.Int("c", -1, "")
	privileged := flag.Bool("privileged", false, "")
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
*/
	//host := flag.Arg(0)
	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = count
	pinger.Interval = interval
	pinger.Timeout = timeout
	pinger.SetPrivileged(debug)

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		fmt.Printf("Failed to ping target host: %s", err)
	}
}

var (
    count          = flag.Int(     "count",       -1,                                     "number of ping")
    interval       = flag.Duration("interval",    time.Second,                            "interval of ping")

)

func main() {
	flag.Usage = usagePing
	flag.Parse()
	if( len(os.Args)>1 ) {
	    mainPing( *count, *interval, time.Second*100000, false, os.Args[1] )
	} else {
		usagePing()
	}
}