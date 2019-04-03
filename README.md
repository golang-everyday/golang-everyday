# golang-everyday

[![jaywcjlove/sb](https://jaywcjlove.github.io/sb/ico/awesome.svg)](https://github.com/sindresorhus/awesome) [![jaywcjlove/sb](https://jaywcjlove.github.io/sb/lang/chinese.svg)](README-zh.md) [![License](https://img.shields.io/github/license/golang-everyday/golang-everyday.svg)](https://jitpack.io/#Coder-zheng/blockchain-note)  [![Stars](https://img.shields.io/github/stars/golang-everyday/golang-everyday.svg)](https://jitpack.io/#golang-everyday/golang-everyday)  [![Forks](https://img.shields.io/github/forks/golang-everyday/golang-everyday.svg)](https://jitpack.io/#golang-everyday/golang-everyday) [![Issues](https://img.shields.io/github/issues/golang-everyday/golang-everyday.svg)](https://jitpack.io/#golang-everyday/golang-everyday)
[![Author](https://img.shields.io/badge/Author-GolangEverydayGroup-black.svg?)](https://github.com/golang-everyday)
[![Author](https://img.shields.io/badge/QQ-812397431-yellow.svg?)](http://wpa.qq.com/msgrd?v=3&uin=812397431&site=qq&menu=yes)



## 目录

* [Golang语法](#Golang语法)

* [Golang开发技巧](#Golang开发技巧)
* [编程题](#编程题)
* [书籍下载](#书籍下载)
* [开源项目](#开源项目)
* [博客链接](#博客链接)
* [开发工具](#开发工具)
* [贡献者](#贡献者)



## Golang语法

## Golang开发技巧

- 多个 if 语句尽量折叠成 switch
- 尽量用 `chan struct{}` 来传递信号, `chan bool` 表达的不够清楚
- 总是把 for-select 换成一个函数
- 分组定义 `const` 类型声明和 `var` 逻辑类型声明
- 为整型常量值实现 `Stringer` 接口
- 用 defer 来检查你的错误

## 编程题

**[⬆ 返回顶部](#目录)**

## 经典面试题

* 使用两个 goroutine 交替打印一个打印数字一个字母   最后结果   12AB34CD56EF78GH910IJ [解答](https://play.golang.org/p/CWWN5kl8Mpx)

**[⬆ 返回顶部](#目录)**

## 书籍下载

**[⬆ 返回顶部](#目录)**

## 开源项目

**[⬆ 返回顶部](#目录)**

## 博客链接

**[⬆ 返回顶部](#目录)**

## 开发工具

* **vscode 插件**

  | 插件名称               | 插件描述                                                     |
  | ---------------------- | ------------------------------------------------------------ |
  | GitLens                | 非常方便的查看文件代码的 commit 信息（提交时间，提交人等）。 |
  | Code Runner            | 针对非常多的语言而快速方便执行的小插件。                     |
  | filesize               | 在 VSCode 底部工具栏，非常方便的显示文件大小。               |
  | Go                     | Go 语言插件。                                                |
  | Terminal               | 命令行工具插件。                                             |
  | Vim                    | Vim 插件                                                     |
  | VSCode Great Icons     | VSCode 美化不同的文件。                                      |
  | WakaTime               | 统计项目代码的时间。                                         |
  | BetterComments         | 代码注释                                                     |
  | Beautify               | 格式化 js ，json，html 代码。                                |
  | Auto Import            | 自动倒包                                                     |
  | Bookmarks              | 好用的书签                                                   |
  | Bracket Pair Colorizer | 多种颜色括号，结构清晰明了                                   |
  | Code Runner            | 一键运行代码                                                 |

* 建议配置CDPATH 环境变量，这样我们在任何地方都能进入 github 下的项目了

  ![temp_paste_image_4486809a7891f1f021f88d41c9750b7c](https://ws3.sinaimg.cn/large/006tKfTcly1g1ofrduix4g30ep09vdfv.gif)

* [vscode代码补全](https://github.com/Microsoft/vscode-go/blob/master/snippets/go.json)

**[⬆ 返回顶部](#目录)**

