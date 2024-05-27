#!/bin/sh

# 获取系统架构信息
ARCH=$(uname -m)

echo "平台: ${ARCH}"

# 获取系统位数
BITS=$(getconf LONG_BIT)

# 使用 readelf 获取 ELF 标头信息
if [ "$ARCH" = "armv7l" ]; then
  ARM_VERSION=$(readelf -A /proc/self/exe | grep Tag_CPU_arch: | awk '{print $2}')
fi

# 获取最新的标签名称
LATEST_TAG=$(curl -s https://api.github.com/repos/AirportR/miaospeed/releases/latest | grep 'tag_name' | cut -d '"' -f 4)

if [ "$ARCH" = "x86_64" ] && [ "$BITS" = "64" ]; then
  echo "架构: linux/amd64"
  wget -O /opt/miaospeed.gz https://github.com/AirportR/miaospeed/releases/download/"$LATEST_TAG"/miaospeed-linux-amd64-"$LATEST_TAG".gz
  gzip -d /opt/miaospeed.gz
  chmod +x /opt/miaospeed
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  echo "架构: linux/arm64"
  wget -O /opt/miaospeed.gz https://github.com/AirportR/miaospeed/releases/download/"$LATEST_TAG"/miaospeed-linux-arm64-"$LATEST_TAG".gz
  gzip -d /opt/miaospeed.gz
  chmod +x /opt/miaospeed
elif [ "$ARCH" = "armv7l" ] && [ "$ARM_VERSION" = "v6KZ" ]; then
  echo "架构: linux/arm/v6"
  wget -O /opt/miaospeed.gz https://github.com/AirportR/miaospeed/releases/download/"$LATEST_TAG"/miaospeed-linux-armv6-"$LATEST_TAG".gz
  gzip -d /opt/miaospeed.gz
  chmod +x /opt/miaospeed
elif [ "$ARCH" = "armv7l" ] && [ "$ARM_VERSION" = "v7" ]; then
  echo "架构: linux/arm/v7"
  wget -O /opt/miaospeed.gz https://github.com/AirportR/miaospeed/releases/download/"$LATEST_TAG"/miaospeed-linux-armv7-"$LATEST_TAG".gz
  gzip -d /opt/miaospeed.gz
  chmod +x /opt/miaospeed
elif [ "$ARCH" = "x86_64" ] && [ "$BITS" = "32" ]; then
  echo "架构: linux/386"
  wget -O /opt/miaospeed.gz https://github.com/AirportR/miaospeed/releases/download/"$LATEST_TAG"/miaospeed-linux-386-"$LATEST_TAG".gz
  gzip -d /opt/miaospeed.gz
  chmod +x /opt/miaospeed
fi