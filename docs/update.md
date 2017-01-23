# To update the AV run the following:

```bash
$ docker run --name=comodo malice/comodo update
```

## Then to use the updated AVG container:

```bash
$ docker commit comodo malice/comodo:updated
$ docker rm comodo # clean up updated container
$ docker run --rm malice/comodo:updated EICAR
```
