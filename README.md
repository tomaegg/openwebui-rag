## Milvus Management

浏览器进入milvus [管理界面](http://localhost:9091/webui)

## Developement

```shell
go tool task ...
```

## TODO

- [x] 基于python实现文本分割, 写一个起一个简单的web服务, http+json传输即可
- [x] 完善go操纵milvus的逻辑, 增加插入和查询操作
- [ ] 修改cmd/embed部分, 支持批量处理文件, 调用python的split api, 向量插入milvus

### 方案1

- [ ] 编写go-server由openwebui调用, 入参topic, server根据topic从milvus中检索内容,
直接让openwebui那边连接到ollama, 出10个选择题

### 方案2

- [ ] 编写go-server由openwebui调用, 入参topic, server根据topic从milvus中检索内容, 出题功能需要在server调用ollama实现.
