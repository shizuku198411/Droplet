#!/bin/bash

CID=123456
ROOTFS=/
CWD=/
CMD='/usr/bin/python3 -m http.server 8777'
#CMD='/bin/sh -c "echo hello world!; sleep 10"'
#CMD='/bin/sh'
HOSTNAME=mycontainer

IF_NAME=eth0
IF_ADDR=10.166.0.1/24
IF_GW=10.166.0.254
DNS=8.8.8.8

IMAGE_LAYER=/image/path
UPPER_DIR=/upper/path
WORK_DIR=/work/path
MERGE_DIR=/merge/path

OUTDIR=/etc/raind/container/$CID

./bin/droplet spec \
  --rootfs "$ROOTFS" \
  --cwd "$CWD" \
  --command "$CMD" \
  --hostname "$HOSTNAME" \
  --if_name "$IF_NAME" \
  --if_addr "$IF_ADDR" \
  --if_gateway "$IF_GW" \
  --dns "$DNS" \
  --image_layer "$IMAGE_LAYER" \
  --upper_dir "$UPPER_DIR" \
  --work_dir "$WORK_DIR" \
  --merge_dir "$MERGE_DIR" \
  --output "$OUTDIR"