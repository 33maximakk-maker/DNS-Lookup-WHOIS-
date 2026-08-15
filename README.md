# 🌐 DNS Lookup + WHOIS — узнай всё о домене

> «DNS — это адрес, WHOIS — это паспорт»

**DNS Lookup + WHOIS** — это набор консольных утилит для получения детальной информации о домене: DNS-записи (A, AAAA, MX, NS, TXT, CNAME, SOA) и WHOIS-данные (регистратор, дата регистрации, истечения, контакты).  
Программа позволяет быстро проверить настройки домена и его владельца.

## 🚀 Особенности
- 🔍 Получение всех основных DNS-записей.
- 📋 WHOIS-информация (регистратор, даты, контакты).
- 🎨 Цветной вывод в терминале.
- 📊 Вывод в удобной таблице.
- 💾 Сохранение результатов в JSON и CSV.
- 🌍 Поддержка всех TLD.
- ⚡ Асинхронная/конкурентная обработка нескольких доменов.
- 🖥️ Кроссплатформенная поддержка.

## 🛠️ Установка и запуск

Для работы требуется установленный `whois` (на Linux/macOS) или доступ к публичным WHOIS-серверам.

| OS | Команда установки whois |
|----|--------------------------|
| **Linux (Debian/Ubuntu)** | `sudo apt install whois` |
| **macOS (Homebrew)** | `brew install whois` |
| **Windows** | Скачайте утилиту или используйте библиотеки |

### Запуск

Для каждого языка — минимальные зависимости.

| Язык       | Зависимости                          | Команда запуска                         |
|------------|--------------------------------------|-----------------------------------------|
| Python     | `whois`, `dnspython`                 | `python dns_whois.py example.com`       |
| Go         | `github.com/likexian/whois`          | `go run dns_whois.go example.com`       |
| JavaScript | `whois`, `dns` (встроенный)          | `node dns_whois.js example.com`         |
| Java       | `whois` библиотека (например, `net.whois`) | `javac -cp .:whois.jar ... && java ...` |
| C#         | `Whois` NuGet пакет                  | `dotnet run example.com`                |
| Rust       | `whois`, `trust-dns-resolver`        | `cargo run -- example.com`              |
| Ruby       | `whois` gem                          | `ruby dns_whois.rb example.com`         |
| PHP        | `whois` библиотека                  | `php dns_whois.php example.com`         |

## 📖 Пример использования

```bash
$ python dns_whois.py example.com
Вывод:

text
🌐 DNS Lookup + WHOIS (Python)
Домен: example.com

DNS-записи:
─────────────────────────────────────────
A     93.184.216.34
AAAA  2606:2800:220:1:248:1893:25c8:1946
MX    0 . (no MX record)
NS    a.iana-servers.net
NS    b.iana-servers.net
TXT   "v=spf1 -all"
─────────────────────────────────────────

WHOIS-информация:
─────────────────────────────────────────
Регистратор:   IANA
Дата регистрации: 1995-08-14
Дата истечения:  2024-08-13
DNS-серверы:    a.iana-servers.net, b.iana-servers.net
─────────────────────────────────────────
💾 Сохранено: example_com.json
💾 Сохранено: example_com.csv
🤝 Вклад
Принимаются улучшения, новые языки, фичи.

📜 Лицензия
MIT — используйте свободно.

Автор: Ваш покорный слуга
