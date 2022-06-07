package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vote_cj/dbpool"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type Competitor struct {
	CandId      int    `db:"cand_id" json:"cid"`
	GroupId     int    `db:"group_id" json:"gid"`
	Author      string `db:"author" json:"name"`
	Title       string `db:"title" json:"teamer"`
	Info 				string `db:"info" json:"info"`
	Summary     string `db:"summary" json:"content"`
}

func GetCompetitors(c *gin.Context){
	var res []Competitor
	err := dbpool.Mysqlpool.Select(&res, `Select cand_id,group_id,author,title,info,summary from candidate_info`)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(540, gin.H{"msg": "网络错误"})
		c.Abort()
		return
	}
	resstr, _ := json.Marshal(res)
	c.String(http.StatusOK, "var competitors = "+string(resstr))
}

type Comp struct{
	CandId      int    `db:"cand_id" json:"cid"`
	Vote				int 	 `db:"vote" json:"vote"`
}

func GetComp(c *gin.Context){
	rc := dbpool.Redispool.Get()
	defer rc.Close()
	var res []Comp

	rre, err := redis.Bytes(rc.Do("Get", "/cands"))
	if err == nil{
		err = json.Unmarshal(rre, &res)
		if err == nil{
			c.JSON(http.StatusOK, res)
			c.Abort()
			return
		}else{
			log.Println(err)
		}
	}else{
		log.Println(err)
	}
	
	err = dbpool.Mysqlpool.Select(&res, `Select a.cand_id, b.vote from candidate_info a join candidate_vote b on a.cand_id=b.cand_id order by link asc`)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(540, gin.H{"msg": "网络错误"})
		c.Abort()
		return
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		rc.Do("Set", "/cands", string(resJson))
		rc.Do("Expire", "/cands", 5)
	}
	
	c.JSON(http.StatusOK, res)
}