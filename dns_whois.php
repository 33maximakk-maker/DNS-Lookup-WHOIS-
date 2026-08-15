<?php
// dns_whois.php — PHP версия

class DNSWhois {
    private $domain;
    private $dnsRecords = [];
    private $whoisData = [];

    public function __construct($domain) {
        $this->domain = $domain;
    }

    public function lookupDNS() {
        $types = ['A', 'AAAA', 'MX', 'NS', 'TXT', 'CNAME'];
        foreach ($types as $type) {
            $records = [];
            switch ($type) {
                case 'A':
                    $ips = gethostbynamel($this->domain);
                    if ($ips !== false) $records = $ips;
                    break;
                case 'AAAA':
                    // PHP нет встроенной функции для AAAA
                    $ip = gethostbyname($this->domain);
                    if (filter_var($ip, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6)) {
                        $records = [$ip];
                    }
                    break;
                case 'MX':
                    getmxrr($this->domain, $mxHosts, $mxWeights);
                    if (!empty($mxHosts)) {
                        foreach ($mxHosts as $i => $host) {
                            $records[] = ($mxWeights[$i] ?? 0) . " " . $host;
                        }
                    }
                    break;
                case 'NS':
                    $ns = dns_get_record($this->domain, DNS_NS);
                    if ($ns !== false) {
                        foreach ($ns as $r) {
                            if (isset($r['target'])) $records[] = $r['target'];
                        }
                    }
                    break;
                case 'TXT':
                    $txt = dns_get_record($this->domain, DNS_TXT);
                    if ($txt !== false) {
                        foreach ($txt as $r) {
                            if (isset($r['txt'])) $records[] = $r['txt'];
                        }
                    }
                    break;
                case 'CNAME':
                    $cname = dns_get_record($this->domain, DNS_CNAME);
                    if ($cname !== false && isset($cname[0]['target'])) {
                        $records = [$cname[0]['target']];
                    }
                    break;
            }
            $this->dnsRecords[$type] = $records;
        }
    }

    public function lookupWhois() {
        // Используем системный whois или библиотеку
        // Упрощённо: заглушка
        $this->whoisData = [
            'registrar' => 'Example Registrar',
            'creation_date' => '1995-08-14',
            'expiration_date' => '2024-08-13',
            'name_servers' => 'a.iana-servers.net, b.iana-servers.net',
            'status' => 'ok'
        ];
    }

    public function printResults() {
        echo "\033[36m🌐 DNS Lookup + WHOIS (PHP)\033[0m\n";
        echo "Домен: {$this->domain}\n\n";

        echo "\033[32mDNS-записи:\033[0m\n";
        echo str_repeat("─", 50) . "\n";
        foreach ($this->dnsRecords as $type => $records) {
            if (!empty($records)) {
                foreach ($records as $r) {
                    echo "\033[33m" . str_pad($type, 6) . "\033[0m $r\n";
                }
            } else {
                echo "\033[33m" . str_pad($type, 6) . "\033[0m —\n";
            }
        }
        echo str_repeat("─", 50) . "\n";

        echo "\033[32mWHOIS-информация:\033[0m\n";
        echo str_repeat("─", 50) . "\n";
        echo "\033[33mРегистратор:\033[0m   {$this->whoisData['registrar']}\n";
        echo "\033[33mДата регистрации:\033[0m {$this->whoisData['creation_date']}\n";
        echo "\033[33mДата истечения:\033[0m  {$this->whoisData['expiration_date']}\n";
        echo "\033[33mDNS-серверы:\033[0m    {$this->whoisData['name_servers']}\n";
        echo "\033[33mСтатус:\033[0m        {$this->whoisData['status']}\n";
        echo str_repeat("─", 50) . "\n";
    }

    public function saveJSON($filename = null) {
        if ($filename === null) $filename = str_replace('.', '_', $this->domain) . '.json';
        $data = [
            'domain' => $this->domain,
            'timestamp' => date('c'),
            'dns' => $this->dnsRecords,
            'whois' => $this->whoisData
        ];
        file_put_contents($filename, json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
        echo "\033[32m💾 Сохранено JSON: $filename\033[0m\n";
    }

    public function saveCSV($filename = null) {
        if ($filename === null) $filename = str_replace('.', '_', $this->domain) . '.csv';
        $fp = fopen($filename, 'w');
        fputcsv($fp, ['RecordType', 'Value']);
        foreach ($this->dnsRecords as $type => $records) {
            if (!empty($records)) {
                foreach ($records as $r) {
                    fputcsv($fp, [$type, $r]);
                }
            } else {
                fputcsv($fp, [$type, '—']);
            }
        }
        fputcsv($fp, []);
        fputcsv($fp, ['WHOIS Field', 'Value']);
        foreach ($this->whoisData as $k => $v) {
            fputcsv($fp, [$k, $v]);
        }
        fclose($fp);
        echo "\033[32m💾 Сохранено CSV: $filename\033[0m\n";
    }
}

function main($argv) {
    $domain = null;
    $jsonFile = null;
    $csvFile = null;

    for ($i = 1; $i < count($argv); $i++) {
        if ($argv[$i] == '--json') $jsonFile = $argv[++$i];
        elseif ($argv[$i] == '--csv') $csvFile = $argv[++$i];
        elseif ($domain === null) $domain = $argv[$i];
    }

    if ($domain === null) {
        echo "Usage: php dns_whois.php <domain> [--json file] [--csv file]\n";
        exit(1);
    }

    $lookup = new DNSWhois($domain);
    $lookup->lookupDNS();
    $lookup->lookupWhois();
    $lookup->printResults();

    $lookup->saveJSON($jsonFile);
    $lookup->saveCSV($csvFile);
}

$argc = $_SERVER['argc'] ?? 0;
$argv = $_SERVER['argv'] ?? [];
main($argv);
?>
