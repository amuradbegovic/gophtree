package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type MenuItem struct {
	Type        byte
	DisplayName string
	Selector    string
	Host        string
	Port        int
}

func (item MenuItem) URL() string {
	url := "gopher://" + item.Host
	if item.Port != 70 {
		url += ":" + strconv.Itoa(item.Port)
	}
	if item.Type != 0 {
		url += "/" + string(item.Type)
	}

	if item.Selector != "/" { // ugly hack
		url += item.Selector
	}

	return url
}

func (item MenuItem) HTML() string {
	url := item.URL()
	html := fmt.Sprintf("<a href=\"%s\">", url)
	if item.DisplayName == "" {
		item.DisplayName = url
	}
	html += item.DisplayName + "</a>"
	return html
}

func (item MenuItem) String() string {
	return fmt.Sprintf("%c%s\t%s\t%s\t%d\n", item.Type, item.DisplayName, item.Selector, item.Host, item.Port)
}

// If successful, returns a MenuItem struct with an empty DisplayName field.
// Other fields are filled with information from the given URL string.
func URLToMenuItem(url string) (MenuItem, error) {
	var addr, host, selector string
	var itemType byte
	var port int

	errInvalid := errors.New("invalid Gopher URL")

	// trim trailing spaces
	url = strings.TrimSpace(url)

	// Strip the "gopher://" URL prefix if found. i
	// If not found, assume it is a gopher server because of the nature of this program.
	url = strings.TrimPrefix(url, "gopher://")

	// Split the URL string with forward slash as the separator
	fields := strings.Split(url, "/")
	if len(fields) < 1 {
		return MenuItem{}, errInvalid
	}

	addr = fields[0]

	// split addr into host and port
	// keep in mind that IPv6 addresses can contain ':' characters before the port
	// only split at the last occurence
	splitIndex := strings.LastIndexByte(addr, ':')
	if splitIndex != -1 {
		var err error
		port, err = strconv.Atoi(addr[splitIndex+1:])
		if err != nil {
			return MenuItem{}, errInvalid
		}

		host = addr[:splitIndex]

		host = strings.TrimPrefix(host, "[")
		host = strings.TrimSuffix(host, "]")
	} else {
		host = addr
		port = 70
	}

	// if the URL contains only the server address, it is probably a link to the root menu
	//itemType = '1'

	if len(fields) == 2 {
		if len(fields[1]) != 1 {
			return MenuItem{}, errInvalid
		} else {
			itemType = fields[1][0]
		}
	}

	if len(fields) >= 3 {
		// merge all remaining fields into selector
		selector = strings.Join(fields[2:], "/")
		selector = "/" + selector
	} else {
		selector = "/"
	}

	return MenuItem{itemType, "", selector, host, port}, nil
}

func getGopher(item MenuItem) (string, error) {
	// connect to the server
	address := item.Host + ":" + strconv.Itoa(item.Port)
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return "", err
	}
	defer connection.Close()

	// send selector extracted from the URL to the server
	_, err = connection.Write([]byte(item.Selector + "\n"))

	if err != nil {
		return "", err
	}

	// finally, read server's response
	response := ""
	reader := bufio.NewReader(connection)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		response += scanner.Text() + "\n"
	}

	return response, nil // return the response
}

func lineToMenuItem(line string) MenuItem {

	fields := strings.Split(line, "\t")

	if len(fields) >= 4 {
		itemType := fields[0][0]
		displayName := fields[0][1:]
		selector := fields[1]
		host := fields[2]
		port, _ := strconv.Atoi(fields[3])

		return MenuItem{itemType, displayName, selector, host, port}
	} else {
		return MenuItem{'i', line, "", "", 0}
	}
}

func responseToMenu(response string) []MenuItem {
	menu := make([]MenuItem, 0)

	scanner := bufio.NewScanner(strings.NewReader(response))
	for scanner.Scan() {
		menu = append(menu, lineToMenuItem(scanner.Text()))
	}
	//fmt.Println(menu)

	return menu
}
