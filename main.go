package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

func main() {
	server, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		if err := conn.Close(); err != nil {
			log.Println(err)
		}

		processData(conn)

		break
	}
}

func processData(conn io.Reader) {
	buff := bufio.NewScanner(conn)

	buff.Split(split)

	for buff.Scan() {
		line := buff.Text()

		if line == "" {
			break
		}

		parseLine(line)
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

	re := regexp.MustCompile("/([0-9]{2})([0-9]{2})/")
	match := re.FindStringSubmatch(parts[3])

	if flag == "T" {
		stats.clock = match[1] + ":" + match[2]
		stats.clockMode = ":"
		stats.clockMin = match[1]
		stats.clockSec = match[2]

	} else {
		stats.clock = match[1] + "." + match[2][0:1]
		stats.clockMode = "."
		stats.clockMin = match[1]
		stats.clockSec = match[2][0:1]
	}
	stats.clockStatus = "Running" // TODO: make it dynamic from clock data
	stats.period = parts[6]
	stats.homeScore = parts[4]
	stats.guestScore = parts[5]
	stats.gamePeriod = parts[6]
	stats.homeShots = parts[7]
	stats.guestShots = parts[8]
	stats.bot = parts[1]
	stats.hPlayer1 = parts[9][0:2]
	stats.hPlayer1Clock = parts[9][2:3] + ":" + parts[9][3:]
	stats.hPlayer1ClockMin = parts[9][2:3]
	stats.hPlayer1ClockSec = parts[9][3:]
	stats.hPlayer2 = parts[10][0:2]
	stats.hPlayer2Clock = parts[10][2:3] + ":" + parts[9][3:]
	stats.hPlayer2ClockMin = parts[10][2:3]
	stats.hPlayer2ClockSec = parts[10][3:]
	stats.vPlayer1 = parts[12][0:2]
	stats.vPlayer1Clock = parts[11][2:3] + ":" + parts[11][3:]
	stats.vPlayer1Min = parts[11][2:3]
	stats.vPlayer1Sec = parts[11][3:]
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
