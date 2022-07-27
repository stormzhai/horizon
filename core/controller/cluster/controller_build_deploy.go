package cluster

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"g.hz.netease.com/horizon/core/common"
	"g.hz.netease.com/horizon/pkg/cluster/code"
	codemodels "g.hz.netease.com/horizon/pkg/cluster/code"
	"g.hz.netease.com/horizon/pkg/cluster/registry"
	"g.hz.netease.com/horizon/pkg/cluster/tekton"
	prmodels "g.hz.netease.com/horizon/pkg/pipelinerun/models"
	regionmodels "g.hz.netease.com/horizon/pkg/region/models"
	"g.hz.netease.com/horizon/pkg/util/wlog"

	"github.com/mozillazg/go-pinyin"
)

func (c *controller) BuildDeploy(ctx context.Context, clusterID uint,
	r *BuildDeployRequest) (_ *BuildDeployResponse, err error) {
	const op = "cluster controller: build deploy"
	defer wlog.Start(ctx, op).StopPrint()

	currentUser, err := common.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cluster, err := c.clusterMgr.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return nil, err
	}

	var gitRef, gitRefType = cluster.GitRef, cluster.GitRefType
	if r.Git != nil {
		if r.Git.Commit != "" {
			gitRefType = codemodels.GitRefTypeCommit
			gitRef = r.Git.Commit
		} else if r.Git.Tag != "" {
			gitRefType = codemodels.GitRefTypeTag
			gitRef = r.Git.Tag
		} else if r.Git.Branch != "" {
			gitRefType = codemodels.GitRefTypeBranch
			gitRef = r.Git.Branch
		}
	}

	commit, err := c.commitGetter.GetCommit(ctx, cluster.GitURL, gitRefType, gitRef)
	if err != nil {
		return nil, err
	}

	regionEntity, err := c.regionMgr.GetRegionEntity(ctx, cluster.RegionName)
	if err != nil {
		return nil, err
	}

	// 1. create project in harbor
	harbor := c.registryFty.GetByHarborConfig(ctx, &registry.HarborConfig{
		Server:          regionEntity.Harbor.Server,
		Token:           regionEntity.Harbor.Token,
		PreheatPolicyID: regionEntity.Harbor.PreheatPolicyID,
	})
	if _, err := harbor.CreateProject(ctx, application.Name); err != nil {
		return nil, err
	}

	// 2. update image in git repo
	imageURL := assembleImageURL(regionEntity, application.Name, cluster.Name, gitRef, commit.ID)

	configCommit, err := c.clusterGitRepo.GetConfigCommit(ctx, application.Name, cluster.Name)
	if err != nil {
		return nil, err
	}

	// 3. add pipelinerun in db
	pr := &prmodels.Pipelinerun{
		ClusterID:        clusterID,
		Action:           prmodels.ActionBuildDeploy,
		Status:           string(prmodels.StatusCreated),
		Title:            r.Title,
		Description:      r.Description,
		GitURL:           cluster.GitURL,
		GitRefType:       gitRefType,
		GitRef:           gitRef,
		GitCommit:        commit.ID,
		ImageURL:         imageURL,
		LastConfigCommit: configCommit.Master,
		ConfigCommit:     configCommit.Gitops,
	}
	prCreated, err := c.pipelinerunMgr.Create(ctx, pr)
	if err != nil {
		return nil, err
	}

	// 4. create pipelinerun in k8s
	tektonClient, err := c.tektonFty.GetTekton(cluster.EnvironmentName)
	if err != nil {
		return nil, err
	}
	clusterFiles, err := c.clusterGitRepo.GetCluster(ctx,
		application.Name, cluster.Name, cluster.Template)
	if err != nil {
		return nil, err
	}

	_, err = tektonClient.CreatePipelineRun(ctx, &tekton.PipelineRun{
		Application:   application.Name,
		ApplicationID: application.ID,
		Cluster:       cluster.Name,
		ClusterID:     cluster.ID,
		Environment:   cluster.EnvironmentName,
		Git: tekton.PipelineRunGit{
			URL:       cluster.GitURL,
			Subfolder: cluster.GitSubfolder,
			Commit:    commit.ID,
		},
		ImageURL:         imageURL,
		Operator:         currentUser.GetEmail(),
		PipelinerunID:    prCreated.ID,
		PipelineJSONBlob: clusterFiles.PipelineJSONBlob,
		Region:           cluster.RegionName,
		RegionID:         regionEntity.ID,
		Template:         cluster.Template,
	})
	if err != nil {
		return nil, err
	}

	return &BuildDeployResponse{
		PipelinerunID: prCreated.ID,
	}, nil
}

func assembleImageURL(regionEntity *regionmodels.RegionEntity,
	application, cluster, branch, commit string) string {
	// domain is harbor server
	domain := strings.TrimPrefix(regionEntity.Harbor.Server, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	// time now
	timeFormat := "20060102150405"
	timeStr := time.Now().Format(timeFormat)

	// normalize branch
	args := pinyin.Args{
		Fallback: func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		},
	}
	normalizedBranch := strings.Join(pinyin.LazyPinyin(branch, args), "")
	normalizedBranch = regexp.MustCompile(`[^a-zA-Z0-9_.-]`).ReplaceAllString(normalizedBranch, "_")

	return fmt.Sprintf("%v/%v/%v:%v-%v-%v",
		domain, application, cluster, normalizedBranch, commit[:8], timeStr)
}

func (c *controller) GetDiff(ctx context.Context, clusterID uint, refType, ref string) (_ *GetDiffResponse, err error) {
	const op = "cluster controller: get diff"
	defer wlog.Start(ctx, op).StopPrint()

	// 1. get cluster
	cluster, err := c.clusterMgr.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	// 2. get application
	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return nil, err
	}

	// 3. get code commit
	var commit *code.Commit
	if ref != "" {
		commit, err = c.commitGetter.GetCommit(ctx, cluster.GitURL, refType, ref)
		if err != nil {
			return nil, err
		}
	}

	// 4.  get config diff
	diff, err := c.clusterGitRepo.CompareConfig(ctx, application.Name, cluster.Name, nil, nil)
	if err != nil {
		return nil, err
	}
	return ofClusterDiff(cluster.GitURL, refType, ref, commit, diff), nil
}

func ofClusterDiff(gitURL, refType, ref string, commit *code.Commit, diff string) *GetDiffResponse {
	var codeInfo *CodeInfo

	// TODO: support any gitlab or gitlab not only internal
	if commit != nil {
		// git@github.com:demo/demo.git
		var historyLink string
		if strings.HasPrefix(gitURL, common.InternalGitSSHPrefix) {
			httpURL := common.InternalSSHToHTTPURL(gitURL)
			historyLink = httpURL + common.CommitHistoryMiddle + ref
		}
		codeInfo = &CodeInfo{
			CommitID:  commit.ID,
			CommitMsg: commit.Message,
			Link:      historyLink,
		}
		switch refType {
		case codemodels.GitRefTypeTag:
			codeInfo.Tag = ref
		case codemodels.GitRefTypeBranch:
			codeInfo.Branch = ref
		}
	}

	return &GetDiffResponse{
		CodeInfo:   codeInfo,
		ConfigDiff: diff,
	}
}
