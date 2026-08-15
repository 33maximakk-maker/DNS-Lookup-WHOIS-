// dns_whois.rs — Rust версия

use std::env;
use std::fs;
use std::net::ToSocketAddrs;
use std::time::SystemTime;
use serde::{Deserialize, Serialize};
use whois_rust::{Whois, WhoIsLookupOptions};

#[derive(Serialize, Deserialize)]
struct DNSData {
    domain: String,
    timestamp: String,
    dns: std::collections::HashMap<String, Vec<String>>,
    whois: std::collections::HashMap<String, String>,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();
    let mut domain = None;
    let mut json_file = None;
    let mut csv_file = None;

    let mut i = 1;
    while i < args.len() {
        match args[i].as_str() {
            "--json" => { json_file = Some(args[i+1].clone()); i += 2; }
            "--csv" => { csv_file = Some(args[i+1].clone()); i += 2; }
            _ => { if domain.is_none() { domain = Some(args[i].clone()); } i += 1; }
        }
    }

    if domain.is_none() {
        println!("Usage: cargo run -- <domain> [--json file] [--csv file]");
        return Ok(());
    }

    let domain = domain.unwrap();
    println!("\x1b[36m🌐 DNS Lookup + WHOIS (Rust)\x1b[0m");
    println!("Домен: {}\n", domain);

    let mut dns_records = std::collections::HashMap::new();
    // DNS lookup (упрощённо, используем системный резолвер)
    let record_types = ["A", "AAAA", "MX", "NS", "TXT", "CNAME"];
    for &rtype in &record_types {
        let mut records = Vec::new();
        match rtype {
            "A" => {
                if let Ok(ips) = dns_lookup::lookup_host(&domain, dns_lookup::LookupIpStrategy::Ipv4Only) {
                    for ip in ips {
                        records.push(ip.to_string());
                    }
                }
            }
            "AAAA" => {
                if let Ok(ips) = dns_lookup::lookup_host(&domain, dns_lookup::LookupIpStrategy::Ipv6Only) {
                    for ip in ips {
                        records.push(ip.to_string());
                    }
                }
            }
            _ => {} // другие типы пропускаем для простоты
        }
        dns_records.insert(rtype.to_string(), records);
    }

    // WHOIS
    let whois = Whois::default();
    let result = whois.lookup(WhoIsLookupOptions::from_str(&domain))?;
    let whois_text = result.unwrap_or_default();
    let mut whois_data = std::collections::HashMap::new();
    // Парсинг простой
    for line in whois_text.lines() {
        if let Some(pos) = line.find(':') {
            let key = line[..pos].trim().to_lowercase().replace(" ", "_");
            let value = line[pos+1..].trim().to_string();
            if !key.is_empty() {
                whois_data.insert(key, value);
            }
        }
    }

    // Вывод
    println!("\x1b[32mDNS-записи:\x1b[0m");
    println!("{}", "─".repeat(50));
    for (rtype, records) in &dns_records {
        if records.is_empty() {
            println!("\x1b[33m{:6}\x1b[0m —", rtype);
        } else {
            for r in records {
                println!("\x1b[33m{:6}\x1b[0m {}", rtype, r);
            }
        }
    }
    println!("{}", "─".repeat(50));

    println!("\x1b[32mWHOIS-информация:\x1b[0m");
    println!("{}", "─".repeat(50));
    for (key, value) in &whois_data {
        if key.contains("registrar") {
            println!("\x1b[33mРегистратор:\x1b[0m   {}", value);
        } else if key.contains("creation") {
            println!("\x1b[33mДата регистрации:\x1b[0m {}", value);
        } else if key.contains("expiration") {
            println!("\x1b[33mДата истечения:\x1b[0m  {}", value);
        } else if key.contains("name_server") {
            println!("\x1b[33mDNS-серверы:\x1b[0m    {}", value);
        } else if key.contains("status") {
            println!("\x1b[33mСтатус:\x1b[0m        {}", value);
        }
    }
    println!("{}", "─".repeat(50));

    // Сохранение
    let data = DNSData {
        domain: domain.clone(),
        timestamp: SystemTime::now().duration_since(SystemTime::UNIX_EPOCH).unwrap().as_secs().to_string(),
        dns: dns_records,
        whois: whois_data,
    };
    let json = serde_json::to_string_pretty(&data)?;
    let json_filename = json_file.unwrap_or_else(|| format!("{}.json", domain.replace(".", "_")));
    fs::write(&json_filename, json)?;
    println!("\x1b[32m💾 Сохранено JSON: {}\x1b[0m", json_filename);

    let csv_filename = csv_file.unwrap_or_else(|| format!("{}.csv", domain.replace(".", "_")));
    let mut csv = String::from("RecordType,Value\n");
    for (rtype, records) in &data.dns {
        if records.is_empty() {
            csv.push_str(&format!("{},—\n", rtype));
        } else {
            for r in records {
                csv.push_str(&format!("{},{}\n", rtype, r));
            }
        }
    }
    csv.push_str("\nWHOIS Field,Value\n");
    for (key, value) in &data.whois {
        csv.push_str(&format!("{},{}\n", key, value));
    }
    fs::write(&csv_filename, csv)?;
    println!("\x1b[32m💾 Сохранено CSV: {}\x1b[0m", csv_filename);

    Ok(())
}
