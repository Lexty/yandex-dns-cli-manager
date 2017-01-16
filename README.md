# Yandex DNS CLI manager

### Get started

    go get github.com/lexty/yandex-dns-cli-manager
    yandex-dns-cli-manager get-token
    # Open URL and follow the instructions
    yandex-dns-cli-manager settings --admin-token <YOUR_ADMIN_TOKEN> --domain <YOR_DOMAIN>
    yandex-dns-cli-manager list

### Usage
```
Yandex DNS CLI manager allows you to change the DNS settings of your domain on pdd.yandex.ru

Usage:
  yandex-dns-cli-manager [command]

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

Use "yandex-dns-cli-manager [command] --help" for more information about a command.
```

### License

MIT