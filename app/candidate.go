package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vote_cj/dbpool"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type Candidate struct {
	CandId      int    `db:"cand_id"`
	GroupId     int    `db:"group_id"`
	CategoryId  int    `db:"category_id"`
	Area        string `db:"area"`
	Author      string `db:"author"`
	Corporation string `db:"corporation"`
	Title       string `db:"title"`
	Info 				string `db:"info"`
	Summary     string `db:"summary"`
	Images      string `db:"images"`
	Order       int    `db:"order"`
	Link        string `db:"link"`
	Vote        int    `db:"vote"`
}

type CandidateMini struct {
	CandId     int `db:"cand_id"`
	GroupId    int `db:"group_id"`
	CategoryId int `db:"category_id"`
	Vote       int `db:"vote"`
}

type CandidateList struct {
	CandId     string `db:"cand_id"`
	GroupId    string `db:"group_id"`
	CategoryId string `db:"category_id"`
}

type ResCand struct{
	CandId int `db:"cand_id" json:"cid"`
	GroupId int `db:"group_id" json:"type"`
	Vote int `db:"vote" json:"poll"`
}



func init() {
	rc := dbpool.Redispool.Get()
	defer rc.Close()
	rows, err := dbpool.Mysqlpool.Queryx(`SELECT cand_id,group_id,category_id FROM candidate_info ORDER BY cand_id;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cl CandidateList
		_ = rows.StructScan(&cl)
		_, err = rc.Do("hmset", "cands:"+cl.CandId, "CandId", cl.CandId, "GroupId", cl.GroupId, "CategoryId", cl.CategoryId)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("candidates参照表加载至内存：成功")
}

func GetCands(c *gin.Context) {
	rc := dbpool.Redispool.Get()
	defer rc.Close()
	res := make(map[int]map[int][]CandidateMini)
	strs, err := redis.String(rc.Do("Get", "/cands"))
	if err != nil {
		fmt.Println("缓存失效：" + err.Error())
	} else {
		err = json.Unmarshal([]byte(strs), &res)
		if err != nil {
			fmt.Println("解析错误：" + err.Error())
		} else {
			c.JSON(http.StatusOK, res)
			return
		}
	}
	rows, err := dbpool.Mysqlpool.Queryx(`SELECT a.cand_id,a.group_id,b.vote FROM candidate_info a LEFT JOIN candidate_vote b ON a.cand_id=b.cand_id ORDER BY a.cand_id;`)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(540, gin.H{"msg": "网络错误"})
		c.Abort()
		return
	}
	defer rows.Close()
	for rows.Next() {
		var cm CandidateMini
		_ = rows.StructScan(&cm)
		if _, ok := res[cm.GroupId]; !ok {
			res[cm.GroupId] = make(map[int][]CandidateMini)
		}

		if _, ok := res[cm.GroupId][cm.CategoryId]; !ok {
			res[cm.GroupId][cm.CategoryId] = []CandidateMini{}
		}
		res[cm.GroupId][cm.CategoryId] = append(res[cm.GroupId][cm.CategoryId], cm)
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		rc.Do("Set", "/cands", string(resJson))
		rc.Do("Expire", "/cands", 2)
	}
	c.JSON(http.StatusOK, res)
}

func GetCandsComplete(c *gin.Context) {
	res := make(map[int]map[int][]Candidate)
	rows, err := dbpool.Mysqlpool.Queryx(`SELECT a.*,b.vote FROM candidate_info a 
	LEFT JOIN candidate_vote b ON a.cand_id=b.cand_id ORDER BY a.cand_id;`)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(540, gin.H{"msg": "网络错误"})
		c.Abort()
		return
	}
	defer rows.Close()
	for rows.Next() {
		var c Candidate
		_ = rows.StructScan(&c)
		if _, ok := res[c.GroupId]; !ok {
			res[c.GroupId] = make(map[int][]Candidate)
		}

		if _, ok := res[c.GroupId][c.CategoryId]; !ok {
			res[c.GroupId][c.CategoryId] = []Candidate{}
		}
		res[c.GroupId][c.CategoryId] = append(res[c.GroupId][c.CategoryId], c)
	}

	resstr, _ := json.Marshal(res)

	c.String(http.StatusOK, "var candidates = "+string(resstr))
}
