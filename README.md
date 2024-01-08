# SLAPC

SSH into server, See logs and pull code.
Written in go, currently mainly for laravel applications on ubuntu servers.

## Setup

Create a new json file called config.json in the same directory as the executable and fill in the following data.

```json
{
    "known_hosts": "path/to/known_hosts",
    "private_key": "path/to/private_key",
    "servers": [
        {
            "name": "any identifier",
            "host": "1.2.3.4",
            "port": "22",
            "user": "root",
            "pass": "password",
            "path": "/var/www/html"
        }
    ]
}
```

Add as many server entries as you want to to the servers array.

## Todo

When the config.json file is not found, open the default text editor automatically and save to the right path that slapc will read.
