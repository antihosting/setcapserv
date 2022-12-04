# trd

Trigger daemon (TRD) watching for the file updates, and if file changed it automatically runs commands.
Useful tool for automatic low port opening permissions, restarting service. 
```
setcap CAP_NET_BIND_SERVICE=+eip file_name
file_name restart
```

This daemon could be run in root, especially for setcap command.

### Usage

In order to start daemon you need to run

```
./trd -c "echo %1" -c "echo %1" watch_file_path
```

If you want to run daemon in foreground mode, then include option `-f`

```
./trd -f -c "echo %1" watch_file_path
```

Example of setcap and restart
```
./trd -c "setcap CAP_NET_BIND_SERVICE=+eip %1" -c "%1 restart" watch_file_path
```




