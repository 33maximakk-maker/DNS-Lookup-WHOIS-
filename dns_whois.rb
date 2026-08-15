# dns_whois.rb — Ruby версия

require 'resolv'
require 'whois'
require 'whois-parser'
require 'json'
require 'csv'
require 'time'
require 'optparse'

class DNSWhois
  def initialize(domain)
    @domain = domain
    @dns_records = {}
    @whois_data = {}
  end

  def lookup_dns
    resolver = Resolv::DNS.new
    types = %w[A AAAA MX NS TXT CNAME]
    types.each do |type|
      records = []
      begin
        case type
        when 'A'
          records = resolver.getaddresses(@domain).map(&:to_s).select { |ip| ip =~ /^\d+\./ }
        when 'AAAA'
          records = resolver.getaddresses(@domain).map(&:to_s).select { |ip| ip.include?(':') }
        when 'MX'
          records = resolver.getresources(@domain, Resolv::DNS::Resource::IN::MX)
                    .map { |mx| "#{mx.preference} #{mx.exchange.to_s}" }
        when 'NS'
          records = resolver.getresources(@domain, Resolv::DNS::Resource::IN::NS)
                    .map { |ns| ns.name.to_s }
        when 'TXT'
          records = resolver.getresources(@domain, Resolv::DNS::Resource::IN::TXT)
                    .map { |txt| txt.data.join(' ') }
        when 'CNAME'
          cname = resolver.getresource(@domain, Resolv::DNS::Resource::IN::CNAME)
          records = [cname.name.to_s] if cname
        end
      rescue
        records = []
      end
      @dns_records[type] = records
    end
  end

  def lookup_whois
    begin
      w = Whois::Client.new.lookup(@domain)
      @whois_data = {
        'registrar' => w.registrar.name,
        'creation_date' => w.creation_date.to_s,
        'expiration_date' => w.expiration_date.to_s,
        'name_servers' => w.name_servers.join(', '),
        'status' => w.status.join(', ')
      }
    rescue => e
      @whois_data = { 'error' => e.message }
    end
  end

  def print_results
    puts "\e[36m🌐 DNS Lookup + WHOIS (Ruby)\e[0m"
    puts "Домен: #{@domain}\n"

    puts "\e[32mDNS-записи:\e[0m"
    puts "─" * 50
    @dns_records.each do |type, records|
      if records.any?
        records.each do |r|
          puts "\e[33m#{type.ljust(6)}\e[0m #{r}"
        end
      else
        puts "\e[33m#{type.ljust(6)}\e[0m —"
      end
    end
    puts "─" * 50

    puts "\e[32mWHOIS-информация:\e[0m"
    puts "─" * 50
    if @whois_data['error']
      puts "\e[31mОшибка WHOIS: #{@whois_data['error']}\e[0m"
    else
      puts "\e[33mРегистратор:\e[0m   #{@whois_data['registrar']}"
      puts "\e[33mДата регистрации:\e[0m #{@whois_data['creation_date']}"
      puts "\e[33mДата истечения:\e[0m  #{@whois_data['expiration_date']}"
      puts "\e[33mDNS-серверы:\e[0m    #{@whois_data['name_servers']}"
      puts "\e[33mСтатус:\e[0m        #{@whois_data['status']}"
    end
    puts "─" * 50
  end

  def save_json(filename)
    filename ||= @domain.gsub('.', '_') + '.json'
    data = {
      domain: @domain,
      timestamp: Time.now.iso8601,
      dns: @dns_records,
      whois: @whois_data
    }
    File.write(filename, JSON.pretty_generate(data))
    puts "\e[32m💾 Сохранено JSON: #{filename}\e[0m"
  end

  def save_csv(filename)
    filename ||= @domain.gsub('.', '_') + '.csv'
    CSV.open(filename, 'w') do |csv|
      csv << ['RecordType', 'Value']
      @dns_records.each do |type, records|
        if records.any?
          records.each { |r| csv << [type, r] }
        else
          csv << [type, '—']
        end
      end
      csv << []
      csv << ['WHOIS Field', 'Value']
      @whois_data.each { |k, v| csv << [k, v] }
    end
    puts "\e[32m💾 Сохранено CSV: #{filename}\e[0m"
  end
end

def main
  options = {}
  OptionParser.new do |opts|
    opts.banner = "Usage: ruby dns_whois.rb <domain> [--json file] [--csv file]"
    opts.on("--json FILE", "Сохранить JSON") { |v| options[:json] = v }
    opts.on("--csv FILE", "Сохранить CSV") { |v| options[:csv] = v }
  end.parse!

  domain = ARGV[0]
  unless domain
    puts "Usage: ruby dns_whois.rb <domain> [--json file] [--csv file]"
    exit 1
  end

  lookup = DNSWhois.new(domain)
  lookup.lookup_dns
  lookup.lookup_whois
  lookup.print_results

  lookup.save_json(options[:json])
  lookup.save_csv(options[:csv])
end

main if __FILE__ == $0
