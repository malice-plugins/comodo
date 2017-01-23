malice-comodo
=============

[![Circle CI](https://circleci.com/gh/maliceio/malice-comodo.png?style=shield)](https://circleci.com/gh/maliceio/malice-comodo) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/comodo.svg)](https://hub.docker.com/r/malice/comodo/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/comodo.svg)](https://hub.docker.com/r/malice/comodo/) [![Docker Image](https://img.shields.io/badge/docker image-609 MB-blue.svg)](https://hub.docker.com/r/malice/comodo/)

This repository contains a **Dockerfile** of [comodo](https://www.comodo.com) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/comodo/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

-	[ubuntu:precise (*104 MB*\)](https://hub.docker.com/_/ubuntu/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/comodo/) from public [DockerHub](https://hub.docker.com): `docker pull malice/comodo`

### Usage

```
docker run --rm malice/comodo EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/comodo FILE

Usage: comodo [OPTIONS] COMMAND [arg...]

Malice Comodo AntiVirus Plugin

Version: v0.1.0, BuildTime: 20170122

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --table, -t	       output as Markdown table
  --callback, -c	    POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x	       proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --timeout value       malice plugin timeout (in seconds) (default: 60) [$MALICE_TIMEOUT]    
  --elasitcsearch value elasitcsearch address for Malice to store results [$MALICE_ELASTICSEARCH]   
  --help, -h	        show help
  --version, -v	     print the version

Commands:
  update	Update virus definitions
  web       Create a Comodo scan web service  
  help		Shows a list of commands or help for one command

Run 'comodo COMMAND --help' for more information on a command.
```

Sample Output
-------------

### JSON:

```json
{
  "comodo": {
    "infected": true,
    "result": "Malware",
    "engine": "1.1",
    "updated": "20170122"
  }
}
```

### STDOUT (Markdown Table):

---

#### Comodo

| Infected | Result  | Engine    | Updated  |
|----------|---------|-----------|----------|
| true     | Malware | 13.0.3114 | 20170122 |

---

Documentation
-------------

-	[To write results to ElasticSearch](https://github.com/maliceio/malice-comodo/blob/master/docs/elasticsearch.md)
-	[To create a Comodo scan micro-service](https://github.com/maliceio/malice-comodo/blob/master/docs/web.md)
-	[To post results to a webhook](https://github.com/maliceio/malice-comodo/blob/master/docs/callback.md)
-	[To update the AV definitions](https://github.com/maliceio/malice-comodo/blob/master/docs/update.md)

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-comodo/issues/new).

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice-comodo/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice-comodo/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice-comodo/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

MIT Copyright (c) 2016-2017 **blacktop**
