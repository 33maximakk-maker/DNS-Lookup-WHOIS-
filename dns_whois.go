// dns_whois.go — Go версия

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type DNSWhois struct {
	Domain     string
	DNSRecords map[string][]string
	WhoisData  map[string]interface{}
}

func NewDNSWhois(domain string) *DNSWhois {
	return &DNSWhois{
		Domain:     domain,
		DNSRecords: make(map[string][]string),
		WhoisData:  make(map[string]interface{}),
	}
}

func (d *DNSWhois) lookupDNS() {
	recordTypes := []string{"A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA"}
	for _, rtype := range recordTypes {
		switch rtype {
		case "A":
			ips, err := net.LookupIP(d.Domain)
			if err == nil {
				for _, ip := range ips {
					if ipv4 := ip.To4(); ipv4 != nil {
						d.DNSRecords["A"] = append(d.DNSRecords["A"], ipv4.String())
					}
				}
			}
		case "AAAA":
			ips, err := net.LookupIP(d.Domain)
			if err == nil {
				for _, ip := range ips {
					if ip.To4() == nil && ip.To16() != nil {
						d.DNSRecords["AAAA"] = append(d.DNSRecords["AAAA"], ip.String())
					}
				}
			}
		case "MX":
			mx, err := net.LookupMX(d.Domain)
			if err == nil {
				for _, m := range mx {
					d.DNSRecords["MX"] = append(d.DNSRecords["MX"], fmt.Sprintf("%d %s", m.Pref, m.Host))
				}
			}
		case "NS":
			ns, err := net.LookupNS(d.Domain)
			if err == nil {
				for _, n := range ns {
					d.DNSRecords["NS"] = append(d.DNSRecords["NS"], n.Host)
				}
			}
		case "TXT":
			txt, err := net.LookupTXT(d.Domain)
			if err == nil {
				d.DNSRecords["TXT"] = append(d.DNSRecords["TXT"], txt...)
			}
		case "CNAME":
			cname, err := net.LookupCNAME(d.Domain)
			if err == nil {
				d.DNSRecords["CNAME"] = append(d.DNSRecords["CNAME"], cname)
			}
		case "SOA":
			// SOA не поддерживается напрямую, пропускаем
		}
	}
}

func (d *DNSWhois) lookupWhois() {
	raw, err := whois.Whois(d.Domain)
	if err != nil {
		d.WhoisData["error"] = err.Error()
		return
	}
	parsed, err := whoisparser.Parse(raw)
	if err != nil {
		d.WhoisData["error"] = err.Error()
		return
	}
	d.WhoisData["domain"] = parsed.Domain
	d.WhoisData["registrar"] = parsed.Registrar
	d.WhoisData["creation_date"] = parsed.CreatedDate
	d.WhoisData["expiration_date"] = parsed.ExpirationDate
	d.WhoisData["name_servers"] = parsed.NameServers
	d.WhoisData["status"] = parsed.Status
}

func (d *DNSWhois) printResults() {
	fmt.Println("\x1b[36m🌐 DNS Lookup + WHOIS (Go)\x1b[0m")
	fmt.Printf("Домен: %s\n\n", d.Domain)

	fmt.Println("\x1b[32mDNS-записи:\x1b[0m")
	fmt.Println("─" * 50)
	for rtype, records := range d.DNSRecords {
		if len(records) > 0 {
			for _, r := range records {
				fmt.Printf("\x1b[33m%-6s\x1b[0m %s\n", rtype, r)
			}
		} else {
			fmt.Printf("\x1b[33m%-6s\x1b[0m —\n", rtype)
		}
	}
	fmt.Println("─" * 50)

	fmt.Println("\x1b[32mWHOIS-информация:\x1b[0m")
	fmt.Println("─" * 50)
	if err, ok := d.WhoisData["error"]; ok {
		fmt.Printf("\x1b[31mОшибка WHOIS: %v\x1b[0m\n", err)
	} else {
		fmt.Printf("\x1b[33mРегистратор:\x1b[0m   %v\n", d.WhoisData["registrar"])
		fmt.Printf("\x1b[33mДата регистрации:\x1b[0m %v\n", d.WhoisData["creation_date"])
		fmt.Printf("\x1b[33mДата истечения:\x1b[0m  %v\n", d.WhoisData["expiration_date"])
		if ns, ok := d.WhoisData["name_servers"]; ok {
			fmt.Printf("\x1b[33mDNS-серверы:\x1b[0m    %v\n", strings.Join(ns.([]string), ", "))
		}
		if status, ok := d.WhoisData["status"]; ok {
			fmt.Printf("\x1b[33mСтатус:\x1b[0m        %v\n", status)
		}
	}
	fmt.Println("─" * 50)
}

func (d *DNSWhois) saveJSON(filename string) {
	if filename == "" {
		filename = strings.ReplaceAll(d.Domain, ".", "_") + ".json"
	}
	data := map[string]interface{}{
		"domain":    d.Domain,
		"timestamp": time.Now().Format(time.RFC3339),
		"dns":       d.DNSRecords,
		"whois":     d.WhoisData,
	}
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(filename, jsonData, 0644)
	fmt.Printf("\x1b[32m💾 Сохранено JSON: %s\x1b[0m\n", filename)
}

func (d *DNSWhois) saveCSV(filename string) {
	if filename == "" {
		filename = strings.ReplaceAll(d.Domain, ".", "_") + ".csv"
	}
	file, _ := os.Create(filename)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"RecordType", "Value"})
	for rtype, records := range d.DNSRecords {
		if len(records) > 0 {
			for _, r := range records {
				writer.Write([]string{rtype, r})
			}
		} else {
			writer.Write([]string{rtype, "—"})
		}
	}
	writer.Write([]string{})
	writer.Write([]string{"WHOIS Field", "Value"})
	for k, v := range d.WhoisData {
		if k == "error" {
			writer.Write([]string{k, fmt.Sprintf("%v", v)})
			continue
		}
		str := fmt.Sprintf("%v", v)
		if k == "name_servers" {
			if ns, ok := v.([]string); ok {
				str = strings.Join(ns, ", ")
			}
		}
		writer.Write([]string{k, str})
	}
	fmt.Printf("\x1b[32m💾 Сохранено CSV: %s\x1b[0m\n", filename)
}

func main() {
	domain := flag.String("domain", "", "Домен для проверки")
	jsonFile := flag.String("json", "", "Сохранить в JSON файл")
	csvFile := flag.String("csv", "", "Сохранить в CSV файл")
	flag.Parse()

	if *domain == "" && flag.NArg() > 0 {
		*domain = flag.Arg(0)
	}
	if *domain == "" {
		fmt.Println("Usage: go run dns_whois.go <domain> [--json file] [--csv file]")
		os.Exit(1)
	}

	d := NewDNSWhois(*domain)
	d.lookupDNS()
	d.lookupWhois()
	d.printResults()

	if *jsonFile != "" {
		d.saveJSON(*jsonFile)
	} else {
		d.saveJSON("")
	}
	if *csvFile != "" {
		d.saveCSV(*csvFile)
	} else {
		d.saveCSV("")
	}
}
