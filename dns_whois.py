
### 1. `dns_whois.py` (Python)

```python
# dns_whois.py — Python версия

import sys
import json
import csv
import argparse
import dns.resolver
import whois
from datetime import datetime
from colorama import init, Fore, Style

init(autoreset=True)

class DNSWhoisLookup:
    def __init__(self, domain):
        self.domain = domain
        self.dns_records = {}
        self.whois_data = {}

    def lookup_dns(self):
        """Получает DNS-записи для домена."""
        record_types = ['A', 'AAAA', 'MX', 'NS', 'TXT', 'CNAME', 'SOA']
        resolver = dns.resolver.Resolver()
        resolver.timeout = 5
        resolver.lifetime = 5

        for rtype in record_types:
            try:
                answers = resolver.resolve(self.domain, rtype)
                self.dns_records[rtype] = [str(r) for r in answers]
            except (dns.resolver.NXDOMAIN, dns.resolver.NoAnswer, dns.resolver.Timeout):
                self.dns_records[rtype] = []

    def lookup_whois(self):
        """Получает WHOIS-информацию для домена."""
        try:
            w = whois.whois(self.domain)
            self.whois_data = {
                'registrar': w.registrar,
                'creation_date': str(w.creation_date) if w.creation_date else None,
                'expiration_date': str(w.expiration_date) if w.expiration_date else None,
                'name_servers': w.name_servers,
                'status': w.status,
                'emails': w.emails if hasattr(w, 'emails') else None,
            }
        except Exception as e:
            self.whois_data = {'error': str(e)}

    def print_results(self):
        """Выводит результаты в терминал."""
        print(f"{Fore.CYAN}🌐 DNS Lookup + WHOIS (Python)")
        print(f"Домен: {self.domain}")
        print()

        # DNS
        print(f"{Fore.GREEN}DNS-записи:")
        print("─" * 50)
        for rtype, records in self.dns_records.items():
            if records:
                for r in records:
                    print(f"{Fore.YELLOW}{rtype:<6}{Style.RESET_ALL} {r}")
            else:
                print(f"{Fore.YELLOW}{rtype:<6}{Style.RESET_ALL} —")
        print("─" * 50)

        # WHOIS
        print(f"{Fore.GREEN}WHOIS-информация:")
        print("─" * 50)
        if 'error' in self.whois_data:
            print(f"{Fore.RED}Ошибка WHOIS: {self.whois_data['error']}")
        else:
            w = self.whois_data
            print(f"{Fore.YELLOW}Регистратор:{Style.RESET_ALL}   {w.get('registrar', 'N/A')}")
            print(f"{Fore.YELLOW}Дата регистрации:{Style.RESET_ALL} {w.get('creation_date', 'N/A')}")
            print(f"{Fore.YELLOW}Дата истечения:{Style.RESET_ALL}  {w.get('expiration_date', 'N/A')}")
            ns = w.get('name_servers')
            if ns:
                print(f"{Fore.YELLOW}DNS-серверы:{Style.RESET_ALL}    {', '.join(ns) if isinstance(ns, list) else ns}")
            else:
                print(f"{Fore.YELLOW}DNS-серверы:{Style.RESET_ALL}    N/A")
            if w.get('status'):
                print(f"{Fore.YELLOW}Статус:{Style.RESET_ALL}        {', '.join(w['status']) if isinstance(w['status'], list) else w['status']}")
        print("─" * 50)

    def save_json(self, filename=None):
        if filename is None:
            filename = f"{self.domain.replace('.', '_')}.json"
        data = {
            'domain': self.domain,
            'timestamp': datetime.now().isoformat(),
            'dns': self.dns_records,
            'whois': self.whois_data
        }
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print(f"{Fore.GREEN}💾 Сохранено JSON: {filename}")

    def save_csv(self, filename=None):
        if filename is None:
            filename = f"{self.domain.replace('.', '_')}.csv"
        with open(filename, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerow(['RecordType', 'Value'])
            for rtype, records in self.dns_records.items():
                if records:
                    for r in records:
                        writer.writerow([rtype, r])
                else:
                    writer.writerow([rtype, '—'])
            writer.writerow([])
            writer.writerow(['WHOIS Field', 'Value'])
            for k, v in self.whois_data.items():
                if isinstance(v, list):
                    v = ', '.join(v)
                writer.writerow([k, v])
        print(f"{Fore.GREEN}💾 Сохранено CSV: {filename}")

def main():
    parser = argparse.ArgumentParser(description='DNS Lookup + WHOIS')
    parser.add_argument('domain', help='Домен для проверки')
    parser.add_argument('--json', help='Сохранить в JSON файл')
    parser.add_argument('--csv', help='Сохранить в CSV файл')
    args = parser.parse_args()

    lookup = DNSWhoisLookup(args.domain)
    lookup.lookup_dns()
    lookup.lookup_whois()
    lookup.print_results()

    if args.json:
        lookup.save_json(args.json)
    else:
        lookup.save_json()
    if args.csv:
        lookup.save_csv(args.csv)
    else:
        lookup.save_csv()

if __name__ == "__main__":
    main()
