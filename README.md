# FIND
## Introduction
A tiny command tool for finding and recording useful information.
It seems like a client for a simple key-value database using local data file which you can directly read and write(though in most case you needn't).

## Features
### Orders
#### Find
Example:
```shell
find keyword1 keyword2
```
It'll print the notes whose key **contains** keyword1 **and** keyword2.

#### Add
Example:
```shell
add keyword:content
```
It'll append 'keyword:content' to local data file, so you can find it by the keyword(or part of it) next time.

This order asynchronously updates the backup if the backup service is available.

#### Delete
Example:
```shell
del keyword
```
It'll remove the note whose key **equals** to keyword after a confirmation.

If you don't want the confirmation, try '-f'(means fast) option like this:
```shell
del -f keyword
```
It'll remove the note whose key **equals** to keyword without confirmation.

If you want to batch delete, try '-a'(means all) option like this:
```shell
del -a keyword1 keyword2
```
It'll remove the notes whose key **contains** keyword1 **and** keyword2 after a confirmation.

Certainly you can use '-f' and '-a' at the same time(but be careful).

This order asynchronously updates the backup if the backup service is available.

#### Modify
Example:
```shell
mod keyword:content
```
It'll delete the old note whose key **equals** to keyword without confirmation, and then add 'keyword:content' to local data file.

If the old note doesn't exist, this order is equivalent to add.

This order asynchronously updates the backup if the backup service is available.

#### Weather
Example:
```shell
weather 北京市昌平区
```
It'll print the weather of 北京市昌平区.

Support china address only.

#### Exit
Example:
```shell
exit
```
It'll simply exit the program.

### Backup
FIND only support redis backup service for now and there is no public service provided(I'm sorry /(ㄒoㄒ)/~~).

Feature of backup can be used to sync your note among multiple computers.

### Remind
If the reminder is enabled(which is default), and your note contains 'todo' in key and 'remind@' in value like this:
```shell
todo:fix bug remind@14:00
```
FIND(must be running) will remind you at today's 14:00.

The reminder accept 2 formats of remind-time, which accurate to minutes:
```text
15:04
2006-01-02 15:04
```

FIND support multiple ways to remind you, for now the list is:
1. desktop notification(windows only, which is default)
2. email notification(qq only)

You need an auth code of qq email if you want to use email notification, refer to [qq email auth code](https://service.mail.qq.com/cgi-bin/help?subtype=1&id=28&no=1001256).