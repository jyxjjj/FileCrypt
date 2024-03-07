package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	DiffTimestamp int64 = 1702310400000
)

var (
	mutex          *sync.Mutex
	MachineID      int64 = 1
	LastSerialTime       = time.Now().UnixNano() / 1e6
	Serial         int64 = 0
	Socket         string
	Server         net.Listener
	Verbose        = false
)

func init() {
	log.Println("DESMG SerialNumber Server")
	mutex = new(sync.Mutex)
	MachineID = getMachineID()
	Socket = getSocket()
	isVerbose()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go handleSignal(sigCh)
}

func isVerbose() {
	for _, arg := range os.Args {
		if arg == "-v" {
			Verbose = true
			break
		}
	}
}

func getMachineID() int64 {
	machineIDEnv := os.Getenv("MACHINE_ID")
	if machineIDEnv == "" {
		machineIDEnv = "1"
	}
	machineID, err := strconv.ParseInt(machineIDEnv, 10, 64)
	if err != nil {
		return 1
	}
	return machineID
}

func getSocket() string {
	socket := os.Getenv("SOCKET_FILENAME")
	if socket == "" {
		socket = "DESMG-SerialNumber.sock"
	}
	return socket
}

func handleSignal(sigCh chan os.Signal) {
	_ = <-sigCh
	_ = Server.Close()
	_ = os.Remove(Socket)
	log.Println("Exiting...")
	os.Exit(0)
}

func main() {
	removeSocketFile()
	Server = startServer()
	for {
		conn, err := Server.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	serialNumber := GetSerialNumber()
	_, _ = conn.Write([]byte("DESMG-SNID::" + serialNumber))
	_ = conn.Close()
}

func startServer() net.Listener {
	server, err := net.Listen("unix", Socket)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Server started: %s", Socket)
	return server
}

func removeSocketFile() {
	err := os.Remove(Socket)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error: %v", err)
	}
}

func GetSerialNumber() string {
	mutex.Lock()
	defer mutex.Unlock()
	now := time.Now().UnixNano() / 1e6
	diff := now - DiffTimestamp
	if diff == LastSerialTime {
		Serial = (Serial + 1) & 4095
		if Serial == 0 {
			for diff <= LastSerialTime {
				now = time.Now().UnixNano() / 1e6
				diff = now - DiffTimestamp
			}
		}
	} else {
		Serial = 0
	}
	LastSerialTime = diff
	ID := 3<<60 | diff<<20 | MachineID<<12 | Serial
	if Verbose == true {
		bin := strconv.FormatInt(ID, 2)
		if len(bin) < 64 {
			for i := 0; i < 64-len(bin); i++ {
				bin = "0" + bin
			}
		}
		bin = bin[:1] + " " + bin[1:]
		bin = bin[:4] + " " + bin[4:]
		bin = bin[:45] + " " + bin[45:]
		bin = bin[:54] + " " + bin[54:]
		log.Printf("| %16d %4d %4d | %72s %32X(%d) %32d(%d)", diff, MachineID, Serial, bin, ID, len(fmt.Sprintf("%X", ID)), ID, len(fmt.Sprintf("%d", ID)))
	}
	SN := fmt.Sprintf("%X", ID)
	return SN
}
