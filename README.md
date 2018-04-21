# Backup Rsync
 
Push backups via rsync

## Usage

Cron every hour

```
backup-rsync \
-logtostderr \
-v=0 \
-one-time=false \
-wait=1h \
-host=backupserver.example.com \
-user=backup \
-port=22 \
-privatekey=~/.ssh/id_rsa \
-source=/opt/apache-maven-3.3.9/ \
-target=/backup/
```

Run one time

```
backup-rsync \
-logtostderr \
-v=0 \
-one-time=true \
-host=backupserver.example.com \
-user=backup \
-port=22 \
-privatekey=~/.ssh/id_rsa \
-source=/opt/apache-maven-3.3.9/ \
-target=/backup/
```

## Sampe sudoers file for target /backup and user backup

```
backup ALL=NOPASSWD: /bin/ln -s [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] /backup/current
backup ALL=NOPASSWD: /bin/ln -s empty /backup/current
backup ALL=NOPASSWD: /bin/ls /backup/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] 
backup ALL=NOPASSWD: /bin/ls /backup/current 
backup ALL=NOPASSWD: /bin/mkdir -p /backup/empty
backup ALL=NOPASSWD: /bin/mkdir -p /backup/incomplete/*
backup ALL=NOPASSWD: /bin/mv /backup/incomplete /backup/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]
backup ALL=NOPASSWD: /bin/rm /backup/current
backup ALL=NOPASSWD: /bin/rmdir /backup/empty
backup ALL=NOPASSWD: /usr/bin/rsync --server -logDtprze.iLsfxC --log-format=X --delete-excluded --numeric-ids --link-dest /backup/current/* . /backup/incomplete/*
```
