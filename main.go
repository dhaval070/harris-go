package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

func main() {
	server, err := net.Listen("tcp", ":30000")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening")

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("connected")
		processData(conn)

		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}
}

func processData(conn io.Reader) {
	buff := bufio.NewScanner(conn)

	buff.Split(split)

	log.Println("scanning")

	for buff.Scan() {
		line := buff.Text()
		log.Println(line)

		if line == "" {
			break
		}

		stats := parseLine(line)
		d, err := json.Marshal(stats)
		if err != nil {
			panic(err)
		}
		log.Println(string(d))
	}
}

func parseLine(line string) *Stats {
	pos := strings.Index(line, "@@")
	if pos == -1 {
		return nil
	}
	var stats Stats

	line = line[pos:]
	parts := strings.Split(line, "/")

	for i, s := range parts {
		parts[i] = reverse(s)
	}
	log.Println(parts)

	var flag string
	if len(parts) > 13 {
		flag = parts[13]
	}

	re := regexp.MustCompile("([0-9]{2})([0-9]{2})")
	match := re.FindStringSubmatch(parts[3])

	if flag == "T" {
		stats.Clock = match[1] + ":" + match[2]
		stats.ClockMode = ":"
		stats.ClockMin = match[1]
		stats.ClockSec = match[2]

	} else {
		stats.Clock = match[1] + "." + match[2][0:1]
		stats.ClockMode = "."
		stats.ClockMin = match[1]
		stats.ClockSec = match[2][0:1]
	}
	stats.ClockStatus = "Running" // TODO: make it dynamic from Clock data
	stats.Period = parts[6]
	stats.HomeScore = parts[4]
	stats.GuestScore = parts[5]
	stats.GamePeriod = parts[6]
	stats.HomeShots = parts[7]
	stats.GuestShots = parts[8]
	stats.Bot = parts[1]
	stats.HPlayer1 = parts[9][0:2]
	stats.HPlayer1Clock = parts[9][2:3] + ":" + parts[9][3:]
	stats.HPlayer1ClockMin = parts[9][2:3]
	stats.HPlayer1ClockSec = parts[9][3:]
	stats.HPlayer2 = parts[10][0:2]
	stats.HPlayer2Clock = parts[10][2:3] + ":" + parts[9][3:]
	stats.HPlayer2ClockMin = parts[10][2:3]
	stats.HPlayer2ClockSec = parts[10][3:]
	stats.VPlayer1 = parts[12][0:2]
	stats.VPlayer1Clock = parts[11][2:3] + ":" + parts[11][3:]
	stats.VPlayer1Min = parts[11][2:3]
	stats.VPlayer1Sec = parts[11][3:]
	return &stats
}

func reverse(str string) string {
	chars := []rune(str)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	str := string(data)

	pos := strings.Index(str, "##")

	if pos == -1 {
		return 0, nil, nil
	}
	advance = pos + 2

	token = []byte(str[0 : pos+2])

	return advance, token, err
}
