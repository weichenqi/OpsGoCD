# OpsGoCD
The one with customizable and flexible CD system


### 项目名字来源
 
一个ops用golang写的一个CD（CI）系统，功能未完善且比较初级，不喜勿喷。

### 功能：

可以扩充大量工作节点，并可以高度自定义流水线环节的CD系统，目前功能主要从提交代码到上传镜像仓库为止，相关状态数据已经入库，可以通过开发界面，加入审批等功能，形成一个可视化的CI/CD系统。

### 设计图
     补充ing
     
### 部署：
  1.部署mysql，版本随意，主流版本即可，导入三个sql。
  
  2.部署redis，版本随意，主流版本即可。
  
  3.修改数据库相关信息：
  
  mysql：配置在cicdSerevr/mypack_mysql/myMysql.go中

    func InitConn() *sql.DB {
      // Set up the DSN
      dsn := "cicd_user:******@tcp(127.0.0.1:3306)/wcicd"
      // Open a connection to the database
      sqldb, err := sql.Open("mysql", dsn)
      if err != nil {
        dlogger.Error(err.Error())
        os.Exit(1)
    }
  
  
  redis：配置在cicdSerevr/mypack_redis/myRedis.go中
  
    func InitRedisConn() *redis.Client {
      rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "******", // no password set
        DB:       0,        // use default DB
        PoolSize: 8,
      })
      err := rdb.Ping(ctx).Err()
      if err != nil {
        dlogger.Error(err.Error())
        os.Exit(1)
      }
      return rdb
  }
  
  
  不提出配置文件的目的在于一个二进制包就能启动，不依赖于其他文件。
  
  4.目前只支持java项目，maven的私服setting.xml文件需要自己准备，不同的maven镜像的工作目录和config目录可能所有区别，需要修改相关路径等。
  目前支持在ci.yaml自定义镜像的工作目录，不支持定义settings所在的目录。
  
  需要匹配setting.xml路径涉及的文件在Agent机器上，具体文件在cicdAgent/deploy_code/deployCode.go
  
  func (di *DeployBasicInfo) DeployJar() (deployStat int) {
  ...

    resp, err := cli.ContainerCreate(ctx, &container.Config{
      Image:      di.DeployBaseImage, //"maven:3.3-jdk-8"
      Cmd:        di.DeployCmd,       //"mvn package -Dmaven.test.skip=true"
      Tty:        false,
      WorkingDir: di.WorkDir,
    }, &container.HostConfig{
      AutoRemove: false,
      //settings.xml in /tmp/cicdAgent/conf
      Binds: []string{codeDir + ":" + di.WorkDir, "/tmp/cicdAgent/conf:/usr/share/maven/ref/"}},
      nil, nil, deployImgName)
 ...
	}
  
  5.自动构建依赖git仓库中有ci.yaml文件存在，目的是为了定义项目构建的命令。
  
    build:
    language: java
    image: maven:3.3-jdk-8
    workdir: /usr/src/mymaven
    script: mvn package -Dmaven.test.skip=true
    
    
  6.agent需要安装docker环境，定义好镜像仓库等。
  
  
  7.分别编译运行server端和agent端。
  
  
  8.go初学者，没有做test case，抱歉。
  
  9.待完成工作，支持golang和python。
  
  10.有bug欢迎提issue

  
