# cmd/job

注册定时任务请参考 [demo.go](./demo.go)。

查看所有定时任务
```bash
go run main.go job list
```

执行一次某个任务
```bash
go run main.go job once foo

go run main.go job once --http=true baz
```

调度所有定时任务
```bash
go run main.go job
```
