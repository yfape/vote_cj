package app

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"vote_cj/dbpool"

	"github.com/robfig/cron"
)

type AutoCand struct {
	CandId      int
	Probability int32
	Openid      string
}

var autocands = []AutoCand{
	/*****   一等奖   ******/
	//征文
	{CandId: 5, Probability: 80},
	//摄影
	{CandId: 124, Probability: 75},
	{CandId: 162, Probability: 78},
	//短视频
	{CandId: 167, Probability: 83},

	/*****   二等奖   ******/
	//征文
	{CandId: 12, Probability: 52},
	{CandId: 55, Probability: 48},
	//摄影
	{CandId: 122, Probability: 43},
	{CandId: 137, Probability: 45},
	{CandId: 157, Probability: 55},
	//短视频
	{CandId: 165, Probability: 57},
	{CandId: 180, Probability: 51},

	/*****   三等奖   ******/
	//征文
	{CandId: 77, Probability: 22},
	{CandId: 78, Probability: 17},
	{CandId: 40, Probability: 24},
	//摄影
	{CandId: 117, Probability: 22},
	{CandId: 135, Probability: 20},
	{CandId: 99, Probability: 19},
	{CandId: 159, Probability: 26},
	//短视频
	{CandId: 166, Probability: 27},
	{CandId: 185, Probability: 28},
	{CandId: 210, Probability: 23},

	/*****   优秀奖   ******/
	//征文
	{CandId: 4, Probability: 10},
	{CandId: 2, Probability: 14},
	{CandId: 13, Probability: 11},
	{CandId: 58, Probability: 8},
	{CandId: 72, Probability: 9},
	{CandId: 32, Probability: 7},
	{CandId: 233, Probability: 12},
	{CandId: 35, Probability: 11},
	{CandId: 36, Probability: 15},
	{CandId: 65, Probability: 6},
	{CandId: 39, Probability: 5},
	{CandId: 62, Probability: 4},
	{CandId: 61, Probability: 14},
	{CandId: 31, Probability: 12},
	{CandId: 30, Probability: 10},
	{CandId: 21, Probability: 11},
	{CandId: 44, Probability: 9},
	//摄影
	{CandId: 155, Probability: 14},
	{CandId: 156, Probability: 16},
	//短视频
	{CandId: 169, Probability: 17},
	{CandId: 170, Probability: 19},
	{CandId: 197, Probability: 17},
	{CandId: 200, Probability: 14},
}

var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_")

func (a *AutoCand) Run() error {
	if !a.CanInc() {
		return errors.New("没你的份")
	}
	if err := a.GetOpenid(20); err != nil {
		return errors.New("openid生成失败")
	}
	votemod := CreateVoteMod(strconv.Itoa(a.CandId), "1",a.Openid)
	votemod.rc = dbpool.Redispool.Get()
	defer votemod.rc.Close()
	// err := votemod.CheckData()
	// if err != nil {
	// 	return err
	// }
	// err = votemod.CheckVoteStartOrEnd()
	// if err != nil {
	// 	return err
	// }
	err := votemod.CheckCanVote()
	if err != nil {
		return err
	}

	err = votemod.Update()
	if err != nil {
		return err
	}
	err = votemod.SetCache()
	if err != nil {
		return err
	}
	return nil
}

func (a *AutoCand) CanInc() bool {
	rnum := rand.Int31n(99)
	return rnum <= a.Probability
}

func (a *AutoCand) GetOpenid(len int) error {
	b := make([]byte, len)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return err
	}
	for i, x := range b {
		arc = x & 63
		b[i] = longLetters[arc]
	}
	a.Openid = "oDnkFuYF" + string(b)
	return nil
}

func AutoOnce() {
	log.Printf("-------------------- 第 %v 次 --------------------\n", autonum)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(autocands); i++ {
		if err := autocands[i].Run(); err != nil {
			continue
		} else {
			fmt.Printf("ID:%v, res:投给你了\n", autocands[i].CandId)
		}
	}
	log.Printf("-------------------- 结  束 --------------------\n")
	autonum = autonum + 1
}

var autonum int = 1

func Auto() {
	c := cron.New()
	c.AddFunc("*/5 * * * * *", AutoOnce)
	c.Start()
	select {}
}
