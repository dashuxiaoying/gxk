package main

import (
	"net/http"
	"github.com/mozillazg/request"
)

func main(){
	bb := new(http.Client)
	reqcc := request.NewRequest(bb)
	reqcc.Cookies=map[string]string{
		"JSESSIONID":"08C2D42163AD77AC9B8BFF711194C570-o1",
	}
	reqcc.Get("http://hnzj.user.ghlearning.com/my/record")


}