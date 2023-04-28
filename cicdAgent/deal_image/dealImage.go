package deal_image

import (
	"bufio"
	"bytes"
	dlogger "cicdAgent/log/zlog"
	sp "cicdAgent/struct_pack"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"strings"
	"time"
)

// harbor auth info
const S_URL = "ip/xxx/"
const S_USERNAME = "******"
const S_PASSWORD = "******"

type DeployBasicInfo sp.DeployBasicInfo

func (di *DeployBasicInfo) BuildImage() (imageName string) {
	fmt.Println("start building image！")
	codeDir := "/tmp/" + di.ProjectName + "-" + di.Branch
	ts := time.Now().Format("20060102150405.000")
	imageName = S_URL + strings.ToLower(di.ProjectName) + ":" + di.Branch + "-" + ts
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		//fmt.Println(err)
		dlogger.Error(err.Error())
		return ""
	}

	tar, err := archive.TarWithOptions(codeDir, &archive.TarOptions{})
	if err != nil {
		dlogger.Error(err.Error())
		return ""
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		NoCache:    true,
		Tags:       []string{imageName},
		BuildArgs:  map[string]*string{"projectName": &di.ProjectName},
		Remove:     true,
	}
	res, err := dockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		dlogger.Error(err.Error())
		return ""
	}

	//write log to file
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		dlogger.Info(scanner.Text())
	}
	return
}

// push image to harbor
func (di *DeployBasicInfo) PushImage(imageName string) (pushImageName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		dlogger.Error(err.Error())
		return
	}
	defer cli.Close()

	authConfig := types.AuthConfig{
		Username: S_USERNAME,
		Password: S_PASSWORD,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		//fmt.Println(err)
		dlogger.Error(err.Error())
		return
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePush(ctx, imageName, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		dlogger.Error("push image err " + err.Error())
		return
	}
	//write log to file
	defer out.Close()
	var buf bytes.Buffer
	io.Copy(&buf, out)
	asString := string(buf.Bytes())
	resultList := strings.Split(asString, "\n")
	for _, value := range resultList {
		dlogger.Info(value)
	}
	pushImageName = imageName
	return
}

func DoImageTask(di *DeployBasicInfo) (result int, pushImageName string) {
	f := "/tmp/" + di.TaskId + "_deploy.log"
	dlogger.InitLogger(f)
	buildName := di.BuildImage()
	if buildName != "" {
		pushImageName = di.PushImage(buildName)
		return 1, pushImageName
	} else {
		dlogger.Error("build image err！！")
		return 0, ""
	}
}
