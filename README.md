![](https://user-images.githubusercontent.com/22412567/89914023-fb3a6e80-dc26-11ea-82ba-5ed80e2ffb69.jpg)

# qwquiver

> A website for exploring and analyzing exam results.

## Quick Start

```bash
# 运行服务器
$ qwquiver serve

# 考试管理
$ qwquiver exam

# 列出所有考试
$ qwquiver exam list

# 导入工具帮助文档
$ qwquiver import -h

# 导入考试成绩 xlsx
$ qwquiver import "20200811.xlsx" --exam-name "期末考试" --exam-conf "{"Grp":"高中","Label":"期末考试","Subj":["YW","SX","YY","WL","HX","SW"],"SubjFullScore":{"HX":100,"SW":100,"SX":150,"WL":100,"YW":150,"YY":150},"Date":"202008","Note":"备注"}"

# 编辑考试数据
$ qwquiver exam config set -h

# 获取考试数据
$ qwquiver exam config get -h
```

## Build Setup

```bash
# clone qwquiver project
$ git clone --recurse-submodules https://github.com/qwqcode/qwquiver.git

# install frontend dependencies
$ cd frontend
$ yarn install
$ cd ../

# build qwquiver
$ cd build
$ ./build_frontend.sh # build frontend
$ ./build_win.sh # build for windows
```
