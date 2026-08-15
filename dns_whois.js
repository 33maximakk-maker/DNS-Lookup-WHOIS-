// dns_whois.js — JavaScript версия

const dns = require('dns');
const { promisify } = require('util');
const whois = require('whois');
const fs = require('fs');
const path = require('path');

const resolve4 = promisify(dns.resolve4);
const resolve6 = promisify(dns.resolve6);
const resolveMx = promisify(dns.resolveMx);
const resolveNs = promisify(dns.resolveNs);
const resolveTxt = promisify(dns.resolveTxt);
const resolveCname = promisify(dns.resolveCname);

class DNSWhois {
    constructor(domain) {
        this.domain = domain;
        this.dnsRecords = {};
        this.whoisData = {};
    }

    async lookupDNS() {
        const recordTypes = ['A', 'AAAA', 'MX', 'NS', 'TXT', 'CNAME'];
        for (const type of recordTypes) {
            try {
                let result;
                switch (type) {
                    case 'A':
                        result = await resolve4(this.domain);
                        break;
                    case 'AAAA':
                        result = await resolve6(this.domain);
                        break;
                    case 'MX':
                        result = await resolveMx(this.domain);
                        result = result.map(m => `${m.priority} ${m.exchange}`);
                        break;
                    case 'NS':
                        result = await resolveNs(this.domain);
                        break;
                    case 'TXT':
                        result = await resolveTxt(this.domain);
                        result = result.map(t => t.join(''));
                        break;
                    case 'CNAME':
                        result = await resolveCname(this.domain);
                        result = [result];
                        break;
                }
                this.dnsRecords[type] = result || [];
            } catch (err) {
                this.dnsRecords[type] = [];
            }
        }
    }

    lookupWhois() {
        return new Promise((resolve, reject) => {
            whois.lookup(this.domain, (err, data) => {
                if (err) {
                    this.whoisData = { error: err.message };
                    return resolve();
                }
                this.whoisData = { raw: data };
                // Парсинг упрощённый, для демонстрации
                const lines = data.split('\n');
                const fields = ['Registrar', 'Creation Date', 'Expiration Date', 'Name Server', 'Status'];
                for (const line of lines) {
                    for (const field of fields) {
                        if (line.startsWith(field)) {
                            const parts = line.split(':');
                            if (parts.length > 1) {
                                const key = field.replace(/ /g, '_').toLowerCase();
                                this.whoisData[key] = parts.slice(1).join(':').trim();
                            }
                        }
                    }
                }
                resolve();
            });
        });
    }

    printResults() {
        console.log('\x1b[36m🌐 DNS Lookup + WHOIS (JavaScript)\x1b[0m');
        console.log(`Домен: ${this.domain}\n`);

        console.log('\x1b[32mDNS-записи:\x1b[0m');
        console.log('─'.repeat(50));
        for (const [type, records] of Object.entries(this.dnsRecords)) {
            if (records.length > 0) {
                for (const r of records) {
                    console.log(`\x1b[33m${type.padEnd(6)}\x1b[0m ${r}`);
                }
            } else {
                console.log(`\x1b[33m${type.padEnd(6)}\x1b[0m —`);
            }
        }
        console.log('─'.repeat(50));

        console.log('\x1b[32mWHOIS-информация:\x1b[0m');
        console.log('─'.repeat(50));
        if (this.whoisData.error) {
            console.log(`\x1b[31mОшибка WHOIS: ${this.whoisData.error}\x1b[0m`);
        } else {
            const fields = ['registrar', 'creation_date', 'expiration_date', 'name_server', 'status'];
            const labels = ['Регистратор:', 'Дата регистрации:', 'Дата истечения:', 'DNS-серверы:', 'Статус:'];
            for (let i = 0; i < fields.length; i++) {
                const val = this.whoisData[fields[i]] || 'N/A';
                console.log(`\x1b[33m${labels[i]}\x1b[0m ${val}`);
            }
        }
        console.log('─'.repeat(50));
    }

    saveJSON(filename) {
        if (!filename) filename = this.domain.replace(/\./g, '_') + '.json';
        const data = {
            domain: this.domain,
            timestamp: new Date().toISOString(),
            dns: this.dnsRecords,
            whois: this.whoisData
        };
        fs.writeFileSync(filename, JSON.stringify(data, null, 2));
        console.log(`\x1b[32m💾 Сохранено JSON: ${filename}\x1b[0m`);
    }

    saveCSV(filename) {
        if (!filename) filename = this.domain.replace(/\./g, '_') + '.csv';
        let csv = 'RecordType,Value\n';
        for (const [type, records] of Object.entries(this.dnsRecords)) {
            if (records.length > 0) {
                for (const r of records) {
                    csv += `${type},${r}\n`;
                }
            } else {
                csv += `${type},—\n`;
            }
        }
        csv += '\nWHOIS Field,Value\n';
        for (const [key, val] of Object.entries(this.whoisData)) {
            if (key === 'raw') continue;
            csv += `${key},${val}\n`;
        }
        fs.writeFileSync(filename, csv);
        console.log(`\x1b[32m💾 Сохранено CSV: ${filename}\x1b[0m`);
    }
}

async function main() {
    const args = process.argv.slice(2);
    let domain = null;
    let jsonFile = null;
    let csvFile = null;

    for (let i = 0; i < args.length; i++) {
        if (args[i] === '--json') jsonFile = args[++i];
        else if (args[i] === '--csv') csvFile = args[++i];
        else if (!domain) domain = args[i];
    }

    if (!domain) {
        console.log('Usage: node dns_whois.js <domain> [--json file] [--csv file]');
        process.exit(1);
    }

    const lookup = new DNSWhois(domain);
    await lookup.lookupDNS();
    await lookup.lookupWhois();
    lookup.printResults();

    lookup.saveJSON(jsonFile);
    lookup.saveCSV(csvFile);
}

main().catch(console.error);
