#!/bin/bash -e

LOOPDEV=/dev/loop0
PARTDEV=/dev/loop0p1
TARGET=/mnt
IMAGE=$1

losetup $LOOPDEV $IMAGE
echo "* Dispositivo: $LOOPDEV -> $IMAGE"

partprobe /dev/loop0
echo "* Partição:    $PARTDEV"

mount -o ro,sync $PARTDEV /mnt
echo "* Montado:     $TARGET"
