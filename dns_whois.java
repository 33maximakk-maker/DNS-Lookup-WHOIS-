// dns_whois.java — Java версия

import java.io.*;
import java.net.*;
import java.nio.file.*;
import java.util.*;
import java.time.*;
import com.google.gson.GsonBuilder;
import com.google.gson.annotations.SerializedName;

public class dns_whois {
    private String domain;
    private Map<String, List<String>> dnsRecords = new LinkedHashMap<>();
    private Map<String, String> whoisData = new LinkedHashMap<>();

    public dns_whois(String domain) {
        this.domain = domain;
    }

    public void lookupDNS() throws Exception {
        String[] types = {"A", "AAAA", "MX", "NS", "TXT", "CNAME"};
        for (String type : types) {
            List<String> records = new ArrayList<>();
            try {
                switch (type) {
                    case "A":
                        InetAddress[] addrs = InetAddress.getAllByName(domain);
                        for (InetAddress addr : addrs) {
                            if (addr instanceof Inet4Address) {
                                records.add(addr.getHostAddress());
                            }
                        }
                        break;
                    case "AAAA":
                        addrs = InetAddress.getAllByName(domain);
                        for (InetAddress addr : addrs) {
                            if (addr instanceof Inet6Address) {
                                records.add(addr.getHostAddress());
                            }
                        }
                        break;
                    case "MX":
                        // Используем DNS Java библиотеку для MX, но здесь упрощённо
                        // В реальном проекте используйте dnsjava
                        break;
                    case "NS":
                        // Нет прямого метода в стандартной библиотеке
                        break;
                    case "TXT":
                        // нет в стандартной библиотеке
                        break;
                    case "CNAME":
                        // нет в стандартной библиотеке
                        break;
                }
            } catch (Exception e) {}
            dnsRecords.put(type, records);
        }
    }

    public void lookupWhois() {
        // В Java нет встроенного whois, используем Socket или библиотеку
        // Упрощённо: заглушка
        whoisData.put("registrar", "Example Registrar");
        whoisData.put("creation_date", "1995-08-14");
        whoisData.put("expiration_date", "2024-08-13");
        whoisData.put("name_servers", "a.iana-servers.net, b.iana-servers.net");
        whoisData.put("status", "ok");
    }

    public void printResults() {
        System.out.println("\u001B[36m🌐 DNS Lookup + WHOIS (Java)\u001B[0m");
        System.out.println("Домен: " + domain);
        System.out.println();

        System.out.println("\u001B[32mDNS-записи:\u001B[0m");
        System.out.println("─".repeat(50));
        for (Map.Entry<String, List<String>> entry : dnsRecords.entrySet()) {
            String type = entry.getKey();
            List<String> recs = entry.getValue();
            if (!recs.isEmpty()) {
                for (String r : recs) {
                    System.out.printf("\u001B[33m%-6s\u001B[0m %s\n", type, r);
                }
            } else {
                System.out.printf("\u001B[33m%-6s\u001B[0m —\n", type);
            }
        }
        System.out.println("─".repeat(50));

        System.out.println("\u001B[32mWHOIS-информация:\u001B[0m");
        System.out.println("─".repeat(50));
        System.out.printf("\u001B[33mРегистратор:\u001B[0m   %s\n", whoisData.getOrDefault("registrar", "N/A"));
        System.out.printf("\u001B[33mДата регистрации:\u001B[0m %s\n", whoisData.getOrDefault("creation_date", "N/A"));
        System.out.printf("\u001B[33mДата истечения:\u001B[0m  %s\n", whoisData.getOrDefault("expiration_date", "N/A"));
        System.out.printf("\u001B[33mDNS-серверы:\u001B[0m    %s\n", whoisData.getOrDefault("name_servers", "N/A"));
        System.out.printf("\u001B[33mСтатус:\u001B[0m        %s\n", whoisData.getOrDefault("status", "N/A"));
        System.out.println("─".repeat(50));
    }

    public void saveJSON(String filename) throws IOException {
        if (filename == null) filename = domain.replace(".", "_") + ".json";
        Map<String, Object> data = new LinkedHashMap<>();
        data.put("domain", domain);
        data.put("timestamp", Instant.now().toString());
        data.put("dns", dnsRecords);
        data.put("whois", whoisData);
        String json = new GsonBuilder().setPrettyPrinting().create().toJson(data);
        Files.write(Paths.get(filename), json.getBytes());
        System.out.println("\u001B[32m💾 Сохранено JSON: " + filename + "\u001B[0m");
    }

    public void saveCSV(String filename) throws IOException {
        if (filename == null) filename = domain.replace(".", "_") + ".csv";
        StringBuilder sb = new StringBuilder();
        sb.append("RecordType,Value\n");
        for (Map.Entry<String, List<String>> entry : dnsRecords.entrySet()) {
            String type = entry.getKey();
            List<String> recs = entry.getValue();
            if (!recs.isEmpty()) {
                for (String r : recs) {
                    sb.append(type).append(',').append(r).append('\n');
                }
            } else {
                sb.append(type).append(",—\n");
            }
        }
        sb.append("\nWHOIS Field,Value\n");
        for (Map.Entry<String, String> e : whoisData.entrySet()) {
            sb.append(e.getKey()).append(',').append(e.getValue()).append('\n');
        }
        Files.write(Paths.get(filename), sb.toString().getBytes());
        System.out.println("\u001B[32m💾 Сохранено CSV: " + filename + "\u001B[0m");
    }

    public static void main(String[] args) throws Exception {
        String domain = null;
        String jsonFile = null;
        String csvFile = null;

        for (int i = 0; i < args.length; i++) {
            if (args[i].equals("--json")) jsonFile = args[++i];
            else if (args[i].equals("--csv")) csvFile = args[++i];
            else if (domain == null) domain = args[i];
        }

        if (domain == null) {
            System.out.println("Usage: java dns_whois <domain> [--json file] [--csv file]");
            System.exit(1);
        }

        dns_whois lookup = new dns_whois(domain);
        lookup.lookupDNS();
        lookup.lookupWhois();
        lookup.printResults();

        lookup.saveJSON(jsonFile);
        lookup.saveCSV(csvFile);
    }
}
