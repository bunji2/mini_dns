// mini_dns.go
// Author: Bunji Square
// Version: 1.1.a
// Usage: mini_dns -h

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"

	//"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/miekg/dns"
)

const (
	version       = "1.1.a"
	flag1_name    = "zone_file"
	flag1_default = "./zone.txt"
	flag1_desc    = "zone file"
	flag2_name    = "port"
	flag2_default = 53
	flag2_desc    = "port number"
	flag3_name    = "server"
	flag3_default = ""
	flag3_desc    = "addr:port   of dns server"
)

func main() {
	log.Println("mini_dns", version)
	// 実行ファイルのディレクトリに移動する
	p, err := os.Executable()
	if err != nil {
		log.Panicln(err)
	}
	os.Chdir(filepath.Dir(p))

	zone_file := flag.String(flag1_name, flag1_default, flag1_desc)
	port := flag.Int(flag2_name, flag2_default, flag2_desc)
	server := flag.String(flag3_name, flag3_default, flag3_desc)
	flag.Parse()

	log.Println("conf: zone_file =", *zone_file)
	if *server != "" {
		log.Println("conf: dns server =", *server)
	} else {
		log.Println("conf: no dns server")
	}
	fp, err := os.Open(*zone_file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	reader := bufio.NewReaderSize(fp, 2048)

	rrs := []dns.RR{}
	/*
		to := dns.ParseZone(reader, "", "myzone")
		for x := range to {
			fmt.Println("#", x.RR)
			rrs = append(rrs, dns.Copy(x.RR))
		}
	*/
	to := dns.NewZoneParser(reader, "", "myzone")
	for x, ok := to.Next(); ok; x, ok = to.Next() {
		fmt.Println("#", x)
		rrs = append(rrs, dns.Copy(x))
	}

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true

		for _, q := range r.Question {
			answers := []dns.RR{}
			//fmt.Println("q =", q)
			for _, rr := range rrs {
				//fmt.Println("rr =", rr)
				rh := rr.Header()
				if q.Name == rh.Name && q.Qtype == rh.Rrtype && q.Qclass == rh.Class {
					answers = append(answers, rr)
				}
			}
			if len(answers) == 0 && *server != "" {
				answers = append(answers, resolver(*server, q.Name, q.Qtype)...)
			}
			m.Answer = append(m.Answer, answers...)
		}
		w.WriteMsg(m)
		log.Println(r)
		log.Println(m)
	})

	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(*port), Net: "udp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping\n", s)
}

func resolver(server, fqdn string, r_type uint16) []dns.RR {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	//m1.RecursionDesired = true
	m1.SetQuestion(fqdn, r_type)

	in, err := dns.Exchange(m1, server)
	if err == nil {
		return in.Answer
	}
	return []dns.RR{}
}
