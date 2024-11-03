#!/bin/bash
umount /mnt
kpartx -d /dev/loop0
losetup -d /dev/loop0
