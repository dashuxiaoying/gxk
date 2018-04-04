package main

import (
	"net/http"
	"github.com/mozillazg/request"
	"time"
	"net/url"
	"log"
	"encoding/json"
	"sync"
	"github.com/segmentio/go-prompt"
	"compress/gzip"
	"os"
	"io"
	"strings"
)

const (
	UserAgents           = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.110 Safari/537.36"
	infoUrl              = "http://hnzj.user.ghlearning.com/train/cms/study/find-my-study-info.gson"
	classUrl             = "http://hnzj.user.ghlearning.com/train/class/my-class.gson?myClassId="
	courseurl            = "http://hnzj.user.ghlearning.com/train/cms/my-video/count-timing-nosign.gson"
	getprocess           = "http://hnzj.user.ghlearning.com/train/class/calculate-my-class-learn-speed.gson"
	updateurl            = "http://hnzj.user.ghlearning.com/train/class/update-my-class.gson"
	vcodeurl             = "http://hnzj.user.ghlearning.com/train/login/get-login-vcode.gson"
	loginurl             = "http://hnzj.user.ghlearning.com/train/login/login.gson"
	zyurl                = "http://hnzj.user.ghlearning.com/train/cms/order/is-free-order.gson"
	zyclassId            = "68d0baea-01f2-49ac-b88a-f53b379d0d95"
	mainclassId          = "c7821735-63ab-4df8-9456-6d1c559ca8b3"
	savemyexamrecordsurl = "http://hnzj.user.ghlearning.com/train/cms/paper/save-my-exam-record.gson?isMakeUp=0&examType=1"
	paperlisturl         = "http://hnzj.user.ghlearning.com/train/cms/paper/find-paper-list.gson?pageNo=1&pageSize=10&cmsType=0"
	paperrelisturl       = "http://hnzj.user.ghlearning.com/train/cms/paper/find-paper-list.gson?pageNo=1&pageSize=10&cmsType=2"
	examUrl              = "http://hnzj.user.ghlearning.com/train/cms/paper/submit-paper.gson"
	starturl             = "http://hnzj.user.ghlearning.com/train/cms/paper/start-do-paper-or-test.gson"
	examparams           = `
	[{"answers":[{"optionId": "e7b31dd5-49eb-4d7e-be4c-a84fb6b0b973"}], "paperQuestionId": "2a26985b-b511-4a29-ae8c-fada8a6c2c93"},
{"answers":[{"optionId": "d9847a74-e411-49ed-896e-ea1a5fd7a245"}], "paperQuestionId": "9c01d73f-e6fd-4ab8-a9fc-179f3ad99cac"},
{"answers":[{"optionId": "2b2fdfeb-8014-4fc0-82e0-f4947b2a8f8b"}], "paperQuestionId": "88fe4c6a-035b-4a1c-a274-d9a4c257dca1"},
{"answers":[{"optionId": "d8d672db-ca18-4ffd-9af2-787d97175d5c"}], "paperQuestionId": "cb02209d-4e57-4ecf-b769-0bd34a801e42"},
{"answers":[{"optionId": "3aaa8ace-901e-474e-8c65-7a58747a3755"}], "paperQuestionId": "4b806f28-87ad-4ff9-b256-41a7a0473882"},
{"answers":[{"optionId": "085a4fc6-c43f-42b3-b8bb-c0a16aa85d7d"}], "paperQuestionId": "c2c28b99-4dfd-406e-be1b-2e1e7547e684"},
{"answers":[{"optionId": "900aa629-ec32-4754-95c4-c638440c2691"}], "paperQuestionId": "3141cb88-3a41-4a3c-af8e-b1751c23686f"},
{"answers":[{"optionId": "7f4e3976-59ee-44aa-acee-f500a06e9aa0"}], "paperQuestionId": "209ec057-9841-448a-952a-d5c005050adf"},
{"answers":[{"optionId": "9928b16e-d2f6-491a-a09a-14db7643c689"}], "paperQuestionId": "718806c2-b254-4054-8caf-e777e687560a"},
{"answers":[{"optionId": "114560c1-ffa5-4b11-8393-8efde6628e64"}], "paperQuestionId": "18e3afa8-e2b6-4639-92f3-4407b0f8f919"},
{"answers":[{"optionId": "05abb2ca-5a9e-4ed7-a44f-e43ee17e614a"}], "paperQuestionId": "e123430b-fe61-4fca-90ab-044a70224ef5"},
{"answers":[{"optionId": "5ae8538a-1a10-4ddb-a3db-64bb7df0ce02"}], "paperQuestionId": "272f595b-1c9e-4825-bafa-161b1af5b269"},
{"answers":[{"optionId": "9a80c7ed-683c-4718-943f-97b36e54a158"}], "paperQuestionId": "809cbaa4-a1ce-4cd8-9f94-908542f18f57"},
{"answers":[{"optionId": "2b21d8dd-6b3d-4aaa-a466-b3062a6932f2"}], "paperQuestionId": "13267a15-b0cf-4370-822b-ac330614682d"},
{"answers":[{"optionId": "5883cfcc-ece5-4a01-bb71-92204d194708"}], "paperQuestionId": "db7d4d49-3e9c-4cb1-af0b-558a0c0ddcfb"},
{"answers":[{"optionId": "da8a9187-6a3a-4223-8264-c24cd307144b"}], "paperQuestionId": "e3fe284f-8f65-45a1-af32-408d45f2c540"},
{"answers":[{"optionId": "4657778d-7e40-47f2-9e0b-e15b51f9283a"}], "paperQuestionId": "0964e53f-85a0-42a2-8c79-21182ad68acb"},
{"answers":[{"optionId": "99cb45c3-24a8-4d1e-bdfb-dc4b313a32f2"}], "paperQuestionId": "13542050-b1a1-4f58-ae64-adf51b689527"},
{"answers":[{"optionId": "dc537ad0-e9b3-4d3d-a045-fbdaeb1ce7f1"}], "paperQuestionId": "e72de0af-4405-4928-9b87-6b95c8678b84"},
{"answers":[{"optionId": "49d207c3-c4f2-4abd-9e6e-f0e6f2113a94"}], "paperQuestionId": "1ee2a444-c462-43b9-9328-295197744703"},
{"answers":[{"optionId": "d9c84411-4e9b-4a12-b6d9-6c90c68b5680"}, {"optionId": "298892b6-c247-4f5f-936a-e1f9f5a7cdea"}, {"optionId": "c0fec689-f805-4b6c-b750-daa8f61f629b"}], "paperQuestionId": "c2df87fc-9d2d-4251-8793-db9aef00dce5"},
{"answers":[{"optionId": "0a178c1b-c546-42e3-b405-5450bc7038c5"}, {"optionId": "6bbf0bdc-5cf6-4e15-8b2f-33f7e107a519"}, {"optionId": "1acaafbb-7aa9-4eec-acb0-7b1c10dce3bc"}], "paperQuestionId": "ad042c9a-a4fe-4ef3-8cfc-600f0686b20d"},
{"answers":[{"optionId": "e92a2838-11f0-4fff-9591-290b712fef0a"}, {"optionId": "9c619065-02a4-4e29-b0d9-7299d22c48ab"}, {"optionId": "5576a739-c48f-4416-abc2-48f4c62edf88"}, {"optionId": "cca2c8f5-4bbb-4a1d-b491-1ca897ee3986"}], "paperQuestionId": "843243e1-2c39-4674-9e3b-0369b40bdff0"},
{"answers":[{"optionId": "6a485e23-f542-411f-84ed-486e96ea1247"}, {"optionId": "ceec5215-2f46-4775-ae3b-de730652f8ad"}], "paperQuestionId": "606d2be2-0c99-4440-8a0e-26e55a80b745"},
{"answers":[{"optionId": "85e58505-235d-40bd-893d-55252c8b2ff9"}, {"optionId": "5109db23-d94f-4a5e-b9fb-2230625305d3"}, {"optionId": "32df67af-3c17-4ccc-94e8-ff9228e051be"}, {"optionId": "39c1840c-9260-4884-befa-1966af05c446"}], "paperQuestionId": "e91739da-cb83-497c-8097-4aaf3fc0739a"},
{"answers":[{"optionId": "b223391e-f739-4198-bceb-629051f15162"}, {"optionId": "d920c21c-d7fa-4d1d-b403-90b1008cd243"}, {"optionId": "7a51741c-6dd9-4ea4-b286-4962f0faaccf"}], "paperQuestionId": "af93f72a-4bf5-4993-ab0c-889c53bdb0e0"},
{"answers":[{"optionId": "6b57462c-9cd7-4f21-a477-751dd193e9c5"}, {"optionId": "955a5353-73f8-4b29-9963-57b689a4da53"}], "paperQuestionId": "59513f8b-bfff-412d-a8aa-52ea62543dd4"},
{"answers":[{"optionId": "17960476-f0e4-4bfd-9e26-63906b3c7782"}, {"optionId": "429b5e2d-a669-4df9-8dd4-9a6de141c58e"}, {"optionId": "36f87e86-6068-4fe5-af3a-62d2763f9607"}], "paperQuestionId": "b9f95e9a-ce85-40d7-a069-dc754504a847"},
{"answers":[{"optionId": "81dc1fc3-f51e-469b-935a-85e5a785ad5e"}, {"optionId": "335e3d40-79b2-4ddd-ad56-694007117b25"}, {"optionId": "f17446d4-606d-4b88-84b8-29ff19e8b11e"}], "paperQuestionId": "0ab3161d-747d-4e43-a88d-b9a66c9ad7ce"},
{"answers":[{"optionId": "b19bf120-c437-46c4-b550-260d5671dab1"}, {"optionId": "a74ceed8-2fe8-4444-b7d2-b89f43286cc0"}, {"optionId": "7936b752-85e5-4fac-803d-6e706340c9aa"}, {"optionId": "e47b993e-eac7-4bde-9884-6f6dfd917558"}], "paperQuestionId": "a4a57601-2c3a-4918-8d5b-ac2505efe0b1"}
]
	`
)

