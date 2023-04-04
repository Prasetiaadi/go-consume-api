package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var BASE_URL = "https://jsonplaceholder.typicode.com"

func Index(w http.ResponseWriter, r *http.Request) {
	var posts []PostStruct

	response, err := http.Get(BASE_URL + "/posts")
	if err != nil {
		log.Print(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&posts); err != nil {
		log.Print(err)
	}

	data := map[string]interface{}{
		"posts": posts,
	}

	temp, _ := template.ParseFiles("views/index.html")
	temp.Execute(w, data)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var post PostStruct
	var data map[string]interface{}

	id := r.URL.Query().Get("id")
	if id != "" {
		res, err := http.Get(BASE_URL + "/posts/" + id)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&post); err != nil {
			log.Print(err)
		}

		data = map[string]interface{}{
			"post": post,
		}
	}

	temp, _ := template.ParseFiles("views/create.html")
	temp.Execute(w, data)
}

func Store(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.Form.Get("post_id")

	idInt, _ := strconv.ParseInt(id, 10, 64)

	newPost := PostStruct{
		Id:     idInt,
		Title:  r.Form.Get("post_title"),
		Body:   r.Form.Get("post_body"),
		UserId: 1,
	}

	jsonValue, _ := json.Marshal(newPost)
	buff := bytes.NewBuffer(jsonValue)

	var req *http.Request
	var err error

	if id != "" {
		// update
		fmt.Println("update")
		req, err = http.NewRequest(http.MethodPut, BASE_URL+"/posts/"+id, buff)
		if err != nil {
			log.Print(err)
			fmt.Println(req.Response.StatusCode)
		}
	} else {
		// create
		fmt.Println("create")
		req, err = http.NewRequest(http.MethodPost, BASE_URL+"/posts", buff)
		if err != nil {
			log.Print(err)
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var postResponse PostStruct

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&postResponse); err != nil {
		log.Print(err)
	}

	fmt.Println(res.StatusCode)
	// fmt.Println(postResponse)

	// mengembalikan ke halaman home
	if res.StatusCode == 200 || res.StatusCode == 201 {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	fmt.Println("delete")
	reqDlt, err := http.NewRequest(http.MethodDelete, BASE_URL+"/posts/"+id, nil)
	if err != nil {
		log.Print(err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(reqDlt)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)

	if res.StatusCode == 200 || res.StatusCode == 201 {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}
