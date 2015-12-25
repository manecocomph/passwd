# passwd
存储泄漏的用户名和密码到MongoDB 并提供简单的页面查询 a golang project

clone 后, 到 passwd/passwd/ 路径下  git install  编译打包，打包文件为 passwd

命令行执行:
xiatian@kali:~/goWork/bin$ ./passwd -h
Usage of ./passwd:
  -i	import files to MongoDB
  -s	start the server (default true)
  -t	test if configuration is fine

有2个配置文件：
config.json  配置MongoDB的连接参数(默认没有用户名/密码) 和本地服务器的端口
import.cfg   配置要导入的源数据的文件名和格式

