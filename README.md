# trd

Trigger daemon (TRD) watching the file updates, and if the file changed automatically run commands.
Modern and fancy `Control-M` for automatic jobs/commands triggering if some files get changed.

Useful tool for automatic tasking and low port permissions opening in linux systems, restarting services automatically and unpacking packages. 
```
setcap CAP_NET_BIND_SERVICE=+eip file_name
file_name restart
```

This daemon could be run in root, especially for setcap command.

### Usage

The special string `%1` is using in command line expression and would be replaced by the actual watching file.

In order to start daemon you need to run:

```
./trd -c "echo %1" -c "echo %1" watch_file_path
```

If you want to run daemon in foreground mode, then include option `-f`

```
./trd -f -c "echo %1" watch_file_path
```

In order to verbose all events incluse option `-v`
```
./trd -v -f -c "echo %1" watch_file_path
```

Example of unpacking executable under user permissions:
```
./trd -c "gzip -d -f %1" /my/path/distr.gz
```

Example of setcap and restart under root permissions:
```
./trd -c "setcap CAP_NET_BIND_SERVICE=+eip %1" -c "%1 restart" /my/path/service_linux
```

Commands would be executed in the sequence as in command line.



