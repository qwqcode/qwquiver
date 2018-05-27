QwQuery
============================

<p align="center"><img src="./docs/screenshots/bootstrap.png"></p>

QwQuery 是一个 Material Design 风格的学生考试成绩统计站

> 未经允许代码和衍生品不得用于商业用途，侵权必究

Features
------------
- Material Design
- 根据 `学生姓名` `学校` `班级` 查询
- 每个科目成绩排序，`正序` `倒序`
- 表格
    - 固定表头
    - 全屏显示
    - 手动调整字体大小
    - 隐藏指定字段
    - 设定每页显示项目数量
- 快速查看每组数据的平均分
- 直接保存为 `.xls` 电子表格文件
- 直接打印
- 显示考生数
- 趋势，统计图
    - 历史成绩
    - 总分趋势
    - 市平均分趋势
- 控制台功能
    - 构建排名（录取人数风格）
    - [统一原本不同的学校名和班级名（姓名作为桥梁）](./docs/两堆数据统一原本不同的学校名和班级名的方法.md)
- 适配手机端 响应式页面

Quick Start
------------

```sh
$ composer install
$ php -r "copy('config/db.example.php', 'config/db.php');copy('config/params.example.php', 'config/params.php');"
```
Donate
------------
如果您觉得我的项目对您有帮助，并且您愿意给予我一点小小的支持，您可以通过以下方式向我捐赠，这样可以维持项目持续地发展，非常感谢！ヽ(•̀ω•́ )ゝ

If you are enjoying this app, please consider making a donation to keep it alive.

| Alipay | Wechat | 
| :------: | :------: | 
| <img width="150" src="./docs/donate/alipay.png"> | <img width="150" src="./docs/donate/wechat.png"> | 

捐赠者的名字将保存于 [捐赠者列表](https://github.com/Zneiat/donate-qwqaq)，非常感谢你们的支持

Author
------------
[ZNEIAT](http://www.qwqaq.com)

Using
------------
- [yiisoft/yii2](https://github.com/yiisoft/yii2)
- [antvis/g2](https://github.com/antvis/g2)

Screenshots
------------

<p align="center">
<img src="./docs/screenshots/home.png">
<img src="./docs/screenshots/charts.png">
<img src="./docs/screenshots/full-screen.png">
<img src="./docs/screenshots/phone.png">
</p>
