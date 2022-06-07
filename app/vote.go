package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vote_cj/config"
	"vote_cj/dbpool"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
)

func Vote(c *gin.Context) {
	openid, bo := c.Get("openid")
	if !bo {
		c.JSON(440, gin.H{"msg": "非法请求"})
		c.Abort()
		return
	}
	votemod := CreateVoteMod(c.Param("cand_id"), c.Param("group_id"), openid.(string))
	votemod.rc = dbpool.Redispool.Get()

	defer votemod.rc.Close()
	err := votemod.CheckData()
	if err != nil {
		c.JSON(440, gin.H{"msg": "非法请求"})
		c.Abort()
		return
	}
	err = votemod.CheckVoteStartOrEnd()
	if err != nil {
		c.JSON(441, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	err = votemod.CheckCanVote()
	if err != nil {
		c.JSON(442, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	err = votemod.SetCache()
	if err != nil {
		log.Println(err.Error())
		c.JSON(541, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	err = votemod.Update()
	if err != nil {
		c.JSON(443, gin.H{"msg": "投票失败"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "投票成功"})
}

type VoteMod struct {
	Openid     string
	CandId     string
	GroupId	 string
	CurVoteNum int
	Key        string
	rc         redis.Conn
}

func CreateVoteMod(CandId string, GroupId string, Openid string) *VoteMod {
	var votemod VoteMod
	votemod.CandId = CandId
	votemod.Openid = Openid
	votemod.GroupId = GroupId
	votemod.CurVoteNum = 0
	return &votemod
}

func (v *VoteMod) CheckData() error {
	_, err := strconv.Atoi(v.CandId)
	return err
}

func (v *VoteMod) CheckVoteStartOrEnd() error {
	now := time.Now().Unix()
	if now < config.System.Start {
		return errors.New("活动未开始")
	} else if now >= config.System.End {
		return errors.New("活动已结束")
	}
	return nil
}

func (v *VoteMod) CheckCanVote() error {
	date := carbon.Now().ToFormatString("md")
	var sb strings.Builder
	sb.WriteString("vote_")
	sb.WriteString(date)
	sb.WriteString("_")
	sb.WriteString(v.GroupId)
	sb.WriteString(":")
	sb.WriteString(v.Openid)
	v.Key = sb.String()

	// ci, err := redis.String(v.rc.Do("hget", "cands:"+v.CandId, "CategoryId"))
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return errors.New("投票失败")
	// }
	// v.CategoryId = ci

	res, err := redis.Int(v.rc.Do("hget", v.Key, v.CandId))
	if err != nil {
		res = 0
	}
	v.CurVoteNum = res + 1
	if v.CurVoteNum > config.System.Maxvote {
		return errors.New("已投10票，明天再来")
	}
	return nil
}

func (v *VoteMod) Update() error {
	var sbd strings.Builder
	sbd.WriteString(`UPDATE candidate_vote SET vote = vote + 1 WHERE cand_id=`)
	sbd.WriteString(v.CandId)
	rows, err := dbpool.Mysqlpool.Query(sbd.String())
	if err != nil {
		return err
	}
	rows.Close()

	var txt = fmt.Sprintf("INSERT INTO vote_log(openid, cand_id, create_time) values ('%v',%v,%v)", v.Openid, v.CandId, time.Now().Unix())
	rows, err = dbpool.Mysqlpool.Query(txt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	rows.Close()
	return nil
}

func (v *VoteMod) SetCache() error {
	_, err := v.rc.Do("hmset", v.Key, v.CandId, v.CurVoteNum)
	if err != nil{
		return errors.New("投票失败")
	}
	alen, err := redis.Int(v.rc.Do("hlen", v.Key))
	if err != nil{
		alen = 0
	}
	if alen > 5{
		_,_ = v.rc.Do("hdel", v.Key, v.CandId)
		return errors.New("同一类型只能投5名")
	}

	return nil
}
