# observer

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/imkira/go-observer/blob/master/LICENSE.txt)
[![GoDoc](https://godoc.org/github.com/imkira/go-observer?status.svg)](https://godoc.org/github.com/imkira/go-observer)
[![Build Status](http://img.shields.io/travis/imkira/go-observer.svg?style=flat)](https://travis-ci.org/imkira/go-observer)
[![Coverage](https://codecov.io/gh/imkira/go-observer/branch/master/graph/badge.svg)](https://codecov.io/gh/imkira/go-observer)
[![codebeat badge](https://codebeat.co/badges/28bdd579-8b34-4940-a3e0-35ac52794a42)](https://codebeat.co/projects/github-com-imkira-go-observer)
[![goreportcard](https://goreportcard.com/badge/github.com/imkira/go-observer)](https://goreportcard.com/report/github.com/imkira/go-observer)

本项目从 [go-observer](https://github.com/imkira/go-observer) 分叉并简化而来，用于发送广播。

## 问题

在 Go 语言中，使用 channel 给各个 goroutine 发送消息的方式如下：

```go
for _, channel := range channels {
  channel <- value
}
```

这种方式存在两个缺陷：

- 一旦某个 channel 处于阻塞状态，此 channel 以及后续，都无法收到消息。
- 这个方法的复杂度是 O(N)，观察者越多，越消耗资源。

## 解决方法

```go
type state struct {
  value interface{}
  next  *state
  done  chan struct{}
}
```

- value: 用于记录信息
- next: 指向下一个 state
- done: 当此 state.next 指向的新 state 时，关闭 done。利用关闭后的 channel 总是可以获取信息的特性，告诉 observer 还有后续 state。

包中含有两个接口

- Property: 相等于 publisher
- Stream: 相等于 observer

## 内存占用

Stream 的内存占用取决于其长度，由读取速度最慢的 observer 决定。使用时，请确保各个 observer 都能尽快地读取。

## 使用方法

首先，需要安装

```text
go get -u github.com/aQuaYi/observer
```

然后导入

```go
import "github.com/aQuaYi/observer"
```

The following example creates one property that is updated every second by one
or more publishers, and observed by one or more observers.

### 文档

更多使用方法，可以查看[在线文档](https://godoc.org/github.com/aQuaYi/observer).

### 示例：Property 和 Publisher

```go
val := 1
prop := observer.NewProperty(val) // 创建了一个 Property，具有初始值 1
for {
  time.Sleep(time.Second)
  val += 1
  fmt.Printf("will publish value: %d\n", val)
  prop.Update(val) // 每次都利用 val 的新值更新 Property
}
```

注意：

- Property 是线程安全的，可以复制到多个 Publisher 同时更新

### 示例: Observer

```go
stream := prop.Observe() // 从 Property 生成 Stream
for {
  val := stream.Value().(int) // 获取 Stream 中第一个 state 的值
  stream.Wait() // 当 stream 到达尾部时候，会发生阻塞
}
```

注意：

- Stream **不**是线程安全的。必须使用 ```Property.Observe()``` 或 ```Stream.Clone()``` 创建新的 Stream。

### 实例

请前往
[examples/multiple.go](https://github.com/aQuaYi/observer/blob/master/examples/multiple.go)
查看多个 Publisher 和多个 Observer 的简单例子。
