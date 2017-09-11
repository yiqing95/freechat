## 静态文件服务

- [servefile 和 fileserver 的区别](https://stackoverflow.com/questions/28793619/golang-what-to-use-http-servefile-or-http-fileserver)

- 实现[野路子](http://www.jb51.net/article/56869.htm)

~~~go

    func staticResource(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        request_type := path[strings.LastIndex(path, "."):]
        switch request_type {
        case ".css":
                w.Header().Set("content-type", "text/css")
        case ".js":
                w.Header().Set("content-type", "text/javascript")
        default:
        } 
        fin, err := os.Open(*realPath + path)
        defer fin.Close()
        if err != nil {
                log.Fatal("static resource:", err)
        } 
        fd, _ := ioutil.ReadAll(fin)
        w.Write(fd)
    }

    http.HandleFunc("/", staticResource)
~~~

- [也是自己实现的](https://github.com/YuriyNasretdinov/social-net/blob/part1/main.go)