# go-hypixel-api

---

## SDK 文档
- [go-hypixel-api-docs](https://sn0wo2.github.io/go-hypixel-api-docs/)

---

### README 语言

[English](README.md)  
-> 简体中文

---

[![ci](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/ci.yml/badge.svg)](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/ci.yml)
[![lint](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/lint.yml/badge.svg)](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/lint.yml)
[![CodeQL Advanced](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/codeql.yml/badge.svg)](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/codeql.yml)
[![Dependabot Updates](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/Sn0wo2/go-hypixel-api/actions/workflows/dependabot/dependabot-updates)
[![Go Report Card](https://goreportcard.com/badge/github.com/Sn0wo2/go-hypixel-api)](https://goreportcard.com/report/github.com/Sn0wo2/go-hypixel-api)
[![GitHub release](https://img.shields.io/github/v/release/Sn0wo2/go-hypixel-api?color=blue)](https://github.com/Sn0wo2/go-hypixel-api/releases)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)

Golang 编写的 Hypixel API SDK

---

## 简单快捷
- go-hypixel-api 的 API 做到了尽可能简洁, 快速的高效进行开发, 内置了所有 Hypixel API 调用路径并封装成函数

## 最小化
- go-hypixel-api 在设计之初就希望在功能齐全情况下尽可能减小包尺寸和依赖

## 支持API速率限制
- 通过一系列 Hypixel API 返回响应头自动计算剩余请求次数并自动进行堵塞, 无需手动设置

## 自由度
- 我们希望给开发者提供更高的自由度, 使得 go-hypixel-api 的所有设置几乎都可以自主调节, 并暴露底层函数

## 缓存策略支持
- 尽管 go-hypixel-api 并不考虑内置缓存模块, 但是我们提供了较为完善的 Hook 和 Callback 机制, 允许开发者通过自己的缓存策略进行存储

## 快速适配
- 我们会尽快跟进 Hypixel API 的变化, 并在底层实现充足的自由度, 使得开发者在项目未跟进时也能快速进行适配

---

📕 [Hypixel API](https://api.hypixel.net)

⚠️ [Hypixel API 政策](https://developer.hypixel.net/policies)

## 许可证

[Apache-2.0](LICENSE)

* go-hypixel-api 不隶属于 [Hypixel Inc](https://www.hypixel.net) 也未被 [Hypixel Inc](https://www.hypixel.net) 所认可
