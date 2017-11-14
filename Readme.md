# cloudgo

---

## 1. 概览
基于现有 web 库，编写一个类似 [cloudgo](http://blog.csdn.net/pmlpml/article/details/78404838#t4) 的简单 web 应用。

#### 程序逻辑：
> hello user 的 web 服务。

#### 使用方法
```
Usage of cloudgo:
  -p, --port string   PORT for httpd listening (default "8080")
```

## 2. 使用的框架
程序使用了[fasthttp](https://github.com/valyala/fasthttp)库以及支持fasthttp的[fasthttprouter](https://github.com/buaazp/fasthttprouter)路由库，结合起来使用，是非常不错的web server框架，号称比官方库[net/http](https://golang.org/pkg/net/http/)快10倍以上。
快的原因
当然，`fasthttp`也有缺陷：

 - 现在还不支持HTTP/2.0，也不支持WebSocket。
 - 不提供与官方库[net/http](https://golang.org/pkg/net/http/)相同的API。

但在十倍的性能提升下，决定尝试一下`fasthttp`。

## 3. 测试
####`curl`测试
测试结果如下。显然请求得到了正确回复：`{"Test":"Hello testuser"}`。
```
$ curl -v http://localhost:9090/hello/testuser
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 9090 (#0)
> GET /hello/testuser HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9090
> Accept: */*
>
< HTTP/1.1 200 OK
* Server fasthttp is not blacklisted
< Server: fasthttp
< Date: Tue, 14 Nov 2017 15:26:54 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 26
<
{"Test":"Hello testuser"}
* Connection #0 to host localhost left intact
```

####`ab`测试
参数说明参考了[CentOS服务器Http压力测试之ab](http://linux.it.net.cn/CentOS/fast/2015/0715/16393.html)。
测试结果如下。总体而言，在并发个数为100，总请求个数为1000的情况下，服务器平均每秒收到`3725.74`个请求，服务器平均请求等待时间为`0.268`ms，用户平均请求等待时间为`26.840`ms。50％的用户请求等待时间小于`26`ms，最长的用户请求等待时间为`46`ms。
```
$ ab -n1000 -c100 http://localhost:9090/hello/your      ###-n 执行的请求数量 -c 并发请求个数
This is ApacheBench, Version 2.3 <$Revision: 1604373 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/
Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests
Server Software:        fasthttp
Server Hostname:        localhost
Server Port:            9090
Document Path:          /hello/your                     ###请求的资源
Document Length:        22 bytes                        ###文档返回的长度，不包括相应头
Concurrency Level:      100                             ###并发个数
Time taken for tests:   0.268 seconds                   ###总请求时间
Complete requests:      1000
Failed requests:        0
Total transferred:      176000 bytes                
HTML transferred:       22000 bytes                 
Requests per second:    3725.74 [#/sec] (mean)              ###吞吐率:平均每秒的请求数
Time per request:       26.840 [ms] (mean)                  ###用户平均请求等待时间
Time per request:       0.268 [ms] (mean, across all concurrent requests)   ###服务器平均请求等待时间
Transfer rate:          640.36 [Kbytes/sec] received        ###传输速率
Connection Times (ms)                       ###用户请求响应时间可分成网络链接，系统处理和等待三部分。
              min  mean[+/-sd] median   max                 ###标准差(数值越大系统响应时间越不稳定)
Connect:        0   12   5.2     13      21
Processing:     1   14   6.6     14      31
Waiting:        1    9   5.2      9      22
Total:          9   26   7.6     26      46
Percentage of the requests served within a certain time (ms)
  50%     26                                        ###50%的请求都在 26 ms内完成
  66%     28                                        ###66%的请求都在 28 ms内完成
  75%     30                                        ###75%的请求都在 30 ms内完成
  80%     32                                        ###80%的请求都在 32 ms内完成
  90%     36
  95%     40
  98%     42
  99%     44
 100%     46 (longest request)
```
增加并发请求个数为之前十倍的测试结果如下。而在并发个数为1000，总请求个数为10000的情况下，服务器平均每秒收到`6125.06`个请求(之前的`1.64`倍)，服务器平均请求等待时间为`0.163`ms(之前的`0.608`倍)，用户平均请求等待时间为`163.264`ms(之前的`6.08`倍)。50％的用户请求等待时间小于`23`ms(之前的`0.88`倍)，最长的用户请求等待时间为`1179`ms(之前的`25.63`倍)。
```
$ ab -n10000 -c1000 http://localhost:9090/hello/your
... ...
Concurrency Level:      1000
Time taken for tests:   1.633 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1760000 bytes
HTML transferred:       220000 bytes
Requests per second:    6125.06 [#/sec] (mean)
Time per request:       163.264 [ms] (mean)
Time per request:       0.163 [ms] (mean, across all concurrent requests)
Transfer rate:          1052.74 [Kbytes/sec] received
Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0  112 290.7     11    1084
Processing:     3   31  28.8     13     216
Waiting:        0   25  24.6     11     210
Total:          4  142 309.3     23    1179
Percentage of the requests served within a certain time (ms)
  50%     23
  66%     57
  75%    111
  80%    119
  90%    162
  95%   1136
  98%   1162
  99%   1168
 100%   1179 (longest request)
```
明显看到，并发个数增加的情况下，用户平均请求等待时间也随着增加，约一半用户的请求等待时间并没有受到影响，而另一半用户却需要经受长时间的排队等待，请求等待时间显著增长。用户请求等待时间大部分都花费在网络链接上了，由Connect均值为`112`(之前的`9.33`倍)，Processing均值为`31`(之前的`2.21`倍)，Waiting均值为`25`(之前的`2.77`倍)可看出。

但总体而言，`fasthttp`的表现还是非常可观的。
