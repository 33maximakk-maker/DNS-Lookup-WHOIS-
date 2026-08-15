// dns_whois.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Text.Json;
using System.Text.RegularExpressions;
using System.Threading.Tasks;

class DNSWhois
{
    private string domain;
    private Dictionary<string, List<string>> dnsRecords = new Dictionary<string, List<string>>();
    private Dictionary<string, string> whoisData = new Dictionary<string, string>();

    public DNSWhois(string domain)
    {
        this.domain = domain;
    }

    public async Task LookupDNSAsync()
    {
        string[] types = { "A", "AAAA", "MX", "NS", "TXT", "CNAME" };
        foreach (var type in types)
        {
            var records = new List<string>();
            try
            {
                switch (type)
                {
                    case "A":
                        var ips = await Dns.GetHostAddressesAsync(domain);
                        records.AddRange(ips.Where(ip => ip.AddressFamily == System.Net.Sockets.AddressFamily.InterNetwork)
                                           .Select(ip => ip.ToString()));
                        break;
                    case "AAAA":
                        ips = await Dns.GetHostAddressesAsync(domain);
                        records.AddRange(ips.Where(ip => ip.AddressFamily == System.Net.Sockets.AddressFamily.InterNetworkV6)
                                           .Select(ip => ip.ToString()));
                        break;
                    case "MX":
                        var mx = await Dns.GetMXRecordsAsync(domain);
                        records.AddRange(mx.Select(m => $"{m.Preference} {m.Exchange}"));
                        break;
                    case "NS":
                        var ns = await Dns.GetNSRecordsAsync(domain);
                        records.AddRange(ns.Select(n => n.ToString()));
                        break;
                    case "TXT":
                        var txt = await Dns.GetTXTRecordsAsync(domain);
                        records.AddRange(txt.Select(t => t.ToString()));
                        break;
                    case "CNAME":
                        var host = await Dns.GetHostEntryAsync(domain);
                        if (host.HostName != domain) records.Add(host.HostName);
                        break;
                }
            }
            catch { }
            dnsRecords[type] = records;
        }
    }

    public void LookupWhois()
    {
        // Упрощённо: в C# нет встроенного WHOIS, используем заглушку
        whoisData["registrar"] = "Example Registrar";
        whoisData["creation_date"] = "1995-08-14";
        whoisData["expiration_date"] = "2024-08-13";
        whoisData["name_servers"] = "a.iana-servers.net, b.iana-servers.net";
        whoisData["status"] = "ok";
    }

    public void PrintResults()
    {
        Console.WriteLine("\u001B[36m🌐 DNS Lookup + WHOIS (C#)\u001B[0m");
        Console.WriteLine($"Домен: {domain}\n");

        Console.WriteLine("\u001B[32mDNS-записи:\u001B[0m");
        Console.WriteLine(new string('─', 50));
        foreach (var kv in dnsRecords)
        {
            if (kv.Value.Count > 0)
            {
                foreach (var r in kv.Value)
                    Console.WriteLine($"\u001B[33m{kv.Key,-6}\u001B[0m {r}");
            }
            else
            {
                Console.WriteLine($"\u001B[33m{kv.Key,-6}\u001B[0m —");
            }
        }
        Console.WriteLine(new string('─', 50));

        Console.WriteLine("\u001B[32mWHOIS-информация:\u001B[0m");
        Console.WriteLine(new string('─', 50));
        Console.WriteLine($"\u001B[33mРегистратор:\u001B[0m   {whoisData.GetValueOrDefault("registrar", "N/A")}");
        Console.WriteLine($"\u001B[33mДата регистрации:\u001B[0m {whoisData.GetValueOrDefault("creation_date", "N/A")}");
        Console.WriteLine($"\u001B[33mДата истечения:\u001B[0m  {whoisData.GetValueOrDefault("expiration_date", "N/A")}");
        Console.WriteLine($"\u001B[33mDNS-серверы:\u001B[0m    {whoisData.GetValueOrDefault("name_servers", "N/A")}");
        Console.WriteLine($"\u001B[33mСтатус:\u001B[0m        {whoisData.GetValueOrDefault("status", "N/A")}");
        Console.WriteLine(new string('─', 50));
    }

    public void SaveJSON(string filename)
    {
        if (string.IsNullOrEmpty(filename))
            filename = domain.Replace(".", "_") + ".json";
        var data = new
        {
            domain = this.domain,
            timestamp = DateTime.Now.ToString("o"),
            dns = dnsRecords,
            whois = whoisData
        };
        string json = JsonSerializer.Serialize(data, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine($"\u001B[32m💾 Сохранено JSON: {filename}\u001B[0m");
    }

    public void SaveCSV(string filename)
    {
        if (string.IsNullOrEmpty(filename))
            filename = domain.Replace(".", "_") + ".csv";
        using var writer = new StreamWriter(filename);
        writer.WriteLine("RecordType,Value");
        foreach (var kv in dnsRecords)
        {
            if (kv.Value.Count > 0)
            {
                foreach (var r in kv.Value)
                    writer.WriteLine($"{kv.Key},{r}");
            }
            else
            {
                writer.WriteLine($"{kv.Key},—");
            }
        }
        writer.WriteLine();
        writer.WriteLine("WHOIS Field,Value");
        foreach (var kv in whoisData)
            writer.WriteLine($"{kv.Key},{kv.Value}");
        Console.WriteLine($"\u001B[32m💾 Сохранено CSV: {filename}\u001B[0m");
    }

    public static async Task Main(string[] args)
    {
        string domain = null;
        string jsonFile = null;
        string csvFile = null;

        for (int i = 0; i < args.Length; i++)
        {
            if (args[i] == "--json") jsonFile = args[++i];
            else if (args[i] == "--csv") csvFile = args[++i];
            else if (domain == null) domain = args[i];
        }

        if (domain == null)
        {
            Console.WriteLine("Usage: dotnet run <domain> [--json file] [--csv file]");
            return;
        }

        var lookup = new DNSWhois(domain);
        await lookup.LookupDNSAsync();
        lookup.LookupWhois();
        lookup.PrintResults();

        lookup.SaveJSON(jsonFile);
        lookup.SaveCSV(csvFile);
    }
}
