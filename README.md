# Momotok-Server

## 抖音项目服务端

具体功能内容参考飞书说明文档

工程除了go.mod包含的依赖外，还使用ffmpeg以生成视频缩略图，需要另外安装并配置进%path%在服务器上  
ffmpeg地址:https://ffbinaries.com/downloads  
其中go-ffmpeg为私人库，编译或开发时不一定能同步到依赖，需要手动操作

```shell
go build && ./Momotok-Server
```
### 快速开始
release中已发布第一个稳定可用版本，解压缩rar文件到空白文件夹，按照**readme.md**的说明建库，然后按照自己的情况修改服务器配置文件（**要记得修改DNS和静态文件URL**），一切准备就绪之后双击 _Momotok-Server.exe_ 启动服务器。APK安装包也包含在压缩包中可供Android设备安装。

### 功能说明

* 各种数据都保存在MySQL数据库中
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/hashed256_video_name 即可，也可以在服务器配置文件修改地址
* 服务器存在配置文件，在system的config.yaml中定义了多个可修改选项

### 正在开发
 1. 代码与数据库优化，找到并修复潜在BUG 

### 目前进度
- [x] basic apis
  - [x] controller.Feed
  - [x] controller.UserInfo
  - [x] controller.Register
  - [x] controller.Login
  - [x] controller.Publish
  - [x] controller.PublishList
- [x] extra apis - I
  - [x] apiRouter.POST("/favorite/action/", controller.FavoriteAction)
  - [x] apiRouter.GET("/favorite/list/", controller.FavoriteList)
  - [x] apiRouter.POST("/comment/action/", controller.CommentAction)
  - [x] apiRouter.GET("/comment/list/", controller.CommentList)
- [x] extra apis - II
  - [x] apiRouter.POST("/relation/action/", controller.RelationAction)
  - [x] apiRouter.GET("/relation/follow/list/", controller.FollowList)
  - [x] apiRouter.GET("/relation/follower/list/", controller.FollowerList)
  - [x] apiRouter.GET("/relation/friend/list/", controller.FriendList)
  - [x] apiRouter.GET("/message/chat/", controller.MessageChat)
  - [x] apiRouter.POST("/message/action/", controller.MessageAction) 

### 建库说明
````mysql
create table user
(
  id                   int auto_increment
    primary key,
  username             varchar(50)                         not null,
  ip                   varchar(15)                         null,
  password             varchar(60)                         null,
  created_at           timestamp default CURRENT_TIMESTAMP not null,
  total_received_likes int       default 0                 null,
  work_count           int       default 0                 null,
  total_likes          int       default 0                 null,
  follow_count         int       default 0                 null,
  follower_count       int       default 0                 null,
  constraint name
    unique (username)
)
  engine = InnoDB;

create table video
(
  id              int auto_increment
    primary key,
  author_id       int                                 null,
  play_url        varchar(255)                        null,
  cover_url       varchar(255)                        null,
  favourite_count int           default 0             null,
  comment_count   int           default 0             null,
  title           varchar(72)                         null,
  publish_time    timestamp default CURRENT_TIMESTAMP null,
  constraint cover_url
    unique (cover_url),
  constraint play_url
    unique (play_url),
  constraint video_ibfk_1
    foreign key (author_id) references user (id)
)
  engine = InnoDB;

create index author_id
  on video (author_id);

create table likes
(
    id       int auto_increment
        primary key,
    video_id int                                 null,
    liked_at timestamp default CURRENT_TIMESTAMP not null,
    user_id  int                                 null,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (video_id) REFERENCES video(id)
)ENGINE = InnoDB;

create index user_id
  on likes (user_id);

create index video_id
  on likes (video_id);

create table comments
(
    id              int auto_increment primary key,
    video_id        int ,
    commenter_id    int ,
    content         text  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ,
    create_date     timestamp default CURRENT_TIMESTAMP ,
    FOREIGN KEY (commenter_id) REFERENCES user(id),
    FOREIGN KEY (video_id) REFERENCES video(id)
)ENGINE = InnoDB;

create index commenter_id
  on comments (commenter_id);

create index video_id
  on comments (video_id);

create table follow_list
(
  id                  int auto_increment primary key,
  follower_uid        int ,
  following_uid       int ,
  FOREIGN KEY (follower_uid) REFERENCES user(id),
  FOREIGN KEY (following_uid) REFERENCES user(id)
)ENGINE = InnoDB;

create table messages
(
  id           int auto_increment
    primary key,
  sender_id    int,
  retriever_id int,
  created_at   bigint,
  message      text,
  constraint messages_retriever_fk
    foreign key (retriever_id) references user (id),
  constraint messages_sender_fk
    foreign key (sender_id) references user (id)
)
  engine = InnoDB;
````
