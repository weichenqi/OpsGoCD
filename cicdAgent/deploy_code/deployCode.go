package deploy_code

import (
	"bytes"
	dlogger "cicdAgent/log/zlog"
	sp "cicdAgent/struct_pack"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	yaml "gopkg.in/yaml.v2"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// gitlab auth
const GitName = "******"
const GitPassword = "******"

type DeployBasicInfo sp.DeployBasicInfo
type Conf sp.Ci

func (di *DeployBasicInfo) CloneCode() (cloneState int) {
	codeDir := "/tmp/" + di.ProjectName + "-" + di.Branch
	dr, err := os.Open(codeDir)
	defer dr.Close()
	if err != nil {
		dlogger.Error("code dir not exits")
	} else {
		names, err := dr.Readdirnames(-1)
		if err != nil {
			//del old,clone the new code
			os.Remove(codeDir)
		} else {
			for _, name := range names {
				err = os.RemoveAll(codeDir + "/" + name)
			}
		}
	}
	dlogger.Info("start clone code")
	_, err = git.PlainClone(codeDir, false, &git.CloneOptions{
		URL: "http://" + GitName + ":" + GitPassword + "@" + strings.Replace(di.GithttpUrl, "http://", "", -1),
		//URL: di.GitSshUrl,
		ReferenceName: plumbing.ReferenceName(di.BranchName),
		Progress:      os.Stdout,
	})
	if err != nil {
		//fmt.Println(err, di.ProjectName, di.Branch, "pull code err")
		dlogger.Error("pull code err " + err.Error())
		return 0
	}
	dlogger.Info("pull code finished")
	return 1
}

func (di *DeployBasicInfo) YamlFileParser() (*Conf, error) {
	fileName := "/tmp/" + di.ProjectName + "-" + di.Branch + "/" + "ci.yaml"
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var conf Conf
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		dlogger.Error("yaml file err")
		return nil, fmt.Errorf("in file %q: %v", fileName, err)
	}
	return &conf, nil

}

func (di *DeployBasicInfo) DeployJar() (deployStat int) {
	deployStat = 0
	dlogger.Info("start deploy Jar ")
	codeDir := "/tmp/" + di.ProjectName + "-" + di.Branch
	timeUnixNano := time.Now().UnixNano()
	ts := strconv.FormatInt(timeUnixNano, 10)
	deployImgName := "deployjar-" + di.ProjectName + "-" + di.Branch + "-" + ts
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		//fmt.Println(err)
		dlogger.Error("deploy Jar err " + err.Error())
	}
	defer cli.Close()

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
	if err != nil {
		deployStat = 0
		dlogger.Error("deploy Jar err" + err.Error())
		//fmt.Println(err)
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		//panic(err)
		dlogger.Error("deploy Jar err" + err.Error())
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			deployStat = 0
			dlogger.Error("deploy Jar err" + err.Error())
			return
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		deployStat = 0
		dlogger.Error("deploy Jar err" + err.Error())
		return
	}

	//write log to file
	var buf bytes.Buffer
	io.Copy(&buf, out)
	asString := string(buf.Bytes())
	//logger.Info(asString)
	resultList := strings.Split(asString, "\n")
	for _, value := range resultList {
		dlogger.Info(value)
	}
	deployStat = 1

	//Remove the container for deploy Jar
	containerName := deployImgName
	// Get a list of containers that match the given name
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.Arg("name", containerName))})
	if err != nil {
		//fmt.Println(err)
		dlogger.Error(err.Error())
	}

	// If there are no containers with the given name, return an error
	if len(containers) == 0 {
		dlogger.Error("no containers with the given name found " + err.Error())
	}
	containerID := containers[0].ID
	// Remove the container
	if err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{}); err != nil {
		//fmt.Println(err)
		dlogger.Error(err.Error())
	}
	return
}

func (di *DeployBasicInfo) DeployPy() (deployStat int) {
	fmt.Println("Hello golang，waiting")
	return
}

func (di *DeployBasicInfo) DeployGo() (deployStat int) {
	fmt.Println("Hello python,waiting")
	return
}

func DoDeployTask(di *DeployBasicInfo) (s int) {
	f := "/tmp/" + di.TaskId + "_deploy.log"
	dlogger.InitLogger(f)
	cloneState := di.CloneCode()
	if cloneState == 1 {
		conf, _ := di.YamlFileParser()
		di.DeployCmd = strings.Split(conf.Build.Script, " ")
		di.DeployBaseImage = conf.Build.Image
		di.Language = conf.Build.Language
		di.WorkDir = conf.Build.Workdir
		switch di.Language {
		case "java":
			s = di.DeployJar()
		case "golang":
			s = di.DeployGo()
		case "python":
			s = di.DeployPy()
		default:
			fmt.Println("Can not find suitable language")
		}
	}

	return
}
