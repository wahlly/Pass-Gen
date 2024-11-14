package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var usage = `
Usage
----
--get platform=[string] - Gets saved password for a platform
--set platform=[string] len=[int] level=(basic|medium|hard|xhard) - Creates and saves a password
`

var ErrUsage = errors.New(usage)

var pattern = regexp.MustCompile(`\S+=\S+`)

type level int

const (
	_ level = iota
	level_basic
	level_medium
	level_hard
	level_xhard
)

var level_key = map[string]level{
	"basic": level_basic,
	"medium": level_medium,
	"hard": level_hard,
	"xhard": level_xhard,
}

type commands struct{
	get, set bool
}

func createCommands() (c commands) {
	flag.BoolVar(&c.get, "get", false, "get password for platform")
	flag.BoolVar(&c.set, "set", false, "set password for platform")
	flag.Parse()
	return
}

func (c commands) exec(store *store) (string, error) {
	switch {
	case c.get:
		return c.getPassword(store)
	case c.set:
		return c.setPassword(store)
	default:
		return "", ErrUsage
	}
}

func (c commands) getPassword(store *store) (string, error) {
	params, err := c.parse()
	if err != nil {
		return "", err
	}

	return store.find(params["platform"])
}

func (c commands) setPassword(store *store) (string, error) {
	params, err := c.parse()
	if err != nil {
		return "", err
	}

	var password string

	n, err := strconv.Atoi(params["len"])
	if err != nil {
		return "", err
	}

	if n < 8 {
		return "", fmt.Errorf("password len cannot be less than 8")
	}

	switch level_key[params["level"]] {
	case level_basic:
		password = basicPassword(n)
	case level_medium:
		password = mediumPassword(n)
	case level_hard:
		password = hardPassword(n)
	case level_xhard:
		password = xhardPassword(n)
	default:
		return "", ErrUsage
	}

	password = shuffleStr(password)

	if err := store.add(params["platform"], password); err != nil {
		return "", err
	}

	return password, nil
}

func (c commands) parse() (map[string]string, error) {
	args := flag.Args()

	if len(args) == 0 {
		return nil, ErrUsage
	}

	params := make(map[string]string)

	for i := range args {
		if !pattern.MatchString(args[i]) {
			return nil, ErrUsage
		}

		parts := strings.Split(args[i], "=")
		params[parts[0]] = parts[1]
	}

	return params, nil
}

func main() {
	store, err := newStore()
	if err != nil {
		log.Fatalf("could not initialize store: %v", err)
	}

	c := createCommands()

	password, err := c.exec(store)
	if err != nil {
		log.Fatalf("could not execute flag commands: %v", err)
	}

	fmt.Printf("password: %s\n", password)
}