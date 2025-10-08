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

func GetRanking(c *gin.Context) {
	since := c.Query("since")
	until := c.Query("until")
	repos := c.Query("repos")
	branches := c.Query("branches")
	authors := c.Query("authors")
	if since == "" {
		since = GetConfig().Insight.Since
	}
	ranking, err := gitinsight.GetRanking(since, until, repos, branches, authors)
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
				"since": since,
				"until": until,
			},
			"data": ranking,
		})
	}
}

func GetRepoBranches(c *gin.Context) {
	since := c.Query("since")
	until := c.Query("until")
	if since == "" {
		since = GetConfig().Insight.Since
	}
	repo := c.Query("repo")
	branches, err := gitinsight.GetRepoBranches(since, until, repo)
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
				"since": since,
				"until": until,
			},
			"data": branches,
		})
	}
}

func GetContributors(c *gin.Context) {
	since := c.Query("since")
	until := c.Query("until")
	if since == "" {
		since = GetConfig().Insight.Since
	}
	contributors, err := gitinsight.GetAuthors(since, until)
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
				"since": since,
				"until": until,
			},
			"data": contributors,
		})
	}
}

func GetCommits(c *gin.Context) {
	since := c.Query("since")
	until := c.Query("until")
	repos := c.Query("repos")
	branches := c.Query("branches")
	authors := c.Query("authors")
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
	count, err := gitinsight.CountCommitLogs(&gitinsight.CommitLogFilter{
		DateFrom:   since,
		DateTo:     until,
		RepoUrl:    repos,
		BranchName: branches,
		Nickname:   authors,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	commits, err := gitinsight.GetCommitLogs(&gitinsight.CommitLogFilter{
		DateFrom:   since,
		DateTo:     until,
		RepoUrl:    repos,
		BranchName: branches,
		Nickname:   authors,
	}, offset, limit)

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
				"offset": offset,
				"limit":  limit,
				"since":  since,
				"until":  until,
				"total":  count,
			},
			"data": commits,
		})
	}
}

func GetCommitHeatmap(c *gin.Context) {
	since := c.Query("since")
	until := c.Query("until")
	repos := c.Query("repos")
	branches := c.Query("branches")
	authors := c.Query("authors")
	messageType := c.Query("messageType")
	isMerge := c.Query("isMerge")
	if since == "" {
		since = GetConfig().Insight.Since
	}
	data, err := gitinsight.GetCommitHeatmapData(&gitinsight.CommitLogFilter{
		DateFrom:    since,
		DateTo:      until,
		RepoUrl:     repos,
		BranchName:  branches,
		Nickname:    authors,
		MessageType: messageType,
		IsMerge:     isMerge,
	})
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
				"since": since,
				"until": until,
			},
			"data": data,
		})
	}
}

func GetCommitPeriod(c *gin.Context) {
	since := c.Query("since")
	until := c.Query("until")
	repos := c.Query("repos")
	branches := c.Query("branches")
	authors := c.Query("authors")
	messageType := c.Query("messageType")
	isMerge := c.Query("isMerge")
	period := c.Query("period")

	if since == "" {
		since = GetConfig().Insight.Since
	}
	data, err := gitinsight.GetCommitStatsByPeriodAndUser(&gitinsight.CommitLogFilter{
		DateFrom:    since,
		DateTo:      until,
		RepoUrl:     repos,
		BranchName:  branches,
		Nickname:    authors,
		MessageType: messageType,
		IsMerge:     isMerge,
	}, period)
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
				"since": since,
				"until": until,
			},
			"data": data,
		})
	}
}
