# DO-dyndns
Automatically update DNS record via Digitaloceans API

This program lets you update a DNS entry of domains you manage via Digitalocean.
It queries the current IP from icanhazip.com (currently IPv4 only), and then updates the domain record if necessary.

## Usage

There are 2 modes do-dyndns can run in.

First: Update IP once and then immediately exit. For example, set up a cronjob, which then checks the IP every hour.
(If the IP has not changed, the program will not call the DO API)

Second: Start a small Webserver, which listens on a predefined port, and updates when someone accesses the server.
This can be used for example with pfsense/opnsense router, where one can configure custom dyndns services.
Just run the program somwhere behind the router, then point the dyndns config to the listening port.
Because do-dyndns currently gets the IP from icanhazip.com, it is not needed to use authentication (for now).
Nevertheless, one should not run this open to the internet, so that noone can DOS the server and/or the DO-API.

## Configuration
It uses a simple json settings file. Per default, the server looks at "/usr/local/etc/do-dyndns.json" or "./do-dyndns.json".
Alternatively, provide a different path as fist command line argument.

### Options:

* domain: the base domain, set up via digitalocean. (example.com)
* subdomain: the subdomain which should be updated with the IP Address. (dyndns; so the updated URL would be dyndns.example.com)
* token: the dyndns API token.
    * Needs Read/Write permissions to update the record, so this file should be guarded appropiately. (eg. `chmod 600 ./do-dyndns.json`, `chmod root:wheel ./do.dyndns.json`)
* httpServer: boolean, wether do-dyndns should update the record once, or start a http server.
* httpPort: port on which the server should listen on.
