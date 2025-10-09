package server

import (
	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/gin-gonic/gin"
	"github.com/robotism/gitinsight/gitinsight"
)

func RegisterRoute(g *gin.RouterGroup) {
	g.GET("/commits", GetCommits)
	g.GET("/contributors", GetContributors)
	g.GET("/branches", GetRepoBranches)
	g.GET("/ranking", GetRanking)
	g.GET("/heatmap", GetCommitHeatmap)
	g.GET("/period", GetCommitPeriod)
}

func getFilterFromContext(c *gin.Context) *gitinsight.CommitLogFilter {
	since := c.Query("since")
	until := c.Query("until")
	repos := c.Query("repos")
	branches := c.Query("branches")
	authors := c.Query("authors")
	isMerge := c.Query("isMerge")
	messageType := c.Query("messageType")
	period := c.Query("period")

	commitHash := c.Query("commitHash")
	leEffective := c.Query("leEffective")
	geEffective := c.Query("geEffective")

	offset, err := xcast.ToIntE(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := xcast.ToIntE(c.Query("limit"))
	if err != nil {
		limit = 50
	}
	if since == "" {
		since = GetConfig().Insight.Since
	}
	filter := &gitinsight.CommitLogFilter{
		Offset:      offset,
		Limit:       limit,
		DateFrom:    since,
		DateTo:      until,
		RepoUrl:     repos,
		BranchName:  branches,
		CommitHash:  commitHash,
		Nickname:    authors,
		IsMerge:     isMerge,
		MessageType: messageType,
		Period:      period,
		LeEffective: leEffective,
		GeEffective: geEffective,
	}
	return filter
}

func GetRanking(c *gin.Context) {
	filter := getFilterFromContext(c)
	ranking, err := gitinsight.GetRanking(filter)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"meta": gin.H{
				"since": filter.DateFrom,
				"until": filter.DateTo,
			},
			"data": ranking,
		})
	}
}

func GetRepoBranches(c *gin.Context) {
	filter := getFilterFromContext(c)
	branches, err := gitinsight.GetRepoBranches(filter)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"meta": gin.H{
				"since": filter.DateFrom,
				"until": filter.DateTo,
			},
			"data": branches,
		})
	}
}

func GetContributors(c *gin.Context) {
	filter := getFilterFromContext(c)
	contributors, err := gitinsight.GetAuthors(filter)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"meta": gin.H{
				"since": filter.DateFrom,
				"until": filter.DateTo,
			},
			"data": contributors,
		})
	}
}

func GetCommits(c *gin.Context) {
	filter := getFilterFromContext(c)
	count, err := gitinsight.CountCommitLogs(filter)

	config := []gitinsight.Auth{}
	for _, auth := range GetConfig().Insight.Auths {
		config = append(config, gitinsight.Auth{
			Domain:        auth.Domain,
			CommitUrlTmpl: auth.CommitUrlTmpl,
		})
	}

	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	commits, err := gitinsight.GetCommitLogs(filter)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"meta": gin.H{
				"offset": filter.Offset,
				"limit":  filter.Limit,
				"since":  filter.DateFrom,
				"until":  filter.DateTo,
				"total":  count,
				"config": config,
			},
			"data": commits,
		})
	}
}

func GetCommitHeatmap(c *gin.Context) {
	filter := getFilterFromContext(c)
	data, err := gitinsight.GetCommitHeatmapData(filter)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"meta": gin.H{
				"since": filter.DateFrom,
				"until": filter.DateTo,
			},
			"data": data,
		})
	}
}

func GetCommitPeriod(c *gin.Context) {
	filter := getFilterFromContext(c)
	data, err := gitinsight.GetCommitStatsByPeriodAndUser(filter)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"meta": gin.H{
				"since": filter.DateFrom,
				"until": filter.DateTo,
			},
			"data": data,
		})
	}
}
