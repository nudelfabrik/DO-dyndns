# GOdyndns
Automatically update DNS record via Digitalocean or Gandi.net API

This program lets you update a DNS entry of domains you manage via Digitalocean or Gandis v5 API
It queries the current IP from icanhazip.com (currently IPv4 only), and then updates the domain record if necessary.

## Usage

There are 2 modes GOdyndns can run in.

First: Update IP once and then immediately exit. For example, set up a cronjob, which then checks the IP every hour.
(If the IP has not changed, the program will not call the API)

Second: Start a small Webserver, which listens on a predefined port, and updates when someone accesses the server.
This can be used for example with pfsense/opnsense router, where one can configure custom dyndns services.
Just run the program somwhere behind the router, then point the dyndns config to the listening port.
Because godyndns currently gets the IP from icanhazip.com, it is not needed to use authentication (for now).
Nevertheless, one should not run this open to the internet, so that noone can DOS the server and/or the DO-API.

## Configuration
It uses a simple json settings file. Per default, the server looks at "/usr/local/etc/godyndns.json" or "./godyndns.json".
Alternatively, provide a different path with -f flag.

### Options:

* domain: the base domain, set up via digitalocean or Gandi. (example.com)
* subdomain: the subdomain which should be updated with the IP Address. (dyndns; so the updated URL would be dyndns.example.com)
* token: the API token.
    * Needs Read/Write permissions to update the record, so this file should be guarded appropiately. (eg. `chmod 600 ./godyndns.json`, `chmod root:wheel ./godyndns.json`)
* httpServer: boolean, wether do-dyndns should update the record once, or start a http server.
* httpPort: port on which the server should listen on.
