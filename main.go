// Downloads torrents from the command-line.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/anacrolix/tagflag"
	log "github.com/inconshreveable/log15"
	optimizedconn "github.com/johannwagner/scion-optimized-connection/pkg"
	"github.com/lucas-clemente/quic-go"
	"github.com/netsec-ethz/scion-apps/pkg/appnet"
	"github.com/netsec-ethz/scion-apps/pkg/appnet/appquic"
	// "github.com/netsec-ethz/scion-apps/pkg/appnet"
)

var (
	// Don't verify the server's cert, as we are not using the TLS PKI.
	TLSCfg = &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"h3"}}
)

var flags = struct {
	IsServer  bool
	NumCons   int
	StartPort int
	Local     string
	Remote    string
	tagflag.StartPos
}{
	IsServer:  true,
	NumCons:   1,
	StartPort: 43000,
	Local:     "127.0.0.1",
	Remote:    "127.0.0.1",
}

const (
	PacketSize int64 = 1500
	NumPackets int64 = 4000000
)

func LogFatal(msg string, a ...interface{}) {
	log.Crit(msg, a...)
	os.Exit(1)
}

func Check(e error) {
	if e != nil {
		LogFatal("Fatal error. Exiting.", "err", e)
	}
}

func main() {
	if err := mainErr(); err != nil {
		log.Info("error in main: %v", err)
		os.Exit(1)
	}
}

func mainErr() error {
	tagflag.Parse(&flags)

	certs := appquic.GetDummyTLSCerts()
	TLSCfg.Certificates = certs

	startPort := uint16(flags.StartPort)
	var i uint16
	i = 0
	var wg sync.WaitGroup
	for i < uint16(flags.NumCons) {
		go func(wg *sync.WaitGroup, startPort uint16, i uint16) {
			defer wg.Done()
			if flags.IsServer {
				runServer(flags.Local, startPort+i)
			} else {
				runClient(flags.Local, flags.Remote, int(startPort+i))
			}
		}(&wg, startPort, i)
		i++
	}
	wg.Wait()
	time.Sleep(time.Minute * 5)

	return nil
}

func runClient(local, remote string, startPort int) {
	localAddr, err := appnet.ResolveUDPAddr(local)
	Check(err)
	udpAddr := net.UDPAddr{
		IP:   localAddr.Host.IP,
		Port: startPort,
	}
	remoteAddr, err := appnet.ResolveUDPAddr(remote)
	Check(err)
	remoteAddr.Host.Port = startPort
	fmt.Printf("Dial from %s to %s\n", udpAddr.String(), remoteAddr.String())
	sconn, err := optimizedconn.Dial(&udpAddr, remoteAddr)
	Check(err)
	host := appnet.MangleSCIONAddr(remote)
	// DCConn, err := appnet.DialAddr(serverAddr)

	sess, err := quic.Dial(sconn, remoteAddr, host, TLSCfg, &quic.Config{
		KeepAlive: true,
	})
	Check(err)
	fmt.Println("GOT SESSION")
	conn, sessErr := sess.OpenStreamSync(context.Background())
	Check(sessErr)
	fmt.Println("GOT STREAM")
	// clientCCAddr := CCConn.LocalAddr().(*net.UDPAddr)
	sb := make([]byte, PacketSize)
	// Data channel connection
	// DCConn, err := appnet.DefNetwork().Dial(
	//	context.TODO(), "udp", clientCCAddr, serverAddr, addr.SvcNone)
	// Check(err)
	var i int64 = 0
	start := time.Now()
	for i < NumPackets {
		// Compute how long to wait

		// PrgFill(bwp.PrgKey, int(i*bwp.PacketSize), sb)
		// Place packet number at the beginning of the packet, overwriting some PRG data
		_, err := conn.Write(sb)
		Check(err)
		i++
	}
	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s\n", elapsed)
}

func runServer(local string, port uint16) error {
	serverAddr, err := appnet.ResolveUDPAddr(local)
	Check(err)
	fmt.Printf("Listen on Port %d", port)
	udpAddr := net.UDPAddr{
		IP:   serverAddr.Host.IP,
		Port: int(port),
	}
	conn, err := optimizedconn.Listen(&udpAddr)
	Check(err)

	qConn, listenErr := quic.Listen(conn, TLSCfg, &quic.Config{KeepAlive: true})
	Check(listenErr)

	var numPacketsReceived int64
	numPacketsReceived = 0
	recBuf := make([]byte, PacketSize+1000)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Printf("Received %d packets\n", numPacketsReceived)
	}()
	x, err := qConn.Accept(context.Background())
	Check(err)
	DCConn, err := x.AcceptStream(context.Background())
	Check(err)
	for numPacketsReceived < NumPackets {
		_, err := DCConn.Read(recBuf)

		// Ignore errors, todo: detect type of error and quit if it was because of a SetReadDeadline
		if err != nil {
			fmt.Println(err)
			continue
		}
		/*
			if int64(n) != PacketSize {
				// The packet has incorrect size, do not count as a correct packet
				// fmt.Println("Incorrect size.", n, "bytes instead of", PacketSize)
				continue
			}*/
		// fmt.Printf("Read packet of size %d\n", n)
		numPacketsReceived++
		// fmt.Printf("Received %d packets\n", numPacketsReceived)
	}

	fmt.Printf("Received all packets")
	return nil
}