var (
	reqcc      *request.Request
	jobChan    chan url.Values
	JSESSIONID string
	wg         *sync.WaitGroup
)

type B struct {
	VideoSpeed  float32 `json:"videoSpeed"`
	CourseSpeed float32 `json:"courseSpeed"`
}

type A struct {
	RespCode  int    `json:"respCode"`
	RespDesc  string `json:"respDesc"`
	Attribute B      `json:"attribute"`
}

type M struct {
	RespCode  int
	RespDesc  string
	Attribute C `json:"attribute"`
}

type C struct {
	Data D `json:"data"`
}

type D struct {
	MyClassId           string `json:"myClassId"`
	MyClassCourseRPList []E    `json:"myClassCourseRPList"`
}

type E struct {
	MyClassCourseId string `json:"myClassCourseId"`
	VideoRPs        []F    `json:"videoRPs"`
}

type F struct {
	MyClassCourseVideoId string `json:"myClassCourseVideoId"`
	VideoSource          string `json:"videoSource"`
}

type WatchInfo struct {
	Playduration int
	Vid          string
	Timestamp    int64
}

type N struct {
	RespCode  int
	RespDesc  string
	Attribute StudyInfo `json:"attribute"`
}
type StudyInfo struct {
	List []map[string]interface{} `json:"list"`
}

