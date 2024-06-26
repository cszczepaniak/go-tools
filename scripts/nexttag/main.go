package main

import (
	"bufio"
	"bytes"
	"cmp"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	latest := flag.Bool("latest", false, "If true, print the most recent tag that exists.")
	next := flag.Bool("next", false, "If true, print what would be the next patch version.")
	flag.Parse()

	if latest == next {
		panic("must set either latest or next, but not both")
	}

	out, err := exec.Command("git", "tag").CombinedOutput()
	if err != nil {
		panic(err)
	}

	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		panic("expect to find at least one tag")
	}

	sc := bufio.NewScanner(bytes.NewReader(trimmed))

	highest := semver{}
	for sc.Scan() {
		s, err := parseSemver(sc.Text())
		if err != nil {
			fmt.Println("malformed semver: ", err)
			continue
		}

		if s.cmp(highest) > 0 {
			highest = s
		}
	}

	if *next {
		highest.patch++
	}

	fmt.Print(highest)
}

type semver struct {
	major int
	minor int
	patch int
}

func (s semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

func (s semver) cmp(other semver) int {
	if s.major != other.major {
		return cmp.Compare(s.major, other.major)
	}

	if s.minor != other.minor {
		return cmp.Compare(s.minor, other.minor)
	}

	return cmp.Compare(s.patch, other.patch)
}

func parseSemver(s string) (semver, error) {
	if !strings.HasPrefix(s, "v") {
		return semver{}, errors.New("semver should start with 'v'")
	}

	parts := strings.Split(strings.TrimPrefix(s, "v"), ".")
	if len(parts) != 3 {
		return semver{}, errors.New("semver should have three parts")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return semver{}, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return semver{}, err
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return semver{}, err
	}

	return semver{
		major: major,
		minor: minor,
		patch: patch,
	}, nil
}

func compareSemver(a, b string) {
}
