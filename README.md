# Prometheus Summary

```
______                      _____
| ___ \                    /  ___|
| |_/ / __ ___  _ __ ___   \ `--. _   _ _ __ ___  _ __ ___   __ _ _ __ _   _
|  __/ '__/ _ \| '_ ` _ \   `--. \ | | | '_ ` _ \| '_ ` _ \ / _` | '__| | | |
| |  | | | (_) | | | | | | /\__/ / |_| | | | | | | | | | | | (_| | |  | |_| |
\_|  |_|  \___/|_| |_| |_| \____/ \__,_|_| |_| |_|_| |_| |_|\__,_|_|   \__, |
                                                                        __/ |
                                                                       |___/
```

A lazy tool written by Golang to export Prometheus summary in different format:

- JSON.
- YAML.
- CSV.
- Plain Text table.

```bash
+-----------------------+----------------------------+--------+-------------------+--------------------------+---------------------------+-----------------------+------------------+--------------------------------+
|         NAME          |          ADDRESS           | STATUS | STORAGE RETENTION | NUMBER OF ACTIVE TARGETS | NUMBER OF DROPPED TARGETS | NUMBER OF TIME SERIES | NUMBER OF CHUNKS | NUMBER OF INGESTED SAMPLES PER |
|                       |                            |        |                   |                          |                           |                       |                  |            SECONDS             |
+-----------------------+----------------------------+--------+-------------------+--------------------------+---------------------------+-----------------------+------------------+--------------------------------+
|          prometheus_1 |    http://fakeaddress:9091 |     OK |               90d |                      897 |                         0 |                     0 |                0 |                                |
|          prometheus_2 |    http://fakeaddress:9092 |     OK |               60d |                      829 |                      1698 |               2387664 |          2387664 |                                |
+-----------------------+----------------------------+--------+-------------------+--------------------------+---------------------------+-----------------------+------------------+--------------------------------+

```

## How to use

- Get the executable file from [bin folder](./bin/prom-summary)
- Check out the usage.

```bash
bin/prom-summary --help                                                                                                                                                        prom-summary/prom-summary -> master ? ! |•
usage: prom-summary [<flags>]

A lazy tool written by Golang to export Prometheus summary.

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).
  --config.file="etc/config.yml"
          Prom-summary configuration file path.
```

- Prepare the config file, you can find the sample config file [here](./etc/config.yml).
- Run it!

```bash
bin/prom-summary --config.file /tmp/config.yml

+-----------------------+----------------------------+--------+-------------------+--------------------------+---------------------------+-----------------------+------------------+--------------------------------+
|         NAME          |          ADDRESS           | STATUS | STORAGE RETENTION | NUMBER OF ACTIVE TARGETS | NUMBER OF DROPPED TARGETS | NUMBER OF TIME SERIES | NUMBER OF CHUNKS | NUMBER OF INGESTED SAMPLES PER |
|                       |                            |        |                   |                          |                           |                       |                  |            SECONDS             |
+-----------------------+----------------------------+--------+-------------------+--------------------------+---------------------------+-----------------------+------------------+--------------------------------+
|          prometheus_1 |    http://fakeaddress:9091 |     OK |               90d |                      897 |                         0 |                     0 |                0 |                                |
|          prometheus_2 |    http://fakeaddress:9092 |     OK |               60d |                      829 |                      1698 |               2387664 |          2387664 |                                |
+-----------------------+----------------------------+--------+-------------------+--------------------------+---------------------------+-----------------------+------------------+--------------------------------+

```
