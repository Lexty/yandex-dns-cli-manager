# Yandex DNS CLI manager

```
Yandex DNS CLI manager allows you to change the DNS settings of your domain on pdd.yandex.ru

Usage:
  yandexdns [command]

Available Commands:
  add         Add a new DNS record
  delete      Delete the DNS record by ID
  edit        Edit DNS record
  get-token   Instruction for getting token
  list        The list of the DNS records
  settings    Show or change settings
  version     Print the version of YandexDns

Flags:
  -a, --admin-token="": admin's token
	  --config="": config file (default is $HOME/.yandexdns.json)
  -d, --domain="": domain name
  -h, --help[=false]: help for yandexdns

Use "yandexdns [command] --help" for more information about a command.
```