type P struct {
	RespCode  int
	RespDesc  string
	Attribute map[string]float32 `json:"attribute"`
}

type Q struct {
	RespCode  int
	RespDesc  string
	Attribute map[string]G `json:"attribute"`
}

type G struct {
	List    []map[string]interface{} `json:"list"`
	PaperId string                   `json:"paperId"`
}

type Y struct {
	RespCode  int
	RespDesc  string
	Attribute map[string]string `json:"attribute"`
}

func main() {
	log.Println("服务启动")

	//登录
	idcode := prompt.StringRequired("输入身份证")
	password := prompt.String("密码[默认后六位]")
	idcode = strings.ToUpper(idcode)
	if len(password) == 0 {
		log.Println("使用默认密码")
		password = idcode[12:]
	}
	ac := new(http.Client)
	reqcc = request.NewRequest(ac)
	reqcc.Headers = map[string]string{
		"User-Agent": UserAgents,
	}

	res, _ := reqcc.Get(vcodeurl)
	defer res.Body.Close()
	img, _ := os.Create("image.png")
	read, _ := gzip.NewReader(res.Body)
	io.Copy(img, read)
	img.Close()
	cookies := res.Cookies()
	for _, c := range cookies {
		if c.Name == "JSESSIONID" {
			JSESSIONID = c.Value
			break
		}
	}
	code := prompt.String("输入验证码:")

	reqcc.Data = map[string]string{
		"account":    strings.TrimSpace(idcode),
		"password":   strings.TrimSpace(password),
		"vcode":      strings.TrimSpace(code),
		"roleType":   "1",
		"isRemember": "0",
	}
	reqcc.Cookies = map[string]string{
		"JSESSIONID": JSESSIONID,
	}
	r, _ := reqcc.Post(loginurl)
	rrr, _ := r.Text()
	iiresp := A{}
	json.Unmarshal([]byte(rrr), &iiresp)
	if iiresp.RespCode == -1 {
		log.Println("##########账号或密码错误，登录失败##########")
		return
	}

	defer func() {
		log.Println("idcode:" + idcode)
	}()

	wg = &sync.WaitGroup{}
	ok := run()
	wg.Wait()
	log.Println("################all learn finish###################")

	if ok == false {
		return
	}

	//考试
	log.Println("#################start exam###############")

	r1, _ := reqcc.Get(paperlisturl)
	r1t, _ := r1.Content()
	r1s := Q{}
	json.Unmarshal(r1t, &r1s)
	papernum := len(r1s.Attribute["data"].List)
	if papernum == 0 {
		log.Println("##############没有找到试卷，尝试在补考中寻找################")
		r1, _ := reqcc.Get(paperrelisturl)
		r1t, _ := r1.Content()
		r1s := Q{}
		json.Unmarshal(r1t, &r1s)
		papernum = len(r1s.Attribute["data"].List)
		if papernum == 0 {
			log.Println("##############补考试卷也没找到，可能已经完成考试，退出################################")
			return
		}
	}
	myExamId := r1s.Attribute["data"].List[0]["myExamId"].(string)
	paperId := r1s.Attribute["data"].PaperId
	examres, _ := reqcc.Get(savemyexamrecordsurl + "&myExamId=" + myExamId)
	exambytes, _ := examres.Content()
	exams := Y{}
	json.Unmarshal(exambytes, &exams)
	myExamRecordId := exams.Attribute["data"]

	reqcc.Get(starturl + "?paperId=" + paperId + "&myExamRecordId=" + myExamRecordId)

	encodeurl := url.QueryEscape(examparams)
	examtrueurl := examUrl + "?myExamId=" + myExamId + "&myExamRecordId=" + myExamRecordId + "&answerList=" + encodeurl

	reqcc.Get(examtrueurl)
	log.Println("################finish exam#####################")
	return

}

