<h4 align="center">GithubFast 是一个用python编写的更新github相关域名的dns解析。</h4>
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/wjlin0/GitHubFast?filename=go.mod" alt="">
    <a href="https://github.com/wjlin0/GithubFast">
    <img src=" https://img.shields.io/github/stars/wjlin0/GitHubFast" alt="">
    </a>
<a href="https://github.com/wjlin0/GithubFast/releases"><img src="https://img.shields.io/github/downloads/wjlin0/GithubFast/total" alt=""></a> 
    <a href="https://github.com/wjlin0/GithubFast">
    <img src="https://img.shields.io/github/last-commit/wjlin0/GithubFast" alt="">
    </a>

<a href="https://wjlin0.com/"><img src="https://img.shields.io/badge/wjlin0-blog-green" alt=""></a>


</p>
# 介绍

由于github dns污染的问题，导致国内访问速度较慢，为了解决这个问题做了这个工具
# 安装
```shell
git clone https://github.com/wjlin0/GithubFast
cd GithubFast && go mod tidy && go build
./GithubFast
```
OR
```shell
go install github.com/wjlin0/GithubFast
GithubFast
```

# 运行
```shell
# 确保 是管理员权限 也就是最高权限运行
GithubFast
```