package slo

import (
	"context"
	"g.hz.netease.com/horizon/pkg/pipeline/manager"
	pipelinerunmodels "g.hz.netease.com/horizon/pkg/pipelinerun/models"
)

const (
	BuildTask            = "build"
	BuildTaskDisplayName = "构建"
	GitStep              = "git"
	CompileStep          = "compile"
	ImageStep            = "image"

	DeployTask            = "deploy"
	DeployTaskDisplayName = "发布（环境准备）"
	DeployStep            = "deploy"
)

var (
	envMapping = map[string][]string{
		"test":   {"perf", "reg", "test"},
		"online": {"pre", "online"},
	}

	// 构建RT临界值
	_buildRT uint = 60
	// 构建可用率SLO
	_buildRequestSLO = 99.99

	// 发布RT临界值
	_deployRT uint = 30
	// 发布可用率SLO
	_deployRequestSLO = 99.99
)

type Controller interface {
	PipelineSLO(ctx context.Context, environment string, start, end int64) (pipelineSLOs []*PipelineSLO, err error)
}

type controller struct {
	pipelineManager manager.Manager
}

func (c controller) PipelineSLO(ctx context.Context, environment string,
	start, end int64) (pipelineSLOs []*PipelineSLO, err error) {
	slos, err := c.pipelineManager.ListPipelineSLOsByEnvsAndTimeRange(ctx, envMapping[environment], start, end)
	if err != nil {
		return nil, err
	}

	pipelineSLOMap := map[string]*PipelineSLO{
		BuildTask: {
			Name:                BuildTask,
			DisplayName:         BuildTaskDisplayName,
			Count:               0,
			RequestAvailability: 0,
			RequestSlo:          _buildRequestSLO,
			RTAvailability:      0,
			RT:                  _buildRT,
		},
		DeployTask: {
			Name:                DeployTask,
			DisplayName:         DeployTaskDisplayName,
			Count:               0,
			RequestAvailability: 0,
			RequestSlo:          _deployRequestSLO,
			RTAvailability:      0,
			RT:                  _deployRT,
		},
	}

	buildTaskCount, buildSuccessCount, buildRTSuccessCount := 0, 0, 0
	deployTaskCount, deploySuccessCount, deployRTSuccessCount := 0, 0, 0
	for _, slo := range slos {
		if build, ok := slo.Tasks[BuildTask]; ok {
			buildTaskCount++
			if build.Result == pipelinerunmodels.ResultOK {
				buildSuccessCount++
				// 这里注意是用整体Task的耗时减掉compile step的耗时，这样的结果更加准确，包含了Pod启动准备所需的时间
				if build.Duration-slo.Tasks[BuildTask].Steps[CompileStep].Duration < pipelineSLOMap[BuildTask].RT {
					buildRTSuccessCount++
				}
			} else {
				// 如果是compile失败了，slo维度也认为是成功的
				if compile, ok := slo.Tasks[BuildTask].Steps[CompileStep]; ok && compile.Result == pipelinerunmodels.ResultFailed {
					buildSuccessCount++
				}
			}
		}
		if deploy, ok := slo.Tasks[DeployTask]; ok {
			deployTaskCount++
			if deploy.Result == pipelinerunmodels.ResultOK {
				deploySuccessCount++
				if deploy.Duration < pipelineSLOMap[DeployTask].RT {
					deployRTSuccessCount++
				}
			}
		}
	}

	pipelineSLOMap[BuildTask].Count = buildTaskCount
	pipelineSLOMap[BuildTask].RequestAvailability = float64(buildSuccessCount) * 100 / float64(buildTaskCount)
	pipelineSLOMap[BuildTask].RTAvailability = float64(buildRTSuccessCount) * 100 / float64(buildTaskCount)

	pipelineSLOMap[DeployTask].Count = deployTaskCount
	pipelineSLOMap[DeployTask].RequestAvailability = float64(deploySuccessCount) * 100 / float64(deployTaskCount)
	pipelineSLOMap[DeployTask].RTAvailability = float64(deployRTSuccessCount) * 100 / float64(deployTaskCount)

	for _, slo := range pipelineSLOMap {
		pipelineSLOs = append(pipelineSLOs, slo)
	}
	return pipelineSLOs, nil
}

// NewController initializes a new group controller
func NewController() Controller {
	return &controller{
		pipelineManager: manager.Mgr,
	}
}
