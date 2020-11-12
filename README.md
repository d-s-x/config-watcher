### config-watcher
This is a program to trigger a reload when Kubernetes ConfigMaps are updated.
It watches mounted volume dirs and executes a command when a config map has been changed.

### Example

```
config-watcher                   \
  --volume-dir=/etc/config       \
  --command=/bin/sh              \
  --argument=reload-prometheus.sh
```

### License

This project is [Apache Licensed](LICENSE.txt)

