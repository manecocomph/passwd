# passwd
存储泄漏的用户名和密码到MongoDB 并提供简单的页面查询 a golang project

clone 后, 到 passwd/passwd/ 路径下  git install  编译打包，打包文件为 passwd<br>
命令行执行:<br>
kaili@kaili: ./passwd -h<br>
Usage of ./passwd:<br>
&nbsp;&nbsp;&nbsp;&nbsp;-i import files to MongoDB<br>  
&nbsp;&nbsp;&nbsp;&nbsp;-s start the server (default true)<br>
&nbsp;&nbsp;&nbsp;&nbsp;-t test if configuration is fine<br>

有2个配置文件：<br>
&nbsp;&nbsp;&nbsp;&nbsp;config.json  配置MongoDB的连接参数(默认没有用户名/密码) 和本地服务器的端口<br>
&nbsp;&nbsp;&nbsp;&nbsp;import.cfg   配置要导入的源数据的文件名和格式<br>
