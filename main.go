package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// permanode is sha1-6f498bda032b38cc94f6267a40f994c6788e9b41

var (
	flagVerbose = flag.Bool("v", false, "Be verbose")
	flagInput = flag.String("input", "", "text file to read as input")
)

var input []byte

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("nope")
	}
	var err error
	if *flagInput != "" {
		input, err = ioutil.ReadFile(*flagInput)
		if err != nil {
			log.Fatal(err)
		}
	}
	addMember(args[0])
}

func camtoolSearchBlobs() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		l := sc.Text()
		if !strings.Contains(l, `"blob":`) {
			continue
		}
		fields := strings.Fields(l)
		if len(fields) != 2 {
			continue
		}
		fmt.Printf("%s ", strings.Replace(fields[1], `"`, "", -1))
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}

func addMember(containerSet string) {
	// TODO(mpl): figure out how to read directly from the piped output when it is large
	var inputRd io.Reader = os.Stdin
	if *flagInput != "" {
		inputRd = bytes.NewReader(input)
	}
	sc := bufio.NewScanner(inputRd)
	count := 0
	for sc.Scan() {
		l := sc.Text()
		if !strings.Contains(l, `"blob":`) {
			continue
		}
		fields := strings.Fields(l)
		if len(fields) != 2 {
			continue
		}

		ref := strings.Replace(fields[1], `"`, "", -1)
		cmdstr := "attr -add " + containerSet + " camliMember " + ref
		if *flagVerbose {
			println(cmdstr)
		}
		if out, err := exec.Command("camput", strings.Fields(cmdstr)...).Output(); err != nil {
			log.Fatalf("camput error: %v, %v", err, string(out))
		}
		count++
	}
	if *flagVerbose {
		println(count)
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scanning error: %v", err)
	}
}

func addMemberInDocker() {
	//	var vals []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		l := sc.Text()
		if !strings.Contains(l, `"blob":`) {
			continue
		}
		fields := strings.Fields(l)
		if len(fields) != 2 {
			continue
		}

		ref := strings.Replace(fields[1], `"`, "", -1)
		cmdstr := "run --rm -v /home/mpl/.config/camlistore/other:/home/camli/.config/camlistore camlistore/world camput attr -add sha1-e0d659e2da43e09470dd43919c3db16c53eba5a6 camliMember " + ref
		println(cmdstr)
		if err := exec.Command("docker", strings.Fields(cmdstr)...).Run(); err != nil {
			log.Fatal(err)
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}
