# GOFS
gofs支持简单的路径挂载, 文件基本操作(CRUD). 通过集成aplayer和dplayer实现了部分音频和视频的在线播放, 目前只开放MP3、MP4、flac格式使用播放器, 其他格式使用浏览器默认方式打开. windows安装完成后会自动生成一个名字为GoFileService的服务, Linux目前为手动安装使用.
* DEMO: http://demo.sw0810.com:810
* 配置文件: {gofs}/conf/gofs.json
* 默认端口: 8080
* 默认用户: admin 密码空
* 默认访问地址: 127.0.0.1:8080

# 基本功能
* 文件虚拟路径挂载
![挂载路径](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/mount.png "挂载路径")
* 文件的基本操作(CRUD)
![文件的基本操作](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/filelist.gif "文件的基本操作")
* 图片浏览
![图片浏览](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/picture.png "图片浏览")
* 视频在线播放
![视频在线播放](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/video.png "视频在线播放")
* 音频在线播放
![音频在线播放](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/music.png "音频在线播放")
* 用户管理(CRUD)
![用户管理](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/user.png "用户管理")
* 文件权限(CRUD)
![文件权限](https://github.com/wupeng364/GoFileService/blob/master/readme/imgs/filepermission.png "文件权限管理")