func run() bool {
	//课程学习
	iresp, err := reqcc.Get(infoUrl)
	if err != nil {
		log.Println("##########获取课程信息失败###########")
		return false
	}
	itxt, err := iresp.Text()
	if err != nil {
		log.Println("##########获取课程信息失败###########")
		return false
	}

	myClassId := ""

	studyInfo := N{}
	json.Unmarshal([]byte(itxt), &studyInfo)
	list := studyInfo.Attribute.List
	if len(list) == 0 {
		log.Println("#######没有选择任何课程#########")
		return false
	}
	buy := false
	for aa := 0; aa < len(list); aa++ {
		if list[aa]["classId"].(string) == mainclassId {
			buy = true
			break
		}
	}
	if buy == false {
		log.Println("#######必修课未选，请先交费#########")
		return false
	}
	if len(list) == 1 {
		log.Println("########选修课未选择，尝试选课#######")
		reqcc.Data = map[string]string{
			"classId": zyclassId,
		}
		chooseresp, _ := reqcc.Post(zyurl)
		ab, _ := chooseresp.Text()
		log.Println("########选修成功########", ab)
	}
	//再次查询获取最新课程表
	time.Sleep(5 * time.Second)
	iresp, err = reqcc.Get(infoUrl + "?t=" + time.Now().String())
	itxts,_:=iresp.Content()
	studyInfos := N{}
	json.Unmarshal(itxts, &studyInfos)
	lists := studyInfos.Attribute.List

	for x := 0; x < len(lists); x++ {
		myClassId = lists[x]["myClassId"].(string)
		log.Println(myClassId)
		status := lists[x]["status"].(float64)
		if status == float64(2) {
			log.Println(myClassId, "课程已经学完")
			continue
		}

		cresp, err := reqcc.Get(classUrl + myClassId)
		if err != nil {
			log.Println("##########获取班级信息失败###########")
			return false
		}
		ctxt, err := cresp.Text()
		if err != nil {
			log.Println("##########获取班级信息失败###########")
			return false
		}

		cc := M{}
		json.Unmarshal([]byte(ctxt), &cc)
		myClassCourseId := ""
		videoRPs := []F{}
		myClassCourseVideoId := ""
		vid := ""
		watchInfo := WatchInfo{}
		MyClassCourseRPList := cc.Attribute.Data.MyClassCourseRPList
		for i := 0; i < len(MyClassCourseRPList); i++ {
			myClassCourseId = MyClassCourseRPList[i].MyClassCourseId
			videoRPs = MyClassCourseRPList[i].VideoRPs

			for j := 0; j < len(videoRPs); j++ {
				myClassCourseVideoId = videoRPs[j].MyClassCourseVideoId
				vid = videoRPs[j].VideoSource
				u, _ := url.Parse(courseurl)
				q := u.Query()
				q.Set("myClassId", myClassId)
				q.Set("myClassCourseId", myClassCourseId)
				q.Set("myClassCourseVideoId", myClassCourseVideoId)
				watchInfo.Timestamp = time.Now().Unix()
				watchInfo.Playduration = 0
				watchInfo.Vid = vid
				w, _ := json.Marshal(watchInfo)
				q.Set("watchInfo", string(w))
				u.RawQuery = q.Encode()

				wg.Add(1)
				go worker(u, myClassId, myClassCourseId, myClassCourseVideoId, wg)
			}
		}
	}
	return true
}

func worker(u *url.URL, myClassId, myClassCourseId, myClassCourseVideoId string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		resp, err := reqcc.Get(u.String())
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		res, err := resp.Text()
		if err != nil {
			continue
			return
		}
		log.Println(res)
		r := A{}
		json.Unmarshal([]byte(res), &r)
		if r.Attribute.CourseSpeed == 100.0 {
			reqcc.Get(getprocess + "?myClassId=" + myClassId)
			reqcc.Data = map[string]string{
				"myClassId":            myClassId,
				"myClassCourseId":      myClassCourseId,
				"myClassCourseVideoId": myClassCourseVideoId,
			}
			ccc, _ := reqcc.Post(updateurl)
			log.Println(ccc.Text())
			log.Println("#################learn finish#######################")
			break
		}
		respp, _ := reqcc.Get(getprocess + "?myClassId=" + myClassId)
		if err != nil {
			continue
		}
		defer respp.Body.Close()
		res, err = respp.Text()
		if err != nil {
			continue
		}
		log.Println(res)
		dd := P{}
		json.Unmarshal([]byte(res), &dd)
		if dd.Attribute["data"] == 100.0 {
			log.Println("#################learn finish#######################")
			break
		}

		time.Sleep(5 * time.Second)
	}

	return

}
