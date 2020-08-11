![](https://user-images.githubusercontent.com/22412567/89914023-fb3a6e80-dc26-11ea-82ba-5ed80e2ffb69.jpg)

# qwquiver

> A website for exploring and analyzing exam results.

## Features
- 多平台支持 (Win, macOS, Linux)
- Material Design 风格
- 多条件查询
  - 支持正则表达式
  - 数据模糊查询
  - 根据学校班级查询
- 单科成绩
  - 多视角排名
  - 数据排序
- 表格
  - 固定表头
  - 大屏幕展示
  - 字体大小手动调整
  - 自定义字段显示隐藏
  - 每页项目数量设定
- 快速查看数据平均分
- 下载为 `.xlsx` 表格文件
- 在线直接打印
- 考生数统计
- 趋势，统计图
  - 历史成绩
  - 总分趋势
  - 平均分趋势
  - 多科对比
- 命令行功能
  - 考试管理
  - 成绩导入（从电子表格文件）
  - 构建排名（高考录取人数排名风格）
- 适配手机端 响应式页面

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
$ qwquiver import "20200811.xlsx" --exam-name "期末考试" --exam-conf '{"Grp":"高中","Label":"期末考试","Subj":["YW","SX","YY","WL","HX","SW"],"SubjFullScore":{"HX":100,"SW":100,"SX":150,"WL":100,"YW":150,"YY":150},"Date":"202008","Note":"备注"}'

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

## Screenshots

![](https://user-images.githubusercontent.com/22412567/89917968-0b088180-dc2c-11ea-882e-204382d49818.png)
![](https://user-images.githubusercontent.com/22412567/89917972-0ba11800-dc2c-11ea-9793-d584f1837361.png)
![](https://user-images.githubusercontent.com/22412567/89917976-0c39ae80-dc2c-11ea-8064-dc4a19f79bd8.png)
![](https://user-images.githubusercontent.com/22412567/89917979-0cd24500-dc2c-11ea-9139-cc751d95e07b.png)
![](https://user-images.githubusercontent.com/22412567/89917982-0d6adb80-dc2c-11ea-92c4-cb86774728c6.png)
